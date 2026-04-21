-- +goose Up
-- +goose StatementBegin
CREATE TABLE utilization_metrics (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id),
  department_id UUID NOT NULL REFERENCES departments(id),
  category_id UUID NOT NULL REFERENCES equipment_categories(id),
  date DATE NOT NULL,
  total_hours_available NUMERIC(10, 2) NOT NULL DEFAULT 0,
  total_hours_in_use NUMERIC(10, 2) NOT NULL DEFAULT 0,
  total_hours_in_maintenance NUMERIC(10, 2) NOT NULL DEFAULT 0,
  sharing_requests_sent INTEGER NOT NULL DEFAULT 0,
  sharing_requests_received INTEGER NOT NULL DEFAULT 0,
  sharing_requests_fulfilled INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(tenant_id, department_id, category_id, date)
);

CREATE INDEX idx_utilization_date ON utilization_metrics(tenant_id, date);
CREATE INDEX idx_utilization_dept ON utilization_metrics(department_id, date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS utilization_metrics;
-- +goose StatementEnd
