-- name: CreateEquipment :one
INSERT INTO equipment (id, manufacturer, model_number, equipment_type, btu, efficiency_rating, price, equipment_length, equipment_width, equipment_height)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10
)
RETURNING *;

-- name: GetEquipment :many
SELECT * FROM equipment
WHERE equipment_type = $1
    AND equipment_width = $2;