package bitcoinactions

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/hkparker/crypto-proxy/parsers/bitcoin"

	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

func init() {
	add(BitcoinAction{
		Scope: []string{"filterload"},
		Do:    saveBloomFilters,
	})
}

func saveBloomFilters(message bitcoin.BitcoinMessage) {
	filterload, err := bitcoin.ParseFilterLoad(message)
	if err != nil {
		log.WithFields(log.Fields{
			"at":    "bitcoinactions.saveBloomFilters",
			"error": err.Error(),
		}).Error("error parsing filterload message")
		return
	}
	dir := "bloom_filters/"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0700)
		if err != nil {
			log.WithFields(log.Fields{
				"at":    "bitcoinactions.saveBloomFilters",
				"error": err.Error(),
				"path":  dir,
			}).Error("error creating directory for bloom filters")
		}
	}
	path := dir + uuid.NewV4().String()
	data, err := json.MarshalIndent(filterload, "", "  ")
	if err != nil {
		log.WithFields(log.Fields{
			"at":    "bitcoinactions.saveBloomFilters",
			"error": err.Error(),
		}).Error("error marshalling filterload message")
	}
	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		log.WithFields(log.Fields{
			"at":    "bitcoinactions.saveBloomFilters",
			"error": err.Error(),
			"path":  path,
		}).Error("error saving bloom filter file")
	}
}
