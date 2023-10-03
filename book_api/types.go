package bookapi
type Author struct {
	FirstName string
	LastName  string
}

type Book struct {
	ISBN string `json:"isbn"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Author  Author `json:"author"`
}