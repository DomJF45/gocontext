package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

// replicate http handler
func main() {
	start := time.Now()
	ctx := context.WithValue(context.Background(), "foo", "bar")
	userID := 10
	val, err := fetchUserData(ctx, userID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Result: ", val)
	fmt.Println("Took: ", time.Since(start))
}

type Response struct {
	Value int
	Error error
}

func fetchUserData(ctx context.Context, userID int) (int, error) {
	val := ctx.Value("foo")
	fmt.Println(val)
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	resch := make(chan Response)

	go func() {
		val, err := fetchSlowStuff()
		resch <- Response{
			Value: val,
			Error: err,
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return -1, fmt.Errorf("fetching data form third party took too long")
		case resp := <-resch:
			return resp.Value, resp.Error
		}
	}
}

// minimize or reduce functions that take a long time
func fetchSlowStuff() (int, error) {
	time.Sleep(time.Millisecond * 200)
	return 666, nil
}
