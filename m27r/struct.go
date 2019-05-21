package m27r

type Character struct {
	Intact bool // indicator if any data missing

	Comics      []int32 // list of comic id
	Description string  // short bio or description
	Events      []int32 // list of event id
	ID          int32
	Modified    string
	Name        string
	Series      []int32 // list of series id
	Stories     []int32 // list of story id
	Thumbnail   string  // url of thumbnail image
	URLs        []*URL  // list of resource urls
}

type Comic struct {
	Intact bool // indicator if any data missing

	Characters         []int32 // list of character id
	CollectedIssues    []int32 // list of comic id
	Collections        []int32 // list of comic id
	Creators           []int32 // list of creator id
	Dates              []*ComicDate
	Description        string
	DigitalID          int32
	EAN                string
	Events             []int32 // list of event id
	Format             string
	ID                 int32
	Images             []string
	ISSN               string
	IssueNumber        float64
	Modified           string
	PageCount          int32
	Prices             []*ComicPrice
	SeriesID           int32
	Stories            []int32 // list of story id
	TextObjects        []*TextObject
	Thumbnail          string // url of thumbnail image
	Title              string
	UPC                string
	URLs               []*URL // list of resource urls
	VariantDescription string
	Variants           []int32 // list of comic id
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
