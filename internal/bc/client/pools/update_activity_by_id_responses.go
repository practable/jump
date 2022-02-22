// Code generated by go-swagger; DO NOT EDIT.

package pools

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/practable/relay/internal/bc/models"
)

// UpdateActivityByIDReader is a Reader for the UpdateActivityByID structure.
type UpdateActivityByIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdateActivityByIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUpdateActivityByIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewUpdateActivityByIDUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewUpdateActivityByIDNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewUpdateActivityByIDInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewUpdateActivityByIDOK creates a UpdateActivityByIDOK with default headers values
func NewUpdateActivityByIDOK() *UpdateActivityByIDOK {
	return &UpdateActivityByIDOK{}
}

/* UpdateActivityByIDOK describes a response with status code 200, with default header values.

OK
*/
type UpdateActivityByIDOK struct {
	Payload *models.ID
}

func (o *UpdateActivityByIDOK) Error() string {
	return fmt.Sprintf("[PUT /pools/{pool_id}/activities/{activity_id}][%d] updateActivityByIdOK  %+v", 200, o.Payload)
}
func (o *UpdateActivityByIDOK) GetPayload() *models.ID {
	return o.Payload
}

func (o *UpdateActivityByIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ID)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateActivityByIDUnauthorized creates a UpdateActivityByIDUnauthorized with default headers values
func NewUpdateActivityByIDUnauthorized() *UpdateActivityByIDUnauthorized {
	return &UpdateActivityByIDUnauthorized{}
}

/* UpdateActivityByIDUnauthorized describes a response with status code 401, with default header values.

Unauthorized
*/
type UpdateActivityByIDUnauthorized struct {
	Payload interface{}
}

func (o *UpdateActivityByIDUnauthorized) Error() string {
	return fmt.Sprintf("[PUT /pools/{pool_id}/activities/{activity_id}][%d] updateActivityByIdUnauthorized  %+v", 401, o.Payload)
}
func (o *UpdateActivityByIDUnauthorized) GetPayload() interface{} {
	return o.Payload
}

func (o *UpdateActivityByIDUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateActivityByIDNotFound creates a UpdateActivityByIDNotFound with default headers values
func NewUpdateActivityByIDNotFound() *UpdateActivityByIDNotFound {
	return &UpdateActivityByIDNotFound{}
}

/* UpdateActivityByIDNotFound describes a response with status code 404, with default header values.

Not Found
*/
type UpdateActivityByIDNotFound struct {
	Payload interface{}
}

func (o *UpdateActivityByIDNotFound) Error() string {
	return fmt.Sprintf("[PUT /pools/{pool_id}/activities/{activity_id}][%d] updateActivityByIdNotFound  %+v", 404, o.Payload)
}
func (o *UpdateActivityByIDNotFound) GetPayload() interface{} {
	return o.Payload
}

func (o *UpdateActivityByIDNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateActivityByIDInternalServerError creates a UpdateActivityByIDInternalServerError with default headers values
func NewUpdateActivityByIDInternalServerError() *UpdateActivityByIDInternalServerError {
	return &UpdateActivityByIDInternalServerError{}
}

/* UpdateActivityByIDInternalServerError describes a response with status code 500, with default header values.

Internal Error
*/
type UpdateActivityByIDInternalServerError struct {
	Payload interface{}
}

func (o *UpdateActivityByIDInternalServerError) Error() string {
	return fmt.Sprintf("[PUT /pools/{pool_id}/activities/{activity_id}][%d] updateActivityByIdInternalServerError  %+v", 500, o.Payload)
}
func (o *UpdateActivityByIDInternalServerError) GetPayload() interface{} {
	return o.Payload
}

func (o *UpdateActivityByIDInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
