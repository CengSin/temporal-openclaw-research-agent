package activity

type ArticleRequest struct {
	Topic string `json:"topic"`
}

type ArticleResult struct {
	Content string `json:"content"`
}
