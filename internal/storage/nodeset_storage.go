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
func GetNodeSet(ctx context.Context, uid string) (*v1.NodeSet, error) {
ns, err := LoadNodeSet(ctx, uid)
if err != nil {
return nil, err
}

if err := ResolveNodeSet(ctx, ns); err != nil {
return nil, err
}

return ns, nil
}

// ListNodeSets retrieves all NodeSets and resolves them
func ListNodeSets(ctx context.Context) ([]*v1.NodeSet, error) {
nss, err := LoadAllNodeSets(ctx)
if err != nil {
return nil, err
}

for _, ns := range nss {
if err := ResolveNodeSet(ctx, ns); err != nil {
return nil, err
}
}

return nss, nil
}