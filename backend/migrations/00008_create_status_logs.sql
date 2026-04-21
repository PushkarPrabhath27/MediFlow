-- +goose Up
-- +goose StatementBegin
CREATE TABLE equipment_status_logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  equipment_id UUID NOT NULL REFERENCES equipment_items(id),
  tenant_id UUID NOT NULL REFERENCES tenants(id),
  old_status VARCHAR(50),
  new_status VARCHAR(50) NOT NULL,
  old_department_id UUID,
  new_department_id UUID,
  old_location_id UUID,
  new_location_id UUID,
  changed_by_user_id UUID REFERENCES users(id),
  reason TEXT,
  changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_status_logs_equipment ON equipment_status_logs(equipment_id);
CREATE INDEX idx_status_logs_changed_at ON equipment_status_logs(changed_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS equipment_status_logs;
-- +goose StatementEnd
