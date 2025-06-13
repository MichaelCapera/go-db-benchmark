# Benchmark de MySQL vs PostgreSQL con Go

Este repositorio contiene pruebas de rendimiento entre MySQL y PostgreSQL para inserciones masivas y concurrentes, usando Go.

## Objetivo

Medir y comparar el tiempo de inserción de:
- 100,000 usuarios (secuencial)
- 10,000 transacciones (concurrente con 5 workers)

## Resultados (local, Ryzen 5, SSD, 32 GB RAM)

| Operación | MySQL | PostgreSQL |
|----------|-------|------------|
| Insertar 100k usuarios | 15.27s | 9.18s |
| 10k transacciones concurrentes | 591ms | 489ms |

## Requisitos

- Go >= 1.18
- MySQL y PostgreSQL instalados localmente

## Cómo ejecutar

```bash
go run main.go


# go-db-benchmark

> Benchmarking insert performance in MySQL vs PostgreSQL using Go with concurrent transactions.

---

## 📋 Overview

This project compares the insert performance of **MySQL** and **PostgreSQL** using Go. It simulates both **sequential user inserts** and **concurrent transaction inserts**, using Goroutines to replicate real-world load scenarios.

The goal is to evaluate which database performs better under typical backend workloads, especially in user-driven applications like CRMs, payment processors, or transactional platforms.

---

## 🚀 Technologies

- **Go (Golang)**
- **MySQL** 8+
- **PostgreSQL** 13+
- **Go modules**

---

## 📁 Structure

```text
.
├── main.go                 # Benchmarking code
├── go.mod / go.sum        # Go dependencies
├── schema_mysql.sql       # MySQL table creation
├── schema_postgres.sql    # PostgreSQL table creation
├── results.md             # Performance output and comparison
└── README.md              # This file
