// Package access provides an API for accessing jump.
// Jump is aware of the roles of hosts and clients in
// establishing ssh connections, and will allocate unique
// connections to new clients. It requires crossbar to trigger
// the SSH host to connect AFTER the client has connected,
// because SSH is a server-speaks-first protocol. Hence access
// does not need to transmit the URI of the unique connection to the host
// because shellbar will do this when the client makes its
// websocket connection. There is no guarantee a host is connected
// at any given time, and if it drops its management channel
// which is connected to the base session_id, then it cannot be
// reached. As crossbar puts a websocket wrapper around the
// already-encrypted TCP/IP, the communication remains encrypted
// end-to-end. For more details on SSH security properties, see
// https://docstore.mik.ua/orelly/networking_2ndEd/ssh/ch03_01.htm
package access

import (
	"context"
	"flag"
	"fmt"
	"regexp"
	"time"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/runtime/security"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/practable/jump/internal/access/models"
	"github.com/practable/jump/internal/access/restapi"
	"github.com/practable/jump/internal/access/restapi/operations"
	"github.com/practable/jump/internal/permission"
	"github.com/practable/jump/internal/ttlcode"
	log "github.com/sirupsen/logrus"
)

//Config represents configuration of the relay & lets configuration be passed as argument to permit testing
type Config struct {

	// Audience must match the host in token
	Audience string

	// ExchangeCode swaps a code for the associated Token
	CodeStore *ttlcode.CodeStore

	// ConnectionType not needed because we already set that in the swagger file - ConnectHandler only handles connect/

	// Listen is the port this service listens on
	Listen int

	// Secret is used to validate tokens
	Secret string

	//Target is the FQDN of the relay instance
	Target string
}

func getPrefixFromPath(path string) string {

	re := regexp.MustCompile(`^\/([\w\%-]*)\/`)

	matches := re.FindStringSubmatch(path)
	if len(matches) < 2 {
		return ""
	}

	// matches[0] = "/{prefix}/"
	// matches[1] = "{prefix}"
	return matches[1]

}

// API starts the API
// Inputs
// @closed - channel will be closed when server shutsdown
// @wg - waitgroup, we must wg.Done() when we are shutdown
// @port - where to listen locally
// @host - external FQDN of the host (for checking against tokens) e.g. https://relay-access.practable.io
// @target - FQDN of the relay instance e.g. wss://relay.practable.io
// @secret- HMAC shared secret which incoming tokens will be signed with
// @cs - pointer to the CodeStore this API shares with the shellbar websocket relay
// @options - for future backwards compatibility (no options currently available)
func API(ctx context.Context, config Config) { // port int, host, secret, target string, cs *ttlcode.CodeStore) {

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	//create new service API
	api := operations.NewAccessAPI(swaggerSpec)
	server := restapi.NewServer(api)

	//parse flags
	flag.Parse()

	// set the port this service will run on
	server.Port = config.Listen

	// set the Authorizer
	api.BearerAuth = validateHeader(config.Secret, config.Audience)

	// set the Handlers
	api.ConnectHandler = operations.ConnectHandlerFunc(func(params operations.ConnectParams, principal interface{}) middleware.Responder {

		token, ok := principal.(*jwt.Token)
		if !ok {
			c := "401"
			m := "token not jwt"
			return operations.NewConnectUnauthorized().WithPayload(&models.Error{Code: &c, Message: &m})
		}

		// save checking for key existence individually by checking all at once
		claims, ok := token.Claims.(*permission.Token)

		if !ok {
			c := "401"
			m := "token claims incorrect type"
			return operations.NewConnectUnauthorized().WithPayload(&models.Error{Code: &c, Message: &m})
		}

		if !permission.HasRequiredClaims(*claims) {
			c := "401"
			m := "token missing required claims"
			return operations.NewConnectUnauthorized().WithPayload(&models.Error{Code: &c, Message: &m})
		}

		if params.HostID == "" {

			c := "401"
			m := "path missing host_id"
			return operations.NewConnectUnauthorized().WithPayload(&models.Error{Code: &c, Message: &m})
		}

		if claims.Topic != params.HostID {

			c := "401"
			m := "host_id does not match topic in token"
			return operations.NewConnectUnauthorized().WithPayload(&models.Error{Code: &c, Message: &m})
		}

		// Now we check the scopes ....
		// If "host" is present, then we connect to the base session
		// If "client" is present, then we connect to a unique sub-session
		//  Scopes are modified to be read, write
		// If both scopes are offered, then the behaviour depends on the routing
		// default to treating as a host
		// unless a ConnectionID present in query e.g.
		// &connection_id=134234234324
		// in which case, distinguishing between host and client is irrelevant

		// make scope map
		sm := make(map[string]bool)

		for _, v := range claims.Scopes {
			sm[v] = true
		}

		// check scopes
		hasClientScope := false
		hasHostScope := false
		hasStatsScope := false

		for _, scope := range claims.Scopes {
			if scope == "host" {
				hasHostScope = true
			}
			if scope == "client" {
				hasClientScope = true
			}
			if scope == "stats" {
				hasStatsScope = true
			}
		}

		if hasStatsScope && params.HostID != "stats" {
			c := "401"
			m := "path not valid for stats scope"
			return operations.NewConnectUnauthorized().WithPayload(&models.Error{Code: &c, Message: &m})
		}

		if !hasStatsScope && params.HostID == "stats" {
			c := "401"
			m := "path not valid without stats scope"
			return operations.NewConnectUnauthorized().WithPayload(&models.Error{Code: &c, Message: &m})
		}

		if !(hasClientScope || hasHostScope || hasStatsScope) {
			c := "401"
			m := "missing client, host or stats scope"
			return operations.NewConnectUnauthorized().WithPayload(&models.Error{Code: &c, Message: &m})
		}

		if hasClientScope && hasHostScope {
			c := "401"
			m := "can only have client or host scope, not both"
			return operations.NewConnectUnauthorized().WithPayload(&models.Error{Code: &c, Message: &m})
		}

		topic := claims.Topic
		topicSalt := ""
		alertHost := false

		if hasClientScope { //need a new unique connection
			topic = topic + "/" + uuid.New().String()
			topicSalt = uuid.New().String()
			alertHost = true
		}

		// set read, write scopes as needed by crossbar
		if hasClientScope || hasHostScope || hasStatsScope {
			sm["read"] = true
			sm["write"] = true // we may want later to be able to request an immediate report
		}

		// convert scopes to array of strings
		ss := []string{}

		for k := range sm {
			ss = append(ss, k)
		}

		// Shellbar will take care of alerting the admin channel of
		// the new connection for protocol timing reasons
		// Because ssh is "server speaks first", we want to bridge
		// to the server only when client already in place and
		// listening. There are no further hits on the access endpoint
		// though - the rest is done via websockets
		// hence no handler is needed for https://{access-host}/connect/{host_id}/{connection_id}

		pt := permission.NewToken(
			config.Target,
			claims.ConnectionType,
			topic,
			ss,
			claims.IssuedAt.Unix(),
			claims.NotBefore.Unix(),
			claims.ExpiresAt.Unix(),
		)

		permission.SetTopicSalt(&pt, topicSalt)
		permission.SetAlertHost(&pt, alertHost)

		code := config.CodeStore.SubmitToken(pt)

		uri := config.Target + "/" + claims.ConnectionType + "/" + topic + "?code=" + code

		return operations.NewConnectOK().WithPayload(
			&operations.ConnectOKBody{
				URI: uri,
			})

	})

	go func() {
		<-ctx.Done()
		if err := server.Shutdown(); err != nil {
			log.Fatalln(err)
		}

	}()

	//serve API
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}

	<-ctx.Done()

}

// ValidateHeader checks the bearer token.
// wrap the secret so we can get it at runtime without using global
func validateHeader(secret, host string) security.TokenAuthentication {

	return func(bearerToken string) (interface{}, error) {
		// For apiKey security syntax see https://swagger.io/docs/specification/2-0/authentication/
		claims := &permission.Token{}

		token, err := jwt.ParseWithClaims(bearerToken, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method was %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil {
			log.WithFields(log.Fields{"error": err, "token": bearerToken}).Info(err.Error())
			return nil, fmt.Errorf("error reading token was %s", err.Error())
		}

		if !token.Valid { //checks iat, nbf, exp
			log.Info("Token invalid")
			return nil, fmt.Errorf("token invalid")
		}

		if cc, ok := token.Claims.(*permission.Token); ok {

			if !cc.RegisteredClaims.VerifyAudience(host, true) {
				log.WithFields(log.Fields{"aud": cc.RegisteredClaims.Audience, "host": host}).Error("aud does not match this host")
				return nil, fmt.Errorf("aud %s does not match this host %s", cc.RegisteredClaims.Audience, host)
			}

		} else {
			log.WithFields(log.Fields{"token": bearerToken, "host": host}).Error("error parsing token")
			return nil, err
		}

		return token, nil
	}
}

// Token returns a signed token
func Token(audience, ct, topic, secret string, scopes []string, iat, nbf, exp int64) (string, error) {

	var claims permission.Token
	claims.IssuedAt = jwt.NewNumericDate(time.Unix(iat, 0))
	claims.NotBefore = jwt.NewNumericDate(time.Unix(nbf, 0))
	claims.ExpiresAt = jwt.NewNumericDate(time.Unix(exp, 0))
	claims.Audience = jwt.ClaimStrings{audience}
	claims.Topic = topic
	claims.ConnectionType = ct // e.g. connect
	claims.Scopes = scopes     // e.g. "host", "client", or "stats"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))

}
