package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/juanjcsr/twittlks/auth"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(batchCmd)
	rootCmd.AddCommand(fullCmd)
}

var rootCmd = &cobra.Command{
	Use:   "twittlks",
	Short: "Twittlks is a Twitter like's fetcher",
}

func Execute() {
	err := godotenv.Load()
	viper.SetConfigName("tokens")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Load config and access tokens from file

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func setupViperConfig() (*auth.AccessTokens, error) {
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
