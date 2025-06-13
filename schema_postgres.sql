CREATE TABLE usuarios (
    id INT PRIMARY KEY,
    nombre VARCHAR(100)
);

CREATE TABLE transacciones (
    id SERIAL PRIMARY KEY,
    usuario_id INT,
    monto DECIMAL(10,2),
    fecha TIMESTAMP,
    estado VARCHAR(50),
    FOREIGN KEY (usuario_id) REFERENCES usuarios(id)
);

