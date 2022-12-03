package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
	movieService "webApp/services"
)

type MovieDTO struct {
	ID         string
	Name       string
	CreateDate string
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	movies := movieService.GetMovies()

	var movieDtos []MovieDTO
	for i := 0; i < len(movies); i++ {
		movieDto := MovieDTO{}
		movieDto.ID = movies[i].ID
		movieDto.Name = movies[i].Name
		movieDto.CreateDate = FormatMovieDate(movies[i].CreateDate)

		movieDtos = append(movieDtos, movieDto)
	}

	var index = template.Must(template.ParseFiles("views/index.html"))
	if err := index.Execute(w, movieDtos); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	var add = template.Must(template.ParseFiles("views/add.html"))
	err := add.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func SaveHandler(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("name")

	//To add movie with new ID
	id := movieService.GetCount() + 1

	movie := movieService.Movie{}
	movie.ID = fmt.Sprint(id)
	movie.Name = name
	movie.CreateDate = time.Now()

	key := fmt.Sprintf("%s%d", "mv", id)
	movieService.AddMovie(movie, key)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func FormatMovieDate(date time.Time) string {
	result := date.Format(time.RFC822)
	return result
}

func main() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/add", AddHandler)
	http.HandleFunc("/save", SaveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
