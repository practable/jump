// Code generated by go-swagger; DO NOT EDIT.

package groups

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"

	"github.com/practable/relay/internal/booking/models"
)

// NewAddPoolsByGroupIDParams creates a new AddPoolsByGroupIDParams object
//
// There are no default values defined in the spec.
func NewAddPoolsByGroupIDParams() AddPoolsByGroupIDParams {

	return AddPoolsByGroupIDParams{}
}

// AddPoolsByGroupIDParams contains all the bound params for the add pools by group ID operation
// typically these are obtained from a http.Request
//
// swagger:parameters addPoolsByGroupID
type AddPoolsByGroupIDParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: path
	*/
	GroupID string
	/*
	  Required: true
	  In: body
	*/
	Pools models.IDList
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewAddPoolsByGroupIDParams() beforehand.
func (o *AddPoolsByGroupIDParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rGroupID, rhkGroupID, _ := route.Params.GetOK("group_id")
	if err := o.bindGroupID(rGroupID, rhkGroupID, route.Formats); err != nil {
		res = append(res, err)
	}

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.IDList
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("pools", "body", ""))
			} else {
				res = append(res, errors.NewParseError("pools", "body", "", err))
			}
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			ctx := validate.WithOperationRequest(context.Background())
			if err := body.ContextValidate(ctx, route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Pools = body
			}
		}
	} else {
		res = append(res, errors.Required("pools", "body", ""))
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindGroupID binds and validates parameter GroupID from path.
func (o *AddPoolsByGroupIDParams) bindGroupID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route
	o.GroupID = raw

	return nil
}
