package storage

import (
"context"
"fmt"
"time"

"github.com/user/node-service/apis/example.fabrica.dev/v1"
)

// ResolveNodeSet filters nodes in-memory by matching the NodeSetSpec.Selector against SMD labels
func ResolveNodeSet(ctx context.Context, ns *v1.NodeSet) error {
smdNodes, err := smdClient.ListNodes(ctx)
if err != nil {
return fmt.Errorf("failed to fetch nodes from SMD: %w", err)
}

var resolved []string
for _, sn := range smdNodes {
match := true
for k, v := range ns.Spec.Selector {
if lv, ok := sn.Labels[k]; !ok || lv != v {
match = false
break
}
}
if match {
resolved = append(resolved, sn.ID)
}
}

ns.Status.ResolvedXNames = resolved
ns.Status.NodeCount = len(resolved)
ns.Status.LastResolved = time.Now().UTC().Format(time.RFC3339)

return nil
}

// GetNodeSet retrieves a NodeSet and resolves it
func GetNodeSet(ctx context.Context, uidOrName string) (*v1.NodeSet, error) {
// 1. Get the base nodeset from local storage (FileBackend).
// Try searching by UID first.
ns, err := LoadNodeSet(ctx, uidOrName)
if err != nil {
// Try searching by name if UID lookup failed.
nss, loadErr := LoadAllNodeSets(ctx)
if loadErr != nil {
return nil, fmt.Errorf("failed to list all nodesets: %w", loadErr)
}
for _, item := range nss {
if item.Metadata.Name == uidOrName {
ns = item
err = nil
break
}
}
}

if err != nil {
return nil, fmt.Errorf("nodeset not found: %w", err)
}

if ns == nil {
return nil, fmt.Errorf("nodeset not found: %s", uidOrName)
}

// 2. Resolve the dynamic status from SMD in-memory
if err := ResolveNodeSet(ctx, ns); err != nil {
return nil, fmt.Errorf("failed to resolve nodeset from SMD: %w", err)
}

return ns, nil
}

// ListNodeSets retrieves all NodeSets and resolves them
func ListNodeSets(ctx context.Context) ([]*v1.NodeSet, error) {
	// Original Fabrica logic: Read all from local file storage first
	nss, err := LoadAllNodeSets(ctx)
	if err != nil {
		return nil, err
	}

	// Dynamic logic: Augment each loaded object with state from SMD
	for _, ns := range nss {
		if err := ResolveNodeSet(ctx, ns); err != nil {
			return nil, err
		}
	}

	return nss, nil
}
