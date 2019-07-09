package main

import (
	"io"
	"net"
	"strconv"

	"github.com/hkparker/crypto-proxy/actions"

	log "github.com/sirupsen/logrus"
)

func intercept(crypto string, client, node net.Conn) {
	multiWriter := io.MultiWriter(actions.GetActionWriter(crypto), node)

	go func() {
		if bytes_copied, err := io.Copy(multiWriter, client); err != nil {
			log.WithFields(log.Fields{
				"at":           "intercept",
				"error":        err.Error(),
				"bytes_copied": bytes_copied,
				"client":       client.RemoteAddr().String(),
			}).Error("error copying data from client connection into MultiWriter")
		}
	}()

	go func() {
		if bytes_copied, err := io.Copy(client, node); err != nil {
			log.WithFields(log.Fields{
				"at":             "intercept",
				"error":          err.Error(),
				"bytes_copied":   bytes_copied,
				"node_address":   node.RemoteAddr().String(),
				"client_address": client.RemoteAddr().String(),
			}).Error("error copying data from node back to client")
		}
	}()
}

func openProxy(dial string, port int) (net.Conn, error) {
	node, err := net.Dial("tcp", formatConnectionString(dial, port))
	if err != nil {
		log.WithFields(log.Fields{
			"at":    "openProxy",
			"error": err.Error(),
		}).Error("error dialing node")
	} else {
		log.WithFields(log.Fields{
			"at":      "openProxy",
			"address": dial,
			"port":    port,
		}).Info("connected to node")
	}
	return node, err
}

func openBind(bind string, port int) net.Listener {
	log.WithFields(log.Fields{
		"at":      "openBind",
		"address": bind,
		"port":    port,
	}).Info("creating listener")
	listener, err := net.Listen("tcp", formatConnectionString(bind, port))
	if err != nil {
		log.WithFields(log.Fields{
			"at":    "openBind",
			"error": err.Error(),
		}).Fatal("error binding proxy address")
	}
	return listener
}

func formatConnectionString(addr string, port int) string {
	return addr + ":" + strconv.Itoa(port)
}
