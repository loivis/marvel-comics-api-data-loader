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

type URL struct {
	Type string
	URL  string
}
