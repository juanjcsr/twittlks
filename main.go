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
	"github.com/juanjcsr/twittlks/lks/db"
	"github.com/juanjcsr/twittlks/lks/s3batch"
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

	dburl := viper.GetString("db.url")
	if err != nil {
		log.Fatalln("no url defined in tokens.yaml")
	}

	ctx := context.Background()
	s3c, err := s3batch.NewAWSClient("twittlks")
	if err != nil {
		log.Fatalln(err)
	}
	v := viper.GetViper()
	d := OpenDB(dburl)
	dblast, err := d.GetLastInsertedTuit(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	//
	// InitialLoad(ctx, v, d)

	//
	BatchLoad(ctx, ac, v, d, dblast, s3c)
}

func InitialLoad(ctx context.Context, v *viper.Viper, d *db.DBClient) error {
	// FETCH TUITS
	_, err := SaveLikedToDB(ctx, "fulltuits.jsonl", true, v, d)
	return err
}

func BatchLoad(ctx context.Context, ac auth.AuthClient, v *viper.Viper, d *db.DBClient, dblast string, c *s3batch.S3Client) {
	lksClient := lks.NewLKSClient(ac, v)
	last, err := GetLastWeekLikedTwits(lksClient, v, dblast)
	if err != nil {
		log.Println(err)
	}
	if dblast != last {
		_, err = SaveLikedToDB(ctx, lksClient.GetConfigCurrentPartFilename(), false, v, d)
		if err != nil {
			log.Fatalln(err)
		}
		// c, err := s3batch.NewAWSClient("twittlks")
		// if err != nil {
		// 	log.Fatalln(err)
		// }
		err = c.UploadFile(ctx, "twittlks", "part_twitts", lksClient.GetConfigCurrentPartFilename())
		if err != nil {
			log.Fatalln(err)
		}
	}

}

func OpenDB(u string) *db.DBClient {
	d := db.OpenSQLConn(u)
	d.OpenBUN()
	return d
}

func SaveLikedToDB(ctx context.Context, filename string, newDB bool, v *viper.Viper, d *db.DBClient) (string, error) {
	tl, err := db.ReadLineFromFile(filename)
	if err != nil {
		return "", err
	}

	if err = d.CreateTables(ctx, newDB); err != nil {
		return "", err
	}
	lastTL, err := d.SaveTuitsToDB(tl, ctx)
	viper.Set("tuits.last_saved_tuit", lastTL)
	viper.WriteConfig()
	if err != nil {
		log.Fatalf("last inserted tuit: %s, err: %s", lastTL, err)
		return lastTL, err
	}
	return lastTL, nil
}

func GetLastWeekLikedTwits(lksClient *lks.LksClient, v *viper.Viper, last string) (string, error) {
	if last == "" {
		return "", fmt.Errorf("need to load first liked tuits to db")
	}
	last, err := lksClient.FetchLksCurrentWeekFromConfig(last)
	if err != nil {
		return last, err
	}
	v.Set("tuits.last_liked_tuit", last)
	v.WriteConfig()
	return last, nil
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
