CREATE TABLE IF NOT EXISTS schematic_unlocks (
                                                 schematic_class_name TEXT REFERENCES schematics(class_name) ON DELETE CASCADE,
    unlock_type TEXT,
    unlock_data TEXT,
    PRIMARY KEY (schematic_class_name, unlock_type, unlock_data)
    );