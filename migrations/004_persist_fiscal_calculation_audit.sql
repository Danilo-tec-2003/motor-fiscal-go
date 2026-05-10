ALTER TABLE fiscal_simulations
    ADD COLUMN IF NOT EXISTS rule_id BIGINT REFERENCES fiscal_rules(id),
    ADD COLUMN IF NOT EXISTS rule_code VARCHAR(80),
    ADD COLUMN IF NOT EXISTS rule_status VARCHAR(30),
    ADD COLUMN IF NOT EXISTS calculation_basis VARCHAR(30);

CREATE INDEX IF NOT EXISTS idx_fiscal_simulations_rule_id
ON fiscal_simulations (rule_id);

CREATE INDEX IF NOT EXISTS idx_fiscal_simulations_rule_code
ON fiscal_simulations (rule_code);

CREATE TABLE IF NOT EXISTS fiscal_simulation_tax_details (
    id BIGSERIAL PRIMARY KEY,
    fiscal_simulation_id BIGINT NOT NULL REFERENCES fiscal_simulations(id) ON DELETE CASCADE,
    tax_name VARCHAR(20) NOT NULL,
    base_value NUMERIC(15,2) NOT NULL,
    base_reduction_rate NUMERIC(7,4) NOT NULL DEFAULT 0,
    effective_base_value NUMERIC(15,2) NOT NULL,
    rate NUMERIC(7,4) NOT NULL,
    amount NUMERIC(15,2) NOT NULL,
    formula TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_fiscal_simulation_tax_details_simulation_id
ON fiscal_simulation_tax_details (fiscal_simulation_id);

CREATE INDEX IF NOT EXISTS idx_fiscal_simulation_tax_details_tax_name
ON fiscal_simulation_tax_details (tax_name);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'ck_fiscal_simulation_tax_details_tax_name'
    ) THEN
        ALTER TABLE fiscal_simulation_tax_details
            ADD CONSTRAINT ck_fiscal_simulation_tax_details_tax_name
            CHECK (tax_name IN ('ICMS', 'IBS', 'CBS'));
    END IF;

    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'ck_fiscal_simulation_tax_details_base_reduction_rate'
    ) THEN
        ALTER TABLE fiscal_simulation_tax_details
            ADD CONSTRAINT ck_fiscal_simulation_tax_details_base_reduction_rate
            CHECK (base_reduction_rate >= 0 AND base_reduction_rate <= 100);
    END IF;
END $$;
