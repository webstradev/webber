package main

import (
	"log"

	"github.com/westradev/webber/webber"
)

func main() {
	userData := map[string]string{
		"name": "Erik",
		"age":  "27",
	}

	wb, err := webber.New()
	if err != nil {
		log.Fatal(err)
	}

	id, err := wb.Insert("users", userData)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(id)
}
