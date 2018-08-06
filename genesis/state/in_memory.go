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
	"fmt"
	"sync"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type inMemoryState struct {
	sync.Mutex
	peers map[strfmt.UUID]Peer
}

// NewInMemoryState creates a new in memory state
func NewInMemoryState() State {
	state := inMemoryState{
		peers: make(map[strfmt.UUID]Peer),
	}
	go state.garbageCollect()
	return State(&state)
}

func (im *inMemoryState) RegisterPeer(name string, uri strfmt.URI) (*Peer, error) {
	im.Lock()
	defer im.Unlock()

	id := strfmt.UUID(uuid.NewV4().String())

	log.Debugf("Registering peer '%v' with id '%v'", name, id)
	peer := Peer{
		PeerInfo: PeerInfo{
			Id:            id,
			LastContactAt: time.Now(),
		},
		name: name,
		uri:  uri,
	}

	im.peers[id] = peer
	go im.broadcastUpdate()
	return &peer, nil
}

func (im *inMemoryState) ListPeers() ([]Peer, error) {
	im.Lock()
	defer im.Unlock()

	peers := make([]Peer, 0)

	for _, v := range im.peers {
		peers = append(peers, v)
	}

	return peers, nil
}

func (im *inMemoryState) RemovePeer(id strfmt.UUID) error {
	im.Lock()
	defer im.Unlock()

	_, ok := im.peers[id]

	if ok {
		delete(im.peers, id)
	}

	go im.broadcastUpdate()

	return nil
}

func (im *inMemoryState) UpdateLastContact(id strfmt.UUID, contactAt time.Time) error {
	log.Debugf("Updating last contact for %v", id)

	im.Lock()
	defer im.Unlock()

	peer, ok := im.peers[id]

	if ok {
		peer.LastContactAt = contactAt
		im.peers[id] = peer
		return nil
	}
	return fmt.Errorf("no such peer exists")
}

func (im *inMemoryState) garbageCollect() {
	for {
		time.Sleep(1 * time.Second)
		deletedSome := false

		im.Lock()
		for key, peer := range im.peers {
			peerTimesOutAt := peer.PeerInfo.LastContactAt.Add(time.Second * 60)
			if time.Now().After(peerTimesOutAt) {
				log.Infof("garbage collecting peer %v", peer.ID)
				delete(im.peers, key)
				deletedSome = true
			}
		}
		im.Unlock()

		if deletedSome {
			im.broadcastUpdate()
		}
	}
}

func (im *inMemoryState) broadcastUpdate() {
	log.Info("Broadcasting peer update")
	im.Lock()
	defer im.Unlock()

	peers := make([]Peer, 0)

	for _, peer := range im.peers {
		peers = append(peers, peer)
	}

	for _, peer := range peers {
		go broadcastUpdate(peer, peers)
	}
}
