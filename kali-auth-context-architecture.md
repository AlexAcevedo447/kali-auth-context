# 🧱 Kali Auth Context --- Arquitectura del Microservicio Identity

## 📌 Descripción General

Este proyecto forma parte de un ERP modular construido bajo arquitectura
de microservicios. El microservicio actual corresponde al **Identity
Service**, encargado de autenticación, autorización y gestión de
identidad.

El sistema inicia como un servicio de facturación, pero está diseñado
para escalar a un ERP completo.

------------------------------------------------------------------------

## 🧠 Principios Arquitectónicos

-   SOLID
-   Domain Driven Design (DDD)
-   Arquitectura Hexagonal (Ports & Adapters)
-   CQRS (Command Query Responsibility Segregation)
-   Clean Architecture mindset
-   Dependency Injection Wire (Go way)

------------------------------------------------------------------------

## ⚙️ Tecnologías Seleccionadas

### Backend

-   **Go (Golang)** → lenguaje principal
-   **Fiber** → framework HTTP
-   **pgx + pgxpool** → driver PostgreSQL de alto rendimiento
-   **PostgreSQL** → base de datos principal

### Calidad de Código

-   gofumpt → formateo
-   goimports → organización de imports
-   golangci-lint → linting profesional

------------------------------------------------------------------------

## 🗂️ Estructura Base del Proyecto


    cmd/
      api/
        main.go   # composition root

    internal/
      identity/
        domain/
        application/
          ports/
        infrastructure/
          postgres/
          http/

------------------------------------------------------------------------

## 🔌 Arquitectura Hexagonal

### Dominio

Contiene entidades puras y lógica de negocio sin dependencias externas.

### Puertos (Ports)

Interfaces que definen contratos:

-   UserCommandRepository
-   UserQueryRepository
-   TransactionManager

Separados por CQRS:

-   Commands → escritura
-   Queries → lectura

### Adaptadores (Infrastructure)

Implementaciones concretas:

-   PostgreSQL repositories
-   HTTP handlers (Fiber)

------------------------------------------------------------------------

## 🔄 Flujo de Inicialización (main.go)

1.  Load config (.env)
2.  Create logger
3.  Connect DB (pgxpool)
4.  Crear repositories
5.  Crear usecases
6.  Crear handlers
7.  Registrar rutas
8.  Start server

El main.go actúa como **Composition Root** para la inyección de
dependencias.

------------------------------------------------------------------------

## 🗄️ Configuración

Variables de entorno actuales:

-   APP_PORT
-   DB_HOST
-   DB_USER
-   DB_PASS
-   DB_NAME

Conexión construida mediante DSN keyword format para pgx.

------------------------------------------------------------------------

## 🚀 Estado Actual

✅ Servidor Fiber levantado\
✅ Conexión PostgreSQL funcionando\
✅ Arquitectura base definida\
🔄 Pendiente iniciar capa de dominio (Identity)

------------------------------------------------------------------------

## 🎯 Objetivo del Identity Service

-   Registro de usuarios
-   Login
-   Manejo de roles y permisos
-   Tokens JWT (futuro)
-   Alta concurrencia
-   Preparado para escalabilidad horizontal

------------------------------------------------------------------------

## 🧩 Notas Importantes

-   No se utilizan monolitos.
-   Cada microservicio tendrá su propia base de datos.
-   Se prioriza claridad arquitectónica antes de escalar
    funcionalidades.
-   DI manual para mantener control explícito de dependencias.

------------------------------------------------------------------------

## 📎 Uso de este archivo

Este archivo sirve como **contexto portable** para continuar el
desarrollo del proyecto en futuras sesiones. Puede cargarse nuevamente
para restaurar el estado conceptual del sistema.
