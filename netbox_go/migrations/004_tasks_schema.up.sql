-- Миграция 004: Схема для системы управления задачами с ролями, группами и видами работ
-- Эта миграция создаёт таблицы для управления задачами, группами пользователей, видами работ и назначениями

-- ============================================
-- Виды работ (Work Types)
-- ============================================
CREATE TABLE IF NOT EXISTS tasks_worktype (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tasks_worktype_name ON tasks_worktype(name);

-- ============================================
-- Группы (Groups)
-- ============================================
CREATE TYPE group_type_enum AS ENUM ('assignee', 'reviewer');

CREATE TABLE IF NOT EXISTS tasks_group (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    type group_type_enum NOT NULL,
    description TEXT,
    shift_start TIME,  -- Начало смены (для групп-смен)
    shift_end TIME,    -- Конец смены (для групп-смен)
    work_days TEXT[],  -- Дни недели работы (например: {'monday', 'tuesday'})
    created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tasks_group_name ON tasks_group(name);
CREATE INDEX IF NOT EXISTS idx_tasks_group_type ON tasks_group(type);

-- ============================================
-- Участники групп (Group Members)
-- ============================================
CREATE TABLE IF NOT EXISTS tasks_group_member (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id UUID NOT NULL REFERENCES tasks_group(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,  -- Ссылка на users_user (из account модуля)
    created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    added_by UUID NOT NULL,  -- Кто добавил участника
    UNIQUE(group_id, user_id)  -- Один пользователь может быть в группе только один раз
);

CREATE INDEX IF NOT EXISTS idx_tasks_group_member_group ON tasks_group_member(group_id);
CREATE INDEX IF NOT EXISTS idx_tasks_group_member_user ON tasks_group_member(user_id);

-- ============================================
-- Компетенции групп (Group Work Types)
-- Связь многие-ко-многим между группами и видами работ
-- ============================================
CREATE TABLE IF NOT EXISTS tasks_group_worktype (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id UUID NOT NULL REFERENCES tasks_group(id) ON DELETE CASCADE,
    work_type_id UUID NOT NULL REFERENCES tasks_worktype(id) ON DELETE CASCADE,
    created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    added_by UUID NOT NULL,  -- Кто добавил компетенцию
    UNIQUE(group_id, work_type_id)  -- Одна компетенция может быть указана только один раз
);

CREATE INDEX IF NOT EXISTS idx_tasks_group_worktype_group ON tasks_group_worktype(group_id);
CREATE INDEX IF NOT EXISTS idx_tasks_group_worktype_worktype ON tasks_group_worktype(work_type_id);

-- ============================================
-- Задачи (Tasks)
-- ============================================
CREATE TYPE task_status_enum AS ENUM (
    'draft',         -- Черновик
    'assigned',      -- Назначена
    'in_progress',   -- В работе
    'completed',     -- Выполнена
    'under_review',  -- На проверке
    'approved',      -- Принята
    'rejected',      -- Отклонена
    'cancelled'      -- Отменена
);

CREATE TABLE IF NOT EXISTS tasks_task (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(200) NOT NULL,
    description TEXT,
    work_type_id UUID NOT NULL REFERENCES tasks_worktype(id),
    status task_status_enum NOT NULL DEFAULT 'draft',
    priority INTEGER NOT NULL DEFAULT 3 CHECK (priority >= 1 AND priority <= 5),
    created_by_id UUID NOT NULL,  -- Ссылка на users_user (создатель задачи)
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    due_date TIMESTAMP WITH TIME ZONE,  -- Дедлайн
    completed_at TIMESTAMP WITH TIME ZONE,  -- Когда завершена
    reviewed_at TIMESTAMP WITH TIME ZONE,   -- Когда проверена
    review_comment TEXT  -- Комментарий проверяющего
);

CREATE INDEX IF NOT EXISTS idx_tasks_task_status ON tasks_task(status);
CREATE INDEX IF NOT EXISTS idx_tasks_task_work_type ON tasks_task(work_type_id);
CREATE INDEX IF NOT EXISTS idx_tasks_task_created_by ON tasks_task(created_by_id);
CREATE INDEX IF NOT EXISTS idx_tasks_task_due_date ON tasks_task(due_date);
CREATE INDEX IF NOT EXISTS idx_tasks_task_priority ON tasks_task(priority DESC);

-- ============================================
-- Назначения задач (Task Assignments)
-- Связывает задачи с пользователями или группами в определённых ролях
-- ============================================
CREATE TYPE task_role_enum AS ENUM ('creator', 'assignee', 'reviewer');
CREATE TYPE assignee_type_enum AS ENUM ('user', 'group');

CREATE TABLE IF NOT EXISTS tasks_task_assignment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id UUID NOT NULL REFERENCES tasks_task(id) ON DELETE CASCADE,
    role task_role_enum NOT NULL,
    assignee_type assignee_type_enum NOT NULL,
    user_id UUID,  -- Заполняется если assignee_type = 'user'
    group_id UUID,  -- Заполняется если assignee_type = 'group'
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,  -- Кто создал назначение
    actual_executor_id UUID,  -- Фактический исполнитель (если задача назначена группе)
    
    -- Проверка: если assignee_type = 'user', то user_id должен быть указан
    CONSTRAINT check_user_assignee CHECK (
        (assignee_type = 'user' AND user_id IS NOT NULL) OR
        (assignee_type = 'group' AND group_id IS NOT NULL)
    ),
    -- Проверка: если assignee_type = 'group', то group_id должен быть указан
    CONSTRAINT check_group_assignee CHECK (
        (assignee_type = 'group' AND group_id IS NOT NULL) OR
        (assignee_type = 'user' AND user_id IS NOT NULL)
    )
);

CREATE INDEX IF NOT EXISTS idx_tasks_task_assignment_task ON tasks_task_assignment(task_id);
CREATE INDEX IF NOT EXISTS idx_tasks_task_assignment_role ON tasks_task_assignment(role);
CREATE INDEX IF NOT EXISTS idx_tasks_task_assignment_user ON tasks_task_assignment(user_id);
CREATE INDEX IF NOT EXISTS idx_tasks_task_assignment_group ON tasks_task_assignment(group_id);

-- Уникальное ограничение: одна роль одного типа назначения на задачу
-- (можно иметь одного creator, одного assignee, одного reviewer)
CREATE UNIQUE INDEX IF NOT EXISTS idx_tasks_task_assignment_unique 
ON tasks_task_assignment(task_id, role);

-- ============================================
-- Комментарии к задачам (Task Comments)
-- ============================================
CREATE TABLE IF NOT EXISTS tasks_task_comment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id UUID NOT NULL REFERENCES tasks_task(id) ON DELETE CASCADE,
    author_id UUID NOT NULL,  -- Автор комментария
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tasks_task_comment_task ON tasks_task_comment(task_id);
CREATE INDEX IF NOT EXISTS idx_tasks_task_comment_author ON tasks_task_comment(author_id);

-- ============================================
-- Вложения задач (Task Attachments)
-- ============================================
CREATE TABLE IF NOT EXISTS tasks_task_attachment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id UUID NOT NULL REFERENCES tasks_task(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    file_type VARCHAR(100),  -- MIME тип файла
    uploaded_by UUID NOT NULL,
    uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tasks_task_attachment_task ON tasks_task_attachment(task_id);
CREATE INDEX IF NOT EXISTS idx_tasks_task_attachment_uploaded_by ON tasks_task_attachment(uploaded_by);

-- ============================================
-- Триггеры для обновления updated_at
-- ============================================

-- Функция для обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Триггер для tasks_worktype
DROP TRIGGER IF EXISTS update_tasks_worktype_updated_at ON tasks_worktype;
CREATE TRIGGER update_tasks_worktype_updated_at
    BEFORE UPDATE ON tasks_worktype
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Триггер для tasks_group
DROP TRIGGER IF EXISTS update_tasks_group_updated_at ON tasks_group;
CREATE TRIGGER update_tasks_group_updated_at
    BEFORE UPDATE ON tasks_group
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Триггер для tasks_task
DROP TRIGGER IF EXISTS update_tasks_task_updated_at ON tasks_task;
CREATE TRIGGER update_tasks_task_updated_at
    BEFORE UPDATE ON tasks_task
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Триггер для tasks_task_comment
DROP TRIGGER IF EXISTS update_tasks_task_comment_updated_at ON tasks_task_comment;
CREATE TRIGGER update_tasks_task_comment_updated_at
    BEFORE UPDATE ON tasks_task_comment
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- Начальные данные (опционально)
-- ============================================

-- Добавим базовые виды работ
INSERT INTO tasks_worktype (id, name, description) VALUES
    ('00000000-0000-0000-0000-000000000001', 'Физическое подключение сервера', 'Монтаж серверного оборудования в стойку, подключение кабелей питания и сети'),
    ('00000000-0000-0000-0000-000000000002', 'Программная настройка сервера', 'Установка ОС, настройка сетевого стека, установка необходимого ПО'),
    ('00000000-0000-0000-0000-000000000003', 'Диагностика оборудования', 'Выявление неисправностей оборудования, тестирование компонентов')
ON CONFLICT (id) DO NOTHING;

-- Комментарии к миграции
COMMENT ON TABLE tasks_worktype IS 'Справочник видов работ';
COMMENT ON TABLE tasks_group IS 'Группы пользователей (ответственные или проверяющие)';
COMMENT ON TABLE tasks_group_member IS 'Участники групп';
COMMENT ON TABLE tasks_group_worktype IS 'Компетенции групп (какие виды работ может выполнять группа)';
COMMENT ON TABLE tasks_task IS 'Задачи';
COMMENT ON TABLE tasks_task_assignment IS 'Назначения ролей в задачах (кто создатель, ответственный, проверяющий)';
COMMENT ON TABLE tasks_task_comment IS 'Комментарии к задачам';
COMMENT ON TABLE tasks_task_attachment IS 'Вложения к задачам';

COMMENT ON COLUMN tasks_group.type IS 'Тип группы: assignee (ответственные) или reviewer (проверяющие)';
COMMENT ON COLUMN tasks_group.shift_start IS 'Начало смены для групп-смен';
COMMENT ON COLUMN tasks_group.shift_end IS 'Конец смены для групп-смен';
COMMENT ON COLUMN tasks_group.work_days IS 'Дни недели работы группы';
COMMENT ON COLUMN tasks_task.status IS 'Статус задачи';
COMMENT ON COLUMN tasks_task.priority IS 'Приоритет задачи (1-5, где 5 - наивысший)';
COMMENT ON COLUMN tasks_task_assignment.role IS 'Роль в задаче: creator, assignee, reviewer';
COMMENT ON COLUMN tasks_task_assignment.assignee_type IS 'Тип назначения: user или group';
COMMENT ON COLUMN tasks_task_assignment.actual_executor_id IS 'Фактический исполнитель (если задача назначена группе)';
