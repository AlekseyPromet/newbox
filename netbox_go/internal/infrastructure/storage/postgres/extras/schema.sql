-- Extras Module Schema

-- EventRule
CREATE TABLE extras_eventrule (
    id UUID PRIMARY KEY,
    name VARCHAR(150) NOT NULL UNIQUE,
    description VARCHAR(200),
    event_types TEXT[] NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    conditions JSONB,
    action_type VARCHAR(30) NOT NULL,
    action_object_type VARCHAR(255) NOT NULL,
    action_object_id BIGINT,
    action_data JSONB,
    comments TEXT,
    owner_id UUID REFERENCES core_user(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_extras_eventrule_action_object ON extras_eventrule(action_object_type, action_object_id);

-- EventRule Object Types (Many-to-Many)
CREATE TABLE extras_eventrule_object_types (
    event_rule_id UUID REFERENCES extras_eventrule(id) ON DELETE CASCADE,
    content_type VARCHAR(255) NOT NULL,
    PRIMARY KEY (event_rule_id, content_type)
);

-- Webhook
CREATE TABLE extras_webhook (
    id UUID PRIMARY KEY,
    name VARCHAR(150) NOT NULL UNIQUE,
    description VARCHAR(200),
    payload_url VARCHAR(500) NOT NULL,
    http_method VARCHAR(30) NOT NULL,
    http_content_type VARCHAR(100) NOT NULL,
    additional_headers TEXT,
    body_template TEXT,
    secret VARCHAR(255),
    ssl_verification BOOLEAN NOT NULL DEFAULT TRUE,
    ca_file_path VARCHAR(4096),
    owner_id UUID REFERENCES core_user(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- CustomLink
CREATE TABLE extras_customlink (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    link_text TEXT NOT NULL,
    link_url TEXT NOT NULL,
    weight SMALLINT NOT NULL DEFAULT 100,
    group_name VARCHAR(50),
    button_class VARCHAR(30) NOT NULL,
    new_window BOOLEAN NOT NULL DEFAULT FALSE,
    owner_id UUID REFERENCES core_user(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- CustomLink Object Types (Many-to-Many)
CREATE TABLE extras_customlink_object_types (
    custom_link_id UUID REFERENCES extras_customlink(id) ON DELETE CASCADE,
    content_type VARCHAR(255) NOT NULL,
    PRIMARY KEY (custom_link_id, content_type)
);

-- ExportTemplate
CREATE TABLE extras_exporttemplate (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(200),
    template_code TEXT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_extension VARCHAR(10) NOT NULL,
    as_attachment BOOLEAN NOT NULL DEFAULT TRUE,
    owner_id UUID REFERENCES core_user(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ExportTemplate Object Types (Many-to-Many)
CREATE TABLE extras_exporttemplate_object_types (
    export_template_id UUID REFERENCES extras_exporttemplate(id) ON DELETE CASCADE,
    content_type VARCHAR(255) NOT NULL,
    PRIMARY KEY (export_template_id, content_type)
);

-- SavedFilter
CREATE TABLE extras_savedfilter (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    description VARCHAR(200),
    user_id UUID REFERENCES core_user(id) ON DELETE SET NULL,
    weight SMALLINT NOT NULL DEFAULT 100,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    shared BOOLEAN NOT NULL DEFAULT TRUE,
    parameters JSONB NOT NULL,
    owner_id UUID REFERENCES core_user(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- SavedFilter Object Types (Many-to-Many)
CREATE TABLE extras_savedfilter_object_types (
    saved_filter_id UUID REFERENCES extras_savedfilter(id) ON DELETE CASCADE,
    content_type VARCHAR(255) NOT NULL,
    PRIMARY KEY (saved_filter_id, content_type)
);

-- TableConfig
CREATE TABLE extras_tableconfig (
    id UUID PRIMARY KEY,
    object_type VARCHAR(255) NOT NULL,
    "table" VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(200),
    user_id UUID REFERENCES core_user(id) ON DELETE SET NULL,
    weight SMALLINT NOT NULL DEFAULT 1000,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    shared BOOLEAN NOT NULL DEFAULT TRUE,
    columns TEXT[] NOT NULL,
    ordering TEXT[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ImageAttachment
CREATE TABLE extras_imageattachment (
    id UUID PRIMARY KEY,
    object_type VARCHAR(255) NOT NULL,
    object_id BIGINT NOT NULL,
    image TEXT NOT NULL,
    image_height SMALLINT NOT NULL,
    image_width SMALLINT NOT NULL,
    name VARCHAR(50),
    description VARCHAR(200),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_extras_imageattachment_object ON extras_imageattachment(object_type, object_id);

-- JournalEntry
CREATE TABLE extras_journalentry (
    id UUID PRIMARY KEY,
    assigned_object_type VARCHAR(255) NOT NULL,
    assigned_object_id BIGINT NOT NULL,
    created_by UUID REFERENCES core_user(id) ON DELETE SET NULL,
    kind VARCHAR(30) NOT NULL,
    comments TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_extras_journalentry_object ON extras_journalentry(assigned_object_type, assigned_object_id);

-- Bookmark
CREATE TABLE extras_bookmark (
    id UUID PRIMARY KEY,
    object_type VARCHAR(255) NOT NULL,
    object_id BIGINT NOT NULL,
    user_id UUID NOT NULL REFERENCES core_user(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (object_type, object_id, user_id)
);
CREATE INDEX idx_extras_bookmark_object ON extras_bookmark(object_type, object_id);
