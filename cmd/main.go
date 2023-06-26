package main

import (
	"log"

	"github.com/westradev/webber/webber"
)

func main() {
	// user := map[string]string{
	// 	"name": "Erik",
	// 	"age":  "27",
	// }

	wb, err := webber.New()
	if err != nil {
		log.Fatal(err)
	}

	coll, err := wb.CreateCollection("users")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(coll)
}
