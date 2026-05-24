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

CREATE TABLE IF NOT EXISTS buildings (
                                         class_name TEXT PRIMARY KEY,
                                         display_name TEXT NOT NULL,
                                         description TEXT,
                                         power_consumption FLOAT,
                                         power_consumption_exponent FLOAT
);

CREATE TABLE IF NOT EXISTS recipes (
                                       class_name TEXT PRIMARY KEY,
                                       display_name TEXT NOT NULL,
                                       produced_in JSONB,
                                       duration FLOAT,
                                       manufactoring_menu_priority INT
);

CREATE TABLE IF NOT EXISTS recipe_ingredients (
                                                  recipe_class_name TEXT REFERENCES recipes(class_name) ON DELETE CASCADE,
    item_class_name TEXT REFERENCES items(class_name) ON DELETE CASCADE,
    amount FLOAT,
    PRIMARY KEY (recipe_class_name, item_class_name)
    );

CREATE TABLE IF NOT EXISTS recipe_products (
                                               recipe_class_name TEXT REFERENCES recipes(class_name) ON DELETE CASCADE,
    item_class_name TEXT REFERENCES items(class_name) ON DELETE CASCADE,
    amount FLOAT,
    PRIMARY KEY (recipe_class_name, item_class_name)
    );

CREATE TABLE IF NOT EXISTS schematics (
                                          class_name TEXT PRIMARY KEY,
                                          display_name TEXT NOT NULL,
                                          description TEXT,
                                          schematic_type TEXT,
                                          time_to_complete FLOAT
);

CREATE TABLE IF NOT EXISTS schematic_costs (
                                               schematic_class_name TEXT REFERENCES schematics(class_name) ON DELETE CASCADE,
    item_class_name TEXT REFERENCES items(class_name) ON DELETE CASCADE,
    amount FLOAT,
    PRIMARY KEY (schematic_class_name, item_class_name)
    );

CREATE TABLE IF NOT EXISTS schematic_unlocks (
                                                 schematic_class_name TEXT REFERENCES schematics(class_name) ON DELETE CASCADE,
    unlock_type TEXT,
    unlock_data TEXT,
    PRIMARY KEY (schematic_class_name, unlock_type, unlock_data)
    );