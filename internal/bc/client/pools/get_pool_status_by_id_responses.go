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

// GetPoolStatusByIDReader is a Reader for the GetPoolStatusByID structure.
type GetPoolStatusByIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetPoolStatusByIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetPoolStatusByIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewGetPoolStatusByIDUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetPoolStatusByIDInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewGetPoolStatusByIDOK creates a GetPoolStatusByIDOK with default headers values
func NewGetPoolStatusByIDOK() *GetPoolStatusByIDOK {
	return &GetPoolStatusByIDOK{}
}

/* GetPoolStatusByIDOK describes a response with status code 200, with default header values.

GetPoolStatusByIDOK get pool status by Id o k
*/
type GetPoolStatusByIDOK struct {
	Payload *models.Status
}

func (o *GetPoolStatusByIDOK) Error() string {
	return fmt.Sprintf("[GET /pools/{pool_id}/status][%d] getPoolStatusByIdOK  %+v", 200, o.Payload)
}
func (o *GetPoolStatusByIDOK) GetPayload() *models.Status {
	return o.Payload
}

func (o *GetPoolStatusByIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Status)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetPoolStatusByIDUnauthorized creates a GetPoolStatusByIDUnauthorized with default headers values
func NewGetPoolStatusByIDUnauthorized() *GetPoolStatusByIDUnauthorized {
	return &GetPoolStatusByIDUnauthorized{}
}

/* GetPoolStatusByIDUnauthorized describes a response with status code 401, with default header values.

Unauthorized
*/
type GetPoolStatusByIDUnauthorized struct {
	Payload interface{}
}

func (o *GetPoolStatusByIDUnauthorized) Error() string {
	return fmt.Sprintf("[GET /pools/{pool_id}/status][%d] getPoolStatusByIdUnauthorized  %+v", 401, o.Payload)
}
func (o *GetPoolStatusByIDUnauthorized) GetPayload() interface{} {
	return o.Payload
}

func (o *GetPoolStatusByIDUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetPoolStatusByIDInternalServerError creates a GetPoolStatusByIDInternalServerError with default headers values
func NewGetPoolStatusByIDInternalServerError() *GetPoolStatusByIDInternalServerError {
	return &GetPoolStatusByIDInternalServerError{}
}

/* GetPoolStatusByIDInternalServerError describes a response with status code 500, with default header values.

Internal Error
*/
type GetPoolStatusByIDInternalServerError struct {
	Payload interface{}
}

func (o *GetPoolStatusByIDInternalServerError) Error() string {
	return fmt.Sprintf("[GET /pools/{pool_id}/status][%d] getPoolStatusByIdInternalServerError  %+v", 500, o.Payload)
}
func (o *GetPoolStatusByIDInternalServerError) GetPayload() interface{} {
	return o.Payload
}

func (o *GetPoolStatusByIDInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
