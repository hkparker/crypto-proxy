package bitcoin

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"io"

	log "github.com/sirupsen/logrus"
)

const BTC_MAINNET = 0xD9B4BEF9
const BTC_TESTNET = 0xDAB5BFFA
const BTC_TESTNET3 = 0x0709110B
const BTC_NAMECOIN = 0xFEB4BEF9

type BitcoinMessage struct {
	Network  uint32
	Command  string
	Length   uint32
	Checksum uint32
	Payload  []byte
}

func ReadBitcoinMessage(reader io.Reader) (BitcoinMessage, error) {
	message := BitcoinMessage{}

	network, err := scanNextMagicBytes(reader)
	if err != nil {
		return message, err
	}
	message.Network = network

	command, err := readCommand(reader)
	if err != nil {
		return message, err
	}
	message.Command = command

	length, err := readUint32(reader)
	if err != nil {
		return message, err
	}
	message.Length = length

	checksum, err := readUint32(reader)
	if err != nil {
		return message, err
	}
	message.Checksum = checksum

	payload, err := readBytes(reader, length)
	if err != nil {
		return message, err
	}
	message.Payload = payload

	err = checkChecksum(message.Payload, checksum)
	if err != nil {
		log.WithFields(log.Fields{
			"at":      "bitcoin.ReadBitcoinMessage",
			"command": command,
			"error":   err.Error(),
		}).Error("checksum mismatch on message")
		return message, err
	}

	return message, nil
}

func scanNextMagicBytes(reader io.Reader) (uint32, error) {
	buffer := make([]byte, 4)
	for {
		next := make([]byte, 1)
		_, err := reader.Read(next)
		if err != nil && err != io.EOF {
			return 0, err
		}
		buffer[0] = buffer[1]
		buffer[1] = buffer[2]
		buffer[2] = buffer[3]
		buffer[3] = next[0]
		match, network, err := matchMagicBytes(buffer)
		if err != nil {
			return 0, err
		}
		if match {
			return network, nil
		}
	}
}

func matchMagicBytes(buffer []byte) (bool, uint32, error) {
	var network uint32
	err := binary.Read(bytes.NewReader(buffer), binary.LittleEndian, &network)
	if err != nil {
		return false, 0, err
	}
	if network == BTC_MAINNET || network == BTC_TESTNET || network == BTC_TESTNET3 || network == BTC_NAMECOIN {
		return true, network, nil
	}
	return false, network, nil
}

func readCommand(reader io.Reader) (string, error) {
	strBytes := make([]byte, 12)
	_, err := reader.Read(strBytes)
	if err != nil {
		return "", err
	}
	strBytes = bytes.TrimRight(strBytes, string([]byte{0x00}))
	return string(strBytes), nil
}

func readUint32(reader io.Reader) (uint32, error) {
	intBytes := make([]byte, 4)
	_, err := reader.Read(intBytes)
	if err != nil {
		return 0, err
	}
	var data uint32
	err = binary.Read(bytes.NewReader(intBytes), binary.LittleEndian, &data)
	if err != nil {
		return 0, err
	}
	return data, nil
}

func readBytes(reader io.Reader, count uint32) ([]byte, error) {
	data := make([]byte, count)
	_, err := io.ReadFull(reader, data)
	return data, err
}

func checkChecksum(payload []byte, checksum uint32) error {
	innerSha := sha256.Sum256(payload)
	outerSha := sha256.Sum256(innerSha[:])
	firstFour := outerSha[0:4]
	var actualChecksum uint32
	err := binary.Read(bytes.NewReader(firstFour), binary.LittleEndian, &actualChecksum)
	if err != nil {
		return errors.New("unable to read sha bytes as little endian uint32")
	}
	if checksum != actualChecksum {
		return errors.New("checksum mismatch")
	}
	return nil
}
