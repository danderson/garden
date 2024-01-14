package views

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"sync"
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
