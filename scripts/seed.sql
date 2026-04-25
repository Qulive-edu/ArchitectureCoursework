-- Используйте price_per_hour вместо price
INSERT INTO places (title, address, floor_type, description, image, "price_per_hour")
VALUES
('Спортзал №1', 'ул. Ленина, 10', 'Большой спортивный зал для игр', 'https://volleygrad.ru/wp-content/uploads/2021/03/photo_2022-11-09_21-14-33.jpg', 1500),
('Футбольное поле', 'пр. Победы, 22', 'Открытое футбольное поле', 'https://findsport.ru/userfiles/images/620x365/1654861015_b97bdb3f0aad94547f4bdbd9d6ef73c4.jpg', 2500),
('Баскетбольная площадка', 'ул. Молодёжная, 7', 'Закрытая площадка для баскетбола', 'https://terball.ru/sites/default/files/styles/crop_1000_800/public/2020-05/basketball-rent-gallery-6.jpg?itok=5LPm2h2C', 2000);

INSERT INTO time_slots (place_id, start_time, end_time, is_available)
VALUES
(1, NOW() + INTERVAL '1 hour', NOW() + INTERVAL '2 hours', true),
(1, NOW() + INTERVAL '3 hours', NOW() + INTERVAL '4 hours', true),
(2, NOW() + INTERVAL '1 day', NOW() + INTERVAL '1 day 2 hours', true),
(3, NOW() + INTERVAL '2 days', NOW() + INTERVAL '2 days 2 hours', true);