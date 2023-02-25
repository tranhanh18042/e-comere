CREATE DATABASE IF NOT EXISTS `service_customer`;

USE `service_customer`;

CREATE TABLE IF NOT EXISTS `customer` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `status` INT NOT NULL,
    `username` VARCHAR(255) NOT NULL,
    `password` VARCHAR(255) NOT NULL,
    `first_name` VARCHAR(255) NOT NULL,
    `last_name` VARCHAR(255) NOT NULL,
    `address` VARCHAR(255) NOT NULL,
    `phone_number` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL
) ENGINE=INNODB;
