// Code generated by go-swagger; DO NOT EDIT.

package groups

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/practable/relay/internal/bc/models"
)

// AddNewGroupReader is a Reader for the AddNewGroup structure.
type AddNewGroupReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *AddNewGroupReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewAddNewGroupOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewAddNewGroupUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewAddNewGroupInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewAddNewGroupOK creates a AddNewGroupOK with default headers values
func NewAddNewGroupOK() *AddNewGroupOK {
	return &AddNewGroupOK{}
}

/* AddNewGroupOK describes a response with status code 200, with default header values.

AddNewGroupOK add new group o k
*/
type AddNewGroupOK struct {
	Payload *models.ID
}

func (o *AddNewGroupOK) Error() string {
	return fmt.Sprintf("[POST /groups][%d] addNewGroupOK  %+v", 200, o.Payload)
}
func (o *AddNewGroupOK) GetPayload() *models.ID {
	return o.Payload
}

func (o *AddNewGroupOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ID)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewAddNewGroupUnauthorized creates a AddNewGroupUnauthorized with default headers values
func NewAddNewGroupUnauthorized() *AddNewGroupUnauthorized {
	return &AddNewGroupUnauthorized{}
}

/* AddNewGroupUnauthorized describes a response with status code 401, with default header values.

Unauthorized
*/
type AddNewGroupUnauthorized struct {
	Payload interface{}
}

func (o *AddNewGroupUnauthorized) Error() string {
	return fmt.Sprintf("[POST /groups][%d] addNewGroupUnauthorized  %+v", 401, o.Payload)
}
func (o *AddNewGroupUnauthorized) GetPayload() interface{} {
	return o.Payload
}

func (o *AddNewGroupUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewAddNewGroupInternalServerError creates a AddNewGroupInternalServerError with default headers values
func NewAddNewGroupInternalServerError() *AddNewGroupInternalServerError {
	return &AddNewGroupInternalServerError{}
}

/* AddNewGroupInternalServerError describes a response with status code 500, with default header values.

AddNewGroupInternalServerError add new group internal server error
*/
type AddNewGroupInternalServerError struct {
	Payload interface{}
}

func (o *AddNewGroupInternalServerError) Error() string {
	return fmt.Sprintf("[POST /groups][%d] addNewGroupInternalServerError  %+v", 500, o.Payload)
}
func (o *AddNewGroupInternalServerError) GetPayload() interface{} {
	return o.Payload
}

func (o *AddNewGroupInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
