CREATE TABLE `order`(
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
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
ALTER TABLE
    `order` ADD PRIMARY KEY `order_id_primary`(`id`);
CREATE TABLE `order_log`(
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `order_id` INT NOT NULL,
    `status` INT NOT NULL,
    `event_time` TIMESTAMP NOT NULL
);
ALTER TABLE
    `order_log` ADD PRIMARY KEY `order_log_id_primary`(`id`);
ALTER TABLE
    `order_log` ADD CONSTRAINT `order_log_order_id_foreign` FOREIGN KEY(`order_id`) REFERENCES `order`(`id`);