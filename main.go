package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func cleanTablesMySQL(db *sql.DB) error {
	_, err := db.Exec("SET FOREIGN_KEY_CHECKS=0;")
	if err != nil {
		return err
	}
	_, err = db.Exec("TRUNCATE TABLE transacciones;")
	if err != nil {
		return err
	}
	_, err = db.Exec("TRUNCATE TABLE usuarios;")
	if err != nil {
		return err
	}
	_, err = db.Exec("SET FOREIGN_KEY_CHECKS=1;")
	return err
}

func cleanTablesPostgres(db *sql.DB) error {
	_, err := db.Exec("TRUNCATE TABLE transacciones, usuarios RESTART IDENTITY CASCADE;")
	return err
}

func insertDummyUsersMySQL(db *sql.DB, totalUsers int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO usuarios(id, nombre) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i := 0; i < totalUsers; i++ {
		_, err := stmt.Exec(i, fmt.Sprintf("Usuario %d", i))
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func insertDummyUsersPostgres(db *sql.DB, totalUsers int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO usuarios(id, nombre) VALUES ($1, $2)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i := 0; i < totalUsers; i++ {
		_, err := stmt.Exec(i, fmt.Sprintf("Usuario %d", i))
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func insertConcurrent(db *sql.DB, engine string, workers int, insertsPerWorker int) {
	done := make(chan bool)
	start := time.Now()

	for w := 0; w < workers; w++ {
		go func(workerID int) {
			tx, err := db.Begin()
			if err != nil {
				log.Fatalf("[%s] Worker %d: error al iniciar tx: %v", engine, workerID, err)
			}
			var stmt *sql.Stmt
			if engine == "mysql" {
				stmt, err = tx.Prepare("INSERT INTO transacciones(usuario_id, monto, fecha, estado) VALUES (?, ?, NOW(), ?)")
			} else {
				stmt, err = tx.Prepare("INSERT INTO transacciones(usuario_id, monto, fecha, estado) VALUES ($1, $2, NOW(), $3)")
			}
			if err != nil {
				log.Fatalf("[%s] Worker %d: error al preparar stmt: %v", engine, workerID, err)
			}
			defer stmt.Close()

			for i := 0; i < insertsPerWorker; i++ {
				_, err := stmt.Exec(i%100000, 100.0, "pendiente")
				if err != nil {
					log.Fatalf("[%s] Worker %d: error en exec: %v", engine, workerID, err)
				}
			}
			if err := tx.Commit(); err != nil {
				log.Fatalf("[%s] Worker %d: error al commit: %v", engine, workerID, err)
			}
			done <- true
		}(w)
	}

	for i := 0; i < workers; i++ {
		<-done
	}

	elapsed := time.Since(start)
	totalInserts := workers * insertsPerWorker
	fmt.Printf("[%s] Insertadas %d filas en %s (%.2f inserts/segundo)\n",
		engine, totalInserts, elapsed, float64(totalInserts)/elapsed.Seconds())
}

func main() {
	mysqlDSN := "root:1234@tcp(localhost:3306)/benchmarkdb"
	pgDSN := "host=localhost port=5432 user=postgres password=1234 dbname=benchmarkdb sslmode=disable"

	dbMySQL, err := sql.Open("mysql", mysqlDSN)
	if err != nil {
		log.Fatal("Error al conectar MySQL:", err)
	}
	defer dbMySQL.Close()

	dbPG, err := sql.Open("postgres", pgDSN)
	if err != nil {
		log.Fatal("Error al conectar PostgreSQL:", err)
	}
	defer dbPG.Close()

	totalUsers := 100000

	fmt.Println("Limpiando tablas en MySQL...")
	if err := cleanTablesMySQL(dbMySQL); err != nil {
		log.Fatal("Error limpiando tablas MySQL:", err)
	}

	fmt.Println("Limpiando tablas en PostgreSQL...")
	if err := cleanTablesPostgres(dbPG); err != nil {
		log.Fatal("Error limpiando tablas PostgreSQL:", err)
	}

	fmt.Println("Insertando usuarios en MySQL...")
	startUsersMySQL := time.Now()
	err = insertDummyUsersMySQL(dbMySQL, totalUsers)
	if err != nil {
		log.Fatal("Error insertando usuarios en MySQL:", err)
	}
	fmt.Println("Usuarios insertados en MySQL en:", time.Since(startUsersMySQL))

	fmt.Println("Insertando usuarios en PostgreSQL...")
	startUsersPG := time.Now()
	err = insertDummyUsersPostgres(dbPG, totalUsers)
	if err != nil {
		log.Fatal("Error insertando usuarios en PostgreSQL:", err)
	}
	fmt.Println("Usuarios insertados en PostgreSQL en:", time.Since(startUsersPG))

	fmt.Println("Insertando transacciones concurrentes en MySQL...")
	insertConcurrent(dbMySQL, "mysql", 5, 2000)

	fmt.Println("Insertando transacciones concurrentes en PostgreSQL...")
	insertConcurrent(dbPG, "postgresql", 5, 2000)
}
