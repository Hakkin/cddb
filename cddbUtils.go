package cddb

import (
	"fmt"
	"github.com/hakkin/cddb/gracenote"
	"log"
	"strconv"
	"strings"
)

func cddbStatus(errorCode int, errorMessage string, endResponse bool) string {
	var endCharacter string
	if endResponse {
		endCharacter = "."
	}
	return fmt.Sprintf("%v %v%v", errorCode, errorMessage, endCharacter)
}

func queryResponse(albums []gracenote.Album) (response string, err error) {
	if len(albums) == 0 {
		log.Println("query no match found")
		return cddbStatus(202, "No match found", true), nil
	}

	responseString := []string{}
	responseString = append(responseString, cddbStatus(211, "Found inexact matches, list follows (until terminating marker)", false))

	for _, album := range albums {
		responseString = append(responseString, "Misc "+album.GN_ID+" "+album.Artist+" / "+album.Title)
	}

	responseString = append(responseString, ".")

	return strings.Join(responseString, "\r\n"), nil
}

func readResponse(albums []gracenote.Album, readCmd ReadCmd) (response string, err error) {
	if len(albums) == 0 {
		log.Println("read no match found")
		return cddbStatus(401, readCmd.category+" "+readCmd.discID+" Specified CDDB entry not found", true), nil
	}

	album := albums[0]

	responseString := []string{}
	responseString = append(responseString, cddbStatus(210, readCmd.category+" "+album.GN_ID+" CD database entry follows (until terminating `.')", false))
	responseString = append(responseString, "DISCID="+album.GN_ID)
	responseString = append(responseString, "DTITLE="+album.Artist+" / "+album.Title)
	responseString = append(responseString, "DYEAR="+strconv.Itoa(album.Date))
	responseString = append(responseString, "DGENRE="+album.Genre)

	for i, v := range album.Tracks {
		title := v.Title
		if v.Artist != "" {
			title = v.Artist + " / " + title
		}

		responseString = append(responseString, "TTITLE"+strconv.Itoa(i)+"="+title)
	}

	responseString = append(responseString, "EXTD=")

	for i := range album.Tracks {
		responseString = append(responseString, "EXTT"+strconv.Itoa(i)+"=")
	}

	responseString = append(responseString, "PLAYORDER=")

	responseString = append(responseString, ".\r\n")

	return strings.Join(responseString, "\r\n"), nil
}
