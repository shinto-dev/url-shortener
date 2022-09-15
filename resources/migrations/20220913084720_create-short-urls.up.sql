create table short_urls
(
    id           INT PRIMARY KEY AUTO_INCREMENT,
    created_at   TIMESTAMP,
    updated_at   TIMESTAMP,
    deleted_at   TIMESTAMP,
    original_url VARCHAR(255) NOT NULL,
    short_path   VARCHAR(16)  NOT NULL,
    not_archived BOOLEAN GENERATED ALWAYS AS (IF(deleted_at IS NULL, 1, NULL)) VIRTUAL,
    constraint unique_short_path unique (short_path, not_archived)
);
