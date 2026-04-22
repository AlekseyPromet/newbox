-- Migration: 001_create_tokens_table.sql
-- ============================================
-- Enum для версий токенов
-- ============================================

CREATE TYPE users_token_version AS ENUM ('v1', 'v2');

-- ============================================
-- Таблица токенов пользователей
-- ============================================

CREATE TABLE users_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version users_token_version NOT NULL DEFAULT 'v2',
    user_id UUID NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
    description VARCHAR(200) DEFAULT '',
    created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    expires TIMESTAMP WITH TIME ZONE,
    last_used TIMESTAMP WITH TIME ZONE,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    write_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    
    -- Поля для v1 токенов (plaintext)
    plaintext VARCHAR(40) UNIQUE,
    
    -- Поля для v2 токенов
    key VARCHAR(12) UNIQUE,
    pepper_id SMALLINT,
    hmac_digest VARCHAR(64),
    
    -- Allowed IP addresses (CIDR notation)
    allowed_ips INET[]
);

-- Индексы
CREATE INDEX idx_users_tokens_user_id ON users_tokens(user_id);
CREATE INDEX idx_users_tokens_created ON users_tokens(created DESC);
CREATE INDEX idx_users_tokens_expires ON users_tokens(expires);
CREATE INDEX idx_users_tokens_enabled ON users_tokens(enabled);
CREATE INDEX idx_users_tokens_key ON users_tokens(key) WHERE key IS NOT NULL;
CREATE INDEX idx_users_tokens_plaintext ON users_tokens(plaintext) WHERE plaintext IS NOT NULL;

-- Комментарий
COMMENT ON TABLE users_tokens IS 'API tokens for user authentication';
COMMENT ON COLUMN users_tokens.version IS 'Token version: v1 (legacy) or v2 (secure)';
COMMENT ON COLUMN users_tokens.plaintext IS 'Plaintext token value for v1 tokens only';
COMMENT ON COLUMN users_tokens.key IS 'Public key identifier for v2 tokens';
COMMENT ON COLUMN users_tokens.pepper_id IS 'ID of the pepper used for HMAC hashing (v2 only)';
COMMENT ON COLUMN users_tokens.hmac_digest IS 'SHA256 HMAC digest of token (v2 only)';
COMMENT ON COLUMN users_tokens.allowed_ips IS 'List of allowed client IP networks (CIDR)';

-- ============================================
-- Проверка целостности для версий токенов
-- ============================================

ALTER TABLE users_tokens ADD CONSTRAINT enforce_version_dependent_fields CHECK (
    (version = 'v1' AND 
     key IS NULL AND 
     pepper_id IS NULL AND 
     hmac_digest IS NULL AND 
     plaintext IS NOT NULL)
    OR
    (version = 'v2' AND 
     key IS NOT NULL AND 
     pepper_id IS NOT NULL AND 
     hmac_digest IS NOT NULL AND 
     plaintext IS NULL)
);

-- Migration: 002_create_owners_tables.sql
-- ============================================
-- Таблица групп владельцев (OwnerGroup)
-- ============================================

CREATE TABLE users_owner_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description VARCHAR(200) DEFAULT '',
    created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Индексы
CREATE INDEX idx_users_owner_groups_name ON users_owner_groups(name);

-- Комментарии
COMMENT ON TABLE users_owner_groups IS 'Groups of owners for object ownership';
COMMENT ON COLUMN users_owner_groups.name IS 'Unique name of the owner group';
COMMENT ON COLUMN users_owner_groups.description IS 'Description of the owner group';

-- ============================================
-- Таблица владельцев (Owner)
-- ============================================

CREATE TABLE users_owners (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description VARCHAR(200) DEFAULT '',
    group_id UUID REFERENCES users_owner_groups(id) ON DELETE SET NULL,
    created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Индексы
CREATE INDEX idx_users_owners_name ON users_owners(name);
CREATE INDEX idx_users_owners_group_id ON users_owners(group_id);

-- Комментарии
COMMENT ON TABLE users_owners IS 'Owners for object ownership (users or groups)';
COMMENT ON COLUMN users_owners.name IS 'Unique name of the owner';
COMMENT ON COLUMN users_owners.description IS 'Description of the owner';
COMMENT ON COLUMN users_owners.group_id IS 'Reference to owner group (if this is a group owner)';

-- ============================================
-- Связь владельцев с пользователями (M2M)
-- ============================================

CREATE TABLE users_owners_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL REFERENCES users_owners(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
    created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT users_owners_users_unique UNIQUE (owner_id, user_id)
);

-- Индексы
CREATE INDEX idx_users_owners_users_owner_id ON users_owners_users(owner_id);
CREATE INDEX idx_users_owners_users_user_id ON users_owners_users(user_id);

-- Комментарии
COMMENT ON TABLE users_owners_users IS 'Many-to-many relationship between owners and users';

-- ============================================
-- Связь владельцев с группами пользователей (M2M)
-- ============================================

CREATE TABLE users_owners_user_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL REFERENCES users_owners(id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES users_groups(id) ON DELETE CASCADE,
    created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT users_owners_user_groups_unique UNIQUE (owner_id, group_id)
);

-- Индексы
CREATE INDEX idx_users_owners_user_groups_owner_id ON users_owners_user_groups(owner_id);
CREATE INDEX idx_users_owners_user_groups_group_id ON users_owners_user_groups(group_id);

-- Комментарии
COMMENT ON TABLE users_owners_user_groups IS 'Many-to-many relationship between owners and user groups';
