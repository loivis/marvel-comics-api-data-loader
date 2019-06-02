package marvel

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// Character .
type Character struct {
	Comics      *ComicList
	Description string
	Events      *EventList
	ID          int
	Modified    string
	Name        string
	ResourceURI string
	Series      *SeriesList
	Stories     *StoryList
	Thumbnail   *Image
	URLs        []*URL
}

// CharacterList .
type CharacterList struct {
	Available     int
	CollectionURI string
	Items         []*CharacterSummary
	Returned      int
}

// CharacterSummary .
type CharacterSummary struct {
	Name        string
	ResourceURI string
	Role        string
}

// Comic .
type Comic struct {
	Characters         *CharacterList
	CollectedIssues    []*ComicSummary
	Collections        []*ComicSummary
	Creators           *CreatorList
	Dates              []*ComicDate
	Description        string
	DiamondCode        string
	DigitalID          int
	Ean                string
	Events             *EventList
	Format             string
	ID                 int
	Images             []*Image
	ISBN               string
	ISSN               string
	IssueNumber        float64
	Modified           string
	PageCount          int
	Prices             []*ComicPrice
	ResourceURI        string
	Series             *SeriesSummary
	Stories            *StoryList
	TextObjects        []*TextObject
	Thumbnail          *Image
	Title              string
	UPC                string
	URLs               []*URL
	VariantDescription string
	Variants           []*ComicSummary
}

// UnmarshalJSON handles mismatched data type in api response.
func (comic *Comic) UnmarshalJSON(data []byte) error {
	typ := reflect.TypeOf(comic)

	type Alias Comic

	aux := &struct {
		DiamondCode interface{}
		ISBN        interface{}
		*Alias
	}{
		Alias: (*Alias)(comic),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		if v, ok := err.(*json.UnmarshalTypeError); ok {
			v.Struct = typ.String()
		}
		return err
	}

	switch v := aux.DiamondCode.(type) {
	case nil:
	case string:
		comic.DiamondCode = v
	case int:
		comic.DiamondCode = strconv.FormatInt(int64(v), 10)
	case float64:
		comic.DiamondCode = strconv.FormatFloat(float64(v), 'f', -1, 64)
	default:
		return &json.UnmarshalTypeError{
			Value:  fmt.Sprintf("%T", v),
			Type:   typ,
			Struct: typ.String(),
			Field:  "DiamondCode",
		}
	}

	switch v := aux.ISBN.(type) {
	case nil:
	case string:
		comic.ISBN = v
	case int:
		comic.ISBN = strconv.FormatInt(int64(v), 10)
	case float64:
		comic.ISBN = strconv.FormatFloat(float64(v), 'f', -1, 64)
	default:
		return &json.UnmarshalTypeError{
			Value:  fmt.Sprintf("%T", v),
			Type:   typ,
			Struct: typ.String(),
			Field:  "ISBN",
		}
	}

	return nil
}

// ComicDate .
type ComicDate struct {
	Date string
	Type string
}

// ComicList .
type ComicList struct {
	Available     int
	CollectionURI string
	Items         []*ComicSummary
	Returned      int
}

// ComicPrice .
type ComicPrice struct {
	Price float32
	Type  string
}

// ComicSummary .
type ComicSummary struct {
	Name        string
	ResourceURI string
}

// Creator .
type Creator struct {
	Comics      *ComicList
	Events      *EventList
	FirstName   string
	FullName    string
	ID          int
	LastName    string
	MiddleName  string
	Modified    string
	ResourceURI string
	Series      *SeriesList
	Stories     *StoryList
	Suffix      string
	Thumbnail   *Image
	URLs        []*URL
}

// UnmarshalJSON handles mismatched data type in api response.
func (creator *Creator) UnmarshalJSON(data []byte) error {
	typ := reflect.TypeOf(creator)

	type Alias Creator

	aux := &struct {
		LastName interface{}
		Suffix   interface{}
		*Alias
	}{
		Alias: (*Alias)(creator),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		if v, ok := err.(*json.UnmarshalTypeError); ok {
			v.Struct = typ.String()
		}
		return err
	}

	switch v := aux.LastName.(type) {
	case nil:
	case string:
		creator.LastName = v
	case int:
		creator.LastName = strconv.FormatInt(int64(v), 10)
	case float64:
		creator.LastName = strconv.FormatFloat(float64(v), 'f', -1, 64)
	default:
		return &json.UnmarshalTypeError{
			Value:  fmt.Sprintf("%T", v),
			Type:   typ,
			Struct: typ.String(),
			Field:  "LastName",
		}
	}

	switch v := aux.Suffix.(type) {
	case nil:
	case string:
		creator.Suffix = v
	case int:
		creator.Suffix = strconv.FormatInt(int64(v), 10)
	case float64:
		creator.Suffix = strconv.FormatFloat(float64(v), 'f', -1, 64)
	default:
		return &json.UnmarshalTypeError{
			Value:  fmt.Sprintf("%T", v),
			Type:   typ,
			Struct: typ.String(),
			Field:  "Suffix",
		}
	}

	return nil
}

// CreatorList .
type CreatorList struct {
	Available     int
	CollectionURI string
	Items         []*CreatorSummary
	Returned      int
}

// CreatorSummary .
type CreatorSummary struct {
	Name        string
	ResourceURI string
	Role        string
}

// Event .
type Event struct {
	Characters  *CharacterList
	Comics      *ComicList
	Creators    *CreatorList
	Description string
	End         string
	ID          int
	Modified    string
	Next        *EventSummary
	Previous    *EventSummary
	ResourceURI string
	Series      *SeriesList
	Start       string
	Stories     *StoryList
	Thumbnail   *Image
	Title       string
	URLs        []*URL
}

// EventList .
type EventList struct {
	Available     int
	CollectionURI string
	Items         []*EventSummary
	Returned      int
}

// EventSummary .
type EventSummary struct {
	Name        string
	ResourceURI string
}

// Series .
type Series struct {
	Characters  *CharacterList
	Comics      *ComicList
	Creators    *CreatorList
	Description string
	EndYear     int
	Events      *EventList
	ID          int
	Modified    string
	Next        *SeriesSummary
	Previous    *SeriesSummary
	Rating      string
	ResourceURI string
	StartYear   int
	Stories     *StoryList
	Thumbnail   *Image
	Title       string
	URLs        []*URL
}

// SeriesList .
type SeriesList struct {
	Available     int
	CollectionURI string
	Items         []*SeriesSummary
	Returned      int
}

// SeriesSummary .
type SeriesSummary struct {
	Name        string
	ResourceURI string
}

// Story .
type Story struct {
	Characters    *CharacterList
	Comics        *ComicList
	Creators      *CreatorList
	Description   string
	Events        *EventList
	ID            int
	Modified      string
	Originalissue *ComicSummary
	ResourceURI   string
	Series        *SeriesList
	Thumbnail     *Image
	Title         string
	Type          string
}

// UnmarshalJSON handles mismatched data type in api response.
func (story *Story) UnmarshalJSON(data []byte) error {
	typ := reflect.TypeOf(story)

	type Alias Story

	aux := &struct {
		Description interface{}
		Title       interface{}
		*Alias
	}{
		Alias: (*Alias)(story),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		if v, ok := err.(*json.UnmarshalTypeError); ok {
			v.Struct = typ.String()
		}
		return err
	}

	switch v := aux.Description.(type) {
	case nil:
	case string:
		story.Description = v
	case int:
		story.Description = strconv.FormatInt(int64(v), 10)
	case float64:
		story.Description = strconv.FormatFloat(float64(v), 'f', -1, 64)
	default:
		return &json.UnmarshalTypeError{
			Value:  fmt.Sprintf("%T", v),
			Type:   typ,
			Struct: typ.String(),
			Field:  "Description",
		}
	}

	switch v := aux.Title.(type) {
	case nil:
	case string:
		story.Title = v
	case int:
		story.Title = strconv.FormatInt(int64(v), 10)
	case float64:
		story.Title = strconv.FormatFloat(float64(v), 'f', -1, 64)
	default:
		return &json.UnmarshalTypeError{
			Value:  fmt.Sprintf("%T", v),
			Type:   typ,
			Struct: typ.String(),
			Field:  "Title",
		}
	}

	return nil
}

// StoryList .
type StoryList struct {
	Available     int
	CollectionURI string
	Items         []*StorySummary
	Returned      int
}

// StorySummary .
type StorySummary struct {
	Name        string
	ResourceURI string
	Type        string
}

// Image .
type Image struct {
	Extension string
	Path      string
}

// TextObject .
type TextObject struct {
	Language string
	Text     string
	Type     string
}

// URL .
type URL struct {
	Type string
	URL  string
}
