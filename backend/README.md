# Backend for Sports Places Booking (Go)

This is a minimal backend implementing places, time slots and bookings.
Tech: Go, PostgreSQL, chi router.

Run with:
1. docker compose up --build
2. Apply migrations placed in ./migrations are mounted into Postgres at container start

API:
- POST /auth/register {name,email,password}
- POST /auth/login {email,password}
- GET /places
- GET /places/{id}
- GET /places/{id}/slots
- POST /bookings (Authorization: Bearer <token>) body {place_id, slot_id}
- GET /bookings/my (Authorization: Bearer <token>)

