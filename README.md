# Checkly API

API responsável por monitorar URLs, verificar sua disponibilidade
periodicamente e avaliar sua saúde com base em uma máquina de estados.

O sistema executa checks HTTP, avalia o estado de cada URL
(Healthy, Degraded, Recovering, Down) e agenda automaticamente
as próximas verificações.

---

## 🧠 Visão Geral

O fluxo principal da aplicação funciona da seguinte forma:

1. Uma URL é cadastrada para monitoramento
2. O sistema executa um check HTTP inicial
3. A URL recebe um estado inicial (Healthy ou Degraded)
4. Novos checks são agendados automaticamente
5. A cada verificação, o estado da URL é reavaliado

A lógica de negócio é isolada em services de aplicação
e o domínio é mantido independente de infraestrutura.

---

## 🏗️ Arquitetura

O projeto segue princípios de **Clean Architecture**, com separação clara
de responsabilidades:

internal/
modules/
urls/
domain/ → Entidades e regras centrais
application/ → Services e orquestração
presentation/ → HTTP handlers e DTOs
external/ → Infraestrutura (DB, etc)

---

### Camadas

- **Domain**

  - Estados da URL
  - Entidades
  - Regras centrais do negócio

- **Application**

  - Orquestração de casos de uso
  - Máquina de estados
  - Coordenação entre domínio e persistência

- **Presentation**

  - Endpoints HTTP
  - Validação de entrada
  - DTOs

- **External**
  - Implementações de repositórios
  - Banco de dados
  - Integrações externas

---

## 📡 API HTTP

A API expõe endpoints REST para gerenciamento e monitoramento de URLs.

A documentação completa da API está disponível via **Swagger**:

GET /swagger/

Exemplos de endpoints:

- `POST /urls`
- `GET /urls`
- `PUT /urls/{id}`

> O Swagger é a fonte de verdade do contrato HTTP.

---

## 🧪 Testes

Os testes são focados principalmente em:

- Services de aplicação
- Regras de domínio
- Casos de uso

Para executar todos os testes:

bash
go test ./...

---

🚀 Como rodar o projeto
Pré-requisitos

Go 1.22+

Executar localmente
go run cmd/api/main.go

---

📚 Documentação

Swagger → contrato HTTP e exemplos de uso

godoc → regras de domínio e serviços de aplicação

README → visão geral e onboarding

📌 Observações

Regras de negócio não estão acopladas à camada HTTP

Services de aplicação não dependem de frameworks

Repositórios são acessados via interfaces
