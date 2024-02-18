package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-pdf/fpdf"
	"github.com/gorilla/mux"
)

type WorkEntry struct {
	Date       string  `json:"date"`
	Start      string  `json:"start"`
	End        string  `json:"end"`
	TotalHours float64 `json:"totalHours"`
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set headers
		w.Header().Set("Access-Control-Allow-Origin", "*")                                                                                   // Allow any origin
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")                                                    // Allow specific methods
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization") // Allow specific headers

		// If it's a preflight OPTIONS request, send a 200 response
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func main() {
	router := mux.NewRouter()

	handler := corsMiddleware(router)

	// router.HandleFunc("/entries", GetEntries).Methods("GET")
	// router.HandleFunc("/entry", CreateEntry).Methods("POST")
	router.HandleFunc("/generate", GeneratedPdf).Methods("POST")

	log.Fatal(http.ListenAndServe(":1997", handler))
}

func GeneratedPdf(w http.ResponseWriter, r *http.Request) {
	var entries []WorkEntry
	err := json.NewDecoder(r.Body).Decode(&entries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 12)

	header := []string{"Datum", "Start", "Eind", "Gewerkt"}
	eachColWidth := 45.0
	headerHeight := 10.0
	for _, h := range header {
		pdf.CellFormat(eachColWidth, headerHeight, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 12) // Reset font for table content

	var totalHours float64
	var entryCountPerPage int = 0
	var subTotal float64
	for _, entry := range entries {
		if entryCountPerPage >= 24 {
			// Add subtotal before adding a new page
			pdf.CellFormat(0, 10, fmt.Sprintf("Subtotaal: %.2f", subTotal), "0", 0, "R", false, 0, "")
			pdf.AddPage()
			entryCountPerPage = 0 // Reset entry count for the new page
			// Re-add the header on the new page
			for _, h := range header {
				pdf.CellFormat(eachColWidth, headerHeight, h, "1", 0, "C", false, 0, "")
			}
			pdf.Ln(-1)
			subTotal = 0
		}

		pdf.CellFormat(eachColWidth, 10, entry.Date, "1", 0, "", false, 0, "")
		pdf.CellFormat(eachColWidth, 10, entry.Start, "1", 0, "", false, 0, "")
		pdf.CellFormat(eachColWidth, 10, entry.End, "1", 0, "", false, 0, "")
		pdf.CellFormat(eachColWidth, 10, fmt.Sprintf("%.2f", entry.TotalHours), "1", 0, "", false, 0, "")
		pdf.Ln(-1)
		totalHours += entry.TotalHours
		subTotal += entry.TotalHours
		entryCountPerPage++
	}

	// Check if we've added any entries to ensure we don't add a subtotal unnecessarily
	if entryCountPerPage > 0 {
		pdf.CellFormat(0, 10, fmt.Sprintf("Subtotaal: %.2f", subTotal), "0", 0, "R", false, 0, "")
		pdf.Ln(5)
	}

	// Total at the end of the document, ensuring it's always added at the bottom
	pdf.CellFormat(0, 10, fmt.Sprintf("Totaal: %.2f", totalHours), "0", 0, "R", false, 0, "")

	// Output PDF
	err = pdf.Output(w)
	if err != nil {
		fmt.Println("Failed to generate PDF:", err)
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename=\"work_entries.pdf\"")
	}
}
