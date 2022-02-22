// Code generated by go-swagger; DO NOT EDIT.

package pools

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/practable/relay/internal/booking/models"
)

// AddActivityByPoolIDOKCode is the HTTP code returned for type AddActivityByPoolIDOK
const AddActivityByPoolIDOKCode int = 200

/*AddActivityByPoolIDOK add activity by pool Id o k

swagger:response addActivityByPoolIdOK
*/
type AddActivityByPoolIDOK struct {

	/*
	  In: Body
	*/
	Payload *models.ID `json:"body,omitempty"`
}

// NewAddActivityByPoolIDOK creates AddActivityByPoolIDOK with default headers values
func NewAddActivityByPoolIDOK() *AddActivityByPoolIDOK {

	return &AddActivityByPoolIDOK{}
}

// WithPayload adds the payload to the add activity by pool Id o k response
func (o *AddActivityByPoolIDOK) WithPayload(payload *models.ID) *AddActivityByPoolIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add activity by pool Id o k response
func (o *AddActivityByPoolIDOK) SetPayload(payload *models.ID) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddActivityByPoolIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// AddActivityByPoolIDUnauthorizedCode is the HTTP code returned for type AddActivityByPoolIDUnauthorized
const AddActivityByPoolIDUnauthorizedCode int = 401

/*AddActivityByPoolIDUnauthorized Unauthorized

swagger:response addActivityByPoolIdUnauthorized
*/
type AddActivityByPoolIDUnauthorized struct {

	/*
	  In: Body
	*/
	Payload interface{} `json:"body,omitempty"`
}

// NewAddActivityByPoolIDUnauthorized creates AddActivityByPoolIDUnauthorized with default headers values
func NewAddActivityByPoolIDUnauthorized() *AddActivityByPoolIDUnauthorized {

	return &AddActivityByPoolIDUnauthorized{}
}

// WithPayload adds the payload to the add activity by pool Id unauthorized response
func (o *AddActivityByPoolIDUnauthorized) WithPayload(payload interface{}) *AddActivityByPoolIDUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add activity by pool Id unauthorized response
func (o *AddActivityByPoolIDUnauthorized) SetPayload(payload interface{}) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddActivityByPoolIDUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// AddActivityByPoolIDNotFoundCode is the HTTP code returned for type AddActivityByPoolIDNotFound
const AddActivityByPoolIDNotFoundCode int = 404

/*AddActivityByPoolIDNotFound Not Available

swagger:response addActivityByPoolIdNotFound
*/
type AddActivityByPoolIDNotFound struct {

	/*
	  In: Body
	*/
	Payload interface{} `json:"body,omitempty"`
}

// NewAddActivityByPoolIDNotFound creates AddActivityByPoolIDNotFound with default headers values
func NewAddActivityByPoolIDNotFound() *AddActivityByPoolIDNotFound {

	return &AddActivityByPoolIDNotFound{}
}

// WithPayload adds the payload to the add activity by pool Id not found response
func (o *AddActivityByPoolIDNotFound) WithPayload(payload interface{}) *AddActivityByPoolIDNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add activity by pool Id not found response
func (o *AddActivityByPoolIDNotFound) SetPayload(payload interface{}) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddActivityByPoolIDNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// AddActivityByPoolIDInternalServerErrorCode is the HTTP code returned for type AddActivityByPoolIDInternalServerError
const AddActivityByPoolIDInternalServerErrorCode int = 500

/*AddActivityByPoolIDInternalServerError Internal Error

swagger:response addActivityByPoolIdInternalServerError
*/
type AddActivityByPoolIDInternalServerError struct {

	/*
	  In: Body
	*/
	Payload interface{} `json:"body,omitempty"`
}

// NewAddActivityByPoolIDInternalServerError creates AddActivityByPoolIDInternalServerError with default headers values
func NewAddActivityByPoolIDInternalServerError() *AddActivityByPoolIDInternalServerError {

	return &AddActivityByPoolIDInternalServerError{}
}

// WithPayload adds the payload to the add activity by pool Id internal server error response
func (o *AddActivityByPoolIDInternalServerError) WithPayload(payload interface{}) *AddActivityByPoolIDInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add activity by pool Id internal server error response
func (o *AddActivityByPoolIDInternalServerError) SetPayload(payload interface{}) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddActivityByPoolIDInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}
