package main

import (
	"log"

	"github.com/westradev/webber/webber"
)

func main() {
	userData := webber.M{
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

	log.Println("Inserted user with id: ", id)

	results, err := wb.Find("users", webber.Filter{})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v\n", results)
}
