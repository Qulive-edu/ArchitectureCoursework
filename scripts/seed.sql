-- Используйте price_per_hour вместо price
INSERT INTO places (title, address, floor_type, description, image, "price_per_hour")
VALUES
('Спортзал №1', 'ул. Ленина, 10', 'Паркет', 'Большой спортивный зал для игр', 'https://volleygrad.ru/wp-content/uploads/2021/03/photo_2022-11-09_21-14-33.jpg', 1500),
('Футбольное поле', 'пр. Победы, 22', 'Искусственная трава', 'Открытое футбольное поле', 'https://findsport.ru/userfiles/images/620x365/1654861015_b97bdb3f0aad94547f4bdbd9d6ef73c4.jpg', 2500),
('Баскетбольная площадка', 'ул. Молодёжная, 7', 'Паркет', 'Закрытая площадка для баскетбола', 'https://terball.ru/sites/default/files/styles/crop_1000_800/public/2020-05/basketball-rent-gallery-6.jpg?itok=5LPm2h2C', 2000);

INSERT INTO time_slots (place_id, start_time, end_time, is_available)
VALUES
(1, '2026-01-01 10:00', '2026-01-01 12:00', true),
(1, '2026-01-02 14:00', '2026-01-02 15:00', true),
(2, '2026-01-03 16:00', '2026-01-03 18:00', true),
(3, '2026-02-01 09:00', '2026-02-01 11:00', true);