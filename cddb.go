package cddb

import (
	"bufio"
	"fmt"
	"github.com/hakkin/cddb/abstract"
	"github.com/hakkin/cddb/gracenote"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	
	"golang.org/x/net/context"
)

func CddbHttp(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	ctx := abstract.GetContext(r)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	var reader io.Reader
	switch r.Method {
	case http.MethodGet:
		reader = strings.NewReader(r.URL.Query().Get("cmd"))
	case http.MethodPost:
		reader = r.Body
	default:
		fmt.Fprint(w, cddbStatus(530, "Unsupported method", true))
		return
	}

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)

	scanner.Scan()
	if scanner.Text() != "cddb" {
		fmt.Fprint(w, cddbStatus(500, "Unknown command", true))
		return
	}

	scanner.Scan()
	command := scanner.Text()
	cmdArray := []string{}
	for scanner.Scan() {
		cmdArray = append(cmdArray, scanner.Text())
	}
	switch command {
	case "query":
		queryCmd, err := createQueryCmd(cmdArray)
		if err != nil {
			log.Println(err)
			fmt.Fprint(w, err)
			return
		}
		for i, v := range path {
			switch i {
			case 1:
				queryCmd.language = v
			case 2:
				queryCmd.country = v
			}
		}
		response, err := Query(ctx, queryCmd)
		if err != nil {
			log.Println(err)
			fmt.Fprint(w, err)
			return
		}
		fmt.Fprint(w, response)
	case "read":
		readCmd, err := createReadCmd(cmdArray)
		if err != nil {
			log.Println(err)
			fmt.Fprint(w, err)
			return
		}
		for i, v := range path {
			switch i {
			case 1:
				readCmd.language = v
			case 2:
				readCmd.country = v
			}
		}
		response, err := Read(ctx, readCmd)
		if err != nil {
			log.Println(err)
			fmt.Fprint(w, err)
			return
		}
		fmt.Fprint(w, response)
	default:
		fmt.Fprint(w, cddbStatus(500, "Unknown command", true))
		return
	}
}

func Query(ctx context.Context, queryCmd QueryCmd) (response string, err error) {
	query := gracenote.Queries{Language: queryCmd.language, Country: queryCmd.country}
	query.Auth = gracenote.Auth{Client: cddbConfig.Client, User: cddbConfig.User}
	query.Query = gracenote.Query{Command: "ALBUM_TOC"}

	var offsetsString = []string{}
	for i := range queryCmd.offsets {
		offset := strconv.Itoa(queryCmd.offsets[i])
		offsetsString = append(offsetsString, offset)
	}
	query.Query.TOC = gracenote.TOC{Offsets: strings.Join(offsetsString, " ")}

	albums, err := gracenote.QueryAlbum(ctx, query)
	if err != nil {
		return "", err
	}

	response, err = queryResponse(albums)
	if err != nil {
		return "", err
	}

	return response, nil
}

func Read(ctx context.Context, readCmd ReadCmd) (response string, err error) {
	query := gracenote.Queries{Language: readCmd.language, Country: readCmd.country}
	query.Auth = gracenote.Auth{Client: cddbConfig.Client, User: cddbConfig.User}
	query.Query = gracenote.Query{Command: "ALBUM_FETCH"}
	query.Query.GN_ID = readCmd.discID

	albums, err := gracenote.QueryAlbum(ctx, query)
	if err != nil {
		return "", err
	}

	response, err = readResponse(albums, readCmd)
	if err != nil {
		return "", nil
	}

	return response, nil
}
