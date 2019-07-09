package bitcoinactions

import (
	"github.com/hkparker/crypto-proxy/helpers"
	"github.com/hkparker/crypto-proxy/parsers/bitcoin"
)

var actions []BitcoinAction

type BitcoinAction struct {
	Scope []string
	Do    func(bitcoin.BitcoinMessage)
}

func add(action BitcoinAction) {
	actions = append(actions, action)
}

func Dispatch(message bitcoin.BitcoinMessage) {
	for _, action := range actions {
		if helpers.StringsInclude(action.Scope, message.Command) {
			action.Do(message)
		}
	}
}
