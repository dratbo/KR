CREATE TABLE IF NOT EXISTS recipe_ingredients (
                                                  recipe_class_name TEXT REFERENCES recipes(class_name) ON DELETE CASCADE,
    item_class_name TEXT REFERENCES items(class_name) ON DELETE CASCADE,
    amount FLOAT,
    PRIMARY KEY (recipe_class_name, item_class_name)
    );