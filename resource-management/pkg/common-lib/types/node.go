package types

import (
	"fmt"
	"strconv"

	"resource-management/pkg/common-lib/types/location"
)

// TODO - add more fields for minimal node record
type Node struct {
	id              string
	resourceVersion string
	label           string
	loc             *location.Location
}

func NewNode(id, rv, label string, location *location.Location) *Node {
	return &Node{
		id:              id,
		resourceVersion: rv,
		label:           label,
		loc:             location,
	}
}

func (n *Node) Copy() *Node {
	return &Node{
		id:              n.id,
		resourceVersion: n.resourceVersion,
		label:           n.label,
		loc:             n.loc,
	}
}

func (n *Node) GetId() string {
	return n.id
}

func (n *Node) GetLocation() *location.Location {
	return n.loc
}

func (n *Node) GetResourceVersion() uint64 {
	rv, err := strconv.ParseUint(n.resourceVersion, 10, 64)
	if err != nil {
		fmt.Printf("Unable to convert resource version %s to uint64\n", n.resourceVersion)
		return 0
	}
	return rv
}

type HardwareConfig struct {
	ConfigId string
}
