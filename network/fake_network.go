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
)

// FakeNetwork is an empty fake struct :)
type FakeNetwork struct {
	// nothing here :)
}

// IsReady checks if the fake network is ready and always returns false
func (fn FakeNetwork) IsReady() bool {
	return false
}

// GetStatus gets the fake network status
func (fn FakeNetwork) GetStatus() string {
	return "not configured"
}

// ListPeers gets a fake listing of peers and always returns an error
func (fn FakeNetwork) ListPeers() ([]Peer, error) {
	return nil, fmt.Errorf("cannot list peers, because there is no network configured")
}

// UpdatePeers fakes updating of peers and always returns an error
func (fn FakeNetwork) UpdatePeers(newPeers []Peer) error {
	return fmt.Errorf("cannot update peers, because there is no network configured")
}
