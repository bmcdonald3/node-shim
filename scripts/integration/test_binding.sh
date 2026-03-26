#!/bin/bash
trap 'kill $(jobs -p) 2>/dev/null' EXIT

# 1. Start a mock downstream server on port 8081 to simulate metadata and boot services
cat << 'EOF' > mock_downstream.go
package main
import ("log"; "net/http")
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	log.Fatal(http.ListenAndServe(":8081", nil))
}
EOF
go run mock_downstream.go &
sleep 2

# 2. Start node-service on port 8080
METADATA_SERVICE_URL=http://localhost:8081 BOOT_SERVICE_URL=http://localhost:8081 go run ../../cmd/server/*.go &
sleep 3

# 3. Create a NodeSet
curl -s -X POST http://localhost:8080/nodesets -H "Content-Type: application/json" -d '{"apiVersion":"v1","kind":"NodeSet","metadata":{"name":"test-nodes"},"spec":{"selector":{"role":"compute"}}}' > /dev/null

# 4. Create a ProfileBinding
HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST http://localhost:8080/profilebindings -H "Content-Type: application/json" -d '{"apiVersion":"v1","kind":"ProfileBinding","metadata":{"name":"test-bind"},"spec":{"targetRef":{"kind":"NodeSet","name":"test-nodes"},"profile":"v2"}}')

if [ "$HTTP_STATUS" -eq 201 ]; then
    echo "SUCCESS: ProfileBinding created and materialized."
    exit 0
else
    echo "FAIL: API returned HTTP $HTTP_STATUS. The write-through logic to downstream services failed."
    exit 1
fi