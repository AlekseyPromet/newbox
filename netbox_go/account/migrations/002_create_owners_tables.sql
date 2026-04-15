-- Migration: 002_create_owners_tables.sql
-- Source: netbox/users/migrations/0015_owner.py
-- Description: Creates Owner and OwnerGroup tables for object ownership
-- Date: 2024-01-15

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
