package bookapi



type BookStorage interface {
	AddBook(Book)
	GetBookByISBN(string) Book
	GetBooks() []Book
}

type BookMemoryStorage struct {
	books []Book
}

func (s *BookMemoryStorage) AddBook(newBook Book) {
	s.books = append(s.books, newBook)
}

func (s *BookMemoryStorage) GetBookByISBN(ISBN string) Book {
	for i := 0; i < len(s.books); i++ {
		if s.books[i].ISBN == ISBN {
			return s.books[i]
		}
	}
	return Book{}
}

func (s *BookMemoryStorage) GetBooks()  ([]Book) {
	books := make([]Book, len(s.books))
	copy(books, s.books)
	return books
}