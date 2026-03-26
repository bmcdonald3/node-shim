package storage

import (
"context"
"fmt"
"os"

"github.com/user/node-service/apis/example.fabrica.dev/v1"
"github.com/user/node-service/internal/client"
)

var smdClient *client.SMDClient

func init() {
smdURL := os.Getenv("SMD_URL")
smdClient = client.NewSMDClient(smdURL)
}

// GetNode retrieves a Node, composing it from local storage and SMD
func GetNode(ctx context.Context, uid string) (*v1.Node, error) {
// 1. Get the base node from local storage (FileBackend)
node, err := LoadNode(ctx, uid)
if err != nil {
return nil, err
}

// 2. Fetch hardware state from SMD
smdNodes, err := smdClient.ListNodes(ctx)
if err != nil {
return nil, fmt.Errorf("failed to fetch nodes from SMD: %w", err)
}

// 3. Find the specific node in SMD data and map it
for _, sn := range smdNodes {
if sn.ID == node.Spec.XName {
node.Status.InventoryStatus = sn.State
node.Status.LastDiscovery = sn.LastUpdate.String()
// Merge labels from SMD into Node Spec
if node.Spec.Labels == nil {
node.Spec.Labels = make(map[string]string)
}
for k, v := range sn.Labels {
node.Spec.Labels[k] = v
}
return node, nil
}
}

return node, nil
}

// ListNodes retrieves all Nodes, composing them from local storage and SMD
func ListNodes(ctx context.Context) ([]*v1.Node, error) {
// 1. Get base nodes from local storage
nodes, err := LoadAllNodes(ctx)
if err != nil {
return nil, err
}

// 2. Fetch all nodes from SMD
smdNodes, err := smdClient.ListNodes(ctx)
if err != nil {
return nil, fmt.Errorf("failed to fetch nodes from SMD: %w", err)
}

// Create a map for quick lookup
smdMap := make(map[string]client.SMDNode)
for _, sn := range smdNodes {
smdMap[sn.ID] = sn
}

// 3. Compose each node
for _, node := range nodes {
if sn, ok := smdMap[node.Spec.XName]; ok {
node.Status.InventoryStatus = sn.State
node.Status.LastDiscovery = sn.LastUpdate.String()
// Merge labels from SMD into Node Spec
if node.Spec.Labels == nil {
node.Spec.Labels = make(map[string]string)
}
for k, v := range sn.Labels {
node.Spec.Labels[k] = v
}
}
}

return nodes, nil
}