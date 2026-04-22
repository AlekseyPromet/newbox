-- EventRule
-- name: GetEventRule :id
SELECT * FROM extras_eventrule WHERE id = $1;
-- name: CreateEventRule :event_rule
INSERT INTO extras_eventrule (id, name, description, event_types, enabled, conditions, action_type, action_object_type, action_object_id, action_data, comments, owner_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING *;
-- name: UpdateEventRule :event_rule
UPDATE extras_eventrule 
SET name = $2, description = $3, event_types = $4, enabled = $5, conditions = $6, action_type = $7, action_object_type = $8, action_object_id = $9, action_data = $10, comments = $11, owner_id = $12, updated_at = $14
WHERE id = $1
RETURNING *;
-- name: DeleteEventRule :id
DELETE FROM extras_eventrule WHERE id = $1;
-- name: ListEventRules :offset :limit
SELECT * FROM extras_eventrule ORDER BY name LIMIT $2 OFFSET $1;

-- Webhook
-- name: GetWebhook :id
SELECT * FROM extras_webhook WHERE id = $1;
-- name: CreateWebhook :webhook
INSERT INTO extras_webhook (id, name, description, payload_url, http_method, http_content_type, additional_headers, body_template, secret, ssl_verification, ca_file_path, owner_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING *;
-- name: UpdateWebhook :webhook
UPDATE extras_webhook 
SET name = $2, description = $3, payload_url = $4, http_method = $5, http_content_type = $6, additional_headers = $7, body_template = $8, secret = $9, ssl_verification = $10, ca_file_path = $11, owner_id = $12, updated_at = $14
WHERE id = $1
RETURNING *;
-- name: DeleteWebhook :id
DELETE FROM extras_webhook WHERE id = $1;
-- name: ListWebhooks :offset :limit
SELECT * FROM extras_webhook ORDER BY name LIMIT $2 OFFSET $1;

-- CustomLink
-- name: GetCustomLink :id
SELECT * FROM extras_customlink WHERE id = $1;
-- name: CreateCustomLink :custom_link
INSERT INTO extras_customlink (id, name, enabled, link_text, link_url, weight, group_name, button_class, new_window, owner_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;
-- name: UpdateCustomLink :custom_link
UPDATE extras_customlink 
SET name = $2, enabled = $3, link_text = $4, link_url = $5, weight = $6, group_name = $7, button_class = $8, new_window = $9, owner_id = $10, updated_at = $12
WHERE id = $1
RETURNING *;
-- name: DeleteCustomLink :id
DELETE FROM extras_customlink WHERE id = $1;
-- name: ListCustomLinks :offset :limit
SELECT * FROM extras_customlink ORDER BY group_name, weight, name LIMIT $2 OFFSET $1;

-- ExportTemplate
-- name: GetExportTemplate :id
SELECT * FROM extras_exporttemplate WHERE id = $1;
-- name: CreateExportTemplate :export_template
INSERT INTO extras_exporttemplate (id, name, description, template_code, mime_type, file_name, file_extension, as_attachment, owner_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;
-- name: UpdateExportTemplate :export_template
UPDATE extras_exporttemplate 
SET name = $2, description = $3, template_code = $4, mime_type = $5, file_name = $6, file_extension = $7, as_attachment = $8, owner_id = $9, updated_at = $11
WHERE id = $1
RETURNING *;
-- name: DeleteExportTemplate :id
DELETE FROM extras_exporttemplate WHERE id = $1;
-- name: ListExportTemplates :offset :limit
SELECT * FROM extras_exporttemplate ORDER BY name LIMIT $2 OFFSET $1;

-- SavedFilter
-- name: GetSavedFilter :id
SELECT * FROM extras_savedfilter WHERE id = $1;
-- name: CreateSavedFilter :saved_filter
INSERT INTO extras_savedfilter (id, name, slug, description, user_id, weight, enabled, shared, parameters, owner_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;
-- name: UpdateSavedFilter :saved_filter
UPDATE extras_savedfilter 
SET name = $2, slug = $3, description = $4, user_id = $5, weight = $6, enabled = $7, shared = $8, parameters = $9, owner_id = $10, updated_at = $12
WHERE id = $1
RETURNING *;
-- name: DeleteSavedFilter :id
DELETE FROM extras_savedfilter WHERE id = $1;
-- name: ListSavedFilters :offset :limit
SELECT * FROM extras_savedfilter ORDER BY weight, name LIMIT $2 OFFSET $1;

-- TableConfig
-- name: GetTableConfig :id
SELECT * FROM extras_tableconfig WHERE id = $1;
-- name: CreateTableConfig :table_config
INSERT INTO extras_tableconfig (id, object_type, table, name, description, user_id, weight, enabled, shared, columns, ordering, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING *;
-- name: UpdateTableConfig :table_config
UPDATE extras_tableconfig 
SET object_type = $2, table = $3, name = $4, description = $5, user_id = $6, weight = $7, enabled = $8, shared = $9, columns = $10, ordering = $11, updated_at = $13
WHERE id = $1
RETURNING *;
-- name: DeleteTableConfig :id
DELETE FROM extras_tableconfig WHERE id = $1;
-- name: ListTableConfigs :offset :limit
SELECT * FROM extras_tableconfig ORDER BY weight, name LIMIT $2 OFFSET $1;

-- ImageAttachment
-- name: GetImageAttachment :id
SELECT * FROM extras_imageattachment WHERE id = $1;
-- name: CreateImageAttachment :image_attachment
INSERT INTO extras_imageattachment (id, object_type, object_id, image, image_height, image_width, name, description, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;
-- name: UpdateImageAttachment :image_attachment
UPDATE extras_imageattachment 
SET object_type = $2, object_id = $3, image = $4, image_height = $5, image_width = $6, name = $7, description = $8, updated_at = $10
WHERE id = $1
RETURNING *;
-- name: DeleteImageAttachment :id
DELETE FROM extras_imageattachment WHERE id = $1;
-- name: ListImageAttachments :offset :limit
SELECT * FROM extras_imageattachment ORDER BY name, id LIMIT $2 OFFSET $1;

-- JournalEntry
-- name: GetJournalEntry :id
SELECT * FROM extras_journalentry WHERE id = $1;
-- name: CreateJournalEntry :journal_entry
INSERT INTO extras_journalentry (id, assigned_object_type, assigned_object_id, created_by, kind, comments, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;
-- name: UpdateJournalEntry :journal_entry
UPDATE extras_journalentry 
SET assigned_object_type = $2, assigned_object_id = $3, created_by = $4, kind = $5, comments = $6, updated_at = $8
WHERE id = $1
RETURNING *;
-- name: DeleteJournalEntry :id
DELETE FROM extras_journalentry WHERE id = $1;
-- name: ListJournalEntries :offset :limit
SELECT * FROM extras_journalentry ORDER BY created_at DESC LIMIT $2 OFFSET $1;

-- Bookmark
-- name: GetBookmark :id
SELECT * FROM extras_bookmark WHERE id = $1;
-- name: CreateBookmark :bookmark
INSERT INTO extras_bookmark (id, object_type, object_id, user_id, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
-- name: DeleteBookmark :id
DELETE FROM extras_bookmark WHERE id = $1;
-- name: ListBookmarks :offset :limit
SELECT * FROM extras_bookmark ORDER BY created_at, id LIMIT $2 OFFSET $1;
