-- sqlc 使用的 schema（与 doc/schema.sql 保持一致，不含触发器）
CREATE TABLE IF NOT EXISTS users (
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(100) NOT NULL DEFAULT '',
    email      VARCHAR(255) NOT NULL DEFAULT '',
    phone      VARCHAR(20)  DEFAULT '',
    password   VARCHAR(255) NOT NULL DEFAULT '',
    status     SMALLINT     NOT NULL DEFAULT 1,
    created_at TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uniq_email UNIQUE (email)
);

CREATE TABLE IF NOT EXISTS drama_info (
    id                      BIGSERIAL PRIMARY KEY,
    drama_no                VARCHAR(32)   NOT NULL DEFAULT '',
    title                   VARCHAR(200)  NOT NULL DEFAULT '',
    outline                 TEXT,
    cover_image             VARCHAR(500)  DEFAULT '',
    characters              VARCHAR(1000) NOT NULL DEFAULT '',
    character_relation_desc VARCHAR(1000) NOT NULL DEFAULT '',
    status                  SMALLINT      DEFAULT 1,
    task_no                 VARCHAR(64)  NOT NULL DEFAULT '',
    create_by               VARCHAR(64)  NOT NULL DEFAULT '',
    update_by               VARCHAR(64)  NOT NULL DEFAULT '',
    create_at               BIGINT       NOT NULL DEFAULT 0,
    update_at               BIGINT       NOT NULL DEFAULT 0,
    deleted                 BOOLEAN      NOT NULL DEFAULT FALSE,
    CONSTRAINT uniq_drama_no UNIQUE (drama_no)
);
