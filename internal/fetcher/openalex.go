package fetcher

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (APIFetcher *APIFetcher) Fetch(ctx context.Context, fetchReq FetchRequest) (FetchResult, error) {
	endpoint, err := url.Parse(APIFetcher.APISourceInfo.BaseURL + "/works")
	if err != nil {
    	return FetchResult{}, err
	}

	q := endpoint.Query()
	q.Set("search", APIFetcher.APISourceInfo.Search)
	q.Set("per-page", "3")
	q.Set("api_key", APIFetcher.APISourceInfo.Apikey)
	// q.Set("select", "id,doi,title")

	endpoint.RawQuery = q.Encode()
	requestURL := endpoint.String()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
    	return FetchResult{}, err
	}

	resp, err := APIFetcher.FetcherClient.HttpClient.Do(req)
	if err != nil {
    	return FetchResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return FetchResult{}, fmt.Errorf("OpenAlex 请求失败: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return FetchResult{}, err
	}

	fmt.Println(string(body))

	return FetchResult{}, err
}