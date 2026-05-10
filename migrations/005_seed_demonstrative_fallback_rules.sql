-- ============================================================
-- 005_seed_demonstrative_fallback_rules.sql
-- Cobertura operacional demonstrativa para o MVP.
--
-- Importante:
-- Estas regras NAO representam uma base fiscal oficial.
-- Elas existem para permitir demonstracao nacional do fluxo:
-- preview -> emissao -> calculo oficial -> auditoria.
--
-- Em producao, regras fallback devem ser substituidas por regras
-- especificas validadas por fonte legal/contabil ou por provider fiscal.
-- ============================================================

INSERT INTO fiscal_rules (
    rule_code,
    rule_version,
    description,
    priority,
    status,
    calculation_basis,
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
    'FALLBACK_INTERESTADUAL_PJ_2026_01',
    '2026.01',
    'Regra fallback demonstrativa para frete interestadual PJ. Usada apenas quando nao existir regra especifica por UF.',
    900,
    'PENDING_REVIEW',
    'FREIGHT_VALUE',
    'BR',
    'BR',
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
    'FALLBACK_INTERESTADUAL_PF_2026_01',
    '2026.01',
    'Regra fallback demonstrativa para frete interestadual PF. Usada apenas quando nao existir regra especifica por UF.',
    900,
    'PENDING_REVIEW',
    'FREIGHT_VALUE',
    'BR',
    'BR',
    'INTERESTADUAL',
    'PF',
    12.00,
    3.60,
    0.90,
    '6351',
    '2026-01-01',
    '2026-12-31',
    TRUE
),
(
    'FALLBACK_INTERNA_PJ_2026_01',
    '2026.01',
    'Regra fallback demonstrativa para frete interno PJ. Usada apenas quando nao existir regra especifica por UF.',
    900,
    'PENDING_REVIEW',
    'FREIGHT_VALUE',
    'BR',
    'BR',
    'INTERNA',
    'PJ',
    18.00,
    3.60,
    0.90,
    '5351',
    '2026-01-01',
    '2026-12-31',
    TRUE
),
(
    'FALLBACK_INTERNA_PF_2026_01',
    '2026.01',
    'Regra fallback demonstrativa para frete interno PF. Usada apenas quando nao existir regra especifica por UF.',
    900,
    'PENDING_REVIEW',
    'FREIGHT_VALUE',
    'BR',
    'BR',
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
ON CONFLICT (rule_code) DO UPDATE SET
    description = EXCLUDED.description,
    priority = EXCLUDED.priority,
    status = EXCLUDED.status,
    calculation_basis = EXCLUDED.calculation_basis,
    icms_rate = EXCLUDED.icms_rate,
    ibs_rate = EXCLUDED.ibs_rate,
    cbs_rate = EXCLUDED.cbs_rate,
    cfop = EXCLUDED.cfop,
    valid_from = EXCLUDED.valid_from,
    valid_to = EXCLUDED.valid_to,
    active = EXCLUDED.active,
    updated_at = NOW();

INSERT INTO fiscal_rule_conditions (
    fiscal_rule_id,
    field_name,
    operator,
    field_value
)
SELECT
    fr.id,
    condition.field_name,
    'EQUALS',
    condition.field_value
FROM fiscal_rules fr
JOIN (
    VALUES
        ('FALLBACK_INTERESTADUAL_PJ_2026_01', 'operation_type', 'INTERESTADUAL'),
        ('FALLBACK_INTERESTADUAL_PJ_2026_01', 'customer_type', 'PJ'),
        ('FALLBACK_INTERESTADUAL_PF_2026_01', 'operation_type', 'INTERESTADUAL'),
        ('FALLBACK_INTERESTADUAL_PF_2026_01', 'customer_type', 'PF'),
        ('FALLBACK_INTERNA_PJ_2026_01', 'operation_type', 'INTERNA'),
        ('FALLBACK_INTERNA_PJ_2026_01', 'customer_type', 'PJ'),
        ('FALLBACK_INTERNA_PF_2026_01', 'operation_type', 'INTERNA'),
        ('FALLBACK_INTERNA_PF_2026_01', 'customer_type', 'PF')
) AS condition(rule_code, field_name, field_value)
    ON condition.rule_code = fr.rule_code
ON CONFLICT DO NOTHING;

INSERT INTO fiscal_rule_taxes (
    fiscal_rule_id,
    tax_name,
    rate,
    base_reduction_rate,
    calculation_order
)
SELECT
    fr.id,
    tax.tax_name,
    tax.rate,
    0,
    tax.calculation_order
FROM fiscal_rules fr
JOIN (
    VALUES
        ('FALLBACK_INTERESTADUAL_PJ_2026_01', 'ICMS', 12.00, 1),
        ('FALLBACK_INTERESTADUAL_PJ_2026_01', 'IBS', 3.60, 2),
        ('FALLBACK_INTERESTADUAL_PJ_2026_01', 'CBS', 0.90, 3),
        ('FALLBACK_INTERESTADUAL_PF_2026_01', 'ICMS', 12.00, 1),
        ('FALLBACK_INTERESTADUAL_PF_2026_01', 'IBS', 3.60, 2),
        ('FALLBACK_INTERESTADUAL_PF_2026_01', 'CBS', 0.90, 3),
        ('FALLBACK_INTERNA_PJ_2026_01', 'ICMS', 18.00, 1),
        ('FALLBACK_INTERNA_PJ_2026_01', 'IBS', 3.60, 2),
        ('FALLBACK_INTERNA_PJ_2026_01', 'CBS', 0.90, 3),
        ('FALLBACK_INTERNA_PF_2026_01', 'ICMS', 18.00, 1),
        ('FALLBACK_INTERNA_PF_2026_01', 'IBS', 3.60, 2),
        ('FALLBACK_INTERNA_PF_2026_01', 'CBS', 0.90, 3)
) AS tax(rule_code, tax_name, rate, calculation_order)
    ON tax.rule_code = fr.rule_code
ON CONFLICT (fiscal_rule_id, tax_name) DO UPDATE SET
    rate = EXCLUDED.rate,
    base_reduction_rate = EXCLUDED.base_reduction_rate,
    calculation_order = EXCLUDED.calculation_order;

INSERT INTO fiscal_rule_sources (
    fiscal_rule_id,
    source_type,
    source_status,
    title,
    reference,
    notes
)
SELECT
    fr.id,
    'INTERNAL_NOTE',
    'PENDING_CONFIRMATION',
    'Regra fallback demonstrativa para MVP',
    fr.rule_code,
    'Regra operacional demonstrativa. Nao utilizar como base fiscal oficial sem validacao contabil e fonte legal confirmada.'
FROM fiscal_rules fr
WHERE fr.rule_code IN (
    'FALLBACK_INTERESTADUAL_PJ_2026_01',
    'FALLBACK_INTERESTADUAL_PF_2026_01',
    'FALLBACK_INTERNA_PJ_2026_01',
    'FALLBACK_INTERNA_PF_2026_01'
)
ON CONFLICT DO NOTHING;

INSERT INTO fiscal_rule_accounting_reviews (
    fiscal_rule_id,
    review_status,
    notes
)
SELECT
    fr.id,
    'PENDING_REVIEW',
    'Fallback demonstrativo criado para cobertura operacional do MVP. Exige validacao fiscal antes de uso produtivo.'
FROM fiscal_rules fr
WHERE fr.rule_code IN (
    'FALLBACK_INTERESTADUAL_PJ_2026_01',
    'FALLBACK_INTERESTADUAL_PF_2026_01',
    'FALLBACK_INTERNA_PJ_2026_01',
    'FALLBACK_INTERNA_PF_2026_01'
)
  AND NOT EXISTS (
      SELECT 1
      FROM fiscal_rule_accounting_reviews ar
      WHERE ar.fiscal_rule_id = fr.id
        AND ar.review_status = 'PENDING_REVIEW'
  );
