# Guia da Estrutura do Projeto

Este documento resume a responsabilidade de cada diretório do projeto `motor-fiscal-go`.

## Estrutura geral

```text
motor-fiscal-go/
├── cmd/api
├── internal/config
├── internal/handler
├── internal/service
├── internal/client
├── internal/dto
├── internal/model
├── internal/errors
├── internal/validator
├── internal/logger
├── internal/repository
├── internal/cache
├── docs
├── scripts
└── tests
```

## Responsabilidade dos diretórios

| Diretório | Responsabilidade |
|---|---|
| `cmd/api` | Ponto de entrada da aplicação. Contém o `main.go`, responsável por iniciar o servidor HTTP, carregar configurações e registrar rotas. |
| `internal/config` | Configurações da aplicação: porta, variáveis de ambiente, tokens, banco, Redis, URLs externas e timeouts. |
| `internal/handler` | Camada HTTP. Recebe requisições, interpreta dados de entrada, chama os services e devolve respostas padronizadas. |
| `internal/service` | Camada de regra de negócio. Orquestra os fluxos fiscais, valida decisões e coordena chamadas para clients, repository e cache. |
| `internal/client` | Comunicação com APIs externas, como TecnoSpeed, ViaCEP, SEFAZ, serviços de rota e outras integrações. |
| `internal/dto` | Estruturas de entrada e saída em JSON. Representa contratos de request e response entre sistemas. |
| `internal/model` | Estruturas internas do domínio fiscal e logístico usadas dentro da aplicação. |
| `internal/errors` | Padronização de erros, códigos internos, mensagens e conversão para respostas HTTP. |
| `internal/validator` | Validações de payload, campos obrigatórios, formatos, limites e regras simples de entrada. |
| `internal/logger` | Configuração e centralização dos logs estruturados da aplicação. |
| `internal/repository` | Comunicação com banco de dados, como PostgreSQL, quando houver necessidade de persistência. |
| `internal/cache` | Integração com cache, como Redis, para dados temporários, tokens ou consultas frequentes. |
| `docs` | Documentação técnica, decisões arquiteturais, contratos, fluxos e exemplos de payload. |
| `scripts` | Scripts auxiliares para setup, build, testes, Docker ou execução local. |
| `tests` | Testes de integração, testes de contrato e cenários ponta a ponta. |

## Fluxo principal entre camadas

```text
handler -> service -> client/repository/cache
```

## Regras de responsabilidade

- `handler` deve lidar apenas com entrada e saída HTTP.
- `service` deve concentrar as regras de negócio e decisões do fluxo.
- `client` deve conhecer os detalhes das APIs externas.
- `repository` deve conhecer os detalhes do banco de dados.
- `cache` deve centralizar o uso de Redis ou outro mecanismo de cache.
- `dto` deve representar dados trafegados entre sistemas.
- `model` deve representar conceitos internos da aplicação.
- `errors` deve padronizar os erros retornados pela API.
- `validator` deve validar dados antes que a regra de negócio seja executada.
- `logger` deve centralizar a forma como a aplicação registra logs.

