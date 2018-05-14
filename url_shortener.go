package main

import ("fmt"
		"database/sql"
		_ "github.com/go-sql-driver/mysql"
		"strings"
		"math"
		// "bufio"
		// "os"
		"net/http"
		"log"
		"html/template"
		// "io/ioutil"
	)

// all possible characters
var characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var tpl *template.Template


func init() {
	// intitalize template
	tpl = template.Must(template.ParseGlob("template/*.html"))
}



func id_to_shortener(id int) string{
	var shorturl = ""
	if id == 0{
		return string(characters[0])
	} else{
		for id > 0{
			// convert urlid to shortenedurl by modular division on int id and using result to retrieve index from 
			// character string to concatentate to shorturl to produce the shortened url

			shorturl = string(characters[id%62]) + shorturl
			id = id/62		
			// fmt.Println(id)
		}
		return shorturl
	}
}

func shortened_to_id(shortened string) int{

	var id, i, strlength = 0, 0, len(shortened)-1
	// convert shortenedurl to url id by looping through the shortenedurl and finding that character's index on the 
	// character string. Then multiplied the character index by 62 to the i power, where i represents the place on the 
	// shorturl
	for i <= strlength{
		id += strings.Index(characters, string(shortened[i])) * int(math.Pow(62, float64(strlength - i)))
        fmt.Println(id, string(shortened[i]), strlength)
        i++
	}

	return id
	
}

func index_handle(w http.ResponseWriter, r *http.Request){
	// initial template to load 
	tpl.ExecuteTemplate(w, "index.html", nil)
}

func shortener_handle(w http.ResponseWriter, r *http.Request){
	var s int
	// var largest string
	db, err := sql.Open("mysql", "root:password@tcp(172.17.0.2:3306)/urlshortener")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

// get the largest urlid to get the new shortened url
	stmt, er := db.Prepare("SELECT urlid FROM url ORDER BY urlid DESC LIMIT 0, 1")
	if er != nil{
		log.Fatal(er)
	}
	er = stmt.QueryRow().Scan(&s)
	switch{
		case er == sql.ErrNoRows:
        	fmt.Printf("no ID.")
		case er != nil:
		    fmt.Printf("nil return")
		default:
		    fmt.Printf(string(s))
	}


// on post, insert new longurl and the shortenedurl
	if r.Method == "POST"{
		r.ParseForm()
		var longurl = r.Form["longurl"]
		
		shortenedid :=  id_to_shortener(s+1)
		fmt.Println(longurl, s)

		_, er := db.Exec("insert into url (longurl, shorturl) values (?, ?)", longurl[0], shortenedid)
		switch{
			case er == sql.ErrNoRows:
	        	fmt.Printf("no ID.")
			case er != nil:
			    fmt.Printf("nil return")
			default:
			    fmt.Printf("successful insert")
		}
		db.Close()
		tpl.ExecuteTemplate(w, "index.html", shortenedid)

	}
}

func short_to_url_handler(w http.ResponseWriter, r *http.Request) {
	var longurl string
	db, err := sql.Open("mysql", "root:password@tcp(172.17.0.2:3306)/urlshortener")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	//get shorturl from parameter in the requested url 
	shorturl := r.URL.Query().Get("short")
	log.Println("short to url", string(shorturl))
	urlid := shortened_to_id(string(shorturl))

	//retrieve longurl from db by converting shorturl to urlid 
	stmt, er := db.Prepare("select longurl from url where urlid=?")
	if er != nil{
		log.Fatal(er)
	}


	er = stmt.QueryRow(urlid).Scan(&longurl)


	http.Redirect(w, r, longurl, 301)

	tpl.ExecuteTemplate(w, "index.html", shorturl)


}


func main() {

	http.HandleFunc("/", index_handle)
	http.HandleFunc("/shortener_handle", shortener_handle)
	http.HandleFunc("/shorttolong/", short_to_url_handler)

	http.ListenAndServe(":8000", nil)

	// fmt.Println(shortened_to_id("AnNB"))
	fmt.Println(characters)


}

