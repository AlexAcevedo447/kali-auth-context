# 🧱 Kali Auth Context --- Arquitectura del Microservicio Identity

## 📌 Descripción General

Este proyecto forma parte de un ERP modular construido bajo arquitectura
de microservicios. El microservicio actual corresponde al **Identity
Service**, encargado de autenticación, autorización y gestión de
identidad.

El sistema inicia como un servicio de facturación, pero está diseñado
para escalar a un ERP completo, priorizando rendimiento, claridad
arquitectónica y mantenibilidad.

------------------------------------------------------------------------

## 🧠 Principios Arquitectónicos

-   SOLID\
-   Domain Driven Design (DDD)\
-   Arquitectura Hexagonal (Ports & Adapters)\
-   CQRS (Command Query Responsibility Segregation)\
-   Clean Architecture mindset\
-   Dependency Injection con Wire\
-   Separación estricta de responsabilidades\
-   Production-ready mindset (graceful shutdown + limpieza de recursos)

------------------------------------------------------------------------

## ⚙️ Tecnologías Seleccionadas

### Backend

-   **Go (Golang)** → lenguaje principal\
-   **Fiber** → framework HTTP\
-   **Bun ORM** → ORM sobre `database/sql`\
-   **PostgreSQL** → base de datos principal

> Se elimina `pgxpool` directo y se utiliza Bun como capa ORM escalable,
> manteniendo alto rendimiento y mejor experiencia de desarrollo.

------------------------------------------------------------------------

## 🗂️ Estructura Base del Proyecto

    cmd/
      api/
        main.go   # Composition Root

    internal/
      domain/
        identity/

      application/
        commands/
        queries/
        ports/

      infrastructure/
        config/
        db/
        persistence/
          postgres/
        http/
        logger/

------------------------------------------------------------------------

## 🔌 Arquitectura Hexagonal

### Dominio

-   Entidades puras (`User`)
-   Value Objects (`UserId`)
-   Sin dependencias externas
-   Sin lógica técnica (hash, DB, HTTP)

------------------------------------------------------------------------

### Application Layer

#### Commands (escritura)

-   1 command por caso de uso
-   1 repository por interfaz
-   DTOs separados del dominio
-   Validaciones y hashing en la capa de aplicación

#### Queries (lectura)

-   Separadas de commands (CQRS real)
-   1 repository por interfaz
-   Orientadas a lectura eficiente

------------------------------------------------------------------------

### Ports

Cada responsabilidad tiene su propio contrato:

#### Commands

-   `ICreateUserCommandRepository`
-   `IUpdateUserCommandRepository`
-   `IDeleteUserCommandRepository`

#### Queries

-   `IGetUserByIdQueryRepository`
-   `IGetUserByEmailQueryRepository`
-   `IListUsersQueryRepository`

#### Servicios transversales

-   `IPasswordHasher`
-   `IUUIDProvider`

No se combinan responsabilidades en una misma implementación.

------------------------------------------------------------------------

### Infrastructure

Implementaciones concretas:

-   Repositories con Bun
-   HTTP handlers (Fiber)
-   Logger estructurado
-   Config loader
-   DB connection factory

------------------------------------------------------------------------

## 🔄 Flujo de Inicialización (main.go)

El `main.go` actúa como **Composition Root**.

Flujo actual:

1.  Load config (.env)
2.  Crear logger
3.  Crear conexión Bun (`*bun.DB`)
4.  Registrar `defer` para limpieza
5.  Inyectar dependencias (Wire)
6.  Crear HTTP server
7.  Start server
8.  Graceful shutdown
9.  Cierre ordenado de recursos

La conexión a base de datos:

-   Se crea una vez
-   Vive durante todo el proceso
-   Se cierra cuando el servicio termina

Nunca se abre ni se cierra por request.

------------------------------------------------------------------------

## 🔐 Seguridad y Futuro

-   Hashing de contraseñas vía `IPasswordHasher`
-   JWT validation externa (Keycloak)
-   Preparado para roles y permisos
-   Alta concurrencia
-   Preparado para escalabilidad horizontal

------------------------------------------------------------------------

## 🚀 Estado Actual

✅ Arquitectura hexagonal definida\
✅ CQRS real separado\
✅ Bun integrado como ORM\
✅ Connection lifecycle controlado desde main\
✅ Limpieza de recursos con `defer`\
🔄 Implementación completa de repositories en progreso

------------------------------------------------------------------------

## 🧩 Decisiones Clave

-   No monolito.
-   Cada microservicio tiene su propia base de datos.
-   No mezclar commands y queries.
-   No abrir/cerrar DB por request.
-   No lógica técnica en dominio.
-   Todo recurso externo debe cerrarse correctamente.
-   Arquitectura pensada para entorno crítico (identity + billing).

------------------------------------------------------------------------

## 📎 Uso de este archivo

Este archivo sirve como **contexto portable actualizado** para continuar
el desarrollo en futuras sesiones.

Puede cargarse nuevamente para restaurar el estado conceptual y técnico
del sistema sin perder coherencia arquitectónica.
