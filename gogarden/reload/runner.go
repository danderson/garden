package reload

import (
	"context"
	"errors"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

type Runner struct {
	command  []string
	ctx      context.Context
	shutdown context.CancelFunc

	cmd *exec.Cmd

	mu     sync.Mutex
	cancel context.CancelFunc
}

func NewRunner(ctx context.Context, command ...string) (*Runner, error) {
	ret := &Runner{
		command: command,
		ctx:     ctx,
	}

	go ret.run()
	return ret, nil
}

func (r *Runner) Restart() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.cancel != nil {
		r.cancel()
	}
}

func (r *Runner) resetCancel() context.Context {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.cancel != nil {
		r.cancel()
	}
	ctx, cancel := context.WithCancel(r.ctx)
	r.cancel = cancel
	return ctx
}

func (r *Runner) run() {
	for {
		if err := r.runOnce(); err != nil {
			log.Print(err)
		}
		if r.ctx.Err() != nil {
			return
		}
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
	log.Print("starting command")
	if err := r.cmd.Run(); errors.Is(err, fs.ErrNotExist) {
		log.Printf("Command %s not found, waiting for notification", r.command[0])
		<-r.ctx.Done()
	} else if err != nil {
		return err
	}
	return nil
}
