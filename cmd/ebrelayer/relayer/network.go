package relayer

// ------------------------------------------------------------
//    Network
//
//    Validates input and initializes a websocket Ethereum
//    client.
// ------------------------------------------------------------

import (
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/golang/glog"
)

// IsWebsocketURL return true if the given URL is a websocket URL
func IsWebsocketURL(rawurl string) bool {
	u, err := url.Parse(rawurl)
	if err != nil {
		log.Infof("Error while parsing URL: %v", err)
		return false
	}
	if u.Scheme == "ws" || u.Scheme == "wss" {
		return true
	}
	return false
}

// SetupWebsocketEthClient returns an websocket ethclient if URL is valid.
func SetupWebsocketEthClient(ethURL string) (*ethclient.Client, error) {
	if ethURL == "" {
		return nil, nil
	}

	if !IsWebsocketURL(ethURL) {
		return nil, fmt.Errorf(
			"In valid websocket eth client URL: %v",
			ethURL,
		)
	}

	client, err := ethclient.Dial(ethURL)
	if err != nil {
		return nil, fmt.Errorf("error dialing websocket client")
	}

	return client, nil
}
