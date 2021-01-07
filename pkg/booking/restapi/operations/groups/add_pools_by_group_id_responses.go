// Code generated by go-swagger; DO NOT EDIT.

package groups

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/timdrysdale/relay/pkg/booking/models"
)

// AddPoolsByGroupIDOKCode is the HTTP code returned for type AddPoolsByGroupIDOK
const AddPoolsByGroupIDOKCode int = 200

/*AddPoolsByGroupIDOK add pools by group Id o k

swagger:response addPoolsByGroupIdOK
*/
type AddPoolsByGroupIDOK struct {

	/*
	  In: Body
	*/
	Payload models.Idlist `json:"body,omitempty"`
}

// NewAddPoolsByGroupIDOK creates AddPoolsByGroupIDOK with default headers values
func NewAddPoolsByGroupIDOK() *AddPoolsByGroupIDOK {

	return &AddPoolsByGroupIDOK{}
}

// WithPayload adds the payload to the add pools by group Id o k response
func (o *AddPoolsByGroupIDOK) WithPayload(payload models.Idlist) *AddPoolsByGroupIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add pools by group Id o k response
func (o *AddPoolsByGroupIDOK) SetPayload(payload models.Idlist) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddPoolsByGroupIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.Idlist{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// AddPoolsByGroupIDUnauthorizedCode is the HTTP code returned for type AddPoolsByGroupIDUnauthorized
const AddPoolsByGroupIDUnauthorizedCode int = 401

/*AddPoolsByGroupIDUnauthorized Unauthorized

swagger:response addPoolsByGroupIdUnauthorized
*/
type AddPoolsByGroupIDUnauthorized struct {

	/*
	  In: Body
	*/
	Payload interface{} `json:"body,omitempty"`
}

// NewAddPoolsByGroupIDUnauthorized creates AddPoolsByGroupIDUnauthorized with default headers values
func NewAddPoolsByGroupIDUnauthorized() *AddPoolsByGroupIDUnauthorized {

	return &AddPoolsByGroupIDUnauthorized{}
}

// WithPayload adds the payload to the add pools by group Id unauthorized response
func (o *AddPoolsByGroupIDUnauthorized) WithPayload(payload interface{}) *AddPoolsByGroupIDUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add pools by group Id unauthorized response
func (o *AddPoolsByGroupIDUnauthorized) SetPayload(payload interface{}) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddPoolsByGroupIDUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// AddPoolsByGroupIDNotFoundCode is the HTTP code returned for type AddPoolsByGroupIDNotFound
const AddPoolsByGroupIDNotFoundCode int = 404

/*AddPoolsByGroupIDNotFound Not Found

swagger:response addPoolsByGroupIdNotFound
*/
type AddPoolsByGroupIDNotFound struct {
}

// NewAddPoolsByGroupIDNotFound creates AddPoolsByGroupIDNotFound with default headers values
func NewAddPoolsByGroupIDNotFound() *AddPoolsByGroupIDNotFound {

	return &AddPoolsByGroupIDNotFound{}
}

// WriteResponse to the client
func (o *AddPoolsByGroupIDNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// AddPoolsByGroupIDInternalServerErrorCode is the HTTP code returned for type AddPoolsByGroupIDInternalServerError
const AddPoolsByGroupIDInternalServerErrorCode int = 500

/*AddPoolsByGroupIDInternalServerError add pools by group Id internal server error

swagger:response addPoolsByGroupIdInternalServerError
*/
type AddPoolsByGroupIDInternalServerError struct {

	/*
	  In: Body
	*/
	Payload interface{} `json:"body,omitempty"`
}

// NewAddPoolsByGroupIDInternalServerError creates AddPoolsByGroupIDInternalServerError with default headers values
func NewAddPoolsByGroupIDInternalServerError() *AddPoolsByGroupIDInternalServerError {

	return &AddPoolsByGroupIDInternalServerError{}
}

// WithPayload adds the payload to the add pools by group Id internal server error response
func (o *AddPoolsByGroupIDInternalServerError) WithPayload(payload interface{}) *AddPoolsByGroupIDInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add pools by group Id internal server error response
func (o *AddPoolsByGroupIDInternalServerError) SetPayload(payload interface{}) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddPoolsByGroupIDInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}
