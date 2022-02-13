package lks

import (
	"fmt"
	"log"
	"time"
)

func (l *LksClient) FetchLksCurrentWeekFromConfig() (string, error) {
	lastLiked := l.config.LastLikedTuit
	if lastLiked == "" {
		return "", fmt.Errorf("no prev. tuit history")
	}
	log.Println(lastLiked)
	newTuits := []TuitLike{}
	l.config.LastPage = ""
	res, err := FetchFromPage(l, l.config, &newTuits, false)
	if err != nil {
		return lastLiked, err
	}

	if err = appendTuitsLikeSliceToFile(newTuits, l.GetConfigCurrentPartFilename()); err != nil {
		return lastLiked, err
	}
	if len(*res) > 0 {
		lastLiked = (*res)[0].ID
	}
	return lastLiked, nil
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
