// Code generated by go-swagger; DO NOT EDIT.

package graphql

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/creativesoftwarefdn/weaviate/models"
)

// WeaviateGraphqlBatchPostReader is a Reader for the WeaviateGraphqlBatchPost structure.
type WeaviateGraphqlBatchPostReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *WeaviateGraphqlBatchPostReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewWeaviateGraphqlBatchPostOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 401:
		result := NewWeaviateGraphqlBatchPostUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 403:
		result := NewWeaviateGraphqlBatchPostForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 422:
		result := NewWeaviateGraphqlBatchPostUnprocessableEntity()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewWeaviateGraphqlBatchPostOK creates a WeaviateGraphqlBatchPostOK with default headers values
func NewWeaviateGraphqlBatchPostOK() *WeaviateGraphqlBatchPostOK {
	return &WeaviateGraphqlBatchPostOK{}
}

/*WeaviateGraphqlBatchPostOK handles this case with default header values.

Succesful query (with select).
*/
type WeaviateGraphqlBatchPostOK struct {
	Payload *models.GraphQLResponse
}

func (o *WeaviateGraphqlBatchPostOK) Error() string {
	return fmt.Sprintf("[POST /graphql/batch][%d] weaviateGraphqlBatchPostOK  %+v", 200, o.Payload)
}

func (o *WeaviateGraphqlBatchPostOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GraphQLResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewWeaviateGraphqlBatchPostUnauthorized creates a WeaviateGraphqlBatchPostUnauthorized with default headers values
func NewWeaviateGraphqlBatchPostUnauthorized() *WeaviateGraphqlBatchPostUnauthorized {
	return &WeaviateGraphqlBatchPostUnauthorized{}
}

/*WeaviateGraphqlBatchPostUnauthorized handles this case with default header values.

Unauthorized or invalid credentials.
*/
type WeaviateGraphqlBatchPostUnauthorized struct {
}

func (o *WeaviateGraphqlBatchPostUnauthorized) Error() string {
	return fmt.Sprintf("[POST /graphql/batch][%d] weaviateGraphqlBatchPostUnauthorized ", 401)
}

func (o *WeaviateGraphqlBatchPostUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewWeaviateGraphqlBatchPostForbidden creates a WeaviateGraphqlBatchPostForbidden with default headers values
func NewWeaviateGraphqlBatchPostForbidden() *WeaviateGraphqlBatchPostForbidden {
	return &WeaviateGraphqlBatchPostForbidden{}
}

/*WeaviateGraphqlBatchPostForbidden handles this case with default header values.

The used API-key has insufficient permissions.
*/
type WeaviateGraphqlBatchPostForbidden struct {
}

func (o *WeaviateGraphqlBatchPostForbidden) Error() string {
	return fmt.Sprintf("[POST /graphql/batch][%d] weaviateGraphqlBatchPostForbidden ", 403)
}

func (o *WeaviateGraphqlBatchPostForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewWeaviateGraphqlBatchPostUnprocessableEntity creates a WeaviateGraphqlBatchPostUnprocessableEntity with default headers values
func NewWeaviateGraphqlBatchPostUnprocessableEntity() *WeaviateGraphqlBatchPostUnprocessableEntity {
	return &WeaviateGraphqlBatchPostUnprocessableEntity{}
}

/*WeaviateGraphqlBatchPostUnprocessableEntity handles this case with default header values.

Request body contains well-formed (i.e., syntactically correct), but semantically erroneous. Are you sure the class is defined in the configuration file?
*/
type WeaviateGraphqlBatchPostUnprocessableEntity struct {
	Payload *models.ErrorResponse
}

func (o *WeaviateGraphqlBatchPostUnprocessableEntity) Error() string {
	return fmt.Sprintf("[POST /graphql/batch][%d] weaviateGraphqlBatchPostUnprocessableEntity  %+v", 422, o.Payload)
}

func (o *WeaviateGraphqlBatchPostUnprocessableEntity) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
