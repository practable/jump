package chanio

import (
	"context"
	"fmt"
	"io"
)

// ChanIO represents a ChanIO that can io.Copy data to/from channels
type ChanIO struct {
	// Send channel is read by ChanIO and the data sent to Write/WriteFrom
	ToRead    chan []byte
	FromWrite chan []byte
	ctx       context.Context
}

// New returns a pointer to a ChanIO that uses supplied channels
// it is usually the case the channels will already exist, e.g. from a reconws instance
func New(ctx context.Context, toRead, fromWrite chan []byte) *ChanIO {
	return &ChanIO{
		ToRead:    toRead,
		FromWrite: fromWrite,
		ctx:       ctx,
	}
}

func (c *ChanIO) Write(p []byte) (n int, err error) {

	select {
	case <-c.ctx.Done():
		return 0, io.EOF
	case c.FromWrite <- p:
		return len(p), nil
	}

}

func (c *ChanIO) Read(p []byte) (n int, err error) {

	select {
	case <-c.ctx.Done():
		return 0, io.EOF
	case data := <-c.ToRead:
		n := copy(data, p)
		var err error
		if n < len(data) {
			err = fmt.Errorf("supplied byte slice of %d too small for message of size %d", len(p), len(data))
		}
		return n, err
	}
}
