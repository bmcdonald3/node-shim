// Copyright © 2025 OpenCHAMI a Series of LF Projects, LLC
//
// SPDX-License-Identifier: MIT

package v1

import (
	"context"
	"github.com/openchami/fabrica/pkg/fabrica"
)

// Campaign represents a campaign resource
type Campaign struct {
	APIVersion string           `json:"apiVersion"`
	Kind       string           `json:"kind"`
	Metadata   fabrica.Metadata `json:"metadata"`
	Spec       CampaignSpec   `json:"spec" validate:"required"`
	Status     CampaignStatus `json:"status,omitempty"`
}

// CampaignSpec defines the desired state of Campaign
type CampaignSpec struct {
TargetRef BindingTarget `json:"targetRef" validate:"required"`
Profile   string        `json:"profile" validate:"required"`
BatchSize int           `json:"batchSize" validate:"required,gt=0"`
}

// CampaignStatus defines the observed state of Campaign
type CampaignStatus struct {
Phase         string   `json:"phase,omitempty"`
Message       string   `json:"message,omitempty"`
TotalNodes    int      `json:"totalNodes,omitempty"`
AffectedNodes int      `json:"affectedNodes,omitempty"`
Bindings      []string `json:"bindings,omitempty"`
}

// Validate implements custom validation logic for Campaign
func (r *Campaign) Validate(ctx context.Context) error {
	// Add custom validation logic here
	// Example:
	// if r.Spec.Description == "forbidden" {
	//     return errors.New("description 'forbidden' is not allowed")
	// }

	return nil
}
// GetKind returns the kind of the resource
func (r *Campaign) GetKind() string {
	return "Campaign"
}

// GetName returns the name of the resource
func (r *Campaign) GetName() string {
	return r.Metadata.Name
}

// GetUID returns the UID of the resource
func (r *Campaign) GetUID() string {
	return r.Metadata.UID
}

// IsHub marks this as the hub/storage version
func (r *Campaign) IsHub() {}
