ALTER TABLE mapping.mappings DROP CONSTRAINT IF EXISTS mappings_token_key;
ALTER TABLE mapping.mappings DROP COLUMN IF EXISTS token;
ALTER TABLE mapping.kinds DROP COLUMN IF EXISTS short_name;
