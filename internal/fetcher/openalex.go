package fetcher

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	//"io"
)

//这是收到响应后要提取出的信息
type OpenAlexArticleInfo struct {
    OpenAlexID string `json:"openalex_id"`
    DOI        string `json:"doi"`
    OAStatus   string `json:"oa_status"`
    IsOA       bool   `json:"is_oa"`
    Version    string `json:"version,omitempty"`
    License    string `json:"license,omitempty"`
}

func (APIFetcher *APIFetcher) Fetch(ctx context.Context, fetchReq FetchRequest) (FetchResult, error) {
	endpoint, err := url.Parse(APIFetcher.APISourceInfo.BaseURL + "/works")
	if err != nil {
    	return FetchResult{}, err
	}

	query := endpoint.Query()
	query.Set("search", APIFetcher.APISourceInfo.Search)
	//query.Set("per-page", "3")
	query.Set("api_key", APIFetcher.APISourceInfo.Apikey)
	query.Set("select", "id,doi,title,abstract_inverted_index")

	endpoint.RawQuery = query.Encode()
	requestURL := endpoint.String()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
    	return FetchResult{}, err
	}

	/*urltest := "https://link.springer.com/article/10.1007/s10676-010-9253-3"
	resptest, err:= APIFetcher.FetcherClient.HttpClient.Get(urltest)
	defer resptest.Body.Close()

	body, err := io.ReadAll(resptest.Body)
	if err != nil {
    	fmt.Println("读取失败:", err)
	}

	fmt.Printf("状态码: %d\n", resptest.StatusCode)
	fmt.Printf("最终URL: %s\n", resptest.Request.URL)
	fmt.Printf("Content-Type: %s\n", resptest.Header.Get("Content-Type"))
	fmt.Printf("长度: %d\n", len(body))
	fmt.Printf("前 500 字节:\n%s\n", body)
*/


	resp, err := APIFetcher.FetcherClient.HttpClient.Do(req)
	if err != nil {
    	return FetchResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return FetchResult{}, fmt.Errorf("OpenAlex 请求失败: %d", resp.StatusCode)
	}

	/*body, err := io.ReadAll(resp.Body)
	if err != nil {
		return FetchResult{}, err
	}*/

	ParseOpenalex(resp.Body)
	//fmt.Println(string(body))

	return FetchResult{}, err
}