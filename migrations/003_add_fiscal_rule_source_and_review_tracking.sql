CREATE TABLE IF NOT EXISTS fiscal_rule_sources (
    id BIGSERIAL PRIMARY KEY,
    fiscal_rule_id BIGINT NOT NULL REFERENCES fiscal_rules(id) ON DELETE CASCADE,
    source_type VARCHAR(40) NOT NULL,
    source_status VARCHAR(30) NOT NULL DEFAULT 'PENDING_CONFIRMATION',
    title VARCHAR(255) NOT NULL,
    reference VARCHAR(120) NOT NULL,
    url TEXT,
    published_at DATE,
    notes TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_fiscal_rule_sources_reference
ON fiscal_rule_sources (
    fiscal_rule_id,
    source_type,
    reference
);

CREATE INDEX IF NOT EXISTS idx_fiscal_rule_sources_rule_id
ON fiscal_rule_sources (fiscal_rule_id);

CREATE INDEX IF NOT EXISTS idx_fiscal_rule_sources_status
ON fiscal_rule_sources (source_status);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'ck_fiscal_rule_sources_source_type'
    ) THEN
        ALTER TABLE fiscal_rule_sources
            ADD CONSTRAINT ck_fiscal_rule_sources_source_type
            CHECK (
                source_type IN (
                    'FEDERAL_LAW',
                    'STATE_LAW',
                    'CONFAZ',
                    'SEFAZ',
                    'SENATE_RESOLUTION',
                    'ACCOUNTING_GUIDANCE',
                    'INTERNAL_NOTE'
                )
            );
    END IF;

    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'ck_fiscal_rule_sources_source_status'
    ) THEN
        ALTER TABLE fiscal_rule_sources
            ADD CONSTRAINT ck_fiscal_rule_sources_source_status
            CHECK (source_status IN ('PENDING_CONFIRMATION', 'CONFIRMED', 'REPLACED'));
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS fiscal_rule_accounting_reviews (
    id BIGSERIAL PRIMARY KEY,
    fiscal_rule_id BIGINT NOT NULL REFERENCES fiscal_rules(id) ON DELETE CASCADE,
    review_status VARCHAR(30) NOT NULL DEFAULT 'PENDING_REVIEW',
    reviewer_name VARCHAR(120),
    reviewer_role VARCHAR(80),
    reviewed_at TIMESTAMP,
    notes TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_fiscal_rule_accounting_reviews_rule_id
ON fiscal_rule_accounting_reviews (fiscal_rule_id);

CREATE INDEX IF NOT EXISTS idx_fiscal_rule_accounting_reviews_status
ON fiscal_rule_accounting_reviews (review_status);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'ck_fiscal_rule_accounting_reviews_status'
    ) THEN
        ALTER TABLE fiscal_rule_accounting_reviews
            ADD CONSTRAINT ck_fiscal_rule_accounting_reviews_status
            CHECK (review_status IN ('PENDING_REVIEW', 'APPROVED', 'REJECTED', 'CHANGES_REQUESTED'));
    END IF;
END $$;

UPDATE fiscal_rules
SET
    status = 'PENDING_REVIEW',
    description = CASE
        WHEN description = '' THEN 'Regra demonstrativa para desenvolvimento, pendente de fonte oficial e validacao contabil.'
        ELSE description
    END
WHERE NOT EXISTS (
    SELECT 1
    FROM fiscal_rule_accounting_reviews ar
    WHERE ar.fiscal_rule_id = fiscal_rules.id
      AND ar.review_status = 'APPROVED'
);

INSERT INTO fiscal_rule_sources (
    fiscal_rule_id,
    source_type,
    source_status,
    title,
    reference,
    notes
)
SELECT
    id,
    'INTERNAL_NOTE',
    'PENDING_CONFIRMATION',
    'Fonte oficial pendente de confirmacao',
    'A_CONFIRMAR',
    'Registro criado para deixar explicito que a regra ainda precisa de fonte normativa oficial antes de uso produtivo.'
FROM fiscal_rules
ON CONFLICT DO NOTHING;

INSERT INTO fiscal_rule_accounting_reviews (
    fiscal_rule_id,
    review_status,
    notes
)
SELECT
    fr.id,
    'PENDING_REVIEW',
    'Regra aguardando validacao contabil antes de uso produtivo.'
FROM fiscal_rules fr
WHERE NOT EXISTS (
    SELECT 1
    FROM fiscal_rule_accounting_reviews ar
    WHERE ar.fiscal_rule_id = fr.id
      AND ar.review_status = 'PENDING_REVIEW'
);
