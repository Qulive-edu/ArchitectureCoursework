import api from "./axios";

export const fetchPlaces = () => api.get("/places");
export const fetchPlace = (id) => api.get(`/places/${id}`);
export const fetchSlots = (id) => api.get(`/places/${id}/slots`);
