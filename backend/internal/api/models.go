package api

import (
	"github.com/google/uuid"
)

// CalculateRequest is the form data from frontend
type CalculateRequest struct {
	SquareFootage int    `json:"square_footage" validate:"required,min=500,max=10000"`
	CurrentSystem string `json:"current_system"`
	HomeAge       string `json:"home_age"`
	// Add other form fields as needed
	ExtraData map[string]interface{} `json:"extra_data,omitempty"`
}

// CalculateResponse returns the calculation results
type CalculateResponse struct {
	LeadID        uuid.UUID `json:"lead_id"`
	SquareFootage int       `json:"square_footage"`
	MinBTU        int       `json:"min_btu"`
	MaxBTU        int       `json:"max_btu"`
	MinTons       float64   `json:"min_tons"`
	MaxTons       float64   `json:"max_tons"`
}

// SystemsResponse returns compatible HVAC systems
type SystemsResponse struct {
	LeadID  uuid.UUID `json:"lead_id"`
	Systems []System  `json:"systems"`
}

// System represents a complete HVAC system bundle
type System struct {
	Furnace    EquipmentPiece `json:"furnace"`
	Condenser  EquipmentPiece `json:"condenser"`
	Coil       EquipmentPiece `json:"coil"`
	TotalPrice float64        `json:"total_price"`
}

// EquipmentPiece is a single piece of equipment
type EquipmentPiece struct {
	Model      string  `json:"model"`
	BTU        int     `json:"btu,omitempty"`
	Efficiency float64 `json:"efficiency"`
}
