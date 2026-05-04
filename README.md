# Motor Fiscal

Microserviço fiscal em Go para cálculo, validação e simulação tributária em operações de frete.

## Contexto
O sistema principal de fretes é responsável pela operação logística. O Motor Fiscal centraliza regras fiscais para evitar lógica espalhada no Java/JSP/PostgreSQL.

## Problema Resolvido
Regras de ICMS, IBS, CBS, CFOP e validações fiscais precisam ser versionadas, testáveis e rastreáveis.

## Arquitetura
Sistema Java → API REST Go → regras/cache/banco/API externa → resposta fiscal.

## Principais Endpoints
- GET /health
- POST /api/v1/tax/simulate
- POST /api/v1/tax/compare
- POST /api/v1/tax/batch
- POST /api/v1/cte/validate

## Estrutura
Descrever cmd, internal, handler, service, client, dto, model, config, errors, validator, logger, tests.

## Variáveis De Ambiente
- APP_PORT
- DATABASE_URL
- REDIS_ADDR
- INTERNAL_API_KEY
- EXTERNAL_API_BASE_URL
- EXTERNAL_API_TOKEN
- HTTP_TIMEOUT_SECONDS

## Como Rodar
- go run ./cmd/api
- docker compose up

## Como Testar
- go test ./...

## Roadmap
- Health
- Simulação fiscal
- Regras versionadas
- Validação CT-e
- Cache Redis
- Histórico PostgreSQL
- Integração Java
- APIs externas futuras

## Checklist
- Sem float para dinheiro
- Logs com correlation ID
- Token fora do Git
- Testes de cálculo
- Tratamento de erro padronizado
