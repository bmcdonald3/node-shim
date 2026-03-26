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
