package lks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"

	"github.com/juanjcsr/twittlks/auth"
)

const maxResults = "100"

type LksClient struct {
	client auth.AuthClient
	config LksConfig
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
