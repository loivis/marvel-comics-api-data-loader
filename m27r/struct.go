package m27r

type Character struct {
	Intact bool `json:"intact"` // indicator if any data missing

	Comics      []int32 `json:"comics"`      // list of comic id
	Description string  `json:"description"` // short bio or description
	Events      []int32 `json:"events"`      // list of event id
	ID          int32   `json:"id"`
	Modified    string  `json:"modified"`
	Name        string  `json:"name"`
	Series      []int32 `json:"series"`    // list of series id
	Stories     []int32 `json:"stories"`   // list of story id
	Thumbnail   string  `json:"thumbnail"` // url of thumbnail image
	URLs        []*URL  `json:"urls"`      // list of resource urls
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

type ComicDate struct {
	Date string
	Type string
}

type ComicPrice struct {
	Price float32
	Type  string
}

type TextObject struct {
	Language string
	Text     string
	Type     string
}

type URL struct {
	Type string
	URL  string
}
