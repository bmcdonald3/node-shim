#!/bin/bash
trap 'kill $(jobs -p) 2>/dev/null' EXIT
cat << 'EOF' > mock_downstream.go
package main
import ("log"; "net/http")
func main() { http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }); log.Fatal(http.ListenAndServe(":8081", nil)) }
EOF
go run smd_mock.go &
go run mock_downstream.go &
sleep 2
SMD_URL=http://localhost:2379 METADATA_SERVICE_URL=http://localhost:8081 BOOT_SERVICE_URL=http://localhost:8081 go run ../../cmd/server/*.go &
sleep 3
curl -s -X POST http://localhost:8080/nodesets -H "Content-Type: application/json" -d '{"apiVersion":"v1","kind":"NodeSet","metadata":{"name":"campaign-nodes"},"spec":{"selector":{"role":"compute"}}}' > /dev/null
RESPONSE=$(curl -s -X POST http://localhost:8080/campaigns -H "Content-Type: application/json" -d '{"apiVersion":"v1","kind":"Campaign","metadata":{"name":"canary-rollout"},"spec":{"targetRef":{"kind":"NodeSet","name":"campaign-nodes"},"profile":"v3","batchSize":2}}')
HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST http://localhost:8080/campaigns -H "Content-Type: application/json" -d '{"apiVersion":"v1","kind":"Campaign","metadata":{"name":"canary-rollout"},"spec":{"targetRef":{"kind":"NodeSet","name":"campaign-nodes"},"profile":"v3","batchSize":2}}')

if [ "$HTTP_STATUS" -eq 201 ]; then
    echo "SUCCESS: Campaign created."
    exit 0
else
    echo "FAIL: API returned HTTP $HTTP_STATUS."
    echo "Response: $RESPONSE"
    exit 1
fi
