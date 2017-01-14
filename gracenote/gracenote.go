package gracenote

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/hakkin/cddb/abstract"
	"net/http"
	"strings"
	
	"golang.org/x/net/context"
)

func getEndpoint(clientID string) string {
	clientNumber := strings.Split(clientID, "-")[0]
	return fmt.Sprintf("https://c%v.web.cddbp.net/webapi/xml/1.0/", clientNumber)
}

func QueryAlbum(ctx context.Context, query Queries) (albums []Album, err error) {
	endpoint := getEndpoint(query.Auth.Client)

	buffer := &bytes.Buffer{}

	enc := xml.NewEncoder(buffer)

	err = enc.Encode(query)
	if err != nil {
		return nil, err
	}
	
	client := abstract.GetClient(ctx)
	response, err := client.Post(endpoint, "application/xml", buffer)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Gracenote returned %v", response.StatusCode)
	}

	dec := xml.NewDecoder(response.Body)

	var r responses
	err = dec.Decode(&r)
	if err != nil {
		return nil, err
	}

	if r.Message != "" && r.Responses[0].Status != "OK" {
		return nil, fmt.Errorf("%v: %v", r.Responses[0].Status, r.Message)
	}

	return r.Responses[0].Albums, nil
}
