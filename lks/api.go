package lks

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/juanjcsr/twittlks/auth"
)

type TuitLike struct {
	Source               string               `json:"source,omitempty"`
	AuthorID             string               `json:"author_id,omitempty"`
	Attachments          Attachments          `json:"attachments,omitempty"`
	CreatedAt            time.Time            `json:"created_at,omitempty"`
	PossiblySensitive    bool                 `json:"possibly_sensitive,omitempty"`
	GeoID                GeoID                `json:"geo,omitempty"`
	Entities             Entities             `json:"entities,omitempty"`
	Text                 string               `json:"text,omitempty"`
	ID                   string               `json:"id,omitempty" bun:",pk"`
	ConversationID       string               `json:"conversation_id,omitempty"`
	Lang                 string               `json:"lang,omitempty"`
	ReplySettings        string               `json:"reply_settings,omitempty"`
	ReferencedTweets     []ReferencedTweetsID `json:"referenced_tweets,omitempty"`
	ReferencedTweetsList []ReferencedTweets   `json:"referenced_tweets_list"`
	InReplyToUserID      string               `json:"in_reply_to_user_id,omitempty"`
	Author               Users                `json:"author"`
	MediaData            []Media              `json:"media"`
	Places               Place                `json:"place"`
	PlaceID              string
	Raw                  json.RawMessage
}

type Data struct {
	Source            string               `json:"source,omitempty"`
	AuthorID          string               `json:"author_id,omitempty"`
	Attachments       Attachments          `json:"attachments,omitempty"`
	CreatedAt         time.Time            `json:"created_at,omitempty"`
	PossiblySensitive bool                 `json:"possibly_sensitive,omitempty"`
	GeoID             GeoID                `json:"geo,omitempty"`
	Entities          Entities             `json:"entities,omitempty"`
	Text              string               `json:"text,omitempty"`
	ID                string               `json:"id,omitempty"`
	ConversationID    string               `json:"conversation_id,omitempty"`
	Lang              string               `json:"lang,omitempty"`
	ReplySettings     string               `json:"reply_settings,omitempty"`
	ReferencedTweets  []ReferencedTweetsID `json:"referenced_tweets,omitempty"`
	InReplyToUserID   string               `json:"in_reply_to_user_id,omitempty"`
}

type TwitLikesWrapper struct {
	Data     []Data   `json:"data,omitempty"`
	Includes Includes `json:"includes,omitempty"`
	Meta     Meta     `json:"meta,omitempty"`
}
type Attachments struct {
	MediaKeys []string `json:"media_keys,omitempty"`
}

type Urls struct {
	Start       int    `json:"start,omitempty"`
	End         int    `json:"end,omitempty"`
	URL         string `json:"url,omitempty"`
	ExpandedURL string `json:"expanded_url,omitempty"`
	DisplayURL  string `json:"display_url,omitempty"`
	Status      int    `json:"status,omitempty"`
	UnwoundURL  string `json:"unwound_url,omitempty"`
}

type Annotations struct {
	Start          int     `json:"start,omitempty"`
	End            int     `json:"end,omitempty"`
	Probability    float64 `json:"probability,omitempty"`
	Type           string  `json:"type,omitempty"`
	NormalizedText string  `json:"normalized_text,omitempty"`
}

type ReferencedTweetsID struct {
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
}
type Mentions struct {
	Start        int    `json:"start,omitempty"`
	End          int    `json:"end,omitempty"`
	Username     string `json:"username,omitempty"`
	ID           string `json:"id,omitempty"`
	MentionsType string
}

type Entities struct {
	Mentions    []Mentions    `json:"mentions,omitempty"`
	Urls        []Urls        `json:"urls,omitempty"`
	Annotations []Annotations `json:"annotations,omitempty"`
	Description Description   `json:"description,omitempty"`
	URL         URL           `json:"url,omitempty"`
}

type Media struct {
	DurationMs      int    `json:"duration_ms,omitempty"`
	Height          int    `json:"height,omitempty"`
	Width           int    `json:"width,omitempty"`
	PreviewImageURL string `json:"preview_image_url,omitempty"`
	MediaKey        string `json:"media_key,omitempty"`
	Type            string `json:"type,omitempty"`
	URL             string `json:"url,omitempty"`
	TuitID          string
	TuitLike        *TuitLike `bun:"rel:belongs-to,join:tuit_id=id"`
}
type URL struct {
	Urls []Urls `json:"urls,omitempty"`
}

type PublicMetrics struct {
	FollowersCount int `json:"followers_count,omitempty"`
	FollowingCount int `json:"following_count,omitempty"`
	TweetCount     int `json:"tweet_count,omitempty"`
	ListedCount    int `json:"listed_count,omitempty"`
}

type Hashtags struct {
	Start int    `json:"start,omitempty"`
	End   int    `json:"end,omitempty"`
	Tag   string `json:"tag,omitempty"`
}
type Description struct {
	Hashtags []Hashtags `json:"hashtags,omitempty"`
	Mentions []Mentions `json:"mentions,omitempty"`
}

type Users struct {
	Username        string        `json:"username,omitempty"`
	PinnedTweetID   string        `json:"pinned_tweet_id,omitempty"`
	Description     string        `json:"description,omitempty"`
	URL             string        `json:"url,omitempty"`
	ProfileImageURL string        `json:"profile_image_url,omitempty"`
	Protected       bool          `json:"protected,omitempty"`
	Location        string        `json:"location,omitempty"`
	ID              string        `json:"id,omitempty" bun:",pk"`
	CreatedAt       time.Time     `json:"created_at,omitempty"`
	Verified        bool          `json:"verified,omitempty"`
	Entities        Entities      `json:"entities,omitempty"`
	PublicMetrics   PublicMetrics `json:"public_metrics,omitempty"`
	Name            string        `json:"name,omitempty"`
	TuitLike        []*TuitLike   `bun:"rel:has-many,join:id=author"`

	Mentions []Mentions `bun:"rel:has-many,join:id="`
}

type ReferencedTweets struct {
	Entities             Entities             `json:"entities"`
	Source               string               `json:"source"`
	ReferencedTweetsList []ReferencedTweetsID `json:"referenced_tweets"`
	AuthorID             string               `json:"author_id"`
	InReplyToUserID      string               `json:"in_reply_to_user_id"`
	CreatedAt            time.Time            `json:"created_at"`
	PossiblySensitive    bool                 `json:"possibly_sensitive"`
	Text                 string               `json:"text"`
	ID                   string               `json:"id"`
	ConversationID       string               `json:"conversation_id"`
	Lang                 string               `json:"lang"`
	ReplySettings        string               `json:"reply_settings"`
}
type Includes struct {
	Media  []Media            `json:"media"`
	Users  []Users            `json:"users"`
	Tweets []ReferencedTweets `json:"tweets"`
	Places []Place            `json:"places,omitempty"`
}
type Meta struct {
	ResultCount   int    `json:"result_count"`
	NextToken     string `json:"next_token"`
	PreviousToken string `json:"previous_token"`
}

type GeoID struct {
	PlaceID string `json:"place_id"`
}

type Place struct {
	CountryCode string      `json:"country_code"`
	ID          string      `json:"id" bun:",pk"`
	Geo         Geo         `json:"geo"`
	Country     string      `json:"country"`
	FullName    string      `json:"full_name"`
	Name        string      `json:"name"`
	PlaceType   string      `json:"place_type"`
	TuitLike    []*TuitLike `bun:"has-many,join:id=place_id"`
}

type Geo struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
	Bbox        []float64 `json:"bbox"`
}

func NewLKSClient(ac auth.AuthClient) *LksClient {
	return &LksClient{
		client: ac,
	}
}

func (t *TwitLikesWrapper) ToTuitLikeList() []TuitLike {
	tlList := []TuitLike{}
	for _, tuit := range t.Data {
		tl := TuitLike{}
		for _, user := range t.Includes.Users {
			if tuit.AuthorID == user.ID {
				tl.Author = user
			}
		}

		for _, media := range t.Includes.Media {
			for _, tm := range tuit.Attachments.MediaKeys {
				if tm == media.MediaKey {
					tl.MediaData = append(tl.MediaData, media)
				}
			}
		}

		for _, ts := range t.Includes.Tweets {
			for _, rf := range tuit.ReferencedTweets {
				if rf.ID == ts.ID {
					tl.ReferencedTweetsList = append(tl.ReferencedTweetsList, ts)
				}
			}
		}

		for _, p := range t.Includes.Places {
			if p.ID == tuit.GeoID.PlaceID {
				tl.Places = p
			}
		}
		tl.Source = tuit.Source
		tl.AuthorID = tuit.AuthorID
		tl.Attachments = tuit.Attachments
		tl.CreatedAt = tuit.CreatedAt
		tl.PossiblySensitive = tuit.PossiblySensitive
		tl.GeoID = tuit.GeoID
		tl.Entities = tuit.Entities
		tl.Text = tuit.Text
		tl.ID = tuit.ID
		tl.ConversationID = tuit.ConversationID
		tl.Lang = tuit.Lang
		tl.ReplySettings = tuit.ReplySettings
		tl.ReferencedTweets = tuit.ReferencedTweets
		tl.InReplyToUserID = tuit.InReplyToUserID

		tlList = append(tlList, tl)
	}
	return tlList
}

func (t *TuitLike) ToJSON() string {
	b, err := json.Marshal(t)
	if err != nil {
		fmt.Println(err)
	}
	return string(b)
}
