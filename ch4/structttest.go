package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Vector struct {
	X, Y int
}
type segment struct {
	start Vector
	end   Vector
}
type colorVector struct {
	Vector
	Color uint8
}

type Movie struct {
	Actors []string
	Title  string `json:"movie_title"`
	Year   int
}

func main() {
	var cv = colorVector{Color: 255}
	cv.X = 5
	var gwtwActors = []string{"harrison ford", "robert plant", "emma stone"}
	gwtw := Movie{gwtwActors, "gone with the wind", 1995}
	var movies = []Movie{gwtw}
	moviesJSON, err := json.MarshalIndent(movies, "", "  ")
	if err != nil {
		fmt.Printf("JSON Marshal fail:%v", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", moviesJSON)
}
