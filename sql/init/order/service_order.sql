CREATE DATABASE IF NOT EXISTS `service_order`;

USE `service_order`;

CREATE TABLE IF NOT EXISTS `order`(
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `status` INT NOT NULL,
    `customer_id` INT NOT NULL,
    `item_id` INT NOT NULL,
    `item_quantity` INT NOT NULL,
    `address` VARCHAR(255) NOT NULL,
    `item_amount` INT NOT NULL,
    `ship_fee` INT NOT NULL,
    `total_amount` INT NOT NULL,
    `discount_amount` INT NOT NULL
);
CREATE TABLE IF NOT EXISTS `order_log`(
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `order_id` INT UNSIGNED NOT NULL,
    `status` INT NOT NULL,
    `event_time` TIMESTAMP NOT NULL,
    FOREIGN KEY (order_id)
        REFERENCES `order`(id)
        ON DELETE CASCADE
);
