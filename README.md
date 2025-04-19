# crype

A simple but secure crypto payments gateway, built with Golang.

## Features

- gRPC API for creating and managing crypto payment orders
- PostgreSQL database for order and payment address storage
- Blockchain wallet generation (currently supports USDC on Base)
- Dockerized development environment
- Code generation for database models (Jet) and gRPC (Protobuf)

## Project Structure

- `server/` - Main gRPC server implementation
- `api/proto/` - Protobuf definitions for gRPC services
- `api/gen_proto/` - Generated Go code from Protobuf
- `api/gen_jet/` - Generated database models and queries (Jet)
- `client/` - (Placeholder for client code)
- `utils/` - Utilities for blockchain and database
- `sql/` - Database schema
- `scripts/` - Helper scripts for code and DB generation

## Getting Started

### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- [Jet](https://github.com/go-jet/jet) CLI (for DB codegen)
- [protoc](https://grpc.io/docs/protoc-installation/) (for gRPC codegen)

### Setup

1. Copy `.env.example` to `.env` and fill in your environment variables.
2. Start the database:
   ```sh
   docker compose up -d
   ```
3. Generate database models (Jet):
   ```sh
   ./scripts/generate_jet.sh
   ```
4. Generate gRPC code:
   ```sh
   ./scripts/generate_proto.sh
   ```
5. Run the server:
   ```sh
   go run crype/server
   ```

### API

- gRPC endpoint: `localhost:8080`
- Service: `OrderService`
  - `CreateOrder(amount, currency)` → returns order ID, payment address, timestamps

See `api/proto/order.proto` for details.

### Database

- PostgreSQL, schema auto-initialized from `sql/schema.sql`
- Tables: `orders`, `payment_addresses`

### Blockchain

- Wallets generated for supported currencies (see `utils/blockchain.go`)
- Currently supports: `USDC_BASE`

## Development

- Use the provided scripts in `scripts/` for code generation and DB management.
- To reset the database:
  ```sh
  ./scripts/redeploy_db.sh
  ```

## License

MIT
