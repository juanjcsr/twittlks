package lks

import (
	"log"
	"time"
)

func (l *LksClient) FetchLksCurrentWeekFromConfig(lastliked string) (string, error) {
	// lastLiked := ""
	newTuits := []TuitLike{}
	l.config.LastPage = ""
	res, err := FetchFromPage(l, l.config, &newTuits, false, lastliked)
	if err != nil {
		return "", err
	}

	if err = appendTuitsLikeSliceToFile(newTuits, l.GetConfigCurrentPartFilename()); err != nil {
		return "", err
	}
	if len(*res) > 0 {
		lastliked = (*res)[0].ID
	}
	return lastliked, nil
}

func FetchFromPage(lksclient *LksClient, c *LksConfig, tl *[]TuitLike, found bool, lastliked string) (*[]TuitLike, error) {
	if found {
		return tl, nil
	}
	lt, err := getPagedTuits(c.UserID, c.LastPage, lksclient)
	if err != nil {
		return tl, err
	}
	serverTL := lt.ToTuitLikeList()
	for _, t := range serverTL {
		if t.ID == lastliked {
			log.Println("got to previous tuit")
			found = true
			break
		}
		*tl = append(*tl, t)
	}
	c.LastPage = lt.Meta.NextToken
	time.Sleep(5 * time.Second)
	tl, err = FetchFromPage(lksclient, c, tl, found, lastliked)
	return tl, err
}
