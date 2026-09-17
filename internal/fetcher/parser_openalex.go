package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"secmind/internal/article"
	"time"
	//"secmind/configs"
)

//以下是网站的响应信息
type OpenAlexResponse struct {
    Meta    OpenAlexMeta    `json:"meta"`
    Results []OpenAlexWork  `json:"results"`
}

type OpenAlexMeta struct {
    Count   int `json:"count"`
    Page    int `json:"page"`
    PerPage int `json:"per_page"`
}

type OpenAlexWork struct {
    ID              string `json:"id"`
    DOI             string `json:"doi"`
    Title           string `json:"title"`
    PublicationDate string `json:"publication_date"`
    PublicationYear int    `json:"publication_year"`

    PrimaryLocation *OpenAlexLocation `json:"primary_location"`
    BestOALocation  *OpenAlexLocation `json:"best_oa_location"`

    OpenAccess struct {
        IsOA     bool   `json:"is_oa"`
        OAStatus string `json:"oa_status"`
        OAURL    string `json:"oa_url"`
    } `json:"open_access"`

	AbstractInvertedIndex map[string][]int `json:"abstract_inverted_index,omitempty"`
}

type OpenAlexLocation struct {
    LandingPageURL string `json:"landing_page_url"`
    PDFURL         string `json:"pdf_url"`
    IsOA           bool   `json:"is_oa"`
    Version        string `json:"version"`
    License        string `json:"license"`
}


func ParseOpenalex(reader io.Reader) ([]article.FeedArticle, error) {
	var resp OpenAlexResponse

	err := json.NewDecoder(reader).Decode(&resp)
	if err != nil {
        return nil, fmt.Errorf("解析 OpenAlex响应 JSON 失败: %w", err)
    }

	var articles []article.FeedArticle

	for _, work := range resp.Results {
		openalexArticle := article.FeedArticle{
			Source: "openalex",
			Guid: work.ID,
			Title: work.Title,
			PublishedAt: parseOpenAlexData(work.PublicationDate),
			OpenAlex: &article.OpenAlexArticleInfo{
				OpenAlexID: work.ID,
				DOI: work.DOI,
				IsOA: work.OpenAccess.IsOA,
				OAStatus: work.OpenAccess.OAStatus,
			},
		}
		openalexArticle.Link = pickLink(work)
	}
}

func parseOpenAlexData(s string) time.Time {
	if s == "" {
		return time.Time{}
	}

	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}
	}

	return  t
}