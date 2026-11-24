ALTER TABLE locations
    DROP CONSTRAINT IF EXISTS locations_coordinates_pair,
    DROP CONSTRAINT IF EXISTS locations_longitude_range,
    DROP CONSTRAINT IF EXISTS locations_latitude_range,
    DROP COLUMN IF EXISTS longitude,
    DROP COLUMN IF EXISTS latitude;
