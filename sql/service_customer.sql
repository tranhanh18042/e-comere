CREATE TABLE `customer`(
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `status` INT NOT NULL,
    `username` VARCHAR(255) NOT NULL,
    `password` VARCHAR(255) NOT NULL,
    `first_name` VARCHAR(255) NOT NULL,
    `last_name` INT NOT NULL,
    `address` VARCHAR(255) NOT NULL,
    `phone_number` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL
);
ALTER TABLE
    `customer` ADD PRIMARY KEY `customer_id_primary`(`id`);