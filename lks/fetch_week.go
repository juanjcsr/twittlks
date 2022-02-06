package lks

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

func FetchLksCurrentWeekFromConfig(lksclient *LksClient, v *viper.Viper) error {
	c := NewLksConfig(v)
	if c.LastLikedTuit == "" {
		return fmt.Errorf("no prev. tuit history")
	}
	log.Println(c.LastLikedTuit)
	newTuits := []TuitLike{}
	c.LastPage = ""
	res, err := FetchFromPage(lksclient, c, &newTuits, false)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(len(*res))
	return nil
}

func FetchFromPage(lksclient *LksClient, c *LksConfig, tl *[]TuitLike, found bool) (*[]TuitLike, error) {
	if found {
		return tl, nil
	}
	lt, err := getPagedTuits(c.UserID, c.LastPage, lksclient)
	if err != nil {
		return tl, err
	}
	serverTL := lt.ToTuitLikeList()
	for _, t := range serverTL {
		if t.ID == c.LastLikedTuit {
			log.Println("got to previous tuit")
			found = true
			break
		}
		*tl = append(*tl, t)
		log.Printf("new tuit: %s", t.Text)
	}
	c.LastPage = lt.Meta.NextToken
	time.Sleep(5 * time.Second)
	tl, err = FetchFromPage(lksclient, c, tl, found)
	return tl, err
}
