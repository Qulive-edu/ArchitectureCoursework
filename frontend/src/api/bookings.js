import api from "./axios";

export const createBooking = (place_id, slot_id) =>
  api.post("/bookings", { place_id, slot_id });

export const fetchMyBookings = () => api.get("/bookings/my");
