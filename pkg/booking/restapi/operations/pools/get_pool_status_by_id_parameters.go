// Code generated by go-swagger; DO NOT EDIT.

package pools

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NewGetPoolStatusByIDParams creates a new GetPoolStatusByIDParams object
// no default values defined in spec.
func NewGetPoolStatusByIDParams() GetPoolStatusByIDParams {

	return GetPoolStatusByIDParams{}
}

// GetPoolStatusByIDParams contains all the bound params for the get pool status by ID operation
// typically these are obtained from a http.Request
//
// swagger:parameters getPoolStatusByID
type GetPoolStatusByIDParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  In: query
	*/
	Duration *float64
	/*
	  Required: true
	  In: path
	*/
	PoolID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetPoolStatusByIDParams() beforehand.
func (o *GetPoolStatusByIDParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qDuration, qhkDuration, _ := qs.GetOK("duration")
	if err := o.bindDuration(qDuration, qhkDuration, route.Formats); err != nil {
		res = append(res, err)
	}

	rPoolID, rhkPoolID, _ := route.Params.GetOK("pool_id")
	if err := o.bindPoolID(rPoolID, rhkPoolID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindDuration binds and validates parameter Duration from query.
func (o *GetPoolStatusByIDParams) bindDuration(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertFloat64(raw)
	if err != nil {
		return errors.InvalidType("duration", "query", "float64", raw)
	}
	o.Duration = &value

	return nil
}

// bindPoolID binds and validates parameter PoolID from path.
func (o *GetPoolStatusByIDParams) bindPoolID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.PoolID = raw

	return nil
}
