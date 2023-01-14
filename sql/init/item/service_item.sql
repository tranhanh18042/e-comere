CREATE DATABASE IF NOT EXISTS `service_item`;

USE `service_item`;

CREATE TABLE IF NOT EXISTS `provider`(
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `provider_name` VARCHAR(255) NOT NULL,
    `phone_number` VARCHAR(255) NOT NULL,
    `address` VARCHAR(255) NOT NULL
) ENGINE=INNODB;

CREATE TABLE IF NOT EXISTS `warehouse`(
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `warehouse_name` VARCHAR(255) NOT NULL,
    `address` VARCHAR(255) NOT NULL,
    `phone_number` VARCHAR(255) NOT NULL
) ENGINE=INNODB;

CREATE TABLE IF NOT EXISTS `item`(
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `warehouse_id` INT UNSIGNED NOT NULL,
    `provider_id` INT UNSIGNED NOT NULL,
    `quantity` INT NOT NULL,
    `status` INT NOT NULL,
    `item_name` VARCHAR(255) NOT NULL,
    `unit_price` INT NOT NULL,
    `description` VARCHAR(255) NOT NULL,
    FOREIGN KEY (provider_id)
        REFERENCES `provider`(id)
        ON DELETE CASCADE,
    FOREIGN KEY (warehouse_id)
        REFERENCES `warehouse`(id)
        ON DELETE CASCADE
) ENGINE=INNODB;
