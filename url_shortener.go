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
		"io/ioutil"
	)

// db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/")
// if err != nil {
// 	panic(err.Error())
// }
// defer db.Close()
// fmt.Println("connected")
var characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
var tpl *template.Template
type Page struct{
	Title string
	Body[] byte
}

func init() {
	tpl = template.Must(template.ParseGlob("template/*.html"))

	// db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/")
	// if err != nil {
	// 	panic(err.Error())
	// }
	// defer db.Close()


}

func (p *Page) save() error {
    filename := p.Title + ".txt"
    return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

func id_to_shortener(id int) string{
	var shorturl = ""
	if id == 0{
		return string(characters[0])
	} else{
		for id > 0{
			shorturl = string(characters[id%62]) + shorturl
			id = id/62		
			// fmt.Println(id)
		}
		return shorturl
	}
}

func shortened_to_id(shortened string) int{

	var id, i, strlength = 0, 0, len(shortened)-1
	for i <= strlength{
		id += strings.Index(characters, string(shortened[i])) * int(math.Pow(62, float64(strlength - i)))
        fmt.Println(id, string(shortened[i]), strlength)
        i++
	}

	return id
	
}

func index_handle(w http.ResponseWriter, r *http.Request){
	tpl.ExecuteTemplate(w, "index.html", nil)
}

func shortener_handle(w http.ResponseWriter, r *http.Request){
	var s int
	// var largest string
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/urlshortener")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()


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

	if r.Method == "POST"{
		r.ParseForm()
		var longurl = r.Form["longurl"]
		
		shortenedid :=  id_to_shortener(s+1)
		fmt.Println(longurl, s)
		// db.Query("insert into url (longurl, shorturl) values (?, ?)")
		// db.Close()
		// stmt, er = db.Prepare("insert into url (longurl, shorturl) values (?, ?)")
		// if er != nil{
		// 	log.Fatal(er)
		// }

		// er = stmt.Exec(longurl[0], shortenedid)
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
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/urlshortener")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	shorturl := r.URL.Query().Get("short")
	log.Println("short to url", string(shorturl))
	urlid := shortened_to_id(string(shorturl))


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

