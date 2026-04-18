-- init schema for places booking

CREATE TABLE IF NOT EXISTS places (
    id SERIAL PRIMARY KEY,
    title VARCHAR(127) NOT NULL,
    address VARCHAR(255) NOT NULL,
    floor_type VARCHAR(255) NOT NULL,
    description TEXT,
    image VARCHAR(255),
    price_per_hour INT NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(127) NOT NULL,
    email VARCHAR(127) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS time_slots (
    id SERIAL PRIMARY KEY,
    place_id INT REFERENCES places(id),
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    is_available BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS bookings (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    place_id INT REFERENCES places(id),
    slot_id INT REFERENCES time_slots(id),
    status VARCHAR(64) DEFAULT 'confirmed',
    created_at TIMESTAMP DEFAULT NOW()
);
