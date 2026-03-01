package main

import (
	"database/sql"
	"log"
	"net/http"
)

var db *sql.DB

func pingHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{
		"message": "API funcionando correctamente",
	})
}

func main() {
	db, err := conectarBD()
	if err != nil {
		log.Fatal("Error conectando BD: ", err)
	}
	defer db.Close()
	log.Println("✅ Conexión a MySQL OK")

	a := api{db: db}

	http.HandleFunc("/api/ping", a.ping)
	http.HandleFunc("/api/usuarios", a.users)
	http.HandleFunc("/api/productos", a.products)
	http.HandleFunc("/api/pedidos", a.orders)
	http.HandleFunc("/api/productos/", a.productByID) // ojo el slash final
	http.HandleFunc("/api/resumen", a.resumen)

	// Servir frontend estático (web/index.html)
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	log.Println("🚀 Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
