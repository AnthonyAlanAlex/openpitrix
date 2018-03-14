// Code generated by go-swagger; DO NOT EDIT.

package runtime_env_manager

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"openpitrix.io/openpitrix/test/models"
)

// DeleteRuntimeEnvCredentialReader is a Reader for the DeleteRuntimeEnvCredential structure.
type DeleteRuntimeEnvCredentialReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteRuntimeEnvCredentialReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewDeleteRuntimeEnvCredentialOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewDeleteRuntimeEnvCredentialOK creates a DeleteRuntimeEnvCredentialOK with default headers values
func NewDeleteRuntimeEnvCredentialOK() *DeleteRuntimeEnvCredentialOK {
	return &DeleteRuntimeEnvCredentialOK{}
}

/*DeleteRuntimeEnvCredentialOK handles this case with default header values.

DeleteRuntimeEnvCredentialOK delete runtime env credential o k
*/
type DeleteRuntimeEnvCredentialOK struct {
	Payload *models.OpenpitrixDeleteRuntimeEnvCredentialResponse
}

func (o *DeleteRuntimeEnvCredentialOK) Error() string {
	return fmt.Sprintf("[DELETE /v1/runtime_env_credentials][%d] deleteRuntimeEnvCredentialOK  %+v", 200, o.Payload)
}

func (o *DeleteRuntimeEnvCredentialOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.OpenpitrixDeleteRuntimeEnvCredentialResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}