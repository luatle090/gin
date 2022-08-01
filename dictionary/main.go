package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var dict map[Dictionary]string

type Dictionary struct {
	Lang, Word string
}

type Word struct {
	Words Lang `json:"words"`
}

type Lang struct {
	En []Mean `json:"en"`
	Fr []Mean `json:"fr"`
}

type Mean struct {
	Word string `json:"word"`
	Mean string `json:"mean"`
}

func SearchEnToFr(c *gin.Context) {
	word := c.Params.ByName("word")
	mean, ok := dict[Dictionary{"en", word}]
	if !ok {
		c.JSON(http.StatusNoContent, gin.H{
			"mean": "not found",
		})
	}

	c.JSON(200, gin.H{
		"mean": mean,
	})
}

func SearchFrToEn(c *gin.Context) {
	word := c.Params.ByName("word")
	mean, ok := dict[Dictionary{"fr", word}]
	if !ok {
		c.JSON(http.StatusNoContent, gin.H{
			"mean": "not found",
		})
	}

	c.JSON(200, gin.H{
		"mean": mean,
	})
}

func LoadTuDien() {
	file, err := ioutil.ReadFile("data/TuDien.json")
	if err != nil {
		log.Println(err)
	}

	var d Word
	err = json.Unmarshal(file, &d)
	if err != nil {
		log.Println(err)
	}

	for _, v := range d.Words.En {
		dict[Dictionary{"en", v.Word}] = v.Mean
	}

	for _, v := range d.Words.Fr {
		dict[Dictionary{"fr", v.Word}] = v.Mean
	}

}

func init() {
	dict = make(map[Dictionary]string)
	LoadTuDien()
}

func main() {
	r := gin.Default()
	r.GET("/entofr/:word", SearchEnToFr)
	r.GET("/frtoen/:word", SearchFrToEn)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
