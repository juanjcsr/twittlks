package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/juanjcsr/twittlks/auth"
	"github.com/juanjcsr/twittlks/lks"
	"github.com/spf13/viper"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	tokens, err := setupViperConfig()

	if err != nil {
		log.Println(err)
		tokens = runAuth()
	}
	// authClient := auth.NewAuthClient(*tokens)
	ac := *auth.NewAuthClient(*tokens)
	*tokens = ac.GetTokens()
	viper.Set("app.expires", tokens.ExpiresIn)
	viper.Set("app.token_type", tokens.TokenType)
	viper.Set("app.access_token", tokens.AccessToken)
	viper.Set("app.refresh_token", tokens.RefreshToken)
	viper.Set("app.scope", tokens.Scope)
	viper.Set("app.granted_date", tokens.GrantedDate)
	viper.WriteConfig()

	lt, err := GetAuthedUserLikes("6846262", ac)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(lt)
}

func runAuth() *auth.AccessTokens {
	srvExitDone := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	srvExitDone.Add(1)
	scopes := []string{"tweet.read", "users.read", "like.read", "offline.access"}
	s := auth.NewAuthServer(ctx, 8080, scopes, cancel)
	s.OpenBrowserForLogin()
	s.StartServer(srvExitDone)
	srvExitDone.Wait()
	fmt.Println(s.Tokens)
	s.Tokens.GrantedDate = time.Now()
	return &s.Tokens
}

func GetAuthedUserLikes(userID string, ac auth.AuthClient) (*lks.TwitLikesWrapper, error) {
	params := url.Values{}
	params.Set("max_results", "5")
	params.Set("user.fields", "created_at,description,entities,id,location,name,pinned_tweet_id,profile_image_url,protected,public_metrics,url,username,verified,withheld")
	params.Set("place.fields", "country,country_code,full_name,geo,id,name,place_type")
	params.Set("media.fields", "duration_ms,height,media_key,preview_image_url,type,url,width,alt_text")
	params.Set("tweet.fields", "attachments,author_id,conversation_id,created_at,entities,geo,id,in_reply_to_user_id,lang,possibly_sensitive,referenced_tweets,reply_settings,source,text,withheld")
	params.Set("expansions", "attachments.poll_ids,attachments.media_keys,author_id,entities.mentions.username,geo.place_id,in_reply_to_user_id,referenced_tweets.id,referenced_tweets.id.author_id")
	u := fmt.Sprintf("https://api.twitter.com/2/users/%s/liked_tweets?%s", userID, params.Encode())
	res, err := ac.Get(u, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	// body, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	lt := &lks.TwitLikesWrapper{}
	json.NewDecoder(res.Body).Decode(lt)
	return lt, nil
}

func setupViperConfig() (*auth.AccessTokens, error) {
	viper.SetConfigName("tokens")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file, default \n", err)
		os.Exit(1)
	}

	expiresIn := viper.GetInt("app.expires")
	tokenType := viper.GetString("app.token_type")
	accessToken := viper.GetString("app.access_token")
	refreshToken := viper.GetString("app.refresh_token")
	scope := viper.GetString("app.scope")
	lastDate := viper.GetTime("app.granted_date")
	expired := false

	if expiresIn == 0 || tokenType == "" || accessToken == "" || refreshToken == "" {
		return nil, fmt.Errorf("no config file")
	}

	if lastDate.Add(time.Second * time.Duration(expiresIn)).Before(time.Now()) {
		expired = true
	}
	tokens := auth.AccessTokens{
		TokenType:    tokenType,
		ExpiresIn:    expiresIn,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Scope:        scope,
		GrantedDate:  lastDate,
		Expired:      expired,
	}

	return &tokens, err

}
