package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/apex/log"
)

func convertUser(user string) string {
	return strings.SplitN(user, "@", 2)[0]
}

func main() {
	log.Info("Starting oauth2_proxy-token service")
	terminateChan := make(chan os.Signal, 2)
	signal.Notify(terminateChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-terminateChan
		log.Warn("Detected shutdown request")
		os.Exit(0)
	}()

	configFile := flag.String("config", "tokenconfig.json", "Configuration file in JSON format")
	flag.Parse()
	config := NewConfig()
	err := config.Read(*configFile)
	if err != nil {
		log.WithFields(log.Fields{
			"configFile": *configFile,
			"err":        err,
		}).Fatal("Failed to read configuration")
	} else {
		var mutex = &sync.Mutex{}
		go func() {
			for true {
				mutex.Lock()
				maintainTokens(config)
				mutex.Unlock()
				time.Sleep(time.Duration(config.MaintenanceIntervalSeconds) * time.Second)
			}
		}()
		http.HandleFunc(config.HTTPPath, func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.Header.Get("Authorization"), "Basic") {
				log.Warn("Rejected token attempt via HTTP Basic Authentication request")
				return
			}
			user := convertUser(r.Header.Get(config.HeaderUsername))
			uri := r.Header.Get(config.HeaderURI)
			if uri == "" {
				uri = r.RequestURI
			}
			mutex.Lock()
			t, err := createOrUpdateToken(config, user, uri)
			mutex.Unlock()
			if err != nil {
				t = fmt.Sprintf("ERROR: Failed to generate token -> %s", err)
			}
			w.Write([]byte(t))
		})
		log.WithFields(log.Fields{
			"hostAndPort": config.HTTPHostPort,
			"httpPath":    config.HTTPPath,
		}).Info("Listening for incoming requests")
		err := http.ListenAndServe(config.HTTPHostPort, nil)
		if err != nil {
			log.WithField("err", err).Fatal("Unexpected failure")
		}
	}
	log.Info("Exiting oauth2_proxy-token service")
}
