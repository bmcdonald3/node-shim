#!/bin/bash
# test_nodeset.sh

# Cleanup any lingering background processes on exit
trap 'kill $(jobs -p) 2>/dev/null' EXIT

# Start a dummy SMD server
cat <<EOF > smd_mock.go
package main
import (
	"encoding/json"
	"net/http"
)
func main() {
	http.HandleFunc("/hms/v2/Inventory/Hardware/Nodes", func(w http.ResponseWriter, r *http.Request) {
		nodes := []map[string]interface{}{
			{
				"ID": "x3000c0s1b0n0",
				"State": "Populated",
				"Labels": map[string]string{"role": "compute"},
			},
		}
		json.NewEncoder(w).Encode(nodes)
	})
	http.ListenAndServe(":2379", nil)
}
EOF
go run smd_mock.go &
SMD_PID=\$!

# Start the API server in the background
SMD_URL=http://localhost:2379 go run ../../cmd/server/*.go &

# Wait for the server to be ready
for i in {1..10}; do
    if curl -s http://localhost:8080/health > /dev/null; then
        break
    fi
    sleep 1
done

echo "--> Creating NodeSet"
curl -s -X POST http://localhost:8080/nodesets \
  -H "Content-Type: application/json" \
  -d '{"apiVersion":"v1","kind":"NodeSet","metadata":{"name":"compute-nodes"},"spec":{"selector":{"role":"compute"}}}'

echo "--> Fetching NodeSet"
HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/nodesets/compute-nodes)

if [ "$HTTP_STATUS" -eq 200 ]; then
    echo "SUCCESS: NodeSet retrieved successfully."
    curl -s http://localhost:8080/nodesets/compute-nodes
    exit 0
else
    echo "FAIL: API returned HTTP $HTTP_STATUS"
    exit 1
fi