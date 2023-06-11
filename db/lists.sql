DROP TABLE IF EXISTS lists;
DROP TABLE IF EXISTS list_items;

CREATE TABLE list_items (
    id INT AUTO_INCREMENT NOT NULL,
    name VARCHAR(255) NOT NULL,
    list_id INT NOT NULL,
    PRIMARY KEY (`id`)
);

CREATE TABLE lists (
    id INT AUTO_INCREMENT NOT NULL,
    name VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`)
);

CREATE UNIQUE INDEX list_id_to_id ON list_items (list_id, id);
