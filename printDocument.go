package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PrintParams struct {
	Paper       string
	Orientation int
	// orientation refer to document info
	Pages  string
	Copies int
}

func printDocumentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if len(fileDest) == 0 {
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}

	if r.Body == http.NoBody {
		http.Error(w, "no params", http.StatusBadRequest)
		return
	}
	pdfPath := recentPdf

	// parse custom config
	decoder := json.NewDecoder(r.Body)
	var customParams PrintParams

	err := decoder.Decode(&customParams)
	if err != nil {
		fmt.Print(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if err := PrintDocument(pdfPath,
		map[string]string{
			"orientation-requested": fmt.Sprint(customParams.Orientation),
			"media":                 customParams.Paper,
			"pages-ranges":          customParams.Pages,
			"copies":                fmt.Sprint(customParams.Copies),
		}); err != nil {
		http.Error(w, "error printing"+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "successfully printed")

	// delete file
	// go scheduleFileDeletion(fileDest, 1)
}
