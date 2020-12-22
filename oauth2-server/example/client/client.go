package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	authServerURL = "http://10.0.0.91:30011"
)

var (
	config = oauth2.Config{
		ClientID:     "PtMpAgwHIi.keti-openfx",
		ClientSecret: "R2vx9i1nG",
		RedirectURL:  "http://10.0.0.91:9095/oauth2",
		Scopes:       []string{"admin"},
		Endpoint: oauth2.Endpoint{
			AuthURL:   authServerURL + "/authorize",
			TokenURL:  authServerURL + "/token",
			AuthStyle: 0,
		},
	}
	globalToken *oauth2.Token // Non-concurrent security
)

type access_info struct {
	AccessAPI  []string
	client_id  string
	expires_in int
	user_id    string
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		u := config.AuthCodeURL("xyz")
		http.Redirect(w, r, u, http.StatusFound)
	})

	http.HandleFunc("/oauth2", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		state := r.Form.Get("state")
		if state != "xyz" {
			http.Error(w, "State invalid", http.StatusBadRequest)
			return
		}
		code := r.Form.Get("code")
		if code == "" {
			http.Error(w, "Code not found", http.StatusBadRequest)
			return
		}

		log.Printf("code : %+v", code)
		log.Printf("config : %+v", config)

		token, err := config.Exchange(context.Background(), code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		globalToken = token

		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(token)
	})

	http.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		if globalToken == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		globalToken.Expiry = time.Now()
		token, err := config.TokenSource(context.Background(), globalToken).Token()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		globalToken = token
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(token)
	})

	http.HandleFunc("/try", func(w http.ResponseWriter, r *http.Request) {
		if globalToken == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		resp, err := http.Get(fmt.Sprintf("%s/verify?access_token=%s", authServerURL, globalToken.AccessToken))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()

		// var target access_info
		// json.NewDecoder(resp.Body).Decode(&target)
		// log.Println(target.AccessAPI)

		io.Copy(w, resp.Body)
	})

	http.HandleFunc("/pwd", func(w http.ResponseWriter, r *http.Request) {
		token, err := config.PasswordCredentialsToken(context.Background(), "test", "test")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		globalToken = token
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(token)
	})

	http.HandleFunc("/client", func(w http.ResponseWriter, r *http.Request) {
		cfg := clientcredentials.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			TokenURL:     config.Endpoint.TokenURL,
		}

		token, err := cfg.Token(context.Background())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(token)
	})

	log.Println("Client is running at 9095 port.Please open http://localhost:9095")
	log.Fatal(http.ListenAndServe(":9095", nil))
}
