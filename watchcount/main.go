package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/sync/errgroup"

	"github.com/Khan/genqlient/graphql"
	"github.com/google/uuid"
)

func main() {
	maxConcurrency, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("You must provide a max concurrency as the first argument. e.g: watchcount 1000")
		fmt.Println("Error parsing max concurrency: ", err)
		return
	}

	ctx := context.Background()
	url := "https://yuovision-api-tmp.yuorei.com/graphql"

	client := graphql.NewClient(url, http.DefaultClient)

	resp, err := getVideos(ctx, client)
	if err != nil {
		fmt.Println("Error getting videos: ", err)
		return
	}
	fmt.Println("Got videos")

	eg, _ := errgroup.WithContext(ctx)
	eg.SetLimit(maxConcurrency)
	for {
		for _, video := range resp.Videos {
			eg.Go(func() error {
				return incrementWatchCountFunc(ctx, video.Id)
			})

		}
	}
}

func incrementWatchCountFunc(ctx context.Context, videoID string) error {
	println("Incrementing watch count for video ID: ", videoID)
	url := "https://yuovision-api-tmp.yuorei.com/graphql"

	uuid, err := uuid.NewRandom()
	if err != nil {
		fmt.Printf("UUID generation error: %v\n", err)
		return err
	}
	clientID := "client_" + uuid.String()
	println("Client ID: ", clientID)

	client := graphql.NewClient(url, http.DefaultClient)

	var input IncrementWatchCountInput
	input.UserID = clientID
	input.VideoID = videoID

	resp, err := incrementWatchCount(ctx, client, input)
	if err != nil {
		fmt.Printf("GraphQL mutation error: %v\n", err)
		return err
	}

	println("Watch count incremented to: ", resp.IncrementWatchCount.WatchCount)
	return nil
}
