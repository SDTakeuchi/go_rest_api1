package main

import (
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"math/rand"
	"strconv"
	"github.com/gorilla/mux"
)


// Book struct (Model)
type Book struct {
	ID		string  `json:"id"`
	Isbn	string  `json:"isbn"`
	Title 	string  `json:"title"`
	Author	*Author `json:"author"`
}

type Author struct {
	FirstName	string  `json:"firstname"`
	LastName    string  `json:"lastname"`
}

// Init books var as a slice Book struct
var books []Book

//Get all books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}


//Get single book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r) // get params

	// loop through the books and find with id
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

// Create a new book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book Book //新規作成用の金型を持ってきている
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(10000000)) //Mock ID
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

// Update a nook
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)

			var book Book //新規作成用の金型を持ってきている
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	json.NewEncoder(w).Encode(books)
}

// Delete a book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			 // think like making a new list excluding the matching item
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}


func main() {
	// Init router first
	r := mux.NewRouter()

	//Mock data - @todo - implement DB
	books = append(books, Book{ID:"1", Isbn:"38724", Title:"Progmatic Programmer", Author:&Author{FirstName:"Douglas", LastName:"Takeuchi"}})
	books = append(books, Book{ID:"2", Isbn:"48917", Title:"Automate the stuff", Author:&Author{FirstName:"Kirk", LastName:"Hammet"}})

	// Route handlers \ endpoints
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	fmt.Println("Server is running...")
	log.Fatal(http.ListenAndServe(":8000", r))

}