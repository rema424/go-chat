package element

type Message struct {
	ID   int64  `json:"-" db:"id"`
	User string `json:"user" db:"user"`
	Body string `json:"body" db:"body"`
}
