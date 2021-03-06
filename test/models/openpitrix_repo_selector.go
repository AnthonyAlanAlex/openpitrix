// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// OpenpitrixRepoSelector openpitrix repo selector
// swagger:model openpitrixRepoSelector
type OpenpitrixRepoSelector struct {

	// create time
	CreateTime strfmt.DateTime `json:"create_time,omitempty"`

	// selector key
	SelectorKey string `json:"selector_key,omitempty"`

	// selector value
	SelectorValue string `json:"selector_value,omitempty"`
}

// Validate validates this openpitrix repo selector
func (m *OpenpitrixRepoSelector) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *OpenpitrixRepoSelector) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *OpenpitrixRepoSelector) UnmarshalBinary(b []byte) error {
	var res OpenpitrixRepoSelector
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
