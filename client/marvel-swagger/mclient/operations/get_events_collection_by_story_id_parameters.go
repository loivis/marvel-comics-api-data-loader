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

// NewGetEventsCollectionByStoryIDParams creates a new GetEventsCollectionByStoryIDParams object
// with the default values initialized.
func NewGetEventsCollectionByStoryIDParams() *GetEventsCollectionByStoryIDParams {
	var ()
	return &GetEventsCollectionByStoryIDParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetEventsCollectionByStoryIDParamsWithTimeout creates a new GetEventsCollectionByStoryIDParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetEventsCollectionByStoryIDParamsWithTimeout(timeout time.Duration) *GetEventsCollectionByStoryIDParams {
	var ()
	return &GetEventsCollectionByStoryIDParams{

		timeout: timeout,
	}
}

// NewGetEventsCollectionByStoryIDParamsWithContext creates a new GetEventsCollectionByStoryIDParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetEventsCollectionByStoryIDParamsWithContext(ctx context.Context) *GetEventsCollectionByStoryIDParams {
	var ()
	return &GetEventsCollectionByStoryIDParams{

		Context: ctx,
	}
}

// NewGetEventsCollectionByStoryIDParamsWithHTTPClient creates a new GetEventsCollectionByStoryIDParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetEventsCollectionByStoryIDParamsWithHTTPClient(client *http.Client) *GetEventsCollectionByStoryIDParams {
	var ()
	return &GetEventsCollectionByStoryIDParams{
		HTTPClient: client,
	}
}

/*GetEventsCollectionByStoryIDParams contains all the parameters to send to the API endpoint
for the get events collection by story Id operation typically these are written to a http.Request
*/
type GetEventsCollectionByStoryIDParams struct {

	/*Apikey
	  [Auth] public apikey

	*/
	Apikey string
	/*Characters
	  Return only events which feature the specified characters (accepts a comma-separated list of ids).

	*/
	Characters []int32
	/*Comics
	  Return only events which take place in the specified comics (accepts a comma-separated list of ids).

	*/
	Comics []int32
	/*Creators
	  Return only events which feature work by the specified creators (accepts a comma-separated list of ids).

	*/
	Creators []int32
	/*Hash
	  [Auth] md5 digest of concatenation of ts, private key, public key

	*/
	Hash string
	/*Limit
	  Limit the result set to the specified number of resources.

	*/
	Limit *int32
	/*ModifiedSince
	  Return only events which have been modified since the specified date.

	*/
	ModifiedSince *strfmt.Date
	/*Name
	  Filter the event list by name.

	*/
	Name *string
	/*NameStartsWith
	  Return events with names that begin with the specified string (e.g. Sp).

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
	  Return only events which are part of the specified series (accepts a comma-separated list of ids).

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

// WithTimeout adds the timeout to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) WithTimeout(timeout time.Duration) *GetEventsCollectionByStoryIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) WithContext(ctx context.Context) *GetEventsCollectionByStoryIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) WithHTTPClient(client *http.Client) *GetEventsCollectionByStoryIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithApikey adds the apikey to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) WithApikey(apikey string) *GetEventsCollectionByStoryIDParams {
	o.SetApikey(apikey)
	return o
}

// SetApikey adds the apikey to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) SetApikey(apikey string) {
	o.Apikey = apikey
}

// WithCharacters adds the characters to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) WithCharacters(characters []int32) *GetEventsCollectionByStoryIDParams {
	o.SetCharacters(characters)
	return o
}

// SetCharacters adds the characters to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) SetCharacters(characters []int32) {
	o.Characters = characters
}

// WithComics adds the comics to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) WithComics(comics []int32) *GetEventsCollectionByStoryIDParams {
	o.SetComics(comics)
	return o
}

// SetComics adds the comics to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) SetComics(comics []int32) {
	o.Comics = comics
}

// WithCreators adds the creators to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) WithCreators(creators []int32) *GetEventsCollectionByStoryIDParams {
	o.SetCreators(creators)
	return o
}

// SetCreators adds the creators to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) SetCreators(creators []int32) {
	o.Creators = creators
}

// WithHash adds the hash to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) WithHash(hash string) *GetEventsCollectionByStoryIDParams {
	o.SetHash(hash)
	return o
}

// SetHash adds the hash to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) SetHash(hash string) {
	o.Hash = hash
}

// WithLimit adds the limit to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) WithLimit(limit *int32) *GetEventsCollectionByStoryIDParams {
	o.SetLimit(limit)
	return o
}

// SetLimit adds the limit to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) SetLimit(limit *int32) {
	o.Limit = limit
}

// WithModifiedSince adds the modifiedSince to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) WithModifiedSince(modifiedSince *strfmt.Date) *GetEventsCollectionByStoryIDParams {
	o.SetModifiedSince(modifiedSince)
	return o
}

// SetModifiedSince adds the modifiedSince to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) SetModifiedSince(modifiedSince *strfmt.Date) {
	o.ModifiedSince = modifiedSince
}

// WithName adds the name to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) WithName(name *string) *GetEventsCollectionByStoryIDParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) SetName(name *string) {
	o.Name = name
}

// WithNameStartsWith adds the nameStartsWith to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) WithNameStartsWith(nameStartsWith *string) *GetEventsCollectionByStoryIDParams {
	o.SetNameStartsWith(nameStartsWith)
	return o
}

// SetNameStartsWith adds the nameStartsWith to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) SetNameStartsWith(nameStartsWith *string) {
	o.NameStartsWith = nameStartsWith
}

// WithOffset adds the offset to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) WithOffset(offset *int32) *GetEventsCollectionByStoryIDParams {
	o.SetOffset(offset)
	return o
}

// SetOffset adds the offset to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) SetOffset(offset *int32) {
	o.Offset = offset
}

// WithOrderBy adds the orderBy to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) WithOrderBy(orderBy []string) *GetEventsCollectionByStoryIDParams {
	o.SetOrderBy(orderBy)
	return o
}

// SetOrderBy adds the orderBy to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) SetOrderBy(orderBy []string) {
	o.OrderBy = orderBy
}

// WithSeries adds the series to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) WithSeries(series []int32) *GetEventsCollectionByStoryIDParams {
	o.SetSeries(series)
	return o
}

// SetSeries adds the series to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) SetSeries(series []int32) {
	o.Series = series
}

// WithStoryID adds the storyID to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) WithStoryID(storyID int32) *GetEventsCollectionByStoryIDParams {
	o.SetStoryID(storyID)
	return o
}

// SetStoryID adds the storyId to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) SetStoryID(storyID int32) {
	o.StoryID = storyID
}

// WithTs adds the ts to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) WithTs(ts string) *GetEventsCollectionByStoryIDParams {
	o.SetTs(ts)
	return o
}

// SetTs adds the ts to the get events collection by story Id params
func (o *GetEventsCollectionByStoryIDParams) SetTs(ts string) {
	o.Ts = ts
}

// WriteToRequest writes these params to a swagger request
func (o *GetEventsCollectionByStoryIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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

	var valuesCreators []string
	for _, v := range o.Creators {
		valuesCreators = append(valuesCreators, swag.FormatInt32(v))
	}

	joinedCreators := swag.JoinByFormat(valuesCreators, "")
	// query array param creators
	if err := r.SetQueryParam("creators", joinedCreators...); err != nil {
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