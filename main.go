package main

import (
	"context"
	"fmt"
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

	runAuth()
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
	scopes := []string{"tweet.read", "users.read", "like.read"}
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
	q := r.URL.Query().Get("code")
	if q == "123" {
		s.cancel()
	}
	w.Write([]byte("OK"))
}

// func (s *AuthServer) GetAccessToken(accessCode string) {
// 	creds := s.clientID + ":" + s.clientSecret
// 	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(creds))

// 	http.Post()

// }

func setupScopesURL(scopes []string) string {
	scopeList := strings.Join(scopes, "%20")
	codeStr := "&code_challenge=" + "stringstring" + "&code_challenge_method=plain"
	urlParams := scopeList + codeStr
	return urlParams
}
