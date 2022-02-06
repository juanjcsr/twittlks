package lks

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type LksConfig struct {
	HistoryFile string
	UserID      string
	LastPage    string
	Count       int
	viperConfig *viper.Viper
}

const (
	configFileName   = "app.history_file"
	configTotalCount = "tuits.total_count"
	configLastPage   = "tuits.last_page"
	configUserID     = "tuits.user_id"
)

func NewLksConfig(config *viper.Viper) *LksConfig {
	historyFn := config.GetString(configFileName)
	if historyFn == "" {
		now := time.Now()
		historyFn = fmt.Sprintf("history_%d_%d_%d.jsonl", now.Year(), now.Month(), now.Day())
	}
	userID := config.GetString(configUserID)
	lastPage := config.GetString(configLastPage)
	count := config.GetInt(configTotalCount)
	c := &LksConfig{
		HistoryFile: historyFn,
		UserID:      userID,
		LastPage:    lastPage,
		Count:       count,
		viperConfig: config,
	}

	return c
}

func (c *LksConfig) SaveLastLksState(lastPage string, nextCount int) error {
	tc := c.Count + nextCount
	c.viperConfig.Set(configLastPage, lastPage)
	c.viperConfig.Set(configTotalCount, tc)
	c.LastPage = lastPage
	c.Count = tc
	return c.viperConfig.WriteConfig()
}

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
