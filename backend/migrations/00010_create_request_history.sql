-- +goose Up
-- +goose StatementBegin
CREATE TABLE request_history (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  request_id UUID NOT NULL REFERENCES sharing_requests(id) ON DELETE CASCADE,
  changed_by_user_id UUID REFERENCES users(id),
  old_status VARCHAR(50),
  new_status VARCHAR(50),
  notes TEXT,
  changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS request_history;
-- +goose StatementEnd
