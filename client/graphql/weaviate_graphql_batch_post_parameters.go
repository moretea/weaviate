// Code generated by go-swagger; DO NOT EDIT.

package graphql

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

	models "github.com/creativesoftwarefdn/weaviate/models"
)

// NewWeaviateGraphqlBatchPostParams creates a new WeaviateGraphqlBatchPostParams object
// with the default values initialized.
func NewWeaviateGraphqlBatchPostParams() *WeaviateGraphqlBatchPostParams {
	var ()
	return &WeaviateGraphqlBatchPostParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewWeaviateGraphqlBatchPostParamsWithTimeout creates a new WeaviateGraphqlBatchPostParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewWeaviateGraphqlBatchPostParamsWithTimeout(timeout time.Duration) *WeaviateGraphqlBatchPostParams {
	var ()
	return &WeaviateGraphqlBatchPostParams{

		timeout: timeout,
	}
}

// NewWeaviateGraphqlBatchPostParamsWithContext creates a new WeaviateGraphqlBatchPostParams object
// with the default values initialized, and the ability to set a context for a request
func NewWeaviateGraphqlBatchPostParamsWithContext(ctx context.Context) *WeaviateGraphqlBatchPostParams {
	var ()
	return &WeaviateGraphqlBatchPostParams{

		Context: ctx,
	}
}

// NewWeaviateGraphqlBatchPostParamsWithHTTPClient creates a new WeaviateGraphqlBatchPostParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewWeaviateGraphqlBatchPostParamsWithHTTPClient(client *http.Client) *WeaviateGraphqlBatchPostParams {
	var ()
	return &WeaviateGraphqlBatchPostParams{
		HTTPClient: client,
	}
}

/*WeaviateGraphqlBatchPostParams contains all the parameters to send to the API endpoint
for the weaviate graphql batch post operation typically these are written to a http.Request
*/
type WeaviateGraphqlBatchPostParams struct {

	/*Body
	  The GraphQL query request parameters as batch.

	*/
	Body models.GraphQLQueryBatch

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the weaviate graphql batch post params
func (o *WeaviateGraphqlBatchPostParams) WithTimeout(timeout time.Duration) *WeaviateGraphqlBatchPostParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the weaviate graphql batch post params
func (o *WeaviateGraphqlBatchPostParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the weaviate graphql batch post params
func (o *WeaviateGraphqlBatchPostParams) WithContext(ctx context.Context) *WeaviateGraphqlBatchPostParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the weaviate graphql batch post params
func (o *WeaviateGraphqlBatchPostParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the weaviate graphql batch post params
func (o *WeaviateGraphqlBatchPostParams) WithHTTPClient(client *http.Client) *WeaviateGraphqlBatchPostParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the weaviate graphql batch post params
func (o *WeaviateGraphqlBatchPostParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the weaviate graphql batch post params
func (o *WeaviateGraphqlBatchPostParams) WithBody(body models.GraphQLQueryBatch) *WeaviateGraphqlBatchPostParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the weaviate graphql batch post params
func (o *WeaviateGraphqlBatchPostParams) SetBody(body models.GraphQLQueryBatch) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *WeaviateGraphqlBatchPostParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
