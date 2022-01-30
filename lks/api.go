package lks

import "time"

type TuitLike struct {
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
	Author            Users                `json:"author"`
	MediaData         []Media              `json:"media"`
	Places            Place                `json:"place"`
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
	Start    int    `json:"start,omitempty"`
	End      int    `json:"end,omitempty"`
	Username string `json:"username,omitempty"`
	ID       string `json:"id,omitempty"`
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
	ID              string        `json:"id,omitempty"`
	CreatedAt       time.Time     `json:"created_at,omitempty"`
	Verified        bool          `json:"verified,omitempty"`
	Entities        Entities      `json:"entities,omitempty"`
	PublicMetrics   PublicMetrics `json:"public_metrics,omitempty"`
	Name            string        `json:"name,omitempty"`
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
	ResultCount int    `json:"result_count"`
	NextToken   string `json:"next_token"`
}

type GeoID struct {
	PlaceID string `json:"place_id"`
}

type Place struct {
	CountryCode string `json:"country_code"`
	ID          string `json:"id"`
	Geo         Geo    `json:"geo"`
	Country     string `json:"country"`
	FullName    string `json:"full_name"`
	Name        string `json:"name"`
	PlaceType   string `json:"place_type"`
}

type Geo struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
	Bbox        []float64 `json:"bbox"`
}
