package types

import (
	"encoding/json"
	"io"

	"global-resource-service/resource-management/pkg/common-lib/types"
	"global-resource-service/resource-management/pkg/common-lib/types/event"
)

type RegionNodeEvents [][]*event.NodeEvent

type PostCRVstatus bool

// RRM: Resource Region Manager
//
type ResponseFromRRM struct {
	RegionNodeEvents [][]*event.NodeEvent
	RvMap            types.ResourceVersionMap
	Length           uint64
}

// The type is for pulling data with batch from RRM - Resource Region Manager
//
type PullDataFromRRM struct {
	BatchLength uint64
	DefaultCRV  uint64
	CRV         types.ResourceVersionMap
}

// ToJSON serializes the contents of the collection to JSON
// NewEncoder provides better performance than json.Unmarshal
// as it does not have to buffer the output into an in memory
// slice of bytes. This reduces allocations and the overheads
// of the service
//
// https://golang.org/pkg/encoding/json/#NewEncoder
//
// ToJSON serializes the given interface into a string based JSON format
//
func ToJSON(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)

	return e.Encode(i)
}

func (p *RegionNodeEvents) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)

	return e.Encode(p)
}

func (p *ResponseFromRRM) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)

	return e.Encode(p)
}

func (p *PostCRVstatus) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)

	return e.Encode(p)
}

// FromJSON deserializes the object from JSON string
// in an io.Reader to the given interface
//
func FromJSON(i interface{}, r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(i)
}

func (p *RegionNodeEvents) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)

	return e.Decode(p)
}
