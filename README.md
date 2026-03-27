# ReCircuit

A full-stack refurbished electronics marketplace built with **Go** (backend) and **Vue 3** (frontend). Browse, buy, and sell refurbished phones, laptops, and tablets.

## Quick Start

```bash
docker compose up --build
```

- **API**: http://localhost:8080
- **Frontend**: http://localhost:5173

## Seed Accounts

| Role   | Email              | Password   |
|--------|--------------------|------------|
| Seller | seller@example.com | password1  |
| Buyer  | buyer@example.com  | password1  |

## Local Development (without Docker)

### Prerequisites

- Go 1.22+
- Node.js 20+
- PostgreSQL 16+

### Backend

```bash
cd backend
export DATABASE_URL="postgres://marketplace:marketplace@localhost:5432/marketplace?sslmode=disable"
export JWT_SECRET="dev-secret"
go mod tidy
make migrate-up
make run
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

## API Endpoints

| Method | Path                     | Auth     | Description                  |
|--------|--------------------------|----------|------------------------------|
| POST   | /api/v1/register         | -        | Create account               |
| POST   | /api/v1/login            | -        | Login, get JWT               |
| GET    | /api/v1/categories       | -        | List categories              |
| GET    | /api/v1/products         | -        | List/filter/search products  |
| GET    | /api/v1/products/{id}    | -        | Product detail               |
| POST   | /api/v1/products         | Seller   | Create product               |
| PUT    | /api/v1/products/{id}    | Seller   | Update own product           |
| DELETE | /api/v1/products/{id}    | Seller   | Delete own product           |
| GET    | /api/v1/seller/products  | Seller   | List my listings             |
| GET    | /api/v1/seller/stats     | Seller   | Sales stats                  |
| POST   | /api/v1/orders           | Buyer    | Place order                  |
| GET    | /api/v1/orders           | Auth     | List my orders               |
| GET    | /api/v1/orders/{id}      | Auth     | Order detail                 |
| GET    | /healthz                 | -        | Health check                 |

## Tech Stack

**Backend**: Go 1.22, chi router, pgx, golang-migrate, JWT, bcrypt, slog

**Frontend**: Vue 3, TypeScript, Vite, Tailwind CSS, Pinia, Vue Router

## Project Structure

```
recircuit-demo-app/
├── backend/           # Go REST API
│   ├── cmd/server/    # Entry point
│   ├── internal/      # App code (config, model, repository, service, handler, middleware)
│   └── migrations/    # SQL migrations
├── frontend/          # Vue 3 SPA
│   └── src/           # Components, views, stores, API layer
└── docker-compose.yml
```
