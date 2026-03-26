package storage

import (
	"context"
	"fmt"
	"time"

"github.com/openchami/fabrica/pkg/fabrica"
"github.com/openchami/fabrica/pkg/resource"
"github.com/user/node-service/apis/example.fabrica.dev/v1"
)

// CreateCampaign implements the logic to create profile bindings for a batch of nodes
func CreateCampaign(ctx context.Context, c *v1.Campaign) error {
	// 1. Fetch the target NodeSet
	if c.Spec.TargetRef.Kind != "NodeSet" {
		return fmt.Errorf("campaign target must be a NodeSet, got %s", c.Spec.TargetRef.Kind)
	}

	ns, err := GetNodeSet(ctx, c.Spec.TargetRef.Name)
	if err != nil {
		return fmt.Errorf("failed to fetch target NodeSet %s: %w", c.Spec.TargetRef.Name, err)
	}

	// 2. Select nodes up to BatchSize
	totalNodes := len(ns.Status.ResolvedXNames)
	batchSize := c.Spec.BatchSize
	if batchSize > totalNodes {
		batchSize = totalNodes
	}

	selectedXNames := ns.Status.ResolvedXNames[:batchSize]

	// 3. Create ProfileBindings for each selected node
	var bindingUIDs []string
	for _, xname := range selectedXNames {
		uid, err := resource.GenerateUIDForResource("ProfileBinding")
		if err != nil {
			return fmt.Errorf("failed to generate UID for ProfileBinding: %w", err)
		}

pb := &v1.ProfileBinding{
APIVersion: "v1",
Kind:       "ProfileBinding",
Metadata: fabrica.Metadata{
UID:       uid,
Name:      fmt.Sprintf("%s-%s", c.Metadata.Name, xname),
CreatedAt: time.Now(),
UpdatedAt: time.Now(),
},
			Spec: v1.ProfileBindingSpec{
				TargetRef: v1.BindingTarget{
					Kind: "Node",
					Name: xname,
				},
				Profile: c.Spec.Profile,
			},
		}

		if err := CreateProfileBinding(ctx, pb); err != nil {
			return fmt.Errorf("failed to create ProfileBinding for node %s: %w", xname, err)
		}
		bindingUIDs = append(bindingUIDs, uid)
	}

	// 4. Update Campaign status
	c.Status.Phase = "Completed"
	c.Status.TotalNodes = totalNodes
	c.Status.AffectedNodes = len(selectedXNames)
	c.Status.Bindings = bindingUIDs

	// 5. Save Campaign
	return SaveCampaign(ctx, c)
}