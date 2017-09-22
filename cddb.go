package cddb

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Constants for supported Gracenote commands.
const (
	CommandFetch = "ALBUM_FETCH"
	CommandTOC   = "ALBUM_TOC"
)

// Format string for Gracenote's API endpoint.
//
// The first verb is replaced with the digits of the supplied Client ID that
// precede the hyphen.
const endpointFormat = "https://c%v.web.cddbp.net/webapi/xml/1.0/"

// Auth stores the required authentication strings needed to use
// the Gracenote API.
type Auth struct {
	Client string `xml:"CLIENT"`
	User   string `xml:"USER"`
}

// A client to query Gracenote.
//
// For specifics on the effects of the Language and Country parameters,
// please refer to the follow articles:
// https://developer.gracenote.com/sites/default/files/web/webapi/Content/music-web-api/Setting%20the%20Language%20Preference.html
//
// https://developer.gracenote.com/sites/default/files/web/webapi/Content/music-web-api/Specifying%20a%20Country%20Specific.html
type Client struct {
	// Auth stores the required authentication strings needed to use
	// the Gracenote API.
	//
	// Both Client and User fields of Auth are required, an error will be
	// returned if either are empty.
	Auth Auth `xml:"AUTH"`

	// Language sets the LANG parameter of the Gracenote request.
	//
	// If Language is empty, it will be omitted from the query.
	Language string `xml:"LANG,omitempty"`

	// Country sets the COUNTRY parameter of the Gracenote request.
	//
	// If Country is empty, it will be omitted from the query.
	Country string `xml:"COUNTRY,omitempty"`

	// HTTPClient is the HTTP client used to query the Gracenote API.
	//
	// If HTTPClient is nil, http.DefaultClient is used.
	HTTPClient *http.Client `xml:"-"`
}

// TOC describes the Table of Contents of an Audio CD.
//
// Gracenote expects the TOC to contain the absolute frame offset to
// the start of each track, with the final offset being the total CD
// length in frames. Each offset should be separated by a space.
//
// Gracenote expects the offsets to include the lead-in frames of the CD,
// and will return an error if the TOC contains any offset less than 1.
type TOC struct {
	Offsets string `xml:"OFFSETS,omitempty"`
}

// Gracenote query
//
// Command field is required, an error will be returned if it is empty.
type Query struct {
	Command string `xml:"CMD,attr"`
	ID      string `xml:"GN_ID,omitempty"`
	TOC     *TOC   `xml:"TOC,omitempty"`
}

type queries struct {
	XMLName string `xml:"QUERIES"`
	*Client
	Query *Query `xml:"QUERY"`
}

type Track struct {
	Number int    `xml:"TRACK_NUM"`
	ID     string `xml:"GN_ID"`
	Artist string `xml:"ARTIST"`
	Title  string `xml:"TITLE"`
}

type Album struct {
	ID         string  `xml:"GN_ID"`
	Artist     string  `xml:"ARTIST"`
	Title      string  `xml:"TITLE"`
	Date       int     `xml:"DATE"`
	Genre      string  `xml:"GENRE"`
	TrackCount int     `xml:"TRACK_COUNT"`
	Tracks     []Track `xml:"TRACK"`
}

type Response struct {
	Status string  `xml:"STATUS,attr"`
	Albums []Album `xml:"ALBUM"`
}

type responses struct {
	Message  string   `xml:"MESSAGE"`
	Response Response `xml:"RESPONSE"`
}

func (q *Query) SetCommand(command string) {
	q.Command = command
}

func (q *Query) SetID(id string) {
	q.ID = id
}

func (q *Query) SetTOC(offsets string) {
	if q.TOC == nil {
		q.TOC = &TOC{}
	}
	q.TOC.Offsets = offsets
}

func (c *Client) SetAuth(client string, user string) {
	c.Auth.Client = client
	c.Auth.User = user
}

func (c *Client) SetLanguage(language string) {
	c.Language = language
}

func (c *Client) SetCountry(country string) {
	c.Country = country
}

func (c *Client) SetHTTPClient(httpClient *http.Client) {
	c.HTTPClient = httpClient
}

func (c *Client) Do(q *Query) (*Response, error) {
	if c.Auth.Client == "" {
		return nil, errors.New("cddb: Missing Client ID")
	}
	if c.Auth.User == "" {
		return nil, errors.New("cddb: Missing User ID")
	}
	if q.Command == "" {
		return nil, errors.New("cddb: Missing command")
	}

	query := &queries{Client: c, Query: q}

	buffer := &bytes.Buffer{}
	err := xml.NewEncoder(buffer).Encode(query)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(endpointFormat, strings.Split(c.Auth.Client, "-")[0])
	httpClient := c.httpClient()
	resp, err := httpClient.Post(endpoint, "application/xml", buffer)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cddb: Gracenote returned status code %v", resp.StatusCode)
	}

	var r responses
	err = xml.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, err
	}

	if r.Response.Status == "ERROR" {
		return nil, fmt.Errorf("cddb: Gracenote returned error: %v", r.Message)
	}

	return &r.Response, nil
}

func (c *Client) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return http.DefaultClient
}
