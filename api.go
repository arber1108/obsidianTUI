package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var store = Store{}

type Store struct {
	apiKey string
	vaults []string
}

type Directory struct {
	Files []string `json:"files"`
	Path  string
}

const basePath = "https://127.0.0.1:27124"

func setApiKey(s string) {
	store.apiKey = s
}

func checkApiKey(apiKey string) (int, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest("GET", basePath+"/vault/", nil)
	if err != nil {
		fmt.Printf("Unable to create request: %v", err)
		os.Exit(1)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	res, err := client.Do(req)
	if err != nil {
		return res.StatusCode, err
	}

	return res.StatusCode, nil

}

func getDirectory(path string) Directory {
	var dir Directory
	dir.Path = path

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	url := basePath + "/vault/"
	if path != "" {
		url += path
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Unable to create request: %v", err)
		os.Exit(1)
	}
	req.Header.Set("Authorization", "Bearer "+store.apiKey)

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Unable to get directory: %v", err)
		os.Exit(1)
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Unable to read response body: %v", err)
		os.Exit(1)
	}

	err = json.Unmarshal(resBody, &dir)
	if err != nil {
		fmt.Printf("Unable to parse response: %v", err)
		os.Exit(1)
	}
	return dir
}

func getFile(path string) string {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest("GET", basePath+"/vault/"+path, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("Authorization", "Bearer "+store.apiKey)
	req.Header.Set("Accept", "text/markdown")

	res, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ""
	}

	return string(body)
}

func (d Directory) isDirectory(index int) bool {
	return strings.HasSuffix(d.Files[index], "/")
}

func (d Directory) parentPath() string {
	p := strings.TrimSuffix(d.Path, "/")
	i := strings.LastIndex(p, "/")
	if i < 0 {
		return ""
	}
	return p[:i+1]
}

func (d Directory) isMarkdown(index int) bool {
	return strings.HasSuffix(d.Files[index], ".md")
}
