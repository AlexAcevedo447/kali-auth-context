# Token Validation Strategy (Eficiente y Escalable)

## Objetivo
Evitar validar tokens por llamada remota en cada request. Cada servicio valida localmente el JWT y decide si atiende o rechaza.

## Patrón recomendado
1. Auth service emite access token JWT corto (5-15 min).
2. Cada microservicio valida firma + claims localmente.
3. Se usa refresh token solo en auth service.
4. Eventos async (Rabbit/Kafka) para revocaciones o cambios críticos.

## Validaciones mínimas en cada servicio
- Firma
- exp / nbf / iat
- iss
- aud
- sub
- tid (tenant)

## Estado actual en kali-auth-context
El middleware JWT está conectado de extremo a extremo:

### Endpoints PÚBLICOS (sin autenticación)
- POST /api/v1/auth/login → genera access_token

### Endpoints PROTEGIDOS (requieren JWT válido)
- POST /api/v1/auth/authorize
- /api/v1/users/* → CRUD de usuarios
- /api/v1/tenants/* → CRUD de tenants
- /api/v1/roles/* → CRUD de roles
- /api/v1/permissions/* → CRUD de permisos
- /api/v1/rbac/* → gestión de roles/permisos por usuario

El middleware valida en cada request protegido:
- Presencia de Authorization: Bearer header
- Firma HS256
- Expiración (exp)
- Emisor (iss)
- Audiencia (aud)
- Claims mínimos: sub (userId), tid (tenantId)

En caso de fallo, responde 401 Unauthorized.

## Implementación en Fiber (middleware reusable)
Ubicación: internal/infrastructure/http/middleware/jwt_auth.go

Uso base en otro servicio:

```go
import "kali-auth-context/internal/infrastructure/http/middleware"

protected := v1.Group("", middleware.NewJWTAuthMiddleware(middleware.JWTAuthConfig{
    Secret:    container.Config.JWTSecret,
    Issuer:    container.Config.JWTIssuer,
    Audience:  container.Config.JWTAudience,
    ClockSkew: 30 * time.Second,
}))

protected.Get("/me", func(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "user_id":  c.Locals("auth.sub"),
        "tenant_id": c.Locals("auth.tid"),
        "email":    c.Locals("auth.email"),
    })
})
```

El middleware almacena datos normalizados en locals:
- auth.sub → userId
- auth.tid → tenantId
- auth.email → email
- auth.claims → mapa completo de claims JWT

## Cuándo usar Rabbit o gRPC
- Rabbit/Kafka: sí para revocación, auditoría, sincronización de permisos y sesiones.
- gRPC introspection: solo para rutas de alto riesgo donde quieras validación online obligatoria.

## Checklist de operación
- JWT_SECRET robusto en producción (no usar "change-me-dev-secret").
- Rotación de secreto o migración a firma asimétrica (RS256/ES256) en evolución.
- TTL corto de access token (5-15 min recomendado).
- Logs con subject y tenant para trazabilidad.
- Rechazar tokens sin claims requeridos.
- Actualizar la estrategia si se agrega refresh token o multi-tenancy compleja.

