package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/juanjcsr/twittlks/auth"
	"github.com/juanjcsr/twittlks/lks"
	"github.com/juanjcsr/twittlks/lks/db"
	"github.com/juanjcsr/twittlks/lks/s3batch"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var batchCmd = &cobra.Command{
	Use:   "batch",
	Short: "Fetch liked twitts from the current day",
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

		if err := BatchLoad(ctx, ac, v, d, dblast, s3c); err != nil {
			return err
		}
		return nil
	},
}

func BatchLoad(ctx context.Context, ac auth.AuthClient, v *viper.Viper, d *db.DBClient, dblast string, c *s3batch.S3Client) error {
	lksClient := lks.NewLKSClient(ac, v)
	last, err := GetLastWeekLikedTwits(lksClient, v, dblast)
	if err != nil {
		return err
	}
	if dblast != last {
		_, err = SaveLikedToDB(ctx, lksClient.GetConfigCurrentPartFilename(), false, v, d)
		if err != nil {
			return err
		}
		err = c.UploadFile(ctx, "twittlks", "part_twitts", lksClient.GetConfigCurrentPartFilename())
		if err != nil {
			return err
		}
	}
	return nil
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
