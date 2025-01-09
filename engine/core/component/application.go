package core

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

type Application struct {
	Router     *Router
	Components map[string]ComponentInterface
	HotReload  *HotReload
	Config     *ApplicationConfig
	server     *http.Server
}

type ApplicationConfig struct {
	Port            string
	TemplatesDir    string
	StaticDir       string
	EnableHotReload bool
}

func NewApplication() *Application {
	return &Application{
		Router:     NewRouter(),
		Components: make(map[string]ComponentInterface),
		Config: &ApplicationConfig{
			Port:            "8080",
			TemplatesDir:    "./templates",
			StaticDir:       "./static",
			EnableHotReload: true,
		},
	}
}

func (app *Application) Configure(config ApplicationConfig) {
	app.Config = &config
	if config.EnableHotReload {
		hotReload, err := NewHotReload()
		if err != nil {
			log.Printf("Warning: Hot reload initialization failed: %v", err)
		} else {
			app.HotReload = hotReload
		}
	}
}

func (app *Application) RegisterComponent(path string, component ComponentInterface) {
	app.Components[path] = component
	app.Router.AddRoute(path, component)
	if app.HotReload != nil {
		componentPath := fmt.Sprintf("./components/%s", path)
		app.HotReload.Watch(componentPath)
	}
}

func (app *Application) handleRequest(w http.ResponseWriter, r *http.Request) {
	component, exists := app.Router.routes[r.URL.Path]
	if !exists {
		http.NotFound(w, r)
		return
	}
	if component.GetLifeCycle().BeforeMount != nil {
		component.GetLifeCycle().BeforeMount()
	}
	w.Header().Set("Content-Type", "text/html")
	err := component.GetTemplate().Execute(w, component.GetState())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if component.GetLifeCycle().AfterMount != nil {
		component.GetLifeCycle().AfterMount()
	}
}

func (app *Application) Run(port ...string) error {
	if len(port) > 0 {
		app.Config.Port = port[0]
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.handleRequest)
	fs := http.FileServer(http.Dir(app.Config.StaticDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	if app.HotReload != nil {
		mux.HandleFunc("/ws", app.handleHotReloadWs)
	}

	app.server = &http.Server{
		Addr:    ":" + app.Config.Port,
		Handler: mux,
	}
	log.Printf("Server starting on http://localhost:%s", app.Config.Port)
	return app.server.ListenAndServe()
}

func (app *Application) handleHotReloadWs(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()
	for {
		select {
		case event := <-app.HotReload.watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				message := map[string]string{
					"type": "reload",
					"file": event.Name,
				}
				if err := conn.WriteJSON(message); err != nil {
					log.Printf("WebSocket write failed: %v", err)
					return
				}
			}
		case err := <-app.HotReload.watcher.Errors:
			log.Printf("Watcher error: %v", err)
		}
	}
}
