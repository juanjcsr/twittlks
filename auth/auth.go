package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/browser"
)

const serverURL = "https://a19a-200-52-55-42.ngrok.io"

type AuthServer struct {
	server       *http.Server
	client       *http.Client
	ctx          context.Context
	cancel       context.CancelFunc
	port         string
	clientID     string
	clientSecret string
	scopes       []string
	Tokens       AccessTokens
}

type AuthClient struct {
	client       *http.Client
	tokens       AccessTokens
	clientID     string
	clientSecret string
}

type AccessTokens struct {
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	AccessToken  string    `json:"access_token"`
	Scope        string    `json:"scope"`
	RefreshToken string    `json:"refresh_token"`
	GrantedDate  time.Time `json:"granted_date"`
	Expired      bool
}

func NewAuthServer(ctx context.Context, port int, scopes []string, cancel context.CancelFunc) *AuthServer {

	p := strconv.Itoa(port)

	mux := http.NewServeMux()
	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + p,
	}
	authServer := AuthServer{
		server:       srv,
		ctx:          ctx,
		cancel:       cancel,
		port:         p,
		clientID:     os.Getenv("CLIENT_ID"),
		clientSecret: os.Getenv("CLIENT_SECRET"),
		client:       &http.Client{},
		scopes:       scopes,
	}

	mux.HandleFunc("/auth_callback", authServer.authCallbackStopEndpoint)
	return &authServer
}

func NewAuthClient(tokens AccessTokens) *AuthClient {

	ac := &AuthClient{
		tokens:       tokens,
		client:       &http.Client{},
		clientID:     os.Getenv("CLIENT_ID"),
		clientSecret: os.Getenv("CLIENT_SECRET"),
	}
	if tokens.Expired {
		err := ac.getRefreshToken()
		if err != nil {
			log.Panicln(err)
		}
	}
	return ac
}

func (c *AuthClient) Get(URL string, params url.Values) (*http.Response, error) {
	r, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Add("Authorization", "Bearer "+c.tokens.AccessToken)
	res, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("not authorized")
	}
	return res, nil
}

func (c *AuthClient) getRefreshToken() error {
	t, err := getTokens(c.clientID, c.clientSecret, "", &c.tokens, true)
	if err != nil {
		return err
	}
	c.tokens = *t
	return nil
}

func (s *AuthServer) OpenBrowserForLogin() {
	authURL := s.setupOAuth2URL(s.scopes)

	browser.OpenURL(authURL)
}
func (s *AuthServer) setupOAuth2URL(scopes []string) string {
	const oauth2URL = "https://twitter.com/i/oauth2/authorize"
	responseType := "response_type=" + "code"
	clientID := "client_id=" + s.clientID
	state := "state=" + "state"
	redirectURI := "redirect_uri=" + url.QueryEscape(serverURL+"/auth_callback")
	scopeList := "scope=" + setupScopesURL(scopes)
	code := "code_challenge=" + "stringstring" + "&code_challenge_method=plain"
	authURL := fmt.Sprintf("%s?%s&%s&%s&%s&%s&%s", oauth2URL, responseType, clientID, state, redirectURI, scopeList, code)
	return authURL
}

func (s *AuthServer) GetRefreshToken(refreshToken string) (*AccessTokens, error) {
	tokens, err := getTokens(s.clientID, s.clientSecret, "", &s.Tokens, true)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (s *AuthServer) GetAccessToken(accessCode string) (*AccessTokens, error) {
	t, err := getTokens(s.clientID, s.clientSecret, accessCode, nil, false)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func setupScopesURL(scopes []string) string {
	scopeList := strings.Join(scopes, "%20")
	urlParams := scopeList
	return urlParams
}

func (s *AuthServer) authCallbackStopEndpoint(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "123" {
		s.cancel()
	}
	log.Println(code)
	log.Println(state)
	//verify state and get access code if successful
	ts, err := s.GetAccessToken(code)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	s.Tokens = *ts
	w.Write([]byte("OK"))
	s.cancel()

}

func (s *AuthServer) StartServer(wg *sync.WaitGroup) {
	go func() {
		defer wg.Done() // let done main fn we are done
		log.Printf("Starting server at port %s", s.port)
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("failed to run server: %v", err)
		}
	}()
	<-s.ctx.Done()
	s.server.Shutdown(s.ctx)
	log.Println("finished server")

}

func (ac *AuthClient) GetTokens() AccessTokens {
	return ac.tokens
}

func getTokens(clientID string, clientSecret string, accessCode string, tokens *AccessTokens, refresh bool) (*AccessTokens, error) {
	creds := clientID + ":" + clientSecret
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(creds))
	authEndpoint := "https://api.twitter.com/2/oauth2/token"

	data := url.Values{}
	ed := ""
	if refresh {
		data.Set("refresh_token", tokens.RefreshToken)
		data.Set("grant_type", "refresh_token")
	} else {
		data.Set("grant_type", "authorization_code")
		data.Set("code_verifier", "stringstring")
		data.Set("code", accessCode)
		ed = "&redirect_uri=" + serverURL + "/auth_callback"
	}
	data.Set("client_id", clientID)

	r, err := http.NewRequest("POST", authEndpoint, strings.NewReader(data.Encode()+ed))
	if err != nil {
		return nil, err
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Authorization", "Basic "+encodedCredentials)
	c := http.Client{}
	res, err := c.Do(r)
	if err != nil {
		return nil, err
	}
	// log.Println(res.Status)
	defer res.Body.Close()
	t := &AccessTokens{}
	t.GrantedDate = time.Now()
	t.Expired = false
	// body, err := ioutil.ReadAll(res.Body)
	json.NewDecoder(res.Body).Decode(t)
	if err != nil {
		return nil, err
	}
	// log.Println(t)

	return t, nil
}
