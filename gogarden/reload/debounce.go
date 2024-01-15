package reload

import (
	"log"
	"time"
)

type debouncer struct {
	fireDelay   time.Duration
	repeatDelay time.Duration
	notify      func()

	input chan struct{}
}

func newDebouncer(fireDelay, repeatDelay time.Duration, notify func()) *debouncer {
	ret := &debouncer{
		fireDelay:   fireDelay,
		repeatDelay: repeatDelay,
		notify:      notify,
		input:       make(chan struct{}, 1),
	}
	go ret.run()
	return ret
}

func (d *debouncer) Notify() {
	select {
	case d.input <- struct{}{}:
	default:
	}
}

func (d *debouncer) Close() error {
	close(d.input)
	return nil
}

func (d *debouncer) run() {
	var soonestNotify time.Time
	for {
		_, ok := <-d.input
		if !ok {
			return
		}
		delay := max(d.fireDelay, time.Until(soonestNotify))
		time.Sleep(delay)
		select {
		case _, ok := <-d.input:
			if !ok {
				return
			}
		default:
		}
		soonestNotify = time.Now().Add(d.repeatDelay)
		log.Printf("notify after %s", delay)
		d.notify()
	}
}
