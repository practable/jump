/*
package pipe relays a connection from one port to another.

It is intended for testing how we handle ssh connections
by making a direct pipe from one port to another, isolating
our investigation down to how the data crosses the boundary
into our code and out again, both for the convenience and
partitioning out the influence (if any) of the hub.

*/

package pipe

import "context"

type Config struct {
	Listen int
	Target int
}

type Pipe struct {
	Listen int
	Target int
}

func New(config Config) *Pipe {

	return &Pipe{
		Listen: config.Listen,
		Target: config.Target,
	}

}

func (p *Pipe) Run(ctx context.Context) {
	<-ctx.Done()
	return
}
