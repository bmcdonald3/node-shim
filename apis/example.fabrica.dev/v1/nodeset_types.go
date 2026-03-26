// Copyright © 2025 OpenCHAMI a Series of LF Projects, LLC
//
// SPDX-License-Identifier: MIT

package v1

import (
	"context"
	"github.com/openchami/fabrica/pkg/fabrica"
)

// NodeSet represents a nodeset resource
type NodeSet struct {
	APIVersion string           `json:"apiVersion"`
	Kind       string           `json:"kind"`
	Metadata   fabrica.Metadata `json:"metadata"`
	Spec       NodeSetSpec   `json:"spec" validate:"required"`
	Status     NodeSetStatus `json:"status,omitempty"`
}

type NodeSetSpec struct {
	// Selector is a map of labels to match nodes against in SMD
	Selector map[string]string `json:"selector,omitempty"`
	// XNames allows for an explicit list of nodes instead of a selector
	XNames []string `json:"xnames,omitempty"`
	// Partitions limits the scope of the NodeSet to specific SMD partitions
	Partitions []string `json:"partitions,omitempty"`
}

type NodeSetStatus struct {
	// ResolvedXNames is the list of xnames currently matching the Spec criteria
	ResolvedXNames []string `json:"resolvedXNames,omitempty"`
	// NodeCount is the current count of resolved nodes
	NodeCount int `json:"nodeCount"`
	// LastResolved is the timestamp of the last in-memory filtering operation
	LastResolved string `json:"lastResolved,omitempty"`
}

// Validate implements custom validation logic for NodeSet
func (r *NodeSet) Validate(ctx context.Context) error {
	// Add custom validation logic here
	// Example:
	// if r.Spec.Description == "forbidden" {
	//     return errors.New("description 'forbidden' is not allowed")
	// }

	return nil
}
// GetKind returns the kind of the resource
func (r *NodeSet) GetKind() string {
	return "NodeSet"
}

// GetName returns the name of the resource
func (r *NodeSet) GetName() string {
	return r.Metadata.Name
}

// GetUID returns the UID of the resource
func (r *NodeSet) GetUID() string {
	return r.Metadata.UID
}

// IsHub marks this as the hub/storage version
func (r *NodeSet) IsHub() {}
