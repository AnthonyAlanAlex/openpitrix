// Code generated by go-swagger; DO NOT EDIT.

package runtime_env_manager

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewDescribeRuntimeEnvsParams creates a new DescribeRuntimeEnvsParams object
// with the default values initialized.
func NewDescribeRuntimeEnvsParams() *DescribeRuntimeEnvsParams {
	var ()
	return &DescribeRuntimeEnvsParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewDescribeRuntimeEnvsParamsWithTimeout creates a new DescribeRuntimeEnvsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewDescribeRuntimeEnvsParamsWithTimeout(timeout time.Duration) *DescribeRuntimeEnvsParams {
	var ()
	return &DescribeRuntimeEnvsParams{

		timeout: timeout,
	}
}

// NewDescribeRuntimeEnvsParamsWithContext creates a new DescribeRuntimeEnvsParams object
// with the default values initialized, and the ability to set a context for a request
func NewDescribeRuntimeEnvsParamsWithContext(ctx context.Context) *DescribeRuntimeEnvsParams {
	var ()
	return &DescribeRuntimeEnvsParams{

		Context: ctx,
	}
}

// NewDescribeRuntimeEnvsParamsWithHTTPClient creates a new DescribeRuntimeEnvsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewDescribeRuntimeEnvsParamsWithHTTPClient(client *http.Client) *DescribeRuntimeEnvsParams {
	var ()
	return &DescribeRuntimeEnvsParams{
		HTTPClient: client,
	}
}

/*DescribeRuntimeEnvsParams contains all the parameters to send to the API endpoint
for the describe runtime envs operation typically these are written to a http.Request
*/
type DescribeRuntimeEnvsParams struct {

	/*LimitValue
	  The uint32 value.

	*/
	LimitValue *int64
	/*OffsetValue
	  The uint32 value.

	*/
	OffsetValue *int64
	/*Owner*/
	Owner []string
	/*RuntimeEnvID*/
	RuntimeEnvID []string
	/*SearchWord*/
	SearchWord *string
	/*Selector*/
	Selector *string
	/*Status*/
	Status []string
	/*VerboseValue
	  The uint32 value.

	*/
	VerboseValue *int64

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the describe runtime envs params
func (o *DescribeRuntimeEnvsParams) WithTimeout(timeout time.Duration) *DescribeRuntimeEnvsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the describe runtime envs params
func (o *DescribeRuntimeEnvsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the describe runtime envs params
func (o *DescribeRuntimeEnvsParams) WithContext(ctx context.Context) *DescribeRuntimeEnvsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the describe runtime envs params
func (o *DescribeRuntimeEnvsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the describe runtime envs params
func (o *DescribeRuntimeEnvsParams) WithHTTPClient(client *http.Client) *DescribeRuntimeEnvsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the describe runtime envs params
func (o *DescribeRuntimeEnvsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithLimitValue adds the limitValue to the describe runtime envs params
func (o *DescribeRuntimeEnvsParams) WithLimitValue(limitValue *int64) *DescribeRuntimeEnvsParams {
	o.SetLimitValue(limitValue)
	return o
}

// SetLimitValue adds the limitValue to the describe runtime envs params
func (o *DescribeRuntimeEnvsParams) SetLimitValue(limitValue *int64) {
	o.LimitValue = limitValue
}

// WithOffsetValue adds the offsetValue to the describe runtime envs params
func (o *DescribeRuntimeEnvsParams) WithOffsetValue(offsetValue *int64) *DescribeRuntimeEnvsParams {
	o.SetOffsetValue(offsetValue)
	return o
}

// SetOffsetValue adds the offsetValue to the describe runtime envs params
func (o *DescribeRuntimeEnvsParams) SetOffsetValue(offsetValue *int64) {
	o.OffsetValue = offsetValue
}

// WithOwner adds the owner to the describe runtime envs params
func (o *DescribeRuntimeEnvsParams) WithOwner(owner []string) *DescribeRuntimeEnvsParams {
	o.SetOwner(owner)
	return o
}

// SetOwner adds the owner to the describe runtime envs params
func (o *DescribeRuntimeEnvsParams) SetOwner(owner []string) {
	o.Owner = owner
}

// WithRuntimeEnvID adds the runtimeEnvID to the describe runtime envs params
func (o *DescribeRuntimeEnvsParams) WithRuntimeEnvID(runtimeEnvID []string) *DescribeRuntimeEnvsParams {
	o.SetRuntimeEnvID(runtimeEnvID)
	return o
}

// SetRuntimeEnvID adds the runtimeEnvId to the describe runtime envs params
func (o *DescribeRuntimeEnvsParams) SetRuntimeEnvID(runtimeEnvID []string) {
	o.RuntimeEnvID = runtimeEnvID
}

// WithSearchWord adds the searchWord to the describe runtime envs params
func (o *DescribeRuntimeEnvsParams) WithSearchWord(searchWord *string) *DescribeRuntimeEnvsParams {
	o.SetSearchWord(searchWord)
	return o
}

// SetSearchWord adds the searchWord to the describe runtime envs params
func (o *DescribeRuntimeEnvsParams) SetSearchWord(searchWord *string) {
	o.SearchWord = searchWord
}

// WithSelector adds the selector to the describe runtime envs params
func (o *DescribeRuntimeEnvsParams) WithSelector(selector *string) *DescribeRuntimeEnvsParams {
	o.SetSelector(selector)
	return o
}

// SetSelector adds the selector to the describe runtime envs params
func (o *DescribeRuntimeEnvsParams) SetSelector(selector *string) {
	o.Selector = selector
}

// WithStatus adds the status to the describe runtime envs params
func (o *DescribeRuntimeEnvsParams) WithStatus(status []string) *DescribeRuntimeEnvsParams {
	o.SetStatus(status)
	return o
}

// SetStatus adds the status to the describe runtime envs params
func (o *DescribeRuntimeEnvsParams) SetStatus(status []string) {
	o.Status = status
}

// WithVerboseValue adds the verboseValue to the describe runtime envs params
func (o *DescribeRuntimeEnvsParams) WithVerboseValue(verboseValue *int64) *DescribeRuntimeEnvsParams {
	o.SetVerboseValue(verboseValue)
	return o
}

// SetVerboseValue adds the verboseValue to the describe runtime envs params
func (o *DescribeRuntimeEnvsParams) SetVerboseValue(verboseValue *int64) {
	o.VerboseValue = verboseValue
}

// WriteToRequest writes these params to a swagger request
func (o *DescribeRuntimeEnvsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.LimitValue != nil {

		// query param limit.value
		var qrLimitValue int64
		if o.LimitValue != nil {
			qrLimitValue = *o.LimitValue
		}
		qLimitValue := swag.FormatInt64(qrLimitValue)
		if qLimitValue != "" {
			if err := r.SetQueryParam("limit.value", qLimitValue); err != nil {
				return err
			}
		}

	}

	if o.OffsetValue != nil {

		// query param offset.value
		var qrOffsetValue int64
		if o.OffsetValue != nil {
			qrOffsetValue = *o.OffsetValue
		}
		qOffsetValue := swag.FormatInt64(qrOffsetValue)
		if qOffsetValue != "" {
			if err := r.SetQueryParam("offset.value", qOffsetValue); err != nil {
				return err
			}
		}

	}

	valuesOwner := o.Owner

	joinedOwner := swag.JoinByFormat(valuesOwner, "")
	// query array param owner
	if err := r.SetQueryParam("owner", joinedOwner...); err != nil {
		return err
	}

	valuesRuntimeEnvID := o.RuntimeEnvID

	joinedRuntimeEnvID := swag.JoinByFormat(valuesRuntimeEnvID, "")
	// query array param runtime_env_id
	if err := r.SetQueryParam("runtime_env_id", joinedRuntimeEnvID...); err != nil {
		return err
	}

	if o.SearchWord != nil {

		// query param search_word
		var qrSearchWord string
		if o.SearchWord != nil {
			qrSearchWord = *o.SearchWord
		}
		qSearchWord := qrSearchWord
		if qSearchWord != "" {
			if err := r.SetQueryParam("search_word", qSearchWord); err != nil {
				return err
			}
		}

	}

	if o.Selector != nil {

		// query param selector
		var qrSelector string
		if o.Selector != nil {
			qrSelector = *o.Selector
		}
		qSelector := qrSelector
		if qSelector != "" {
			if err := r.SetQueryParam("selector", qSelector); err != nil {
				return err
			}
		}

	}

	valuesStatus := o.Status

	joinedStatus := swag.JoinByFormat(valuesStatus, "")
	// query array param status
	if err := r.SetQueryParam("status", joinedStatus...); err != nil {
		return err
	}

	if o.VerboseValue != nil {

		// query param verbose.value
		var qrVerboseValue int64
		if o.VerboseValue != nil {
			qrVerboseValue = *o.VerboseValue
		}
		qVerboseValue := swag.FormatInt64(qrVerboseValue)
		if qVerboseValue != "" {
			if err := r.SetQueryParam("verbose.value", qVerboseValue); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
