INSERT INTO places (title, address, floor_type, description, image, "price_per_hour")
VALUES
('Парковка №1', 'ул. Ленина, 10', 'Большая парковка', 'https://avtoshkola.msk.ru/wp-content/uploads/2022/03/parkovka-elochkoi.jpg', 1500),
('Парковка №2', 'пр. Победы, 22', 'Открытая парковка', 'https://blog.idn500.ru/upload/iblock/3b3/ercbd0j26hu5aspiiy159xx09j1n4z20/ploskostnye_parkovki.jpg', 2500),
('Парковка №3', 'ул. Молодёжная, 7', 'Закрытая парковка', 'https://rosdorznakservis.ru/files/uploads/statii/parkovka-i-avtostoyanka-razbiraemsya-gde-ostavit-mashinu/parkovka.jpg', 2000);

INSERT INTO time_slots (place_id, start_time, end_time, is_available)
VALUES
(1, NOW() + INTERVAL '1 hour', NOW() + INTERVAL '2 hours', true),
(1, NOW() + INTERVAL '3 hours', NOW() + INTERVAL '4 hours', true),
(2, NOW() + INTERVAL '1 day', NOW() + INTERVAL '1 day 2 hours', true),
(3, NOW() + INTERVAL '2 days', NOW() + INTERVAL '2 days 2 hours', true);