ALTER TABLE mapping.kinds ADD COLUMN IF NOT EXISTS mask VARCHAR(255) NOT NULL DEFAULT '';

UPDATE mapping.kinds SET mask = '^[А-ЯЁ][а-яёА-ЯЁ\- ]{1,99}$' WHERE name = 'name';
UPDATE mapping.kinds SET mask = '^\+7\d{10}$' WHERE name = 'phone';
UPDATE mapping.kinds SET mask = '^\d{4} \d{6}$' WHERE name = 'passport';
UPDATE mapping.kinds SET mask = '^[\w.+-]+@[\w-]+\.[a-zA-Z]{2,}$' WHERE name = 'email';
UPDATE mapping.kinds SET mask = '^[А-ЯЁа-яё0-9.,\- ]{5,200}$' WHERE name = 'address';
UPDATE mapping.kinds SET mask = '^\d{2}\.\d{2}\.\d{4}$' WHERE name = 'birth_date';
UPDATE mapping.kinds SET mask = '^\d{3}-\d{3}-\d{3} \d{2}$' WHERE name = 'snils';
UPDATE mapping.kinds SET mask = '^\d{10}(\d{2})?$' WHERE name = 'inn';
UPDATE mapping.kinds SET mask = '^\d{2} \d{2} \d{6}$' WHERE name = 'driver_license';
UPDATE mapping.kinds SET mask = '^\d{4} \d{4} \d{4} \d{4}$' WHERE name = 'bank_card';
UPDATE mapping.kinds SET mask = '^\d{20}$' WHERE name = 'account_number';
UPDATE mapping.kinds SET mask = '^((25[0-5]|2[0-4]\d|1?\d?\d)\.){3}(25[0-5]|2[0-4]\d|1?\d?\d)$' WHERE name = 'ip_address';
