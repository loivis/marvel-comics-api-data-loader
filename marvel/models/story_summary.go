// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// StorySummary story summary
// swagger:model StorySummary
type StorySummary struct {

	// The canonical name of the story.
	Name string `json:"name,omitempty"`

	// The path to the individual story resource.
	ResourceURI string `json:"resourceURI,omitempty"`

	// The type of the story (interior or cover).
	Type string `json:"type,omitempty"`
}

// Validate validates this story summary
func (m *StorySummary) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *StorySummary) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *StorySummary) UnmarshalBinary(b []byte) error {
	var res StorySummary
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
