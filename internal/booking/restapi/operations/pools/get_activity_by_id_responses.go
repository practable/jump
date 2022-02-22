// Code generated by go-swagger; DO NOT EDIT.

package pools

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/practable/relay/internal/booking/models"
)

// GetActivityByIDOKCode is the HTTP code returned for type GetActivityByIDOK
const GetActivityByIDOKCode int = 200

/*GetActivityByIDOK OK

swagger:response getActivityByIdOK
*/
type GetActivityByIDOK struct {

	/*
	  In: Body
	*/
	Payload *models.Description `json:"body,omitempty"`
}

// NewGetActivityByIDOK creates GetActivityByIDOK with default headers values
func NewGetActivityByIDOK() *GetActivityByIDOK {

	return &GetActivityByIDOK{}
}

// WithPayload adds the payload to the get activity by Id o k response
func (o *GetActivityByIDOK) WithPayload(payload *models.Description) *GetActivityByIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get activity by Id o k response
func (o *GetActivityByIDOK) SetPayload(payload *models.Description) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetActivityByIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetActivityByIDUnauthorizedCode is the HTTP code returned for type GetActivityByIDUnauthorized
const GetActivityByIDUnauthorizedCode int = 401

/*GetActivityByIDUnauthorized Unauthorized

swagger:response getActivityByIdUnauthorized
*/
type GetActivityByIDUnauthorized struct {

	/*
	  In: Body
	*/
	Payload interface{} `json:"body,omitempty"`
}

// NewGetActivityByIDUnauthorized creates GetActivityByIDUnauthorized with default headers values
func NewGetActivityByIDUnauthorized() *GetActivityByIDUnauthorized {

	return &GetActivityByIDUnauthorized{}
}

// WithPayload adds the payload to the get activity by Id unauthorized response
func (o *GetActivityByIDUnauthorized) WithPayload(payload interface{}) *GetActivityByIDUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get activity by Id unauthorized response
func (o *GetActivityByIDUnauthorized) SetPayload(payload interface{}) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetActivityByIDUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// GetActivityByIDNotFoundCode is the HTTP code returned for type GetActivityByIDNotFound
const GetActivityByIDNotFoundCode int = 404

/*GetActivityByIDNotFound Not Found

swagger:response getActivityByIdNotFound
*/
type GetActivityByIDNotFound struct {

	/*
	  In: Body
	*/
	Payload interface{} `json:"body,omitempty"`
}

// NewGetActivityByIDNotFound creates GetActivityByIDNotFound with default headers values
func NewGetActivityByIDNotFound() *GetActivityByIDNotFound {

	return &GetActivityByIDNotFound{}
}

// WithPayload adds the payload to the get activity by Id not found response
func (o *GetActivityByIDNotFound) WithPayload(payload interface{}) *GetActivityByIDNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get activity by Id not found response
func (o *GetActivityByIDNotFound) SetPayload(payload interface{}) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetActivityByIDNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// GetActivityByIDInternalServerErrorCode is the HTTP code returned for type GetActivityByIDInternalServerError
const GetActivityByIDInternalServerErrorCode int = 500

/*GetActivityByIDInternalServerError Internal Error

swagger:response getActivityByIdInternalServerError
*/
type GetActivityByIDInternalServerError struct {

	/*
	  In: Body
	*/
	Payload interface{} `json:"body,omitempty"`
}

// NewGetActivityByIDInternalServerError creates GetActivityByIDInternalServerError with default headers values
func NewGetActivityByIDInternalServerError() *GetActivityByIDInternalServerError {

	return &GetActivityByIDInternalServerError{}
}

// WithPayload adds the payload to the get activity by Id internal server error response
func (o *GetActivityByIDInternalServerError) WithPayload(payload interface{}) *GetActivityByIDInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get activity by Id internal server error response
func (o *GetActivityByIDInternalServerError) SetPayload(payload interface{}) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetActivityByIDInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}