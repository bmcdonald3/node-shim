package reconciler

import (
"context"
"fmt"
"log"
"os"
"time"

"github.com/openchami/fabrica/pkg/reconcile"
"github.com/user/node-service/apis/example.fabrica.dev/v1"
"github.com/user/node-service/internal/client"
"github.com/user/node-service/internal/storage"
)

// ProfileBindingReconciler reconciles ProfileBinding resources
type ProfileBindingReconciler struct {
	reconcile.BaseReconciler
	metadataClient *client.MetadataClient
	bootClient     *client.BootClient
}

// NewProfileBindingReconciler creates a new ProfileBinding reconciler
func NewProfileBindingReconciler() *ProfileBindingReconciler {
metadataURL := os.Getenv("METADATA_SERVICE_URL")
bootURL := os.Getenv("BOOT_SERVICE_URL")

return &ProfileBindingReconciler{
metadataClient: client.NewMetadataClient(metadataURL),
bootClient:     client.NewBootClient(bootURL),
}
}

// GetResourceKind returns the resource kind this reconciler handles
func (r *ProfileBindingReconciler) GetResourceKind() string {
return "ProfileBinding"
}

// Reconcile implements the reconciliation logic for ProfileBinding
func (r *ProfileBindingReconciler) Reconcile(ctx context.Context, resource interface{}) (reconcile.Result, error) {
// Type assert to ProfileBinding
pb, ok := resource.(*v1.ProfileBinding)
if !ok {
return reconcile.Result{}, fmt.Errorf("expected *v1.ProfileBinding, got %T", resource)
}

log.Printf("Reconciling ProfileBinding: %s", pb.Metadata.UID)

// Skip if already synced
if pb.Status.Phase == "Synced" {
log.Printf("ProfileBinding %s already synced", pb.Metadata.UID)
return reconcile.Result{}, nil
}

	// 2. Resolve affected xnames based on target type
	var xnames []string
	switch pb.Spec.TargetRef.Kind {
	case "Node":
		// Try to load the Node resource
		node, err := storage.LoadNode(ctx, pb.Spec.TargetRef.Name)
		if err == nil {
			xnames = append(xnames, node.Spec.XName)
		} else {
			// Assume Name is already an xname (for campaign use case)
			xnames = append(xnames, pb.Spec.TargetRef.Name)
		}
	case "NodeSet":
		// Resolve the NodeSet to get xnames
		ns, err := storage.GetNodeSet(ctx, pb.Spec.TargetRef.Name)
		if err != nil {
			log.Printf("Failed to resolve NodeSet %s: %v, will retry", pb.Spec.TargetRef.Name, err)
			return reconcile.Result{RequeueAfter: 5 * time.Second}, fmt.Errorf("failed to load target nodeset: %w", err)
		}
		xnames = ns.Status.ResolvedXNames
	default:
		log.Printf("Unsupported target kind: %s", pb.Spec.TargetRef.Kind)
		return reconcile.Result{}, fmt.Errorf("unsupported target kind: %s", pb.Spec.TargetRef.Kind)
	}

	// 3. Materialize to downstream services
	for _, xname := range xnames {
		// Update metadata-service
		if err := r.metadataClient.UpdateProfile(ctx, xname, pb.Spec.Profile); err != nil {
			log.Printf("Failed to update metadata-service for %s: %v, will retry", xname, err)
			return reconcile.Result{RequeueAfter: 5 * time.Second}, fmt.Errorf("failed to update metadata-service for %s: %w", xname, err)
		}

		// Update boot-service
		if err := r.bootClient.UpdateProfile(ctx, xname, pb.Spec.Profile); err != nil {
			log.Printf("Failed to update boot-service for %s: %v, will retry", xname, err)
			return reconcile.Result{RequeueAfter: 5 * time.Second}, fmt.Errorf("failed to update boot-service for %s: %w", xname, err)
		}
	}

	// 4. Update status to Synced
	pb.Status.Phase = "Synced"
	pb.Status.AffectedNodes = len(xnames)

	if err := storage.SaveProfileBinding(ctx, pb); err != nil {
		log.Printf("Failed to update ProfileBinding status: %v, will retry", err)
		return reconcile.Result{RequeueAfter: 2 * time.Second}, fmt.Errorf("failed to update status: %w", err)
	}

log.Printf("Successfully reconciled ProfileBinding %s (affected %d nodes)", pb.Metadata.UID, len(xnames))
return reconcile.Result{}, nil
}
