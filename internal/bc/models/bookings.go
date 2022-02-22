// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Bookings details of bookings held
//
// Contains credentials to access currently booked activities and info on max concurrent sessions
//
// swagger:model bookings
type Bookings struct {

	// Array of activities, including credentials, sufficient to permit access to the activities
	// Required: true
	Activities []*Activity `json:"activities"`

	// if true, no new bookings are being accepted
	Locked bool `json:"locked,omitempty"`

	// Maximum concurrent bookings permitted
	// Required: true
	Max *int64 `json:"max"`

	// message of the day, typically to explain any current or upcoming locked periods
	// Example: Preview Access: sessions will be available from 10am - 5pm GMT today
	Msg string `json:"msg,omitempty"`
}

// Validate validates this bookings
func (m *Bookings) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateActivities(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMax(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Bookings) validateActivities(formats strfmt.Registry) error {

	if err := validate.Required("activities", "body", m.Activities); err != nil {
		return err
	}

	for i := 0; i < len(m.Activities); i++ {
		if swag.IsZero(m.Activities[i]) { // not required
			continue
		}

		if m.Activities[i] != nil {
			if err := m.Activities[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("activities" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *Bookings) validateMax(formats strfmt.Registry) error {

	if err := validate.Required("max", "body", m.Max); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this bookings based on the context it is used
func (m *Bookings) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateActivities(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Bookings) contextValidateActivities(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Activities); i++ {

		if m.Activities[i] != nil {
			if err := m.Activities[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("activities" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *Bookings) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Bookings) UnmarshalBinary(b []byte) error {
	var res Bookings
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}