CREATE TABLE IF NOT EXISTS fiscal_rules (
    id BIGSERIAL PRIMARY KEY,
    rule_version VARCHAR(20) NOT NULL,
    origin_uf CHAR(2) NOT NULL,
    destination_uf CHAR(2) NOT NULL,
    operation_type VARCHAR(30) NOT NULL,
    customer_type VARCHAR(2) NOT NULL,
    icms_rate NUMERIC(5,2) NOT NULL,
    ibs_rate NUMERIC(5,2) NOT NULL,
    cbs_rate NUMERIC(5,2) NOT NULL,
    cfop VARCHAR(10) NOT NULL,
    valid_from DATE NOT NULL,
    valid_to DATE NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS fiscal_simulations (
    id BIGSERIAL PRIMARY KEY,
    freight_id BIGINT NOT NULL,
    operation_date DATE NOT NULL,
    origin_uf CHAR(2) NOT NULL,
    destination_uf CHAR(2) NOT NULL,
    freight_value NUMERIC(15,2) NOT NULL,
    icms_rate NUMERIC(5,2) NOT NULL,
    icms_amount NUMERIC(15,2) NOT NULL,
    ibs_rate NUMERIC(5,2) NOT NULL,
    ibs_amount NUMERIC(15,2) NOT NULL,
    cbs_rate NUMERIC(5,2) NOT NULL,
    cbs_amount NUMERIC(15,2) NOT NULL,
    total_tax NUMERIC(15,2) NOT NULL,
    total_with_tax NUMERIC(15,2) NOT NULL,
    cfop VARCHAR(10) NOT NULL,
    rule_version VARCHAR(20) NOT NULL,
    from_cache BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_fiscal_rules_unique_rule
ON fiscal_rules (
    rule_version,
    origin_uf,
    destination_uf,
    operation_type,
    customer_type,
    valid_from,
    valid_to
);

INSERT INTO fiscal_rules (
    rule_version,
    origin_uf,
    destination_uf,
    operation_type,
    customer_type,
    icms_rate,
    ibs_rate,
    cbs_rate,
    cfop,
    valid_from,
    valid_to,
    active
) VALUES
(
    '2026.01',
    'PE',
    'SP',
    'INTERESTADUAL',
    'PJ',
    12.00,
    3.60,
    0.90,
    '6351',
    '2026-01-01',
    '2026-12-31',
    TRUE
),
(
    '2026.01',
    'PE',
    'PE',
    'INTERNA',
    'PF',
    18.00,
    3.60,
    0.90,
    '5351',
    '2026-01-01',
    '2026-12-31',
    TRUE
)
ON CONFLICT DO NOTHING;
