ALTER TABLE fiscal_rules
    ADD COLUMN IF NOT EXISTS rule_code VARCHAR(80),
    ADD COLUMN IF NOT EXISTS description VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS priority INTEGER NOT NULL DEFAULT 100,
    ADD COLUMN IF NOT EXISTS status VARCHAR(30) NOT NULL DEFAULT 'APPROVED',
    ADD COLUMN IF NOT EXISTS calculation_basis VARCHAR(30) NOT NULL DEFAULT 'FREIGHT_VALUE';

UPDATE fiscal_rules
SET rule_code = CONCAT(
    'RULE_',
    origin_uf,
    '_',
    destination_uf,
    '_',
    operation_type,
    '_',
    customer_type,
    '_',
    REPLACE(rule_version, '.', '_'),
    '_',
    TO_CHAR(valid_from, 'YYYY_MM_DD')
)
WHERE rule_code IS NULL;

ALTER TABLE fiscal_rules
    ALTER COLUMN rule_code SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS uq_fiscal_rules_rule_code
ON fiscal_rules (rule_code);

CREATE INDEX IF NOT EXISTS idx_fiscal_rules_engine_lookup
ON fiscal_rules (
    status,
    active,
    valid_from,
    valid_to,
    priority
);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'ck_fiscal_rules_status'
    ) THEN
        ALTER TABLE fiscal_rules
            ADD CONSTRAINT ck_fiscal_rules_status
            CHECK (status IN ('DRAFT', 'PENDING_REVIEW', 'APPROVED', 'INACTIVE'));
    END IF;

    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'ck_fiscal_rules_calculation_basis'
    ) THEN
        ALTER TABLE fiscal_rules
            ADD CONSTRAINT ck_fiscal_rules_calculation_basis
            CHECK (calculation_basis IN ('FREIGHT_VALUE'));
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS fiscal_rule_conditions (
    id BIGSERIAL PRIMARY KEY,
    fiscal_rule_id BIGINT NOT NULL REFERENCES fiscal_rules(id) ON DELETE CASCADE,
    field_name VARCHAR(60) NOT NULL,
    operator VARCHAR(30) NOT NULL,
    field_value VARCHAR(120) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_fiscal_rule_conditions
ON fiscal_rule_conditions (
    fiscal_rule_id,
    field_name,
    operator,
    field_value
);

CREATE INDEX IF NOT EXISTS idx_fiscal_rule_conditions_rule_id
ON fiscal_rule_conditions (fiscal_rule_id);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'ck_fiscal_rule_conditions_operator'
    ) THEN
        ALTER TABLE fiscal_rule_conditions
            ADD CONSTRAINT ck_fiscal_rule_conditions_operator
            CHECK (operator IN ('EQUALS', 'NOT_EQUALS', 'IN', 'BETWEEN'));
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS fiscal_rule_taxes (
    id BIGSERIAL PRIMARY KEY,
    fiscal_rule_id BIGINT NOT NULL REFERENCES fiscal_rules(id) ON DELETE CASCADE,
    tax_name VARCHAR(20) NOT NULL,
    rate NUMERIC(7,4) NOT NULL,
    base_reduction_rate NUMERIC(7,4) NOT NULL DEFAULT 0,
    calculation_order INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_fiscal_rule_taxes
ON fiscal_rule_taxes (
    fiscal_rule_id,
    tax_name
);

CREATE INDEX IF NOT EXISTS idx_fiscal_rule_taxes_rule_id
ON fiscal_rule_taxes (fiscal_rule_id);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'ck_fiscal_rule_taxes_tax_name'
    ) THEN
        ALTER TABLE fiscal_rule_taxes
            ADD CONSTRAINT ck_fiscal_rule_taxes_tax_name
            CHECK (tax_name IN ('ICMS', 'IBS', 'CBS'));
    END IF;

    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'ck_fiscal_rule_taxes_rate'
    ) THEN
        ALTER TABLE fiscal_rule_taxes
            ADD CONSTRAINT ck_fiscal_rule_taxes_rate
            CHECK (rate >= 0);
    END IF;

    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'ck_fiscal_rule_taxes_base_reduction_rate'
    ) THEN
        ALTER TABLE fiscal_rule_taxes
            ADD CONSTRAINT ck_fiscal_rule_taxes_base_reduction_rate
            CHECK (base_reduction_rate >= 0 AND base_reduction_rate <= 100);
    END IF;
END $$;

INSERT INTO fiscal_rule_conditions (
    fiscal_rule_id,
    field_name,
    operator,
    field_value
)
SELECT id, 'origin_uf', 'EQUALS', origin_uf
FROM fiscal_rules
ON CONFLICT DO NOTHING;

INSERT INTO fiscal_rule_conditions (
    fiscal_rule_id,
    field_name,
    operator,
    field_value
)
SELECT id, 'destination_uf', 'EQUALS', destination_uf
FROM fiscal_rules
ON CONFLICT DO NOTHING;

INSERT INTO fiscal_rule_conditions (
    fiscal_rule_id,
    field_name,
    operator,
    field_value
)
SELECT id, 'operation_type', 'EQUALS', operation_type
FROM fiscal_rules
ON CONFLICT DO NOTHING;

INSERT INTO fiscal_rule_conditions (
    fiscal_rule_id,
    field_name,
    operator,
    field_value
)
SELECT id, 'customer_type', 'EQUALS', customer_type
FROM fiscal_rules
ON CONFLICT DO NOTHING;

INSERT INTO fiscal_rule_taxes (
    fiscal_rule_id,
    tax_name,
    rate,
    calculation_order
)
SELECT id, 'ICMS', icms_rate, 1
FROM fiscal_rules
ON CONFLICT DO NOTHING;

INSERT INTO fiscal_rule_taxes (
    fiscal_rule_id,
    tax_name,
    rate,
    calculation_order
)
SELECT id, 'IBS', ibs_rate, 2
FROM fiscal_rules
ON CONFLICT DO NOTHING;

INSERT INTO fiscal_rule_taxes (
    fiscal_rule_id,
    tax_name,
    rate,
    calculation_order
)
SELECT id, 'CBS', cbs_rate, 3
FROM fiscal_rules
ON CONFLICT DO NOTHING;
