package main

import (
	"encoding/json"
	"net/http"
)

func getAllPrintersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	printers, err := GetAllPrinters()
	if err != nil {
		http.Error(w, "Error getting printers", http.StatusInternalServerError)
		return
	}

	data := map[string][]string{"printers": printers}

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonData)

}
