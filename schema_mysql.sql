CREATE TABLE usuarios (
    id INT PRIMARY KEY,
    nombre VARCHAR(100)
);

CREATE TABLE transacciones (
    id INT AUTO_INCREMENT PRIMARY KEY,
    usuario_id INT,
    monto DECIMAL(10,2),
    fecha DATETIME,
    estado VARCHAR(50),
    FOREIGN KEY (usuario_id) REFERENCES usuarios(id)
);
