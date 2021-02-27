package main

import (
	"net/http"
	"time"

	"github.com/imega/daemon"
	"github.com/imega/daemon/configuring/env"
	httpserver "github.com/imega/daemon/http-server"
	"github.com/imega/daemon/logging"
)

const shutdownTimeout = 15 * time.Second

func main() {
	log := logging.New(logging.Config{
		Channel: "stock-miner",
		Level:   "debug",
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	h := httpserver.New("stock-miner", log, mux)

	cr := env.Once(h.WatcherConfigFunc)

	d, err := daemon.New(log, cr)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("stock-miner is started")

	if err := d.Run(shutdownTimeout); err != nil {
		log.Errorf("failed to loop until shutdown: %s", err)
	}

	log.Info("stock-miner is stopped")
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
