-- Migration: 001_create_tokens_table.sql
-- Source: netbox/account/migrations/0001_initial.py + users migrations
-- Description: Creates the tokens table for API authentication (v1 and v2 support)
-- Date: 2024-01-15

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
