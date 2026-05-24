CREATE TABLE IF NOT EXISTS tasks (
                                     id BIGSERIAL PRIMARY KEY,
                                     user_id BIGINT NOT NULL,
                                     title TEXT NOT NULL,
                                     description TEXT,
                                     status TEXT DEFAULT 'pending',
                                     created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    );

CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks(user_id);