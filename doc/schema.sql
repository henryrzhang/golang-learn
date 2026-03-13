-- PostgreSQL 建表语句（从 MySQL 迁移）
-- 用户表
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

COMMENT ON TABLE users IS '用户表';
COMMENT ON COLUMN users.id IS '主键 Id';
COMMENT ON COLUMN users.name IS '用户名';
COMMENT ON COLUMN users.email IS '邮箱';
COMMENT ON COLUMN users.phone IS '手机号';
COMMENT ON COLUMN users.password IS '密码(bcrypt加密)';
COMMENT ON COLUMN users.status IS '状态: 1-正常, 0-禁用';
COMMENT ON COLUMN users.created_at IS '创建时间';
COMMENT ON COLUMN users.updated_at IS '更新时间';

-- 自动更新 updated_at 的触发器
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_users_updated_at ON users;
CREATE TRIGGER trigger_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE PROCEDURE update_updated_at_column();

-- 剧集主表
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

COMMENT ON TABLE drama_info IS '剧集主表';
COMMENT ON COLUMN drama_info.id IS '主键 Id';
COMMENT ON COLUMN drama_info.drama_no IS '剧集编号';
COMMENT ON COLUMN drama_info.title IS '剧集标题';
COMMENT ON COLUMN drama_info.outline IS '剧集大纲';
COMMENT ON COLUMN drama_info.cover_image IS '封面图URL';
COMMENT ON COLUMN drama_info.characters IS '剧本主角id列表，逗号分隔';
COMMENT ON COLUMN drama_info.character_relation_desc IS '角色关系描述';
COMMENT ON COLUMN drama_info.status IS '状态 1-草稿 2-提交处理中 3-成功 4-失败';
COMMENT ON COLUMN drama_info.task_no IS '关联的任务编号';
COMMENT ON COLUMN drama_info.create_by IS '创建人';
COMMENT ON COLUMN drama_info.update_by IS '更新人';
COMMENT ON COLUMN drama_info.create_at IS '创建时间';
COMMENT ON COLUMN drama_info.update_at IS '更新时间';
COMMENT ON COLUMN drama_info.deleted IS '是否删除: false-否，true-是';
