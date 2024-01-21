package views

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	"go.universe.tf/garden/types"
)

var (
	staticSlug string
	staticOnce sync.Once
)

var s = fmt.Sprint
var f = fmt.Sprintf

type signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

func atoi[T signed](i T) string {
	return strconv.FormatInt(int64(i), 10)
}

func static(filename string) string {
	staticOnce.Do(func() {
		var slug [8]byte
		if _, err := io.ReadFull(rand.Reader, slug[:]); err != nil {
			panic("no random bytes available")
		}
		staticSlug = hex.EncodeToString(slug[:])
	})

	return fmt.Sprintf("/static/%s/%s", staticSlug, filename)
}

func i64True(b *int64) bool {
	return b != nil && *b > 0
}

func i64False(b *int64) bool {
	return b != nil && *b == 0
}

func def(s, def string) string {
	if s != "" {
		return s
	}
	return def
}

func date(t time.Time) string {
	return t.Format("2006-01-02")
}

func daysBetween(a, b time.Time) int {
	if a.After(b) {
		a, b = b, a
	}

	// Make sure we're examining both t and current time in the local
	// timezone.
	a = a.In(types.Pacific)
	b = b.In(types.Pacific)
	// Shift both times to midnight of their respective day. Note we
	// don't use Truncate, because that operates on the timestamp not
	// the presentation time, and we want to truncate to midnight
	// local time.
	a = time.Date(a.Year(), a.Month(), a.Day(), 0, 0, 0, 0, types.Pacific)
	b = time.Date(b.Year(), b.Month(), b.Day(), 0, 0, 0, 0, types.Pacific)

	// Now we can reliably figure the number of days since t, where a
	// day is defined as a transition through midnight not
	// 24*time.Hour.
	return int(b.Sub(a) / (24 * time.Hour))
}

func daysAgo(t time.Time) string {
	d := daysBetween(time.Now(), t)

	switch {
	case d == 0:
		return "today"
	case d == 1:
		return "yesterday"
	case d < 365:
		return fmt.Sprintf("%d days ago", d)
	case d < 2*365:
		return "1 year ago"
	default:
		return fmt.Sprintf("%d years ago", d/365)
	}
}

func dateToday() string {
	return time.Now().In(types.Pacific).Format("2006-01-02")
}
