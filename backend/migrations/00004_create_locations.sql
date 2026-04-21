-- +goose Up
-- +goose StatementBegin
CREATE TABLE locations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id),
  department_id UUID NOT NULL REFERENCES departments(id),
  name VARCHAR(255) NOT NULL,
  type VARCHAR(50), -- room, bay, corridor, storage
  floor VARCHAR(50),
  building VARCHAR(100),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS locations;
-- +goose StatementEnd
