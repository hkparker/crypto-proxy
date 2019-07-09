package actions

import (
	"io"

	"github.com/hkparker/crypto-proxy/actions/bitcoin"
	"github.com/hkparker/crypto-proxy/parsers/bitcoin"

	log "github.com/sirupsen/logrus"
)

func GetActionWriter(crypto string) io.Writer {
	reader, writer := io.Pipe()

	switch crypto {
	case "bitcoin":
		go dispatchBitcoinMessages(reader)
	default:
		log.WithFields(log.Fields{
			"at":     "actions.GetActionWriter",
			"crypto": crypto,
		}).Fatal("unsupported crypto")
	}

	return writer
}

func dispatchBitcoinMessages(reader io.Reader) {
	for {
		btcMessage, err := bitcoin.ReadBitcoinMessage(reader)
		if err != nil {
			log.WithFields(log.Fields{
				"at":    "actions.dispatchBitcoinMessages",
				"error": err.Error(),
			}).Error("error parsing bitcoin message")
		} else {
			bitcoinactions.Dispatch(btcMessage)
		}
	}
}
