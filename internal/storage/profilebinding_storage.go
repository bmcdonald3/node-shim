package storage

import (
"context"
"fmt"
"os"

"github.com/user/node-service/apis/example.fabrica.dev/v1"
"github.com/user/node-service/internal/client"
)

var (
metadataClient *client.MetadataClient
bootClient     *client.BootClient
)

func init() {
metadataURL := os.Getenv("METADATA_URL")
if metadataURL == "" {
metadataURL = "http://localhost:8081"
}
metadataClient = client.NewMetadataClient(metadataURL)

bootURL := os.Getenv("BOOT_URL")
if bootURL == "" {
bootURL = "http://localhost:8082"
}
bootClient = client.NewBootClient(bootURL)
}

// CreateProfileBinding implements write-through logic to metadata and boot services
func CreateProfileBinding(ctx context.Context, pb *v1.ProfileBinding) error {
// 1. Resolve affected xnames
var xnames []string
switch pb.Spec.TargetRef.Kind {
case "Node":
// For a single node, we need to find its xname.
// Since LoadNode uses UID, we need to load it first.
node, err := LoadNode(ctx, pb.Spec.TargetRef.Name)
if err != nil {
return fmt.Errorf("failed to load target node: %w", err)
}
xnames = append(xnames, node.Spec.XName)
case "NodeSet":
ns, err := GetNodeSet(ctx, pb.Spec.TargetRef.Name)
if err != nil {
return fmt.Errorf("failed to load target nodeset: %w", err)
}
xnames = ns.Status.ResolvedXNames
default:
return fmt.Errorf("unsupported target kind: %s", pb.Spec.TargetRef.Kind)
}

// 2. Materialize to downstream services
for _, xname := range xnames {
if err := metadataClient.UpdateProfile(ctx, xname, pb.Spec.Profile); err != nil {
return fmt.Errorf("failed to update metadata-service for %s: %w", xname, err)
}
if err := bootClient.UpdateProfile(ctx, xname, pb.Spec.Profile); err != nil {
return fmt.Errorf("failed to update boot-service for %s: %w", xname, err)
}
}

// 3. Persist to local storage if downstream updates succeeded
pb.Status.Phase = "Synced"
pb.Status.AffectedNodes = len(xnames)

return SaveProfileBinding(ctx, pb)
}