create database if not exists go_microservice;

create user if not exists 'gouser'@'localhost';
grant all privileges on go_microservice.* to 'gouser'@'localhost';
alter user 'gouser'@'localhost' identified by 'secret';

use go_microservice;

CREATE TABLE products (
                          id mediumint not null auto_increment,
                          name varchar(255) not null,
                          price float(3,2) default 0.00,
                          sku char(11) not null,
                          updatedAt datetime default now() not null on update now(),
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

insert into products (id, name, price, sku, updatedAt)
values
    (1, 'Latte', 1.49, '123-456-789', now()),
    (2, 'Espresso', 0.99, '000-000-001', now());
