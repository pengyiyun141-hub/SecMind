package article

import (
	"time"
)

type FeedArticle struct {
	Source      string		`json:"source"`
	Id          int			`json:"id"`
	Title       string		`json:"title"`
	Link        string		`json:"link"`
	Description string		`json:"description,omitempty"`
	ScmFid		string		`json:"scmfid"`
	FetchedAt   time.Time   `json:"fetched_at"`
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
