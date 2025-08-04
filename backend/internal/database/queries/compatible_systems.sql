-- name: FindCompatibleSystems :many
SELECT
  f.id   AS furnace_id,
  f.manufacturer   AS furnace_manufacturer,
  f.btu            AS furnace_btu,
  f.efficiency_rating AS furnace_afue,
  f.price          AS furnace_price,

  c.id   AS condenser_id,
  c.manufacturer   AS condenser_manufacturer,
  c.btu            AS condenser_btu,
  c.efficiency_rating AS condenser_afue,
  c.price          AS condenser_price,

  co.id  AS coil_id,
  co.manufacturer  AS coil_manufacturer,
  co.btu           AS coil_btu,
  co.efficiency_rating AS coil_afue,
  co.price         AS coil_price,

  CAST((f.price + c.price + co.price) AS DECIMAL) AS total_price
FROM equipment AS f
  JOIN equipment AS c  ON c.equipment_type = 'outdoor_condenser'
  JOIN equipment AS co ON co.equipment_type = 'evaporator_coil'
WHERE
  f.equipment_type = 'furnace'
  AND f.equipment_width = $1
  AND co.equipment_width = $1
  AND c.btu >= $2
  AND c.btu <= $3
  AND co.btu >= $2
  AND co.btu <= $3
ORDER BY total_price ASC;
;
