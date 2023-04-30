package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type StringArray []string

func (a *StringArray) UnmarshalJSON(data []byte) error {
	if data[0] == '"' {
		var s string
		err := json.Unmarshal(data, &s)
		if err != nil {
			return err
		}
		*a = []string{s}
		return nil
	}

	if data[0] == '{' {
		var m map[string]interface{}
		err := json.Unmarshal(data, &m)
		if err != nil {
			return err
		}
		if name, ok := m["name"].(string); ok {
			*a = []string{name}
			return nil
		}
	}

	var arr []string
	err := json.Unmarshal(data, &arr)
	if err != nil {
		return err
	}
	*a = arr
	return nil
}

type DialogflowRequest struct {
	QueryResult struct {
		Parameters struct {
			Genre  StringArray `json:"genre"`
			Author interface{} `json:"person"`
		} `json:"parameters"`
	} `json:"queryResult"`
}

type DialogflowResponse struct {
	FulfillmentText string `json:"fulfillmentText"`
}

type Book struct {
	Title  string
	Author string
	Genre  string
}

var books = []Book{
	{Title: "Introduction to Algorithms", Author: "Thomas H. Cormen, Charles E. Leiserson, Ronald L. Rivest, and Clifford Stein", Genre: "Computer Science"},
	{Title: "A Top-down Approach Featuring the Internet", Author: "Keith W. Ross", Genre: "Computer Science"},
	{Title: "Introduction To The Theory Of Computation", Author: "Michael Sipser", Genre: "Computer Science"},
	{Title: "Dune", Author: "Frank Herbert", Genre: "Science Fiction"},
	{Title: "The Lord of the Rings", Author: "J.R.R. Tolkien", Genre: "Fantasy"},
	{Title: "Neuromancer", Author: "William Gibson", Genre: "Science Fiction"},
	{Title: "The Name of the Wind", Author: "Patrick Rothfuss", Genre: "Fantasy"},
	{Title: "The Story of Art", Author: "E.H. Gombrich", Genre: "Art History"},
	{Title: "The Photographer's Eye: Composition and Design for Better", Author: "Brandon Stanton", Genre: "Photography"},
	{Title: "The Art Lesson", Author: "Phaidon Editors", Genre: "Art History"},
	{Title: "Rich Dad Poor Dad", Author: "Robert T. Kiyosaki", Genre: "Personal Finance"},
	{Title: "The Lean Startup", Author: "Eric Ries", Genre: "Entrepreneurship"},
	{Title: "Good to Great", Author: "Jim Collins", Genre: "Management"},
	{Title: "The 7 Habits of Highly Effective People", Author: "Stephen R. Covey", Genre: "Self-help"},
	{Title: "The Girl with the Dragon Tattoo", Author: "Stieg Larsson", Genre: "Mystery"},
	{Title: "Gone Girl", Author: "Gillian Flynn", Genre: "Thriller"},
	{Title: "The Silence of the Lambs", Author: "Thomas Harris", Genre: "Thriller"},
	{Title: "The Da Vinci Code", Author: "Dan Brown", Genre: "Mystery"},
	{Title: "Big Little Lies", Author: "Liane Moriarty", Genre: "Mystery"},
	{Title: "Mastering the Art of French Cooking", Author: "Julia Child, Simone Beck, Louisette Bertholle", Genre: "Cooking"},
	{Title: "Veganomicon: The Ultimate Vegan Cookbook", Author: "Isa Chandra Moskowitz, Terry Hope Romero", Genre: "Cooking"},
	{Title: "How to Cook Everything", Author: "Mark Bittman", Genre: "Cooking"},
	{Title: "To Kill a Mockingbird", Author: "Harper Lee", Genre: "Fiction"},
	{Title: "1984", Author: "George Orwell", Genre: "Fiction"},
	{Title: "Pride and Prejudice", Author: "Jane Austen", Genre: "Fiction"},
	{Title: "The Catcher in the Rye", Author: "J.D. Salinger", Genre: "Fiction"},
	{Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Genre: "Fiction"},
	{Title: "Sapiens: A Brief History of Humankind", Author: "Yuval Noah Harari", Genre: "History"},
	{Title: "Guns, Germs, and Steel: The Fates of Human Societies", Author: "Jared Diamond", Genre: "History"},
	{Title: "The Rise and Fall of the Third Reich", Author: "William L. Shirer", Genre: "History"},
	{Title: "A People's History of the United States", Author: "Howard Zinn", Genre: "History"},
	{Title: "The Gulag Archipelago", Author: "Aleksandr Solzhenitsyn", Genre: "History"},
	{Title: "Dog Man", Author: "Dav Pilkey", Genre: "Kids"},
	{Title: "Zita the Spacegirl", Author: "Ben Hatke", Genre: "Kids"},
	{Title: "Big Nate Make a Splash", Author: "Lincoln Peirce", Genre: "Kids"},
	{Title: "Pride and Prejudice", Author: "Jane Austen", Genre: "Romance"},
	{Title: "Outlander", Author: "Diana Gabaldon", Genre: "Romance"},
	{Title: "The Notebook", Author: "Nicholas Sparks", Genre: "Romance"},
	{Title: "Jane Eyre", Author: "Charlotte Bronte", Genre: "Romance"},
	{Title: "The Bridges of Madison County", Author: "Robert James Waller", Genre: "Romance"},
	{Title: "The Autobiography of Benjamin Franklin", Author: "Benjamin Franklin", Genre: "Biography"},
	{Title: "Steve Jobs", Author: "Walter Isaacson", Genre: "Biography"},
	{Title: "I Know Why the Caged Bird Sings", Author: "Maya Angelou", Genre: "Biography"},
	{Title: "The Glass Castle", Author: "Jeannette Walls", Genre: "Biography"},
	{Title: "The Immortal Life of Henrietta Lacks", Author: "Rebecca Skloot", Genre: "Biography"},
}

func getBookRecommendationByAuthor(author string) string {
	filteredBooks := []Book{}
	for _, book := range books {
		if strings.Contains(strings.ToLower(book.Author), strings.ToLower(author)) {
			filteredBooks = append(filteredBooks, book)
		}
	}

	if len(filteredBooks) == 0 {
		return fmt.Sprintf("I couldn't find any books by %s. Please try a different author.", author)
	}

	rand.Seed(time.Now().UnixNano())
	randomBook := filteredBooks[rand.Intn(len(filteredBooks))]
	return fmt.Sprintf("I recommend reading '%s' by %s. It's a great %s book.", randomBook.Title, randomBook.Author, randomBook.Genre)
}

func getBookRecommendation(genres, authors []string) string {
	if len(authors) > 0 {
		return getBookRecommendationByAuthor(authors[0])
	}

	lowerCaseGenres := make([]string, len(genres))
	for i, g := range genres {
		lowerCaseGenres[i] = strings.ToLower(g)
	}

	filteredBooks := []Book{}
	for _, book := range books {
		lowerCaseBookGenre := strings.ToLower(book.Genre)

		genreMatch := false
		for _, g := range lowerCaseGenres {
			if lowerCaseBookGenre == g {
				genreMatch = true
				break
			}
		}

		if genreMatch {
			filteredBooks = append(filteredBooks, book)
		}
	}

	if len(filteredBooks) == 0 {
		return "I couldn't find any books matching your criteria. Please try a different genre or author."
	}

	rand.Seed(time.Now().UnixNano())
	randomBook := filteredBooks[rand.Intn(len(filteredBooks))]
	return fmt.Sprintf("I recommend reading '%s' by %s. It's a great %s book.", randomBook.Title, randomBook.Author, randomBook.Genre)
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	var req DialogflowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	genres := req.QueryResult.Parameters.Genre
	var authors []string

	switch author := req.QueryResult.Parameters.Author.(type) {
	case string:
		authors = []string{author}
	case []interface{}:
		for _, a := range author {
			if aStr, ok := a.(string); ok {
				authors = append(authors, aStr)
			} else if aMap, ok := a.(map[string]interface{}); ok {
				if name, ok := aMap["name"].(string); ok {
					authors = append(authors, name)
				}
			} else if aList, ok := a.([]interface{}); ok {
				for _, listItem := range aList {
					if listItemStr, ok := listItem.(string); ok {
						authors = append(authors, listItemStr)
					} else if listItemMap, ok := listItem.(map[string]interface{}); ok {
						if name, ok := listItemMap["name"].(string); ok {
							authors = append(authors, name)
						}
					}
				}
			}
		}
	case map[string]interface{}:
		if name, ok := author["name"].(string); ok {
			authors = []string{name}
		}
	}

	bookRecommendation := getBookRecommendation(genres, authors)
	response := DialogflowResponse{
		FulfillmentText: bookRecommendation,
	}
	fmt.Printf("DialogflowRequest: %#v\n", req)

	fmt.Printf("Genres: %v\n", genres)
	fmt.Printf("Authors: %v\n", authors)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", webhookHandler).Methods(http.MethodPost)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Listening on port %s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}
