// Copyright © 2025 OpenCHAMI a Series of LF Projects, LLC
//
// SPDX-License-Identifier: MIT

package v1

import (
	"context"
	"github.com/openchami/fabrica/pkg/fabrica"
)

// ProfileBinding represents a profilebinding resource
type ProfileBinding struct {
	APIVersion string           `json:"apiVersion"`
	Kind       string           `json:"kind"`
	Metadata   fabrica.Metadata `json:"metadata"`
	Spec       ProfileBindingSpec   `json:"spec" validate:"required"`
	Status     ProfileBindingStatus `json:"status,omitempty"`
}

type ProfileBindingSpec struct {
	// TargetRef identifies what this binding applies to
	TargetRef BindingTarget `json:"targetRef" validate:"required"`
	// Profile is the name of the profile to apply
	Profile string `json:"profile" validate:"required"`
	// BootProfileOverride optionally forces a specific boot configuration
	BootProfileOverride string `json:"bootProfileOverride,omitempty"`
	// ConfigGroupOverrides optionally appends specific metadata groups
	ConfigGroupOverrides []string `json:"configGroupOverrides,omitempty"`
}

type BindingTarget struct {
	// Kind must be either "Node" or "NodeSet"
	Kind string `json:"kind" validate:"oneof=Node NodeSet"`
	// Name is the name of the target resource
	Name string `json:"name" validate:"required"`
}

type ProfileBindingStatus struct {
	// Phase tracks the materialization (Pending, Synced, Failed)
	Phase string `json:"phase,omitempty"`
	// Message provides details on sync failures (e.g., "boot-service unreachable")
	Message string `json:"message,omitempty"`
	// AffectedNodes is a count of nodes currently impacted by this binding
	AffectedNodes int `json:"affectedNodes"`
}

// Validate implements custom validation logic for ProfileBinding
func (r *ProfileBinding) Validate(ctx context.Context) error {
	// Add custom validation logic here
	// Example:
	// if r.Spec.Description == "forbidden" {
	//     return errors.New("description 'forbidden' is not allowed")
	// }

	return nil
}
// GetKind returns the kind of the resource
func (r *ProfileBinding) GetKind() string {
	return "ProfileBinding"
}

// GetName returns the name of the resource
func (r *ProfileBinding) GetName() string {
	return r.Metadata.Name
}

// GetUID returns the UID of the resource
func (r *ProfileBinding) GetUID() string {
	return r.Metadata.UID
}

// IsHub marks this as the hub/storage version
func (r *ProfileBinding) IsHub() {}
