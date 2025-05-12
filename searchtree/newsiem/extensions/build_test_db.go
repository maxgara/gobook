package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
)

const first1000 = false

// build a random list of filenames for each word in the dictionary.
type filemap struct {
	word  string
	files []string
}

func createDB() []filemap {
	db := []filemap{}
	dict := loadDict()
	for _, w := range dict {
		entry := filemap{word: w, files: randomFileNames()}
		db = append(db, entry)
	}
	return db
}

func randomFileNames() []string {
	files := []string{}
	count := rand.Uint32() % 100
	for i := range count {
		files = append(files, fmt.Sprintf("random%d", i))
	}
	return files
}

// load a lot of words for testing
func loadDict() []string {
	bytes, _ := os.ReadFile("british-english-insane.txt")
	words := strings.Split(string(bytes), "\n")
	if first1000 {
		words = words[:1000]
	}
	return words
}
