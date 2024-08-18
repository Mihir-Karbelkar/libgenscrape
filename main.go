package main

import (
	"fmt"
	"libgenscrape/views"
	"libgenscrape/views/components"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	component := views.Index()
	component.Render(r.Context(), w)
}

type SearchPayload struct {
	Query string `json:"query"`
}

func GetDomain(urlToParse string) string {
	url, err := url.Parse(urlToParse)
	if err != nil {
		log.Fatal(err)
	}
	parts := strings.Split(url.Hostname(), ".")
	domain := parts[len(parts)-2] + "." + parts[len(parts)-1]
	return domain
}

func Search(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseMultipartForm(0)
	query := r.FormValue("query")
	books := GetLists(query)
	component := components.LibgenList(books)
	component.Render(r.Context(), w)
}

func Download(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseMultipartForm(0)
	mirrorUrl := r.FormValue("mirror")
	extension := r.FormValue("extension")
	bookName := r.FormValue("bookName")
	downloadUrl := GetDownloadUrl(mirrorUrl)

	component := components.Download(downloadUrl, bookName+"."+extension)
	component.Render(r.Context(), w)
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.POST("/search", Search)
	router.POST("/download", Download)
	log.Fatal(http.ListenAndServe(":8080", router))
	fmt.Println("Server started at: 8080")

}
