package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
	"github.com/rs/cors"

	"github.com/QuoineFinancial/liquid-chain/api/chain"
	"github.com/QuoineFinancial/liquid-chain/api/resource"
	"github.com/QuoineFinancial/liquid-chain/db"
	"github.com/QuoineFinancial/liquid-chain/storage"
)

// API contains all info to serve an api server
type API struct {
	url        string
	rpcServer  *rpc.Server
	httpServer *http.Server
	Router     *mux.Router

	tmAPI resource.TendermintAPI
	meta  *storage.MetaStorage
	state *storage.StateStorage
	chain *storage.ChainStorage
}

// NewAPI return an new instance of API
func NewAPI(url, tmURL, rootDir string, metaDB, stateDB, chainDB db.Database) *API {
	api := &API{
		url:   url,
		tmAPI: resource.NewTendermintAPI(rootDir, tmURL),
		meta:  storage.NewMetaStorage(metaDB),
		state: storage.NewStateStorage(stateDB),
		chain: storage.NewChainStorage(chainDB),
	}
	api.setupServer()
	api.registerServices()
	api.setupRouter()
	return api
}

func (api *API) setupServer() {
	server := rpc.NewServer()
	server.RegisterCodec(json2.NewCodec(), "application/json")
	api.rpcServer = server
}

func (api *API) setupRouter() {
	if api.rpcServer == nil {
		panic("api.setupRouter call without api.server")
	}
	api.Router = mux.NewRouter()
	api.Router.Handle("/", api.rpcServer).Methods("POST")
	api.httpServer = &http.Server{
		Handler: cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
			AllowedMethods:   []string{"POST", "DELETE", "PUT", "GET", "HEAD", "OPTIONS"},
		}).Handler(api.Router),
		Addr: api.url,
	}
}

func (api *API) registerServices() {
	if api.rpcServer == nil {
		panic("api.registerServices call without api.server")
	}
	if err := api.rpcServer.RegisterService(chain.NewService(api.tmAPI, api.meta, api.state, api.chain), "chain"); err != nil {
		panic(err)
	}
}

// Serve starts the server to serve request
func (api *API) Serve() error {
	log.Println("Server is ready at", api.url)
	err := api.httpServer.ListenAndServe()
	return err
}

// Close will immediately stop the server without waiting for any active connection to complete
// For gracefully shutdown please implement another function and use Server.Shutdown()
func (api *API) Close() {
	log.Println("Closing server")
	if api.httpServer != nil {
		err := api.httpServer.Close()
		if err != nil {
			panic(err)
		}
	}
}
