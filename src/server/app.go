package server

import (
	"anet"
	"config"
	"db"
	"log"
	"protocol"
)

const (
	MAXN_INTERNAL_EVENTS = 1024
	MAXN_EXTERNAL_EVENTS = 1024
)

type App struct {
	cfg            config.Config
	closed         bool
	server         *anet.Server
	events         chan anet.Event
	users          map[int32]*UserSession
	sessionMapping map[int32]int32
}

func NewApp(cfg config.Config) (*App, error) {
	app := new(App)
	app.cfg = cfg

	app.users = make(map[int32]*UserSession)
	app.sessionMapping = make(map[int32]int32)
	app.events = make(chan anet.Event, MAXN_EXTERNAL_EVENTS)

	netCfg := cfg.GetSection("network")
	listenSection := netCfg.GetSection("listen")
	addr := listenSection.GetStr("4client")
	log.Printf("init network %s for %s ok.", "4client", addr)

	app.server = anet.NewServer("tcp4", addr, protocol.Proto{}, app.events)

	app.init(cfg)
	return app, nil
}

func (app *App) init(cfg config.Config) {
	dbCfg := cfg.GetSection("database")
	host := dbCfg.GetStr("host")
	port := dbCfg.GetInt16("port")
	user := dbCfg.GetStr("username")
	pass := dbCfg.GetStr("password")
	dbname := dbCfg.GetStr("dbname")
	log.Printf("connect mysql: host=%s, port=%d, user=%s, pass=%s, dbname=%s",
		host, port, user, pass, dbname)
	db.Open(host, port, user, pass, dbname)
}

func (app *App) Close() {
	if app.closed {
		return
	}
	app.server.Close()
	db.Close()
	close(app.events)
}

func (app *App) Run() {
	if err := app.server.ListenAndServe(); err != nil {
		log.Printf("ListenAndServe() error: %s", err)
	} else {
		app.mainLoop()
	}
}

func (app *App) mainLoop() {
	for {
		select {
		case ev, ok := <-app.events:
			if ok {
				app.onEvent(ev)
			}
		}
	}
}
