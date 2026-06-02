package model

type Article struct {
	Source      string
	Id          int
	Title       string
	Link        string
	SourceCount int
	Description string
}

type ScreenedArticle struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	EngTitle string `json:"engtitle"`
	Link     string `json:"link"`
	Source   string `json:"source"`
	Reason   string `json:"reason"`
}
