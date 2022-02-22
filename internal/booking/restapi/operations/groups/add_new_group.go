// Code generated by go-swagger; DO NOT EDIT.

package groups

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// AddNewGroupHandlerFunc turns a function with the right signature into a add new group handler
type AddNewGroupHandlerFunc func(AddNewGroupParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn AddNewGroupHandlerFunc) Handle(params AddNewGroupParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// AddNewGroupHandler interface for that can handle valid add new group params
type AddNewGroupHandler interface {
	Handle(AddNewGroupParams, interface{}) middleware.Responder
}

// NewAddNewGroup creates a new http.Handler for the add new group operation
func NewAddNewGroup(ctx *middleware.Context, handler AddNewGroupHandler) *AddNewGroup {
	return &AddNewGroup{Context: ctx, Handler: handler}
}

/* AddNewGroup swagger:route POST /groups groups addNewGroup

groups

Create new group

*/
type AddNewGroup struct {
	Context *middleware.Context
	Handler AddNewGroupHandler
}

func (o *AddNewGroup) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewAddNewGroupParams()
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