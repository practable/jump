// Code generated by go-swagger; DO NOT EDIT.

package admin

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetStoreStatusHandlerFunc turns a function with the right signature into a get store status handler
type GetStoreStatusHandlerFunc func(GetStoreStatusParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetStoreStatusHandlerFunc) Handle(params GetStoreStatusParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetStoreStatusHandler interface for that can handle valid get store status params
type GetStoreStatusHandler interface {
	Handle(GetStoreStatusParams, interface{}) middleware.Responder
}

// NewGetStoreStatus creates a new http.Handler for the get store status operation
func NewGetStoreStatus(ctx *middleware.Context, handler GetStoreStatusHandler) *GetStoreStatus {
	return &GetStoreStatus{Context: ctx, Handler: handler}
}

/* GetStoreStatus swagger:route GET /admin/status admin getStoreStatus

Get current status

Get the current status (number of groups, pools, bookings, time til last booking finished etc)

*/
type GetStoreStatus struct {
	Context *middleware.Context
	Handler GetStoreStatusHandler
}

func (o *GetStoreStatus) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetStoreStatusParams()
	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal interface{}
	if uprinc != nil {
		principal = uprinc.(interface{}) // this is really a interface{}, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}