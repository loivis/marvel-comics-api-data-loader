// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetCreatorStoryCollectionParams creates a new GetCreatorStoryCollectionParams object
// with the default values initialized.
func NewGetCreatorStoryCollectionParams() *GetCreatorStoryCollectionParams {
	var ()
	return &GetCreatorStoryCollectionParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetCreatorStoryCollectionParamsWithTimeout creates a new GetCreatorStoryCollectionParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetCreatorStoryCollectionParamsWithTimeout(timeout time.Duration) *GetCreatorStoryCollectionParams {
	var ()
	return &GetCreatorStoryCollectionParams{

		timeout: timeout,
	}
}

// NewGetCreatorStoryCollectionParamsWithContext creates a new GetCreatorStoryCollectionParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetCreatorStoryCollectionParamsWithContext(ctx context.Context) *GetCreatorStoryCollectionParams {
	var ()
	return &GetCreatorStoryCollectionParams{

		Context: ctx,
	}
}

// NewGetCreatorStoryCollectionParamsWithHTTPClient creates a new GetCreatorStoryCollectionParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetCreatorStoryCollectionParamsWithHTTPClient(client *http.Client) *GetCreatorStoryCollectionParams {
	var ()
	return &GetCreatorStoryCollectionParams{
		HTTPClient: client,
	}
}

/*GetCreatorStoryCollectionParams contains all the parameters to send to the API endpoint
for the get creator story collection operation typically these are written to a http.Request
*/
type GetCreatorStoryCollectionParams struct {

	/*Apikey
	  [Auth] public apikey

	*/
	Apikey string
	/*Characters
	  Return only stories which feature the specified characters (accepts a comma-separated list of ids).

	*/
	Characters []int32
	/*Comics
	  Return only stories contained in the specified comics (accepts a comma-separated list of ids).

	*/
	Comics []int32
	/*CreatorID
	  The ID of the creator.

	*/
	CreatorID int32
	/*Events
	  Return only stories which take place during the specified events (accepts a comma-separated list of ids).

	*/
	Events []int32
	/*Hash
	  [Auth] md5 digest of concatenation of ts, private key, public key

	*/
	Hash string
	/*Limit
	  Limit the result set to the specified number of resources.

	*/
	Limit *int32
	/*ModifiedSince
	  Return only stories which have been modified since the specified date.

	*/
	ModifiedSince *strfmt.Date
	/*Offset
	  Skip the specified number of resources in the result set.

	*/
	Offset *int32
	/*OrderBy
	  Order the result set by a field or fields. Add a "-" to the value sort in descending order. Multiple values are given priority in the order in which they are passed.

	*/
	OrderBy []string
	/*Series
	  Return only stories contained the specified series (accepts a comma-separated list of ids).

	*/
	Series []int32
	/*Ts
	  [Auth] timestamp

	*/
	Ts string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get creator story collection params
func (o *GetCreatorStoryCollectionParams) WithTimeout(timeout time.Duration) *GetCreatorStoryCollectionParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get creator story collection params
func (o *GetCreatorStoryCollectionParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get creator story collection params
func (o *GetCreatorStoryCollectionParams) WithContext(ctx context.Context) *GetCreatorStoryCollectionParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get creator story collection params
func (o *GetCreatorStoryCollectionParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get creator story collection params
func (o *GetCreatorStoryCollectionParams) WithHTTPClient(client *http.Client) *GetCreatorStoryCollectionParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get creator story collection params
func (o *GetCreatorStoryCollectionParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithApikey adds the apikey to the get creator story collection params
func (o *GetCreatorStoryCollectionParams) WithApikey(apikey string) *GetCreatorStoryCollectionParams {
	o.SetApikey(apikey)
	return o
}

// SetApikey adds the apikey to the get creator story collection params
func (o *GetCreatorStoryCollectionParams) SetApikey(apikey string) {
	o.Apikey = apikey
}

// WithCharacters adds the characters to the get creator story collection params
func (o *GetCreatorStoryCollectionParams) WithCharacters(characters []int32) *GetCreatorStoryCollectionParams {
	o.SetCharacters(characters)
	return o
}

// SetCharacters adds the characters to the get creator story collection params
func (o *GetCreatorStoryCollectionParams) SetCharacters(characters []int32) {
	o.Characters = characters
}

// WithComics adds the comics to the get creator story collection params
func (o *GetCreatorStoryCollectionParams) WithComics(comics []int32) *GetCreatorStoryCollectionParams {
	o.SetComics(comics)
	return o
}

// SetComics adds the comics to the get creator story collection params
func (o *GetCreatorStoryCollectionParams) SetComics(comics []int32) {
	o.Comics = comics
}

// WithCreatorID adds the creatorID to the get creator story collection params
func (o *GetCreatorStoryCollectionParams) WithCreatorID(creatorID int32) *GetCreatorStoryCollectionParams {
	o.SetCreatorID(creatorID)
	return o
}

// SetCreatorID adds the creatorId to the get creator story collection params
func (o *GetCreatorStoryCollectionParams) SetCreatorID(creatorID int32) {
	o.CreatorID = creatorID
}

// WithEvents adds the events to the get creator story collection params
func (o *GetCreatorStoryCollectionParams) WithEvents(events []int32) *GetCreatorStoryCollectionParams {
	o.SetEvents(events)
	return o
}

// SetEvents adds the events to the get creator story collection params
func (o *GetCreatorStoryCollectionParams) SetEvents(events []int32) {
	o.Events = events
}

// WithHash adds the hash to the get creator story collection params
func (o *GetCreatorStoryCollectionParams) WithHash(hash string) *GetCreatorStoryCollectionParams {
	o.SetHash(hash)
	return o
}

// SetHash adds the hash to the get creator story collection params
func (o *GetCreatorStoryCollectionParams) SetHash(hash string) {
	o.Hash = hash
}

// WithLimit adds the limit to the get creator story collection params
func (o *GetCreatorStoryCollectionParams) WithLimit(limit *int32) *GetCreatorStoryCollectionParams {
	o.SetLimit(limit)
	return o
}

// SetLimit adds the limit to the get creator story collection params
func (o *GetCreatorStoryCollectionParams) SetLimit(limit *int32) {
	o.Limit = limit
}

// WithModifiedSince adds the modifiedSince to the get creator story collection params
func (o *GetCreatorStoryCollectionParams) WithModifiedSince(modifiedSince *strfmt.Date) *GetCreatorStoryCollectionParams {
	o.SetModifiedSince(modifiedSince)
	return o
}

// SetModifiedSince adds the modifiedSince to the get creator story collection params
func (o *GetCreatorStoryCollectionParams) SetModifiedSince(modifiedSince *strfmt.Date) {
	o.ModifiedSince = modifiedSince
}

// WithOffset adds the offset to the get creator story collection params
func (o *GetCreatorStoryCollectionParams) WithOffset(offset *int32) *GetCreatorStoryCollectionParams {
	o.SetOffset(offset)
	return o
}

// SetOffset adds the offset to the get creator story collection params
func (o *GetCreatorStoryCollectionParams) SetOffset(offset *int32) {
	o.Offset = offset
}

// WithOrderBy adds the orderBy to the get creator story collection params
func (o *GetCreatorStoryCollectionParams) WithOrderBy(orderBy []string) *GetCreatorStoryCollectionParams {
	o.SetOrderBy(orderBy)
	return o
}

// SetOrderBy adds the orderBy to the get creator story collection params
func (o *GetCreatorStoryCollectionParams) SetOrderBy(orderBy []string) {
	o.OrderBy = orderBy
}

// WithSeries adds the series to the get creator story collection params
func (o *GetCreatorStoryCollectionParams) WithSeries(series []int32) *GetCreatorStoryCollectionParams {
	o.SetSeries(series)
	return o
}

// SetSeries adds the series to the get creator story collection params
func (o *GetCreatorStoryCollectionParams) SetSeries(series []int32) {
	o.Series = series
}

// WithTs adds the ts to the get creator story collection params
func (o *GetCreatorStoryCollectionParams) WithTs(ts string) *GetCreatorStoryCollectionParams {
	o.SetTs(ts)
	return o
}

// SetTs adds the ts to the get creator story collection params
func (o *GetCreatorStoryCollectionParams) SetTs(ts string) {
	o.Ts = ts
}

// WriteToRequest writes these params to a swagger request
func (o *GetCreatorStoryCollectionParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// query param apikey
	qrApikey := o.Apikey
	qApikey := qrApikey
	if qApikey != "" {
		if err := r.SetQueryParam("apikey", qApikey); err != nil {
			return err
		}
	}

	var valuesCharacters []string
	for _, v := range o.Characters {
		valuesCharacters = append(valuesCharacters, swag.FormatInt32(v))
	}

	joinedCharacters := swag.JoinByFormat(valuesCharacters, "")
	// query array param characters
	if err := r.SetQueryParam("characters", joinedCharacters...); err != nil {
		return err
	}

	var valuesComics []string
	for _, v := range o.Comics {
		valuesComics = append(valuesComics, swag.FormatInt32(v))
	}

	joinedComics := swag.JoinByFormat(valuesComics, "")
	// query array param comics
	if err := r.SetQueryParam("comics", joinedComics...); err != nil {
		return err
	}

	// path param creatorId
	if err := r.SetPathParam("creatorId", swag.FormatInt32(o.CreatorID)); err != nil {
		return err
	}

	var valuesEvents []string
	for _, v := range o.Events {
		valuesEvents = append(valuesEvents, swag.FormatInt32(v))
	}

	joinedEvents := swag.JoinByFormat(valuesEvents, "")
	// query array param events
	if err := r.SetQueryParam("events", joinedEvents...); err != nil {
		return err
	}

	// query param hash
	qrHash := o.Hash
	qHash := qrHash
	if qHash != "" {
		if err := r.SetQueryParam("hash", qHash); err != nil {
			return err
		}
	}

	if o.Limit != nil {

		// query param limit
		var qrLimit int32
		if o.Limit != nil {
			qrLimit = *o.Limit
		}
		qLimit := swag.FormatInt32(qrLimit)
		if qLimit != "" {
			if err := r.SetQueryParam("limit", qLimit); err != nil {
				return err
			}
		}

	}

	if o.ModifiedSince != nil {

		// query param modifiedSince
		var qrModifiedSince strfmt.Date
		if o.ModifiedSince != nil {
			qrModifiedSince = *o.ModifiedSince
		}
		qModifiedSince := qrModifiedSince.String()
		if qModifiedSince != "" {
			if err := r.SetQueryParam("modifiedSince", qModifiedSince); err != nil {
				return err
			}
		}

	}

	if o.Offset != nil {

		// query param offset
		var qrOffset int32
		if o.Offset != nil {
			qrOffset = *o.Offset
		}
		qOffset := swag.FormatInt32(qrOffset)
		if qOffset != "" {
			if err := r.SetQueryParam("offset", qOffset); err != nil {
				return err
			}
		}

	}

	valuesOrderBy := o.OrderBy

	joinedOrderBy := swag.JoinByFormat(valuesOrderBy, "")
	// query array param orderBy
	if err := r.SetQueryParam("orderBy", joinedOrderBy...); err != nil {
		return err
	}

	var valuesSeries []string
	for _, v := range o.Series {
		valuesSeries = append(valuesSeries, swag.FormatInt32(v))
	}

	joinedSeries := swag.JoinByFormat(valuesSeries, "")
	// query array param series
	if err := r.SetQueryParam("series", joinedSeries...); err != nil {
		return err
	}

	// query param ts
	qrTs := o.Ts
	qTs := qrTs
	if qTs != "" {
		if err := r.SetQueryParam("ts", qTs); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}