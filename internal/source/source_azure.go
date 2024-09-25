package source

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/Escape-Technologies/cloudfinder/internal/log"
	"github.com/Escape-Technologies/cloudfinder/pkg/provider"
)

type Azure struct{}

type azureJSON struct {
	Values []struct {
		Properties struct {
			AddressPrefixes []string `json:"addressPrefixes"`
		} `json:"properties"`
	} `json:"values"`
}

/*
- Gets https://azservicetags.azurewebsites.net
- Extracts urls https://download.microsoft.com/download/**\/*.json from html
*/
func getAzureFileUrls() ([]string, error) {
	url := "https://azservicetags.azurewebsites.net"

	ctx, cancel := context.WithTimeout(context.Background(), defaultHTTPTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	file, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, errors.New("response is nil")
	}
	if file.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code is %d", file.StatusCode)
	}
	defer file.Body.Close()

	bodyText, err := io.ReadAll(file.Body)
	if err != nil {
		return nil, err
	}

	fileRegex := regexp.MustCompile(`"(https:\/\/download\.microsoft\.com\/download\/.*\.json)"`)

	matches := fileRegex.FindAllStringSubmatch(string(bodyText), -1)

	urls := make([]string, 0)

	for _, match := range matches {
		urls = append(urls, match[1])
	}

	return urls, nil
}

func (a Azure) GetIPRanges() []*IPRange {
	log.Info("Fetching Azure ip ranges from %s", awsFileURL)
	prefixes := make([]string, 0)
	urls, err := getAzureFileUrls()
	if err != nil {
		panic(err)
	}
	for _, url := range urls {
		var azureJSON *azureJSON
		err := LoadFileURLToJSON(url, &azureJSON)
		if err != nil {
			log.Warning("Failed to load file url to json for Azure", err)
			continue
		}
		for _, value := range azureJSON.Values {
			prefixes = append(prefixes, value.Properties.AddressPrefixes...)
		}
	}

	ranges := make([]*IPRange, 0)
	for _, prefix := range prefixes {
		network, cat := ParseCIDR(prefix)
		ranges = append(ranges, &IPRange{
			Network:  network,
			Cat:      cat,
			Provider: provider.Azure,
		})
	}

	return ranges
}
