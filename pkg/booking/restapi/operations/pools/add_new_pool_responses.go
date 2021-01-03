// Code generated by go-swagger; DO NOT EDIT.

package pools

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/timdrysdale/relay/pkg/booking/models"
)

// AddNewPoolOKCode is the HTTP code returned for type AddNewPoolOK
const AddNewPoolOKCode int = 200

/*AddNewPoolOK add new pool o k

swagger:response addNewPoolOK
*/
type AddNewPoolOK struct {

	/*
	  In: Body
	*/
	Payload *models.ID `json:"body,omitempty"`
}

// NewAddNewPoolOK creates AddNewPoolOK with default headers values
func NewAddNewPoolOK() *AddNewPoolOK {

	return &AddNewPoolOK{}
}

// WithPayload adds the payload to the add new pool o k response
func (o *AddNewPoolOK) WithPayload(payload *models.ID) *AddNewPoolOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add new pool o k response
func (o *AddNewPoolOK) SetPayload(payload *models.ID) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddNewPoolOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}