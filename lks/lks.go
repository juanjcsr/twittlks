package lks

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"time"

	"github.com/juanjcsr/twittlks/auth"
	"github.com/juanjcsr/twittlks/lks/s3batch"
	"github.com/spf13/viper"
)

const maxResults = "99"

type LksClient struct {
	client   auth.AuthClient
	config   *LksConfig
	s3api    *s3batch.S3Client
	s3access bool
}

func NewLKSClient(ac auth.AuthClient, v *viper.Viper) *LksClient {
	c := NewLksConfig(v)
	s3api, err := s3batch.NewAWSClient("twittlks")
	s3access := true
	if err != nil {
		s3access = false
	}
	return &LksClient{
		client:   ac,
		config:   c,
		s3api:    s3api,
		s3access: s3access,
	}
}

func (l *LksClient) UploadPartFileToS3(bucket string, path string) error {
	if !l.s3access {
		log.Println("no s3 access")
		return nil
	}
	err := l.s3api.UploadFile(context.TODO(), bucket, path, l.GetConfigCurrentPartFilename())
	if err != nil {
		return err
	}
	return nil
}

func (l *LksClient) GetConfigCurrentPartFilename() string {
	return "part_" + l.config.HistoryFile
}

func (l *LksClient) GetAuthedUserLikesByPage(userID string, page string) (*TwitLikesWrapper, error) {
	params := url.Values{}
	if page != "" {
		params.Set("pagination_token", page)
	}
	params.Set("max_results", maxResults)
	params.Set("user.fields", "created_at,description,entities,id,location,name,pinned_tweet_id,profile_image_url,protected,public_metrics,url,username,verified,withheld")
	params.Set("place.fields", "country,country_code,full_name,geo,id,name,place_type")
	params.Set("media.fields", "duration_ms,height,media_key,preview_image_url,type,url,width,alt_text")
	params.Set("tweet.fields", "attachments,author_id,conversation_id,created_at,entities,geo,id,in_reply_to_user_id,lang,possibly_sensitive,referenced_tweets,reply_settings,source,text,withheld")
	params.Set("expansions", "attachments.poll_ids,attachments.media_keys,author_id,entities.mentions.username,geo.place_id,in_reply_to_user_id,referenced_tweets.id,referenced_tweets.id.author_id")
	u := fmt.Sprintf("https://api.twitter.com/2/users/%s/liked_tweets?%s", userID, params.Encode())
	res, err := l.client.Get(u, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println(string(body))
	s := string(body)
	s = Decode(s)
	lt := &TwitLikesWrapper{}
	// json.NewDecoder([]byte(s)).Decode(lt)
	err = json.Unmarshal([]byte(s), lt)
	if err != nil {
		log.Println(err)
	}
	return lt, nil
}

func (l *LksClient) GetAuthedUserLikes(userID string) (*TwitLikesWrapper, error) {
	params := url.Values{}
	params.Set("max_results", "25")
	params.Set("user.fields", "created_at,description,entities,id,location,name,pinned_tweet_id,profile_image_url,protected,public_metrics,url,username,verified,withheld")
	params.Set("place.fields", "country,country_code,full_name,geo,id,name,place_type")
	params.Set("media.fields", "duration_ms,height,media_key,preview_image_url,type,url,width,alt_text")
	params.Set("tweet.fields", "attachments,author_id,conversation_id,created_at,entities,geo,id,in_reply_to_user_id,lang,possibly_sensitive,referenced_tweets,reply_settings,source,text,withheld")
	params.Set("expansions", "attachments.poll_ids,attachments.media_keys,author_id,entities.mentions.username,geo.place_id,in_reply_to_user_id,referenced_tweets.id,referenced_tweets.id.author_id")
	u := fmt.Sprintf("https://api.twitter.com/2/users/%s/liked_tweets?%s", userID, params.Encode())
	res, err := l.client.Get(u, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println(string(body))
	s := string(body)
	s = Decode(s)
	fmt.Println(s)
	lt := &TwitLikesWrapper{}
	// json.NewDecoder([]byte(s)).Decode(lt)
	json.Unmarshal([]byte(s), lt)
	return lt, nil
}

type LksConfig struct {
	HistoryFile   string
	UserID        string
	LastPage      string
	Count         int
	LastLikedTuit string
	viperConfig   *viper.Viper
}

const (
	configFileName    = "app.history_file"
	configTotalCount  = "tuits.total_count"
	configLastPage    = "tuits.last_page"
	configUserID      = "tuits.user_id"
	configLastLkdTuit = "tuits.last_liked_tuit"
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
	lastLikedTuit := config.GetString(configLastLkdTuit)

	c := &LksConfig{
		HistoryFile:   historyFn,
		UserID:        userID,
		LastPage:      lastPage,
		Count:         count,
		LastLikedTuit: lastLikedTuit,

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
