/* package wsauth provides a websocket client for the practable access-relay protocol

 */

package wsauth

import (
	"net"

	"github.com/gorilla/websocket"
)

func New(url net.URL, token string) (*websocket.Conn, error) {

}
