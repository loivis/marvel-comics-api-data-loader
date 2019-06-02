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

// NewGetSeriesCollectionParams creates a new GetSeriesCollectionParams object
// with the default values initialized.
func NewGetSeriesCollectionParams() *GetSeriesCollectionParams {
	var ()
	return &GetSeriesCollectionParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetSeriesCollectionParamsWithTimeout creates a new GetSeriesCollectionParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetSeriesCollectionParamsWithTimeout(timeout time.Duration) *GetSeriesCollectionParams {
	var ()
	return &GetSeriesCollectionParams{

		timeout: timeout,
	}
}

// NewGetSeriesCollectionParamsWithContext creates a new GetSeriesCollectionParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetSeriesCollectionParamsWithContext(ctx context.Context) *GetSeriesCollectionParams {
	var ()
	return &GetSeriesCollectionParams{

		Context: ctx,
	}
}

// NewGetSeriesCollectionParamsWithHTTPClient creates a new GetSeriesCollectionParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetSeriesCollectionParamsWithHTTPClient(client *http.Client) *GetSeriesCollectionParams {
	var ()
	return &GetSeriesCollectionParams{
		HTTPClient: client,
	}
}

/*GetSeriesCollectionParams contains all the parameters to send to the API endpoint
for the get series collection operation typically these are written to a http.Request
*/
type GetSeriesCollectionParams struct {

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
	/*Events
	  Return only series which have comics that take place during the specified events (accepts a comma-separated list of ids).

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
	  Return only series matching the specified title.

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

// WithTimeout adds the timeout to the get series collection params
func (o *GetSeriesCollectionParams) WithTimeout(timeout time.Duration) *GetSeriesCollectionParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get series collection params
func (o *GetSeriesCollectionParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get series collection params
func (o *GetSeriesCollectionParams) WithContext(ctx context.Context) *GetSeriesCollectionParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get series collection params
func (o *GetSeriesCollectionParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get series collection params
func (o *GetSeriesCollectionParams) WithHTTPClient(client *http.Client) *GetSeriesCollectionParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get series collection params
func (o *GetSeriesCollectionParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithApikey adds the apikey to the get series collection params
func (o *GetSeriesCollectionParams) WithApikey(apikey string) *GetSeriesCollectionParams {
	o.SetApikey(apikey)
	return o
}

// SetApikey adds the apikey to the get series collection params
func (o *GetSeriesCollectionParams) SetApikey(apikey string) {
	o.Apikey = apikey
}

// WithCharacters adds the characters to the get series collection params
func (o *GetSeriesCollectionParams) WithCharacters(characters []int32) *GetSeriesCollectionParams {
	o.SetCharacters(characters)
	return o
}

// SetCharacters adds the characters to the get series collection params
func (o *GetSeriesCollectionParams) SetCharacters(characters []int32) {
	o.Characters = characters
}

// WithComics adds the comics to the get series collection params
func (o *GetSeriesCollectionParams) WithComics(comics []int32) *GetSeriesCollectionParams {
	o.SetComics(comics)
	return o
}

// SetComics adds the comics to the get series collection params
func (o *GetSeriesCollectionParams) SetComics(comics []int32) {
	o.Comics = comics
}

// WithContains adds the contains to the get series collection params
func (o *GetSeriesCollectionParams) WithContains(contains []string) *GetSeriesCollectionParams {
	o.SetContains(contains)
	return o
}

// SetContains adds the contains to the get series collection params
func (o *GetSeriesCollectionParams) SetContains(contains []string) {
	o.Contains = contains
}

// WithCreators adds the creators to the get series collection params
func (o *GetSeriesCollectionParams) WithCreators(creators []int32) *GetSeriesCollectionParams {
	o.SetCreators(creators)
	return o
}

// SetCreators adds the creators to the get series collection params
func (o *GetSeriesCollectionParams) SetCreators(creators []int32) {
	o.Creators = creators
}

// WithEvents adds the events to the get series collection params
func (o *GetSeriesCollectionParams) WithEvents(events []int32) *GetSeriesCollectionParams {
	o.SetEvents(events)
	return o
}

// SetEvents adds the events to the get series collection params
func (o *GetSeriesCollectionParams) SetEvents(events []int32) {
	o.Events = events
}

// WithHash adds the hash to the get series collection params
func (o *GetSeriesCollectionParams) WithHash(hash string) *GetSeriesCollectionParams {
	o.SetHash(hash)
	return o
}

// SetHash adds the hash to the get series collection params
func (o *GetSeriesCollectionParams) SetHash(hash string) {
	o.Hash = hash
}

// WithLimit adds the limit to the get series collection params
func (o *GetSeriesCollectionParams) WithLimit(limit *int32) *GetSeriesCollectionParams {
	o.SetLimit(limit)
	return o
}

// SetLimit adds the limit to the get series collection params
func (o *GetSeriesCollectionParams) SetLimit(limit *int32) {
	o.Limit = limit
}

// WithModifiedSince adds the modifiedSince to the get series collection params
func (o *GetSeriesCollectionParams) WithModifiedSince(modifiedSince *strfmt.Date) *GetSeriesCollectionParams {
	o.SetModifiedSince(modifiedSince)
	return o
}

// SetModifiedSince adds the modifiedSince to the get series collection params
func (o *GetSeriesCollectionParams) SetModifiedSince(modifiedSince *strfmt.Date) {
	o.ModifiedSince = modifiedSince
}

// WithOffset adds the offset to the get series collection params
func (o *GetSeriesCollectionParams) WithOffset(offset *int32) *GetSeriesCollectionParams {
	o.SetOffset(offset)
	return o
}

// SetOffset adds the offset to the get series collection params
func (o *GetSeriesCollectionParams) SetOffset(offset *int32) {
	o.Offset = offset
}

// WithOrderBy adds the orderBy to the get series collection params
func (o *GetSeriesCollectionParams) WithOrderBy(orderBy []string) *GetSeriesCollectionParams {
	o.SetOrderBy(orderBy)
	return o
}

// SetOrderBy adds the orderBy to the get series collection params
func (o *GetSeriesCollectionParams) SetOrderBy(orderBy []string) {
	o.OrderBy = orderBy
}

// WithSeriesType adds the seriesType to the get series collection params
func (o *GetSeriesCollectionParams) WithSeriesType(seriesType *string) *GetSeriesCollectionParams {
	o.SetSeriesType(seriesType)
	return o
}

// SetSeriesType adds the seriesType to the get series collection params
func (o *GetSeriesCollectionParams) SetSeriesType(seriesType *string) {
	o.SeriesType = seriesType
}

// WithStartYear adds the startYear to the get series collection params
func (o *GetSeriesCollectionParams) WithStartYear(startYear *int32) *GetSeriesCollectionParams {
	o.SetStartYear(startYear)
	return o
}

// SetStartYear adds the startYear to the get series collection params
func (o *GetSeriesCollectionParams) SetStartYear(startYear *int32) {
	o.StartYear = startYear
}

// WithStories adds the stories to the get series collection params
func (o *GetSeriesCollectionParams) WithStories(stories []int32) *GetSeriesCollectionParams {
	o.SetStories(stories)
	return o
}

// SetStories adds the stories to the get series collection params
func (o *GetSeriesCollectionParams) SetStories(stories []int32) {
	o.Stories = stories
}

// WithTitle adds the title to the get series collection params
func (o *GetSeriesCollectionParams) WithTitle(title *string) *GetSeriesCollectionParams {
	o.SetTitle(title)
	return o
}

// SetTitle adds the title to the get series collection params
func (o *GetSeriesCollectionParams) SetTitle(title *string) {
	o.Title = title
}

// WithTitleStartsWith adds the titleStartsWith to the get series collection params
func (o *GetSeriesCollectionParams) WithTitleStartsWith(titleStartsWith *string) *GetSeriesCollectionParams {
	o.SetTitleStartsWith(titleStartsWith)
	return o
}

// SetTitleStartsWith adds the titleStartsWith to the get series collection params
func (o *GetSeriesCollectionParams) SetTitleStartsWith(titleStartsWith *string) {
	o.TitleStartsWith = titleStartsWith
}

// WithTs adds the ts to the get series collection params
func (o *GetSeriesCollectionParams) WithTs(ts string) *GetSeriesCollectionParams {
	o.SetTs(ts)
	return o
}

// SetTs adds the ts to the get series collection params
func (o *GetSeriesCollectionParams) SetTs(ts string) {
	o.Ts = ts
}

// WriteToRequest writes these params to a swagger request
func (o *GetSeriesCollectionParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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