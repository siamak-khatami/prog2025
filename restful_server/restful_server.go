package restful_server

import (
	"Code/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func ArticleApiRouter() func(http.ResponseWriter, *http.Request) {

	return func(rw http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodGet:
			fmt.Println("Router received " + r.Method + " method.")
			utils.HandleArticleGet(rw, r)
		case http.MethodPost:
			fmt.Println("Router received " + r.Method + " Method.")
			utils.HandleArticlePost(rw, r)
			// utils.AddArticle(rw, r)
		default:
			fmt.Println("Method not defined (Method: " + r.Method + ")")
			http.Error(rw, "Method "+r.Method+" not supported.", http.StatusNotImplemented)
		}
	}
}

func RestFulHomePage(rw http.ResponseWriter, r *http.Request) {
	// If the method is Get, it will return the content, otherwise, will rise an error
	if r.Method == http.MethodGet {
		// An example Article to show as response
		article := utils.Article{Title: "Article title", Desc: "Article description", Content: "Article content"}
		output, err := json.Marshal(article)
		if err != nil {
			http.Error(rw, "Error during JSON marshaling.", http.StatusInternalServerError)
		}
		// To have a pretty output for existence endpoints, we defined a struct for them.
		endpoints := utils.Endpoints{
			utils.Endpoint{Url: "/", Description: "GET: Welcome to the PROG2005 Cloud Technologies home page!"},
			utils.Endpoint{Url: "/article?id=SomeThing", Description: "GET: returns an article by id (where {id} is the article ID of concern)"},
			utils.Endpoint{Url: "/article", Description: "GET: returns all articles"},
			//utils.Endpoint{Url: "/article/add-random", Description: "GET: Adds random item to the articles list"},
			utils.Endpoint{Url: "/article", Description: "POST: Takes information for a new article and adds it to the memory. Article information should be provided using the schema below (but without the ID field, since that is generated and returned by the service upon adding)." +
				";  Schema: " + string(output) + ". If no input structure is provided, a random new article is added to memory."}}

		rw.Header().Set("Content-Type", "application/json") // Writes the response content type
		json.NewEncoder(rw).Encode(endpoints)
	} else {
		// If the method is not Get.
		http.Error(rw, "Method "+r.Method+" not supported.", http.StatusNotImplemented)
	}
}

func RestFulApi() {
	//Home page which lists available endpoints.
	http.HandleFunc(utils.HomeEndPoint, RestFulHomePage)
	// An endpoint in which returns
	//								- Specific article by id (Method: Get; parameter: is),
	//								- All articles (Method: Get; parameter: no parameter),
	//								- Adds an article (Method: Post, body: Article Struct),
	//								- Adds a random article (Method: Post, body: empty body)
	http.HandleFunc(utils.ArticleEndPoint, ArticleApiRouter())
	Port := utils.GetPort(utils.MainPort)
	fmt.Println("RestFul Api: Listen and server on: http://localhost" + Port)
	log.Fatal(http.ListenAndServe(Port, nil))
}
