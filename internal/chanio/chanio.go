package chanio

import (
	"context"
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
)

// ChanIO represents a ChanIO that can io.Copy data to/from channels
type ChanIO struct {
	// Send channel is read by ChanIO and the data sent to Write/WriteFrom
	ToRead    chan []byte
	FromWrite chan []byte
	ctx       context.Context
	ID        string
}

// New returns a pointer to a ChanIO that uses supplied channels
// it is usually the case the channels will already exist, e.g. from a reconws instance
func New(ctx context.Context, toRead, fromWrite chan []byte, id string) *ChanIO {
	log.Tracef("new chanio(%s)", id)
	return &ChanIO{
		ToRead:    toRead,
		FromWrite: fromWrite,
		ctx:       ctx,
		ID:        id,
	}

}

func (c *ChanIO) Write(p []byte) (n int, err error) {

	select {
	case <-c.ctx.Done():
		log.Tracef("chanio(%s).Read done", c.ID)
		return 0, io.EOF
	case c.FromWrite <- p:
		log.WithField("size", len(p)).Tracef("chanio(%s).Write", c.ID) //even the length is encrypted so don't try to interpret packet
		return len(p), nil
	}

}

func (c *ChanIO) Read(p []byte) (n int, err error) {

	select {
	case <-c.ctx.Done():
		log.Tracef("chanio(%s).Read done", c.ID)
		return 0, io.EOF
	case data := <-c.ToRead:
		n := copy(p, data)
		var err error
		if n < len(data) {
			err = fmt.Errorf("supplied byte slice of %d too small for message of size %d", len(p), len(data))
			log.WithField("error", err.Error()).Errorf("chanio(%s).Read buffer too small", c.ID)
		}
		log.WithField("size", n).Tracef("chanio(%s).Read", c.ID)
		return n, err
	}
}
