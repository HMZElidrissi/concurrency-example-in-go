package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type UserProfile struct {
	ID       int
	Comments []string
	Likes    int
	Friends  []int
}

type Response struct {
	data any
	err  error
}

func handleGetUserProfile(id int) (*UserProfile, error) {
	var (
		respch = make(chan Response, 3)
		wg     = &sync.WaitGroup{}
	)
	// 3 requests inside their own goroutines
	go getComments(id, respch, wg)
	go getLikes(id, respch, wg)
	go getFriends(id, respch, wg)

	// add 3 to the waitgroup
	wg.Add(3)
	wg.Wait()
	close(respch)

	userProfile := &UserProfile{}
	for resp := range respch {
		if resp.err != nil {
			return nil, resp.err
		}
		switch msg := resp.data.(type) {
		case int:
			userProfile.Likes = msg
		case []int:
			userProfile.Friends = msg
		case []string:
			userProfile.Comments = msg

		}
	}
	return userProfile, nil
}

func getComments(id int, respch chan Response, wg *sync.WaitGroup) {
	time.Sleep(time.Millisecond * 200)
	comments := []string{
		"Hey, that was great",
		"great content",
		"alohaaa!",
	}
	respch <- Response{
		data: comments,
		err:  nil,
	}

	wg.Done()
}

func getLikes(id int, respch chan Response, wg *sync.WaitGroup) {
	time.Sleep(time.Millisecond * 200)
	respch <- Response{
		data: 11,
		err:  nil,
	}
	wg.Done()
}

func getFriends(id int, respch chan Response, wg *sync.WaitGroup) {
	time.Sleep(time.Millisecond * 100)
	friendsIds := []int{11, 120, 15}
	respch <- Response{
		data: friendsIds,
		err:  nil,
	}
	wg.Done()
}

func main() {
	start := time.Now()
	userProfile, err := handleGetUserProfile(10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(userProfile)
	fmt.Printf("fetching user took %v\n", time.Since(start))
}
