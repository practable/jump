/*
package pipe relays a connection from one port to another.

It is intended for testing how we handle ssh connections
by making a direct pipe from one port to another, isolating
our investigation down to how the data crosses the boundary
into our code and out again, both for the convenience and
partitioning out the influence (if any) of the hub.

*/

package pipe

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/practable/jump/internal/tcpconnect"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Listen int
	Target int
}

type Pipe struct {
	Left   chan []byte
	Listen int
	Right  chan []byte
	Target int
}

func New(config Config) *Pipe {

	return &Pipe{
		Left:   make(chan []byte, 1),
		Listen: config.Listen,
		Right:  make(chan []byte, 1),
		Target: config.Target,
	}

}

/*

          |------|                   |--------|
          \      \  ----RIGHT--->    \        \
 Target   | host |                   | client |  Listen
 Port     \      \                   \        \  Port
          |------|  <---LEFT-----    |--------|
*/

func (p *Pipe) Run(ctx context.Context) {

	url := ":" + strconv.Itoa(p.Listen)

	listener := tcpconnect.New()

	go listener.Listen(ctx, url, listenHandler(p.Left, p.Right))
	go targetHandler(ctx, p.Target, p.Left, p.Right)

	<-ctx.Done()

}

func (p *Pipe) RunPacket(ctx context.Context) {

	listenURL := ":" + strconv.Itoa(p.Listen)
	targetURL := "localhost:" + strconv.Itoa(p.Target)

	go packetListen(ctx, listenURL, p.Left, p.Right)
	go packetTarget(ctx, targetURL, p.Left, p.Right)

	<-ctx.Done()

}

/* HANDLERS FOR PACKET METHOD */

func packetListen(ctx context.Context, URL string, left, right chan []byte) {
	log.Error("packetListen not implemented")
	<-ctx.Done()
}

func packetTarget(ctx context.Context, URL string, left, right chan []byte) {
	log.Error("packetTarget not implemented")
	<-ctx.Done()
}

/*  HANDLERS FOR ORIGINAL METHOD  */

// listenHandler handles the external client
// it receives data from the external connection and forwards it internally over channel left
// it accepts internal data from channel right and sends to the external connection
func listenHandler(left, right chan []byte) func(context.Context, *tcpconnect.TCPconnect) {

	return func(ctx context.Context, c *tcpconnect.TCPconnect) {

		timeout := time.Second

		id := "listenHandler(" + c.ID + ")"

		go c.HandleConn(ctx, *c.Conn) // messages on c.In and c.Out

		go func() {

			for {
				select {
				case <-ctx.Done():
					log.Infof("%s: context cancelled; done", id)
					return
				case data, ok := <-c.In:
					if !ok {
						log.Debugf("%s: local channel error, closing local read pump", id)
						return
					}
					size := len(data)
					select {
					case left <- data:
						log.WithField("size", size).Debugf("%s: sent %d-bytes left", id, size)
					case <-time.After(timeout):
						log.WithField("size", size).Debugf("%s: timeout waiting to send %d-bytes left", id, size)
					}

				}

			}

		}()

		go func() {

			for {
				select {
				case <-ctx.Done():
					log.Infof("%s: context cancelled; done", id)
					return
				case data, ok := <-right:
					if !ok {
						log.Debugf("%s: relay channel error, closing relay read pump", id)
						return
					}
					size := len(data)
					select {
					case c.Out <- data:
						log.WithField("size", size).Debugf("%s: sent %d-bytes right", id, size)
					case <-time.After(timeout):
						log.WithField("size", size).Debugf("%s: timeout waiting to send %d-bytes right", id, size)
					}
				}
			}

		}()

		<-ctx.Done()
		log.Infof("%s: done", id)
	}
}

// targetHandler handles the external target, typically port 22 to the localhost
// it receives data from the target and forwards it internally on channel right
// it accepts data internally from channel left and sends it externally to the target
func targetHandler(ctx context.Context, target int, left, right chan []byte) {

	timeout := 1 * time.Second

	id := "targetHandler(" + uuid.New().String()[0:6] + ")"

	log.WithFields(log.Fields{"target": target}).Infof("%s: connecting to target", id)

	t := tcpconnect.New()

	url := "localhost:" + strconv.Itoa(target)

	go t.Dial(ctx, url)

	log.WithFields(log.Fields{"target": target}).Tracef("%s: dialling target", id)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case data, ok := <-t.In:
				log.WithFields(log.Fields{"target": target, "size": len(data)}).Tracef("%s: GOT %d-byte message from target", id, len(data))
				if !ok {
					return
				}
				select {
				case right <- data:
					log.WithFields(log.Fields{"target": target, "size": len(data)}).Tracef("%s: SENT %d-byte message to the right", id, len(data))
				case <-time.After(timeout):
					log.Error("timeout sending message ")
					return
				}
			}

		}
	}()

	go func() {
		for {
			select {

			case <-ctx.Done():
				return
			case data, ok := <-left:
				log.WithFields(log.Fields{"target": target, "size": len(data)}).Tracef("%s: GOT %d-byte message from the left", id, len(data))
				if !ok {
					return
				}
				select {
				case t.Out <- data:
					log.WithFields(log.Fields{"target": target, "size": len(data)}).Tracef("%s: SENT %d-byte Message to target", id, len(data))
				case <-time.After(timeout):

					log.Error("timeout sending message ")
					return
				}
			}

		}
	}()

	<-ctx.Done()
	log.WithFields(log.Fields{"target": target}).Infof("%s: DONE", id)
}
