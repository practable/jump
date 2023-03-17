package chanio

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*

          left                              right

      \------------\                    \------------\
a --> | ToRead     |					| ToRead     | <---- c
      \            \   --- io.Copy ---- \            \
b <-- | From Write |   					| From Write | ----> d
      \------------\					\------------\


*/

// TestIOCopy is not a production use-case, in fact, it is inverted
// but there is less test harness to build this way
// so useful while constructing the ChanIO type
func TestIOCopy(t *testing.T) {

	a := make(chan []byte, 1)
	b := make(chan []byte, 1)
	c := make(chan []byte, 1)
	d := make(chan []byte, 1)

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	left := New(ctx, a, b, "left")
	right := New(ctx, c, d, "right")

	go func() { io.Copy(left, right) }()
	go func() { io.Copy(right, left) }()

	h := []byte(`hello`)

	a <- h

	select {
	case m := <-d:
		assert.Equal(t, m, h)
	case <-time.After(time.Second):
		t.Error("reception timed out")
	}
}
