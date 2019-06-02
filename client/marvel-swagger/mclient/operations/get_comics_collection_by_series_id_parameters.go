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

// NewGetComicsCollectionBySeriesIDParams creates a new GetComicsCollectionBySeriesIDParams object
// with the default values initialized.
func NewGetComicsCollectionBySeriesIDParams() *GetComicsCollectionBySeriesIDParams {
	var ()
	return &GetComicsCollectionBySeriesIDParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetComicsCollectionBySeriesIDParamsWithTimeout creates a new GetComicsCollectionBySeriesIDParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetComicsCollectionBySeriesIDParamsWithTimeout(timeout time.Duration) *GetComicsCollectionBySeriesIDParams {
	var ()
	return &GetComicsCollectionBySeriesIDParams{

		timeout: timeout,
	}
}

// NewGetComicsCollectionBySeriesIDParamsWithContext creates a new GetComicsCollectionBySeriesIDParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetComicsCollectionBySeriesIDParamsWithContext(ctx context.Context) *GetComicsCollectionBySeriesIDParams {
	var ()
	return &GetComicsCollectionBySeriesIDParams{

		Context: ctx,
	}
}

// NewGetComicsCollectionBySeriesIDParamsWithHTTPClient creates a new GetComicsCollectionBySeriesIDParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetComicsCollectionBySeriesIDParamsWithHTTPClient(client *http.Client) *GetComicsCollectionBySeriesIDParams {
	var ()
	return &GetComicsCollectionBySeriesIDParams{
		HTTPClient: client,
	}
}

/*GetComicsCollectionBySeriesIDParams contains all the parameters to send to the API endpoint
for the get comics collection by series Id operation typically these are written to a http.Request
*/
type GetComicsCollectionBySeriesIDParams struct {

	/*Apikey
	  [Auth] public apikey

	*/
	Apikey string
	/*Characters
	  Return only comics which feature the specified characters (accepts a comma-separated list of ids).

	*/
	Characters []int32
	/*Collaborators
	  Return only comics in which the specified creators worked together (for example in which BOTH Stan Lee and Jack Kirby did work).

	*/
	Collaborators []int32
	/*Creators
	  Return only comics which feature work by the specified creators (accepts a comma-separated list of ids).

	*/
	Creators []int32
	/*DateDescriptor
	  Return comics within a predefined date range.

	*/
	DateDescriptor []string
	/*DateRange
	  Return comics within a predefined date range.  Dates must be specified as date1,date2 (e.g. 2013-01-01,2013-01-02).  Dates are preferably formatted as YYYY-MM-DD but may be sent as any common date format.

	*/
	DateRange []int32
	/*DiamondCode
	  Filter by diamond code.

	*/
	DiamondCode *string
	/*DigitalID
	  Filter by digital comic id.

	*/
	DigitalID *int32
	/*Ean
	  Filter by EAN.

	*/
	Ean *string
	/*Events
	  Return only comics which take place in the specified events (accepts a comma-separated list of ids).

	*/
	Events []int32
	/*Format
	  Filter by the issue format (e.g. comic, digital comic, hardcover).

	*/
	Format *string
	/*FormatType
	  Filter by the issue format type (comic or collection).

	*/
	FormatType *string
	/*HasDigitalIssue
	  Include only results which are available digitally.

	*/
	HasDigitalIssue []bool
	/*Hash
	  [Auth] md5 digest of concatenation of ts, private key, public key

	*/
	Hash string
	/*Isbn
	  Filter by ISBN.

	*/
	Isbn *string
	/*Issn
	  Filter by ISSN.

	*/
	Issn *string
	/*IssueNumber
	  Return only issues in series whose issue number matches the input.

	*/
	IssueNumber *int32
	/*Limit
	  Limit the result set to the specified number of resources.

	*/
	Limit *int32
	/*ModifiedSince
	  Return only comics which have been modified since the specified date.

	*/
	ModifiedSince *strfmt.Date
	/*NoVariants
	  Exclude variant comics from the result set.

	*/
	NoVariants []bool
	/*Offset
	  Skip the specified number of resources in the result set.

	*/
	Offset *int32
	/*OrderBy
	  Order the result set by a field or fields. Add a "-" to the value sort in descending order. Multiple values are given priority in the order in which they are passed.

	*/
	OrderBy []string
	/*SeriesID
	  The series ID.

	*/
	SeriesID int32
	/*SharedAppearances
	  Return only comics in which the specified characters appear together (for example in which BOTH Spider-Man and Wolverine appear).

	*/
	SharedAppearances []int32
	/*StartYear
	  Return only issues in series whose start year matches the input.

	*/
	StartYear *int32
	/*Stories
	  Return only comics which contain the specified stories (accepts a comma-separated list of ids).

	*/
	Stories []int32
	/*Title
	  Return only issues in series whose title matches the input.

	*/
	Title *string
	/*TitleStartsWith
	  Return only issues in series whose title starts with the input.

	*/
	TitleStartsWith *string
	/*Ts
	  [Auth] timestamp

	*/
	Ts string
	/*Upc
	  Filter by UPC.

	*/
	Upc *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithTimeout(timeout time.Duration) *GetComicsCollectionBySeriesIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithContext(ctx context.Context) *GetComicsCollectionBySeriesIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithHTTPClient(client *http.Client) *GetComicsCollectionBySeriesIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithApikey adds the apikey to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithApikey(apikey string) *GetComicsCollectionBySeriesIDParams {
	o.SetApikey(apikey)
	return o
}

// SetApikey adds the apikey to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetApikey(apikey string) {
	o.Apikey = apikey
}

// WithCharacters adds the characters to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithCharacters(characters []int32) *GetComicsCollectionBySeriesIDParams {
	o.SetCharacters(characters)
	return o
}

// SetCharacters adds the characters to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetCharacters(characters []int32) {
	o.Characters = characters
}

// WithCollaborators adds the collaborators to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithCollaborators(collaborators []int32) *GetComicsCollectionBySeriesIDParams {
	o.SetCollaborators(collaborators)
	return o
}

// SetCollaborators adds the collaborators to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetCollaborators(collaborators []int32) {
	o.Collaborators = collaborators
}

// WithCreators adds the creators to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithCreators(creators []int32) *GetComicsCollectionBySeriesIDParams {
	o.SetCreators(creators)
	return o
}

// SetCreators adds the creators to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetCreators(creators []int32) {
	o.Creators = creators
}

// WithDateDescriptor adds the dateDescriptor to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithDateDescriptor(dateDescriptor []string) *GetComicsCollectionBySeriesIDParams {
	o.SetDateDescriptor(dateDescriptor)
	return o
}

// SetDateDescriptor adds the dateDescriptor to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetDateDescriptor(dateDescriptor []string) {
	o.DateDescriptor = dateDescriptor
}

// WithDateRange adds the dateRange to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithDateRange(dateRange []int32) *GetComicsCollectionBySeriesIDParams {
	o.SetDateRange(dateRange)
	return o
}

// SetDateRange adds the dateRange to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetDateRange(dateRange []int32) {
	o.DateRange = dateRange
}

// WithDiamondCode adds the diamondCode to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithDiamondCode(diamondCode *string) *GetComicsCollectionBySeriesIDParams {
	o.SetDiamondCode(diamondCode)
	return o
}

// SetDiamondCode adds the diamondCode to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetDiamondCode(diamondCode *string) {
	o.DiamondCode = diamondCode
}

// WithDigitalID adds the digitalID to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithDigitalID(digitalID *int32) *GetComicsCollectionBySeriesIDParams {
	o.SetDigitalID(digitalID)
	return o
}

// SetDigitalID adds the digitalId to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetDigitalID(digitalID *int32) {
	o.DigitalID = digitalID
}

// WithEan adds the ean to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithEan(ean *string) *GetComicsCollectionBySeriesIDParams {
	o.SetEan(ean)
	return o
}

// SetEan adds the ean to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetEan(ean *string) {
	o.Ean = ean
}

// WithEvents adds the events to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithEvents(events []int32) *GetComicsCollectionBySeriesIDParams {
	o.SetEvents(events)
	return o
}

// SetEvents adds the events to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetEvents(events []int32) {
	o.Events = events
}

// WithFormat adds the format to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithFormat(format *string) *GetComicsCollectionBySeriesIDParams {
	o.SetFormat(format)
	return o
}

// SetFormat adds the format to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetFormat(format *string) {
	o.Format = format
}

// WithFormatType adds the formatType to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithFormatType(formatType *string) *GetComicsCollectionBySeriesIDParams {
	o.SetFormatType(formatType)
	return o
}

// SetFormatType adds the formatType to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetFormatType(formatType *string) {
	o.FormatType = formatType
}

// WithHasDigitalIssue adds the hasDigitalIssue to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithHasDigitalIssue(hasDigitalIssue []bool) *GetComicsCollectionBySeriesIDParams {
	o.SetHasDigitalIssue(hasDigitalIssue)
	return o
}

// SetHasDigitalIssue adds the hasDigitalIssue to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetHasDigitalIssue(hasDigitalIssue []bool) {
	o.HasDigitalIssue = hasDigitalIssue
}

// WithHash adds the hash to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithHash(hash string) *GetComicsCollectionBySeriesIDParams {
	o.SetHash(hash)
	return o
}

// SetHash adds the hash to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetHash(hash string) {
	o.Hash = hash
}

// WithIsbn adds the isbn to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithIsbn(isbn *string) *GetComicsCollectionBySeriesIDParams {
	o.SetIsbn(isbn)
	return o
}

// SetIsbn adds the isbn to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetIsbn(isbn *string) {
	o.Isbn = isbn
}

// WithIssn adds the issn to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithIssn(issn *string) *GetComicsCollectionBySeriesIDParams {
	o.SetIssn(issn)
	return o
}

// SetIssn adds the issn to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetIssn(issn *string) {
	o.Issn = issn
}

// WithIssueNumber adds the issueNumber to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithIssueNumber(issueNumber *int32) *GetComicsCollectionBySeriesIDParams {
	o.SetIssueNumber(issueNumber)
	return o
}

// SetIssueNumber adds the issueNumber to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetIssueNumber(issueNumber *int32) {
	o.IssueNumber = issueNumber
}

// WithLimit adds the limit to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithLimit(limit *int32) *GetComicsCollectionBySeriesIDParams {
	o.SetLimit(limit)
	return o
}

// SetLimit adds the limit to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetLimit(limit *int32) {
	o.Limit = limit
}

// WithModifiedSince adds the modifiedSince to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithModifiedSince(modifiedSince *strfmt.Date) *GetComicsCollectionBySeriesIDParams {
	o.SetModifiedSince(modifiedSince)
	return o
}

// SetModifiedSince adds the modifiedSince to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetModifiedSince(modifiedSince *strfmt.Date) {
	o.ModifiedSince = modifiedSince
}

// WithNoVariants adds the noVariants to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithNoVariants(noVariants []bool) *GetComicsCollectionBySeriesIDParams {
	o.SetNoVariants(noVariants)
	return o
}

// SetNoVariants adds the noVariants to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetNoVariants(noVariants []bool) {
	o.NoVariants = noVariants
}

// WithOffset adds the offset to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithOffset(offset *int32) *GetComicsCollectionBySeriesIDParams {
	o.SetOffset(offset)
	return o
}

// SetOffset adds the offset to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetOffset(offset *int32) {
	o.Offset = offset
}

// WithOrderBy adds the orderBy to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithOrderBy(orderBy []string) *GetComicsCollectionBySeriesIDParams {
	o.SetOrderBy(orderBy)
	return o
}

// SetOrderBy adds the orderBy to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetOrderBy(orderBy []string) {
	o.OrderBy = orderBy
}

// WithSeriesID adds the seriesID to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithSeriesID(seriesID int32) *GetComicsCollectionBySeriesIDParams {
	o.SetSeriesID(seriesID)
	return o
}

// SetSeriesID adds the seriesId to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetSeriesID(seriesID int32) {
	o.SeriesID = seriesID
}

// WithSharedAppearances adds the sharedAppearances to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithSharedAppearances(sharedAppearances []int32) *GetComicsCollectionBySeriesIDParams {
	o.SetSharedAppearances(sharedAppearances)
	return o
}

// SetSharedAppearances adds the sharedAppearances to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetSharedAppearances(sharedAppearances []int32) {
	o.SharedAppearances = sharedAppearances
}

// WithStartYear adds the startYear to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithStartYear(startYear *int32) *GetComicsCollectionBySeriesIDParams {
	o.SetStartYear(startYear)
	return o
}

// SetStartYear adds the startYear to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetStartYear(startYear *int32) {
	o.StartYear = startYear
}

// WithStories adds the stories to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithStories(stories []int32) *GetComicsCollectionBySeriesIDParams {
	o.SetStories(stories)
	return o
}

// SetStories adds the stories to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetStories(stories []int32) {
	o.Stories = stories
}

// WithTitle adds the title to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithTitle(title *string) *GetComicsCollectionBySeriesIDParams {
	o.SetTitle(title)
	return o
}

// SetTitle adds the title to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetTitle(title *string) {
	o.Title = title
}

// WithTitleStartsWith adds the titleStartsWith to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithTitleStartsWith(titleStartsWith *string) *GetComicsCollectionBySeriesIDParams {
	o.SetTitleStartsWith(titleStartsWith)
	return o
}

// SetTitleStartsWith adds the titleStartsWith to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetTitleStartsWith(titleStartsWith *string) {
	o.TitleStartsWith = titleStartsWith
}

// WithTs adds the ts to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithTs(ts string) *GetComicsCollectionBySeriesIDParams {
	o.SetTs(ts)
	return o
}

// SetTs adds the ts to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetTs(ts string) {
	o.Ts = ts
}

// WithUpc adds the upc to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) WithUpc(upc *string) *GetComicsCollectionBySeriesIDParams {
	o.SetUpc(upc)
	return o
}

// SetUpc adds the upc to the get comics collection by series Id params
func (o *GetComicsCollectionBySeriesIDParams) SetUpc(upc *string) {
	o.Upc = upc
}

// WriteToRequest writes these params to a swagger request
func (o *GetComicsCollectionBySeriesIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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

	var valuesCollaborators []string
	for _, v := range o.Collaborators {
		valuesCollaborators = append(valuesCollaborators, swag.FormatInt32(v))
	}

	joinedCollaborators := swag.JoinByFormat(valuesCollaborators, "")
	// query array param collaborators
	if err := r.SetQueryParam("collaborators", joinedCollaborators...); err != nil {
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

	valuesDateDescriptor := o.DateDescriptor

	joinedDateDescriptor := swag.JoinByFormat(valuesDateDescriptor, "")
	// query array param dateDescriptor
	if err := r.SetQueryParam("dateDescriptor", joinedDateDescriptor...); err != nil {
		return err
	}

	var valuesDateRange []string
	for _, v := range o.DateRange {
		valuesDateRange = append(valuesDateRange, swag.FormatInt32(v))
	}

	joinedDateRange := swag.JoinByFormat(valuesDateRange, "")
	// query array param dateRange
	if err := r.SetQueryParam("dateRange", joinedDateRange...); err != nil {
		return err
	}

	if o.DiamondCode != nil {

		// query param diamondCode
		var qrDiamondCode string
		if o.DiamondCode != nil {
			qrDiamondCode = *o.DiamondCode
		}
		qDiamondCode := qrDiamondCode
		if qDiamondCode != "" {
			if err := r.SetQueryParam("diamondCode", qDiamondCode); err != nil {
				return err
			}
		}

	}

	if o.DigitalID != nil {

		// query param digitalId
		var qrDigitalID int32
		if o.DigitalID != nil {
			qrDigitalID = *o.DigitalID
		}
		qDigitalID := swag.FormatInt32(qrDigitalID)
		if qDigitalID != "" {
			if err := r.SetQueryParam("digitalId", qDigitalID); err != nil {
				return err
			}
		}

	}

	if o.Ean != nil {

		// query param ean
		var qrEan string
		if o.Ean != nil {
			qrEan = *o.Ean
		}
		qEan := qrEan
		if qEan != "" {
			if err := r.SetQueryParam("ean", qEan); err != nil {
				return err
			}
		}

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

	if o.Format != nil {

		// query param format
		var qrFormat string
		if o.Format != nil {
			qrFormat = *o.Format
		}
		qFormat := qrFormat
		if qFormat != "" {
			if err := r.SetQueryParam("format", qFormat); err != nil {
				return err
			}
		}

	}

	if o.FormatType != nil {

		// query param formatType
		var qrFormatType string
		if o.FormatType != nil {
			qrFormatType = *o.FormatType
		}
		qFormatType := qrFormatType
		if qFormatType != "" {
			if err := r.SetQueryParam("formatType", qFormatType); err != nil {
				return err
			}
		}

	}

	var valuesHasDigitalIssue []string
	for _, v := range o.HasDigitalIssue {
		valuesHasDigitalIssue = append(valuesHasDigitalIssue, swag.FormatBool(v))
	}

	joinedHasDigitalIssue := swag.JoinByFormat(valuesHasDigitalIssue, "")
	// query array param hasDigitalIssue
	if err := r.SetQueryParam("hasDigitalIssue", joinedHasDigitalIssue...); err != nil {
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

	if o.Isbn != nil {

		// query param isbn
		var qrIsbn string
		if o.Isbn != nil {
			qrIsbn = *o.Isbn
		}
		qIsbn := qrIsbn
		if qIsbn != "" {
			if err := r.SetQueryParam("isbn", qIsbn); err != nil {
				return err
			}
		}

	}

	if o.Issn != nil {

		// query param issn
		var qrIssn string
		if o.Issn != nil {
			qrIssn = *o.Issn
		}
		qIssn := qrIssn
		if qIssn != "" {
			if err := r.SetQueryParam("issn", qIssn); err != nil {
				return err
			}
		}

	}

	if o.IssueNumber != nil {

		// query param issueNumber
		var qrIssueNumber int32
		if o.IssueNumber != nil {
			qrIssueNumber = *o.IssueNumber
		}
		qIssueNumber := swag.FormatInt32(qrIssueNumber)
		if qIssueNumber != "" {
			if err := r.SetQueryParam("issueNumber", qIssueNumber); err != nil {
				return err
			}
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

	var valuesNoVariants []string
	for _, v := range o.NoVariants {
		valuesNoVariants = append(valuesNoVariants, swag.FormatBool(v))
	}

	joinedNoVariants := swag.JoinByFormat(valuesNoVariants, "")
	// query array param noVariants
	if err := r.SetQueryParam("noVariants", joinedNoVariants...); err != nil {
		return err
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

	// path param seriesId
	if err := r.SetPathParam("seriesId", swag.FormatInt32(o.SeriesID)); err != nil {
		return err
	}

	var valuesSharedAppearances []string
	for _, v := range o.SharedAppearances {
		valuesSharedAppearances = append(valuesSharedAppearances, swag.FormatInt32(v))
	}

	joinedSharedAppearances := swag.JoinByFormat(valuesSharedAppearances, "")
	// query array param sharedAppearances
	if err := r.SetQueryParam("sharedAppearances", joinedSharedAppearances...); err != nil {
		return err
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

	if o.Upc != nil {

		// query param upc
		var qrUpc string
		if o.Upc != nil {
			qrUpc = *o.Upc
		}
		qUpc := qrUpc
		if qUpc != "" {
			if err := r.SetQueryParam("upc", qUpc); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}