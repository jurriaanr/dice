package hello

import (
	"legion/dice/attack"
	"encoding/json"
	//"fmt"
	"net/http"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Hello, World!")

	results := make([...]int, 5)

	for i := 0; i < 5; i++ {
		results[i] = attack.RoleResult(3, 0, 0, true)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
