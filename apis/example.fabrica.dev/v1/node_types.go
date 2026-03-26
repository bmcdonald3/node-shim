// Copyright © 2025 OpenCHAMI a Series of LF Projects, LLC
//
// SPDX-License-Identifier: MIT

package v1

import (
	"context"
	"github.com/openchami/fabrica/pkg/fabrica"
)

// Node represents a node resource
type Node struct {
	APIVersion string           `json:"apiVersion"`
	Kind       string           `json:"kind"`
	Metadata   fabrica.Metadata `json:"metadata"`
	Spec       NodeSpec   `json:"spec" validate:"required"`
	Status     NodeStatus `json:"status,omitempty"`
}

type NodeSpec struct {
	// XName is the Hardware Management Services (HMS) identifier (e.g., x3000c0s1b0n0)
	XName string `json:"xname" validate:"required"`
	// Role defines the primary function (e.g., compute, service)
	Role string `json:"role,omitempty"`
	// Subrole defines a more specific function
	Subrole string `json:"subrole,omitempty"`
	// Labels are key-value pairs used for NodeSet selection
	Labels map[string]string `json:"labels,omitempty"`
}

type NodeStatus struct {
	// EffectiveProfile is the profile currently applied (Explicit > Binding > Default)
	EffectiveProfile string `json:"effectiveProfile,omitempty"`
	// BootProfileRef is the specific boot configuration currently resolving
	BootProfileRef string `json:"bootProfileRef,omitempty"`
	// ConfigGroups are the metadata-service groups contributing to this node's config
	ConfigGroups []string `json:"configGroups,omitempty"`
	// InventoryStatus is the raw state from SMD (e.g., "Populated", "Empty")
	InventoryStatus string `json:"inventoryStatus,omitempty"`
	// LastDiscovery is when SMD last saw this node
	LastDiscovery string `json:"lastDiscovery,omitempty"`
}

// Validate implements custom validation logic for Node
func (r *Node) Validate(ctx context.Context) error {
	// Add custom validation logic here
	// Example:
	// if r.Spec.Description == "forbidden" {
	//     return errors.New("description 'forbidden' is not allowed")
	// }

	return nil
}
// GetKind returns the kind of the resource
func (r *Node) GetKind() string {
	return "Node"
}

// GetName returns the name of the resource
func (r *Node) GetName() string {
	return r.Metadata.Name
}

// GetUID returns the UID of the resource
func (r *Node) GetUID() string {
	return r.Metadata.UID
}

// IsHub marks this as the hub/storage version
func (r *Node) IsHub() {}
