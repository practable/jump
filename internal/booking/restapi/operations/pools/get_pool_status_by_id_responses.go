// Code generated by go-swagger; DO NOT EDIT.

package pools

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/practable/relay/internal/booking/models"
)

// GetPoolStatusByIDOKCode is the HTTP code returned for type GetPoolStatusByIDOK
const GetPoolStatusByIDOKCode int = 200

/*GetPoolStatusByIDOK get pool status by Id o k

swagger:response getPoolStatusByIdOK
*/
type GetPoolStatusByIDOK struct {

	/*
	  In: Body
	*/
	Payload *models.Status `json:"body,omitempty"`
}

// NewGetPoolStatusByIDOK creates GetPoolStatusByIDOK with default headers values
func NewGetPoolStatusByIDOK() *GetPoolStatusByIDOK {

	return &GetPoolStatusByIDOK{}
}

// WithPayload adds the payload to the get pool status by Id o k response
func (o *GetPoolStatusByIDOK) WithPayload(payload *models.Status) *GetPoolStatusByIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get pool status by Id o k response
func (o *GetPoolStatusByIDOK) SetPayload(payload *models.Status) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetPoolStatusByIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetPoolStatusByIDUnauthorizedCode is the HTTP code returned for type GetPoolStatusByIDUnauthorized
const GetPoolStatusByIDUnauthorizedCode int = 401

/*GetPoolStatusByIDUnauthorized Unauthorized

swagger:response getPoolStatusByIdUnauthorized
*/
type GetPoolStatusByIDUnauthorized struct {

	/*
	  In: Body
	*/
	Payload interface{} `json:"body,omitempty"`
}

// NewGetPoolStatusByIDUnauthorized creates GetPoolStatusByIDUnauthorized with default headers values
func NewGetPoolStatusByIDUnauthorized() *GetPoolStatusByIDUnauthorized {

	return &GetPoolStatusByIDUnauthorized{}
}

// WithPayload adds the payload to the get pool status by Id unauthorized response
func (o *GetPoolStatusByIDUnauthorized) WithPayload(payload interface{}) *GetPoolStatusByIDUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get pool status by Id unauthorized response
func (o *GetPoolStatusByIDUnauthorized) SetPayload(payload interface{}) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetPoolStatusByIDUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// GetPoolStatusByIDInternalServerErrorCode is the HTTP code returned for type GetPoolStatusByIDInternalServerError
const GetPoolStatusByIDInternalServerErrorCode int = 500

/*GetPoolStatusByIDInternalServerError Internal Error

swagger:response getPoolStatusByIdInternalServerError
*/
type GetPoolStatusByIDInternalServerError struct {

	/*
	  In: Body
	*/
	Payload interface{} `json:"body,omitempty"`
}

// NewGetPoolStatusByIDInternalServerError creates GetPoolStatusByIDInternalServerError with default headers values
func NewGetPoolStatusByIDInternalServerError() *GetPoolStatusByIDInternalServerError {

	return &GetPoolStatusByIDInternalServerError{}
}

// WithPayload adds the payload to the get pool status by Id internal server error response
func (o *GetPoolStatusByIDInternalServerError) WithPayload(payload interface{}) *GetPoolStatusByIDInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get pool status by Id internal server error response
func (o *GetPoolStatusByIDInternalServerError) SetPayload(payload interface{}) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetPoolStatusByIDInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}