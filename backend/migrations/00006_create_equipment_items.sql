-- +goose Up
-- +goose StatementBegin
CREATE TABLE equipment_items (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id),
  category_id UUID NOT NULL REFERENCES equipment_categories(id),
  department_id UUID REFERENCES departments(id),
  current_location_id UUID REFERENCES locations(id),
  name VARCHAR(255) NOT NULL,
  model VARCHAR(255),
  manufacturer VARCHAR(255),
  serial_number VARCHAR(255),
  asset_tag VARCHAR(100),
  qr_code VARCHAR(255) UNIQUE,
  status VARCHAR(50) NOT NULL DEFAULT 'available',
  -- status: available, in_use, reserved, in_maintenance, in_transit, missing, decommissioned
  purchase_date DATE,
  purchase_cost DECIMAL(12,2),
  notes TEXT,
  is_shared BOOLEAN NOT NULL DEFAULT true, -- whether this item participates in sharing
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_equipment_tenant_status ON equipment_items(tenant_id, status);
CREATE INDEX idx_equipment_department ON equipment_items(department_id);
CREATE INDEX idx_equipment_category ON equipment_items(category_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS equipment_items;
-- +goose StatementEnd
