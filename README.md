# AI Image Analyzer

Aplicacion full-stack para cargar una imagen, analizar su contenido con IA y mostrar etiquetas con su nivel de confianza.

## Tecnologias

- Backend: Go 1.22 + Gin
- Frontend: React 18 + Vite 5
- Contenerizacion: Docker + Docker Compose
- Web server frontend: Nginx (SPA + proxy a backend)

## Estructura

```text
.
├── backend
│   ├── cmd/api/main.go
│   ├── internal/
│   ├── test/
│   └── Dockerfile
├── frontend
│   ├── src/
│   ├── Dockerfile
│   ├── nginx.conf
│   └── package.json
└── docker-compose.yml
```

## Variables de entorno

### Backend

Archivo: `backend/.env` (puedes copiar desde `backend/.env.example`)

- `PORT`: puerto del backend (default `8080`)
- `GCV_API_KEY`: API key de Google Cloud Vision
- `GCV_API_URL`: URL de Google Cloud Vision (default `https://vision.googleapis.com/v1/images:annotate`)
- `MAX_IMAGE_SIZE`: tamano maximo permitido en bytes (default `5242880`)

### Frontend

Archivo: `frontend/.env` (puedes copiar desde `frontend/.env.example`)

- `VITE_API_BASE_URL`: base URL del backend. En local puedes dejarlo vacio para usar el proxy de Vite
- `VITE_MAX_IMAGE_SIZE_MB`: tamano maximo de archivo en MB para validacion en cliente (default `5`)

## Ejecutar con Docker Compose

Desde la raiz del proyecto:

```bash
docker compose up --build
```

- Frontend: `http://localhost:3000`
- Backend: `http://localhost:8080`

Detener servicios:

```bash
docker compose down
```

## Endpoint principal

`POST /api/analyze`

- Content-Type: `multipart/form-data`
- Campo esperado: `image`

Respuesta esperada:

```json
{
 "tags": [
  { "label": "Perro", "confidence": 0.98 },
  { "label": "Parque", "confidence": 0.91 }
 ]
}
```

## Cobertura de pruebas (backend)

- Componentes criticos (objetivo minimo: 80%):
  - `internal/handler`: 91.8%
  - `internal/service`: 91.8%

Verificar componentes criticos:

```bash
cd backend
go test --count=1 -cover ./internal/handler ./internal/service
```

- Componentes opcionales (referencia): resto de paquetes del backend.

Verificar cobertura global opcional:

```bash
cd backend
go test -cover ./...
```

## Buenas practicas aplicadas

- Variables de entorno para configuracion sensible
- Separacion de responsabilidades en frontend (componentes, servicio API, utilidades)
- Validacion de tipo y tamano de archivo en cliente
- Manejo de estados UX: loading, error, success
- UI responsiva para desktop y mobile
- Contenerizacion de frontend y backend para ejecucion reproducible
