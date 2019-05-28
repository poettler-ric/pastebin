package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	// URL to the github api
	githubAPIURL = "https://api.github.com"
	// content type for the json request
	jsonContentType = "application/json"
	// location for the gist api on api.github.com
	gistLocation = "/gists"
)

// CreateData holds the information to create a gist
type CreateData struct {
	Description string              `json:"description"`
	Public      bool                `json:"public"`
	Files       map[string]FileData `json:"files"`
}

// FileData holds the information about file entries in CreateData
type FileData struct {
	Content string `json:"content"`
}

// CreateResponse holds response information after creating a gist
type CreateResponse struct {
	ID      string `json:"id"`
	HTMLURL string `json:"html_url"`
}

func main() {
	var data CreateData
	data.Public = true
	data.Files = make(map[string]FileData)

	if len(os.Args) <= 1 {
		// post stdin
		content, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("error while reading from stdin: %s", err)
		}
		data.Files["stdin"] = FileData{string(content)}
	} else {
		// post given files
		for _, a := range os.Args[1:] {
			content, err := ioutil.ReadFile(a)
			if err != nil {
				log.Fatalf("error while reading %s: %s", a, err)
			}
			data.Files[filepath.Base(a)] = FileData{string(content)}
		}
	}

	// post gist
	out, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("error while marshalling: %s", err)
	}

	res, err := http.Post(githubAPIURL+gistLocation, jsonContentType, bytes.NewReader(out))
	if err != nil {
		log.Fatalf("error while posting: %s", err)
	}
	defer res.Body.Close()
	if err != nil {
		log.Fatalf("error while reading the response body: %s", err)
	}

	body, err := ioutil.ReadAll(res.Body)

	if res.StatusCode == 201 {
		// print gist url
		var creationInfo CreateResponse
		err = json.Unmarshal(body, &creationInfo)
		if err != nil {
			log.Fatalf("error while parsing json: %s", err)
		}
		fmt.Printf("created gist: %s\n", creationInfo.HTMLURL)
	} else {
		// print error message
		fmt.Printf("failed to create gist:\n%s\n%s\n", res.Status, body)
	}
}
