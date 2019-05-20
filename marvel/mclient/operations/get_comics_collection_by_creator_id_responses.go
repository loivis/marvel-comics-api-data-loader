// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/loivis/marvel-comics-api-data-loader/marvel/models"
)

// GetComicsCollectionByCreatorIDReader is a Reader for the GetComicsCollectionByCreatorID structure.
type GetComicsCollectionByCreatorIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetComicsCollectionByCreatorIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetComicsCollectionByCreatorIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetComicsCollectionByCreatorIDOK creates a GetComicsCollectionByCreatorIDOK with default headers values
func NewGetComicsCollectionByCreatorIDOK() *GetComicsCollectionByCreatorIDOK {
	return &GetComicsCollectionByCreatorIDOK{}
}

/*GetComicsCollectionByCreatorIDOK handles this case with default header values.

No response was specified
*/
type GetComicsCollectionByCreatorIDOK struct {
	Payload *models.ComicDataWrapper
}

func (o *GetComicsCollectionByCreatorIDOK) Error() string {
	return fmt.Sprintf("[GET /v1/public/creators/{creatorId}/comics][%d] getComicsCollectionByCreatorIdOK  %+v", 200, o.Payload)
}

func (o *GetComicsCollectionByCreatorIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ComicDataWrapper)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
