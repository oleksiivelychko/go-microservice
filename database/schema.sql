USE mysql;

CREATE DATABASE IF NOT EXISTS go_microservice;

DROP USER IF EXISTS 'gouser'@'%';
CREATE USER IF NOT EXISTS 'gouser'@'%';
GRANT ALL PRIVILEGES ON go_microservice.* TO 'gouser'@'%';
ALTER USER 'gouser'@'%' IDENTIFIED WITH mysql_native_password BY 'secret';
FLUSH PRIVILEGES;

USE go_microservice;

CREATE TABLE products (
                          id MEDIUMINT NOT NULL AUTO_INCREMENT,
                          name VARCHAR(255) NOT NULL,
                          price FLOAT(3,2) DEFAULT 0.00,
                          sku CHAR(11) NOT NULL,
                          updatedAt DATETIME DEFAULT now() NOT NULL ON UPDATE now(),
                          PRIMARY KEY (id)
);

DELIMITER $$
CREATE TRIGGER sku_check BEFORE INSERT ON products
    FOR EACH ROW
BEGIN
    IF (NEW.sku REGEXP '^([0-9]{3})+-([0-9]{3})+-([0-9]{3})$' ) = 0 THEN
        SIGNAL SQLSTATE '01000'
            SET MESSAGE_TEXT = 'SKU has wrong format!';
    END IF;
END$$
DELIMITER;

INSERT INTO products (id, name, price, sku, updatedAt)
VALUES
    (1, 'Latte', 1.49, '123-456-789', now()),
    (2, 'Espresso', 0.99, '000-000-001', now());
