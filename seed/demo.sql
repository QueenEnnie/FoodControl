INSERT INTO products (name, description, quantity, unit, expiration_date)
VALUES
    ('milk', 'Whole milk', 1.5, 'l', '2026-07-20'),
    ('flour', 'All-purpose flour', 1000, 'g', '2026-12-31'),
    ('eggs', 'Large eggs', 10, 'pcs', '2026-07-05'),
    ('butter', 'Unsalted butter', 200, 'g', '2026-08-01'),
    ('sugar', 'White sugar', 500, 'g', '2027-01-15'),
    ('salt', 'Table salt', 300, 'g', NULL),
    ('black pepper', 'Ground black pepper', 100, 'g', NULL),
    ('olive oil', 'Extra virgin olive oil', 500, 'ml', '2027-03-01'),
    ('rice', 'Long grain rice', 300, 'g', '2027-02-10'),
    ('pasta', 'Penne pasta', 400, 'g', '2027-04-12'),
    ('chicken breast', 'Fresh chicken breast', 600, 'g', '2026-06-25'),
    ('ground beef', 'Fresh ground beef', 300, 'g', '2026-06-23'),
    ('tomato', 'Fresh tomatoes', 6, 'pcs', '2026-06-24'),
    ('cucumber', 'Fresh cucumbers', 3, 'pcs', '2026-06-24'),
    ('onion', 'Yellow onions', 5, 'pcs', '2026-07-10'),
    ('garlic', 'Garlic cloves', 8, 'pcs', '2026-07-15'),
    ('potato', 'White potatoes', 1000, 'g', '2026-07-15'),
    ('carrot', 'Fresh carrots', 500, 'g', '2026-07-10'),
    ('cheese', 'Cheddar cheese, expired demo item', 150, 'g', '2026-05-01'),
    ('yogurt', 'Greek yogurt', 400, 'g', '2026-06-28'),
    ('oats', 'Rolled oats', 500, 'g', '2027-01-01'),
    ('banana', 'Ripe bananas', 4, 'pcs', '2026-06-22'),
    ('apple', 'Red apples', 5, 'pcs', '2026-07-01'),
    ('lettuce', 'Romaine lettuce', 1, 'pcs', '2026-06-23'),
    ('tuna', 'Canned tuna', 2, 'can', '2027-08-01'),
    ('bread', 'Sliced bread loaf', 1, 'pcs', '2026-06-22');

INSERT INTO recipes (name, description)
VALUES
    ('Pancakes', 'Classic breakfast pancakes'),
    ('Omelette', 'Simple butter omelette'),
    ('Chicken Rice Bowl', 'Chicken with rice and garlic'),
    ('Beef Pasta', 'Pasta with beef and tomato sauce'),
    ('Cheese Pasta', 'Creamy pasta that demonstrates expired cheese handling'),
    ('Vegetable Salad', 'Fresh vegetable salad'),
    ('Tuna Sandwich', 'Quick tuna sandwich'),
    ('Apple Oatmeal', 'Warm oatmeal with apple'),
    ('Banana Yogurt Bowl', 'Yogurt bowl with banana and oats'),
    ('Potato Soup', 'Simple vegetable soup'),
    ('Roasted Chicken With Potatoes', 'Oven-style chicken and potatoes'),
    ('Fruit Salad', 'Fruit salad with yogurt'),
    ('Rice Pudding', 'Dessert that requires more rice than available'),
    ('Grilled Cheese Sandwich', 'Sandwich that demonstrates expired cheese handling'),
    ('Big Breakfast Plate', 'Breakfast that requires more eggs than available'),
    ('Avocado Toast', 'Toast that demonstrates a missing product');

INSERT INTO recipe_ingredients (recipe_id, product_name, required_quantity, unit)
SELECT r.id, ingredient.product_name, ingredient.required_quantity, ingredient.unit
FROM recipes r
JOIN (
    VALUES
        ('Pancakes', 'milk', 0.5, 'l'),
        ('Pancakes', 'flour', 300, 'g'),
        ('Pancakes', 'eggs', 2, 'pcs'),
        ('Pancakes', 'sugar', 30, 'g'),
        ('Pancakes', 'butter', 30, 'g'),

        ('Omelette', 'eggs', 3, 'pcs'),
        ('Omelette', 'milk', 0.05, 'l'),
        ('Omelette', 'butter', 20, 'g'),
        ('Omelette', 'salt', 2, 'g'),
        ('Omelette', 'black pepper', 1, 'g'),

        ('Chicken Rice Bowl', 'chicken breast', 300, 'g'),
        ('Chicken Rice Bowl', 'rice', 200, 'g'),
        ('Chicken Rice Bowl', 'olive oil', 20, 'ml'),
        ('Chicken Rice Bowl', 'garlic', 2, 'pcs'),

        ('Beef Pasta', 'pasta', 250, 'g'),
        ('Beef Pasta', 'ground beef', 250, 'g'),
        ('Beef Pasta', 'tomato', 3, 'pcs'),
        ('Beef Pasta', 'onion', 1, 'pcs'),
        ('Beef Pasta', 'garlic', 2, 'pcs'),
        ('Beef Pasta', 'olive oil', 20, 'ml'),

        ('Cheese Pasta', 'pasta', 200, 'g'),
        ('Cheese Pasta', 'cheese', 100, 'g'),
        ('Cheese Pasta', 'butter', 30, 'g'),

        ('Vegetable Salad', 'tomato', 2, 'pcs'),
        ('Vegetable Salad', 'cucumber', 1, 'pcs'),
        ('Vegetable Salad', 'lettuce', 1, 'pcs'),
        ('Vegetable Salad', 'olive oil', 15, 'ml'),
        ('Vegetable Salad', 'salt', 2, 'g'),

        ('Tuna Sandwich', 'bread', 1, 'pcs'),
        ('Tuna Sandwich', 'tuna', 1, 'can'),
        ('Tuna Sandwich', 'cucumber', 1, 'pcs'),
        ('Tuna Sandwich', 'lettuce', 1, 'pcs'),

        ('Apple Oatmeal', 'oats', 80, 'g'),
        ('Apple Oatmeal', 'milk', 0.25, 'l'),
        ('Apple Oatmeal', 'apple', 1, 'pcs'),
        ('Apple Oatmeal', 'sugar', 10, 'g'),

        ('Banana Yogurt Bowl', 'banana', 2, 'pcs'),
        ('Banana Yogurt Bowl', 'yogurt', 200, 'g'),
        ('Banana Yogurt Bowl', 'oats', 50, 'g'),

        ('Potato Soup', 'potato', 500, 'g'),
        ('Potato Soup', 'carrot', 200, 'g'),
        ('Potato Soup', 'onion', 1, 'pcs'),
        ('Potato Soup', 'garlic', 1, 'pcs'),
        ('Potato Soup', 'salt', 3, 'g'),
        ('Potato Soup', 'black pepper', 1, 'g'),

        ('Roasted Chicken With Potatoes', 'chicken breast', 500, 'g'),
        ('Roasted Chicken With Potatoes', 'potato', 700, 'g'),
        ('Roasted Chicken With Potatoes', 'olive oil', 30, 'ml'),
        ('Roasted Chicken With Potatoes', 'garlic', 2, 'pcs'),
        ('Roasted Chicken With Potatoes', 'salt', 3, 'g'),
        ('Roasted Chicken With Potatoes', 'black pepper', 1, 'g'),

        ('Fruit Salad', 'apple', 2, 'pcs'),
        ('Fruit Salad', 'banana', 2, 'pcs'),
        ('Fruit Salad', 'yogurt', 150, 'g'),

        ('Rice Pudding', 'rice', 500, 'g'),
        ('Rice Pudding', 'milk', 0.5, 'l'),
        ('Rice Pudding', 'sugar', 50, 'g'),

        ('Grilled Cheese Sandwich', 'bread', 1, 'pcs'),
        ('Grilled Cheese Sandwich', 'cheese', 100, 'g'),
        ('Grilled Cheese Sandwich', 'butter', 20, 'g'),

        ('Big Breakfast Plate', 'eggs', 12, 'pcs'),
        ('Big Breakfast Plate', 'bread', 1, 'pcs'),
        ('Big Breakfast Plate', 'butter', 20, 'g'),
        ('Big Breakfast Plate', 'tomato', 2, 'pcs'),

        ('Avocado Toast', 'bread', 1, 'pcs'),
        ('Avocado Toast', 'avocado', 1, 'pcs'),
        ('Avocado Toast', 'eggs', 1, 'pcs')
) AS ingredient(recipe_name, product_name, required_quantity, unit)
    ON ingredient.recipe_name = r.name;
