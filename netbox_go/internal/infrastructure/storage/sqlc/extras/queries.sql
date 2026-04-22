-- EventRule
-- name: CreateEventRule :exec
-- desc: Create a new event rule
INSERT INTO extras_eventrule (id, name, description, event_types, enabled, conditions, action_type, action_object_type, action_object_id, action_data, comments, owner_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);

-- name: GetEventRule :one
-- desc: Get an event rule by ID
SELECT * FROM extras_eventrule WHERE id = $1;

-- name: UpdateEventRule :exec
-- desc: Update an existing event rule
UPDATE extras_eventrule 
SET name = $2, description = $3, event_types = $4, enabled = $5, conditions = $6, action_type = $7, action_object_type = $8, action_object_id = $9, action_data = $10, comments = $11, owner_id = $12, updated_at = NOW()
WHERE id = $1;

-- name: DeleteEventRule :exec
-- desc: Delete an event rule
DELETE FROM extras_eventrule WHERE id = $1;

-- name: ListEventRules :many
-- desc: List event rules with pagination
SELECT * FROM extras_eventrule ORDER BY name LIMIT $1 OFFSET $2;

-- name: CreateEventRuleObjectType :exec
-- desc: Associate an object type with an event rule
INSERT INTO extras_eventrule_object_types (event_rule_id, content_type) VALUES ($1, $2);

-- name: ListEventRuleObjectTypes :many
-- desc: List object types for an event rule
SELECT content_type FROM extras_eventrule_object_types WHERE event_rule_id = $1;

-- Webhook
-- name: CreateWebhook :exec
-- desc: Create a new webhook
INSERT INTO extras_webhook (id, name, description, payload_url, http_method, http_content_type, additional_headers, body_template, secret, ssl_verification, ca_file_path, owner_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);

-- name: GetWebhook :one
-- desc: Get a webhook by ID
SELECT * FROM extras_webhook WHERE id = $1;

-- name: UpdateWebhook :exec
-- desc: Update a webhook
UPDATE extras_webhook 
SET name = $2, description = $3, payload_url = $4, http_method = $5, http_content_type = $6, additional_headers = $7, body_template = $8, secret = $9, ssl_verification = $10, ca_file_path = $11, owner_id = $12, updated_at = NOW()
WHERE id = $1;

-- name: DeleteWebhook :exec
-- desc: Delete a webhook
DELETE FROM extras_webhook WHERE id = $1;

-- name: ListWebhooks :many
-- desc: List webhooks with pagination
SELECT * FROM extras_webhook ORDER BY name LIMIT $1 OFFSET $2;

-- CustomLink
-- name: CreateCustomLink :exec
-- desc: Create a new custom link
INSERT INTO extras_customlink (id, name, enabled, link_text, link_url, weight, group_name, button_class, new_window, owner_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: GetCustomLink :one
-- desc: Get a custom link by ID
SELECT * FROM extras_customlink WHERE id = $1;

-- name: UpdateCustomLink :exec
-- desc: Update a custom link
UPDATE extras_customlink 
SET name = $2, enabled = $3, link_text = $4, link_url = $5, weight = $6, group_name = $7, button_class = $8, new_window = $9, owner_id = $10, updated_at = NOW()
WHERE id = $1;

-- name: DeleteCustomLink :exec
-- desc: Delete a custom link
DELETE FROM extras_customlink WHERE id = $1;

-- name: ListCustomLinks :many
-- desc: List custom links with pagination
SELECT * FROM extras_customlink ORDER BY group_name, weight, name LIMIT $1 OFFSET $2;

-- name: CreateCustomLinkObjectType :exec
-- desc: Associate an object type with a custom link
INSERT INTO extras_customlink_object_types (custom_link_id, content_type) VALUES ($1, $2);

-- name: ListCustomLinkObjectTypes :many
-- desc: List object types for a custom link
SELECT content_type FROM extras_customlink_object_types WHERE custom_link_id = $1;

-- ExportTemplate
-- name: CreateExportTemplate :exec
-- desc: Create a new export template
INSERT INTO extras_exporttemplate (id, name, description, template_code, mime_type, file_name, file_extension, as_attachment, owner_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: GetExportTemplate :one
-- desc: Get an export template by ID
SELECT * FROM extras_exporttemplate WHERE id = $1;

-- name: UpdateExportTemplate :exec
-- desc: Update an export template
UPDATE extras_exporttemplate 
SET name = $2, description = $3, template_code = $4, mime_type = $5, file_name = $6, file_extension = $7, as_attachment = $8, owner_id = $9, updated_at = NOW()
WHERE id = $1;

-- name: DeleteExportTemplate :exec
-- desc: Delete an export template
DELETE FROM extras_exporttemplate WHERE id = $1;

-- name: ListExportTemplates :many
-- desc: List export templates with pagination
SELECT * FROM extras_exporttemplate ORDER BY name LIMIT $1 OFFSET $2;

-- name: CreateExportTemplateObjectType :exec
-- desc: Associate an object type with an export template
INSERT INTO extras_exporttemplate_object_types (export_template_id, content_type) VALUES ($1, $2);

-- name: ListExportTemplateObjectTypes :many
-- desc: List object types for an export template
SELECT content_type FROM extras_exporttemplate_object_types WHERE export_template_id = $1;

-- SavedFilter
-- name: CreateSavedFilter :exec
-- desc: Create a new saved filter
INSERT INTO extras_savedfilter (id, name, slug, description, user_id, weight, enabled, shared, parameters, owner_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: GetSavedFilter :one
-- desc: Get a saved filter by ID
SELECT * FROM extras_savedfilter WHERE id = $1;

-- name: UpdateSavedFilter :exec
-- desc: Update a saved filter
UPDATE extras_savedfilter 
SET name = $2, slug = $3, description = $4, user_id = $5, weight = $6, enabled = $7, shared = $8, parameters = $9, owner_id = $10, updated_at = NOW()
WHERE id = $1;

-- name: DeleteSavedFilter :exec
-- desc: Delete a saved filter
DELETE FROM extras_savedfilter WHERE id = $1;

-- name: ListSavedFilters :many
-- desc: List saved filters with pagination
SELECT * FROM extras_savedfilter ORDER BY weight, name LIMIT $1 OFFSET $2;

-- name: CreateSavedFilterObjectType :exec
-- desc: Associate an object type with a saved filter
INSERT INTO extras_savedfilter_object_types (saved_filter_id, content_type) VALUES ($1, $2);

-- name: ListSavedFilterObjectTypes :many
-- desc: List object types for a saved filter
SELECT content_type FROM extras_savedfilter_object_types WHERE saved_filter_id = $1;

-- TableConfig
-- name: CreateTableConfig :exec
-- desc: Create a new table config
INSERT INTO extras_tableconfig (id, object_type, "table", name, description, user_id, weight, enabled, shared, columns, ordering)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);

-- name: GetTableConfig :one
-- desc: Get a table config by ID
SELECT * FROM extras_tableconfig WHERE id = $1;

-- name: UpdateTableConfig :exec
-- desc: Update a table config
UPDATE extras_tableconfig 
SET object_type = $2, "table" = $3, name = $4, description = $5, user_id = $6, weight = $7, enabled = $8, shared = $9, columns = $10, ordering = $11, updated_at = NOW()
WHERE id = $1;

-- name: DeleteTableConfig :exec
-- desc: Delete a table config
DELETE FROM extras_tableconfig WHERE id = $1;

-- name: ListTableConfigs :many
-- desc: List table configs with pagination
SELECT * FROM extras_tableconfig ORDER BY weight, name LIMIT $1 OFFSET $2;

-- ImageAttachment
-- name: CreateImageAttachment :exec
-- desc: Create a new image attachment
INSERT INTO extras_imageattachment (id, object_type, object_id, image, image_height, image_width, name, description)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetImageAttachment :one
-- desc: Get an image attachment by ID
SELECT * FROM extras_imageattachment WHERE id = $1;

-- name: UpdateImageAttachment :exec
-- desc: Update an image attachment
UPDATE extras_imageattachment 
SET object_type = $2, object_id = $3, image = $4, image_height = $5, image_width = $6, name = $7, description = $8, updated_at = NOW()
WHERE id = $1;

-- name: DeleteImageAttachment :exec
-- desc: Delete an image attachment
DELETE FROM extras_imageattachment WHERE id = $1;

-- name: ListImageAttachments :many
-- desc: List image attachments for a specific object
SELECT * FROM extras_imageattachment WHERE object_type = $1 AND object_id = $2 ORDER BY name, id;

-- JournalEntry
-- name: CreateJournalEntry :exec
-- desc: Create a new journal entry
INSERT INTO extras_journalentry (id, assigned_object_type, assigned_object_id, created_by, kind, comments)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetJournalEntry :one
-- desc: Get a journal entry by ID
SELECT * FROM extras_journalentry WHERE id = $1;

-- name: ListJournalEntries :many
-- desc: List journal entries for a specific object
SELECT * FROM extras_journalentry WHERE assigned_object_type = $1 AND assigned_object_id = $2 ORDER BY created_at DESC;

-- Bookmark
-- name: CreateBookmark :exec
-- desc: Create a new bookmark
INSERT INTO extras_bookmark (id, object_type, object_id, user_id)
VALUES ($1, $2, $3, $4);

-- name: DeleteBookmark :exec
-- desc: Delete a bookmark
DELETE FROM extras_bookmark WHERE id = $1;

-- name: ListBookmarks :many
-- desc: List bookmarks for a user
SELECT * FROM extras_bookmark WHERE user_id = $1 ORDER BY created_at;
