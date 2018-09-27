package app

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"time"
	"github.com/jackson6/gcic-server/config"
	"github.com/jackson6/gcic-server/app/handler"
	"gopkg.in/mgo.v2"
	netContext "golang.org/x/net/context"
	//context "github.com/gorilla/context"
	"firebase.google.com/go"
	"google.golang.org/api/option"
	"strings"
	"github.com/jackson6/gcic-server/app/dao"
	"fmt"
	//"github.com/mitchellh/mapstructure"
	"github.com/gorilla/context"
	_ "gopkg.in/cq.v1"
	"database/sql"
	"github.com/googollee/go-socket.io"
	"gopkg.in/mgo.v2/bson"
)

// App has router and db instances
type App struct {
	Router *mux.Router
	MongoDB *mgo.Session
	Secret string
	Firebase *firebase.App
	GraphDB *sql.DB
	Socket *socketio.Server
	Stripe string
}

type CORSRouterDecorator struct {
	R *mux.Router
	S *socketio.Server
}

func (c *CORSRouterDecorator) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Language, Content-Type, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}
	param1 := req.URL.Query().Get("transport")
	if param1 != "" {
		rw.Header().Set("Access-Control-Allow-Credentials", "true")
		c.S.ServeHTTP(rw, req)
	} else {
		c.R.ServeHTTP(rw, req)
	}
}

// Initialize initializes the app with predefined configuration
func (a *App) Initialize(config *config.Config) {
	a.Router = mux.NewRouter()
	a.setRouters()
	a.Secret = config.SECRET
	a.Stripe = config.StripeKey

	session, err := mgo.Dial(config.MongoDB.Server)
	if err != nil {
		log.Fatal("error initializing mongo: %v\n", err)
	}
	a.MongoDB = session
	opt := option.WithCredentialsFile("config/invest-ff3f4-firebase-adminsdk-zgkg5-ae79e82ab1.json")
	app, err := firebase.NewApp(netContext.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing firebase: %v\n", err)
	}

	socket, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	//opts := make(map[string]string)
	//opts["host"] = config.Redis.Host
	//opts["port"] = config.Redis.Port
	//opts["prefix"] = config.Redis.Prefix
	//
	//socket.SetAdaptor(redis.Redis(opts))
	a.Socket = socket

	a.Firebase = app
	db, err := sql.Open("neo4j-cypher", "http://localhost:11001")
	if err != nil {
		log.Fatal(err)
	}
	a.GraphDB = db
	defer db.Close()
}

// Auth middleware
func (a *App) ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				client, err := a.Firebase.Auth(netContext.Background())
				if err != nil {
					handler.RespondError(w, http.StatusOK, handler.UnauthorizedError, err)
					return
				}
				token, err := client.VerifyIDToken(netContext.Background(), bearerToken[1])
				if err != nil {
					log.Println(bearerToken[1])
					handler.RespondError(w, http.StatusOK, handler.UnauthorizedError, err)
					return
				}
				user, err := dao.UserFindByKey(a.MongoDB, &bson.M{"user_id": token.UID})
				if err != nil {
					handler.RespondError(w, http.StatusOK, handler.NotFound, err)
					return
				}
				context.Set(req, "user", user)
				next(w, req)
			} else {
				handler.RespondError(w, http.StatusOK, handler.UnauthorizedError, fmt.Errorf("Invalid token"))
				return
			}
		} else {
			handler.RespondError(w, http.StatusOK, handler.UnauthorizedError, fmt.Errorf("No token provided"))
		}
	})
}

// Get wraps the router for GET method for socket
func (a *App) GetSocket(path string, f *socketio.Server) {
	a.Router.Handle(path, f)
}

// Get wraps the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Post wraps the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Put wraps the router for PUT method
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Delete wraps the router for DELETE method
func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

func (a *App) Run(host string) {
	log.Println("Application started", time.Now().String())
	//log.Fatal(http.ListenAndServe(host, &CORSRouterDecorator{a.Router, a.Socket}))
	log.Fatal(http.ListenAndServeTLS(host, "/etc/ssl/certs/apache-selfsigned.crt", "/etc/ssl/certs/apache-selfsigned.key", &CORSRouterDecorator{a.Router, a.Socket}))
}

// setRouters sets the all required routers
func (a *App) setRouters() {

	a.GetSocket("/socket.io/", a.Socket)
	a.Post("/api/user", a.CreateUser)
	a.Get("/api/user", a.ValidateMiddleware(a.GetUser))

	a.Post("/api/plan/test", a.GetPlan)
	a.Get("/api/plan", a.GetPlans)
	a.Post("/api/plan", a.CreatePlan)
	a.Delete("/api/plan", a.DeletePlan)


}

func (a *App) GetUser(w http.ResponseWriter, r *http.Request) {
	user := new(dao.User)
	userContext := context.Get(r, "user")
	user, err := dao.GetUserStruct(userContext)
	if err != nil {
		handler.RespondError(w, http.StatusOK, handler.InternalError, err)
	}
	handler.GetUserEndpoint(w, r, user)
}

func (a *App) CreateUser(w http.ResponseWriter, r *http.Request) {
	dbSession := a.MongoDB.Copy()
	defer dbSession.Close()
	handler.CreateUserEndPoint(dbSession, a.Stripe, w, r)
}

func (a *App) GetPlans(w http.ResponseWriter, r *http.Request) {
	dbSession := a.MongoDB.Copy()
	defer dbSession.Close()
	handler.GetPlansEndpoint(dbSession, w, r)
}

func (a *App) GetPlan(w http.ResponseWriter, r *http.Request) {
	dbSession := a.MongoDB.Copy()
	defer dbSession.Close()
	handler.GetPlanEndpoint(dbSession, w, r)
}

func (a *App) CreatePlan(w http.ResponseWriter, r *http.Request) {
	dbSession := a.MongoDB.Copy()
	defer dbSession.Close()
	handler.CreatePlanEndpoint(dbSession, w, r)
}

func (a *App) DeletePlan(w http.ResponseWriter, r *http.Request) {
	dbSession := a.MongoDB.Copy()
	defer dbSession.Close()
	handler.DeletePlanEndpoint(dbSession, w, r)
}