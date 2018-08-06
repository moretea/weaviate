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
 */package state

import (
	weaviate_client "github.com/creativesoftwarefdn/weaviate/client"
	weaviate_p2p "github.com/creativesoftwarefdn/weaviate/client/p2_p"
	weaviate_models "github.com/creativesoftwarefdn/weaviate/models"

	"net/url"

	log "github.com/sirupsen/logrus"
)

func broadcastUpdate(peer Peer, peers []Peer) {
	log.Debugf("Broadcasting peer update to %v", peer.ID)
	peerURI, err := url.Parse(string(peer.URI()))

	if err != nil {
		log.Infof("Could not broadcast to peer %v; Peer URI is invalid (%v)", peer.ID, peer.URI())
		return
	}

	transportConfig := weaviate_client.TransportConfig{
		Host:     peerURI.Host,
		BasePath: peerURI.Path,
		Schemes:  []string{peerURI.Scheme},
	}

	peerUpdates := make(weaviate_models.PeerUpdateList, 0)

	for _, peer := range peers {
		peerUpdate := weaviate_models.PeerUpdate{
			URI:  peer.URI(),
			ID:   peer.ID,
			Name: peer.Name(),
		}

		peerUpdates = append(peerUpdates, &peerUpdate)
	}

	client := weaviate_client.NewHTTPClientWithConfig(nil, &transportConfig)
	params := weaviate_p2p.NewWeaviateP2pGenesisUpdateParams()
	params.Peers = peerUpdates
	_, err = client.P2P.WeaviateP2pGenesisUpdate(params)
	if err != nil {
		log.Debugf("failed to update %v, because %v", peer.ID, err)
	}
}
