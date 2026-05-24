CREATE TABLE IF NOT EXISTS schematic_costs (
                                               schematic_class_name TEXT REFERENCES schematics(class_name) ON DELETE CASCADE,
    item_class_name TEXT REFERENCES items(class_name) ON DELETE CASCADE,
    amount FLOAT,
    PRIMARY KEY (schematic_class_name, item_class_name)
    );