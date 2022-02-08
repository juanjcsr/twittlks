package lks

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

func FetchLksHistoryFromConfig(lksclient *LksClient, v *viper.Viper) error {
	// Get the latest run and state
	c := NewLksConfig(v)
	for {
		err := FetchAndSavePage(lksclient, c)
		if err != nil {
			if err.Error() == "no more results" {
				break
			}
			return err
		}
		time.Sleep(30 * time.Second)
	}
	return nil
}

func FetchAndSavePage(lksclient *LksClient, c *LksConfig) error {
	lt, err := getPagedTuits(c.UserID, c.LastPage, lksclient)
	if err != nil {
		return err
	}
	appendTuitsToFile(lt, c.HistoryFile)
	log.Printf("last count: %d", lt.Meta.ResultCount)
	c.SaveLastLksState(lt.Meta.NextToken, lt.Meta.ResultCount)
	return nil
}

func appendTuitsToFile(lt *TwitLikesWrapper, filename string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	tlList := lt.ToTuitLikeList()
	for _, tuit := range tlList {
		f.WriteString(tuit.ToJSON() + "\n")
	}
	defer f.Close()
}

func appendTuitsLikeSliceToFile(tla []TuitLike, filename string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	for _, tuit := range tla {
		f.WriteString(tuit.ToJSON() + "\n")
	}
	defer f.Close()
	return nil
}

func getPagedTuits(user string, page string, lksClient *LksClient) (*TwitLikesWrapper, error) {
	lt, err := lksClient.GetAuthedUserLikesByPage(user, page)
	if err != nil {
		return nil, err
	}
	nextPage := lt.Meta.NextToken
	total := lt.Meta.ResultCount
	if nextPage == "" || total == 0 {
		return nil, fmt.Errorf("no more results")
	}
	return lt, nil
}
