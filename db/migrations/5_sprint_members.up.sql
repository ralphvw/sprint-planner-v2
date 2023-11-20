CREATE TABLE IF NOT EXISTS sprint_members(
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users (id) ON DELETE CASCADE,
  sprint_id INT REFERENCES sprints (id) ON DELETE CASCADE,
  designation VARCHAR,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
)
