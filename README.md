# Sistema de Gestión Ecommerce (Go + MySQL + JSON + Web)

Proyecto Final – Programación Orientada a Objetos  
Universidad Internacional del Ecuador  

Autor: Samuel Nathaniel Guasco Cedeño  
Carrera: Ingeniería en Software  
Semestre: Tercer Semestre  

---

## 1. Descripción General

El presente proyecto consiste en el desarrollo de un Sistema de Gestión Ecommerce implementado en Go (Golang) como backend, utilizando MySQL como sistema gestor de base de datos relacional y JSON como mecanismo de serialización para la comunicación entre cliente y servidor.

El sistema permite gestionar usuarios, productos y pedidos, así como calcular automáticamente los totales de cada pedido. Además, se desarrolló un panel web interactivo que consume los servicios REST del backend.

El proyecto integra conceptos fundamentales de programación orientada a objetos, servicios web, persistencia en bases de datos y arquitectura por capas.

---

## 2. Arquitectura del Proyecto

El sistema está estructurado de forma modular, separando responsabilidades en diferentes archivos y capas:

GestionEcommerceAPI/

cmd/api/
- main.go         → Punto de entrada del servidor
- handlers.go     → Controladores HTTP
- database.go     → Conexión a MySQL
- models.go       → Modelos y estructuras de datos
- queries.go      → Consultas SQL y transacciones
- utils.go        → Funciones auxiliares

web/
- index.html      → Panel web (frontend)

SqlDatabase.sql   → Script completo para crear la base de datos
go.mod
go.sum

Esta estructura permite una clara separación entre la lógica de negocio, el acceso a datos y la capa de presentación.

---

## 3. Base de Datos

Nombre del schema:
ecommerce_db

Tablas implementadas:
- usuarios
- productos
- pedidos
- pedido_detalles

Relaciones:
- Un usuario puede tener múltiples pedidos.
- Un pedido puede tener múltiples registros en pedido_detalles.
- Se utilizan llaves foráneas para garantizar integridad referencial.

Transacciones:
El registro de un pedido se realiza mediante una transacción SQL que garantiza que:
- Se inserta el pedido.
- Se insertan todos los detalles.
- Se calcula el total.
- Se confirma con COMMIT.
- En caso de error, se ejecuta ROLLBACK.

Esto asegura consistencia en la base de datos.

---

## 4. Servicios Web Implementados

Todos los servicios utilizan JSON como formato de intercambio de datos.

GET  /api/ping  
Verifica que la API esté funcionando correctamente.

GET  /api/usuarios  
Lista todos los usuarios registrados.

POST /api/usuarios  
Crea un nuevo usuario.

GET  /api/productos  
Lista todos los productos registrados.

POST /api/productos  
Crea un nuevo producto.

GET  /api/pedidos  
Lista todos los pedidos.

POST /api/pedidos  
Crea un pedido con múltiples ítems.

GET  /api/resumen  
Devuelve estadísticas generales del sistema:
- Total de usuarios
- Total de productos
- Total de pedidos
- Total vendido

Se implementaron más de ocho servicios web, cumpliendo con los requisitos académicos.

---

## 5. Lógica de Negocio

Cálculo de subtotal por ítem:
subtotal = precio_unitario * cantidad

Cálculo de total del pedido:
total = suma de todos los subtotales

El total se almacena en la tabla pedidos y se calcula automáticamente durante la transacción de inserción.

---

## 6. Panel Web

Ubicación:
web/index.html

El panel web permite:

- Visualizar totales generales del sistema.
- Crear usuarios.
- Crear productos.
- Crear pedidos con múltiples ítems en formato JSON.
- Listar usuarios, productos y pedidos.
- Ver el total vendido acumulado.

El frontend se comunica con el backend mediante llamadas fetch utilizando JSON.

---

## 7. Instalación y Ejecución

1. Clonar el repositorio:

git clone https://github.com/SamuelGuasco/GestionEcommerce_GOFinal.git
cd GestionEcommerce_GOFinal

2. Importar la base de datos en MySQL Workbench:

- Ir a Server → Data Import.
- Seleccionar “Import from Self-Contained File”.
- Elegir el archivo SqlDatabase.sql.
- Crear o seleccionar el schema ecommerce_db.
- Ejecutar Start Import.

3. Configurar credenciales en:

cmd/api/database.go

Verificar:
- Host
- Puerto
- Usuario
- Contraseña
- Nombre de la base de datos (ecommerce_db)

4. Ejecutar el servidor:

go run ./cmd/api

Debe mostrarse en consola que la conexión a MySQL es exitosa y que el servidor está corriendo en:

http://localhost:8080

5. Abrir el navegador:

http://localhost:8080

---

## 8. Ejemplos de JSON

Crear usuario:

{
  "id": 10,
  "nombre": "Carlos",
  "email": "carlos@mail.com",
  "activo": true
}

Crear producto:

{
  "id": 1,
  "nombre": "Laptop",
  "precio": 1200
}

Crear pedido:

{
  "usuario_id": 1,
  "items": [
    { "producto_id": 1, "cantidad": 1, "precio_unitario": 1300 },
    { "producto_id": 2, "cantidad": 2, "precio_unitario": 25.50 }
  ]
}

---

## 9. Conceptos Aplicados

- Programación Orientada a Objetos en Go.
- Servicios Web REST.
- Serialización y deserialización con JSON.
- Persistencia en base de datos relacional.
- Transacciones SQL.
- Separación por capas.
- Arquitectura modular.
- Consumo de API desde frontend web.

---

## 10. Cumplimiento de Requisitos

- Proyecto completo en repositorio GitHub.
- Servicios web implementados.
- Serialización mediante JSON.
- Base de datos SQL.
- Panel web funcional.
- Transacciones en inserción de pedidos.
- Integración de conocimientos de múltiples unidades académicas.

---

## 11. Observaciones Finales

El sistema es completamente funcional en entorno local y permite gestionar un flujo completo de ecommerce básico, desde la creación de usuarios y productos hasta el registro de pedidos con cálculo automático de totales.

Este proyecto demuestra la integración entre backend, base de datos y frontend web en un entorno estructurado y profesional.