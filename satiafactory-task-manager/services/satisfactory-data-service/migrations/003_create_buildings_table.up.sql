CREATE TABLE IF NOT EXISTS recipes (
                                       class_name TEXT PRIMARY KEY,
                                       display_name TEXT NOT NULL,
                                       produced_in JSONB,
                                       duration FLOAT,
                                       manufactoring_menu_priority INT
);