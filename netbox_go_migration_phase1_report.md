# Перенос моделей DCIM на Go

## Выполненные работы

### 1. Создана структура проекта в соответствии с архитектурой Hive

```
netbox_go/
├── pkg/
│   └── types/
│       ├── common.go    # Общие типы (ID, Slug, Coordinate и др.)
│       └── errors.go    # Ошибки валидации и бизнес-логики
└── internal/
    └── domain/
        └── dcim/
            ├── enum/
            │   └── status.go  # Перечисления статусов
            └── entity/
                ├── doc.go     # Документация пакета
                ├── sites.go   # Сущности Sites модуля
                └── racks.go   # Сущности Racks модуля
```

### 2. Реализованные типы данных (`pkg/types/common.go`)

- `ID` - UUID идентификатор сущности
- `Status` - статус сущности
- `TimeStamp` - метки времени создания/обновления
- `AuditInfo` - информация об аудите
- `Slug` - slug-строка для URL с валидацией
- `Description`, `Comments` - текстовые поля
- `Image` - изображение
- `Contact`, `Tenant`, `ASN` - справочные типы
- `Coordinate` - GPS координаты с валидацией (-90..90, -180..180)
- `Address`, `TimeZone`, `Facility` - строковые типы

### 3. Реализованные ошибки (`pkg/types/errors.go`)

**Ошибки валидации:**
- `ErrInvalidLatitude` - некорректная широта
- `ErrInvalidLongitude` - некорректная долгота
- `ErrInvalidSlug` - некорректный slug
- `ErrInvalidStatus` - некорректный статус

**Ошибки бизнес-логики:**
- `ErrNotFound` - сущность не найдена
- `ErrAlreadyExists` - сущность уже существует
- `ErrValidationFailed` - ошибка валидации
- `ErrPermissionDenied` - доступ запрещён
- `ErrInvalidOperation` - недопустимая операция
- `ErrConstraintViolation` - нарушение ограничения

### 4. Реализованные перечисления (`internal/domain/dcim/enum/status.go`)

**SiteStatus:**
- `Planned`, `Staging`, `Active`, `Retired`

**LocationStatus:**
- `Planned`, `Staging`, `Active`, `Retired`

**RackStatus:**
- `Reserved`, `Available`, `Planned`, `Active`, `Deprecated`

**DeviceStatus:**
- `Offline`, `Active`, `Planned`, `Staged`, `Failed`, `Inventory`, `Decommissioning`

**RackType:**
- `Cabinet4Post`, `Cabinet2Post`, `OpenFrame`, `Enclosure`

**RackDimensionUnit:**
- `Millimeter`, `Inch`

Каждое перечисление имеет методы:
- `GetAll*Statuses()` - получение всех возможных значений
- `Validate()` - проверка корректности значения

### 5. Реализованные сущности Sites (`internal/domain/dcim/entity/sites.go`)

**Region:**
- Поля: ID, Name, Slug, Description, ParentID, Created, Updated
- Метод `Validate()` - проверка имени и slug

**SiteGroup:**
- Поля: ID, Name, Slug, Description, ParentID, Created, Updated
- Метод `Validate()` - проверка имени и slug

**Site:**
- Поля: ID, Name, Slug, Status, RegionID, GroupID, TenantID, Facility, ASNIDs, TimeZone, PhysicalAddress, ShippingAddress, Latitude, Longitude, Description, Comments, Created, Updated
- Метод `Validate()` - комплексная проверка всех полей включая координаты
- Метод `GetStatusColor()` - получение цвета статуса

**Location:**
- Поля: ID, Name, Slug, SiteID, Status, ParentID, TenantID, Facility, Description, Comments, Created, Updated
- Метод `Validate()` - проверка обязательных полей и статуса
- Метод `GetStatusColor()` - получение цвета статуса

### 6. Реализованные сущности Racks (`internal/domain/dcim/entity/racks.go`)

**RackType:**
- Поля: ID, ManufacturerID, Model, Slug, Description, FormFactor, Width, UHeight, StartingUnit, DescUnits, OuterWidth/Height/Depth, OuterUnit, MountingDepth, Weight, MaxWeight, WeightUnit, Created, Updated
- Метод `Validate()` - проверка всех ограничений (высота 1-1000U, ширина 19/23")
- Метод `FullName()` - полное название с производителем
- Метод `Units()` - генерация списка единиц стойки (сверху вниз или снизу вверх)

**RackRole:**
- Поля: ID, Name, Slug, Color, Description, Created, Updated
- Метод `Validate()` - проверка имени и slug

**Rack:**
- Поля: ID, Name, FacilityID, SiteID, LocationID, TenantID, Status, RoleID, RackTypeID, FormFactor, Width, Serial, AssetTag, Airflow, UHeight, StartingUnit, DescUnits, OuterDimensions, Weight, Description, Comments, Created, Updated
- Метод `Validate()` - комплексная проверка всех полей
- Метод `GetStatusColor()` - получение цвета статуса
- Метод `CopyRackTypeAttrs()` - копирование атрибутов из типа стойки

**RackReservation:**
- Поля: ID, RackID, UserID, TenantID, Units, Description, Created, Updated
- Метод `Validate()` - проверка обязательных полей

### 7. Настроен Go модуль

Файл `go.mod`:
```go
module github.com/AlekseyPromet/netbox_go

go 1.19

require (
	github.com/google/uuid v1.4.0
)
```

### 8. Сборка и тестирование

```bash
cd netbox_go
go mod tidy      # Установка зависимостей
go build ./...   # Успешная сборка без ошибок
go test ./...    # Тесты готовы к добавлению
```

## Следующие шаги

1. **Добавить сущности Devices:**
   - Manufacturer, DeviceType, Device
   - ModuleType, Module
   - DeviceBay, InventoryItem

2. **Добавить сущности Components:**
   - ConsolePort, ConsoleServerPort
   - PowerPort, PowerOutlet
   - Interface, FrontPort, RearPort
   - DeviceBay, InventoryItem, ModuleBay

3. **Добавить сущности Cables:**
   - Cable, CableTermination

4. **Добавить сущности Power:**
   - PowerPanel, PowerFeed

5. **Создать repository layer:**
   - Интерфейсы репозиториев для каждой сущности
   - SQL запросы через sqlc

6. **Создать service layer:**
   - Бизнес-логика для каждой сущности
   - Валидация, транзакции

7. **Создать delivery layer:**
   - HTTP handlers (Echo)
   - GraphQL resolvers (gqlgen)
   - Templates (Go templates + HTMX)

8. **Настроить инфраструктуру:**
   - etcd клиент для кэша и блокировок
   - Viper для конфигурации
   - Zap для логирования

## Соответствие SOLID принципам

✅ **Single Responsibility Principle (SRP):**
- Каждая сущность отвечает только за свои данные
- Метод Validate() инкапсулирует логику валидации
- Отдельные пакеты для типов, перечислений и сущностей

✅ **Open/Closed Principle (OCP):**
- Сущности открыты для расширения через композицию
- Перечисления можно расширять новыми значениями

✅ **Liskov Substitution Principle (LSP):**
- Все сущности следуют единому паттерну
- Метод Validate() возвращает error во всех сущностях

✅ **Interface Segregation Principle (ISP):**
- Маленькие специализированные интерфейсы (будут созданы в repository layer)

✅ **Dependency Inversion Principle (DIP):**
- Сущности не зависят от инфраструктуры
- Зависимости будут внедряться через конструкторы в service layer
