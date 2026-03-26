package storage

import (
"context"
"os"

"github.com/user/node-service/apis/example.fabrica.dev/v1"
"github.com/user/node-service/internal/client"
)

var (
metadataClient *client.MetadataClient
bootClient     *client.BootClient
)

func init() {
metadataURL := os.Getenv("METADATA_SERVICE_URL")
metadataClient = client.NewMetadataClient(metadataURL)

bootURL := os.Getenv("BOOT_SERVICE_URL")
bootClient = client.NewBootClient(bootURL)
}

// CreateProfileBinding persists the resource with "Pending" phase.
// Reconciliation handled by the ProfileBinding reconciler.
func CreateProfileBinding(ctx context.Context, pb *v1.ProfileBinding) error {
pb.Status.Phase = "Pending"
pb.Status.AffectedNodes = 0
return SaveProfileBinding(ctx, pb)
}
