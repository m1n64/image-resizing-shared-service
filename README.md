# 🖼 Image Resizer Service

**Image Resizer** is a lightweight embeddable Go binary that runs a local HTTP + gRPC server for handling images. It accepts uploads, compresses them to WebP, generates thumbnails, and serves them locally via REST.

---

## 🚀 Quick Start

```bash
./imageresizer
```

You can also pass a custom config file (must be located next to the binary):

```bash
./imageresizer -c config.json
```

This will start:

- REST API server on `http://localhost:5689`
- gRPC server on `localhost:50066`

You can run it in the background or as a system service.

---

## 📦 Features

- ✅ Upload images via REST and gRPC
- ✅ Automatically compresses images into WebP format
- ✅ Generates multiple thumbnails per image
- ✅ Asynchronous background processing
- ✅ Serves static images at `/images/*`
- ✅ Stores metadata in embedded SQLite DB (file-based)
- ✅ Standalone binary — no Docker, no external dependencies
- ✅ Clean Architecture (DDD + Ports & Adapters)

> 💡 Looking for a full-featured, production-ready version of this service?
> Check out the full version with PostgreSQL, MinIO, Docker support and more:
> **https://github.com/m1n64/image-resizing-service/**
---

## 🛠️ Tech Stack

- **Go (Golang) 1.23** — core application logic
- **Gin** — HTTP web framework
- **gRPC** — transport layer for image handling
- **GORM** — ORM with SQLite backend (file-based)
- **SQLite in WAL mode** - file-based database in WAL mode for performance
- **imaging/chai2010/webp** — image processing and compression
- **Makefile** — CLI tasks and automation

---

## 📚 Config File Format

You can provide a JSON config file using `-c config.json`. This file must be located next to the binary.

```json
{
  "web_server_port": "5689",
  "grpc_server_port": "50066",
  "rest_upload_token": "",
  "grpc_token": "",
  "db_path": "",
  "image_compression": {
    "lossless": false,
    "quality": 80
  }
}
```

| Field                        | Type   | Description                                   |
|------------------------------|--------|-----------------------------------------------|
| `web_server_port`            | string | REST API port override                        |
| `grpc_server_port`           | string | gRPC server port override                     |
| `rest_upload_token`          | string | Optional API token for REST image upload auth |
| `grpc_token`                 | string | Optional API token for gRPC access            |
| `db_path`                    | string | Path to SQLite DB file (relative or absolute) |
| `image_compression.lossless` | bool   | Whether to use lossless WebP compression      |
| `image_compression.quality`  | number | WebP quality (0–100) if not lossless          |

If the file is missing or partial, missing values are filled from ENV vars or defaults.

---

## 🧪 REST Endpoints

| Method | Path                   | Description                     |
|--------|------------------------|---------------------------------|
| GET    | `/ping`                | Health check                    |
| POST   | `/image/upload`        | Upload image via multipart form |
| POST   | `/image/upload/binary` | Upload image via binary body    |
| GET    | `/image/:id`           | Retrieve image metadata & links |
| GET    | `/images/*`            | Access static images            |

### GET `/ping`
- **Response**:
```json
{
    "message": "pong",
    "timestamp": "2025-04-30T17:43:09.332177+03:00"
}
```

### POST `/image/upload`
- **Headers**: optional `X-API-Key: <token>` if `REST_UPLOAD_TOKEN` is set
- **Body**: multipart form with `file=<image>`
- **Response**:
```json
{
  "id": "510d1f37-751b-48de-a216-165f49d536c2",
  "original_url": "/images/originals/510d1f37-751b-48de-a216-165f49d536c2.jpg",
  "compressed_url": "",
  "size": 410451,
  "mime": "image/jpeg",
  "status": "pending",
  "thumbnails": null
}
```

### POST `/image/upload/binary`
- **Headers**: Content-Type = image/jpeg, image/png, etc. + optional `X-API-Key`
- **Body**: raw binary of image
- **Response**: same as `/image/upload`

### GET `/image/:id`
- **Response**:
```json
{
  "id": "18983da8-933d-4061-b94f-3a97fa52f755",
  "original_url": "/images/originals/18983da8-933d-4061-b94f-3a97fa52f755.jpg",
  "compressed_url": "/images/webp/18983da8-933d-4061-b94f-3a97fa52f755.webp",
  "size": 410451,
  "mime": "image/jpeg",
  "status": "ready",
  "thumbnails": [
    {
      "size": "100x100",
      "url": "/images/thumbnails/18983da8-933d-4061-b94f-3a97fa52f755_100x100.webp",
      "type": "tiny"
    },
    {
      "size": "1024x768",
      "url": "/images/thumbnails/18983da8-933d-4061-b94f-3a97fa52f755_1024x768.webp",
      "type": "xlarge"
    },
    ...
  ]
}
```

### GET `/images/*`
- Serves static files
- Example: `/images/originals/<uuid>.jpg`, `/images/thumbs/<uuid>_200x200.webp`


If `REST_UPLOAD_TOKEN` is set, requests to upload endpoints must include:
```
X-API-Key: your_token_here
```

---

## 📡 gRPC API

Service definition:

```proto
service ImageService {
  rpc UploadImage(UploadImageRequest) returns (ImageResponse) {}
  rpc GetImage(GetImageRequest) returns (ImageResponse) {}
}

message UploadImageRequest {
  bytes data = 1;
}

message GetImageRequest {
  string id = 1;
}

message ThumbnailShort {
  string size = 1;
  string url = 2;
  string type = 3;
}

message ImageResponse {
  string id = 1;
  string original_url = 2;
  optional string compressed_url = 3;
  string status = 4;
  int64 size = 5;
  string mime = 6;
  optional string error_message = 7;
  repeated ThumbnailShort thumbnails = 8;
}
```

If `GRPC_TOKEN` is set, gRPC requests must include:
```
authorization: your_token_here
```

Proto location: `./proto/image_service.proto`

---

## ⚙️ Environment Variables

| Variable            | Default        | Description                                                                                  |
|---------------------|----------------|----------------------------------------------------------------------------------------------|
| `WEB_SERVER_PORT`   | `5689`         | Port for the REST API server                                                                 |
| `GRPC_SERVER_PORT`  | `50066`        | Port for the gRPC server                                                                     |
| `REST_UPLOAD_TOKEN` | _(empty)_      | Token required for REST image upload endpoints (optional)                                    |
| `GRPC_TOKEN`        | _(empty)_      | Token required for gRPC access (optional)                                                    |
| `DB_PATH`           | `data/data.db` | Path to the SQLite database (relative or absolute). Uploads folder will be placed alongside. |

---

## 🗃 File Storage Layout

Images are stored next to the SQLite DB. If `DB_PATH` is set to `./data/data.db`, images will be saved in:

```
data/uploads/originals/
                 /webp/
                 /thumbs/
```

Static image access (the server response returns the non-host portion of the image URL):
```
http://<your_host>:WEB_SERVER_PORT/images/originals/<id>.jpg
http://<your_host>:WEB_SERVER_PORT/images/webp/<id>.webp
http://<your_host>:WEB_SERVER_PORT/images/thumbs/<id>_200x200.webp
```

---

## 🧱 Compilation Instructions (No Docker)

Because the project uses CGO (for WebP support), you must build it on Linux **or** using a Linux environment.

### ✅ Build for local Linux system:

```bash
CGO_ENABLED=1 go build -ldflags "-X main.GinMode=release" -o ./tmp/imageresizer ./cmd/server
```

### ✅ Cross-compile from macOS (via Docker):

Create a `Dockerfile.build`:
```Dockerfile
FROM golang:1.23
RUN apt update && apt install -y libwebp-dev build-essential
WORKDIR /app
COPY . .
ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64 GIN_MODE=release
RUN go build -o imageresizer ./cmd/server
```

Then run:
```bash
docker build -f Dockerfile.build -t imageresizer-builder .
docker run --rm -v "$PWD/tmp:/out" imageresizer-builder cp /app/imageresizer /out/
```

Or use
```bash
docker build -t m1n64/imageresizer:latest .
```
```bash
docker run \
  -p 5689:5689 \
  -p 50066:50066 \
  -v "$(pwd)/data:/app/data" \
  m1n64/imageresizer
```
**Note! This solution is not tested by the author!**

---

## 📋 Makefile Commands

| Command      | Description                                      |
|--------------|--------------------------------------------------|
| `make build` | Build the project (CGO enabled)                  |
| `make run`   | Run the project in dev mode with `go run`        |
| `make prod`  | Build the production binary with release flags   |
| `make dev`   | Run using Air (requires Air and Delve)           |
| `make debug` | Run under Delve debugger (headless on port 2345) |
| `make clean` | Remove the `./tmp` directory                     |
| `make proto` | Compile `.proto` definitions with `protoc`       |

---

## 🐳 Running from Docker Hub

You can use the prebuilt image directly in your Docker Compose setup:

```yaml
version: '3.9'

services:
  imageresizer:
    image: m1n64/imageresizer:latest
    ports:
      - "5689:5689"
      - "50066:50066"
    volumes:
      - imageresizer_data:/app/data
      # - ./config/imageresizer.json:/app/config.json
    # command: imageresizer -c config.json
    environment:
      - WEB_SERVER_PORT=5689
      - GRPC_SERVER_PORT=50066
```

This setup:
- Runs the latest Docker Hub image
- Binds REST and gRPC ports
- Mounts the host `./data` folder to persist uploaded files and the SQLite DB

---

## 🧩 Available SDK

- [PHP](https://github.com/m1n64/image-resizer-sdk) - A Lightweight PHP SDK for the Image Resizer Service
- [PHP (Laravel)](https://github.com/junior-idiot/laravel-imageresizer-sdk) - A Laravel wrapper for PHP SDK

---

## 📌 Notes

- The service is designed to be embedded in desktop or backend applications.
- It runs a local REST + gRPC server for internal communication.
- Files are not intended to be exposed directly on the filesystem.
- Public image URLs are routed through the REST server.

---

## ✅ Example Integration

You can integrate this binary into any local system or backend. Use gRPC or REST to:

- Upload files
- Retrieve URLs for WebP and thumbnails
- Get status and metadata

This approach allows local apps to delegate all image handling (resize, compression, access) to a clean standalone service.

---

## ✨ Planned Features

- Dockerized version of the binary on Docker Hub
- Change SQLite in WAL mode to a faster embedded database (libSQL or similar)
- Embedded file-layer abstraction (MinIO-like behavior without exposing FS)
- Optional encryption or obfuscation for stored image files
- Optional signed URLs with expiration
- Configurable image size presets and quality profiles
- Image deletion and garbage collection
- Built-in health checks and metrics

---

## 🧙‍♂️ Author

Made with ❤️ by the **[Kirill Sakharov](https://github.com/m1n64) ([LinkedIn](https://www.linkedin.com/in/kirill-sakharov-862072227/))**
