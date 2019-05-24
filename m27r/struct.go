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
