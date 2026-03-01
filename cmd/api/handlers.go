package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type api struct {
	db *sql.DB
}

func (a api) ping(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"message": "API funcionando correctamente",
	})
}

// ---- USERS ----
func (a api) users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		list, err := listUsers(a.db)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, list)

	case http.MethodPost:
		var req CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		if req.ID <= 0 || req.Nombre == "" || req.Email == "" {
			writeError(w, http.StatusBadRequest, "Campos requeridos: id, nombre, email")
			return
		}
		if err := createUser(a.db, req); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{"ok": true})

	default:
		writeError(w, http.StatusMethodNotAllowed, "Método no permitido")
	}
}

// ---- PRODUCTS ----
func (a api) products(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		list, err := listProducts(a.db)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, list)

	case http.MethodPost:
		var req CreateProductRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		if req.ID <= 0 || req.Nombre == "" {
			writeError(w, http.StatusBadRequest, "Campos requeridos: id, nombre")
			return
		}
		if err := createProduct(a.db, req); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{"ok": true})

	default:
		writeError(w, http.StatusMethodNotAllowed, "Método no permitido")
	}
}

// ---- ORDERS ----
func (a api) orders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		list, err := listOrders(a.db)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, list)

	case http.MethodPost:
		var req CreateOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}

		// si el precio no viene en el JSON, te falla total.
		// Para hacerlo simple hoy: exigimos precio en cada item.
		id, err := createOrderTx(a.db, req)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{
			"ok":        true,
			"pedido_id": id,
		})

	default:
		writeError(w, http.StatusMethodNotAllowed, "Método no permitido")
	}
}

// ---- PRODUCT BY ID ----
// GET /api/productos/{id}
func (a api) productByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	// r.URL.Path viene tipo: /api/productos/123
	idStr := strings.TrimPrefix(r.URL.Path, "/api/productos/")
	if idStr == "" {
		writeError(w, http.StatusBadRequest, "Falta el id del producto")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "id inválido")
		return
	}

	p, err := getProductByID(a.db, id)
	if err != nil {
		// si no existe, 404
		if err == sql.ErrNoRows {
			writeError(w, http.StatusNotFound, "Producto no encontrado")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, p)
}

// ---- RESUMEN ----
// GET /api/resumen
func (a api) resumen(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	res, err := getResumen(a.db)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, res)
}
