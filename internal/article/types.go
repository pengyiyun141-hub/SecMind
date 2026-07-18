package article

type Article struct {
	Source      string
	Id          int
	Title       string
	Link        string
	Description string
	Filename    string
}

type ScreenedArticle struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	EngTitle string `json:"engtitle"`
	ArticleName string `json:"articlename"`
	Link     string `json:"link"`
	Source   string `json:"source"`
	Reason   string `json:"reason"`
}
