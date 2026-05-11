# Motor Fiscal

API fiscal em Go para calculo, simulacao, validacao e auditoria tributaria em
operacoes de frete.

O Motor Fiscal foi criado para integrar o sistema Java `sistema-fretes` com um
servico especializado em regras fiscais. A proposta e retirar calculos fiscais
espalhados do sistema principal e centralizar a decisao tributaria em uma API
menor, testavel, versionada e explicavel.

## Contexto Da Solucao

Este projeto faz parte de uma solucao formada por tres camadas:

| Projeto | Papel |
|---|---|
| `analise-prs` | Projeto Python usado para diagnosticar dores recorrentes em PRs internos |
| `sistema-fretes` | Sistema Java/JSP responsavel pela operacao de fretes |
| `Motor_fiscal` | API Go responsavel por regras fiscais, calculos e auditoria |

Antes da construcao do Motor Fiscal, foi feito um diagnostico usando Python
sobre PRs internos. A analise buscou identificar padroes de manutencao,
recorrencia de problemas e areas com maior custo tecnico.

Resultado local da analise:

| Indicador | Resultado |
|---|---:|
| PRs coletados | 13.397 |
| Repositorios analisados | 2 |
| Autores identificados | 55 |
| Periodo analisado | 2019-02-01 a 2026-04-15 |
| Registros com qualidade para analise aprofundada | 274 |

Categorias com sinais relevantes:

| Categoria | Quantidade |
|---|---:|
| Banco | 209 |
| Relatorio | 185 |
| Performance | 111 |
| Regra de negocio | 58 |

A leitura tecnica foi que regras de negocio e comportamento fiscal nao deveriam
ficar diluidos no sistema operacional. Por isso surgiu a proposta de criar uma
API fiscal dedicada.

## Problema Resolvido

Em sistemas de frete, dados como CFOP, aliquotas, bases de calculo, totais de
tributos e validacoes fiscais tendem a aparecer em varios pontos do codigo.
Isso cria quatro problemas:

- regra duplicada;
- baixa rastreabilidade;
- maior risco de divergencia entre telas e relatorios;
- dificuldade para explicar tecnicamente um calculo fiscal.

O Motor Fiscal centraliza:

- selecao de regra fiscal vigente;
- matching de condicoes;
- calculo de ICMS, IBS e CBS;
- memoria de calculo por imposto;
- historico de simulacoes;
- status e fonte da regra aplicada;
- contrato HTTP estavel para o sistema Java.

## Status

Projeto em fase de MVP tecnico com motor de regras proprio.

Ja implementado:

- API REST em Go;
- Docker com PostgreSQL;
- autenticacao por `X-API-Key`;
- `X-Correlation-ID` para rastreio;
- erros padronizados em JSON;
- motor de regras fiscais por condicoes;
- suporte a prioridade, vigencia e status da regra;
- calculo auditavel de ICMS, IBS e CBS;
- persistencia da simulacao;
- persistencia dos detalhes de calculo por imposto;
- validacao basica de CT-e;
- contrato OpenAPI;
- testes do motor de regras, calculo fiscal, handlers e middleware.

Aviso importante: as regras fiscais atuais sao demonstrativas e estao marcadas
como `PENDING_REVIEW`. Para uso real, as regras precisam de fonte oficial
confirmada e revisao contabil.

## Stack Tecnica

| Area | Tecnologia |
|---|---|
| Linguagem | Go |
| API | `net/http` |
| Banco | PostgreSQL |
| Driver | `pgx/v5` |
| Decimal monetario | `shopspring/decimal` |
| Container | Docker / Docker Compose |
| Contrato | OpenAPI YAML |
| Testes | `go test` |

## Arquitetura

```text
Motor_fiscal/
|-- cmd/api/
|-- internal/
|   |-- config/
|   |-- dto/
|   |-- errors/
|   |-- handler/
|   |-- middleware/
|   |-- model/
|   |-- repository/
|   |-- service/
|   `-- validator/
|-- migrations/
|-- docs/
|-- Dockerfile
|-- docker-compose.yml
|-- go.mod
`-- README.md
```

Responsabilidades:

| Pasta | Responsabilidade |
|---|---|
| `cmd/api` | Inicializacao da API, banco, services, handlers e servidor HTTP |
| `internal/config` | Leitura de variaveis de ambiente |
| `internal/dto` | Contratos de entrada e saida JSON |
| `internal/errors` | Erros padronizados da API |
| `internal/handler` | Rotas HTTP |
| `internal/middleware` | API Key e Correlation ID |
| `internal/model` | Modelos internos do dominio fiscal |
| `internal/repository` | Persistencia em PostgreSQL |
| `internal/service` | Motor de regras, calculo e validacao de negocio |
| `internal/validator` | Validacao dos payloads |
| `migrations` | Modelo fiscal no banco |
| `docs` | OpenAPI e contrato explicativo |

## Endpoints

| Metodo | Rota | Auth | Papel |
|---|---|---|---|
| `GET` | `/health` | Nao | Verifica disponibilidade da API |
| `POST` | `/api/v1/tax/preview` | Sim | Calcula previa fiscal antes de salvar o frete |
| `POST` | `/api/v1/tax/simulate` | Sim | Calcula e persiste simulacao fiscal |
| `POST` | `/api/v1/tax/compare` | Sim | Compara cenario atual com cenario de reforma |
| `POST` | `/api/v1/tax/batch` | Sim | Processa multiplas simulacoes |
| `POST` | `/api/v1/cte/validate` | Sim | Valida dados minimos para CT-e |

Endpoints futuros documentados no OpenAPI:

- `POST /api/v1/cte/preview`;
- `GET /api/v1/rules`;
- `POST /api/v1/rules`;
- `GET /api/v1/reports/tax-summary`;
- `GET /api/v1/reports/reform-impact`;
- `GET /api/v1/reports/inconsistencies`.

## Fluxo De Integracao Com O Sistema Java

```text
sistema-fretes Java/JSP
  -> FreteBO
  -> MotorFiscalClient
  -> POST /api/v1/tax/simulate
     Headers:
       X-API-Key
       X-Correlation-ID
  -> Motor Fiscal Go
  -> PostgreSQL fiscal
  -> JSON com impostos, CFOP, regra e memoria de calculo
  -> Java persiste resumo na tabela frete
```

O Java nao chama a API pelo navegador. A chamada e feita pelo backend para nao
expor o token interno.

Campos retornados ao Java:

- `cfop`;
- `icms.rate` e `icms.amount`;
- `ibs.rate` e `ibs.amount`;
- `cbs.rate` e `cbs.amount`;
- `total_tax`;
- `total_with_tax`;
- `rule_id`;
- `rule_code`;
- `rule_version`;
- `rule_status`;
- `calculation_basis`;
- `calculation_details`;
- `from_cache`.

## Motor De Regras

O motor nao usa aliquotas fixas no codigo. As regras ficam no banco e sao
buscadas conforme a operacao.

Ordem conceitual:

```text
1. receber payload do frete;
2. validar campos obrigatorios;
3. buscar regras candidatas por vigencia e status;
4. avaliar condicoes da regra;
5. ordenar por prioridade;
6. detectar conflitos;
7. validar se a regra possui impostos obrigatorios;
8. calcular ICMS, IBS e CBS com decimal;
9. salvar simulacao e detalhes;
10. retornar resposta JSON.
```

Condicoes suportadas:

| Operador | Uso |
|---|---|
| `EQUALS` | Campo deve ser igual ao valor configurado |
| `NOT_EQUALS` | Campo deve ser diferente |
| `IN` | Campo deve estar em uma lista |
| `BETWEEN` | Campo numerico/data dentro de faixa |

Campos usados em regras:

- `origin_uf`;
- `destination_uf`;
- `customer_type`;
- `operation_type`;
- `freight_value`;
- `operation_date`.

Status de regra:

| Status | Significado |
|---|---|
| `DRAFT` | Rascunho |
| `PENDING_REVIEW` | Pendente de validacao contabil |
| `APPROVED` | Validada para uso |
| `INACTIVE` | Desativada |

## Banco De Dados

Migrations:

```text
migrations/001_create_fiscal_tables.sql
migrations/002_model_fiscal_rule_engine.sql
migrations/003_add_fiscal_rule_source_and_review_tracking.sql
migrations/004_persist_fiscal_calculation_audit.sql
migrations/005_seed_demonstrative_fallback_rules.sql
```

Tabelas principais:

| Tabela | Papel |
|---|---|
| `fiscal_rules` | Regra fiscal, vigencia, prioridade, status, CFOP e base de calculo |
| `fiscal_rule_conditions` | Condicoes que determinam quando a regra se aplica |
| `fiscal_rule_taxes` | Impostos e aliquotas vinculados a regra |
| `fiscal_rule_sources` | Fonte legal ou fonte de referencia |
| `fiscal_rule_accounting_reviews` | Revisao contabil da regra |
| `fiscal_simulations` | Historico resumido de simulacoes |
| `fiscal_simulation_tax_details` | Memoria de calculo por imposto |

O banco do Motor Fiscal e separado do banco operacional do `sistema-fretes`.
Isso evita misturar dados de operacao com regras fiscais e auditoria tributaria.

## Regras Demonstrativas E Fallback

O MVP possui regras demonstrativas para permitir apresentacao fim a fim.

Como funciona:

- regras especificas por UF tem prioridade maior;
- se nao houver regra especifica, o motor usa regra fallback;
- o fallback considera `operation_type` e `customer_type`;
- regras fallback usam prioridade `900`;
- regras especificas devem usar prioridade menor, por exemplo `100`;
- regras fallback ficam como `PENDING_REVIEW`;
- fontes ficam como `INTERNAL_NOTE` e `PENDING_CONFIRMATION`.

Essa decisao permite demonstrar o fluxo completo sem afirmar que existe uma
base fiscal nacional validada.

Para producao, o caminho correto e:

1. cadastrar regras especificas;
2. informar fonte legal;
3. revisar com contador;
4. mudar status para `APPROVED`;
5. manter historico de vigencia.

## Variaveis De Ambiente

| Variavel | Exemplo | Uso |
|---|---|---|
| `APP_PORT` | `8080` | Porta HTTP da API |
| `APP_ENV` | `development` | Ambiente |
| `INTERNAL_API_KEY` | `dev-token` | Token interno entre Java e Go |
| `DATABASE_URL` | `postgres://motor_fiscal:motor_fiscal@localhost:5433/motor_fiscal?sslmode=disable` | Conexao PostgreSQL |

Use `.env.example` como base e nao versione `.env`.

## Como Rodar Com Docker

Subir:

```bash
docker compose up --build -d
```

Ver logs:

```bash
docker logs -f motor-fiscal
```

Parar:

```bash
docker compose down
```

Resetar banco e migrations:

```bash
docker compose down -v
docker compose up --build -d
```

Use `down -v` com cuidado, porque ele apaga o volume do PostgreSQL.

## Como Rodar Sem Docker

```bash
export APP_PORT=8080
export INTERNAL_API_KEY=dev-token
export DATABASE_URL='postgres://motor_fiscal:motor_fiscal@localhost:5433/motor_fiscal?sslmode=disable'
go run ./cmd/api
```

## Como Testar

Todos os testes:

```bash
GOCACHE=/tmp/go-build-cache go test ./...
```

Testes do motor de regras:

```bash
GOCACHE=/tmp/go-build-cache go test ./internal/service -v
```

Formatacao:

```bash
gofmt -w internal cmd
```

## Exemplo De Uso

Health:

```bash
curl -i http://localhost:8080/health
```

Simular frete:

```bash
curl -i -X POST http://localhost:8080/api/v1/tax/simulate \
  -H "Content-Type: application/json" \
  -H "X-API-Key: dev-token" \
  -H "X-Correlation-ID: req-frete-1001" \
  -d '{
    "freight_id": 1001,
    "operation_date": "2026-03-10",
    "origin_uf": "PE",
    "destination_uf": "SP",
    "freight_value": "3500.00",
    "customer_type": "PJ",
    "operation_type": "INTERESTADUAL"
  }'
```

Resposta esperada:

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
    }
  ],
  "from_cache": false
}
```

Erro sem token:

```bash
curl -i -X POST http://localhost:8080/api/v1/tax/simulate
```

Resposta esperada:

```json
{
  "code": "UNAUTHORIZED",
  "message": "Token ausente ou invalido.",
  "correlation_id": "req-...",
  "details": []
}
```

## Consultas Uteis

```sql
SELECT id, rule_code, status, priority, calculation_basis
FROM fiscal_rules
ORDER BY id;

SELECT id, freight_id, rule_code, rule_status, total_tax, created_at
FROM fiscal_simulations
ORDER BY id DESC;

SELECT fiscal_simulation_id, tax_name, base_value, rate, amount, formula
FROM fiscal_simulation_tax_details
ORDER BY id DESC;
```

Dados locais de conexao:

```text
Host: 127.0.0.1
Port: 5433
Database: motor_fiscal
Username: motor_fiscal
Password: motor_fiscal
SSL: disabled
```

## Documentacao Da API

Contrato OpenAPI:

```text
docs/openapi.yaml
```

Contrato explicativo:

```text
docs/api-contract.md
```

O arquivo OpenAPI pode ser aberto no Swagger Editor ou Swagger UI.


## Decisoes Tecnicas

- Go foi usado por ser simples, rapido e adequado para API pequena.
- Valores monetarios trafegam como string para evitar perda de precisao em JSON.
- Calculos usam `decimal`, nao `float`.
- Regras fiscais ficam no banco, nao hardcoded no codigo.
- Toda simulacao salva historico e memoria de calculo.
- `X-Correlation-ID` permite rastrear a chamada entre Java, Go e banco.
- Regras demonstrativas ficam como `PENDING_REVIEW`.
- Fonte legal e revisao contabil foram modeladas desde o MVP.
