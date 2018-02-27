package common

import (
	"math/rand"
	"time"
)

var wordList = []string{"dolphin", "cat", "fish", "dog", "bird", "elephant", "lion", "tiger"}
var adjList = []string{"good", "clever", "beautiful", "lovely", "delicious", "nice", "honest", "trustful", "brave", "generous", "enthusiastic", "radient", "rapture", "reasonable", "responsible", "romantic", "hearty", "helpful", "honesty", "honey", "humor", "useful", "Brilliant", "United", "beautiful", "smart", "obedient", "objective", "obliging", "diligent"}

/*
RandomName is used to return a random name.
Params:
Return: string
*/
func RandomName() string {
	rand.Seed(time.Now().UnixNano())
	adjStr := adjList[rand.Intn(len(adjList))]
	wordStr := wordList[rand.Intn(len(wordList))]
	return adjStr + "_" + wordStr
}

/*
RandStringBytes is used to return a n bytes random string.
Params:	n int
Return: string
*/
func RandStringBytes(n int) string {
	letterBytes := "1234567890abcdef"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
