// Code generated by go-swagger; DO NOT EDIT.

package pools

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/practable/relay/internal/booking/models"
)

// GetPoolDescriptionByIDOKCode is the HTTP code returned for type GetPoolDescriptionByIDOK
const GetPoolDescriptionByIDOKCode int = 200

/*GetPoolDescriptionByIDOK get pool description by Id o k

swagger:response getPoolDescriptionByIdOK
*/
type GetPoolDescriptionByIDOK struct {

	/*
	  In: Body
	*/
	Payload *models.Description `json:"body,omitempty"`
}

// NewGetPoolDescriptionByIDOK creates GetPoolDescriptionByIDOK with default headers values
func NewGetPoolDescriptionByIDOK() *GetPoolDescriptionByIDOK {

	return &GetPoolDescriptionByIDOK{}
}

// WithPayload adds the payload to the get pool description by Id o k response
func (o *GetPoolDescriptionByIDOK) WithPayload(payload *models.Description) *GetPoolDescriptionByIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get pool description by Id o k response
func (o *GetPoolDescriptionByIDOK) SetPayload(payload *models.Description) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetPoolDescriptionByIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetPoolDescriptionByIDUnauthorizedCode is the HTTP code returned for type GetPoolDescriptionByIDUnauthorized
const GetPoolDescriptionByIDUnauthorizedCode int = 401

/*GetPoolDescriptionByIDUnauthorized Unauthorized

swagger:response getPoolDescriptionByIdUnauthorized
*/
type GetPoolDescriptionByIDUnauthorized struct {

	/*
	  In: Body
	*/
	Payload interface{} `json:"body,omitempty"`
}

// NewGetPoolDescriptionByIDUnauthorized creates GetPoolDescriptionByIDUnauthorized with default headers values
func NewGetPoolDescriptionByIDUnauthorized() *GetPoolDescriptionByIDUnauthorized {

	return &GetPoolDescriptionByIDUnauthorized{}
}

// WithPayload adds the payload to the get pool description by Id unauthorized response
func (o *GetPoolDescriptionByIDUnauthorized) WithPayload(payload interface{}) *GetPoolDescriptionByIDUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get pool description by Id unauthorized response
func (o *GetPoolDescriptionByIDUnauthorized) SetPayload(payload interface{}) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetPoolDescriptionByIDUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// GetPoolDescriptionByIDInternalServerErrorCode is the HTTP code returned for type GetPoolDescriptionByIDInternalServerError
const GetPoolDescriptionByIDInternalServerErrorCode int = 500

/*GetPoolDescriptionByIDInternalServerError Internal Error

swagger:response getPoolDescriptionByIdInternalServerError
*/
type GetPoolDescriptionByIDInternalServerError struct {

	/*
	  In: Body
	*/
	Payload interface{} `json:"body,omitempty"`
}

// NewGetPoolDescriptionByIDInternalServerError creates GetPoolDescriptionByIDInternalServerError with default headers values
func NewGetPoolDescriptionByIDInternalServerError() *GetPoolDescriptionByIDInternalServerError {

	return &GetPoolDescriptionByIDInternalServerError{}
}

// WithPayload adds the payload to the get pool description by Id internal server error response
func (o *GetPoolDescriptionByIDInternalServerError) WithPayload(payload interface{}) *GetPoolDescriptionByIDInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get pool description by Id internal server error response
func (o *GetPoolDescriptionByIDInternalServerError) SetPayload(payload interface{}) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetPoolDescriptionByIDInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}
