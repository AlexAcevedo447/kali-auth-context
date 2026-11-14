# Kali Auth Context

## Propósito
Microservicio de identidad para un ERP modular. Hoy cubre autenticación, autorización, usuarios, tenants y RBAC, con aislamiento por tenant como regla base del dominio y de la aplicación.

## Stack
- Go
- Fiber
- Bun ORM sobre database/sql
- PostgreSQL
- Wire para composición
- Zap logger

## Arquitectura
- DDD + Hexagonal + CQRS
- Commands y queries separados por caso de uso
- Un contrato por responsabilidad en ports
- Dominio sin dependencias técnicas
- HTTP organizado por contexto y por commands/queries
- Handlers SRP: un handler por endpoint/unidad de trabajo

## Estructura real
- cmd/api/main.go: composition root, crea app Fiber, resuelve contenedor y registra router
- internal/domain/identity: user, tenant, role, permission, user_role, role_permission, errores de auth/autorización
- internal/domain/policies: password_policy y authorization_policy
- internal/application:
  - auth/commands: login
  - authorization/queries: authorize
  - user/commands y user/queries
  - tenant/commands y tenant/queries
  - rbac/commands y rbac/queries
- internal/ports: contratos separados para user, tenant, role, permission, user_role, role_permission, password_hasher y uuid_provider
- internal/infrastructure/db:
  - models y mappers
  - repositories separados por contexto y por command/query
- internal/infrastructure/http:
  - contexts/auth, authorization, user, tenant, rbac
  - cada contexto separado en commands/queries
  - routes/router.go registra root, health y api/v1
  - shared/errors.go centraliza mapeo de errores HTTP
- internal/bootstrap/di: ensamblado completo del contenedor y del router

## Reglas funcionales actuales
- User siempre pertenece a un tenant
- Tenant puede estar ACTIVE o SUSPENDED
- Authorization valida:
  - request válida
  - tenant activo
  - usuario perteneciente al tenant
  - permisos efectivos por roles
- Login por correo y contraseña: el servidor resuelve el tenant de la cuenta. Si el
  mismo correo y contraseña existen en varias empresas responde `409 company code required`
  y el cliente reenvía `tenant_code` (código legible de empresa, nunca el UUID).
  Un tenant suspendido responde `403 tenant is suspended`.
- El JWT (RS256) lleva el tenant efectivo en el claim `tid`; los servicios lo validan con
  el JWKS publicado en `/.well-known/jwks.json`. `JWT_ISSUER`/`JWT_AUDIENCE` son identificadores
  lógicos y deben ser idénticos en este servicio y en `kali-invoice-service`
  (valores unificados: `http://localhost:8080` y `account`).
- Aislamiento multi-tenant (`middleware.TenantGuards`):
  - el `tenant_id` de query/body se fija al del token; si el cliente envía otro → `403 tenant mismatch`;
  - roles, permisos y usuarios referenciados por id (ruta, query o body) deben pertenecer al
    tenant efectivo; si no → `404 resource not found` (no se revela que existen);
  - no se puede asignar un rol de una empresa a un usuario de otra;
  - `/api/v1/tenants` solo para el **administrador de plataforma**: usuario de la empresa de
    plataforma (`SEED_MASTER_TENANT_ID`) con el permiso RBAC `tenants/manage` (o comodín; el
    admin sembrado tiene `*/*`). Un usuario de la plataforma sin ese permiso → `403 platform
    administrator required`. Solo ese administrador puede indicar un tenant destino existente al
    crear usuarios/roles; cualquier otro → `403 tenant mismatch`;
  - `ActiveTenant`: toda ruta protegida exige que la empresa del token esté activa (consulta en
    cada petición). Una empresa suspendida → `403 tenant is suspended`, también con tokens emitidos
    antes de la suspensión; al reactivarla, los tokens aún vigentes vuelven a funcionar.
- Tenant con `code` legible y único (`SEED_MASTER_TENANT_CODE` para el maestro; para los
  existentes se deriva del nombre al migrar).
- Administración de empresas (`/api/v1/tenants`, solo plataforma):
  - el nombre se guarda sin espacios a los lados; vacío o solo espacios → `400 tenant name is required`;
  - `PUT` reemplaza nombre, código y estado; un estado distinto de `ACTIVE`/`SUSPENDED` →
    `400 tenant status must be ACTIVE or SUSPENDED` (vacío conserva el comportamiento previo: `ACTIVE`);
  - editar, activar o suspender una empresa inexistente → `404 resource not found`;
  - la empresa de plataforma no se puede suspender (`/suspend` o `PUT` con `SUSPENDED`) →
    `409 platform tenant cannot be suspended`;
  - el login devuelve `is_platform_tenant` (true solo para `SEED_MASTER_TENANT_ID`), informativo
    para la interfaz; la autorización la aplica `PlatformAdmin`.
- `GET /api/v1/auth/session` (protegido): 200 si el token es válido y su empresa está activa.
  `kali-invoice-service` lo consulta con el mismo token para bloquear empresas suspendidas.
- Login mitiga enumeración con fake hash y soporta NeedsRehash
- Password policy exige mínimo 12 caracteres, mayúscula, minúscula, dígito, símbolo y sin espacios
- Password hasher actual: Argon2id con compatibilidad de verificación para hashes bcrypt heredados

## API expuesta
- / and /health
- /api/v1/auth/login
- /api/v1/auth/authorize
- /api/v1/users
- /api/v1/tenants
- /api/v1/roles
- /api/v1/permissions
- /api/v1/rbac

## Estado actual
- Core domain/application/infrastructure implementado
- Repositories Bun implementados por responsabilidad
- Router principal conectado desde main
- HTTP migrado a handlers SRP por contexto
- Paquete legacy de handlers eliminado
- go test ./... en verde

## Notas de continuidad
- Mantener tenant propagation en dominio, aplicación, repositorios y HTTP
- No reagrupar handlers por contexto en una sola clase
- No mezclar commands con queries
- Si se agrega funcionalidad nueva, seguir el patrón: dominio -> ports -> repository -> use case -> handler -> router -> DI
