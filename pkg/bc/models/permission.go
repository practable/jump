// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Permission relay permission
//
// Represents fields needed in relay permission tokens
//
// swagger:model permission
type Permission struct {

	// audience
	// Required: true
	Audience *string `json:"audience"`

	// connection type
	// Required: true
	ConnectionType *string `json:"connection_type"`

	// scopes
	// Required: true
	Scopes []string `json:"scopes"`

	// topic
	// Required: true
	Topic *string `json:"topic"`
}

// Validate validates this permission
func (m *Permission) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAudience(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateConnectionType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateScopes(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTopic(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Permission) validateAudience(formats strfmt.Registry) error {

	if err := validate.Required("audience", "body", m.Audience); err != nil {
		return err
	}

	return nil
}

func (m *Permission) validateConnectionType(formats strfmt.Registry) error {

	if err := validate.Required("connection_type", "body", m.ConnectionType); err != nil {
		return err
	}

	return nil
}

func (m *Permission) validateScopes(formats strfmt.Registry) error {

	if err := validate.Required("scopes", "body", m.Scopes); err != nil {
		return err
	}

	return nil
}

func (m *Permission) validateTopic(formats strfmt.Registry) error {

	if err := validate.Required("topic", "body", m.Topic); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Permission) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Permission) UnmarshalBinary(b []byte) error {
	var res Permission
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
