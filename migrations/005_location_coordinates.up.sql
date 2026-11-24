-- Exact position of a location on the map.
-- Columns are nullable so locations created before this migration stay valid;
-- the pair check guarantees latitude and longitude are always set together.
ALTER TABLE locations
    ADD COLUMN latitude  DOUBLE PRECISION,
    ADD COLUMN longitude DOUBLE PRECISION,
    ADD CONSTRAINT locations_latitude_range   CHECK (latitude  BETWEEN -90  AND 90),
    ADD CONSTRAINT locations_longitude_range  CHECK (longitude BETWEEN -180 AND 180),
    ADD CONSTRAINT locations_coordinates_pair CHECK ((latitude IS NULL) = (longitude IS NULL));
