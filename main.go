package main

import (
	"flag"

	log "github.com/sirupsen/logrus"
)

func main() {
	var crypto string
	var bind string
	var bindPort int
	var dial string
	var dialPort int

	flag.StringVar(&crypto, "crypto", "bitcoin", "cryptocurrency to proxy")
	flag.StringVar(&bind, "bind", "0.0.0.0", "address to bind")
	flag.IntVar(&bindPort, "bind-port", 8333, "port to bind")
	flag.StringVar(&dial, "dial", "", "bitcoin node to proxy to")
	flag.IntVar(&dialPort, "dial-port", 8333, "port to dial for bitcoin node")
	flag.Parse()

	listener := openBind(bind, bindPort)

	for {
		client, err := listener.Accept()
		if err != nil {
			log.WithFields(log.Fields{
				"at":    "main",
				"error": err.Error(),
			}).Error("error accepting connection from client")
			continue
		} else {
			log.WithFields(log.Fields{
				"at":     "main",
				"client": client.RemoteAddr().String(),
			}).Info("new connection from a client")
		}

		node, err := openProxy(dial, dialPort)
		if err != nil {
			log.WithFields(log.Fields{
				"at":    "main",
				"error": err.Error(),
				"node":  dial,
				"port":  dialPort,
			}).Error("error connecting back to actual node")
			client.Close()
			continue
		} else {
			go intercept(crypto, client, node)
		}
	}
}
