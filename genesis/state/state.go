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
	"time"

	"github.com/go-openapi/strfmt"
)

// PeerInfo contains the information including the UUID of the peer
type PeerInfo struct {
	ID            strfmt.UUID
	LastContactAt time.Time
}

// Peer contains individual peer info including the URI
type Peer struct {
	PeerInfo
	name string
	uri  strfmt.URI
}

// Name returns the name of the peer
func (p Peer) Name() string {
	return p.name
}

// URI returns the uri of the peer
func (p Peer) URI() strfmt.URI {
	return p.uri
}

// State is an abstract interface over how the Genesis server should store state.
type State interface {
	RegisterPeer(name string, uri strfmt.URI) (*Peer, error)
	ListPeers() ([]Peer, error)

	// Idempotent remove; removing a non-existing peer should not fail.
	RemovePeer(id strfmt.UUID) error

	UpdateLastContact(id strfmt.UUID, contactTime time.Time) error
}
