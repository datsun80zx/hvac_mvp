-- +goose Up
CREATE TABLE compatible_systems (
    id UUID PRIMARY KEY,
    furnace_id UUID NOT NULL REFERENCES equipment(id),
    condenser_id UUID NOT NULL REFERENCES equipment(id),
    coil_id UUID NOT NULL REFERENCES equipment(id),
    total_price DECIMAL,

    CONSTRAINT unique_equipment_combo UNIQUE(furnace_id, condenser_id, coil_id),
    CONSTRAINT no_self_reference CHECK (
        furnace_id <> condenser_id
        AND condenser_id <> coil_id
        AND furnace_id <> coil_id
        )
);

-- +goose Down
DROP TABLE compatible_systems;