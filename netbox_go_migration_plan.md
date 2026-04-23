# План переписывания NetBox на Go (бэкенд) с сохранением HTMX + Bootstrap 5 (фронтенд)

## Обзор проекта

NetBox — сложная система управления сетевой инфраструктурой (IPAM/DCIM) со следующими основными модулями:
- **DCIM**: Устройства, стойки, порты, кабели, инвентарь
- **IPAM**: IP-адреса, префиксы, VLAN, VRF, пулы адресов
- **Virtualization**: Виртуальные машины, кластеры
- **Circuits**: Цепи подключения, провайдеры
- **Tenancy**: Арендаторы, контакты
- **Users/Auth**: Пользователи, группы, токены, OAuth
- **Extras**: Вебхуки, события, кастомные поля, скрипты, отчёты
- **Core/Wireless/VPN**: Системные функции, беспроводные сети, VPN

---

## Технологический стек

### Бэкенд
| Компонент | Технология | Назначение |
|-----------|------------|------------|
| Язык | Go 1.23+ | Основной язык разработки |
| Web-фреймворк | Echo v4 | HTTP-сервер, роутинг, middleware |
| ORM/SQL | sqlc | Type-safe SQL генерация |
| БД | PostgreSQL 16+ | Хранение данных (сохранение схемы) |
| GraphQL | gqlgen | GraphQL сервер |
| Кэш | go-etcd + etcd v3 | etcd клиент для кэширования и распределённых блокировок (кластер) |
| Конфигурация | Viper | Управление конфигурацией |
| Логирование | Zap | Высокопроизводительное логирование |
| Шаблоны | Go templates | Рендеринг HTML (совместимость с HTMX) |
| Миграции | golang-migrate | Управление миграциями БД |
| Очереди | Asynq | Фоновые задачи (замена RQ) |
| Архитектура приложения | Clean Architecture | Модульная архитектура, чистая архитектура, DDD |

### Фронтенд (сохраняется)
| Компонент | Технология |
|-----------|------------|
| Динамика | HTMX |
| UI-фреймворк | Bootstrap 5 |
| Скрипты | TypeScript (частично) |
| Стили | SCSS |

### Инфраструктура
| Компонент | Технология |
|-----------|------------|
| Контейнеризация | Docker |
| Оркестрация | Kubernetes |
| Деплой | Helm |

---

## Принципы Clean Architecture в архитектуре

### Архитектурный подход Clean Architecture

Проект использует архитектуру **Clean Architecture** — эволюцию чистой архитектуры (Clean Architecture) и предметно-ориентированного проектирования (DDD), адаптированную для Go.

**Ключевые принципы Clean Architecture:**

1. **Модульность по доменам (Domain Modules)**
   - Каждый бизнес-домен (DCIM, IPAM, Auth, Extras) — отдельный модуль
   - Модули инкапсулируют свою бизнес-логику, данные и зависимости
   - Межмодульное взаимодействие только через публичные интерфейсы

2. **Слои внутри модуля (Layered Structure)**
   ```
   module/dcim/
   ├── domain/          # Entities, Value Objects, Domain Services
   ├── application/     # Use Cases, DTOs, Application Services  
   ├── infrastructure/  # Repositories, External adapters
   └── delivery/        # HTTP handlers, GraphQL resolvers, Templates
   ```

3. **Dependency Rule (Правило зависимостей)**
   - Зависимости направлены внутрь: delivery → application → domain
   - Domain слой не зависит ни от чего внешнего
   - Infrastructure реализует интерфейсы из Application/Domain

4. **Ports & Adapters (Гексагональная архитектура)**
   - Ports: интерфейсы для входящих/исходящих операций
   - Adapters: реализации для конкретных технологий (Echo, sqlc, etcd)

5. **CQRS разделение (опционально)**
   - Command handlers для записи
   - Query handlers для чтения
   - Разные модели для write/read операций

### 1. Single Responsibility Principle (SRP)
- Каждый сервис отвечает за одну предметную область (DCIM, IPAM, Auth)
- Отдельные пакеты для: репозиториев, сервисов, хендлеров, DTO
- Разделение логики GraphQL и REST API
- Uber-FX модули изолируют ответственность по доменам

### 2. Open/Closed Principle (OCP)
- Интерфейсы для расширяемости (плагины, бэкенды поиска)
- Dependency Injection через интерфейсы
- Система плагинов через registration pattern
- Uber-FX позволяет добавлять новые модули без изменения существующих

### 3. Liskov Substitution Principle (LSP)
- Единая иерархия моделей через embedding
- Общие интерфейсы для всех типов объектов
- Uber-FX domain entities могут заменяться совместимыми реализациями

### 4. Interface Segregation Principle (ISP)
- Узкоспециализированные интерфейсы (Reader, Writer, Searcher)
- Отдельные интерфейсы для разных типов доступа
- Uber-FX ports определяются конкретно под каждый use case

### 5. Dependency Inversion Principle (DIP)
- Зависимость от абстракций, не реализаций
- Внедрение зависимостей через конструкторы
- Uber-FX инфраструктура зависит от доменных интерфейсов

---

## Архитектура проекта

### Структура с применением Clean Architecture архитектуры

```
netbox-go/
├── cmd/
│   └── netbox/
│       └── main.go              # Точка входа, DI контейнер
├── internal/
│   ├── config/                  # Конфигурация (Viper)
│   │   ├── config.go
│   │   └── settings.go
│   ├── database/                # Подключение к БД
│   │   ├── db.go
│   │   ├── migrations/
│   │   └── queries/             # SQL файлы для sqlc
│   ├── models/                  # Go-структуры (автогенерация sqlc)
│   ├── modules/                 # Clean Architecture модули по доменам
│   │   ├── dcim/                # DCIM модуль
│   │   │   ├── domain/          # Domain слой (Entities, Value Objects)
│   │   │   │   ├── entity/      # Сущности: Device, Rack, Site, Cable
│   │   │   │   ├── valueobject/ # Value Objects: Position, Color
│   │   │   │   ├── service/     # Domain Services
│   │   │   │   └── repository/  # Repository interfaces (Ports)
│   │   │   ├── application/     # Application слой (Use Cases)
│   │   │   │   ├── usecase/     # Use Cases: CreateDevice, TraceCable
│   │   │   │   ├── dto/         # DTO для передачи данных
│   │   │   │   └── service/     # Application Services
│   │   │   ├── infrastructure/  # Infrastructure слой
│   │   │   │   ├── repository/  # Repository implementations (Adapters)
│   │   │   │   └── adapter/     # Внешние адаптеры
│   │   │   └── delivery/        # Delivery слой
│   │   │       ├── http/        # Echo handlers (REST, HTMX)
│   │   │       ├── graphql/     # GraphQL resolvers
│   │   │       └── template/    # Go templates
│   │   ├── ipam/                # IPAM модуль (аналогичная структура)
│   │   │   ├── domain/
│   │   │   ├── application/
│   │   │   ├── infrastructure/
│   │   │   └── delivery/
│   │   ├── users/               # Users/Auth модуль
│   │   │   ├── domain/
│   │   │   ├── application/
│   │   │   ├── infrastructure/
│   │   │   └── delivery/
│   │   ├── extras/              # Extras модуль (webhooks, custom fields)
│   │   │   ├── domain/
│   │   │   ├── application/
│   │   │   ├── infrastructure/
│   │   │   └── delivery/
│   │   └── virtualization/      # Virtualization модуль
│   │       ├── domain/
│   │       ├── application/
│   │       ├── infrastructure/
│   │       └── delivery/
│   ├── shared/                  # Общие компоненты между модулями
│   │   ├── middleware/          # Echo middleware
│   │   │   ├── auth.go
│   │   │   ├── logging.go
│   │   │   ├── request_id.go
│   │   │   └── htmx.go
│   │   ├── auth/                # Аутентификация и авторизация
│   │   │   ├── jwt.go
│   │   │   ├── session.go
│   │   │   ├── oauth.go
│   │   │   └── permissions.go
│   │   ├── cache/               # Кэширование (etcd кластер)
│   │   │   ├── cache.go
│   │   │   ├── client.go        # etcd клиент с поддержкой кластера
│   │   │   └── keys.go
│   │   ├── events/              # Система событий
│   │   │   ├── events.go
│   │   │   ├── rules.go
│   │   │   └── webhooks.go
│   │   ├── jobs/                # Фоновые задачи (Asynq)
│   │   │   ├── processor.go
│   │   │   └── tasks/
│   │   ├── search/              # Поиск и индексация
│   │   │   ├── index.go
│   │   │   └── backends/
│   │   ├── plugins/             # Система плагинов
│   │   │   ├── registry.go
│   │   │   ├── loader.go
│   │   │   └── interfaces.go
│   │   └── templates/           # Go templates (общие)
│   │       ├── base/
│   │       └── functions.go
│   └── utils/                   # Утилиты
│       ├── validation/
│       ├── pagination/
│       └── response/
├── pkg/                         # Публичные пакеты
│   └── api/                     # SDK для клиентов
├── migrations/                  # SQL миграции
├── scripts/                     # Скрипты сборки
├── docker/                      # Docker конфигурации
├── helm/                        # Helm чарты
├── configs/                     # Примеры конфигов
├── go.mod
├── go.sum
├── sqlc.yaml                    # Конфигурация sqlc
├── gqlgen.yml                   # Конфигурация gqlgen
└── Makefile
```

**Преимущества Clean Architecture структуры:**
- Полная изоляция доменов (DCIM не зависит от IPAM напрямую)
- Легкое тестирование каждого слоя отдельно
- Возможность замены инфраструктуры без изменения бизнес-логики
- Масштабируемость команды (разные команды работают с разными модулями)
- Четкое разделение ответственности

---

## Этапы реализации

### Этап 0: Подготовка (2-3 недели)

#### 0.1 Анализ существующей схемы PostgreSQL
- [ ] Экспорт полной схемы БД (`pg_dump --schema-only`)
- [ ] Документирование всех таблиц, индексов, связей
- [ ] Выявление denormalized полей и триггеров
- [ ] Анализ существующих миграций Django

#### 0.2 Проектирование новых таблиц

#### 0.3 Настройка окружения
- [ ] Инициализация Go модуля
- [ ] Настройка `sqlc.yaml` для генерации кода
- [ ] Настройка `gqlgen.yml` для GraphQL
- [ ] Конфигурация Viper (config.yaml, env vars)
- [ ] Настройка Zap logger
- [ ] Docker Compose для локальной разработки (PostgreSQL, etcd кластер)

#### 0.4 CI/CD базовый пайплайн
- [ ] GitHub Actions: build, test, lint
- [ ] Docker image сборка
- [ ] Helm chart базовая структура

---

### Этап 1: Инфраструктурный слой (3-4 недели)

#### 1.1 Конфигурация и логирование
```go
// internal/config/config.go
type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Etcd     EtcdConfig     `mapstructure:"etcd"`
    Auth     AuthConfig     `mapstructure:"auth"`
    Log      LogConfig      `mapstructure:"log"`
}

func Load() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./configs")
    viper.AutomaticEnv()
    // ...
}
```

#### 1.2 Подключение к БД и sqlc
```yaml
# sqlc.yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "internal/database/queries"
    schema: "migrations"
    gen:
      go:
        package: "db"
        out: "internal/database/db"
        sql_package: "database/sql"
```

#### 1.3 etcd клиент (кластер)

**Преимущества etcd перед Redis:**
- Встроенная поддержка кластера (consensus через Raft)
- Сильная согласованность данных (linearizable reads)
- Распределённые блокировки через leases и mutexes
- Watch механизм для real-time уведомлений об изменениях
- Идеально подходит для service discovery и конфигурации

#### 1.4 Middleware Echo
- Request ID генерация (UUID)
- Логирование запросов (Zap)
- Обработка ошибок
- HTMX headers поддержка
- CORS (если нужно)

---

### Этап 2: Система аутентификации и авторизации (4-5 недель)

#### 2.1 Модели пользователей
- Сохранение структуры таблиц `users_user`, `users_group`, `users_objectpermission`
- Генерация sqlc моделей
- Repository паттерн для CRUD операций

#### 2.2 Аутентификация
- [ ] Session-based аутентификация (совместимость с текущей)
- [ ] Token-based API аутентификация
- [ ] OAuth2/OIDC провайдеры (GitHub, Google, Azure AD и т.д.)
- [ ] Remote User Authentication (для reverse proxy)

#### 2.3 Авторизация
- [ ] Система разрешений (permissions)
- [ ] Object-level permissions
- [ ] Constraints на основе атрибутов
- [ ] Role-based access control (RBAC)

#### 2.4 Middleware авторизации

---

### Этап 3: Core модули DCIM (8-10 недель)

#### 3.1 Модели данных
Сохранение структуры таблиц:
- `dcim_device`, `dcim_devicetype`, `dcim_devicerole`
- `dcim_rack`, `dcim_rackreservation`
- `dcim_site`, `dcim_location`
- `dcim_cable`, `dcim_consoleport`, `dcim_interface`
- И другие (~50 таблиц DCIM)

#### 3.2 Сервисный слой в Clean Architecture архитектуре

**Domain слой (entities):**
```go
// internal/modules/dcim/domain/entity/device.go
type Device struct {
    ID          int64
    Name        string
    DeviceType  *DeviceType
    Site        *Site
    Location    *Location
    Rack        *Rack
    Status      DeviceStatus
    // ... другие поля
}

func (d *Device) Validate() error {
    // Domain валидация
}
```

**Application слой (use cases):**
```go
// internal/modules/dcim/application/usecase/create_device.go
type CreateDeviceUseCase interface {
    Execute(ctx context.Context, dto CreateDeviceDTO) (*entity.Device, error)
}

type createDeviceUseCase struct {
    deviceRepo repository.DeviceRepository
    validator  domain.DeviceValidator
    eventPub   events.Publisher
}

func (uc *createDeviceUseCase) Execute(ctx context.Context, dto CreateDeviceDTO) (*entity.Device, error) {
    // Application логика
    device := entity.NewDevice(dto)
    if err := uc.validator.Validate(device); err != nil {
        return nil, err
    }
    
    saved, err := uc.deviceRepo.Create(ctx, device)
    if err != nil {
        return nil, err
    }
    
    uc.eventPub.Publish(ctx, events.DeviceCreated{DeviceID: saved.ID})
    return saved, nil
}
```

**Infrastructure слой (repository implementation):**
```go
// internal/modules/dcim/infrastructure/repository/device_repository.go
type deviceRepository struct {
    db *sqlc.Queries
}

func (r *deviceRepository) GetByID(ctx context.Context, id int64) (*entity.Device, error) {
    row, err := r.db.GetDeviceByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return mapToDevice(row), nil
}
```

**Delivery слой (HTTP handler):**
```go
// internal/modules/dcim/delivery/http/device_handler.go
type DeviceHandler struct {
    createUC usecase.CreateDeviceUseCase
    getUC    usecase.GetDeviceUseCase
    logger   *zap.Logger
}

func (h *DeviceHandler) Create(c echo.Context) error {
    var dto application.CreateDeviceDTO
    if err := c.Bind(&dto); err != nil {
        return err
    }
    
    device, err := h.createUC.Execute(c.Request().Context(), dto)
    if err != nil {
        return err
    }
    
    return c.JSON(http.StatusCreated, device)
}
```

#### 3.3 Бизнес-логика
- [ ] Валидация имен компонентов
- [ ] Tracing кабелей (сложная логика соединений)
- [ ] Power chain calculations
- [ ] Rack elevation calculations
- [ ] Inventory item tree (MPTT эмуляция)

#### 3.4 REST API handlers
```go
// internal/handler/rest/dcim/device_handler.go
func (h *DeviceHandler) RegisterRoutes(g *echo.Group) {
    g.GET("/devices/", h.List)
    g.GET("/devices/:id", h.Get)
    g.POST("/devices/", h.Create)
    g.PUT("/devices/:id", h.Update)
    g.DELETE("/devices/:id", h.Delete)
    g.GET("/devices/:id/trace", h.TraceCable)
}
```

#### 3.5 GraphQL резолверы
```go
// internal/graphql/resolver.go
func (r *queryResolver) Device(ctx context.Context, id int64) (*models.Device, error) {
    return r.deviceService.GetByID(ctx, id)
}

func (r *queryResolver) DeviceList(ctx context.Context, filter DeviceFilter) ([]*models.Device, error) {
    return r.deviceService.List(ctx, filter)
}
```

---

### Этап 4: Модуль IPAM (6-8 недель)

#### 4.1 Модели данных
- `ipam_ipaddress`, `ipam_prefix`, `ipam_vlan`
- `ipam_vrf`, `ipam_routetarget`
- `ipam_aggregate`, `ipam_role`
- `ipam_iprange`, `ipam_fhrpgroup`

#### 4.2 Специфичная логика
- [ ] IP адресация и валидация (netaddr integration)
- [ ] Prefix hierarchy (древовидная структура)
- [ ] Available IP/prefix calculation
- [ ] VLAN groups и scoped VLANs
- [ ] FHRP group assignments

**Пример реализации в Clean Architecture архитектуре:**

```go
// internal/modules/ipam/domain/entity/prefix.go
type Prefix struct {
    ID           int64
    Prefix       net.IPNet
    VRF          *VRF
    Site         *Site
    VLAN         *VLAN
    Status       PrefixStatus
    IsPool       bool
}

func (p *Prefix) GetAvailableIPs() ([]net.IP, error) {
    // Domain логика вычисления доступных IP
}

func (p *Prefix) GetAvailablePrefixes(maskLen int) ([]net.IPNet, error) {
    // Domain логика вычисления доступных подсетей
}
```

```go
// internal/modules/ipam/application/usecase/get_available_ips.go
type GetAvailableIPsUseCase interface {
    Execute(ctx context.Context, prefixID int64) ([]net.IP, error)
}

type getAvailableIPsUseCase struct {
    prefixRepo repository.PrefixRepository
    cache      cache.Cache
}

func (uc *getAvailableIPsUseCase) Execute(ctx context.Context, prefixID int64) ([]net.IP, error) {
    // Проверка кэша
    cacheKey := fmt.Sprintf("ipam:prefix:%d:available_ips", prefixID)
    if cached, err := uc.cache.Get(ctx, cacheKey); err == nil && cached != nil {
        return cached.([]net.IP), nil
    }
    
    // Получение префикса
    prefix, err := uc.prefixRepo.GetByID(ctx, prefixID)
    if err != nil {
        return nil, err
    }
    
    // Вычисление доступных IP
    ips, err := prefix.GetAvailableIPs()
    if err != nil {
        return nil, err
    }
    
    // Кэширование результата
    uc.cache.Set(ctx, cacheKey, ips, 5*time.Minute)
    
    return ips, nil
}
```

---

### Этап 5: Остальные модули (8-10 недель)

#### 5.1 Virtualization
- Кластеры, виртуальные машины
- VM interfaces и disk assignments

#### 5.2 Circuits
- Circuit types, providers, terminations

#### 5.3 Tenancy
- Tenants, contacts, assignments

#### 5.4 Wireless
- Wireless LANs, PHY rates

#### 5.5 VPN
- Tunnels, tunnel groups, IKE policies

---

### Этап 6: Система событий и вебхуков (3-4 недели)

#### 6.1 Event types registry
```go
// internal/events/events.go
type EventType struct {
    Name        string
    Text        string
    Kind        string // info, success, warning, danger
    Destructive bool
}

var eventTypes = map[string]EventType{
    "device.created": {Name: "device.created", Text: "Device Created", Kind: "success"},
    "device.deleted": {Name: "device.deleted", Text: "Device Deleted", Kind: "danger"},
    // ...
}
```

#### 6.2 Event rules
- Условия срабатывания (content type, actions)
- Filters на основе атрибутов

#### 6.3 Webhooks
- HTTP вызовы с retry logic
- Payload templating
- Secret signing

#### 6.4 Фоновая обработка
```go
// internal/jobs/tasks/webhook_task.go
func ProcessWebhook(ctx context.Context, job *asynq.Task) error {
    var payload WebhookPayload
    if err := json.Unmarshal(job.Payload(), &payload); err != nil {
        return err
    }
    // Выполнение webhook вызова
}
```

---

### Этап 7: Кастомные поля и расширения (4-5 недель)

#### 7.1 Custom fields
- Поддержка типов: text, integer, boolean, date, url, select, multi-select
- Привязка к объектам (content type)
- Валидация и фильтрация

#### 7.2 Custom links
- Шаблоны ссылок с переменными
- Группировка

#### 7.3 Export templates
- использовать Go templates
- Экспорт в различные форматы

#### 7.4 Saved filters
- Сохранённые поисковые запросы
- Shared filters

---

### Этап 8: Поиск и индексация (3-4 недели)

#### 8.1 Search backend interface
```go
// internal/search/index.go
type SearchBackend interface {
    Index(ctx context.Context, obj Searchable) error
    Remove(ctx context.Context, objectType string, objectID int64) error
    Search(ctx context.Context, query SearchQuery) (*SearchResult, error)
}
```

#### 8.2 PostgreSQL full-text search реализация
- tsvector генерация
- Weighted search (name, description, custom fields)
- Filtering по типу объекта

#### 8.3 Search indexes registration
```go
// internal/search/backends/postgres.go
type DeviceIndex struct{}

func (i *DeviceIndex) GetFields() []SearchField {
    return []SearchField{
        {Name: "name", Weight: 1.0, Lookup: LookupPartial},
        {Name: "serial", Weight: 0.8, Lookup: LookupExact},
        {Name: "comments", Weight: 0.5, Lookup: LookupPartial},
    }
}
```

---

### Этап 9: HTMX интеграция (4-5 недель)

#### 9.1 Template functions
```go
// internal/templates/functions.go
func RegisterFuncs(funcMap template.FuncMap) {
    funcMap["has_perm"] = hasPermFunc
    funcMap["get_custom_field"] = getCustomFieldFunc
    funcMap["render_component"] = renderComponentFunc
    funcMap["htmx_method"] = htmxMethodFunc
    // Django template tags совместимость
}
```

#### 9.2 HTMX handlers
```go
// internal/handler/htmx/forms.go
func (h *FormHandler) RenderCreateForm(c echo.Context) error {
    ctx := c.Request().Context()
    // Получение данных для формы
    return c.Render(http.StatusOK, "htmx/form.html", data)
}

func (h *FormHandler) SubmitCreate(c echo.Context) error {
    // Обработка POST с валидацией
    // Возврат либо формы с ошибками, либо redirect
}
```

#### 9.3 Компоненты
- [ ] Динамические таблицы с пагинацией
- [ ] Формы с валидацией
- [ ] Modal окна
- [ ] Object selectors (autocomplete)
- [ ] Bulk operations
- [ ] Notifications polling

#### 9.4 Миграция существующих шаблонов
Конвертация Django templates → Go templates:
```django
<!-- Django -->
{% for device in devices %}
  <tr>{{ device.name }}</tr>
{% endfor %}

<!-- Go -->
{{ range .Devices }}
  <tr>{{ .Name }}</tr>
{{ end }}
```

---

### Этап 10: Система плагинов (4-5 недель)

#### 10.1 Plugin interface
```go
// internal/plugins/interfaces.go
type Plugin interface {
    Name() string
    Version() string
    Init(ctx context.Context, cfg map[string]interface{}) error
    RegisterRoutes(g *echo.Group)
    RegisterGraphQLTypes(schema *graphql.Schema)
    RegisterTemplateExtensions() []TemplateExtension
    RegisterMenuItems() []MenuItem
}
```

#### 10.2 Plugin loader
- Загрузка .so файлов или отдельных бинарников
- Изоляция плагинов (опционально gRPC)
- Registry паттерн

#### 10.3 Plugin hooks
- Pre/Post save hooks
- Custom validators
- Additional tabs on objects

---

### Этап 11: Аудит и изменение логов (2-3 недели)

#### 11.1 Change logging
```go
// internal/repository/audit.go
func (r *AuditRepository) LogChange(ctx context.Context, change AuditChange) error {
    // Запись в audit_log таблицу
    // Сравнение старого и нового состояния
}
```

#### 11.2 Object changelog UI
- История изменений объекта
- Diff view
- Filter by user, date, action type

---

### Этап 12: Тестирование (параллельно всем этапам)

#### 12.1 Unit тесты в Clean Architecture архитектуре

**Тестирование Domain слоя:**
```go
// internal/modules/dcim/domain/entity/device_test.go
func TestDevice_Validate(t *testing.T) {
    device := &entity.Device{
        Name: "test-device",
        Status: entity.DeviceStatusActive,
    }
    
    err := device.Validate()
    assert.NoError(t, err)
}
```

**Тестирование Application слоя (Use Cases):**
```go
// internal/modules/dcim/application/usecase/create_device_test.go
func TestCreateDeviceUseCase_Execute(t *testing.T) {
    // Mock dependencies
    mockRepo := &mocks.MockDeviceRepository{}
    mockValidator := &mocks.MockDeviceValidator{}
    mockEventPub := &mocks.MockEventPublisher{}
    
    uc := usecase.NewCreateDeviceUseCase(mockRepo, mockValidator, mockEventPub)
    
    dto := application.CreateDeviceDTO{Name: "test-device"}
    
    mockValidator.On("Validate", mock.Anything).Return(nil)
    mockRepo.On("Create", mock.Anything, mock.Anything).Return(&entity.Device{ID: 1}, nil)
    mockEventPub.On("Publish", mock.Anything, mock.Anything).Return(nil)
    
    result, err := uc.Execute(context.Background(), dto)
    
    assert.NoError(t, err)
    assert.Equal(t, int64(1), result.ID)
    mockValidator.AssertExpectations(t)
    mockRepo.AssertExpectations(t)
    mockEventPub.AssertExpectations(t)
}
```

**Тестирование Infrastructure слоя:**
```go
// internal/modules/dcim/infrastructure/repository/device_repository_test.go
func TestDeviceRepository_GetByID(t *testing.T) {
    // testcontainers для PostgreSQL
    pool, err := dockertest.NewPool("")
    require.NoError(t, err)
    
    resource, _ := pool.Run("postgres", "16", []string{"POSTGRES_PASSWORD=test"})
    defer pool.Purge(resource)
    
    db := setupTestDB(resource.GetHostPort())
    repo := NewDeviceRepository(db)
    
    // Тест с реальной БД
}
```

#### 12.2 Integration тесты
- Тесты с реальной БД (testcontainers)
- API endpoint тесты
- GraphQL resolver тесты

#### 12.3 E2E тесты
- Playwright для HTMX функциональности
- Критические пользовательские сценарии

#### 12.4 Нагрузочное тестирование
- k6 или vegeta
- Benchmark критичных endpoints

---

### Этап 13: Миграция данных и деплой (3-4 недели)

#### 13.1 Миграция данных
- Прямое использование существующей БД (совместимость схем)
- Скрипты проверки целостности данных
- Dual-write период (опционально)

#### 13.2 Docker образ
```dockerfile
FROM golang:1.23-alpine AS builder
# Build stage

FROM alpine:3.20
# Runtime stage с templates
```

#### 13.3 Kubernetes манифесты
- Deployment с readiness/liveness probes
- ConfigMap для конфигурации
- Secrets для чувствительных данных
- HPA для автоскейлинга

#### 13.4 Helm chart
```yaml
# values.yaml
replicaCount: 3
image:
  repository: netbox-go
  tag: latest
database:
  host: postgres
  port: 5432
etcd:
  endpoints:
    - http://etcd-0.etcd:2379
    - http://etcd-1.etcd:2379
    - http://etcd-2.etcd:2379
  dial_timeout: 5s
```

**Docker Compose для etcd кластера (локальная разработка):**
```yaml
version: '3.8'
services:
  etcd-0:
    image: quay.io/coreos/etcd:v3.5.11
    command:
      - etcd
      - --name=etcd-0
      - --advertise-client-urls=http://etcd-0:2379
      - --listen-client-urls=http://0.0.0.0:2379
      - --initial-advertise-peer-urls=http://etcd-0:2380
      - --listen-peer-urls=http://0.0.0.0:2380
      - --initial-cluster=etcd-0=http://etcd-0:2380,etcd-1=http://etcd-1:2380,etcd-2=http://etcd-2:2380
      - --initial-cluster-state=new
    ports:
      - "2379:2379"
      - "2380:2380"
  
  etcd-1:
    image: quay.io/coreos/etcd:v3.5.11
    command:
      - etcd
      - --name=etcd-1
      - --advertise-client-urls=http://etcd-1:2379
      - --listen-client-urls=http://0.0.0.0:2379
      - --initial-advertise-peer-urls=http://etcd-1:2380
      - --listen-peer-urls=http://0.0.0.0:2380
      - --initial-cluster=etcd-0=http://etcd-0:2380,etcd-1=http://etcd-1:2380,etcd-2=http://etcd-2:2380
      - --initial-cluster-state=new
    ports:
      - "2381:2379"
      - "2382:2380"
  
  etcd-2:
    image: quay.io/coreos/etcd:v3.5.11
    command:
      - etcd
      - --name=etcd-2
      - --advertise-client-urls=http://etcd-2:2379
      - --listen-client-urls=http://0.0.0.0:2379
      - --initial-advertise-peer-urls=http://etcd-2:2380
      - --listen-peer-urls=http://0.0.0.0:2380
      - --initial-cluster=etcd-0=http://etcd-0:2380,etcd-1=http://etcd-1:2380,etcd-2=http://etcd-2:2380
      - --initial-cluster-state=new
    ports:
      - "2383:2379"
      - "2384:2380"

#### 13.5 Мониторинг
- Prometheus metrics (Echo prometheus middleware)
- Health check endpoints
- Distributed tracing (OpenTelemetry опционально)

---

## Детальный план по неделям

| Недели | Этап | Основные задачи |
|--------|------|-----------------|
| 1-3 | 0 | Анализ схемы, проектирование, настройка окружения |
| 4-7 | 1 | Инфраструктура: конфиг, БД, etcd кластер, логирование, middleware |
| 8-12 | 2 | Auth: пользователи, сессии, токены, permissions |
| 13-22 | 3 | DCIM: все модели, бизнес-логика, API, GraphQL |
| 23-30 | 4 | IPAM: адреса, префиксы, VLAN, логика availability |
| 31-40 | 5 | Остальные модули: Virtualization, Circuits, Tenancy и др. |
| 41-44 | 6 | События, вебхуки, фоновые задачи |
| 45-49 | 7 | Кастомные поля, export templates, saved filters |
| 50-53 | 8 | Поиск и индексация |
| 54-58 | 9 | HTMX интеграция, миграция шаблонов |
| 59-63 | 10 | Система плагинов |
| 64-66 | 11 | Аудит и changelog |
| 67-70 | 12 | Полное тестирование, оптимизация |
| 71-74 | 13 | Миграция, деплой, документация |

**Общая оценка: 16-18 месяцев для команды из 3-4 разработчиков**

---

## Риски и стратегии минимизации

### Риски

1. **Сложность бизнес-логики**
   - Cable tracing, IP calculations, MPTT trees
   - *Стратегия*: Поэтапная реализация с полным покрытием тестами

2. **Потеря функционала при миграции шаблонов**
   - Django template tags vs Go templates
   - *Стратегия*: Создание совместимых template functions, постепенная миграция

3. **Производительность Go templates vs Django**
   - *Стратегия*: Benchmarking, кэширование rendered fragments

4. **Совместимость плагинов**
   - Python plugins не будут работать
   - *Стратегия*: Четкий plugin API, миграционные гайды, wrapper для простых плагинов

5. **Объем работы**
   - 200+ моделей, тысячи строк бизнес-логики
   - *Стратегия*: Приоритизация по частоте использования, incremental rollout

### Стратегия постепенного внедрения

1. **Phase 1**: Запуск Go бэкенда параллельно с Django
   - Reverse proxy маршрутизирует новые фичи на Go
   - Общие БД и etcd кластер

2. **Phase 2**: Постепенный перенос endpoints
   - Сначала простые CRUD операции
   - Затем сложная бизнес-логика

3. **Phase 3**: Полный переход
   - Django остается только для legacy plugin support
   - Постепенное отключение Django модулей

---

## Критерии приемки

### Функциональные
- [ ] Все основные CRUD операции работают
- [ ] GraphQL API полностью функционален
- [ ] HTMX динамика работает идентично оригиналу
- [ ] Аутентификация и авторизация полностью совместимы
- [ ] Вебхуки и события срабатывают корректно
- [ ] Поиск возвращает релевантные результаты
- [ ] Плагины могут быть подключены

### Нефункциональные
- [ ] Время ответа API < 100ms для 95% запросов
- [ ] Поддержка 1000+ concurrent пользователей
- [ ] Graceful shutdown без потери данных
- [ ] Полное покрытие тестами (>80%)
- [ ] Документация API и разработки

### Совместимость
- [ ] Существующая БД работает без модификаций
- [ ] API backward compatibility (для внешних интеграций)
- [ ] Шаблоны рендерятся корректно
- [ ] Миграционные скрипты проверены

---

## Рекомендации по команде

### Минимальный состав
- 1 Tech Lead / Architect (Go, системы)
- 2-3 Backend разработчика (Go, SQL, GraphQL)
- 1 Frontend разработчик (HTMX, Bootstrap, TypeScript)
- 1 DevOps инженер (Kubernetes, CI/CD)

### Желаемый состав
- 1 Tech Lead
- 4-5 Backend разработчиков
- 2 Frontend разработчика
- 1 QA инженер
- 1 DevOps инженер

**Итого: 8-10 человек для ускорения до 10-12 месяцев**

---

## Следующие шаги

1. **Немедленно**:
   - Создать репозиторий netbox-go
   - Настроить базовую структуру проекта
   - Экспортировать схему PostgreSQL
   - Написать детальные спецификации для каждого модуля

2. **Первый спринт (2 недели)**:
   - Реализовать конфигурацию и подключение к БД
   - Настроить sqlc и сгенерировать первые модели
   - Создать базовый HTTP сервер с health check
   - Настроить CI/CD пайплайн

3. **Второй спринт**:
   - Реализовать аутентификацию
   - Начать миграцию простых моделей (tenancy, circuits)
   - Настроить GraphQL базовую схему
