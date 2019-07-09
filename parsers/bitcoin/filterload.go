package bitcoin

import (
	"bytes"
	"encoding/binary"
	"errors"

	log "github.com/sirupsen/logrus"
)

const BLOOM_UPDATE_NONE = 0x00
const BLOOM_UPDATE_ALL = 0x01
const BLOOM_UPDATE_P2PUBKEY_ONLY = 0x02

type BitcoinFilterLoadMessage struct {
	Message       BitcoinMessage
	Filter        []byte
	HashFunctions uint32
	Tweak         uint32
	Flags         byte
}

func ParseFilterLoad(message BitcoinMessage) (BitcoinFilterLoadMessage, error) {
	filterload := BitcoinFilterLoadMessage{}

	msgLen := len(message.Payload)
	if msgLen < 10 {
		return filterload, errors.New("too few bytes to be a valid filterload payload")
	}

	lastByte := message.Payload[msgLen-1]
	filterload.Flags = lastByte
	if lastByte != BLOOM_UPDATE_NONE && lastByte != BLOOM_UPDATE_ALL && lastByte != BLOOM_UPDATE_P2PUBKEY_ONLY {
		log.WithFields(log.Fields{
			"at":    "bitcoin.ParseFilterLoad",
			"flags": int(lastByte),
		}).Warn("unknown filterload flags")
	}

	tweakBytes := message.Payload[msgLen-5 : msgLen-1]
	var tweak uint32
	err := binary.Read(bytes.NewReader(tweakBytes), binary.LittleEndian, &tweak)
	if err != nil {
		log.WithFields(log.Fields{
			"at":    "bitcoin.ParseFilterLoad",
			"error": err.Error(),
		}).Error("error reading tweaks bytes")
		return filterload, err
	}
	filterload.Tweak = tweak

	hashFunctionsBytes := message.Payload[msgLen-9 : msgLen-5]
	var hashFunctions uint32
	err = binary.Read(bytes.NewReader(hashFunctionsBytes), binary.LittleEndian, &hashFunctions)
	if err != nil {
		log.WithFields(log.Fields{
			"at":    "bitcoin.ParseFilterLoad",
			"error": err.Error(),
		}).Error("error reading hash functions bytes")
		return filterload, err
	}
	filterload.HashFunctions = hashFunctions

	filterload.Filter = message.Payload[0 : msgLen-9]

	message.Payload = []byte{}
	filterload.Message = message

	return filterload, nil
}
