package gracenote

type Track struct {
	TrackNumber int    `xml:"TRACK_NUM"`
	GN_ID       string `xml:"GN_ID"`
	Artist      string `xml:"ARTIST"`
	Title       string `xml:"TITLE"`
}

type Album struct {
	GN_ID      string  `xml:"GN_ID"`
	Artist     string  `xml:"ARTIST"`
	Title      string  `xml:"TITLE"`
	Date       int     `xml:"DATE"`
	Genre      string  `xml:"GENRE"`
	TrackCount int     `xml:"TRACK_COUNT"`
	Tracks     []Track `xml:"TRACK"`
}

type response struct {
	Status string  `xml:"STATUS,attr"`
	Albums []Album `xml:"ALBUM"`
}

type responses struct {
	Message   string     `xml:"MESSAGE"`
	Responses []response `xml:"RESPONSE"`
}

type TOC struct {
	XMLName string `xml:"TOC"`
	Offsets string `xml:"OFFSETS"`
}

type Query struct {
	XMLName string `xml:"QUERY"`
	Command string `xml:"CMD,attr"`
	GN_ID   string `xml:"GN_ID,omitempty"`
	TOC     TOC    `xml:"TOC,omitempty"`
}

type Auth struct {
	XMLName string `xml:"AUTH"`
	Client  string `xml:"CLIENT"`
	User    string `xml:"USER"`
}

type Queries struct {
	XMLName  string `xml:"QUERIES"`
	Auth     Auth   `xml:"AUTH"`
	Language string `xml:"LANG,omitempty"`
	Country  string `xml:"COUNTRY,omitempty"`
	Query    Query  `xml:"QUERY"`
}
