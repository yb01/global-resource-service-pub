/*
Copyright 2014 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package versioned

import (
	"encoding/json"
	"fmt"

	"global-resource-service/resource-management/pkg/common-lib/types"
	"global-resource-service/resource-management/pkg/common-lib/types/event"
)

// Decoder implements the watch.Decoder interface for io.ReadClosers that
// have contents which consist of a series of watchEvent objects encoded
// with the given streaming decoder. The internal objects will be then
// decoded by the embedded decoder.
type Decoder struct {
	decoder *json.Decoder
}

// NewDecoder creates an Decoder for the given writer and codec.
func NewDecoder(decoder *json.Decoder) *Decoder {
	return &Decoder{
		decoder: decoder,
	}
}

// Decode blocks until it can return the next object in the reader. Returns an error
// if the reader is closed or an object can't be decoded.
func (d *Decoder) Decode() (event.EventType, *types.LogicalNode, error) {
	var got event.NodeEvent
	err := d.decoder.Decode(&got)
	if err != nil {
		return "", nil, err
	}

	switch got.Type {
	case event.Added, event.Modified, event.Deleted, event.Error, event.Bookmark:
	default:
		return "", nil, fmt.Errorf("got invalid watch event type: %v", got.Type)
	}

	return got.Type, got.Node, nil
}

// Close closes the underlying r.
func (d *Decoder) Close() {
	return
}
