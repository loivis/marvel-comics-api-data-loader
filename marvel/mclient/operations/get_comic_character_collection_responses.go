// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/loivis/mcapi-loader/marvel/models"
)

// GetComicCharacterCollectionReader is a Reader for the GetComicCharacterCollection structure.
type GetComicCharacterCollectionReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetComicCharacterCollectionReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetComicCharacterCollectionOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetComicCharacterCollectionOK creates a GetComicCharacterCollectionOK with default headers values
func NewGetComicCharacterCollectionOK() *GetComicCharacterCollectionOK {
	return &GetComicCharacterCollectionOK{}
}

/*GetComicCharacterCollectionOK handles this case with default header values.

No response was specified
*/
type GetComicCharacterCollectionOK struct {
	Payload *models.CharacterDataWrapper
}

func (o *GetComicCharacterCollectionOK) Error() string {
	return fmt.Sprintf("[GET /v1/public/comics/{comicId}/characters][%d] getComicCharacterCollectionOK  %+v", 200, o.Payload)
}

func (o *GetComicCharacterCollectionOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.CharacterDataWrapper)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
