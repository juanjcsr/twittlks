package cmd

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/juanjcsr/twittlks/auth"
	"github.com/juanjcsr/twittlks/lks/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fullCmd = &cobra.Command{
	Use:   "full",
	Short: "Fetch all possible history of likes",
	RunE: func(cmd *cobra.Command, args []string) error {
		tokens, err := setupViperConfig()

		// If tokens are missing, run the auth process again
		if err != nil {
			log.Println(err)
			tokens = runAuth()
		}

		ac := *auth.NewAuthClient(*tokens)
		*tokens = ac.GetTokens()
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
		if err != nil {
			log.Fatalln(err)
		}
		v := viper.GetViper()
		d := OpenDB(dburl)

		if err != nil {
			log.Fatalln(err)
		}
		if err := InitialLoad(ctx, v, d); err != nil {
			return err
		}
		return nil
	},
}

func InitialLoad(ctx context.Context, v *viper.Viper, d *db.DBClient) error {
	// FETCH TUITS
	_, err := SaveLikedToDB(ctx, "fulltuits.jsonl", true, v, d)
	return err
}

func OpenDB(u string) *db.DBClient {
	d := db.OpenSQLConn(u)
	d.OpenBUN()
	return d
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
