package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/joho/godotenv"
	"github.com/juanjcsr/twittlks/auth"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//runAuth()
	// GetUserID()
	GetUserLikes("6846262")
}

func runAuth() {
	srvExitDone := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	srvExitDone.Add(1)
	scopes := []string{"tweet.read", "users.read", "like.read", "offline.access"}
	s := auth.NewAuthServer(ctx, 8080, scopes, cancel)
	s.OpenBrowserForLogin()
	s.StartServer(srvExitDone)
	srvExitDone.Wait()
}

func GetUserID() {
	c := &http.Client{}

	r, err := http.NewRequest("GET", "https://api.twitter.com/2/users/me", nil)
	if err != nil {
		log.Fatalln(err)
	}
	r.Header.Add("Authorization", "Bearer "+"BEARERTOKEN")
	res, err := c.Do(r)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(res.Status)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(body))

}

func GetUserLikes(userID string) {
	c := &http.Client{}
	u := fmt.Sprintf("https://api.twitter.com/2/users/%s/liked_tweets", userID)
	r, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Fatalln(err)
	}
	r.Header.Add("Authorization", "Bearer "+"ZHdYWWxmNGhDNW9fZi01eXJaZVc2MFl0SENHeVVfTEluWGRKbTRBbEpKRVNXOjE2NDE2MDg1MDMwOTA6MToxOmF0OjE")
	res, err := c.Do(r)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(res.Status)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(body))

}
