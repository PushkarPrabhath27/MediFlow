-- +goose Up
-- +goose StatementBegin
CREATE TABLE department_min_stock (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id),
  department_id UUID NOT NULL REFERENCES departments(id),
  category_id UUID NOT NULL REFERENCES equipment_categories(id),
  minimum_count INTEGER NOT NULL DEFAULT 1,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(department_id, category_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS department_min_stock;
-- +goose StatementEnd
