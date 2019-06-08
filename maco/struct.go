package maco

type Character struct {
	Intact bool `bson:"intact"` // indicator if any data missing

	Comics      []int  `bson:"comics"`                // list of comic id
	Description string `bson:"description,omitempty"` // short bio or description
	Events      []int  `bson:"events"`                // list of event id
	ID          int    `bson:"id"`
	Modified    string `bson:"modified,omitempty"`
	Name        string `bson:"name,omitempty"`
	Series      []int  `bson:"series"`              // list of series id
	Stories     []int  `bson:"stories"`             // list of story id
	Thumbnail   string `bson:"thumbnail,omitempty"` // url of thumbnail image
	URLs        []*URL `bson:"urls"`                // list of resource urls
}

func (char *Character) Identify() int {
	return int(char.ID)
}

type Comic struct {
	Intact bool `bson:"intact"` // indicator if any data missing

	Characters         []int         `bson:"characters"`       // list of character id
	CollectedIssues    []int         `bson:"collected_issues"` // list of comic id
	Collections        []int         `bson:"collections"`      // list of comic id
	Creators           []int         `bson:"creators"`         // list of creator id
	Dates              []*ComicDate  `bson:"dates"`
	Description        string        `bson:"description,omitempty"`
	DiamondCode        string        `bson:"diamond_code"`
	DigitalID          int           `bson:"digital_id"`
	EAN                string        `bson:"ean,omitempty"`
	Events             []int         `bson:"events"` // list of event id
	Format             string        `bson:"format,omitempty"`
	ID                 int           `bson:"id"`
	Images             []string      `bson:"images"`
	ISBN               string        `bson:"isbn,omitempty"`
	ISSN               string        `bson:"issn,omitempty"`
	IssueNumber        float64       `bson:"issue_number"`
	Modified           string        `bson:"modified,omitempty"`
	PageCount          int           `bson:"page_count"`
	Prices             []*ComicPrice `bson:"prices"`
	SeriesID           int           `bson:"series_id"`
	Stories            []int         `bson:"stories"` // list of story id
	TextObjects        []*TextObject `bson:"text_objects"`
	Thumbnail          string        `bson:"thumbnail,omitempty"` // url of thumbnail image
	Title              string        `bson:"title,omitempty"`
	UPC                string        `bson:"upc,omitempty"`
	URLs               []*URL        `bson:"urls"` // list of resource urls
	VariantDescription string        `bson:"variant_description,omitempty"`
	Variants           []int         `bson:"variants"` // list of comic id
}

func (comic *Comic) Identify() int {
	return int(comic.ID)
}

type Creator struct {
	Intact bool `bson:"intact"` // indicator if any data missing

	Comics     []int  `bson:"comics"` // list of comic id
	Events     []int  `bson:"events"` // list of event id
	FirstName  string `bson:"first_name,omitempty"`
	FullName   string `bson:"full_name,omitempty"`
	ID         int    `bson:"id"`
	LastName   string `bson:"last_name,omitempty"`
	MiddleName string `bson:"middle_name,omitempty"`
	Modified   string `bson:"modified,omitempty"`
	Series     []int  `bson:"series"`  // list of series id
	Stories    []int  `bson:"stories"` // list of story id
	Suffix     string `bson:"suffix,omitempty"`
	Thumbnail  string `bson:"thumbnail,omitempty"` // url of thumbnail image
	URLs       []*URL `bson:"urls"`                // list of resource urls
}

func (creator *Creator) Identify() int {
	return int(creator.ID)
}

type Event struct {
	Intact bool `bson:"intact"` // indicator if any data missing

	Characters  []int  `bson:"characters"` // list of character id
	Comics      []int  `bson:"comics"`     // list of comic id
	Creators    []int  `bson:"creators"`   // list of creator id
	Description string `bson:"description,omitempty"`
	End         string `bson:"end,omitempty"` // The date of publication of the last issue in this event.
	ID          int    `bson:"id"`
	Modified    string `bson:"modified,omitempty"`
	Next        int    `bson:"next"`                // id of the event which follows this event
	Previous    int    `bson:"previous"`            // id of the event which preceded this event
	Series      []int  `bson:"series"`              // list of series id
	Start       string `bson:"start,omitempty"`     // The date of publication of the first issue in this event.
	Stories     []int  `bson:"stories"`             // list of story id
	Thumbnail   string `bson:"thumbnail,omitempty"` // url of thumbnail image
	Title       string `bson:"title,omitempty"`
	URLs        []*URL `bson:"urls"` // list of resource urls
}

func (event *Event) Identify() int {
	return int(event.ID)
}

type Series struct {
	Intact bool `bson:"intact"` // indicator if any data missing

	Characters  []int  `bson:"characters"` // list of character id
	Comics      []int  `bson:"comics"`     // list of comic id
	Creators    []int  `bson:"creators"`   // list of creator id
	Description string `bson:"description,omitempty"`
	EndYear     int    `bson:"end_year"` // The date of publication of the series.
	Events      []int  `bson:"events"`   // list of event id
	ID          int    `bson:"id"`
	Modified    string `bson:"modified,omitempty"`
	Next        int    `bson:"next"`     // id of the series which follows this series
	Previous    int    `bson:"previous"` // id of the series which preceded this series
	Rating      string `bson:"rating,omitempty"`
	StartYear   int    `bson:"start_year"`          // The date of publication of the series.
	Stories     []int  `bson:"stories"`             // list of story id
	Thumbnail   string `bson:"thumbnail,omitempty"` // url of thumbnail image
	Title       string `bson:"title,omitempty"`
	URLs        []*URL `bson:"urls"` // list of resource urls
}

func (series *Series) Identify() int {
	return int(series.ID)
}

type Story struct {
	Intact bool `bson:"intact"` // indicator if any data missing

	Characters    []int  `bson:"characters"` // list of character id
	Comics        []int  `bson:"comics"`     // list of comic id
	Creators      []int  `bson:"creators"`   // list of creator id
	Description   string `bson:"description,omitempty"`
	Events        []int  `bson:"events"` // list of event id
	ID            int    `bson:"id"`
	Modified      string `bson:"modified,omitempty"`
	OriginalIssue int    `json:"original_issue"`      // comic id
	Series        []int  `bson:"series"`              // list of series id
	Thumbnail     string `bson:"thumbnail,omitempty"` // url of thumbnail image
	Title         string `bson:"title,omitempty"`
	Type          string `bson:"type,omitempty"`
}

func (story *Story) Identify() int {
	return int(story.ID)
}

type ComicDate struct {
	Date string `bson:"date"`
	Type string `bson:"type"`
}

type ComicPrice struct {
	Price float32 `bson:"price"`
	Type  string  `bson:"type"`
}

type TextObject struct {
	Language string `bson:"language"`
	Text     string `bson:"text"`
	Type     string `bson:"type"`
}

type URL struct {
	Type string `bson:"type"`
	URL  string `bson:"url"`
}
