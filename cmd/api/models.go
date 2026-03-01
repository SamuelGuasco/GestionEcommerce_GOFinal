package main

import "time"

type User struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
	Activo bool   `json:"activo"`
}

type Product struct {
	ID     int     `json:"id"`
	Nombre string  `json:"nombre"`
	Precio float64 `json:"precio"`
}

// Este lo puedes usar para respuestas (GET), joins, etc.
type OrderItem struct {
	ProductoID int     `json:"producto_id"`
	Cantidad   int     `json:"cantidad"`
	Precio     float64 `json:"precio"` // (puede representar precio_unitario o precio, según tu lógica)
	Subtotal   float64 `json:"subtotal"`
}

type Order struct {
	ID        int64       `json:"id"`
	UsuarioID int         `json:"usuario_id"`
	Total     float64     `json:"total"`
	Fecha     *time.Time  `json:"fecha,omitempty"`
	Items     []OrderItem `json:"items,omitempty"`
}

// ---------- Requests (para POST) ----------

type CreateUserRequest struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
	Activo bool   `json:"activo"`
}

type CreateProductRequest struct {
	ID     int     `json:"id"`
	Nombre string  `json:"nombre"`
	Precio float64 `json:"precio"`
}

// ✅ Item SOLO para crear pedidos (POST)
// Importante: el front manda precio_unitario, por eso este json tag.
type CreateOrderItemRequest struct {
	ProductoID     int     `json:"producto_id"`
	Cantidad       int     `json:"cantidad"`
	PrecioUnitario float64 `json:"precio_unitario"`
}

type CreateOrderRequest struct {
	UsuarioID int                      `json:"usuario_id"`
	Items     []CreateOrderItemRequest `json:"items"`
}
