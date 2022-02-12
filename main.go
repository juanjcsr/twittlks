package main

import (
	"context"
	"fmt"
	"log"
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

	// Load config and access tokens from file
	tokens, err := setupViperConfig()

	// If tokens are missing, run the auth process again
	if err != nil {
		log.Println(err)
		tokens = runAuth()
	}

	ac := *auth.NewAuthClient(*tokens)
	*tokens = ac.GetTokens()

	// Rewrite the tokens in case they were refreshed
	viper.Set("app.expires", tokens.ExpiresIn)
	viper.Set("app.token_type", tokens.TokenType)
	viper.Set("app.access_token", tokens.AccessToken)
	viper.Set("app.refresh_token", tokens.RefreshToken)
	viper.Set("app.scope", tokens.Scope)
	viper.Set("app.granted_date", tokens.GrantedDate)
	viper.WriteConfig()
	v := viper.GetViper()
	// Use the authed http client to create a new LKS client
	lksClient := lks.NewLKSClient(ac)

	err = lks.FetchLksHistoryFromConfig(lksClient, v)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("finished extracting tuits")

	// err = lks.FetchLksCurrentWeekFromConfig(lksClient, v)
	// if err != nil {
	// 	log.Println(err)
	// }
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
