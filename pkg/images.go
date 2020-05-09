package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
)

type ImageService struct {
	secret string
	client *http.Client
}

func NewImageService() *ImageService {
	return &ImageService{
		secret: "563492ad6f91700001000001a62e4b6bf7774ad3bac2149b6ff31fe9",
		client: &http.Client{},
	}
}

type imageSearchResults struct {
	Photos []*imageSearchResult `json:"photos"`
}

type imageSearchResult struct {
	URL     string            `json:"url"`
	Sources map[string]string `json:"src"`
}

// RandomImage searches for a random image for a particular subject and returns
// the bytes of it.
func (s *ImageService) RandomImage(subject string) (string, *[]byte, error) {
	randomNumber := rand.Intn(500)
	url := fmt.Sprintf("https://api.pexels.com/v1/search?query=%s&per_page=1&page=%d", url.QueryEscape(subject), randomNumber)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", s.secret)

	parseFormErr := req.ParseForm()
	if parseFormErr != nil {
		return "", nil, parseFormErr
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", nil, err
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	var results imageSearchResults
	if err := json.Unmarshal(respBody, &results); err != nil {
		return "", nil, err
	}

	if len(results.Photos) == 0 {
		fmt.Println(string(respBody))
		return "", nil, fmt.Errorf("No results when searching for photos")
	}

	imageURL := results.Photos[0].Sources["large2x"]
	if imageURL == "" {
		imageURL = results.Photos[0].Sources["original"]
	}

	fmt.Println(imageURL)
	imageReq, err := http.NewRequest("GET", imageURL, nil)

	imageResp, err := s.client.Do(imageReq)
	if err != nil {
		return "", nil, err
	}

	imageBytes, _ := ioutil.ReadAll(imageResp.Body)
	return imageURL, &imageBytes, nil
}
