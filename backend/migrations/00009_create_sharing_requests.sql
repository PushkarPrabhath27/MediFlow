-- +goose Up
-- +goose StatementBegin
CREATE TABLE sharing_requests (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id),
  requesting_dept_id UUID NOT NULL REFERENCES departments(id),
  requesting_user_id UUID NOT NULL REFERENCES users(id),
  source_dept_id UUID REFERENCES departments(id),
  equipment_id UUID REFERENCES equipment_items(id),
  category_id UUID NOT NULL REFERENCES equipment_categories(id),
  quantity_needed INTEGER NOT NULL DEFAULT 1,
  urgency VARCHAR(20) NOT NULL DEFAULT 'normal', -- low, normal, high, emergency
  reason TEXT,
  status VARCHAR(50) NOT NULL DEFAULT 'pending',
  -- status: pending, matched, approved, declined, in_transit, active, return_pending, completed, cancelled
  needed_by TIMESTAMPTZ,
  expected_return_at TIMESTAMPTZ,
  matched_at TIMESTAMPTZ,
  approved_at TIMESTAMPTZ,
  approved_by_user_id UUID REFERENCES users(id),
  declined_reason TEXT,
  handoff_confirmed_by_source BOOLEAN DEFAULT false,
  handoff_confirmed_by_requester BOOLEAN DEFAULT false,
  handed_off_at TIMESTAMPTZ,
  return_confirmed_by_source BOOLEAN DEFAULT false,
  return_confirmed_by_requester BOOLEAN DEFAULT false,
  returned_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_requests_tenant ON sharing_requests(tenant_id);
CREATE INDEX idx_requests_status ON sharing_requests(status);
CREATE INDEX idx_requests_requesting_dept ON sharing_requests(requesting_dept_id);
CREATE INDEX idx_requests_source_dept ON sharing_requests(source_dept_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sharing_requests;
-- +goose StatementEnd
