-- ============================================================
-- Schema: {{FEATURE_NAME}}
-- Generated from: design/er-diagram.md
-- ============================================================

-- [NEW] ENTITY_A
CREATE TABLE entity_a (
    id          UUID            PRIMARY KEY DEFAULT gen_random_uuid() COMMENT 'Primary key',
    name        VARCHAR(255)    NOT NULL COMMENT 'Display name',
    created_at  TIMESTAMPTZ     NOT NULL DEFAULT now() COMMENT 'Creation timestamp',
    updated_at  TIMESTAMPTZ     NOT NULL DEFAULT now() COMMENT 'Last update timestamp'
) COMMENT 'Entity A description';

-- [NEW] ENTITY_B
CREATE TABLE entity_b (
    id          UUID            PRIMARY KEY DEFAULT gen_random_uuid() COMMENT 'Primary key',
    entity_a_id UUID            NOT NULL COMMENT 'FK to entity_a',
    status      VARCHAR(50)     NOT NULL DEFAULT 'active' COMMENT 'Current status: active|inactive|archived',
    created_at  TIMESTAMPTZ     NOT NULL DEFAULT now() COMMENT 'Creation timestamp',
    updated_at  TIMESTAMPTZ     NOT NULL DEFAULT now() COMMENT 'Last update timestamp',

    CONSTRAINT fk_entity_b_entity_a FOREIGN KEY (entity_a_id) REFERENCES entity_a(id) ON DELETE CASCADE
) COMMENT 'Entity B description';

-- [MODIFIED] existing_table: add new_column
-- ALTER TABLE existing_table
--     ADD COLUMN new_column VARCHAR(100) COMMENT 'New column description';

-- ============================================================
-- Indexes
-- ============================================================
CREATE INDEX idx_entity_b_entity_a_id ON entity_b(entity_a_id);
CREATE INDEX idx_entity_b_status ON entity_b(status);
