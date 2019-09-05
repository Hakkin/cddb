package handler

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/Hakkin/cddb/app/cache"
	"github.com/Hakkin/cddb/app/config"
	"github.com/Hakkin/cddb/app/log"

	"github.com/Hakkin/cddb"
)

const (
	// Error strings returned to the end user
	nmErrStr = "202 No match found."
	seErrStr = "402 Server error."
	csErrStr = "500 Command syntax error."
	ucErrStr = "500 Unknown command."
	umErrStr = "530 Unsupported method."
	// Error format strings for logging
	qsErrFStr = "Query syntax error: %v"
	rsErrFStr = "Read syntax error: %v"
	ucErrFStr = "Unknown command: %v"
)

func CDDB(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	logger := log.WithRequest(r)
	rCache := cache.New(r)

	var reader io.Reader
	switch r.Method {
	case http.MethodGet:
		reader = strings.NewReader(r.URL.Query().Get("cmd"))
	case http.MethodPost:
		reader = r.Body
	default:
		fmt.Fprint(w, umErrStr)
		return
	}

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)

	scanner.Scan()
	if scanner.Text() != "cddb" {
		fmt.Fprint(w, ucErrStr)
		return
	}

	scanner.Scan()
	command := scanner.Text()
	cmdArg := []string{}
	for scanner.Scan() {
		cmdArg = append(cmdArg, scanner.Text())
	}

	var query *cddb.Query
	query, err := makeQuery(command, cmdArg...)
	if err != nil {
		logger.Errorf("%v", err)
		fmt.Fprint(w, csErrStr)
		return
	}

	if command == "read" {
		var ok bool
		query.ID, ok = rCache.Get(query.ID)
		if !ok {
			logger.Errorf("%v", fmt.Errorf(rsErrFStr, cmdArg))
			fmt.Fprint(w, csErrStr)
			return
		}
	}

	client := &cddb.Client{}

	client.SetHTTPClient(&http.Client{Timeout: time.Second * 30})

	client.SetAuth(config.Client, config.User)

	path := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	for i, v := range path {
		switch i {
		case 1:
			client.SetLanguage(v)
		case 2:
			client.SetCountry(v)
		}
	}

	resp, err := client.Do(query)
	if err != nil {
		logger.Errorf("%v", err)
		fmt.Fprint(w, seErrStr)
		return
	}

	if resp.Status == "NO_MATCH" {
		logger.Infof("Query returned no matches")
		fmt.Fprint(w, nmErrStr)
		return
	}

	if command == "query" {
		ids := make([]string, len(resp.Albums))
		for i := range resp.Albums {
			ids[i] = resp.Albums[i].ID
		}

		newIDs, err := rCache.Set(ids...)
		if err != nil {
			logger.Errorf("%v", err)
			fmt.Fprint(w, seErrStr)
			return
		}

		for i := range resp.Albums {
			resp.Albums[i].ID = newIDs[i]
		}
	}

	output, err := parseResponse(command, resp)
	if err != nil {
		logger.Errorf("%v", err)
		fmt.Fprint(w, seErrStr)
		return
	}

	switch command {
	case "query":
		logger.Infof("Query returned %v matches", len(resp.Albums))
	case "read":
		logger.Infof("Read returned %v / %v", resp.Albums[0].Artist, resp.Albums[0].Title)
	}

	fmt.Fprint(w, output)
}

func makeQuery(command string, args ...string) (*cddb.Query, error) {
	var qsErr, rsErr, ucErr error = fmt.Errorf(qsErrFStr, args), fmt.Errorf(rsErrFStr, args), fmt.Errorf(ucErrFStr, command)

	switch command {
	case "query":
		// Syntax:
		// <disc id> <track count> <track 1 offset> [... <track n offset>] <total length of CD in seconds>

		// query command requires at least 4 arguments
		if len(args) < 4 {
			return nil, qsErr
		}

		trackCount, err := strconv.Atoi(args[1])
		if err != nil {
			return nil, qsErr
		}

		// Checking to make sure the track count we got is correct, and that the command syntax is correct
		if len(args[2:len(args)-1]) != trackCount {
			return nil, qsErr
		}

		// We add one to the track count here because Gracenote expects the final offset to be the total CD length in frames
		offsets := make([]int, trackCount+1)

		for i := 0; i < trackCount; i++ {
			offsets[i], err = strconv.Atoi(args[i+2])
			if err != nil {
				return nil, qsErr
			}
		}

		totalSeconds, err := strconv.Atoi(args[len(args)-1])
		if err != nil {
			return nil, qsErr
		}

		// Converting from seconds to CDDA frames, which are 1/75th of a second
		offsets[len(offsets)-1] = totalSeconds * 75

		// Gracenote expects the offsets to include the lead-in frames of the CD
		// Most FreeDB clients don't include these frames, and set the first offset as 0,
		// which will make Gracenote return an error, so we add the standard
		// 150 lead-in frames to the offsets if the first offset is 0
		if offsets[0] == 0 {
			for i := range offsets {
				offsets[i] += 150
			}
		}

		query := &cddb.Query{}

		query.SetCommand(cddb.CommandTOC)

		// Converts offsets from int slice to space separated string
		query.SetTOC(strings.Trim(fmt.Sprint(offsets), "[]"))

		return query, nil
	case "read":
		// Syntax
		// <CD category> <disc id>

		// read command requires 2 arguments
		if len(args) != 2 {
			return nil, rsErr
		}

		query := &cddb.Query{}

		query.SetCommand(cddb.CommandFetch)

		query.SetID(args[1])

		return query, nil
	default:
		return nil, ucErr
	}
}

func parseResponse(command string, r *cddb.Response) (string, error) {
	var ucErr error = fmt.Errorf(ucErrFStr, command)

	const queryTemplateString = `{{/**/ -}}
211 Found inexact matches, list follows (until terminating marker)
{{range . -}}
Misc {{.ID}} {{with .Artist}}{{.}} / {{end}}{{.Title}}
{{end -}}
.
`
	const readTemplateString = `{{/**/ -}}
210 OK, CDDB database entry follows (until terminating marker)
# xmcd
#
# Track frame offsets:
{{range $index, $track := .Tracks -}}
# 0
{{end -}}
#
# Disc length: 0 seconds
#
DISCID={{.ID}}
DTITLE={{with .Artist}}{{.}} / {{end}}{{with .Title}}{{.}}{{end}}
DYEAR={{with .Date}}{{.}}{{end}}
DGENRE={{with .Genre}}{{.}}{{end}}
{{range $index, $track := .Tracks -}}
TTITLE{{$index}}={{with $track.Artist}}{{.}} / {{end}}{{with $track.Title}}{{.}}{{end}}
{{end -}}
EXTD=
{{range $index, $track := .Tracks -}}
EXTT{{$index}}=
{{end -}}
PLAYORDER=
.
`

	var respTemplate *template.Template
	var data interface{}
	var err error

	switch command {
	case "query":
		respTemplate, err = template.New("queryTemplate").Parse(queryTemplateString)
		if err != nil {
			return "", err
		}
		data = r.Albums
	case "read":
		respTemplate, err = template.New("readTemplate").Parse(readTemplateString)
		if err != nil {
			return "", err
		}
		if len(r.Albums) < 1 {
			return "", errors.New("Response contains no albums")
		}
		data = r.Albums[0]
	default:
		return "", ucErr
	}

	response := &bytes.Buffer{}

	err = respTemplate.Execute(response, data)
	if err != nil {
		return "", err
	}

	return response.String(), nil
}
