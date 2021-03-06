// Code generated by go-swagger; DO NOT EDIT.

package keys

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewWeaviateKeysMeChildrenGetParams creates a new WeaviateKeysMeChildrenGetParams object
// with the default values initialized.
func NewWeaviateKeysMeChildrenGetParams() *WeaviateKeysMeChildrenGetParams {

	return &WeaviateKeysMeChildrenGetParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewWeaviateKeysMeChildrenGetParamsWithTimeout creates a new WeaviateKeysMeChildrenGetParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewWeaviateKeysMeChildrenGetParamsWithTimeout(timeout time.Duration) *WeaviateKeysMeChildrenGetParams {

	return &WeaviateKeysMeChildrenGetParams{

		timeout: timeout,
	}
}

// NewWeaviateKeysMeChildrenGetParamsWithContext creates a new WeaviateKeysMeChildrenGetParams object
// with the default values initialized, and the ability to set a context for a request
func NewWeaviateKeysMeChildrenGetParamsWithContext(ctx context.Context) *WeaviateKeysMeChildrenGetParams {

	return &WeaviateKeysMeChildrenGetParams{

		Context: ctx,
	}
}

// NewWeaviateKeysMeChildrenGetParamsWithHTTPClient creates a new WeaviateKeysMeChildrenGetParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewWeaviateKeysMeChildrenGetParamsWithHTTPClient(client *http.Client) *WeaviateKeysMeChildrenGetParams {

	return &WeaviateKeysMeChildrenGetParams{
		HTTPClient: client,
	}
}

/*WeaviateKeysMeChildrenGetParams contains all the parameters to send to the API endpoint
for the weaviate keys me children get operation typically these are written to a http.Request
*/
type WeaviateKeysMeChildrenGetParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the weaviate keys me children get params
func (o *WeaviateKeysMeChildrenGetParams) WithTimeout(timeout time.Duration) *WeaviateKeysMeChildrenGetParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the weaviate keys me children get params
func (o *WeaviateKeysMeChildrenGetParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the weaviate keys me children get params
func (o *WeaviateKeysMeChildrenGetParams) WithContext(ctx context.Context) *WeaviateKeysMeChildrenGetParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the weaviate keys me children get params
func (o *WeaviateKeysMeChildrenGetParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the weaviate keys me children get params
func (o *WeaviateKeysMeChildrenGetParams) WithHTTPClient(client *http.Client) *WeaviateKeysMeChildrenGetParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the weaviate keys me children get params
func (o *WeaviateKeysMeChildrenGetParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *WeaviateKeysMeChildrenGetParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
