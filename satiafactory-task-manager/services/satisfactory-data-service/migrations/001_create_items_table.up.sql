CREATE TABLE IF NOT EXISTS items (
                                     class_name TEXT PRIMARY KEY,
                                     display_name TEXT NOT NULL,
                                     description TEXT,
                                     stack_size INT,
                                     energy_value FLOAT,
                                     form TEXT,
                                     small_icon TEXT,
                                     big_icon TEXT,
                                     categories JSONB
);