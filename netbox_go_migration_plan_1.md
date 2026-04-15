# Этап 1: Анализ текущего состояния проекта NetBox

## 1. Инвентаризация модулей и моделей данных

### Основные модули (Core Modules)

#### DCIM (Data Center Infrastructure Management)
**Объем:** ~8,800 строк кода моделей
**Ключевые сущности:**
- **Sites/Regions/Locations:** Иерархическая структура физических местоположений
- **Racks:** Стойки с позиционированием (elevation), профилями мощности
- **Devices:** Устройства с типами, ролями, производителями
- **Device Components:** 
  - Interfaces (сетевые интерфейсы с IP, VLAN,LAG)
  - Console/Power Ports
  - Front/Rear Ports
  - Device Bays
  - Modules (модульные устройства)
- **Cables:** Трассировка соединений между компонентами
- **Power:** Power Feeds, Panels, распределение питания
- **Templates:** Шаблоны устройств для массового развертывания

#### IPAM (IP Address Management)
**Объем:** ~1,976 строк кода моделей
**Ключевые сущности:**
- **IP Addresses/Prefixes:** Управление IP адресами и подсетями
- **Aggregates:** Глобальные агрегаты адресного пространства
- **VLANs/VLAN Groups:** Виртуальные сети и их группировка
- **VRFs:** Виртуальные таблицы маршрутизации
- **ASNs:** Автономные системы
- **FHRP Groups:** Протоколы резервирования шлюза (HSRP/VRRP)
- **Services:** Сервисы (HTTP, SSH, DNS и т.д.) привязанные к IP

#### Дополнительные модули
- **Circuits:** Провайдеры, цепи, типы цепей
- **Virtualization:** Кластеры, виртуальные машины, диски
- **Tenancy:** Tenants, tenant groups
- **Users & Permissions:** Пользователи, группы, токены
- **Extras:** Custom fields, tags, webhooks, scripts, reports
- **Core:** Jobs, config revisions
- **VPN:** Tunnels, terminations
- **Wireless:** SSID, радио профили

## 2. Архитектура API

### Структура API в каждом модуле
Каждый модуль содержит стандартный набор файлов:
- `api/views.py` - Django REST Framework ViewSets
- `api/serializers_.py` - Сериализаторы для JSON представления
- `api/filtersets.py` - Фильтрация запросов (django-filter)
- `graphql/schema.py` - GraphQL схемы (Graphene)
- `api/urls.py` - Маршрутизация endpoints

### Характеристики API
- **REST API:** Полная CRUD операция для всех сущностей
- **GraphQL:** Гибкие запросы с вложенными данными
- **Фильтрация:** Расширенная фильтрация по всем полям
- **Пагинация:** Limit/offset и cursor-based пагинация
- **Bulk операции:** Массовое создание/обновление/удаление
- **Выбор полей:** `?fields=id,name,status` для оптимизации

### Количество endpoints
- **~200+ REST endpoints** для всех сущностей
- **Единый GraphQL endpoint** `/graphql/`
- **Специализированные endpoints:** 
  - `/api/ipam/prefixes/{id}/available_ips/`
  - `/api/dcim/racks/{id}/elevation/`
  - `/api/dcim/cables/{id}/trace/`

## 3. Модели данных Django

### Базовые классы моделей

#### PrimaryModel
Базовый класс для всех основных сущностей:
```python
class PrimaryModel(ChangeLoggedModel, CustomFieldsMixin, 
                   CustomLinksMixin, TagsMixin, models.Model):
    id = models.AutoField(primary_key=True)
    created = models.DateTimeField(auto_now_add=True)
    last_updated = models.DateTimeField(auto_now=True)
    custom_field_data = models.JSONField(default=dict)
```

#### OrganizationalModel
Для организационных сущностей (Region, Site Group, Tenant Group):
- Наследует PrimaryModel
- Добавляет иерархические связи (parent/children)

#### NestedGroupModel
Для вложенных группировок с MPTT (Modified Preorder Tree Traversal)

#### ChangeLoggedModel
Автоматическое логирование изменений:
- Отслеживание создателя/редактора
- История изменений (ObjectChange модель)

#### CabledObjectModel
Для объектов, участвующих в кабельных соединениях:
- cable_end (A/B сторона)
- link_id, link_peering_type
- Методы для трассировки кабелей

### Миксины (Mixins)
- **CustomFieldsMixin:** Динамические пользовательские поля
- **CustomLinksMixin:** Пользовательские ссылки в UI
- **TagsMixin:** Система тегирования
- **ContactsMixin:** Привязка контактов к объектам
- **ImageAttachmentsMixin:** Вложение изображений

### Особенности моделей
- **Natural Ordering:** Естественная сортировка (rack U1, U2, U10)
- **Status Fields:** Enum поля статуса через Django Choices
- **JSONB:** Хранение custom fields в PostgreSQL JSONB
- **GistIndex:** Гео-индексы для карт и локаций

## 4. Схема PostgreSQL

### Общие паттерны таблиц

#### Стандартные поля
```sql
CREATE TABLE dcim_device (
    id              SERIAL PRIMARY KEY,
    created         TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_updated    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    name            VARCHAR(100),
    slug            VARCHAR(100) UNIQUE,
    status          VARCHAR(50) NOT NULL,
    custom_field_data JSONB NOT NULL DEFAULT '{}',
    -- ... остальные поля
);
```

#### Индексы
- **Primary Keys:** SERIAL или UUID
- **Unique Constraints:** slug, name + site_id
- **Foreign Keys:** Связи между таблицами с ON DELETE CASCADE/SET NULL
- **GistIndex:** Для гео-запросов (координаты сайтов)
- **GIN Index:** Для JSONB полей (custom_field_data)
- **Composite Indexes:** Для частых запросов фильтрации

#### Специфичные типы данных
- **INET/CIDR:** Для IP адресов и префиксов
- **MACADDR:** Для MAC адресов интерфейсов
- **UUID:** Для некоторых сущностей и токенов
- **ARRAY:** Для списков (например, ASN списки)

### Ключевые таблицы
- **core_job:** Фоновые задачи и их статус
- **extras_objectchange:** Audit log всех изменений
- **extras_customfield:** Метаданные пользовательских полей
- **extras_tag:** Теги и их назначения
- **users_token:** API токены пользователей

## 5. Бизнес-логика требующая портирования

### IPAM Calculations
- **Префикс арифметика:** Расчет доступных IP, разделение префиксов
- **Иерархия префиксов:** Автоматическое определение parent/child
- **Доступность:** Поиск свободных префиксов/IP заданного размера
- **Утилиты:** netaddr библиотека для манипуляций с IP

### Cable Tracing
- **Трассировка соединений:** От интерфейса до интерфейса через патч-панели
- **Визуализация пути:** Построение полного пути кабеля
- **Валидация:** Проверка совместимости типов портов
- **Рекурсивный обход:** Графовая навигация по соединениям

### Rack Elevation
- **Позиционирование устройств:** Расчет занимаемых U (units)
- **Визуализация:** Генерация SVG/front-rear views
- **Конфликты:** Обнаружение пересечений устройств
- **Резервирование:** Зарезервированные единицы стойки

### Prefix/VLAN Availability
- **Поиск доступных:** Алгоритмы поиска свободных блоков
- **Утилизация:** Расчет процента использования
- **Предложения:** Рекомендации по размещению

### Natural Ordering
- **Алфавитно-цифровая сортировка:** "interface1" < "interface2" < "interface10"
- **Реализация:** Custom collation или приложение-уровень сортировка

### Permissions & Authorization
- **Object-level permissions:** Доступ к конкретным объектам
- **Action-based:** Разрешения на view/add/change/delete
- **Constraints:** Ограничения через Q-объекты Django
- **Groups & Users:** Ролевая модель доступа

### Webhooks & Events
- **Trigger conditions:** Создание/обновление/удаление объектов
- **Payload serialization:** JSON представление изменений
- **Retry logic:** Повторные попытки при ошибках
- **Secrets:** Подписывание запросов хэшами

### Custom Scripts & Reports
- **Python скрипты:** Пользовательская бизнес-логика
- **Отчеты:** Аудит и валидация данных
- **Планировщик:** Запуск по расписанию или вручную
- **Параметризация:** Входные параметры для скриптов

## 6. Зависимости Python

### Основные зависимости
```txt
Django>=5.0,<5.1                  # Web фреймворк
djangorestframework>=3.14         # REST API
django-filter>=23.5               # Фильтрация запросов
django-mptt>=0.14                 # Деревья (MPTT алгоритм)
django-taggit>=5.0                # Теги
django-cors-headers>=4.3          # CORS заголовки
django-prometheus>=2.3            # Prometheus метрики
graphene-django>=3.2              # GraphQL
psycopg2-binary>=2.9              # PostgreSQL драйвер
netaddr>=0.9                      # IP адресация
Pillow>=10.2                      # Работа с изображениями
pyyaml>=6.0                       # YAML парсинг
requests>=2.31                    # HTTP клиент
social-auth-core>=4.4             # OAuth/Social auth
social-auth-app-django>=5.4       # Django integration
drf-spectacular>=0.27             # OpenAPI документация
svgwrite>=1.4                     # SVG генерация
tzdata>=2024.1                    # Timezone данные
markdown>=3.5                     # Markdown рендеринг
bleach>=6.1                       # HTML санитайзинг
```

### Утилиты и инструменты
```txt
celery>=5.3                       # Task queue
redis>=5.0                        # Redis client (для Celery & cache)
rq>=1.15                          # Alternative task queue
gunicorn>=21.2                    # WSGI server
uvicorn>=0.27                     # ASGI server (для async)
```

### Dev зависимости
```txt
pytest>=8.0                       # Тестирование
pytest-django>=4.7                # Django pytest plugin
factory-boy>=3.3                  # Test fixtures
flake8>=7.0                       # Linting
black>=24.2                       # Code formatting
isort>=5.13                       # Import sorting
mypy>=1.8                         # Type checking
```

## 7. Frontend

### Технологический стек
- **HTMX:** Динамические обновления без full page reload
- **Bootstrap 5:** UI компоненты и сетка
- **TypeScript:** Частичное использование для сложной логики
- **SCSS:** Стилизация с переменными и миксинами
- **Django Templates:** Server-side рендеринг

### Ключевые компоненты
- **Таблицы:** Сортировка, фильтрация, пагинация
- **Формы:** Валидация, динамические поля, bulk редактирование
- **Навигация:** Боковое меню, хлебные крошки
- **Визуализация:** 
  - Rack elevation (SVG)
  - Cable trace diagrams
  - Prefix hierarchy trees
  - Maps с гео-данными

### Интерактивные элементы
- **Bulk operations:** Выделение множества объектов
- **Quick search:** Живой поиск по всем сущностям
- **Modal dialogs:** Подтверждения, быстрые действия
- **Notifications:** Toast уведомления об операциях
- **History:** Отслеживание изменений объекта

### Статические ассеты
- **Изображения:** Логотипы, иконки, превью устройств
- **Шрифты:** Кастомные шрифты для UI
- **JavaScript библиотеки:** 
  - Chart.js (графики)
  - Leaflet (карты)
  - Select2 (выпадающие списки)
  - flatpickr (календари)

## 8. Результат анализа

### Количественные метрики
- **130+ моделей данных** для портирования
- **200+ REST API endpoints**
- **1 единый GraphQL endpoint** с полной схемой
- **~50,000 строк кода** бизнес-логики
- **13 основных модулей** с подмодулями

### Качественные характеристики
- **Сложная бизнес-логика:** Требует тщательного проектирования
- **Высокая связность:** Много межмодульных зависимостей
- **Расширяемость:** Система плагинов и custom fields
- **Производительность:** Оптимизированные запросы к БД
- **Безопасность:** Детальная система разрешений

### Риски миграции
- **Потеря функционала:** Скрипты и отчеты на Python
- **Совместимость:** Обратная совместимость API
- **Данные:** Миграция существующей БД
- **Производительность:** Go vs Django ORM производительность
- **Сообщество:** Потеря экосистемы плагинов

### Рекомендации
1. **Поэтапная миграция:** Начать с наименее зависимых модулей
2. **Двойной запуск:** Параллельная работа Django и Go
3. **API совместимость:** Сохранение контрактов для клиентов
4. **Инструменты:** Автоматическая генерация кода из схем
5. **Тестирование:** Полное покрытие тестами перед релизом

---

*Документ подготовлен для планирования миграции NetBox на Go с сохранением функциональности и архитектуры.*
