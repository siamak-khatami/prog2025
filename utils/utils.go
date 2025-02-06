package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Article struct
type Article struct {
	ID      int    `json:"id,omitempty"`
	Title   string `json:"title"`
	Desc    string `json:"description"`
	Content string `json:"content"`
}

// Endpoint struct
type Endpoint struct {
	Url         string `json:"url"`
	Description string `json:"description"`
}

type Endpoints []Endpoint // Initiating Endpoint slice

// Articles will be saved in a map with key: ID, and content Article struct
// This is a representative of database
type Articles struct {
	collection map[int]Article //`default:"[]Article"`
	IDCounter  int             //`default:"0"`
}

// Articles' struct related methods

type ArticlesInterface interface {
	AddRandom() Articles
	Add() Articles
}

func (arts Articles) Add(a Article) Articles {
	// This function receives an Article and adds it to the caller Articles map.
	// This function gets the id from IDCounter in Articles struct, allocates
	// it to the new article id and then updates the IDCounter.
	a.ID = arts.IDCounter
	arts.IDCounter = arts.IDCounter + 1
	
	arts.collection[a.ID] = a
	return arts
}

func (arts Articles) AddRandom(n int) Articles {
	// This function receives the number of random articles (n) and generates n number of random articles.
	for i:= 0; i < n; i++ {
		new_article := Article{
			ID:      arts.IDCounter,
			Title:   fmt.Sprintf("Article %v", arts.IDCounter),
			Desc:    fmt.Sprintf("Description %v", arts.IDCounter),
			Content: fmt.Sprintf("Content %v", arts.IDCounter),
		}
		arts.collection[new_article.ID] = new_article
		arts.IDCounter = arts.IDCounter + 1
	}
	return arts
}

func Init(n int) Articles {
	// Initialization of Articles struct (database)
	var A Articles
	A.collection = make(map[int]Article) // Collection is a map and we need to make it. 
	A = A.AddRandom(n)
	return A
}

var Arts Articles = Init(5) // Initiating 5 random Articles for mocking purposes

func AddArticle(rw http.ResponseWriter, r *http.Request) {
	// Adds a single article to the Articles' collection
	var art Article
	err := json.NewDecoder(r.Body).Decode(&art)
	if err != nil {
		fmt.Println(art)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		log.Println(err.Error())
	} else {
		Arts = Arts.Add(art)
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(Arts.collection)
	}
}

func AddRandomArticle(rw http.ResponseWriter, r *http.Request) {
	// Adds a random article
	Arts = Arts.AddRandom(1)
	rw.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(rw).Encode(Arts.collection)
	if err != nil {
		http.Error(rw, "Error during JSON encoding (Error: "+err.Error()+")", http.StatusInternalServerError)
	}

}

func AllArticles(rw http.ResponseWriter, r *http.Request) {
	// Returns All articles in the Articles' collection.
	if r.Method == http.MethodGet {
		fmt.Println("Endpoint Hint: All Articles Endpoint")
		rw.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(rw).Encode(Arts.collection)
		if err != nil {
			http.Error(rw, "Error during JSON encoding.", http.StatusInternalServerError)
		}
	}
}

func HandleArticlePost(rw http.ResponseWriter, r *http.Request) {
	// There are twp post functionalities for the POST method.
	//		If body is empty: it will add a random articles.
	//		If body is not empty: it will check and add the article.
	var art Article
	if r.Body == http.NoBody {
		log.Print("A random article generated.")
		AddRandomArticle(rw, r)
	} else {
		log.Print(art)
		AddArticle(rw, r)
	}
}

func HandleArticleGet(rw http.ResponseWriter, r *http.Request) {
	// There are twp post functionalities for the POST method.
	//		If id is empty: it will return all articles
	//		If id is not empty: it will return requested article (if there is any)
	id := r.URL.Query().Get("id")
	if id == "" {
		fmt.Println("Get all articles.")
		AllArticles(rw, r)
	} else {
		fmt.Println("Get article by ID.")
		GetArticleByID(rw, r)
	}

}
func GetArticleByID(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			http.Error(rw, "No proper id passed.", http.StatusBadRequest)
		} else {
			art, ok := Arts.collection[id]
			if ok == false {
				http.Error(rw, "No article with requested id.", http.StatusBadRequest)
			} else {
				rw.Header().Set("Content-Type", "application/json")
				err := json.NewEncoder(rw).Encode(art)
				if err != nil {
					http.Error(rw, "Error during JSON encoding.", http.StatusInternalServerError)
				}
			}
		}
	} else {
		http.Error(rw, "Method "+r.Method+" is allowed.", http.StatusMethodNotAllowed)
	}

}

func GetPort(service string) string {
	var port = os.Getenv("PORT")
	switch service {
	case "main":
		if port == "" {
			port = "8080"
		}
	case "second_port":
		port = "8001"
	}
	return ":" + port
}
