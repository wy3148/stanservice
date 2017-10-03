package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	maxFormSize = int64(10 << 20) //10 M for request body
)

type showImage struct {
	ShowImage string `json:"showImage"`
}

type Show struct {
	Drm          bool      `json:"drm"`
	EpisodeCount int       `json:"episodeCount"`
	Image        showImage `json:"image"`
	Slug         string    `json:"slug"`
	Titile       string    `json:"title"`
}

type ShowRequest struct {
	Payload []Show `json:"payload"`
}

type SimpleShow struct {
	Image string `json:"image"`
	Slug  string `json:"slug"`
	Title string `json:"title"`
}

type ShowResponse struct {
	Response []*SimpleShow `json:"response"`
}

// Don't use it until we add more fields
// type ShowErrorRes struct {
// 	ErrorDetail string `json:"error"`
// }

//Input the request data, we need to filter those elements with drm true and valid episcode count
func filterShow(r *ShowRequest) ([]byte, error) {

	var shows []*SimpleShow
	var res ShowResponse

	if r != nil && r.Payload != nil {
		for _, v := range r.Payload {

			if v.Drm && v.EpisodeCount > 0 {
				matchedItem := &SimpleShow{Image: v.Image.ShowImage, Slug: v.Slug, Title: v.Titile}
				shows = append(shows, matchedItem)
			}
		}
	}

	res.Response = shows
	return json.Marshal(&res)
}

func errorResponse(w http.ResponseWriter, msg string) {
	http.Error(w, `{"error":"`+msg+`"}`, http.StatusBadRequest)
}

func HandleShow(w http.ResponseWriter, req *http.Request) {

	if req.Body == nil {
		errorResponse(w, "empty request body")
		return
	}

	reader := io.LimitReader(req.Body, maxFormSize+1)
	b, err := ioutil.ReadAll(reader)

	if err != nil {
		errorResponse(w, err.Error())
		return
	}

	var showData ShowRequest
	err = json.Unmarshal(b, &showData)

	if err != nil {
		errorResponse(w, err.Error())
		return
	}

	res, err := filterShow(&showData)

	if err != nil {
		errorResponse(w, err.Error())
		return
	}
	w.Write(res)
}

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.HandleFunc("/", HandleShow)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
