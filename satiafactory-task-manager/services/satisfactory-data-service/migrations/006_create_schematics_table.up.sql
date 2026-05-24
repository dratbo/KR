CREATE TABLE IF NOT EXISTS schematics (
                                          class_name TEXT PRIMARY KEY,
                                          display_name TEXT NOT NULL,
                                          description TEXT,
                                          schematic_type TEXT,
                                          time_to_complete FLOAT
);