package domain

type Status string

const (
	NOT_STUDIED Status = "not_studied"
	FAMILIAR    Status = "familiar"
	LEARNED     Status = "learned"
	REVIEWED    Status = "review"
)

type Goal struct {
	ID      int    `json:"id"`
	Text    string `json:"text"`
	Comment string `json:"comment"`
}

type Topic struct {
	ID           int    `json:"id"`
	Key          string `json:"key"`
	Comment      string `json:"comment"`
	Title        string `json:"title"`
	Status       Status `json:"Status"`
	LastReviewed string `json:"last_reviewed"`
}

type Task struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Comment      string `json:"comment"`
	Status       Status `json:"Status"`
	LastReviewed string `json:"last_reviewed"`
}
