// Code generated by go-swagger; DO NOT EDIT.

package access_manager

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

// NewDescribeRolesParams creates a new DescribeRolesParams object
// with the default values initialized.
func NewDescribeRolesParams() *DescribeRolesParams {
	var ()
	return &DescribeRolesParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewDescribeRolesParamsWithTimeout creates a new DescribeRolesParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewDescribeRolesParamsWithTimeout(timeout time.Duration) *DescribeRolesParams {
	var ()
	return &DescribeRolesParams{

		timeout: timeout,
	}
}

// NewDescribeRolesParamsWithContext creates a new DescribeRolesParams object
// with the default values initialized, and the ability to set a context for a request
func NewDescribeRolesParamsWithContext(ctx context.Context) *DescribeRolesParams {
	var ()
	return &DescribeRolesParams{

		Context: ctx,
	}
}

// NewDescribeRolesParamsWithHTTPClient creates a new DescribeRolesParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewDescribeRolesParamsWithHTTPClient(client *http.Client) *DescribeRolesParams {
	var ()
	return &DescribeRolesParams{
		HTTPClient: client,
	}
}

/*DescribeRolesParams contains all the parameters to send to the API endpoint
for the describe roles operation typically these are written to a http.Request
*/
type DescribeRolesParams struct {

	/*ActionBundleID*/
	ActionBundleID []string
	/*Limit*/
	Limit *int64
	/*Offset*/
	Offset *int64
	/*Portal*/
	Portal []string
	/*Reverse*/
	Reverse *bool
	/*RoleID*/
	RoleID []string
	/*RoleName*/
	RoleName []string
	/*SearchWord*/
	SearchWord *string
	/*SortKey*/
	SortKey *string
	/*Status*/
	Status []string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the describe roles params
func (o *DescribeRolesParams) WithTimeout(timeout time.Duration) *DescribeRolesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the describe roles params
func (o *DescribeRolesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the describe roles params
func (o *DescribeRolesParams) WithContext(ctx context.Context) *DescribeRolesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the describe roles params
func (o *DescribeRolesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the describe roles params
func (o *DescribeRolesParams) WithHTTPClient(client *http.Client) *DescribeRolesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the describe roles params
func (o *DescribeRolesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithActionBundleID adds the actionBundleID to the describe roles params
func (o *DescribeRolesParams) WithActionBundleID(actionBundleID []string) *DescribeRolesParams {
	o.SetActionBundleID(actionBundleID)
	return o
}

// SetActionBundleID adds the actionBundleId to the describe roles params
func (o *DescribeRolesParams) SetActionBundleID(actionBundleID []string) {
	o.ActionBundleID = actionBundleID
}

// WithLimit adds the limit to the describe roles params
func (o *DescribeRolesParams) WithLimit(limit *int64) *DescribeRolesParams {
	o.SetLimit(limit)
	return o
}

// SetLimit adds the limit to the describe roles params
func (o *DescribeRolesParams) SetLimit(limit *int64) {
	o.Limit = limit
}

// WithOffset adds the offset to the describe roles params
func (o *DescribeRolesParams) WithOffset(offset *int64) *DescribeRolesParams {
	o.SetOffset(offset)
	return o
}

// SetOffset adds the offset to the describe roles params
func (o *DescribeRolesParams) SetOffset(offset *int64) {
	o.Offset = offset
}

// WithPortal adds the portal to the describe roles params
func (o *DescribeRolesParams) WithPortal(portal []string) *DescribeRolesParams {
	o.SetPortal(portal)
	return o
}

// SetPortal adds the portal to the describe roles params
func (o *DescribeRolesParams) SetPortal(portal []string) {
	o.Portal = portal
}

// WithReverse adds the reverse to the describe roles params
func (o *DescribeRolesParams) WithReverse(reverse *bool) *DescribeRolesParams {
	o.SetReverse(reverse)
	return o
}

// SetReverse adds the reverse to the describe roles params
func (o *DescribeRolesParams) SetReverse(reverse *bool) {
	o.Reverse = reverse
}

// WithRoleID adds the roleID to the describe roles params
func (o *DescribeRolesParams) WithRoleID(roleID []string) *DescribeRolesParams {
	o.SetRoleID(roleID)
	return o
}

// SetRoleID adds the roleId to the describe roles params
func (o *DescribeRolesParams) SetRoleID(roleID []string) {
	o.RoleID = roleID
}

// WithRoleName adds the roleName to the describe roles params
func (o *DescribeRolesParams) WithRoleName(roleName []string) *DescribeRolesParams {
	o.SetRoleName(roleName)
	return o
}

// SetRoleName adds the roleName to the describe roles params
func (o *DescribeRolesParams) SetRoleName(roleName []string) {
	o.RoleName = roleName
}

// WithSearchWord adds the searchWord to the describe roles params
func (o *DescribeRolesParams) WithSearchWord(searchWord *string) *DescribeRolesParams {
	o.SetSearchWord(searchWord)
	return o
}

// SetSearchWord adds the searchWord to the describe roles params
func (o *DescribeRolesParams) SetSearchWord(searchWord *string) {
	o.SearchWord = searchWord
}

// WithSortKey adds the sortKey to the describe roles params
func (o *DescribeRolesParams) WithSortKey(sortKey *string) *DescribeRolesParams {
	o.SetSortKey(sortKey)
	return o
}

// SetSortKey adds the sortKey to the describe roles params
func (o *DescribeRolesParams) SetSortKey(sortKey *string) {
	o.SortKey = sortKey
}

// WithStatus adds the status to the describe roles params
func (o *DescribeRolesParams) WithStatus(status []string) *DescribeRolesParams {
	o.SetStatus(status)
	return o
}

// SetStatus adds the status to the describe roles params
func (o *DescribeRolesParams) SetStatus(status []string) {
	o.Status = status
}

// WriteToRequest writes these params to a swagger request
func (o *DescribeRolesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	valuesActionBundleID := o.ActionBundleID

	joinedActionBundleID := swag.JoinByFormat(valuesActionBundleID, "multi")
	// query array param action_bundle_id
	if err := r.SetQueryParam("action_bundle_id", joinedActionBundleID...); err != nil {
		return err
	}

	if o.Limit != nil {

		// query param limit
		var qrLimit int64
		if o.Limit != nil {
			qrLimit = *o.Limit
		}
		qLimit := swag.FormatInt64(qrLimit)
		if qLimit != "" {
			if err := r.SetQueryParam("limit", qLimit); err != nil {
				return err
			}
		}

	}

	if o.Offset != nil {

		// query param offset
		var qrOffset int64
		if o.Offset != nil {
			qrOffset = *o.Offset
		}
		qOffset := swag.FormatInt64(qrOffset)
		if qOffset != "" {
			if err := r.SetQueryParam("offset", qOffset); err != nil {
				return err
			}
		}

	}

	valuesPortal := o.Portal

	joinedPortal := swag.JoinByFormat(valuesPortal, "multi")
	// query array param portal
	if err := r.SetQueryParam("portal", joinedPortal...); err != nil {
		return err
	}

	if o.Reverse != nil {

		// query param reverse
		var qrReverse bool
		if o.Reverse != nil {
			qrReverse = *o.Reverse
		}
		qReverse := swag.FormatBool(qrReverse)
		if qReverse != "" {
			if err := r.SetQueryParam("reverse", qReverse); err != nil {
				return err
			}
		}

	}

	valuesRoleID := o.RoleID

	joinedRoleID := swag.JoinByFormat(valuesRoleID, "multi")
	// query array param role_id
	if err := r.SetQueryParam("role_id", joinedRoleID...); err != nil {
		return err
	}

	valuesRoleName := o.RoleName

	joinedRoleName := swag.JoinByFormat(valuesRoleName, "multi")
	// query array param role_name
	if err := r.SetQueryParam("role_name", joinedRoleName...); err != nil {
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

	if o.SortKey != nil {

		// query param sort_key
		var qrSortKey string
		if o.SortKey != nil {
			qrSortKey = *o.SortKey
		}
		qSortKey := qrSortKey
		if qSortKey != "" {
			if err := r.SetQueryParam("sort_key", qSortKey); err != nil {
				return err
			}
		}

	}

	valuesStatus := o.Status

	joinedStatus := swag.JoinByFormat(valuesStatus, "multi")
	// query array param status
	if err := r.SetQueryParam("status", joinedStatus...); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
