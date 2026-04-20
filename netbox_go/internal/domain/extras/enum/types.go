// Package enum содержит перечисления для домена Extras
package enum

import "fmt"

// CustomFieldType определяет типы пользовательских полей
type CustomFieldType string

// CustomFieldChoiceSetBaseChoices определяет типы пользовательских полей для выбора по умолчанию
type CustomFieldChoiceSetBaseChoices bool

var ErrInvalidCustomFieldType = fmt.Errorf("invalid custom field type")

const (
	CustomFieldTypeText        CustomFieldType = "text"
	CustomFieldTypeLongText    CustomFieldType = "longtext"
	CustomFieldTypeInteger     CustomFieldType = "integer"
	CustomFieldTypeDecimal     CustomFieldType = "decimal"
	CustomFieldTypeBoolean     CustomFieldType = "boolean"
	CustomFieldTypeDate        CustomFieldType = "date"
	CustomFieldTypeDateTime    CustomFieldType = "datetime"
	CustomFieldTypeURL         CustomFieldType = "url"
	CustomFieldTypeJSON        CustomFieldType = "json"
	CustomFieldTypeSelect      CustomFieldType = "select"
	CustomFieldTypeMultiSelect CustomFieldType = "multiselect"
	CustomFieldTypeObject      CustomFieldType = "object"
	CustomFieldTypeMultiObject CustomFieldType = "multiobject"
)

// Validate проверяает корректность типа поля
func (t CustomFieldType) Validate() error {
	switch t {
	case CustomFieldTypeText, CustomFieldTypeLongText,
		CustomFieldTypeInteger, CustomFieldTypeDecimal,
		CustomFieldTypeBoolean, CustomFieldTypeDate,
		CustomFieldTypeDateTime, CustomFieldTypeURL,
		CustomFieldTypeJSON, CustomFieldTypeSelect,
		CustomFieldTypeMultiSelect, CustomFieldTypeObject,
		CustomFieldTypeMultiObject:
		return nil
	default:
		return ErrInvalidCustomFieldType
	}
}

// CustomFieldFilterLogic определяет логику фильтрации
type CustomFieldFilterLogic string

const (
	CustomFieldFilterDisabled CustomFieldFilterLogic = "disabled"
	CustomFieldFilterLoose    CustomFieldFilterLogic = "loose"
	CustomFieldFilterExact    CustomFieldFilterLogic = "exact"
)

// CustomFieldUIVisible определяет видимость в UI
type CustomFieldUIVisible string

const (
	CustomFieldUIVisibleAlways CustomFieldUIVisible = "always"
	CustomFieldUIVisibleIfSet  CustomFieldUIVisible = "if-set"
	CustomFieldUIVisibleHidden CustomFieldUIVisible = "hidden"
)

// CustomFieldUIEditable определяет редактируемость в UI
type CustomFieldUIEditable string

const (
	CustomFieldUIEditableYes    CustomFieldUIEditable = "yes"
	CustomFieldUIEditableNo     CustomFieldUIEditable = "no"
	CustomFieldUIEditableHidden CustomFieldUIEditable = "hidden"
)

// EventRuleActionType определяет типы действий для правил событий
type EventRuleActionType string

const (
	EventRuleActionWebhook      EventRuleActionType = "webhook"
	EventRuleActionScript       EventRuleActionType = "script"
	EventRuleActionNotification EventRuleActionType = "notification"
)

// JournalEntryKind определяет типы записей журнала
type JournalEntryKind string

const (
	JournalEntryKindInfo    JournalEntryKind = "info"
	JournalEntryKindSuccess JournalEntryKind = "success"
	JournalEntryKindWarning JournalEntryKind = "warning"
	JournalEntryKindDanger  JournalEntryKind = "danger"
)

// BookmarkOrdering определяет порядок сортировки закладок
type BookmarkOrdering string

const (
	BookmarkOrderingNewest         BookmarkOrdering = "-created"
	BookmarkOrderingOldest         BookmarkOrdering = "created"
	BookmarkOrderingAlphabeticalAZ BookmarkOrdering = "name"
	BookmarkOrderingAlphabeticalZA BookmarkOrdering = "-name"
)
