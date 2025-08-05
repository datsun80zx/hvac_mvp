// backend/cmd/seeder/main.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	db "github.com/datsun80zx/hvac_mvp/backend/internal/database/sqlc"
)

type EquipmentData struct {
	Manufacturer     string
	ModelNumber      string
	EquipmentType    string
	BTU              int32
	EfficiencyRating float64
	Length           float64
	Width            float64
	Height           float64
	Price            float64
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Get database URL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	// Connect to database
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer conn.Close()

	// Create queries instance
	queries := db.New(conn)

	// Clear existing data (optional - comment out if you want to append)
	log.Println("Clearing existing equipment data...")
	_, err = conn.Exec("DELETE FROM equipment")
	if err != nil {
		log.Fatal("Failed to clear equipment:", err)
	}

	// Seed equipment data
	log.Println("Seeding equipment data...")

	// FURNACES (various BTU sizes and efficiencies)
	furnaces := []EquipmentData{
		// Standard efficiency furnaces (80% AFUE)
		{Manufacturer: "Goodman", ModelNumber: "GM9C800603AN", EquipmentType: "furnace", BTU: 60000, EfficiencyRating: 80, Length: 33, Width: 17.5, Height: 28, Price: 1200},
		{Manufacturer: "Goodman", ModelNumber: "GM9C800804BN", EquipmentType: "furnace", BTU: 80000, EfficiencyRating: 80, Length: 33, Width: 21, Height: 28, Price: 1400},
		{Manufacturer: "Goodman", ModelNumber: "GM9C801005CN", EquipmentType: "furnace", BTU: 100000, EfficiencyRating: 80, Length: 33, Width: 21, Height: 28, Price: 1600},

		// High efficiency furnaces (96% AFUE)
		{Manufacturer: "Goodman", ModelNumber: "GMVC960603BN", EquipmentType: "furnace", BTU: 60000, EfficiencyRating: 96, Length: 34.5, Width: 17.5, Height: 29.5, Price: 2200},
		{Manufacturer: "Goodman", ModelNumber: "GMVC960804CN", EquipmentType: "furnace", BTU: 80000, EfficiencyRating: 96, Length: 34.5, Width: 21, Height: 29.5, Price: 2500},
		{Manufacturer: "Goodman", ModelNumber: "GMVC961005CN", EquipmentType: "furnace", BTU: 100000, EfficiencyRating: 96, Length: 34.5, Width: 21, Height: 29.5, Price: 2800},

		// Premium furnaces (98% AFUE)
		{Manufacturer: "Carrier", ModelNumber: "59MN7A060F17", EquipmentType: "furnace", BTU: 60000, EfficiencyRating: 98.5, Length: 35, Width: 17.5, Height: 30, Price: 3500},
		{Manufacturer: "Carrier", ModelNumber: "59MN7A080F21", EquipmentType: "furnace", BTU: 80000, EfficiencyRating: 98.5, Length: 35, Width: 21, Height: 30, Price: 3800},
		{Manufacturer: "Trane", ModelNumber: "S9V2C100D5", EquipmentType: "furnace", BTU: 100000, EfficiencyRating: 97, Length: 35, Width: 21, Height: 30, Price: 4200},
	}

	// OUTDOOR CONDENSERS (various tonnages and SEER ratings)
	condensers := []EquipmentData{
		// 2 Ton units (24,000 BTU)
		{Manufacturer: "Goodman", ModelNumber: "GSX130241", EquipmentType: "outdoor_condenser", BTU: 24000, EfficiencyRating: 13, Length: 26, Width: 26, Height: 28, Price: 1100},
		{Manufacturer: "Goodman", ModelNumber: "GSX140241", EquipmentType: "outdoor_condenser", BTU: 24000, EfficiencyRating: 14, Length: 26, Width: 26, Height: 28, Price: 1300},
		{Manufacturer: "Goodman", ModelNumber: "GSXC160241", EquipmentType: "outdoor_condenser", BTU: 24000, EfficiencyRating: 16, Length: 29, Width: 29, Height: 30, Price: 2100},

		// 3 Ton units (36,000 BTU)
		{Manufacturer: "Goodman", ModelNumber: "GSX130361", EquipmentType: "outdoor_condenser", BTU: 36000, EfficiencyRating: 13, Length: 29, Width: 29, Height: 30, Price: 1400},
		{Manufacturer: "Goodman", ModelNumber: "GSX140361", EquipmentType: "outdoor_condenser", BTU: 36000, EfficiencyRating: 14, Length: 29, Width: 29, Height: 30, Price: 1700},
		{Manufacturer: "Goodman", ModelNumber: "GSXC160361", EquipmentType: "outdoor_condenser", BTU: 36000, EfficiencyRating: 16, Length: 35, Width: 35, Height: 36, Price: 2600},

		// 4 Ton units (48,000 BTU)
		{Manufacturer: "Goodman", ModelNumber: "GSX130481", EquipmentType: "outdoor_condenser", BTU: 48000, EfficiencyRating: 13, Length: 35, Width: 35, Height: 36, Price: 1700},
		{Manufacturer: "Goodman", ModelNumber: "GSX140481", EquipmentType: "outdoor_condenser", BTU: 48000, EfficiencyRating: 14, Length: 35, Width: 35, Height: 36, Price: 2100},
		{Manufacturer: "Carrier", ModelNumber: "24ACC448", EquipmentType: "outdoor_condenser", BTU: 48000, EfficiencyRating: 17, Length: 35, Width: 35, Height: 39, Price: 3200},

		// 5 Ton units (60,000 BTU)
		{Manufacturer: "Goodman", ModelNumber: "GSX130601", EquipmentType: "outdoor_condenser", BTU: 60000, EfficiencyRating: 13, Length: 35, Width: 35, Height: 41, Price: 2000},
		{Manufacturer: "Goodman", ModelNumber: "GSX140601", EquipmentType: "outdoor_condenser", BTU: 60000, EfficiencyRating: 14, Length: 35, Width: 35, Height: 41, Price: 2400},
		{Manufacturer: "Trane", ModelNumber: "4TTR6060", EquipmentType: "outdoor_condenser", BTU: 60000, EfficiencyRating: 16, Length: 37, Width: 37, Height: 43, Price: 3800},
	}

	// EVAPORATOR COILS (matching condenser tonnages)
	coils := []EquipmentData{
		// 2 Ton coils (for 17.5" furnaces)
		{Manufacturer: "Goodman", ModelNumber: "CAPF3030A6", EquipmentType: "evaporator_coil", BTU: 24000, EfficiencyRating: 0, Length: 21, Width: 17.5, Height: 14, Price: 450},
		{Manufacturer: "Goodman", ModelNumber: "CHPF3030A6", EquipmentType: "evaporator_coil", BTU: 24000, EfficiencyRating: 0, Length: 21, Width: 17.5, Height: 14, Price: 500},

		// 3 Ton coils (for 17.5" and 21" furnaces)
		{Manufacturer: "Goodman", ModelNumber: "CAPF3636A6", EquipmentType: "evaporator_coil", BTU: 36000, EfficiencyRating: 0, Length: 21, Width: 17.5, Height: 17.5, Price: 550},
		{Manufacturer: "Goodman", ModelNumber: "CAPF3636C6", EquipmentType: "evaporator_coil", BTU: 36000, EfficiencyRating: 0, Length: 24.5, Width: 21, Height: 17.5, Price: 600},
		{Manufacturer: "Goodman", ModelNumber: "CHPF3636C6", EquipmentType: "evaporator_coil", BTU: 36000, EfficiencyRating: 0, Length: 24.5, Width: 21, Height: 17.5, Price: 650},

		// 4 Ton coils (for 21" furnaces)
		{Manufacturer: "Goodman", ModelNumber: "CAPF4860C6", EquipmentType: "evaporator_coil", BTU: 48000, EfficiencyRating: 0, Length: 24.5, Width: 21, Height: 21, Price: 700},
		{Manufacturer: "Goodman", ModelNumber: "CHPF4860C6", EquipmentType: "evaporator_coil", BTU: 48000, EfficiencyRating: 0, Length: 24.5, Width: 21, Height: 21, Price: 750},

		// 5 Ton coils (for 21" furnaces)
		{Manufacturer: "Goodman", ModelNumber: "CAPF6124D6", EquipmentType: "evaporator_coil", BTU: 60000, EfficiencyRating: 0, Length: 24.5, Width: 21, Height: 24.5, Price: 850},
		{Manufacturer: "Goodman", ModelNumber: "CHPF6124D6", EquipmentType: "evaporator_coil", BTU: 60000, EfficiencyRating: 0, Length: 24.5, Width: 21, Height: 24.5, Price: 900},
	}

	// Insert all equipment
	allEquipment := append(append(furnaces, condensers...), coils...)

	for _, eq := range allEquipment {
		err := insertEquipment(queries, eq)
		if err != nil {
			log.Printf("Failed to insert %s: %v", eq.ModelNumber, err)
		} else {
			log.Printf("Inserted %s %s", eq.EquipmentType, eq.ModelNumber)
		}
	}

	// Verify the data
	log.Println("\nVerifying seeded data:")
	verifyData(conn)

	log.Println("\nSeeding complete!")
}

func insertEquipment(queries *db.Queries, eq EquipmentData) error {
	var equipmentType db.NullEquipmentType
	equipmentType.EquipmentType = db.EquipmentType(eq.EquipmentType)
	equipmentType.Valid = true

	params := db.CreateEquipmentParams{
		ID:               uuid.New(),
		Manufacturer:     sql.NullString{String: eq.Manufacturer, Valid: true},
		ModelNumber:      eq.ModelNumber,
		EquipmentType:    equipmentType,
		Btu:              sql.NullInt32{Int32: eq.BTU, Valid: eq.BTU > 0},
		EfficiencyRating: sql.NullString{String: fmt.Sprintf("%.1f", eq.EfficiencyRating), Valid: eq.EfficiencyRating > 0},
		Price:            sql.NullString{String: fmt.Sprintf("%.2f", eq.Price), Valid: true},
		EquipmentLength:  sql.NullString{String: fmt.Sprintf("%.1f", eq.Length), Valid: true},
		EquipmentWidth:   sql.NullString{String: fmt.Sprintf("%.1f", eq.Width), Valid: true},
		EquipmentHeight:  sql.NullString{String: fmt.Sprintf("%.1f", eq.Height), Valid: true},
	}

	_, err := queries.CreateEquipment(context.Background(), params)
	return err
}

func verifyData(conn *sql.DB) {
	// Count equipment by type
	rows, err := conn.Query(`
		SELECT equipment_type, COUNT(*) 
		FROM equipment 
		GROUP BY equipment_type 
		ORDER BY equipment_type
	`)
	if err != nil {
		log.Printf("Failed to verify data: %v", err)
		return
	}
	defer rows.Close()

	fmt.Println("\nEquipment counts by type:")
	for rows.Next() {
		var eqType string
		var count int
		rows.Scan(&eqType, &count)
		fmt.Printf("  %s: %d\n", eqType, count)
	}

	// Test compatible systems query
	fmt.Println("\nTesting compatible systems query for 3-ton systems (21\" width):")
	compatRows, err := conn.Query(`
		SELECT COUNT(*) 
		FROM equipment AS f
		JOIN equipment AS c ON c.equipment_type = 'outdoor_condenser'
		JOIN equipment AS co ON co.equipment_type = 'evaporator_coil'
		WHERE f.equipment_type = 'furnace'
		AND f.equipment_width = '21'
		AND co.equipment_width = '21'
		AND c.btu >= 36000
		AND c.btu <= 36000
		AND co.btu >= 36000
		AND co.btu <= 36000
	`)
	if err != nil {
		log.Printf("Failed to test compatible systems: %v", err)
		return
	}
	defer compatRows.Close()

	if compatRows.Next() {
		var count int
		compatRows.Scan(&count)
		fmt.Printf("  Found %d compatible 3-ton systems\n", count)
	}
}
