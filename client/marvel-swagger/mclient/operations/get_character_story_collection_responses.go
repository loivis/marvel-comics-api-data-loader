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

// GetCharacterStoryCollectionReader is a Reader for the GetCharacterStoryCollection structure.
type GetCharacterStoryCollectionReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetCharacterStoryCollectionReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetCharacterStoryCollectionOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetCharacterStoryCollectionOK creates a GetCharacterStoryCollectionOK with default headers values
func NewGetCharacterStoryCollectionOK() *GetCharacterStoryCollectionOK {
	return &GetCharacterStoryCollectionOK{}
}

/*GetCharacterStoryCollectionOK handles this case with default header values.

No response was specified
*/
type GetCharacterStoryCollectionOK struct {
	Payload *models.StoryDataWrapper
}

func (o *GetCharacterStoryCollectionOK) Error() string {
	return fmt.Sprintf("[GET /v1/public/characters/{characterId}/stories][%d] getCharacterStoryCollectionOK  %+v", 200, o.Payload)
}

func (o *GetCharacterStoryCollectionOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.StoryDataWrapper)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}