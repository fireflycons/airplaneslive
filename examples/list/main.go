package main

import (
	"context"
	"fmt"
	"log"

	"github.com/fireflycons/airplaneslive"
)

func main() {
	ctx := context.Background()
	api, err := airplaneslive.NewApi()
	if err != nil {
		log.Fatalf("Failed to create API instance: %v", err)
	}

	airlines := make([]airplaneslive.Airline, 0, 7000)

	opts := airplaneslive.DefaultListOpts
	for {
		// All search terms empty returns entire list
		list, err := api.Airlines(ctx, "", "", "", "", "", opts)

		if err != nil {
			log.Fatalf("Failed to get airlines: %v", err)
		}

		if len(list.Results) == 0 {
			break
		}

		airlines = append(airlines, list.Results...)
		// Lists are paginated - get next page
		opts = list.NextOffset(opts)
	}

	fmt.Println(len(airlines), "downloaded")
}
