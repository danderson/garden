package reload

import (
	"context"
	"errors"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Runner struct {
	command []string

	cmd  *exec.Cmd
	done chan struct{}

	mu     sync.Mutex
	stop   bool
	cancel context.CancelFunc
}

func NewRunner(command ...string) (*Runner, error) {
	ret := &Runner{
		command: command,
		done:    make(chan struct{}),
	}
	ctx := ret.resetCancel()

	go ret.run(ctx)
	return ret, nil
}

func (r *Runner) Stop() {
	r.requestStop()
	<-r.done
}

func (r *Runner) Restart() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.cancel != nil {
		r.cancel()
	}
}

func (r *Runner) requestStop() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.cancel != nil {
		r.cancel()
	}
	r.stop = true
}

func (r *Runner) stopRequested() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.stop
}

func (r *Runner) resetCancel() context.Context {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.cancel != nil {
		r.cancel()
	}
	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel
	return ctx
}

func (r *Runner) run(initialCtx context.Context) {
	<-initialCtx.Done()
	log.Printf("Starting: %s", strings.Join(r.command, " "))
	for {
		if err := r.runOnce(); err != nil {
			log.Print(err)
		}
		if r.stopRequested() {
			return
		}
		log.Printf("Restarting: %s", strings.Join(r.command, " "))
	}
}

func (r *Runner) runOnce() error {
	ctx := r.resetCancel()
	if r.cmd != nil {
		r.cmd.Wait()
	}
	r.cmd = exec.CommandContext(ctx, r.command[0], r.command[1:]...)
	r.cmd.Stdout = os.Stdout
	r.cmd.Stderr = os.Stderr
	r.cmd.Cancel = func() error {
		r.cmd.Process.Signal(syscall.SIGTERM)
		return nil
	}
	r.cmd.WaitDelay = 2 * time.Second
	if err := r.cmd.Run(); errors.Is(err, fs.ErrNotExist) {
		log.Printf("Command %s not found, waiting for notification", r.command[0])
		<-ctx.Done()
	} else if errors.Is(err, context.Canceled) {
		return nil
	} else if err != nil {
		return err
	}
	return nil
}
