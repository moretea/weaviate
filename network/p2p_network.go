/*                          _       _
 *__      _____  __ ___   ___  __ _| |_ ___
 *\ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
 * \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
 *  \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
 *
 * Copyright © 2016 - 2018 Weaviate. All rights reserved.
 * LICENSE: https://github.com/creativesoftwarefdn/weaviate/blob/develop/LICENSE.md
 * AUTHOR: Bob van Luijt (bob@kub.design)
 * See www.creativesoftwarefdn.org for details
 * Contact: @CreativeSofwFdn / bob@kub.design
 */package network

import (
	"fmt"
	"sync"
	"time"

	"github.com/creativesoftwarefdn/weaviate/messages"
	"github.com/go-openapi/strfmt"

	"net/url"

	genesis_client "github.com/creativesoftwarefdn/weaviate/genesis/client"
	client_ops "github.com/creativesoftwarefdn/weaviate/genesis/client/operations"
	genesisModels "github.com/creativesoftwarefdn/weaviate/genesis/models"
)

const (
	networkStateBootstrapping = "network bootstrapping"
	networkStateFailed        = "network failed"
	networkStateHealthy       = "network healthy"
)

// The real network implementation. Se also `fake_network.go`
type network struct {
	sync.Mutex

	// Peer ID assigned by genesis server
	peerID    strfmt.UUID
	peerName  string
	publicURL strfmt.URI

	state      string
	genesisURL strfmt.URI
	messaging  *messages.Messaging
	client     genesis_client.WeaviateGenesisServer
	peers      []Peer
}

// BootstrapNetwork bootstraps the P2P network or returns error
func BootstrapNetwork(m *messages.Messaging, genesisURL strfmt.URI, publicURL strfmt.URI, peerName string) (*Network, error) {
	if genesisURL == "" {
		return nil, fmt.Errorf("no genesis URL provided in network configuration")
	}

	genesisURI, err := url.Parse(string(genesisURL))
	if err != nil {
		return nil, fmt.Errorf("could not parse genesis URL '%v'", genesisURL)
	}

	if publicURL == "" {
		return nil, fmt.Errorf("no public URL provided in network configuration")
	}

	_, err = url.Parse(string(publicURL))
	if err != nil {
		return nil, fmt.Errorf("could not parse public URL '%v'", publicURL)
	}

	if peerName == "" {
		return nil, fmt.Errorf("no peer name specified in network configuration")
	}

	transportConfig := genesis_client.TransportConfig{
		Host:     genesisURI.Host,
		BasePath: genesisURI.Path,
		Schemes:  []string{genesisURI.Scheme},
	}

	client := genesis_client.NewHTTPClientWithConfig(nil, &transportConfig)

	n := network{
		publicURL:  publicURL,
		peerName:   peerName,
		state:      networkStateBootstrapping,
		genesisURL: genesisURL,
		messaging:  m,
		client:     *client,
		peers:      make([]Peer, 0),
	}

	// Bootstrap the network in the background.
	go n.bootstrap()

	nw := Network(&n)
	return &nw, nil
}

func (n *network) bootstrap() {
	time.Sleep(10) //TODO: Use channel close to listen for when complete configuration is done.
	n.messaging.InfoMessage("Bootstrapping network")

	newPeer := genesisModels.PeerUpdate{
		PeerName: n.peerName,
		PeerURI:  n.publicURL,
	}

	params := client_ops.NewGenesisPeersRegisterParams()
	params.Body = &newPeer
	response, err := n.client.Operations.GenesisPeersRegister(params)
	if err != nil {
		n.messaging.ErrorMessage(fmt.Sprintf("could not register this peer in the network, because: %+v", err))
		n.state = networkStateFailed
	} else {
		n.state = networkStateHealthy
		n.peerID = response.Payload.Peer.ID
		n.messaging.InfoMessage(fmt.Sprintf("registered at Genesis server with id '%v'", n.peerID))
	}

	go n.keepPinging()
}

func (n *network) IsReady() bool {
	return false
}

func (n *network) GetStatus() string {
	return n.state
}

func (n *network) ListPeers() ([]Peer, error) {
	return nil, fmt.Errorf("cannot list peers, because there is no network configured")
}

func (n *network) UpdatePeers(newPeers []Peer) error {
	n.Lock()
	defer n.Unlock()

	n.messaging.InfoMessage(fmt.Sprintf("received updated peer list with %v peers", len(newPeers)))

	n.peers = newPeers

	return nil
}

func (n *network) keepPinging() {
	for {
		time.Sleep(30 * time.Second)
		n.messaging.InfoMessage("pinging Genesis server")

		n.Lock()
		params := client_ops.NewGenesisPeersPingParams()
		params.PeerID = n.peerID
		n.Unlock()
		_, err := n.client.Operations.GenesisPeersPing(params)
		if err != nil {
			n.messaging.InfoMessage(fmt.Sprintf("Could not ping Genesis server; %+v", err))
		}
	}
}
