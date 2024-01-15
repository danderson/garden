package reload

import (
	"context"
	"log"
	"time"
)

type Reloader struct {
	tasks   []*Task
	watcher *watcher
}

func New(tasks []*Task) (*Reloader, error) {

	ret := &Reloader{
		tasks:   tasks,
		watcher: newWatcher(),
	}
	if err := ret.watcher.Start(); err != nil {
		return nil, err
	}
	go ret.run()

	return ret, nil
}

func (r *Reloader) run() {
	for _, task := range r.tasks {
		task.prepare()
	}

	needsRun := true
	for {
		if r.checkTasks() {
			needsRun = true
		}
		if !needsRun {
			log.Print("waiting for changes")
			r.watcher.Wait(context.Background(), 10*time.Millisecond)
		}
		needsRun = false

		var notify []func()
		for _, task := range r.tasks {
			if task.needsRun {
				if err := task.run(); err != nil {
					log.Printf("%s: %v", task.Name, err)
				}
				if n := task.Notify; n != nil {
					notify = append(notify, n)
				}
			}
			if r.checkTasks() {
				needsRun = true
			}
		}

		for _, n := range notify {
			n()
		}
	}
}

func (r *Reloader) checkTasks() bool {
	changes := r.watcher.Changes()
	if len(changes) == 0 {
		return false
	}

	haveRunnable := false
	for _, task := range r.tasks {
		if task.match(changes) {
			haveRunnable = true
		}
	}
	return haveRunnable
}

func (r *Reloader) runTasks(tasks []Task) {
	var ok []func()
	for _, task := range tasks {
		if err := task.run(); err != nil {
			log.Printf("%s: %v", task.Name, err)
		} else if task.Notify != nil {
			ok = append(ok, task.Notify)
		}
	}

	for _, f := range ok {
		f()
	}
}
