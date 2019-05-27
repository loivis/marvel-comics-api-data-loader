package m27r

type Character struct {
	Intact bool `bson:"intact"` // indicator if any data missing

	Comics      []int32 `bson:"comics"`      // list of comic id
	Description string  `bson:"description"` // short bio or description
	Events      []int32 `bson:"events"`      // list of event id
	ID          int32   `bson:"id"`
	Modified    string  `bson:"modified"`
	Name        string  `bson:"name"`
	Series      []int32 `bson:"series"`    // list of series id
	Stories     []int32 `bson:"stories"`   // list of story id
	Thumbnail   string  `bson:"thumbnail"` // url of thumbnail image
	URLs        []*URL  `bson:"urls"`      // list of resource urls
}

func (char *Character) Identify() int32 {
	return char.ID
}

type Comic struct {
	Intact bool `bson:"intact"` // indicator if any data missing

	Characters         []int32       `bson:"characters"`       // list of character id
	CollectedIssues    []int32       `bson:"collected_issues"` // list of comic id
	Collections        []int32       `bson:"collections"`      // list of comic id
	Creators           []int32       `bson:"creators"`         // list of creator id
	Dates              []*ComicDate  `bson:"dates"`
	Description        string        `bson:"description"`
	DigitalID          int32         `bson:"digital_id"`
	EAN                string        `bson:"ean"`
	Events             []int32       `bson:"events"` // list of event id
	Format             string        `bson:"format"`
	ID                 int32         `bson:"id"`
	Images             []string      `bson:"images"`
	ISSN               string        `bson:"issn"`
	IssueNumber        float64       `bson:"issue_number"`
	Modified           string        `bson:"modified"`
	PageCount          int32         `bson:"page_count"`
	Prices             []*ComicPrice `bson:"prices"`
	SeriesID           int32         `bson:"series_id"`
	Stories            []int32       `bson:"stories"` // list of story id
	TextObjects        []*TextObject `bson:"text_objects"`
	Thumbnail          string        `bson:"thumbnail"` // url of thumbnail image
	Title              string        `bson:"title"`
	UPC                string        `bson:"upc"`
	URLs               []*URL        `bson:"urls"` // list of resource urls
	VariantDescription string        `bson:"variant_description"`
	Variants           []int32       `bson:"variants"` // list of comic id
}

func (comic *Comic) Identify() int32 {
	return comic.ID
}

type Creator struct {
	Intact bool `bson:"intact"` // indicator if any data missing

	Comics     []int32 `bson:"comics"` // list of comic id
	Events     []int32 `bson:"events"` // list of event id
	FirtName   string  `bson:"firt_name"`
	FullName   string  `bson:"full_name"`
	ID         int32   `bson:"id"`
	LastName   string  `bson:"last_name"`
	MiddleName string  `bson:"middle_name"`
	Modified   string  `bson:"modified"`
	Series     []int32 `bson:"series"`  // list of series id
	Stories    []int32 `bson:"stories"` // list of story id
	Suffix     string  `bson:"suffix"`
	Thumbnail  string  `bson:"thumbnail"` // url of thumbnail image
	URLs       []*URL  `bson:"urls"`      // list of resource urls
}

func (creator *Creator) Identify() int32 {
	return creator.ID
}

type Event struct {
	Intact bool `bson:"intact"` // indicator if any data missing

	Characters  []int32 `bson:"characters"` // list of character id
	Comics      []int32 `bson:"comics"`     // list of comic id
	Creators    []int32 `bson:"creators"`   // list of creator id
	Description string  `bson:"description"`
	End         string  `bson:"end"` // The date of publication of the last issue in this event.
	ID          int32   `bson:"id"`
	Modified    string  `bson:"modified"`
	Next        int32   `bson:"next"`      // id of the event which follows this event
	Previous    int32   `bson:"previous"`  // id of the event which preceded this event
	Series      []int32 `bson:"series"`    // list of series id
	Start       string  `bson:"start"`     // The date of publication of the first issue in this event.
	Stories     []int32 `bson:"stories"`   // list of story id
	Thumbnail   string  `bson:"thumbnail"` // url of thumbnail image
	Title       string  `bson:"title"`
	URLs        []*URL  `bson:"urls"` // list of resource urls
}

func (event *Event) Identify() int32 {
	return event.ID
}

type Series struct {
	Intact bool `bson:"intact"` // indicator if any data missing

	Characters  []int32 `bson:"characters"` // list of character id
	Comics      []int32 `bson:"comics"`     // list of comic id
	Creators    []int32 `bson:"creators"`   // list of creator id
	Description string  `bson:"description"`
	EndYear     int32   `bson:"end_year"` // The date of publication of the series.
	Events      []int32 `bson:"events"`   // list of event id
	ID          int32   `bson:"id"`
	Modified    string  `bson:"modified"`
	Next        int32   `bson:"next"`     // id of the series which follows this series
	Previous    int32   `bson:"previous"` // id of the series which preceded this series
	Rating      string  `bson:"rating"`
	StartYear   int32   `bson:"start_year"` // The date of publication of the series.
	Stories     []int32 `bson:"stories"`    // list of story id
	Thumbnail   string  `bson:"thumbnail"`  // url of thumbnail image
	Title       string  `bson:"title"`
	URLs        []*URL  `bson:"urls"` // list of resource urls
}

func (series *Series) Identify() int32 {
	return series.ID
}

type Story struct {
	Intact bool `bson:"intact"` // indicator if any data missing

	Characters    []int32 `bson:"characters"` // list of character id
	Comics        []int32 `bson:"comics"`     // list of comic id
	Creators      []int32 `bson:"creators"`   // list of creator id
	Description   string  `bson:"description"`
	Events        []int32 `bson:"events"` // list of event id
	ID            int32   `bson:"id"`
	Modified      string  `bson:"modified"`
	Originalissue int32   `json:"original_issue"` // comic id
	Series        []int32 `bson:"series"`         // list of series id
	Thumbnail     string  `bson:"thumbnail"`      // url of thumbnail image
	Title         string  `bson:"title"`
	Type          string  `bson:"type"`
}

func (story *Story) Identify() int32 {
	return story.ID
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
