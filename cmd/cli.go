package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/juanjcsr/twittlks/auth"
	"github.com/juanjcsr/twittlks/lks/s3batch"
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

	godotenv.Load()
	ctx := context.Background()
	// Load config and access tokens from file
	err := fetchViperConfig(ctx, "tokens.yaml")
	if err != nil {
		log.Panicln(err)
	}

	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
	}
	//save config file to s3
	err = uploadConfigFileToS3(ctx, "tokens.yaml")
	if err != nil {
		log.Printf("could not upload file to s3: %s", err)
	}
}

func fetchViperConfig(ctx context.Context, configName string) error {
	viper.SetConfigName("tokens")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	if _, err := os.Stat(configName); errors.Is(err, os.ErrNotExist) {
		log.Println("no local tokens file, fetching from s3...")
		s3c, err := s3batch.NewAWSClient("twittlks")
		if err != nil {
			return err
		}

		data, err := s3c.GetFile(ctx, s3c.Bucket, "config/tokens.yaml")
		if err != nil {
			return err
		}

		err = os.WriteFile(configName, *data, 0644)
		if err != nil {
			return err
		}
		log.Printf("writen tokens file to: %s", configName)
	}
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

func uploadConfigFileToS3(ctx context.Context, filename string) error {
	log.Println("uploading config file to s3...")
	s3c, err := s3batch.NewAWSClient("twittlks")
	if err != nil {
		return err
	}

	err = s3c.UploadFile(ctx, s3c.Bucket, "config", filename)
	if err != nil {
		return err
	}
	log.Println("finish uploading config file to s3")
	return nil
}

func setupViperConfig() (*auth.AccessTokens, error) {
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
	return &tokens, nil
}
