// Code generated by go-swagger; DO NOT EDIT.

package admin

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/practable/relay/internal/bc/models"
)

// NewImportPoolStoreParams creates a new ImportPoolStoreParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewImportPoolStoreParams() *ImportPoolStoreParams {
	return &ImportPoolStoreParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewImportPoolStoreParamsWithTimeout creates a new ImportPoolStoreParams object
// with the ability to set a timeout on a request.
func NewImportPoolStoreParamsWithTimeout(timeout time.Duration) *ImportPoolStoreParams {
	return &ImportPoolStoreParams{
		timeout: timeout,
	}
}

// NewImportPoolStoreParamsWithContext creates a new ImportPoolStoreParams object
// with the ability to set a context for a request.
func NewImportPoolStoreParamsWithContext(ctx context.Context) *ImportPoolStoreParams {
	return &ImportPoolStoreParams{
		Context: ctx,
	}
}

// NewImportPoolStoreParamsWithHTTPClient creates a new ImportPoolStoreParams object
// with the ability to set a custom HTTPClient for a request.
func NewImportPoolStoreParamsWithHTTPClient(client *http.Client) *ImportPoolStoreParams {
	return &ImportPoolStoreParams{
		HTTPClient: client,
	}
}

/* ImportPoolStoreParams contains all the parameters to send to the API endpoint
   for the import pool store operation.

   Typically these are written to a http.Request.
*/
type ImportPoolStoreParams struct {

	// Poolstore.
	Poolstore *models.Poolstore

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the import pool store params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ImportPoolStoreParams) WithDefaults() *ImportPoolStoreParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the import pool store params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ImportPoolStoreParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the import pool store params
func (o *ImportPoolStoreParams) WithTimeout(timeout time.Duration) *ImportPoolStoreParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the import pool store params
func (o *ImportPoolStoreParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the import pool store params
func (o *ImportPoolStoreParams) WithContext(ctx context.Context) *ImportPoolStoreParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the import pool store params
func (o *ImportPoolStoreParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the import pool store params
func (o *ImportPoolStoreParams) WithHTTPClient(client *http.Client) *ImportPoolStoreParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the import pool store params
func (o *ImportPoolStoreParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithPoolstore adds the poolstore to the import pool store params
func (o *ImportPoolStoreParams) WithPoolstore(poolstore *models.Poolstore) *ImportPoolStoreParams {
	o.SetPoolstore(poolstore)
	return o
}

// SetPoolstore adds the poolstore to the import pool store params
func (o *ImportPoolStoreParams) SetPoolstore(poolstore *models.Poolstore) {
	o.Poolstore = poolstore
}

// WriteToRequest writes these params to a swagger request
func (o *ImportPoolStoreParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if o.Poolstore != nil {
		if err := r.SetBodyParam(o.Poolstore); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
