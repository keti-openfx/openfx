package main

import (
	"bytes"
	crand "crypto/rand"
	"database/sql"
	"encoding/json"
	"log"
	"math"
	"math/big"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-session/session"
	"gopkg.in/oauth2.v4/errors"
	"gopkg.in/oauth2.v4/manage"
	"gopkg.in/oauth2.v4/models"
	"gopkg.in/oauth2.v4/server"
	"gopkg.in/oauth2.v4/store"
)

// userDB mysql
var db *sql.DB
var err error

func main() {
	// open mysql user db
	db, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/OAuth2") // mysql pod에서 넣는 법

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	manager := manage.NewDefaultManager()
	//manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// token store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// generate jwt access token
	// manager.MapAccessGenerate(generates.NewJWTAccessGenerate([]byte("00000000"), jwt.SigningMethodHS512))

	clientStore := store.NewClientStore()
	manager.MapClientStorage(clientStore)

	srv := server.NewServer(server.NewConfig(), manager)

	srv.SetAllowGetAccessRequest(true)

	srv.SetClientInfoHandler(server.ClientFormHandler)

	srv.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
		return username, nil
	})

	srv.SetUserAuthorizationHandler(userAuthorizeHandler)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		//log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		//log.Println("Response Error:", re.Error.Error())
	})

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/signup", signupPage)
	http.HandleFunc("/", homePage)

	http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		store, err := session.Start(r.Context(), w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var form url.Values
		if v, ok := store.Get("ReturnUri"); ok {
			form = v.(url.Values)
		}
		r.Form = form

		store.Delete("ReturnUri")
		store.Save()

		err = srv.HandleAuthorizeRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		// get header dadata
		username := r.FormValue("username")
		password := r.FormValue("password")
		clientID := r.FormValue("client_id")

		var databaseUsername string
		var databasePassword string
		var scope string

		err = db.QueryRow("SELECT scope FROM clients WHERE client_id=?", clientID).Scan(&scope)
		if err != nil {
			w.Write([]byte("This is unregistered Clients information."))
			return
		}

		err := db.QueryRow("SELECT username, password FROM users WHERE username=? and scope=?", username, scope).Scan(&databaseUsername, &databasePassword)
		if err != nil {
			// 알람 메세지 추가
			w.Write([]byte("This is unregistered member information."))
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
		if err != nil {
			w.Write([]byte("The ID or password are not correct."))
			http.Redirect(w, r, "/login", 301)
			return
		}

		r.Form.Set("username", username)
		r.Form.Set("scope", scope)

		var authority string
		err = db.QueryRow("SELECT authority FROM auths WHERE (username=? AND client_id=?)", username, clientID).Scan(&authority)
		if err != nil {
			if err == sql.ErrNoRows {
				_, err = db.Exec("INSERT INTO auths(client_id, username, authority) VALUES(?, ?, ?)", clientID, username, "dev")
				if err != nil {
					http.Error(w, "Server error, unable to create your account.", 500)
					return
				}
			} else {
				w.Write([]byte("can't login because server error."))
				return
			}
		}

		err = srv.HandleTokenRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		token, err := srv.ValidationBearerToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var authority string
		err = db.QueryRow("SELECT authority FROM auths WHERE (username=? AND client_id=?)", token.GetUserID(), token.GetClientID()).Scan(&authority)
		if err != nil {
			// 알람 메세지 추가
			w.Write([]byte("can't login because server error."))
			return
		}

		data := map[string]interface{}{
			"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
			"client_id":  token.GetClientID(),
			"user_id":    token.GetUserID(),
			"scope":      token.GetScope(),
			"grade":      authority,
		}

		log.Printf("[Authorized] Verify / User : %v, Namespaces : %v", token.GetUserID(), token.GetScope())

		e := json.NewEncoder(w)
		w.Header().Set("Content-Type", "application/json")
		e.SetIndent("", "  ")
		e.Encode(data)
	})

	http.HandleFunc("/registerclient", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			http.ServeFile(res, req, "./static/createClient.html")
			return
		}

		scope := req.FormValue("scope")
		domain := req.FormValue("Domain")

		var databaseclientscope string
		var databaseclientid string
		var databaseclientsecret string
		var databaseclientdomain string

		manager.MapClientStorage(clientStore)
		srv.Manager = manager

		err := db.QueryRow("SELECT * FROM clients WHERE scope=?", scope).Scan(&databaseclientid, &databaseclientsecret, &databaseclientdomain, &databaseclientscope)
		switch {
		case err == sql.ErrNoRows: // 결과가 없는 경우
			clientDomain := req.FormValue("Domain")
			clientScope := req.FormValue("scope")
			clientID := RandomString(10) + ".keti-openfx"
			clientSecret := RandomString(10)

			clientStore.Set(clientID, &models.Client{
				ID:     clientID,
				Secret: clientSecret,
				Domain: clientDomain,
			})

			_, err = db.Exec("INSERT INTO clients(client_id, client_secret, domain, scope) VALUES(?, ?, ?, ?)", clientID, clientSecret, clientDomain, clientScope)
			if err != nil {
				http.Error(res, "Server error, unable to create your account.", 500)
				return
			}

			resp, err := http.Post("http://10.0.0.91:31113/api/createns/"+clientScope, "text/plain", bytes.NewBufferString(""))
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			data := map[string]interface{}{
				"clientID":      clientID,
				"clientSecret":  clientSecret,
				"Client Domain": clientDomain,
			}

			// curl -X POST 10.0.0.91:31113/api/createns/{namespaces} 추가 필요

			e := json.NewEncoder(res)
			res.WriteHeader(200)
			res.Header().Set("Content-Type", "application/json")
			e.SetIndent("", "  ")
			e.Encode(data)

		case err != nil:
			res.Write([]byte("can't login because server error."))
			return

		default: // 있는 경우 갱신 필요
			_, err = db.Exec("UPDATE clients SET domain=?, scope=? where client_id=?", domain, scope, databaseclientid)
			if err != nil {
				http.Error(res, "Server error, clients Database Insert error ", 500)
				return
			}

			data := map[string]interface{}{
				"clientID":      databaseclientid,
				"clientSecret":  databaseclientsecret,
				"Client Domain": databaseclientdomain,
			}
			e := json.NewEncoder(res)
			res.WriteHeader(200)
			res.Header().Set("Content-Type", "application/json")
			e.SetIndent("", "  ")
			e.Encode(data)
		}
	})

	log.Println("Server is running at 9096 port.")
	log.Fatal(http.ListenAndServe(":9096", nil))
}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		return
	}

	uid, ok := store.Get("LoggedInUserID")
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}

		store.Set("ReturnUri", r.Form)
		store.Save()

		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	userID = uid.(string)
	store.Delete("LoggedInUserID")
	store.Save()
	return
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		if r.Form == nil {
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		// 로그인 로직 추가
		username := r.FormValue("username")
		password := r.FormValue("password")
		clientid := r.FormValue("client_id")
		clientsecret := r.FormValue("client_secret")

		var databaseUsername string
		var databasePassword string
		var databaseclientscope string

		// client 정보 체크
		err := db.QueryRow("select scope from clients where client_id=? and client_secret=?", clientid, clientsecret).Scan(&databaseclientscope)
		switch {
		case err == sql.ErrNoRows:
			http.Error(w, "Client information that is not registered.", 500)
			return
		case err != nil:
			http.Error(w, "Server error, unable to create your account.", 500)
			return
		default:
			break
		}

		err = db.QueryRow("SELECT username, password FROM users WHERE username=? and scope=?", username, databaseclientscope).Scan(&databaseUsername, &databasePassword)
		if err != nil {
			// 알람 메세지 추가
			w.Write([]byte("can't login because server error."))
			http.Redirect(w, r, "/login", 301)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
		if err != nil {
			// 알람 메세지 추가
			w.Write([]byte("The ID or password are not correct."))
			http.Redirect(w, r, "/login", 301)
			return
		}

		store.Set("LoggedInUserID", r.Form.Get("username"))
		store.Save()

		w.Header().Set("Location", "/auth")
		w.WriteHeader(http.StatusFound)
		return
	}
	outputHTML(w, r, "static/login.html")
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(nil, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := store.Get("LoggedInUserID"); !ok {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	outputHTML(w, r, "static/auth.html")
}

func outputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	seed, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	rand.Seed(seed.Int64())

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func signupPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "./static/signup.html")
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")
	company := req.FormValue("company")
	email := req.FormValue("email")
	clientid := req.FormValue("client_id")
	clientsecret := req.FormValue("client_secret")

	//client 정보 체크
	var databaseclientid string
	var databaseclientsecret string
	var databaseclientdomain string
	var databaseclientscope string

	err := db.QueryRow("select client_id, client_secret, domain, scope from clients where client_id=? and client_secret=?", clientid, clientsecret).Scan(&databaseclientid, &databaseclientsecret, &databaseclientdomain, &databaseclientscope)
	switch {
	case err == sql.ErrNoRows:
		http.Error(res, "Client information that is not registered.", 500)
		return
	case err != nil:
		break
	default:
		break
	}

	// 해당 쿼리문에 정보 추가 필요
	var user string

	err = db.QueryRow("SELECT username FROM users WHERE username=? and scope=?", username, databaseclientscope).Scan(&user)

	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, "Server error,  passwoard encoding problem", 500)
			return
		}
		_, err = db.Exec("INSERT INTO users(scope, username, password, company, email) VALUES(?, ?, ?, ?, ?)", databaseclientscope, username, hashedPassword, company, email)
		if err != nil {
			http.Error(res, "Server error, database input data error ", 500)
			return
		}
		//json 결과 메세지 응답
		enc := json.NewEncoder(res)
		res.WriteHeader(200)
		res.Header().Set("Content-Type", "application/json")
		enc.Encode("User created")
		return

	case err != nil:
		http.Error(res, "Server error, unable to create your account.", 500)
		return
	default:
		http.Error(res, "User who has already registered as a member", 500)
		return
	}
}

func homePage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "./static/index.html")
}
