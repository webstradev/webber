package main

import (
	"log"

	"github.com/westradev/webbr/webbr"
)

func main() {
	userData := webbr.M{
		"name": "Erik",
		"age":  "27",
	}

	wb, err := webbr.New()
	if err != nil {
		log.Fatal(err)
	}

	id, err := wb.Insert("users", userData)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Inserted user with id: ", id)

	results, err := wb.Find("users", webbr.Filter{})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v\n", results)
}
