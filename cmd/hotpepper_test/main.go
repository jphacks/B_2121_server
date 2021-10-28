package main

import (
	"fmt"
	"os"

	"github.com/jphacks/B_2121_server/models"
	"github.com/jphacks/B_2121_server/restaurant_search"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage %s ApiKey keyword", os.Args)
		return
	}

	apiKey := os.Args[1]
	keyword := os.Args[2]
	search := restaurant_search.NewSearchApi(apiKey)
	r, err := search.Search(keyword, &models.Location{
		Latitude:  35.0344823881,
		Longitude: 135.7841217166,
	}, 10)
	if err != nil {
		fmt.Printf("failed to search restaurants: %v", err)
		return
	}
	for _, restaurant := range *r {
		fmt.Printf("%v\n", restaurant)
	}
}
