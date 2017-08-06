package cddb

import (
	"bytes"
	"fmt"
	"github.com/Hakkin/cddb/gracenote"
	"text/template"
)

var queryTemplate, readTemplate *template.Template

func init() {
	const queryTemplateString = `{{/**/ -}}
211 Found inexact matches, list follows (until terminating marker)
{{range . -}}
Misc {{.GN_ID}} {{with .Artist}}{{.}} / {{end}}{{.Title}}
{{end -}}
.
`
	const readTemplateString = `{{/**/ -}}
210 OK, CDDB database entry follows (until terminating marker)
DISCID={{.GN_ID}}
DTITLE={{with .Artist}}{{.}} / {{end}}{{with .Title}}{{.}}{{end}}
DYEAR={{with .Date}}{{.}}{{end}}
DGENRE={{with .Genre}}{{.}}{{end}}
{{range $index, $track := .Tracks -}}
TTITLE{{$index}}={{with $track.Artist}}{{.}} / {{end}}{{with $track.Title}}{{.}}{{end}}
{{end -}}
.
`
	queryTemplate = template.Must(template.New("queryTemplate").Parse(queryTemplateString))
	readTemplate = template.Must(template.New("readTemplate").Parse(readTemplateString))
}

func cddbStatus(errorCode int, errorMessage string, endResponse bool) string {
	var endCharacter string
	if endResponse {
		endCharacter = "."
	}
	return fmt.Sprintf("%v %v%v", errorCode, errorMessage, endCharacter)
}

func queryResponse(albums []gracenote.Album) (response string, err error) {
	if len(albums) == 0 {
		return cddbStatus(202, "No match found", true), nil
	}

	responseBuffer := &bytes.Buffer{}

	err = queryTemplate.Execute(responseBuffer, albums)
	if err != nil {
		return "", err
	}

	return responseBuffer.String(), nil
}

func readResponse(albums []gracenote.Album, readCmd ReadCmd) (response string, err error) {
	if len(albums) == 0 {
		return cddbStatus(401, readCmd.category+" "+readCmd.discID+" Specified CDDB entry not found", true), nil
	}

	album := albums[0]

	responseBuffer := &bytes.Buffer{}

	err = readTemplate.Execute(responseBuffer, album)
	if err != nil {
		return "", err
	}

	return responseBuffer.String(), nil
}
