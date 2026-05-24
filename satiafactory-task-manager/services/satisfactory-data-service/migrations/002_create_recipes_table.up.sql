CREATE TABLE IF NOT EXISTS buildings (
                                         class_name TEXT PRIMARY KEY,
                                         display_name TEXT NOT NULL,
                                         description TEXT,
                                         power_consumption FLOAT,
                                         power_consumption_exponent FLOAT
);