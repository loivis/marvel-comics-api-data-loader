// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// CreatorSummary creator summary
// swagger:model CreatorSummary
type CreatorSummary struct {

	// The full name of the creator.
	Name string `json:"name,omitempty"`

	// The path to the individual creator resource.
	ResourceURI string `json:"resourceURI,omitempty"`

	// The role of the creator in the parent entity.
	Role string `json:"role,omitempty"`
}

// Validate validates this creator summary
func (m *CreatorSummary) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *CreatorSummary) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CreatorSummary) UnmarshalBinary(b []byte) error {
	var res CreatorSummary
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
