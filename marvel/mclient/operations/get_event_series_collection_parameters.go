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

// NewGetEventSeriesCollectionParams creates a new GetEventSeriesCollectionParams object
// with the default values initialized.
func NewGetEventSeriesCollectionParams() *GetEventSeriesCollectionParams {
	var ()
	return &GetEventSeriesCollectionParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetEventSeriesCollectionParamsWithTimeout creates a new GetEventSeriesCollectionParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetEventSeriesCollectionParamsWithTimeout(timeout time.Duration) *GetEventSeriesCollectionParams {
	var ()
	return &GetEventSeriesCollectionParams{

		timeout: timeout,
	}
}

// NewGetEventSeriesCollectionParamsWithContext creates a new GetEventSeriesCollectionParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetEventSeriesCollectionParamsWithContext(ctx context.Context) *GetEventSeriesCollectionParams {
	var ()
	return &GetEventSeriesCollectionParams{

		Context: ctx,
	}
}

// NewGetEventSeriesCollectionParamsWithHTTPClient creates a new GetEventSeriesCollectionParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetEventSeriesCollectionParamsWithHTTPClient(client *http.Client) *GetEventSeriesCollectionParams {
	var ()
	return &GetEventSeriesCollectionParams{
		HTTPClient: client,
	}
}

/*GetEventSeriesCollectionParams contains all the parameters to send to the API endpoint
for the get event series collection operation typically these are written to a http.Request
*/
type GetEventSeriesCollectionParams struct {

	/*Apikey
	  [Auth] public apikey

	*/
	Apikey string
	/*Characters
	  Return only series which feature the specified characters (accepts a comma-separated list of ids).

	*/
	Characters []int32
	/*Comics
	  Return only series which contain the specified comics (accepts a comma-separated list of ids).

	*/
	Comics []int32
	/*Contains
	  Return only series containing one or more comics with the specified format.

	*/
	Contains []string
	/*Creators
	  Return only series which feature work by the specified creators (accepts a comma-separated list of ids).

	*/
	Creators []int32
	/*EventID
	  The event ID.

	*/
	EventID int32
	/*Hash
	  [Auth] md5 digest of concatenation of ts, private key, public key

	*/
	Hash string
	/*Limit
	  Limit the result set to the specified number of resources.

	*/
	Limit *int32
	/*ModifiedSince
	  Return only series which have been modified since the specified date.

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
	/*SeriesType
	  Filter the series by publication frequency type.

	*/
	SeriesType *string
	/*StartYear
	  Return only series matching the specified start year.

	*/
	StartYear *int32
	/*Stories
	  Return only series which contain the specified stories (accepts a comma-separated list of ids).

	*/
	Stories []int32
	/*Title
	  Filter by series title.

	*/
	Title *string
	/*TitleStartsWith
	  Return series with titles that begin with the specified string (e.g. Sp).

	*/
	TitleStartsWith *string
	/*Ts
	  [Auth] timestamp

	*/
	Ts string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get event series collection params
func (o *GetEventSeriesCollectionParams) WithTimeout(timeout time.Duration) *GetEventSeriesCollectionParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get event series collection params
func (o *GetEventSeriesCollectionParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get event series collection params
func (o *GetEventSeriesCollectionParams) WithContext(ctx context.Context) *GetEventSeriesCollectionParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get event series collection params
func (o *GetEventSeriesCollectionParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get event series collection params
func (o *GetEventSeriesCollectionParams) WithHTTPClient(client *http.Client) *GetEventSeriesCollectionParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get event series collection params
func (o *GetEventSeriesCollectionParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithApikey adds the apikey to the get event series collection params
func (o *GetEventSeriesCollectionParams) WithApikey(apikey string) *GetEventSeriesCollectionParams {
	o.SetApikey(apikey)
	return o
}

// SetApikey adds the apikey to the get event series collection params
func (o *GetEventSeriesCollectionParams) SetApikey(apikey string) {
	o.Apikey = apikey
}

// WithCharacters adds the characters to the get event series collection params
func (o *GetEventSeriesCollectionParams) WithCharacters(characters []int32) *GetEventSeriesCollectionParams {
	o.SetCharacters(characters)
	return o
}

// SetCharacters adds the characters to the get event series collection params
func (o *GetEventSeriesCollectionParams) SetCharacters(characters []int32) {
	o.Characters = characters
}

// WithComics adds the comics to the get event series collection params
func (o *GetEventSeriesCollectionParams) WithComics(comics []int32) *GetEventSeriesCollectionParams {
	o.SetComics(comics)
	return o
}

// SetComics adds the comics to the get event series collection params
func (o *GetEventSeriesCollectionParams) SetComics(comics []int32) {
	o.Comics = comics
}

// WithContains adds the contains to the get event series collection params
func (o *GetEventSeriesCollectionParams) WithContains(contains []string) *GetEventSeriesCollectionParams {
	o.SetContains(contains)
	return o
}

// SetContains adds the contains to the get event series collection params
func (o *GetEventSeriesCollectionParams) SetContains(contains []string) {
	o.Contains = contains
}

// WithCreators adds the creators to the get event series collection params
func (o *GetEventSeriesCollectionParams) WithCreators(creators []int32) *GetEventSeriesCollectionParams {
	o.SetCreators(creators)
	return o
}

// SetCreators adds the creators to the get event series collection params
func (o *GetEventSeriesCollectionParams) SetCreators(creators []int32) {
	o.Creators = creators
}

// WithEventID adds the eventID to the get event series collection params
func (o *GetEventSeriesCollectionParams) WithEventID(eventID int32) *GetEventSeriesCollectionParams {
	o.SetEventID(eventID)
	return o
}

// SetEventID adds the eventId to the get event series collection params
func (o *GetEventSeriesCollectionParams) SetEventID(eventID int32) {
	o.EventID = eventID
}

// WithHash adds the hash to the get event series collection params
func (o *GetEventSeriesCollectionParams) WithHash(hash string) *GetEventSeriesCollectionParams {
	o.SetHash(hash)
	return o
}

// SetHash adds the hash to the get event series collection params
func (o *GetEventSeriesCollectionParams) SetHash(hash string) {
	o.Hash = hash
}

// WithLimit adds the limit to the get event series collection params
func (o *GetEventSeriesCollectionParams) WithLimit(limit *int32) *GetEventSeriesCollectionParams {
	o.SetLimit(limit)
	return o
}

// SetLimit adds the limit to the get event series collection params
func (o *GetEventSeriesCollectionParams) SetLimit(limit *int32) {
	o.Limit = limit
}

// WithModifiedSince adds the modifiedSince to the get event series collection params
func (o *GetEventSeriesCollectionParams) WithModifiedSince(modifiedSince *strfmt.Date) *GetEventSeriesCollectionParams {
	o.SetModifiedSince(modifiedSince)
	return o
}

// SetModifiedSince adds the modifiedSince to the get event series collection params
func (o *GetEventSeriesCollectionParams) SetModifiedSince(modifiedSince *strfmt.Date) {
	o.ModifiedSince = modifiedSince
}

// WithOffset adds the offset to the get event series collection params
func (o *GetEventSeriesCollectionParams) WithOffset(offset *int32) *GetEventSeriesCollectionParams {
	o.SetOffset(offset)
	return o
}

// SetOffset adds the offset to the get event series collection params
func (o *GetEventSeriesCollectionParams) SetOffset(offset *int32) {
	o.Offset = offset
}

// WithOrderBy adds the orderBy to the get event series collection params
func (o *GetEventSeriesCollectionParams) WithOrderBy(orderBy []string) *GetEventSeriesCollectionParams {
	o.SetOrderBy(orderBy)
	return o
}

// SetOrderBy adds the orderBy to the get event series collection params
func (o *GetEventSeriesCollectionParams) SetOrderBy(orderBy []string) {
	o.OrderBy = orderBy
}

// WithSeriesType adds the seriesType to the get event series collection params
func (o *GetEventSeriesCollectionParams) WithSeriesType(seriesType *string) *GetEventSeriesCollectionParams {
	o.SetSeriesType(seriesType)
	return o
}

// SetSeriesType adds the seriesType to the get event series collection params
func (o *GetEventSeriesCollectionParams) SetSeriesType(seriesType *string) {
	o.SeriesType = seriesType
}

// WithStartYear adds the startYear to the get event series collection params
func (o *GetEventSeriesCollectionParams) WithStartYear(startYear *int32) *GetEventSeriesCollectionParams {
	o.SetStartYear(startYear)
	return o
}

// SetStartYear adds the startYear to the get event series collection params
func (o *GetEventSeriesCollectionParams) SetStartYear(startYear *int32) {
	o.StartYear = startYear
}

// WithStories adds the stories to the get event series collection params
func (o *GetEventSeriesCollectionParams) WithStories(stories []int32) *GetEventSeriesCollectionParams {
	o.SetStories(stories)
	return o
}

// SetStories adds the stories to the get event series collection params
func (o *GetEventSeriesCollectionParams) SetStories(stories []int32) {
	o.Stories = stories
}

// WithTitle adds the title to the get event series collection params
func (o *GetEventSeriesCollectionParams) WithTitle(title *string) *GetEventSeriesCollectionParams {
	o.SetTitle(title)
	return o
}

// SetTitle adds the title to the get event series collection params
func (o *GetEventSeriesCollectionParams) SetTitle(title *string) {
	o.Title = title
}

// WithTitleStartsWith adds the titleStartsWith to the get event series collection params
func (o *GetEventSeriesCollectionParams) WithTitleStartsWith(titleStartsWith *string) *GetEventSeriesCollectionParams {
	o.SetTitleStartsWith(titleStartsWith)
	return o
}

// SetTitleStartsWith adds the titleStartsWith to the get event series collection params
func (o *GetEventSeriesCollectionParams) SetTitleStartsWith(titleStartsWith *string) {
	o.TitleStartsWith = titleStartsWith
}

// WithTs adds the ts to the get event series collection params
func (o *GetEventSeriesCollectionParams) WithTs(ts string) *GetEventSeriesCollectionParams {
	o.SetTs(ts)
	return o
}

// SetTs adds the ts to the get event series collection params
func (o *GetEventSeriesCollectionParams) SetTs(ts string) {
	o.Ts = ts
}

// WriteToRequest writes these params to a swagger request
func (o *GetEventSeriesCollectionParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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

	valuesContains := o.Contains

	joinedContains := swag.JoinByFormat(valuesContains, "")
	// query array param contains
	if err := r.SetQueryParam("contains", joinedContains...); err != nil {
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

	// path param eventId
	if err := r.SetPathParam("eventId", swag.FormatInt32(o.EventID)); err != nil {
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

	if o.SeriesType != nil {

		// query param seriesType
		var qrSeriesType string
		if o.SeriesType != nil {
			qrSeriesType = *o.SeriesType
		}
		qSeriesType := qrSeriesType
		if qSeriesType != "" {
			if err := r.SetQueryParam("seriesType", qSeriesType); err != nil {
				return err
			}
		}

	}

	if o.StartYear != nil {

		// query param startYear
		var qrStartYear int32
		if o.StartYear != nil {
			qrStartYear = *o.StartYear
		}
		qStartYear := swag.FormatInt32(qrStartYear)
		if qStartYear != "" {
			if err := r.SetQueryParam("startYear", qStartYear); err != nil {
				return err
			}
		}

	}

	var valuesStories []string
	for _, v := range o.Stories {
		valuesStories = append(valuesStories, swag.FormatInt32(v))
	}

	joinedStories := swag.JoinByFormat(valuesStories, "")
	// query array param stories
	if err := r.SetQueryParam("stories", joinedStories...); err != nil {
		return err
	}

	if o.Title != nil {

		// query param title
		var qrTitle string
		if o.Title != nil {
			qrTitle = *o.Title
		}
		qTitle := qrTitle
		if qTitle != "" {
			if err := r.SetQueryParam("title", qTitle); err != nil {
				return err
			}
		}

	}

	if o.TitleStartsWith != nil {

		// query param titleStartsWith
		var qrTitleStartsWith string
		if o.TitleStartsWith != nil {
			qrTitleStartsWith = *o.TitleStartsWith
		}
		qTitleStartsWith := qrTitleStartsWith
		if qTitleStartsWith != "" {
			if err := r.SetQueryParam("titleStartsWith", qTitleStartsWith); err != nil {
				return err
			}
		}

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
