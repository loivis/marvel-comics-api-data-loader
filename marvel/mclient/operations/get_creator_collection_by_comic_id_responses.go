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

// GetCreatorCollectionByComicIDReader is a Reader for the GetCreatorCollectionByComicID structure.
type GetCreatorCollectionByComicIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetCreatorCollectionByComicIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetCreatorCollectionByComicIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetCreatorCollectionByComicIDOK creates a GetCreatorCollectionByComicIDOK with default headers values
func NewGetCreatorCollectionByComicIDOK() *GetCreatorCollectionByComicIDOK {
	return &GetCreatorCollectionByComicIDOK{}
}

/*GetCreatorCollectionByComicIDOK handles this case with default header values.

No response was specified
*/
type GetCreatorCollectionByComicIDOK struct {
	Payload *models.CreatorDataWrapper
}

func (o *GetCreatorCollectionByComicIDOK) Error() string {
	return fmt.Sprintf("[GET /v1/public/comics/{comicId}/creators][%d] getCreatorCollectionByComicIdOK  %+v", 200, o.Payload)
}

func (o *GetCreatorCollectionByComicIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.CreatorDataWrapper)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
