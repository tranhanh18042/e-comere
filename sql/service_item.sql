CREATE TABLE `provider`(
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `provider_name` VARCHAR(255) NOT NULL,
    `phone_number` VARCHAR(255) NOT NULL,
    `address` VARCHAR(255) NOT NULL
);
ALTER TABLE
    `provider` ADD PRIMARY KEY `provider_id_primary`(`id`);
CREATE TABLE `warehouse`(
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `warehouse_name` VARCHAR(255) NOT NULL,
    `address` VARCHAR(255) NOT NULL,
    `phone_number` VARCHAR(255) NOT NULL
);
ALTER TABLE
    `warehouse` ADD PRIMARY KEY `warehouse_id_primary`(`id`);
CREATE TABLE `item`(
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `warehouse_id` INT NOT NULL,
    `provider_id` INT NOT NULL,
    `quantity` INT NOT NULL,
    `status` INT NOT NULL,
    `item_name` VARCHAR(255) NOT NULL,
    `unit_price` INT NOT NULL,
    `description` VARCHAR(255) NOT NULL
);
ALTER TABLE
    `item` ADD PRIMARY KEY `item_id_primary`(`id`);
ALTER TABLE
    `item` ADD CONSTRAINT `item_provider_id_foreign` FOREIGN KEY(`provider_id`) REFERENCES `provider`(`id`);
ALTER TABLE
    `item` ADD CONSTRAINT `item_warehouse_id_foreign` FOREIGN KEY(`warehouse_id`) REFERENCES `warehouse`(`id`);