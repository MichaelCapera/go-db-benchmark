// package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"

// 	_ "github.com/go-sql-driver/mysql"
// 	_ "github.com/lib/pq"
// )

// func testConnection(driver, dsn string) {
// 	db, err := sql.Open(driver, dsn)
// 	if err != nil {
// 		log.Fatalf("[%s] Error al abrir conexión: %v", driver, err)
// 	}
// 	defer db.Close()

// 	err = db.Ping()
// 	if err != nil {
// 		log.Fatalf("[%s] No se pudo conectar: %v", driver, err)
// 	} else {
// 		fmt.Printf("[%s] Conexión exitosa!\n", driver)
// 	}
// }

// func main() {
// 	// Cambia los datos por los que usas localmente
// 	mysqlDSN := "root:1234@tcp(localhost:3306)/benchmarkdb"
// 	pgDSN := "host=localhost port=5432 user=postgres password=1234 dbname=benchmarkdb sslmode=disable"

// 	testConnection("mysql", mysqlDSN)
// 	testConnection("postgres", pgDSN)
// }
