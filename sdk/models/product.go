// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Product Product defines the structure for an API product
//
// swagger:model Product
type Product struct {

	// description
	// Max Length: 10000
	Description string `json:"description,omitempty"`

	// unique identifier
	// Minimum: 1
	ID int64 `json:"id,omitempty"`

	// name
	// Required: true
	// Max Length: 255
	Name *string `json:"name"`

	// price
	// Required: true
	// Minimum: 0.01
	Price *float32 `json:"price"`

	// SKU - in the field of inventory management, a stock keeping unit is a distinct type of item for sale, purchased, or tracked in inventory,
	// such as a product or service, and all attributes associated with the item type that distinguish it from other item types.
	// For a product, these attributes can include manufacturer, description, material, size, color, packaging, and warranty terms.
	// When a business takes inventory of its stock, it counts the quantity it has of each SKU.
	// SKU can also refer to a unique identifier or code, sometimes represented via a barcode for scanning and tracking, that refers to the particular stock keeping unit.
	// These identifiers are not regulated or standardized.
	// When a company receives items from a vendor, it has a choice of maintaining the vendor's SKU or creating its own.
	//
	// Original source: https://en.wikipedia.org/wiki/Stock_keeping_unit
	// Required: true
	// Pattern: [a-z]+-[a-z]+-[a-z]+
	SKU *string `json:"sku"`
}

// Validate validates this product
func (m *Product) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDescription(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePrice(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSKU(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Product) validateDescription(formats strfmt.Registry) error {
	if swag.IsZero(m.Description) { // not required
		return nil
	}

	if err := validate.MaxLength("description", "body", m.Description, 10000); err != nil {
		return err
	}

	return nil
}

func (m *Product) validateID(formats strfmt.Registry) error {
	if swag.IsZero(m.ID) { // not required
		return nil
	}

	if err := validate.MinimumInt("id", "body", m.ID, 1, false); err != nil {
		return err
	}

	return nil
}

func (m *Product) validateName(formats strfmt.Registry) error {

	if err := validate.Required("name", "body", m.Name); err != nil {
		return err
	}

	if err := validate.MaxLength("name", "body", *m.Name, 255); err != nil {
		return err
	}

	return nil
}

func (m *Product) validatePrice(formats strfmt.Registry) error {

	if err := validate.Required("price", "body", m.Price); err != nil {
		return err
	}

	if err := validate.Minimum("price", "body", float64(*m.Price), 0.01, false); err != nil {
		return err
	}

	return nil
}

func (m *Product) validateSKU(formats strfmt.Registry) error {

	if err := validate.Required("sku", "body", m.SKU); err != nil {
		return err
	}

	if err := validate.Pattern("sku", "body", *m.SKU, `[a-z]+-[a-z]+-[a-z]+`); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this product based on context it is used
func (m *Product) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Product) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Product) UnmarshalBinary(b []byte) error {
	var res Product
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
