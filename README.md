# OpenCHAMI Node Service

`node-service` is a Fabrica-based facade providing a node-centric API. It aggregates data from the State Management Database (SMD), `metadata-service`, and `boot-service`.

This service decouples configuration intent from physical SMD group membership. Operators use `ProfileBindings` to assign configuration profiles to nodes or groups of nodes (`NodeSets`). 

## Architecture

The service utilizes Fabrica's declarative Reconciliation Framework. State changes are handled asynchronously:
1. An operator creates a `ProfileBinding`.
2. The storage layer persists the desired state and emits a lifecycle event.
3. The `ProfileBindingReconciler` detects the event, evaluates the target `NodeSet`, and executes the necessary API calls to `metadata-service` and `boot-service`.
4. Upon successful downstream application, the reconciler updates the binding status phase to `Synced`. Transient failures are automatically requeued with exponential backoff.

## Core Resources

* **Node:** A composed, read-only view aggregating hardware identity (SMD), effective config groups (metadata-service), and effective boot parameters (boot-service). Evaluated at request time.
* **NodeSet:** A dynamic grouping primitive. Resolves node membership in-memory by evaluating a label selector against the SMD inventory.
* **ProfileBinding:** Binds a configuration profile to a `Node` or `NodeSet`. Materialized asynchronously by the reconciliation controller.
* **Campaign:** Orchestrates the phased rollout of a profile to a `NodeSet` using configurable batch sizes.

## Environment Configuration

The service requires network access to the downstream OpenCHAMI services. Configure these via environment variables. If unset, they default to standard local hostnames.

* `SMD_URL` (Default: `http://smd:27779`)
* `METADATA_SERVICE_URL` (Default: `http://metadata-service:8080`)
* `BOOT_SERVICE_URL` (Default: `http://boot-service:8080`)

## Running the Service

Start the API server and the reconciliation controller:
```bash
go run ./cmd/server/*.go
```

## Integration Testing
Integration tests run against a local mock server to simulate the downstream services.

To execute the test suite:

```bash
./scripts/integration/test_nodeset.sh
./scripts/integration/test_binding.sh
./scripts/integration/test_campaign.sh
```