# toy-api

A minimal Go HTTP API that exposes endpoints for demonstrating containerization concepts with Docker.

## Getting Started

1. Install [Go](https://go.dev/dl/) 1.21 or newer.
2. Clone the repository and install dependencies (standard library only):
   ```bash
   go mod tidy
   ```
3. Run the server:
   ```bash
   go run ./cmd/server
   ```

The server listens on port `8080` by default. Override the port by setting the `PORT` environment variable, e.g. `PORT=9090`.

## Available Endpoints

- `GET /healthz` – Returns a simple status payload to verify the service is running.
- `GET /api/v1/message` – Returns a friendly message payload suitable for quick demos.
- `GET /api/v1/toys` – Lists the in-memory toy catalog.
- `POST /api/v1/toys` – Adds a new toy to the catalog.
- `GET /api/v1/toys/{id}` – Retrieves a single toy by identifier.
- `PUT /api/v1/toys/{id}` – Replaces a toy with the provided payload.
- `DELETE /api/v1/toys/{id}` – Removes a toy from the catalog.

### Toy Payload

```json
{
  "name": "Rocket Racer",
  "category": "Vehicles",
  "price": 19.99,
  "in_stock": true
}
```

The toy endpoints use an in-memory data store seeded with a few example toys to keep the service stateless—perfect for Docker demos.

## Container Image

Build and run the service in a container:

```bash
docker build -t toy-api .
docker run --rm -p 8080:8080 toy-api
```

Override the port by supplying `-e PORT=9090` and adjusting the published port, for example:

```bash
docker run --rm -p 9090:9090 -e PORT=9090 toy-api
```

Images built on the `main` branch are automatically pushed to GitHub Container Registry under `ghcr.io/<owner>/<repo>` with both `linux/amd64` and `linux/arm64` support. Pull the latest branch image with:

```bash
docker pull ghcr.io/<owner>/<repo>:main
```

## Continuous Integration

GitHub Actions run automatically:

- Pull requests targeting `main` execute `go test ./...` to guard merges.
- Commits on `main` build and publish the multi-architecture Docker image to GHCR to validate the Dockerfile.

## Next Steps

- Containerize the service with Docker and explore build caching layers.
- Add persistence or configuration options as needed for more advanced demos.
- Expand the API with additional resources or filtering endpoints.
