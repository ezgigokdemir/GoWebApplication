package services

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/couchbase/gocb"
)

type Configuration struct {
	Url      string
	Username string
	Password string
}

type Movie struct {
	ID         string    `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	CreateDate time.Time `json:"createdate,omitempty"`
}

type Count struct {
	Value int `json:"count,omitempty"`
}

func (conf Configuration) GetConfig() Configuration {
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&conf)
	if err != nil {
		fmt.Println("error:", err)
	}
	return conf
}

func GetBucket() *gocb.Bucket {
	var conf Configuration
	url := conf.GetConfig().Url
	username := conf.GetConfig().Username
	password := conf.GetConfig().Password
	cluster, cerr := gocb.Connect(url)
	if cerr != nil {
		fmt.Printf("Cluster Open Error: %v", cerr)
	}

	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: username,
		Password: password,
	})

	bucket, berr := cluster.OpenBucket("movie", "")
	if berr != nil {
		fmt.Printf("Bucket Open Error: %v", berr)
	}

	return bucket
}

func AddMovie(movie Movie, key string) {
	bucket := GetBucket()
	bucket.Upsert(key, movie, 0)
}

func GetMovies() []Movie {
	bucket := GetBucket()
	var movies []Movie
	query := gocb.NewN1qlQuery("SELECT m.* FROM movie m")
	rows, _ := bucket.ExecuteN1qlQuery(query, []interface{}{})
	var row interface{}
	for rows.Next(&row) {
		// Convert map to json string
		jsonStr, err := json.Marshal(row)
		if err != nil {
			fmt.Println(err)
		}
		// Convert json string to struct
		var movie Movie
		if err := json.Unmarshal(jsonStr, &movie); err != nil {
			fmt.Println(err)
		}
		movies = append(movies, movie)
	}

	return movies
}

func GetCount() int {
	var count Count
	bucket := GetBucket()
	query := gocb.NewN1qlQuery("SELECT count(*) as Count FROM movie")
	rows, _ := bucket.ExecuteN1qlQuery(query, []interface{}{})
	var row interface{}
	for rows.Next(&row) {
		jsonStr, err := json.Marshal(row)
		if err != nil {
			fmt.Println(err)
		}

		if err := json.Unmarshal(jsonStr, &count); err != nil {
			fmt.Println(err)
		}
	}
	return count.Value
}
