# pipe

SSH connections are breaking with MAC issue after sending large files, preventing automatic admin such as ansible.

Pipe is intended to test alternatives.

## Results

### io.Copy

a direct io.Copy between read and write appears bombproof

### io.Copy to bi-directional channel pair as an intermediary

This exhibits the same failure mechanisms as our relay. It implies that there is some side-effect to do with the timing of io.Copy channels that is resulting in the issue (assumed to be packet combination)


## Ideas

### websocket intermediary

#### Initial stage

io.Copy from conn (tcp) to conn (websocket), set one side up as websocket client, the other as server

If that works, then extend to include

[conn(tcp)<->conn(ws)]<->[conn(ws)<->channels]<->[channels<->conn(ws)]<->[conn(ws)<->conn(tcp)]

if that works, then theoretically, this should work

[conn(tcp)<->conn(ws)]<->[conn(ws)<->channels]<->[channel-based relay]<->[channels<->conn(ws)]<->[conn(ws)<->conn(tcp)]

### Reconnecting

There is work in googles corp-relay to allow ssh connections to be retained over broken connections.

This looks useful, although complicates the implementation initially.

As an intermediate step consider this approach

(a) management connection on host side continues to use reconws, because it does not actually convey any ssh data so can continue in present form. Thus restarting server etc won't prevent the management connection from being made, but it will break any existing connections.
(b) modify the host ssh connections to use some one-off conn-based io.Copy() approach (if that works)
(c) modify the client ssh connections to use the same



