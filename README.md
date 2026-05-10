# Motor Fiscal

Microservico fiscal em Go para calculo, validacao, simulacao e auditoria tributaria em operacoes de frete.

O projeto nasceu para integrar um sistema principal de fretes em Java/JSP com um servico especializado em regras fiscais. A ideia e remover calculos fiscais espalhados pelo sistema principal e centralizar essas decisoes em uma API menor, testavel, versionada e explicavel.

## Status

Projeto em fase de MVP tecnico com motor de regras proprio.

Ja implementado:

- API REST em Go;
- Docker com PostgreSQL;
- autenticacao por `X-API-Key`;
- `X-Correlation-ID` para rastreio;
- erros padronizados em JSON;
- motor de regras fiscais por condicoes;
- calculo auditavel de ICMS, IBS e CBS;
- persistencia da simulacao e dos detalhes do calculo;
- validacao basica de CT-e;
- contrato OpenAPI;
- testes do motor de regras e do calculo fiscal.

Aviso de escopo: as regras fiscais atuais sao demonstrativas e estao marcadas como `PENDING_REVIEW`. Para uso produtivo, as regras devem receber fonte oficial confirmada e validacao contabil. O objetivo deste MVP e demonstrar arquitetura, rastreabilidade, motor de regras e memoria de calculo.

## Problema Resolvido

No sistema de fretes, dados fiscais como CFOP, aliquotas, bases de calculo e totais tributarios nao devem ser digitados manualmente nem ficar duplicados em varias telas ou servlets.

O Motor Fiscal centraliza:

- selecao de regra fiscal vigente;
- calculo de impostos;
- memoria de calculo;
- historico de simulacoes;
- auditoria da regra aplicada;
- contrato HTTP estavel para integracao com Java.

## Fluxo De Integracao

```text
Sistema de Fretes Java/JSP
        |
        | POST /api/v1/tax/simulate
        | X-API-Key
        | X-Correlation-ID
        v
Motor Fiscal Go
        |
        | busca regras candidatas
        v
PostgreSQL: fiscal_rules, fiscal_rule_conditions, fiscal_rule_taxes
        |
        | calcula impostos e salva auditoria
        v
PostgreSQL: fiscal_simulations, fiscal_simulation_tax_details
        |
        | retorna JSON
        v
Sistema de Fretes Java/JSP
```

Na tela de cadastro de frete, o botao **Motor Fiscal** deve chamar um endpoint do proprio sistema Java. Esse endpoint Java chama o Motor Fiscal em Go, recebe o resultado e devolve para a JSP preencher campos readonly, como:

- CFOP;
- aliquota ICMS;
- valor ICMS;
- aliquota IBS;
- valor IBS;
- aliquota CBS;
- valor CBS;
- total de tributos;
- total com tributos;
- regra aplicada;
- status da regra.

## Arquitetura

```text
cmd/api
internal/config
internal/handler
internal/middleware
internal/service
internal/repository
internal/dto
internal/model
internal/errors
internal/validator
docs
migrations
```

Responsabilidades principais:

| Pasta | Responsabilidade |
|---|---|
| `cmd/api` | Ponto de entrada da aplicacao |
| `internal/config` | Leitura de variaveis de ambiente |
| `internal/handler` | Handlers HTTP e rotas |
| `internal/middleware` | API Key e correlation ID |
| `internal/service` | Regras de negocio, matching e calculo |
| `internal/repository` | Acesso ao PostgreSQL |
| `internal/dto` | Contratos de request e response |
| `internal/model` | Modelos internos do dominio fiscal |
| `internal/errors` | Estrutura padrao de erro da API |
| `internal/validator` | Validacao dos payloads |
| `docs` | OpenAPI e contrato explicativo |
| `migrations` | Scripts SQL executados pelo PostgreSQL no Docker |

## Endpoints

| Metodo | Rota | Autenticacao | Responsabilidade |
|---|---|---|---|
| `GET` | `/health` | Nao | Verificar se a API esta online |
| `POST` | `/api/v1/tax/preview` | Sim | Previsualizar calculo fiscal antes de salvar o frete |
| `POST` | `/api/v1/tax/simulate` | Sim | Calcular impostos de um frete |
| `POST` | `/api/v1/tax/compare` | Sim | Comparar modelo atual com cenario da reforma |
| `POST` | `/api/v1/tax/batch` | Sim | Processar multiplos fretes |
| `POST` | `/api/v1/cte/validate` | Sim | Validar dados minimos para CT-e |

Endpoints futuros documentados:

- `POST /api/v1/cte/preview`;
- `GET /api/v1/rules`;
- `POST /api/v1/rules`;
- `GET /api/v1/reports/tax-summary`;
- `GET /api/v1/reports/reform-impact`;
- `GET /api/v1/reports/inconsistencies`.

## Motor De Regras

O motor fiscal nao usa aliquotas fixas no codigo. Ele seleciona regras a partir do banco.

Tabelas principais:

| Tabela | Papel |
|---|---|
| `fiscal_rules` | Regra fiscal, vigencia, prioridade, status e CFOP |
| `fiscal_rule_conditions` | Condicoes para uma regra ser aplicada |
| `fiscal_rule_taxes` | Impostos e aliquotas da regra |
| `fiscal_rule_sources` | Fonte legal ou fonte de referencia da regra |
| `fiscal_rule_accounting_reviews` | Validacao contabil da regra |
| `fiscal_simulations` | Historico resumido das simulacoes |
| `fiscal_simulation_tax_details` | Memoria de calculo por imposto |

Ordem conceitual:

```text
1. buscar regras candidatas por vigencia e status;
2. avaliar condicoes;
3. ordenar por prioridade;
4. detectar conflitos;
5. calcular impostos;
6. salvar auditoria;
7. retornar resposta fiscal.
```

Operadores suportados nas condicoes:

- `EQUALS`;
- `NOT_EQUALS`;
- `IN`;
- `BETWEEN`.

Status de regra:

| Status | Significado |
|---|---|
| `DRAFT` | Rascunho |
| `PENDING_REVIEW` | Pendente de validacao contabil |
| `APPROVED` | Validada para uso |
| `INACTIVE` | Desativada |

### Cobertura Demonstrativa Nacional

Para o MVP, o projeto possui regras fallback demonstrativas que permitem calcular fretes para qualquer combinacao de UF.

Essas regras existem para deixar o fluxo completo apresentavel:

```text
cadastro do frete -> preview fiscal -> emissao do frete -> calculo definitivo -> auditoria
```

Como funciona:

- regras especificas por UF continuam tendo prioridade maior;
- se nao existir regra especifica para a combinacao de origem/destino, o motor usa uma regra fallback;
- o fallback considera `operation_type` e `customer_type`;
- regras fallback usam prioridade `900`;
- regras especificas devem usar prioridade menor, por exemplo `100`;
- as regras fallback ficam como `PENDING_REVIEW`;
- a fonte fica marcada como `INTERNAL_NOTE` e `PENDING_CONFIRMATION`.

Essa escolha evita dois problemas:

- nao travar a demonstracao por falta de uma base fiscal nacional completa;
- nao fingir que regras demonstrativas sao regras fiscais oficiais.

Para producao, o caminho correto e cadastrar regras especificas com fonte legal confirmada, vigencia, revisao contabil e status `APPROVED`.

## Variaveis De Ambiente

| Variavel | Exemplo | Uso |
|---|---|---|
| `APP_PORT` | `8080` | Porta HTTP da API |
| `APP_ENV` | `development` | Ambiente da aplicacao |
| `INTERNAL_API_KEY` | `dev-token` | Token interno entre Java e Go |
| `DATABASE_URL` | `postgres://motor_fiscal:motor_fiscal@localhost:5433/motor_fiscal?sslmode=disable` | Conexao local com PostgreSQL |

Crie seu arquivo local a partir de `.env.example`.

Nao commitar `.env`.

## Como Rodar Com Docker

Subir a aplicacao:

```bash
docker compose up --build -d
```

Ver containers:

```bash
docker ps
```

Ver logs:

```bash
docker logs -f motor-fiscal
```

Parar:

```bash
docker compose down
```

Resetar banco e rodar migrations do zero:

```bash
docker compose down -v
docker compose up --build -d
```

Use `down -v` com cuidado, porque ele apaga o volume do PostgreSQL.

## Como Rodar Sem Docker

Suba um PostgreSQL local e configure:

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

Testes do motor fiscal:

```bash
GOCACHE=/tmp/go-build-cache go test ./internal/service -v
```

Formatar codigo:

```bash
gofmt -w internal cmd
```

## Exemplos De Uso

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

Testar sem token:

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

## DBeaver Ou pgAdmin

Dados de conexao:

```text
Host: 127.0.0.1
Port: 5433
Database: motor_fiscal
Username: motor_fiscal
Password: motor_fiscal
SSL: disabled
```

Consultas uteis:

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

## Documentacao Da API

Contrato OpenAPI:

```text
docs/openapi.yaml
```

Contrato explicativo:

```text
docs/api-contract.md
```

Voce pode abrir `docs/openapi.yaml` no Swagger Editor ou Swagger UI.

## Como Integrar No Sistema Java

Desenho recomendado:

```text
JSP
  -> Servlet Java /motor-fiscal/simular
    -> POST http://localhost:8080/api/v1/tax/simulate
      -> Motor Fiscal Go
    <- JSON fiscal
  <- JSON para a tela
```

Evite chamar o Motor Fiscal direto do JavaScript do navegador, porque isso exporia o `X-API-Key`.

No Java, o token deve ficar no backend, por configuracao de ambiente.

## Decisoes Tecnicas

- Go foi usado por ser simples, rapido, bom para APIs pequenas e facil de empacotar em Docker.
- Valores monetarios trafegam como string para evitar perda de precisao.
- Calculos usam decimal, nao `float`.
- Regras fiscais ficam no banco, nao hardcoded no codigo.
- Toda simulacao salva historico e memoria de calculo.
- `X-Correlation-ID` permite rastrear a mesma operacao entre Java, Go e banco.
- Regras demonstrativas ficam como `PENDING_REVIEW` ate validacao contabil.

## Limitacoes Atuais

- As regras atuais sao exemplos de desenvolvimento.
- Nao existe tela administrativa para cadastrar regras.
- Nao existe cache Redis implementado.
- Nao existe integracao com emissao real de CT-e.
- Nao existe validacao fiscal completa para todos os cenarios brasileiros.
- Fontes legais precisam ser cadastradas e confirmadas antes de uso produtivo.

## Roadmap

- adicionar logs estruturados por request;
- criar endpoints administrativos de regras fiscais;
- criar relatorios fiscais;
- integrar com o sistema Java de fretes;
- adicionar seed de regras reais validadas por contador;
- evoluir validacao CT-e;
- preparar pipeline de CI.

## Checklist Profissional

- [x] API REST em Go
- [x] Docker
- [x] PostgreSQL
- [x] OpenAPI
- [x] API Key
- [x] Correlation ID
- [x] Erros padronizados
- [x] Motor de regras
- [x] Calculo auditavel
- [x] Persistencia da auditoria
- [x] Testes do motor fiscal
- [x] README publico de apresentacao
- [ ] Logs estruturados
- [ ] CI
- [ ] Integracao com sistema Java
- [ ] Regras fiscais reais aprovadas

## Apresentacao Tecnica

Este projeto pode ser apresentado como um microservico fiscal criado para resolver um problema real de acoplamento no sistema de fretes.

Pontos fortes para apresentar:

- separacao entre sistema operacional de frete e decisao fiscal;
- motor de regras com vigencia, prioridade e condicoes;
- auditoria de calculo por imposto;
- historico persistido;
- contrato OpenAPI;
- Docker e PostgreSQL;
- testes automatizados no nucleo fiscal;
- preocupacao com rastreabilidade e validacao contabil.
