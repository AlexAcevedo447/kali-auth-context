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
go build ./... && go vet ./... && go test ./...
```
