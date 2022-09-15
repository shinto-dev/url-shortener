CREATE TABLE unique_ids_64
(
    `id`   bigint(20) unsigned NOT NULL auto_increment,
    `stub` char(1)             NOT NULL default '',
    PRIMARY KEY (`id`),
    UNIQUE KEY `stub` (`stub`)
);

INSERT INTO unique_ids_64 (id, stub)
VALUES (1000000000, 's');
