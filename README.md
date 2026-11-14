# kali-auth-context

Servicio de autenticación y autorización de InvoiceKali (empresas, usuarios, roles,
permisos y emisión de access tokens RS256). Arquitectura y contrato:
[kali-auth-context-architecture.md](kali-auth-context-architecture.md).

## Arranque en desarrollo (Docker)

Lo habitual es levantar todo el sistema desde `invoice-kali-front` con `npm run stack:up`
(ver su README). Para levantar solo este servicio:

```bash
cp .env.example .env.dev        # solo si no existe; valores de ejemplo solo locales
docker compose -f docker-compose.dev.yml --env-file .env.dev up -d --build
curl http://localhost:18080/health
```

Al arrancar:

1. **Claves JWT**: el servicio `jwt-keys` genera `.secrets/jwt_private.pem` y
   `.secrets/jwt_public.pem` (ignorados por git) si no existen. Nunca sobrescribe claves
   existentes. Fuera de Docker: `make jwt-keys-dev` (requiere `openssl`).
2. **Base de datos**: Postgres crea el usuario y la base de `POSTGRES_*` la primera vez
   (volumen vacío). Las **migraciones** se aplican solas al iniciar la API.
3. **Empresa y usuario iniciales**: con `SEED_MASTER_ENABLED=true` se crean (si no existen)
   la empresa de plataforma `master-tenant` (código `master`), el usuario
   `master@local.dev` con la contraseña `SEED_MASTER_PASSWORD` y el rol administrador.
   El valor de `.env.example` es de ejemplo y solo para uso local.

| Recurso | URL |
| --- | --- |
| API | http://localhost:18080 (`/health`) |
| JWKS | http://localhost:18080/.well-known/jwks.json |
| Adminer | http://localhost:18081 |

`JWT_ISSUER` y `JWT_AUDIENCE` deben ser idénticos a los de `kali-invoice-service`
(`http://localhost:8080` y `account`).

## Pruebas

```bash
go build -mod=mod ./... && go vet -mod=mod ./... && go test -mod=mod ./...   # o: make test
```

Se usa `-mod=mod` (dependencias desde la caché de módulos, según `go.mod`/`go.sum`) porque, tras
`docker compose ... up` (o `npm run stack:up` del front), Docker crea en el host la carpeta
`vendor/` del montaje `- /app/vendor`; si esa carpeta queda vacía o desactualizada, Go usa
`-mod=vendor` por defecto y falla con `inconsistent vendoring`. `-mod=mod` no modifica archivos.
