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
- Login es tenant-aware
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
