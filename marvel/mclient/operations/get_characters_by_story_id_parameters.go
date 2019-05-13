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

// NewGetCharactersByStoryIDParams creates a new GetCharactersByStoryIDParams object
// with the default values initialized.
func NewGetCharactersByStoryIDParams() *GetCharactersByStoryIDParams {
	var ()
	return &GetCharactersByStoryIDParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetCharactersByStoryIDParamsWithTimeout creates a new GetCharactersByStoryIDParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetCharactersByStoryIDParamsWithTimeout(timeout time.Duration) *GetCharactersByStoryIDParams {
	var ()
	return &GetCharactersByStoryIDParams{

		timeout: timeout,
	}
}

// NewGetCharactersByStoryIDParamsWithContext creates a new GetCharactersByStoryIDParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetCharactersByStoryIDParamsWithContext(ctx context.Context) *GetCharactersByStoryIDParams {
	var ()
	return &GetCharactersByStoryIDParams{

		Context: ctx,
	}
}

// NewGetCharactersByStoryIDParamsWithHTTPClient creates a new GetCharactersByStoryIDParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetCharactersByStoryIDParamsWithHTTPClient(client *http.Client) *GetCharactersByStoryIDParams {
	var ()
	return &GetCharactersByStoryIDParams{
		HTTPClient: client,
	}
}

/*GetCharactersByStoryIDParams contains all the parameters to send to the API endpoint
for the get characters by story Id operation typically these are written to a http.Request
*/
type GetCharactersByStoryIDParams struct {

	/*Apikey
	  [Auth] public apikey

	*/
	Apikey string
	/*Comics
	  Return only characters which appear in the specified comics (accepts a comma-separated list of ids).

	*/
	Comics []int32
	/*Events
	  Return only characters which appear comics that took place in the specified events (accepts a comma-separated list of ids).

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
	  Return only characters which have been modified since the specified date.

	*/
	ModifiedSince *strfmt.Date
	/*Name
	  Return only characters matching the specified full character name (e.g. Spider-Man).

	*/
	Name *string
	/*NameStartsWith
	  Return characters with names that begin with the specified string (e.g. Sp).

	*/
	NameStartsWith *string
	/*Offset
	  Skip the specified number of resources in the result set.

	*/
	Offset *int32
	/*OrderBy
	  Order the result set by a field or fields. Add a "-" to the value sort in descending order. Multiple values are given priority in the order in which they are passed.

	*/
	OrderBy []string
	/*Series
	  Return only characters which appear the specified series (accepts a comma-separated list of ids).

	*/
	Series []int32
	/*StoryID
	  The story ID.

	*/
	StoryID int32
	/*Ts
	  [Auth] timestamp

	*/
	Ts string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) WithTimeout(timeout time.Duration) *GetCharactersByStoryIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) WithContext(ctx context.Context) *GetCharactersByStoryIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) WithHTTPClient(client *http.Client) *GetCharactersByStoryIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithApikey adds the apikey to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) WithApikey(apikey string) *GetCharactersByStoryIDParams {
	o.SetApikey(apikey)
	return o
}

// SetApikey adds the apikey to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) SetApikey(apikey string) {
	o.Apikey = apikey
}

// WithComics adds the comics to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) WithComics(comics []int32) *GetCharactersByStoryIDParams {
	o.SetComics(comics)
	return o
}

// SetComics adds the comics to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) SetComics(comics []int32) {
	o.Comics = comics
}

// WithEvents adds the events to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) WithEvents(events []int32) *GetCharactersByStoryIDParams {
	o.SetEvents(events)
	return o
}

// SetEvents adds the events to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) SetEvents(events []int32) {
	o.Events = events
}

// WithHash adds the hash to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) WithHash(hash string) *GetCharactersByStoryIDParams {
	o.SetHash(hash)
	return o
}

// SetHash adds the hash to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) SetHash(hash string) {
	o.Hash = hash
}

// WithLimit adds the limit to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) WithLimit(limit *int32) *GetCharactersByStoryIDParams {
	o.SetLimit(limit)
	return o
}

// SetLimit adds the limit to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) SetLimit(limit *int32) {
	o.Limit = limit
}

// WithModifiedSince adds the modifiedSince to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) WithModifiedSince(modifiedSince *strfmt.Date) *GetCharactersByStoryIDParams {
	o.SetModifiedSince(modifiedSince)
	return o
}

// SetModifiedSince adds the modifiedSince to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) SetModifiedSince(modifiedSince *strfmt.Date) {
	o.ModifiedSince = modifiedSince
}

// WithName adds the name to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) WithName(name *string) *GetCharactersByStoryIDParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) SetName(name *string) {
	o.Name = name
}

// WithNameStartsWith adds the nameStartsWith to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) WithNameStartsWith(nameStartsWith *string) *GetCharactersByStoryIDParams {
	o.SetNameStartsWith(nameStartsWith)
	return o
}

// SetNameStartsWith adds the nameStartsWith to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) SetNameStartsWith(nameStartsWith *string) {
	o.NameStartsWith = nameStartsWith
}

// WithOffset adds the offset to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) WithOffset(offset *int32) *GetCharactersByStoryIDParams {
	o.SetOffset(offset)
	return o
}

// SetOffset adds the offset to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) SetOffset(offset *int32) {
	o.Offset = offset
}

// WithOrderBy adds the orderBy to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) WithOrderBy(orderBy []string) *GetCharactersByStoryIDParams {
	o.SetOrderBy(orderBy)
	return o
}

// SetOrderBy adds the orderBy to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) SetOrderBy(orderBy []string) {
	o.OrderBy = orderBy
}

// WithSeries adds the series to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) WithSeries(series []int32) *GetCharactersByStoryIDParams {
	o.SetSeries(series)
	return o
}

// SetSeries adds the series to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) SetSeries(series []int32) {
	o.Series = series
}

// WithStoryID adds the storyID to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) WithStoryID(storyID int32) *GetCharactersByStoryIDParams {
	o.SetStoryID(storyID)
	return o
}

// SetStoryID adds the storyId to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) SetStoryID(storyID int32) {
	o.StoryID = storyID
}

// WithTs adds the ts to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) WithTs(ts string) *GetCharactersByStoryIDParams {
	o.SetTs(ts)
	return o
}

// SetTs adds the ts to the get characters by story Id params
func (o *GetCharactersByStoryIDParams) SetTs(ts string) {
	o.Ts = ts
}

// WriteToRequest writes these params to a swagger request
func (o *GetCharactersByStoryIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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

	var valuesComics []string
	for _, v := range o.Comics {
		valuesComics = append(valuesComics, swag.FormatInt32(v))
	}

	joinedComics := swag.JoinByFormat(valuesComics, "")
	// query array param comics
	if err := r.SetQueryParam("comics", joinedComics...); err != nil {
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

	if o.Name != nil {

		// query param name
		var qrName string
		if o.Name != nil {
			qrName = *o.Name
		}
		qName := qrName
		if qName != "" {
			if err := r.SetQueryParam("name", qName); err != nil {
				return err
			}
		}

	}

	if o.NameStartsWith != nil {

		// query param nameStartsWith
		var qrNameStartsWith string
		if o.NameStartsWith != nil {
			qrNameStartsWith = *o.NameStartsWith
		}
		qNameStartsWith := qrNameStartsWith
		if qNameStartsWith != "" {
			if err := r.SetQueryParam("nameStartsWith", qNameStartsWith); err != nil {
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

	// path param storyId
	if err := r.SetPathParam("storyId", swag.FormatInt32(o.StoryID)); err != nil {
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
