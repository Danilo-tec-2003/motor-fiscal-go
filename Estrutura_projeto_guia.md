# Guia Da Estrutura Do Projeto

Este documento explica a estrutura do projeto `motor-fiscal-go` de forma didatica. Ele serve como guia de estudo, consulta durante o desenvolvimento e apoio para apresentacao tecnica do projeto.

## Visao Geral

O Motor Fiscal e um microservico em Go para centralizar calculos e decisoes fiscais usadas pelo sistema de fretes.

Fluxo principal:

```text
Sistema de Fretes Java/JSP
        |
        | REST/JSON
        v
Motor Fiscal Go
        |
        | regras, calculo e auditoria
        v
PostgreSQL
```

O sistema Java continua responsavel pela operacao de frete. O Motor Fiscal fica responsavel pela regra fiscal, calculo, validacao, rastreabilidade e memoria de calculo.

## Estrutura Atual

```text
motor-fiscal-go/
├── cmd/api
├── internal/config
├── internal/dto
├── internal/errors
├── internal/handler
├── internal/middleware
├── internal/model
├── internal/repository
├── internal/service
├── internal/validator
├── docs
├── migrations
├── Dockerfile
├── docker-compose.yml
├── README.md
├── go.mod
└── go.sum
```

## Responsabilidade Das Pastas

| Pasta | Responsabilidade | Equivalente Em Java/Spring |
|---|---|---|
| `cmd/api` | Ponto de entrada da aplicacao. Inicializa config, banco, services, handlers e servidor HTTP. | Classe `Application` com `main` |
| `internal/config` | Leitura de variaveis de ambiente e configuracoes da aplicacao. | `application.properties`, `@Configuration` |
| `internal/dto` | Estruturas JSON de entrada e saida da API. | DTOs, request/response classes |
| `internal/errors` | Modelo padrao de erro da API. | Exceptions customizadas, error response |
| `internal/handler` | Camada HTTP. Recebe requests, valida metodo, decodifica JSON e chama services. | Controllers |
| `internal/middleware` | Interceptadores HTTP, como API Key e correlation ID. | Filters, Interceptors |
| `internal/model` | Modelos internos do dominio fiscal. | Domain model, entities de dominio |
| `internal/repository` | Acesso ao PostgreSQL. Executa queries e persiste historico. | Repositories, DAOs |
| `internal/service` | Regra de negocio: motor de regras, calculo fiscal, validacao CT-e. | Services |
| `internal/validator` | Validacoes de payload antes da regra de negocio. | Bean Validation, validators |
| `docs` | OpenAPI e contrato explicativo da API. | Swagger/OpenAPI docs |
| `migrations` | Scripts SQL para criar e evoluir banco. | Flyway/Liquibase migrations |

## Fluxo Entre Camadas

```text
handler
  -> validator
  -> service
      -> repository
      -> calculator/motor de regras
  -> response JSON
```

Exemplo real no endpoint `/api/v1/tax/simulate`:

```text
TaxHandler.Simulate
  -> decodifica TaxSimulationRequest
  -> validator.ValidateTaxSimulationRequest
  -> TaxService.Simulate
      -> FiscalRuleService.FindRule
          -> FiscalRuleRepository.FindCandidateRules
          -> avalia condicoes da regra
          -> detecta conflito
      -> CalculateTaxes
      -> FiscalSimulationRepository.Save
  -> retorna TaxSimulationResponse
```

## Arquivos Principais

### `cmd/api/main.go`

Responsavel por:

- carregar configuracao;
- conectar no PostgreSQL;
- criar repositories;
- criar services;
- criar handlers;
- criar router;
- iniciar servidor HTTP.

Este arquivo nao deve conter regra fiscal. Ele apenas monta a aplicacao.

### `internal/config/config.go`

Responsavel por carregar:

- `APP_PORT`;
- `INTERNAL_API_KEY`;
- `DATABASE_URL`.

Se a variavel nao existir, usa valor padrao para desenvolvimento.

### `internal/handler/routes.go`

Registra as rotas:

```text
GET  /health
POST /api/v1/tax/simulate
POST /api/v1/tax/compare
POST /api/v1/tax/batch
POST /api/v1/cte/validate
```

Tambem aplica middlewares:

```text
CorrelationID
APIKey
```

### `internal/handler/tax_handler.go`

Camada HTTP dos endpoints fiscais.

Responsavel por:

- validar metodo HTTP;
- decodificar JSON;
- chamar validator;
- chamar `TaxService`;
- converter erros de negocio em respostas HTTP padronizadas.

Nao deve conter calculo fiscal.

### `internal/service/tax_service.go`

Orquestra o fluxo fiscal:

```text
request
-> busca regra fiscal
-> calcula impostos
-> salva historico
-> retorna response
```

Ele nao deve conhecer detalhes do banco. Para isso existe repository.

### `internal/service/fiscal_rule_service.go`

E o motor de matching de regras fiscais.

Responsavel por:

- receber regras candidatas do repository;
- avaliar condicoes;
- ordenar por prioridade;
- detectar conflitos;
- transformar a regra escolhida em `model.FiscalRule`.

Operadores suportados:

```text
EQUALS
NOT_EQUALS
IN
BETWEEN
```

### `internal/service/tax_calculator.go`

Responsavel pelo calculo dos impostos.

Calcula:

- base de calculo;
- reducao de base, se existir;
- base efetiva;
- aliquota;
- valor do imposto;
- total de tributos;
- total com tributos;
- memoria de calculo.

Exemplo de formula:

```text
effective_base_value * rate / 100
```

### `internal/repository/fiscal_rule_repository.go`

Busca no banco as regras candidatas.

Carrega dados de:

```text
fiscal_rules
fiscal_rule_conditions
fiscal_rule_taxes
```

Nao decide a regra final. Ele apenas busca os dados. Quem decide e o service.

### `internal/repository/fiscal_simulation_repository.go`

Persiste a simulacao fiscal e a auditoria.

Salva dados em:

```text
fiscal_simulations
fiscal_simulation_tax_details
```

Usa transacao para evitar salvar simulacao sem detalhes.

### `internal/middleware/api_key.go`

Valida o header:

```text
X-API-Key
```

Se o token estiver ausente ou invalido, retorna erro JSON padronizado.

### `internal/middleware/correlation_id.go`

Garante rastreabilidade com:

```text
X-Correlation-ID
```

Se o Java enviar um ID, o Go reaproveita. Se nao enviar, o Go gera um.

### `internal/dto`

Contem os contratos JSON.

Arquivos principais:

```text
tax_simulation_request.go
tax_simulation_response.go
tax_batch_request.go
tax_batch_response.go
tax_comparison_response.go
cte_validation_request.go
cte_validation_response.go
```

DTO representa o que trafega pela API.

### `internal/model`

Contem conceitos internos do dominio.

Arquivos principais:

```text
fiscal_rule.go
fiscal_rule_engine.go
```

Model representa o dominio interno, nao necessariamente igual ao JSON.

## Banco De Dados

Principais tabelas:

| Tabela | Papel |
|---|---|
| `fiscal_rules` | Regra fiscal, vigencia, prioridade, status e CFOP |
| `fiscal_rule_conditions` | Condicoes para aplicar uma regra |
| `fiscal_rule_taxes` | Impostos e aliquotas da regra |
| `fiscal_rule_sources` | Fontes legais ou referencias da regra |
| `fiscal_rule_accounting_reviews` | Validacao contabil da regra |
| `fiscal_simulations` | Historico resumido das simulacoes |
| `fiscal_simulation_tax_details` | Memoria de calculo por imposto |

## Migrations

```text
001_create_fiscal_tables.sql
```

Cria tabelas iniciais e regras demonstrativas.

```text
002_model_fiscal_rule_engine.sql
```

Evolui o banco para motor de regras com condicoes e impostos separados.

```text
003_add_fiscal_rule_source_and_review_tracking.sql
```

Adiciona fonte fiscal e validacao contabil.

```text
004_persist_fiscal_calculation_audit.sql
```

Adiciona persistencia da memoria de calculo.

## Exemplo De Fluxo Do Calculo

Entrada:

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

Motor busca regras candidatas:

```text
active = true
status IN ('APPROVED', 'PENDING_REVIEW')
operation_date entre valid_from e valid_to
```

Depois avalia condicoes:

```text
origin_uf = PE
destination_uf = SP
operation_type = INTERESTADUAL
customer_type = PJ
```

Depois calcula:

```text
ICMS = 3500.00 * 12.00 / 100 = 420.00
IBS  = 3500.00 * 3.60 / 100 = 126.00
CBS  = 3500.00 * 0.90 / 100 = 31.50
```

Total:

```text
total_tax = 577.50
total_with_tax = 4077.50
```

## Como Explicar Em Uma Apresentacao

Resumo direto:

> O Motor Fiscal e um microservico em Go que centraliza calculo fiscal de fretes. Ele recebe dados operacionais do sistema Java, seleciona uma regra fiscal versionada no banco, calcula os tributos, salva a memoria de calculo e devolve o resultado para preencher os campos fiscais da tela de frete.

Pontos tecnicos fortes:

- separacao de responsabilidades entre Java e Go;
- API REST com contrato OpenAPI;
- calculo decimal sem `float`;
- regras no banco, nao hardcoded;
- motor de matching por condicoes;
- persistencia de auditoria;
- testes automatizados do nucleo fiscal;
- Docker com PostgreSQL;
- erros padronizados;
- correlation ID.

## O Que Nao Esta No Escopo Atual

- emissao real de CT-e;
- validacao fiscal completa nacional;
- tela administrativa de regras;
- cache Redis;
- logs estruturados avancados;
- integracao final com o sistema Java.

Esses pontos ficam como evolucoes futuras.

