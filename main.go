package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
	"github.com/juanjcsr/twittlks/auth"
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
		viper.Set("app.expires", tokens.ExpiresIn)
		viper.Set("app.token_type", tokens.TokenType)
		viper.Set("app.access_token", tokens.AccessToken)
		viper.Set("app.refresh_token", tokens.RefreshToken)
		viper.Set("app.scope", tokens.Scope)
		viper.WriteConfig()
	}
	authClient := auth.NewAuthClient(*tokens)
	GetAuthedUserLikes("6846262", *authClient)
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
	return &s.Tokens
}

func GetAuthedUserLikes(userID string, ac auth.AuthClient) {
	u := fmt.Sprintf("https://api.twitter.com/2/users/%s/liked_tweets", userID)
	res, err := ac.Get(u, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(body))
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

	expiresInStr := viper.GetString("app.expires")
	tokenType := viper.GetString("app.token_type")
	accessToken := viper.GetString("app.access_token")
	refreshToken := viper.GetString("app.refresh_token")
	scope := viper.GetString("app.scope")

	if expiresInStr == "" || tokenType == "" || accessToken == "" || refreshToken == "" {
		return nil, fmt.Errorf("no config file")
	}
	expiresIn, _ := strconv.Atoi(expiresInStr)
	tokens := auth.AccessTokens{
		TokenType:    tokenType,
		ExpiresIn:    expiresIn,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Scope:        scope,
	}

	return &tokens, err

}
