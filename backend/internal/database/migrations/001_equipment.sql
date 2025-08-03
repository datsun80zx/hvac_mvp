-- +goose Up
CREATE TYPE equipment_type AS ENUM ('furnace', 'outdoor_condenser', 'evaporator_coil');
CREATE TABLE equipment (
    id UUID PRIMARY KEY,
    manufacturer TEXT,
    model_number TEXT UNIQUE NOT NULL,
    equipment_type equipment_type,
    btu INTEGER,
    efficiency_rating DECIMAL,
    equipment_length DECIMAL,
    equipment_width DECIMAL,
    equipment_height DECIMAL,
    price DECIMAL
);

-- +goose Down
DROP TABLE equipment;
DROP TYPE IF EXISTS equipment_type;