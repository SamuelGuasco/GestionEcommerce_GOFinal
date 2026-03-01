package main

import (
	"database/sql"
	"errors"
)

// ---------- USERS ----------

func listUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query(`SELECT id, nombre, email, activo FROM usuarios ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []User
	for rows.Next() {
		var u User
		var activoInt int
		if err := rows.Scan(&u.ID, &u.Nombre, &u.Email, &activoInt); err != nil {
			return nil, err
		}
		u.Activo = (activoInt == 1)
		out = append(out, u)
	}
	return out, rows.Err()
}

func createUser(db *sql.DB, req CreateUserRequest) error {
	activoInt := 0
	if req.Activo {
		activoInt = 1
	}
	_, err := db.Exec(`
		INSERT INTO usuarios (id, nombre, email, activo)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			nombre = VALUES(nombre),
			email  = VALUES(email),
			activo = VALUES(activo)
	`, req.ID, req.Nombre, req.Email, activoInt)
	return err
}

// ---------- PRODUCTS ----------

func listProducts(db *sql.DB) ([]Product, error) {
	rows, err := db.Query(`SELECT id, nombre, precio FROM productos ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Nombre, &p.Precio); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func createProduct(db *sql.DB, req CreateProductRequest) error {
	_, err := db.Exec(
		"INSERT INTO productos (id, nombre, precio) VALUES (?, ?, ?)",
		req.ID,
		req.Nombre,
		req.Precio,
	)
	return err
}

// ---------- ORDERS ----------

func createOrderTx(db *sql.DB, req CreateOrderRequest) (int64, error) {
	if req.UsuarioID <= 0 {
		return 0, errors.New("usuario_id inválido")
	}
	if len(req.Items) == 0 {
		return 0, errors.New("items no puede estar vacío")
	}

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()

	// 1) calcular total
	var total float64
	for _, it := range req.Items {
		if it.ProductoID <= 0 || it.Cantidad <= 0 {
			return 0, errors.New("producto_id/cantidad inválidos")
		}
		// it.Precio = precio unitario que viene en el JSON
		total += it.PrecioUnitario * float64(it.Cantidad)
	}

	// 2) insertar pedido
	res, err := tx.Exec(`INSERT INTO pedidos (usuario_id, total) VALUES (?, ?)`, req.UsuarioID, total)
	if err != nil {
		return 0, err
	}
	pedidoID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	// 3) insertar detalles
	for _, it := range req.Items {
		subtotal := it.PrecioUnitario * float64(it.Cantidad)

		_, err := tx.Exec(`
			INSERT INTO pedido_detalles (pedido_id, producto_id, cantidad, precio_unitario, subtotal)
			VALUES (?, ?, ?, ?, ?)
		`, pedidoID, it.ProductoID, it.Cantidad, it.PrecioUnitario, subtotal)
		if err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return pedidoID, nil
}

func listOrders(db *sql.DB) ([]Order, error) {
	rows, err := db.Query(`SELECT id, usuario_id, total FROM pedidos ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.UsuarioID, &o.Total); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, rows.Err()
}

// ---------- RESUMEN ----------
type Resumen struct {
	Usuarios     int     `json:"usuarios"`
	Productos    int     `json:"productos"`
	Pedidos      int     `json:"pedidos"`
	TotalVendido float64 `json:"total_vendido"`
}

func getResumen(db *sql.DB) (Resumen, error) {
	var r Resumen

	// COUNT usuarios
	if err := db.QueryRow(`SELECT COUNT(*) FROM usuarios`).Scan(&r.Usuarios); err != nil {
		return r, err
	}
	// COUNT productos
	if err := db.QueryRow(`SELECT COUNT(*) FROM productos`).Scan(&r.Productos); err != nil {
		return r, err
	}
	// COUNT pedidos
	if err := db.QueryRow(`SELECT COUNT(*) FROM pedidos`).Scan(&r.Pedidos); err != nil {
		return r, err
	}
	// SUM total pedidos
	if err := db.QueryRow(`SELECT COALESCE(SUM(total), 0) FROM pedidos`).Scan(&r.TotalVendido); err != nil {
		return r, err
	}

	return r, nil
}

// ---------- PRODUCT BY ID ----------
func getProductByID(db *sql.DB, id int) (Product, error) {
	var p Product
	err := db.QueryRow(`SELECT id, nombre, precio FROM productos WHERE id = ?`, id).
		Scan(&p.ID, &p.Nombre, &p.Precio)
	return p, err
}
