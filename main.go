package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/pkg/browser"
)

const serverURL = "https://9ce6-2806-104e-13-4fcd-587d-133e-3b9c-3f5e.ngrok.io"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//runAuth()
	// GetUserID()
	GetUserLikes("USERID")
}

type AuthServer struct {
	server       *http.Server
	client       *http.Client
	ctx          context.Context
	cancel       context.CancelFunc
	port         string
	clientID     string
	clientSecret string
}

func NewAuthServer(ctx context.Context, port int, cancel context.CancelFunc) *AuthServer {

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
	}

	mux.HandleFunc("/auth_callback", authServer.authCallbackStopEndpoint)
	return &authServer
}

func runAuth() {
	srvExitDone := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	srvExitDone.Add(1)
	s := NewAuthServer(ctx, 8080, cancel)
	s.OpenBrowserForLogin()
	s.startServer(srvExitDone)
	srvExitDone.Wait()
}

func (s *AuthServer) OpenBrowserForLogin() {
	scopes := []string{"tweet.read", "users.read", "like.read", "offline.access"}
	authURL := s.setupOAuth2URL(scopes)

	browser.OpenURL(authURL)
}

func (s *AuthServer) startServer(wg *sync.WaitGroup) {
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

func (s *AuthServer) setupOAuth2URL(scopes []string) string {
	const oauth2URL = "https://twitter.com/i/oauth2/authorize"
	responseType := "response_type=" + "code"
	clientID := "client_id=" + s.clientID
	state := "state=" + "state"
	redirectURI := "redirect_uri=" + url.QueryEscape(serverURL+"/auth_callback")
	scopeList := "scope=" + setupScopesURL(scopes)
	authURL := fmt.Sprintf("%s?%s&%s&%s&%s&%s", oauth2URL, responseType, clientID, state, redirectURI, scopeList)
	return authURL
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
	s.GetAccessToken(code)
	w.Write([]byte("OK"))
}

func (s *AuthServer) GetAccessToken(accessCode string) {
	creds := s.clientID + ":" + s.clientSecret
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(creds))
	authEndpoint := "https://api.twitter.com/2/oauth2/token"

	data := url.Values{}
	data.Set("code", accessCode)
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", s.clientID)
	//data.Set("redirect_uri", serverURL)
	data.Set("code_verifier", "stringstring")
	ed := data.Encode() + "&redirect_uri=" + serverURL + "/auth_callback"
	fmt.Println(ed)
	r, err := http.NewRequest("POST", authEndpoint, strings.NewReader(ed))
	if err != nil {
		log.Fatalln(err)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Authorization", "Basic "+encodedCredentials)

	res, err := s.client.Do(r)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.Status)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(body))
}

func setupScopesURL(scopes []string) string {
	scopeList := strings.Join(scopes, "%20")
	codeStr := "&code_challenge=" + "stringstring" + "&code_challenge_method=plain"
	urlParams := scopeList + codeStr
	return urlParams
}

func GetUserID() {
	c := &http.Client{}

	r, err := http.NewRequest("GET", "https://api.twitter.com/2/users/me", nil)
	if err != nil {
		log.Fatalln(err)
	}
	r.Header.Add("Authorization", "Bearer "+"BEARERTOKEN")
	res, err := c.Do(r)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(res.Status)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(body))

}

func GetUserLikes(userID string) {
	c := &http.Client{}
	u := fmt.Sprintf("https://api.twitter.com/2/users/%s/liked_tweets", userID)
	r, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Fatalln(err)
	}
	r.Header.Add("Authorization", "Bearer "+"BEARERTOKEN")
	res, err := c.Do(r)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(res.Status)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(body))

}
