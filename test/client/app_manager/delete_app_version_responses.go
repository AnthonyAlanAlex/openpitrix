// Code generated by go-swagger; DO NOT EDIT.

package app_manager

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"openpitrix.io/openpitrix/test/models"
)

// DeleteAppVersionReader is a Reader for the DeleteAppVersion structure.
type DeleteAppVersionReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteAppVersionReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewDeleteAppVersionOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewDeleteAppVersionOK creates a DeleteAppVersionOK with default headers values
func NewDeleteAppVersionOK() *DeleteAppVersionOK {
	return &DeleteAppVersionOK{}
}

/*DeleteAppVersionOK handles this case with default header values.

DeleteAppVersionOK delete app version o k
*/
type DeleteAppVersionOK struct {
	Payload *models.OpenpitrixDeleteAppVersionResponse
}

func (o *DeleteAppVersionOK) Error() string {
	return fmt.Sprintf("[POST /v1/app_version/action/delete][%d] deleteAppVersionOK  %+v", 200, o.Payload)
}

func (o *DeleteAppVersionOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.OpenpitrixDeleteAppVersionResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
