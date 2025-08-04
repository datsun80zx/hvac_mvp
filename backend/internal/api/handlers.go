package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sqlc-dev/pqtype"

	db "github.com/datsun80zx/hvac_mvp/backend/internal/database/sqlc" // Update with your module name
)

type Handler struct {
	db *db.Queries
}

func NewHandler(database *sql.DB) *Handler {
	return &Handler{
		db: db.New(database),
	}
}

// CalculateHandler handles POST /api/calculate
func (h *Handler) CalculateHandler(w http.ResponseWriter, r *http.Request) {
	var req CalculateRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Calculate BTU range
	// Rule: 400-600 sq ft per ton, 1 ton = 12,000 BTU
	minBTU := (req.SquareFootage * 12000) / 600 // Most efficient
	maxBTU := (req.SquareFootage * 12000) / 400 // Least efficient

	// Round to nearest 6,000 BTU (0.5 ton)
	minBTU = ((minBTU + 5999) / 6000) * 6000
	maxBTU = ((maxBTU + 5999) / 6000) * 6000

	// Create form data JSON
	formData, err := json.Marshal(map[string]interface{}{
		"square_footage": req.SquareFootage,
		"current_system": req.CurrentSystem,
		"home_age":       req.HomeAge,
		"extra_data":     req.ExtraData,
		"timestamp":      time.Now(),
	})
	if err != nil {
		http.Error(w, "Failed to process form data", http.StatusInternalServerError)
		return
	}

	// Store in database
	leadID := uuid.New()
	_, err = h.db.CreateSystemUser(r.Context(), db.CreateSystemUserParams{
		ID:           leadID,
		FormData:     pqtype.NullRawMessage{RawMessage: formData, Valid: true},
		NeededMinBtu: sql.NullString{String: fmt.Sprintf("%d", minBTU), Valid: true},
		NeededMaxBtu: sql.NullString{String: fmt.Sprintf("%d", maxBTU), Valid: true},
	})

	if err != nil {
		http.Error(w, "Failed to save calculation", http.StatusInternalServerError)
		return
	}

	// Return response
	resp := CalculateResponse{
		LeadID:        leadID,
		SquareFootage: req.SquareFootage,
		MinBTU:        minBTU,
		MaxBTU:        maxBTU,
		MinTons:       float64(minBTU) / 12000,
		MaxTons:       float64(maxBTU) / 12000,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// SystemsHandler handles GET /api/systems/{lead_id}
func (h *Handler) SystemsHandler(w http.ResponseWriter, r *http.Request) {
	// Get lead ID from URL
	vars := mux.Vars(r)
	leadID, err := uuid.Parse(vars["lead_id"])
	if err != nil {
		http.Error(w, "Invalid lead ID", http.StatusBadRequest)
		return
	}

	// Get the lead/system user
	systemUser, err := h.db.GetSystemUser(r.Context(), leadID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Lead not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to retrieve lead", http.StatusInternalServerError)
		return
	}

	// Parse BTU values
	minBTU, _ := strconv.Atoi(systemUser.NeededMinBtu.String)
	maxBTU, _ := strconv.Atoi(systemUser.NeededMaxBtu.String)

	// Find compatible systems
	// For MVP, we'll use a standard width - you can make this dynamic later
	standardWidth := sql.NullString{String: "21", Valid: true}

	// Based on your SQLC generated code, the params are positional
	compatibleSystems, err := h.db.FindCompatibleSystems(r.Context(), standardWidth)

	if err != nil {
		http.Error(w, "Failed to find compatible systems", http.StatusInternalServerError)
		return
	}

	// Filter results by BTU range
	var filteredSystems []db.FindCompatibleSystemsRow
	for _, cs := range compatibleSystems {
		if cs.CondenserBtu.Valid &&
			cs.CondenserBtu.Int32 >= int32(minBTU) &&
			cs.CondenserBtu.Int32 <= int32(maxBTU) {
			filteredSystems = append(filteredSystems, cs)
		}
	}

	// Format response
	systems := make([]System, len(filteredSystems))
	for i, cs := range filteredSystems {
		// Get manufacturer and model info
		furnaceModel := ""
		if cs.FurnaceManufacturer.Valid {
			furnaceModel = cs.FurnaceManufacturer.String
		}

		condenserModel := ""
		if cs.CondenserManufacturer.Valid {
			condenserModel = cs.CondenserManufacturer.String
		}
		if cs.CondenserBtu.Valid {
			condenserModel += fmt.Sprintf(" %.1f ton", float64(cs.CondenserBtu.Int32)/12000)
		}

		coilModel := ""
		if cs.CoilManufacturer.Valid {
			coilModel = cs.CoilManufacturer.String + " Coil"
		}

		systems[i] = System{
			Furnace: EquipmentPiece{
				Model:      furnaceModel,
				BTU:        int(cs.FurnaceBtu.Int32),
				Efficiency: parseEfficiency(cs.FurnaceAfue.String),
			},
			Condenser: EquipmentPiece{
				Model:      condenserModel,
				BTU:        int(cs.CondenserBtu.Int32),
				Efficiency: parseEfficiency(cs.CondenserAfue.String), // Note: using AFUE for now
			},
			Coil: EquipmentPiece{
				Model: coilModel,
				BTU:   int(cs.CoilBtu.Int32),
			},
			TotalPrice: float64(cs.TotalPrice), // TotalPrice is already int32 from SQLC
		}
	}

	resp := SystemsResponse{
		LeadID:  leadID,
		Systems: systems,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Helper function to parse efficiency ratings
func parseEfficiency(s string) float64 {
	if s == "" {
		return 0
	}
	val, _ := strconv.ParseFloat(s, 64)
	return val
}
