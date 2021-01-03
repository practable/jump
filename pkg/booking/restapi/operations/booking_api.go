// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/runtime/security"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/timdrysdale/relay/pkg/booking/restapi/operations/groups"
	"github.com/timdrysdale/relay/pkg/booking/restapi/operations/login"
	"github.com/timdrysdale/relay/pkg/booking/restapi/operations/pools"
)

// NewBookingAPI creates a new Booking instance
func NewBookingAPI(spec *loads.Document) *BookingAPI {
	return &BookingAPI{
		handlers:            make(map[string]map[string]http.Handler),
		formats:             strfmt.Default,
		defaultConsumes:     "application/json",
		defaultProduces:     "application/json",
		customConsumers:     make(map[string]runtime.Consumer),
		customProducers:     make(map[string]runtime.Producer),
		PreServerShutdown:   func() {},
		ServerShutdown:      func() {},
		spec:                spec,
		useSwaggerUI:        false,
		ServeError:          errors.ServeError,
		BasicAuthenticator:  security.BasicAuth,
		APIKeyAuthenticator: security.APIKeyAuth,
		BearerAuthenticator: security.BearerAuth,

		JSONConsumer: runtime.JSONConsumer(),

		JSONProducer: runtime.JSONProducer(),

		PoolsAddActivityByPoolIDHandler: pools.AddActivityByPoolIDHandlerFunc(func(params pools.AddActivityByPoolIDParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation pools.AddActivityByPoolID has not yet been implemented")
		}),
		PoolsAddNewPoolHandler: pools.AddNewPoolHandlerFunc(func(params pools.AddNewPoolParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation pools.AddNewPool has not yet been implemented")
		}),
		GroupsGetGroupDescriptionByIDHandler: groups.GetGroupDescriptionByIDHandlerFunc(func(params groups.GetGroupDescriptionByIDParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation groups.GetGroupDescriptionByID has not yet been implemented")
		}),
		GroupsGetGroupIDByNameHandler: groups.GetGroupIDByNameHandlerFunc(func(params groups.GetGroupIDByNameParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation groups.GetGroupIDByName has not yet been implemented")
		}),
		PoolsGetPoolDescriptionByIDHandler: pools.GetPoolDescriptionByIDHandlerFunc(func(params pools.GetPoolDescriptionByIDParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation pools.GetPoolDescriptionByID has not yet been implemented")
		}),
		PoolsGetPoolStatusByIDHandler: pools.GetPoolStatusByIDHandlerFunc(func(params pools.GetPoolStatusByIDParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation pools.GetPoolStatusByID has not yet been implemented")
		}),
		PoolsGetPoolsByGroupIDHandler: pools.GetPoolsByGroupIDHandlerFunc(func(params pools.GetPoolsByGroupIDParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation pools.GetPoolsByGroupID has not yet been implemented")
		}),
		LoginLoginHandler: login.LoginHandlerFunc(func(params login.LoginParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation login.Login has not yet been implemented")
		}),
		PoolsRequestSessionByPoolIDHandler: pools.RequestSessionByPoolIDHandlerFunc(func(params pools.RequestSessionByPoolIDParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation pools.RequestSessionByPoolID has not yet been implemented")
		}),
		PoolsUpdateActivityByIDHandler: pools.UpdateActivityByIDHandlerFunc(func(params pools.UpdateActivityByIDParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation pools.UpdateActivityByID has not yet been implemented")
		}),

		// Applies when the "Authorization" header is set
		BearerAuth: func(token string) (interface{}, error) {
			return nil, errors.NotImplemented("api key auth (Bearer) Authorization from header param [Authorization] has not yet been implemented")
		},
		// default authorizer is authorized meaning no requests are blocked
		APIAuthorizer: security.Authorized(),
	}
}

/*BookingAPI User API for booking experiments */
type BookingAPI struct {
	spec            *loads.Document
	context         *middleware.Context
	handlers        map[string]map[string]http.Handler
	formats         strfmt.Registry
	customConsumers map[string]runtime.Consumer
	customProducers map[string]runtime.Producer
	defaultConsumes string
	defaultProduces string
	Middleware      func(middleware.Builder) http.Handler
	useSwaggerUI    bool

	// BasicAuthenticator generates a runtime.Authenticator from the supplied basic auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	BasicAuthenticator func(security.UserPassAuthentication) runtime.Authenticator
	// APIKeyAuthenticator generates a runtime.Authenticator from the supplied token auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	APIKeyAuthenticator func(string, string, security.TokenAuthentication) runtime.Authenticator
	// BearerAuthenticator generates a runtime.Authenticator from the supplied bearer token auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	BearerAuthenticator func(string, security.ScopedTokenAuthentication) runtime.Authenticator

	// JSONConsumer registers a consumer for the following mime types:
	//   - application/json
	JSONConsumer runtime.Consumer

	// JSONProducer registers a producer for the following mime types:
	//   - application/json
	JSONProducer runtime.Producer

	// BearerAuth registers a function that takes a token and returns a principal
	// it performs authentication based on an api key Authorization provided in the header
	BearerAuth func(string) (interface{}, error)

	// APIAuthorizer provides access control (ACL/RBAC/ABAC) by providing access to the request and authenticated principal
	APIAuthorizer runtime.Authorizer

	// PoolsAddActivityByPoolIDHandler sets the operation handler for the add activity by pool ID operation
	PoolsAddActivityByPoolIDHandler pools.AddActivityByPoolIDHandler
	// PoolsAddNewPoolHandler sets the operation handler for the add new pool operation
	PoolsAddNewPoolHandler pools.AddNewPoolHandler
	// GroupsGetGroupDescriptionByIDHandler sets the operation handler for the get group description by ID operation
	GroupsGetGroupDescriptionByIDHandler groups.GetGroupDescriptionByIDHandler
	// GroupsGetGroupIDByNameHandler sets the operation handler for the get group ID by name operation
	GroupsGetGroupIDByNameHandler groups.GetGroupIDByNameHandler
	// PoolsGetPoolDescriptionByIDHandler sets the operation handler for the get pool description by ID operation
	PoolsGetPoolDescriptionByIDHandler pools.GetPoolDescriptionByIDHandler
	// PoolsGetPoolStatusByIDHandler sets the operation handler for the get pool status by ID operation
	PoolsGetPoolStatusByIDHandler pools.GetPoolStatusByIDHandler
	// PoolsGetPoolsByGroupIDHandler sets the operation handler for the get pools by group ID operation
	PoolsGetPoolsByGroupIDHandler pools.GetPoolsByGroupIDHandler
	// LoginLoginHandler sets the operation handler for the login operation
	LoginLoginHandler login.LoginHandler
	// PoolsRequestSessionByPoolIDHandler sets the operation handler for the request session by pool ID operation
	PoolsRequestSessionByPoolIDHandler pools.RequestSessionByPoolIDHandler
	// PoolsUpdateActivityByIDHandler sets the operation handler for the update activity by ID operation
	PoolsUpdateActivityByIDHandler pools.UpdateActivityByIDHandler
	// ServeError is called when an error is received, there is a default handler
	// but you can set your own with this
	ServeError func(http.ResponseWriter, *http.Request, error)

	// PreServerShutdown is called before the HTTP(S) server is shutdown
	// This allows for custom functions to get executed before the HTTP(S) server stops accepting traffic
	PreServerShutdown func()

	// ServerShutdown is called when the HTTP(S) server is shut down and done
	// handling all active connections and does not accept connections any more
	ServerShutdown func()

	// Custom command line argument groups with their descriptions
	CommandLineOptionsGroups []swag.CommandLineOptionsGroup

	// User defined logger function.
	Logger func(string, ...interface{})
}

// UseRedoc for documentation at /docs
func (o *BookingAPI) UseRedoc() {
	o.useSwaggerUI = false
}

// UseSwaggerUI for documentation at /docs
func (o *BookingAPI) UseSwaggerUI() {
	o.useSwaggerUI = true
}

// SetDefaultProduces sets the default produces media type
func (o *BookingAPI) SetDefaultProduces(mediaType string) {
	o.defaultProduces = mediaType
}

// SetDefaultConsumes returns the default consumes media type
func (o *BookingAPI) SetDefaultConsumes(mediaType string) {
	o.defaultConsumes = mediaType
}

// SetSpec sets a spec that will be served for the clients.
func (o *BookingAPI) SetSpec(spec *loads.Document) {
	o.spec = spec
}

// DefaultProduces returns the default produces media type
func (o *BookingAPI) DefaultProduces() string {
	return o.defaultProduces
}

// DefaultConsumes returns the default consumes media type
func (o *BookingAPI) DefaultConsumes() string {
	return o.defaultConsumes
}

// Formats returns the registered string formats
func (o *BookingAPI) Formats() strfmt.Registry {
	return o.formats
}

// RegisterFormat registers a custom format validator
func (o *BookingAPI) RegisterFormat(name string, format strfmt.Format, validator strfmt.Validator) {
	o.formats.Add(name, format, validator)
}

// Validate validates the registrations in the BookingAPI
func (o *BookingAPI) Validate() error {
	var unregistered []string

	if o.JSONConsumer == nil {
		unregistered = append(unregistered, "JSONConsumer")
	}

	if o.JSONProducer == nil {
		unregistered = append(unregistered, "JSONProducer")
	}

	if o.BearerAuth == nil {
		unregistered = append(unregistered, "AuthorizationAuth")
	}

	if o.PoolsAddActivityByPoolIDHandler == nil {
		unregistered = append(unregistered, "pools.AddActivityByPoolIDHandler")
	}
	if o.PoolsAddNewPoolHandler == nil {
		unregistered = append(unregistered, "pools.AddNewPoolHandler")
	}
	if o.GroupsGetGroupDescriptionByIDHandler == nil {
		unregistered = append(unregistered, "groups.GetGroupDescriptionByIDHandler")
	}
	if o.GroupsGetGroupIDByNameHandler == nil {
		unregistered = append(unregistered, "groups.GetGroupIDByNameHandler")
	}
	if o.PoolsGetPoolDescriptionByIDHandler == nil {
		unregistered = append(unregistered, "pools.GetPoolDescriptionByIDHandler")
	}
	if o.PoolsGetPoolStatusByIDHandler == nil {
		unregistered = append(unregistered, "pools.GetPoolStatusByIDHandler")
	}
	if o.PoolsGetPoolsByGroupIDHandler == nil {
		unregistered = append(unregistered, "pools.GetPoolsByGroupIDHandler")
	}
	if o.LoginLoginHandler == nil {
		unregistered = append(unregistered, "login.LoginHandler")
	}
	if o.PoolsRequestSessionByPoolIDHandler == nil {
		unregistered = append(unregistered, "pools.RequestSessionByPoolIDHandler")
	}
	if o.PoolsUpdateActivityByIDHandler == nil {
		unregistered = append(unregistered, "pools.UpdateActivityByIDHandler")
	}

	if len(unregistered) > 0 {
		return fmt.Errorf("missing registration: %s", strings.Join(unregistered, ", "))
	}

	return nil
}

// ServeErrorFor gets a error handler for a given operation id
func (o *BookingAPI) ServeErrorFor(operationID string) func(http.ResponseWriter, *http.Request, error) {
	return o.ServeError
}

// AuthenticatorsFor gets the authenticators for the specified security schemes
func (o *BookingAPI) AuthenticatorsFor(schemes map[string]spec.SecurityScheme) map[string]runtime.Authenticator {
	result := make(map[string]runtime.Authenticator)
	for name := range schemes {
		switch name {
		case "Bearer":
			scheme := schemes[name]
			result[name] = o.APIKeyAuthenticator(scheme.Name, scheme.In, o.BearerAuth)

		}
	}
	return result
}

// Authorizer returns the registered authorizer
func (o *BookingAPI) Authorizer() runtime.Authorizer {
	return o.APIAuthorizer
}

// ConsumersFor gets the consumers for the specified media types.
// MIME type parameters are ignored here.
func (o *BookingAPI) ConsumersFor(mediaTypes []string) map[string]runtime.Consumer {
	result := make(map[string]runtime.Consumer, len(mediaTypes))
	for _, mt := range mediaTypes {
		switch mt {
		case "application/json":
			result["application/json"] = o.JSONConsumer
		}

		if c, ok := o.customConsumers[mt]; ok {
			result[mt] = c
		}
	}
	return result
}

// ProducersFor gets the producers for the specified media types.
// MIME type parameters are ignored here.
func (o *BookingAPI) ProducersFor(mediaTypes []string) map[string]runtime.Producer {
	result := make(map[string]runtime.Producer, len(mediaTypes))
	for _, mt := range mediaTypes {
		switch mt {
		case "application/json":
			result["application/json"] = o.JSONProducer
		}

		if p, ok := o.customProducers[mt]; ok {
			result[mt] = p
		}
	}
	return result
}

// HandlerFor gets a http.Handler for the provided operation method and path
func (o *BookingAPI) HandlerFor(method, path string) (http.Handler, bool) {
	if o.handlers == nil {
		return nil, false
	}
	um := strings.ToUpper(method)
	if _, ok := o.handlers[um]; !ok {
		return nil, false
	}
	if path == "/" {
		path = ""
	}
	h, ok := o.handlers[um][path]
	return h, ok
}

// Context returns the middleware context for the booking API
func (o *BookingAPI) Context() *middleware.Context {
	if o.context == nil {
		o.context = middleware.NewRoutableContext(o.spec, o, nil)
	}

	return o.context
}

func (o *BookingAPI) initHandlerCache() {
	o.Context() // don't care about the result, just that the initialization happened
	if o.handlers == nil {
		o.handlers = make(map[string]map[string]http.Handler)
	}

	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/pools/{pool_id}/activities"] = pools.NewAddActivityByPoolID(o.context, o.PoolsAddActivityByPoolIDHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/pools"] = pools.NewAddNewPool(o.context, o.PoolsAddNewPoolHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/groups/{group_id}"] = groups.NewGetGroupDescriptionByID(o.context, o.GroupsGetGroupDescriptionByIDHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/groups"] = groups.NewGetGroupIDByName(o.context, o.GroupsGetGroupIDByNameHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/pools/{pool_id}/description"] = pools.NewGetPoolDescriptionByID(o.context, o.PoolsGetPoolDescriptionByIDHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/pools/{pool_id}/status"] = pools.NewGetPoolStatusByID(o.context, o.PoolsGetPoolStatusByIDHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/pools"] = pools.NewGetPoolsByGroupID(o.context, o.PoolsGetPoolsByGroupIDHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/login"] = login.NewLogin(o.context, o.LoginLoginHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/pools/{pool_id}/sessions"] = pools.NewRequestSessionByPoolID(o.context, o.PoolsRequestSessionByPoolIDHandler)
	if o.handlers["PUT"] == nil {
		o.handlers["PUT"] = make(map[string]http.Handler)
	}
	o.handlers["PUT"]["/pools/{pool_id}/activities/{activity_id}"] = pools.NewUpdateActivityByID(o.context, o.PoolsUpdateActivityByIDHandler)
}

// Serve creates a http handler to serve the API over HTTP
// can be used directly in http.ListenAndServe(":8000", api.Serve(nil))
func (o *BookingAPI) Serve(builder middleware.Builder) http.Handler {
	o.Init()

	if o.Middleware != nil {
		return o.Middleware(builder)
	}
	if o.useSwaggerUI {
		return o.context.APIHandlerSwaggerUI(builder)
	}
	return o.context.APIHandler(builder)
}

// Init allows you to just initialize the handler cache, you can then recompose the middleware as you see fit
func (o *BookingAPI) Init() {
	if len(o.handlers) == 0 {
		o.initHandlerCache()
	}
}

// RegisterConsumer allows you to add (or override) a consumer for a media type.
func (o *BookingAPI) RegisterConsumer(mediaType string, consumer runtime.Consumer) {
	o.customConsumers[mediaType] = consumer
}

// RegisterProducer allows you to add (or override) a producer for a media type.
func (o *BookingAPI) RegisterProducer(mediaType string, producer runtime.Producer) {
	o.customProducers[mediaType] = producer
}

// AddMiddlewareFor adds a http middleware to existing handler
func (o *BookingAPI) AddMiddlewareFor(method, path string, builder middleware.Builder) {
	um := strings.ToUpper(method)
	if path == "/" {
		path = ""
	}
	o.Init()
	if h, ok := o.handlers[um][path]; ok {
		o.handlers[method][path] = builder(h)
	}
}
