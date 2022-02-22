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

// NewGetActivityByIDParams creates a new GetActivityByIDParams object
//
// There are no default values defined in the spec.
func NewGetActivityByIDParams() GetActivityByIDParams {

	return GetActivityByIDParams{}
}

// GetActivityByIDParams contains all the bound params for the get activity by ID operation
// typically these are obtained from a http.Request
//
// swagger:parameters getActivityByID
type GetActivityByIDParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: path
	*/
	ActivityID string
	/*True returns all available details, false just description.
	  In: query
	*/
	Details *bool
	/*
	  Required: true
	  In: path
	*/
	PoolID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetActivityByIDParams() beforehand.
func (o *GetActivityByIDParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	rActivityID, rhkActivityID, _ := route.Params.GetOK("activity_id")
	if err := o.bindActivityID(rActivityID, rhkActivityID, route.Formats); err != nil {
		res = append(res, err)
	}

	qDetails, qhkDetails, _ := qs.GetOK("details")
	if err := o.bindDetails(qDetails, qhkDetails, route.Formats); err != nil {
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

// bindActivityID binds and validates parameter ActivityID from path.
func (o *GetActivityByIDParams) bindActivityID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route
	o.ActivityID = raw

	return nil
}

// bindDetails binds and validates parameter Details from query.
func (o *GetActivityByIDParams) bindDetails(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertBool(raw)
	if err != nil {
		return errors.InvalidType("details", "query", "bool", raw)
	}
	o.Details = &value

	return nil
}

// bindPoolID binds and validates parameter PoolID from path.
func (o *GetActivityByIDParams) bindPoolID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route
	o.PoolID = raw

	return nil
}