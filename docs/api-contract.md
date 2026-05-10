# Contrato Da API - Motor Fiscal

Este documento explica, em linguagem mais direta, o contrato REST inicial do
Motor Fiscal. O contrato tecnico principal fica em:

```text
docs/openapi.yaml
```

O projeto usa a abordagem **contract-first**: primeiro definimos os endpoints,
requests, responses e erros; depois implementamos os handlers em Go seguindo
esse contrato.

## Objetivo Da API

O Motor Fiscal sera um microservico em Go chamado pelo sistema principal de
fretes em Java/JSP/PostgreSQL.

Fluxo esperado:

```text
Sistema de Fretes Java
-> Motor Fiscal Go
-> regras fiscais/cache/banco/API externa
-> Motor Fiscal Go
-> Sistema de Fretes Java
```

Responsabilidade do sistema Java:

- manter cadastros, telas, fretes, status, ocorrencias e operacao principal;
- chamar o Motor Fiscal quando precisar de calculo ou validacao fiscal;
- salvar ou apresentar o resultado fiscal recebido.

Responsabilidade do Motor Fiscal:

- calcular ICMS, IBS e CBS;
- determinar CFOP basico pela regra fiscal;
- validar dados minimos para CT-e;
- comparar modelo atual com cenario da reforma tributaria;
- processar calculos em lote;
- usar Redis para cache;
- persistir historico e regras no PostgreSQL.

## Arquivos De Documentacao

```text
docs/openapi.yaml
```

Contrato tecnico em OpenAPI 3.0.3. Pode ser aberto no Swagger Editor ou Swagger
UI.

```text
docs/api-contract.md
```

Documento explicativo para consulta durante o desenvolvimento.

## Autenticacao Inicial

O contrato usa autenticacao simples por header:

```text
X-API-Key: valor-do-token
```

O token real deve vir de variavel de ambiente. Nunca commitar token real.

Variavel sugerida:

```text
INTERNAL_API_KEY
```

O endpoint `GET /health` nao exige token no contrato inicial, porque ele serve
para checagem simples do servico.

## Correlation ID

Todas as chamadas podem enviar:

```text
X-Correlation-ID: req-20260504-0001
```

Esse identificador serve para rastrear a mesma requisicao nos logs do Java, do
Go e de possiveis APIs externas.

Regra recomendada:

- se o Java enviar `X-Correlation-ID`, o Go reaproveita;
- se o Java nao enviar, o Go gera um novo;
- toda resposta deve devolver o mesmo correlation ID.

## Endpoints Do MVP

| Metodo | Rota | Responsabilidade |
|---|---|---|
| GET | `/health` | Verificar se o servico esta online |
| POST | `/api/v1/tax/simulate` | Simular calculo fiscal de um frete |
| POST | `/api/v1/tax/compare` | Comparar modelo atual com reforma tributaria |
| POST | `/api/v1/tax/batch` | Calcular varios fretes em lote |
| POST | `/api/v1/cte/validate` | Validar dados minimos para CT-e |

## Endpoints Futuros

| Metodo | Rota | Responsabilidade |
|---|---|---|
| POST | `/api/v1/cte/preview` | Gerar preview/mock de CT-e sem emissao real |
| GET | `/api/v1/rules` | Listar regras fiscais |
| POST | `/api/v1/rules` | Cadastrar regra fiscal |
| GET | `/api/v1/reports/tax-summary` | Resumo fiscal por periodo |
| GET | `/api/v1/reports/reform-impact` | Impacto estimado da reforma |
| GET | `/api/v1/reports/inconsistencies` | Inconsistencias fiscais registradas |

## Endpoint: GET /health

Responsabilidade:

- confirmar que o Motor Fiscal esta online;
- retornar nome e versao do servico;
- nao depender de banco, Redis ou API externa no MVP.

Response `200`:

```json
{
  "status": "UP",
  "service": "motor-fiscal",
  "version": "1.0.0"
}
```

## Endpoint: POST /api/v1/tax/simulate

Responsabilidade:

- receber dados de um frete;
- validar campos obrigatorios;
- selecionar regra fiscal vigente pelo motor de regras;
- calcular ICMS, IBS e CBS;
- determinar CFOP;
- retornar totais fiscais;
- retornar regra aplicada e detalhes auditaveis do calculo;
- informar se o resultado veio do cache.

Request:

```json
{
  "freight_id": 1001,
  "operation_date": "2026-03-10",
  "origin_uf": "PE",
  "destination_uf": "SP",
  "freight_value": "3500.00",
  "customer_type": "PJ",
  "operation_type": "INTERESTADUAL"
}
```

Campos obrigatorios:

| Campo | Tipo | Observacao |
|---|---|---|
| `freight_id` | integer | ID do frete no sistema Java |
| `operation_date` | string/date | Data usada para regra fiscal vigente |
| `origin_uf` | string | UF de origem |
| `destination_uf` | string | UF de destino |
| `freight_value` | string | Valor monetario com duas casas decimais |
| `customer_type` | string | `PF` ou `PJ` |
| `operation_type` | string | `INTERNA` ou `INTERESTADUAL` |

Response `200`:

```json
{
  "freight_id": 1001,
  "base_value": "3500.00",
  "icms": {
    "rate": "12.00",
    "amount": "420.00"
  },
  "ibs": {
    "rate": "3.60",
    "amount": "126.00"
  },
  "cbs": {
    "rate": "0.90",
    "amount": "31.50"
  },
  "total_tax": "577.50",
  "total_with_tax": "4077.50",
  "cfop": "6351",
  "rule_id": 1,
  "rule_code": "RULE_PE_SP_INTERESTADUAL_PJ_2026_01_2026_01_01",
  "rule_version": "2026.01",
  "rule_status": "PENDING_REVIEW",
  "calculation_basis": "FREIGHT_VALUE",
  "calculation_details": [
    {
      "tax_name": "ICMS",
      "base_value": "3500.00",
      "base_reduction_rate": "0.00",
      "effective_base_value": "3500.00",
      "rate": "12.00",
      "amount": "420.00",
      "formula": "effective_base_value * rate / 100"
    },
    {
      "tax_name": "IBS",
      "base_value": "3500.00",
      "base_reduction_rate": "0.00",
      "effective_base_value": "3500.00",
      "rate": "3.60",
      "amount": "126.00",
      "formula": "effective_base_value * rate / 100"
    },
    {
      "tax_name": "CBS",
      "base_value": "3500.00",
      "base_reduction_rate": "0.00",
      "effective_base_value": "3500.00",
      "rate": "0.90",
      "amount": "31.50",
      "formula": "effective_base_value * rate / 100"
    }
  ],
  "from_cache": false
}
```

Observacao importante:

Os valores acima sao exemplos de desenvolvimento. Aliquotas, CFOPs e regras
reais precisam estar cadastrados no motor de regras, possuir fonte fiscal
confirmada e passar por validacao contabil antes de uso produtivo.

Campos de auditoria:

| Campo | Significado |
|---|---|
| `rule_id` | ID interno da regra aplicada |
| `rule_code` | Codigo estavel da regra fiscal selecionada |
| `rule_status` | Status da regra: `DRAFT`, `PENDING_REVIEW`, `APPROVED` ou `INACTIVE` |
| `calculation_basis` | Base usada no calculo. Hoje: `FREIGHT_VALUE` |
| `calculation_details` | Lista com a memoria de calculo de cada imposto |

## Endpoint: POST /api/v1/tax/compare

Responsabilidade:

- receber os mesmos dados de uma simulacao fiscal;
- calcular ou estimar o cenario atual;
- calcular ou estimar o cenario da reforma com IBS/CBS;
- retornar diferenca e analise resumida.

Request:

Usa o mesmo formato de `/api/v1/tax/simulate`.

Response conceitual:

```json
{
  "freight_id": 1001,
  "current_scenario": {
    "base_value": "3500.00",
    "total_tax": "420.00",
    "total_with_tax": "3920.00"
  },
  "reform_scenario": {
    "base_value": "3500.00",
    "total_tax": "577.50",
    "total_with_tax": "4077.50"
  },
  "difference": "157.50",
  "analysis": "Cenario da reforma apresentou aumento estimado de tributos."
}
```

## Endpoint: POST /api/v1/tax/batch

Responsabilidade:

- receber uma lista de fretes;
- processar calculos fiscais em lote;
- retornar sucesso ou erro por item.

Request:

```json
{
  "items": [
    {
      "freight_id": 1001,
      "operation_date": "2026-03-10",
      "origin_uf": "PE",
      "destination_uf": "SP",
      "freight_value": "3500.00",
      "customer_type": "PJ",
      "operation_type": "INTERESTADUAL"
    }
  ]
}
```

Response conceitual:

```json
{
  "total_items": 1,
  "success_count": 1,
  "error_count": 0,
  "results": [
    {
      "freight_id": 1001,
      "success": true,
      "data": {
        "freight_id": 1001,
        "base_value": "3500.00",
        "icms": {
          "rate": "12.00",
          "amount": "420.00"
        },
        "ibs": {
          "rate": "3.60",
          "amount": "126.00"
        },
        "cbs": {
          "rate": "0.90",
          "amount": "31.50"
        },
        "total_tax": "577.50",
        "total_with_tax": "4077.50",
        "cfop": "6351",
        "rule_id": 1,
        "rule_code": "RULE_PE_SP_INTERESTADUAL_PJ_2026_01_2026_01_01",
        "rule_version": "2026.01",
        "rule_status": "PENDING_REVIEW",
        "calculation_basis": "FREIGHT_VALUE",
        "calculation_details": [
          {
            "tax_name": "ICMS",
            "base_value": "3500.00",
            "base_reduction_rate": "0.00",
            "effective_base_value": "3500.00",
            "rate": "12.00",
            "amount": "420.00",
            "formula": "effective_base_value * rate / 100"
          },
          {
            "tax_name": "IBS",
            "base_value": "3500.00",
            "base_reduction_rate": "0.00",
            "effective_base_value": "3500.00",
            "rate": "3.60",
            "amount": "126.00",
            "formula": "effective_base_value * rate / 100"
          },
          {
            "tax_name": "CBS",
            "base_value": "3500.00",
            "base_reduction_rate": "0.00",
            "effective_base_value": "3500.00",
            "rate": "0.90",
            "amount": "31.50",
            "formula": "effective_base_value * rate / 100"
          }
        ],
        "from_cache": false
      }
    }
  ]
}
```

Regra recomendada:

- se alguns itens falharem e outros passarem, retornar `200` com erro por item;
- se o JSON inteiro estiver invalido, retornar `400`;
- se a lista estiver vazia, retornar `422`.

## Endpoint: POST /api/v1/cte/validate

Responsabilidade:

- validar dados minimos para CT-e;
- verificar valor do frete;
- verificar UF origem/destino;
- verificar remetente e destinatario;
- verificar se ha regra fiscal e CFOP;
- retornar erros e avisos.

Request:

```json
{
  "freight_id": 1001,
  "operation_date": "2026-03-10",
  "origin_uf": "PE",
  "destination_uf": "SP",
  "freight_value": "3500.00",
  "customer_type": "PJ",
  "operation_type": "INTERESTADUAL",
  "sender": {
    "name": "Remetente Exemplo",
    "document": "12345678000199",
    "uf": "PE"
  },
  "recipient": {
    "name": "Destinatario Exemplo",
    "document": "98765432000188",
    "uf": "SP"
  }
}
```

Response `200` quando valido:

```json
{
  "freight_id": 1001,
  "valid": true,
  "cfop": "6351",
  "errors": [],
  "warnings": []
}
```

Response `200` quando invalido, mas a validacao foi executada:

```json
{
  "freight_id": 1001,
  "valid": false,
  "errors": [
    {
      "field": "recipient.document",
      "message": "Documento do destinatario e obrigatorio."
    }
  ],
  "warnings": []
}
```

Use `422` quando o payload estiver incompleto a ponto de impedir a validacao.

## Padrao De Erro

Todos os erros tecnicos ou de validacao devem seguir este formato:

```json
{
  "code": "VALIDATION_ERROR",
  "message": "Payload invalido.",
  "correlation_id": "req-20260504-0001",
  "details": [
    {
      "field": "freight_value",
      "message": "Deve ser maior que zero."
    }
  ]
}
```

Codigos de erro iniciais:

| Codigo | Quando usar |
|---|---|
| `BAD_REQUEST` | JSON invalido ou malformado |
| `UNAUTHORIZED` | `X-API-Key` ausente ou invalido |
| `VALIDATION_ERROR` | Campos invalidos ou inconsistentes |
| `FISCAL_RULE_NOT_FOUND` | Nenhuma regra fiscal vigente encontrada |
| `FISCAL_RULE_CONFLICT` | Mais de uma regra fiscal compativel foi encontrada |
| `FISCAL_RULE_INCOMPLETE` | Regra fiscal sem impostos obrigatorios para calculo |
| `UNSUPPORTED_CALCULATION_BASIS` | Base de calculo fiscal ainda nao suportada |
| `EXTERNAL_TIMEOUT` | Timeout em API externa |
| `EXTERNAL_UNAVAILABLE` | API externa indisponivel |
| `INTERNAL_ERROR` | Erro inesperado no microservico |
| `NOT_IMPLEMENTED` | Endpoint documentado para versao futura |

Status HTTP recomendados:

| Status | Uso |
|---|---|
| `200` | Operacao realizada com sucesso |
| `201` | Recurso criado, como regra fiscal futura |
| `400` | JSON invalido |
| `401` | Token ausente ou invalido |
| `404` | Regra ou recurso nao encontrado |
| `409` | Conflito entre regras fiscais compativeis |
| `422` | Dados em JSON validos, mas fiscalmente invalidos |
| `500` | Erro interno inesperado |
| `501` | Endpoint previsto, mas ainda nao implementado |
| `503` | API externa indisponivel |
| `504` | Timeout externo |

## Decisao Sobre Valores Monetarios

Campos monetarios devem trafegar como string:

```json
{
  "freight_value": "3500.00"
}
```

Motivo:

- evitar perda de precisao em JSON;
- evitar uso de `float` em calculos financeiros;
- facilitar uso de decimal no Go;
- manter padrao previsivel entre Java e Go.

No Go, o calculo deve usar biblioteca decimal, como `shopspring/decimal`, ou
equivalente.

## Regras Fiscais

O calculo fiscal depende de uma regra vigente selecionada pelo motor de regras.
O motor deve:

1. buscar regras candidatas pela vigencia e status;
2. avaliar as condicoes cadastradas em `fiscal_rule_conditions`;
3. ordenar regras compativeis por prioridade;
4. detectar conflito quando regras equivalentes tiverem mesma prioridade;
5. usar os impostos cadastrados em `fiscal_rule_taxes`;
6. retornar uma memoria de calculo auditavel.

A regra deve considerar pelo menos:

- UF origem;
- UF destino;
- tipo de operacao;
- tipo de cliente;
- data da operacao;
- aliquotas;
- CFOP;
- versao da regra;
- vigencia inicial e final;
- status ativo.

Tabelas principais do motor:

```text
fiscal_rules
fiscal_rule_conditions
fiscal_rule_taxes
fiscal_rule_sources
fiscal_rule_accounting_reviews
fiscal_simulations
fiscal_simulation_tax_details
```

Campos principais de `fiscal_rules`:

```text
id
rule_code
rule_version
description
priority
status
calculation_basis
cfop
valid_from
valid_to
active
created_at
updated_at
```

Exemplo conceitual de condicoes:

```text
origin_uf = PE
destination_uf = SP
operation_type = INTERESTADUAL
customer_type = PJ
```

Exemplo conceitual de impostos da regra:

```text
ICMS -> 12.00
IBS  -> 3.60
CBS  -> 0.90
```

Status de regra:

| Status | Uso |
|---|---|
| `DRAFT` | Regra em rascunho |
| `PENDING_REVIEW` | Regra pendente de validacao contabil |
| `APPROVED` | Regra validada para uso |
| `INACTIVE` | Regra desativada |

Auditoria persistida:

Quando `/api/v1/tax/simulate` roda com sucesso, o Motor Fiscal salva:

- resumo em `fiscal_simulations`;
- detalhe por imposto em `fiscal_simulation_tax_details`;
- regra aplicada (`rule_id`, `rule_code`, `rule_status`);
- base, aliquota, valor calculado e formula.

## Cache

O Redis deve ser usado para:

- cache de resultado de simulacao fiscal;
- cache de regra fiscal vigente.

Chaves conceituais:

```text
tax:simulation:{hash_da_entrada}
tax:rule:{origin_uf}:{destination_uf}:{operation_date}:{customer_type}:{operation_type}
```

O response de simulacao deve retornar:

```json
{
  "from_cache": false
}
```

## APIs Externas

Os anexos citam possibilidades futuras como ViaCEP, rota, TecnoSpeed e Mercado
Pago. O contrato oficial dessas APIs externas nao esta documentado nos anexos.

Portanto:

- payloads reais de APIs externas: a confirmar na documentacao oficial;
- headers reais de APIs externas: a confirmar na documentacao oficial;
- tokens reais: sempre por variavel de ambiente;
- integracao real com emissao de CT-e: fora do MVP.

## Ordem De Implementacao Recomendada

1. Validar `docs/openapi.yaml` no Swagger Editor.
2. Implementar `GET /health`.
3. Implementar erro padrao.
4. Implementar DTO de `/api/v1/tax/simulate`.
5. Implementar validator.
6. Implementar service de calculo com valores fixos temporarios.
7. Trocar valores fixos por regras fiscais versionadas.
8. Adicionar Redis.
9. Adicionar historico PostgreSQL.
10. Implementar `/api/v1/cte/validate`.
11. Implementar `/api/v1/tax/compare`.
12. Implementar `/api/v1/tax/batch`.

## Como Validar A Documentacao

Opcao simples:

1. Acesse `https://editor.swagger.io/`.
2. Cole o conteudo de `docs/openapi.yaml`.
3. Corrija qualquer erro de indentacao ou schema.

Opcao local futura com Docker:

```bash
docker run -p 8081:8080 -e SWAGGER_JSON=/docs/openapi.yaml -v ./docs:/docs swaggerapi/swagger-ui
```

Depois acessar:

```text
http://localhost:8081
```

## Checklist Antes De Codar

- [ ] Sei qual endpoint o Java vai chamar primeiro.
- [ ] Sei quais campos sao obrigatorios.
- [ ] Sei qual JSON o Go deve responder.
- [ ] Sei qual formato de erro sera usado.
- [ ] Sei que valores monetarios serao strings.
- [ ] Sei que token real nao entra no Git.
- [ ] Sei que `X-Correlation-ID` deve aparecer nos logs.
- [ ] Sei quais endpoints sao MVP e quais sao futuros.
