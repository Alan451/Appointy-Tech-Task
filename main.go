package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"encoding/json"
	"strings"
	"strconv"
	// "math"
	// "sync"
)

// Mutex Lock used for making the server thread safe
// In a real server we will be defining this var inside the database
var isMutexLockOn bool
var idVariable int
//Article is ...
type Article struct {
	ID  int `json:"id"`
    Title string `json:"title"`
    SubTitle string `json:"subTitle"`
    Content string `json:"content"`
	CreatedAt string `json:"createdAt"`
}
func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
// Articles ... global array
// that we can then populate in our main function
// to simulate a database
var Articles []Article
func articles(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
		case "GET":
			page, ok1 := r.URL.Query()["page"]
			limit, ok2 := r.URL.Query()["limit"]
			if !ok1 || !ok2 || len(page[0])<1 || len(limit[0])<1{
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(Articles)	
				return
			}
			pageNo,err1 := strconv.Atoi(page[0])
			limitValue,err2 := strconv.Atoi(limit[0])
			
			//If limit Value is 0 or invalid parameter for page or limit we return all of the articles
			if err1!=nil || err2!=nil || limitValue==0{
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(Articles)	
				
			}else{
				//Indexing of page starts from 0
				if pageNo*limitValue >= len(Articles){
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprintf(w,"Not Found")
				}else{
					w.WriteHeader(http.StatusOK)
					for i:=pageNo*limitValue; i<min(len(Articles),(pageNo+1)*limitValue);i++{
						json.NewEncoder(w).Encode(Articles[i])
					}
				}
			}
		case "POST":
			//For making the server thread safe
			for isMutexLockOn{
				time.Sleep(1*time.Millisecond)
			}
			isMutexLockOn = true
			w.WriteHeader(http.StatusCreated)
			var a Article
    		// Try to decode the request body into the struct. If there is an error,
   			// respond to the client with the error message and a 400 status code.
			err := json.NewDecoder(r.Body).Decode(&a)
			a.CreatedAt = time.Now().Format(time.RFC850)
			a.ID = idVariable
			if err == nil {
				idVariable++
				Articles = append(Articles,a)
			} else{
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			isMutexLockOn =false
	}
	
}
func getArticleByID(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	id := strings.TrimPrefix(r.URL.Path, "/articles/")
	if len(id) == 0 {
		articles(w, r)
	}else{
		var b bool
		b = false
		id,err := strconv.Atoi(id)
		if err != nil{
			fmt.Fprintf(w,"Not a valid path")
			return
		}
		for _,article := range Articles{
			
			if article.ID == id {
				b=true
				json.NewEncoder(w).Encode(article)
			}
		}
		if !b {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w,"Not Found")
		}else{
			w.WriteHeader(http.StatusOK)
		}
	}
}
func searchArticles(w http.ResponseWriter, r *http.Request){
	keys, ok := r.URL.Query()["q"]
	if !ok || len(keys[0]) < 1 {
		fmt.Fprintln(w,"Url Param 'q' is missing")
        return
	}
	key := keys[0]
	var b bool
	b = false
	for _,article := range Articles{
		if article.Title == key || article.SubTitle==key || article.Content==key{
			b=true
			json.NewEncoder(w).Encode(article)
		}
	}
	if b{
		w.WriteHeader(http.StatusOK)
	}else{
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w,"Not Found")
	}
}
func handleRequests() {
	http.HandleFunc("/articles", articles)
	http.HandleFunc("/articles/", getArticleByID)
	http.HandleFunc("/articles/search", searchArticles)
    log.Fatal(http.ListenAndServe(":8000", nil))
} 


func main() {
	Articles = []Article{
        Article{ID:1 ,Title: "Hello", SubTitle: "Article Description", Content: "Article Content", CreatedAt : time.Now().Format(time.RFC850)},
        Article{ID:2 ,Title: "Hello 2", SubTitle: "Article Description", Content: "Article Content", CreatedAt : time.Now().Format(time.RFC850)},
	}
	idVariable = 3
	isMutexLockOn = false
    handleRequests()
}

