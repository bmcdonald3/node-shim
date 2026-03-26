## OpenCHAMI Node Service

`node-service` is a Fabrica-based facade that presents a coherent node-centric API. It composes data from the State Management Database (SMD), `metadata-service`, and `boot-service`.

This service decouples configuration intent from SMD group membership. Instead of moving nodes between physical inventory groups, operators use `ProfileBindings` to target nodes or groups of nodes (`NodeSets`) with specific configuration profiles. `node-service` acts as a write-through cache, materializing these bindings into the downstream services.

### Core Resources

* **Node:** A composed view aggregating hardware identity (SMD), effective config groups (metadata-service), and effective boot parameters (boot-service). Evaluated at request time.
* **NodeSet:** A dynamic grouping primitive. Resolves node membership in-memory by evaluating a label selector against the SMD inventory.
* **ProfileBinding:** Binds a specific configuration profile to a `Node` or `NodeSet`. Materializes the profile directly into `metadata-service` and `boot-service`.
* **Campaign:** Orchestrates the phased rollout of a profile to a `NodeSet` based on a configured `batchSize`.

### Environment Configuration

The service requires network access to the downstream OpenCHAMI services. Configure these via environment variables. If unset, they default to standard local hostnames.

* `SMD_URL` (Default: `http://smd:27779`)
* `METADATA_SERVICE_URL` (Default: `http://metadata-service:8080`)
* `BOOT_SERVICE_URL` (Default: `http://boot-service:8080`)

### Running the Service

Start the API server:
```bash
go run ./cmd/server/*.go
```