package article

import (
	"time"
)

type FeedArticle struct {
	Source      string		`json:"source"`
	Id          int			`json:"id"`
	Title       string		`json:"title"`
	Link        string		`json:"link"`
	PublishedAt time.Time   `json:"published_at,omitempty"`
	Guid		string		`json:"guid"`
	Abstract	string		`json:"description,omitempty"`
	ScmFid		string		`json:"scmfid"`
	FetchedAt   time.Time   `json:"fetched_at"`

	OpenAlex *OpenAlexArticleInfo `json:"openalex,omitempty"`
}

type OpenAlexArticleInfo struct {
    OpenAlexID string `json:"openalex_id"`
	PDFURL     string `json:"pdf_url,omitempty"`
    DOI        string `json:"doi"`
    OAStatus   string `json:"oa_status"`
    IsOA       bool   `json:"is_oa"`
    Version    string `json:"version,omitempty"`
    License    string `json:"license,omitempty"`
}

type ScreenedArticle struct {
	ID       	int    `json:"id"`
	Title    	string `json:"title"`
	EngTitle 	string `json:"engtitle"`
	ArticleName string `json:"articlename"`
	Link     	string `json:"link"`
	Source   	string `json:"source"`
	Reason   	string `json:"reason"`
}
