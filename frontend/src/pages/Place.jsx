// frontend/src/pages/Place.jsx
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { fetchPlace, fetchSlots } from "../api/places";
import SlotCard from "../components/SlotCard";
import { createBooking } from "../api/bookings";

export default function Place() {
  const { id } = useParams();
  const placeId = Number(id);

  const [place, setPlace] = useState(null);
  const [slots, setSlots] = useState([]); // по умолчанию — пустой массив
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (!placeId || isNaN(placeId)) {
      setError("Неверный ID площадки");
      setLoading(false);
      return;
    }

    const load = async () => {
      try {
        setLoading(true);
        const [placeRes, slotsRes] = await Promise.all([
          fetchPlace(placeId),
          fetchSlots(placeId)
        ]);
        setPlace(placeRes.data);
        setSlots(Array.isArray(slotsRes.data) ? slotsRes.data : []);
      } catch (err) {
        console.error("Ошибка загрузки площадки или слотов:", err);
        setError("Не удалось загрузить данные");
      } finally {
        setLoading(false);
      }
    };

    load();
  }, [placeId]);

  const book = async (slotId) => {
    try {
      await createBooking(placeId, slotId);
      alert("Забронировано!");
    } catch (err) {
      alert("Ошибка бронирования: " + (err.response?.data?.error || err.message));
    }
  };

  if (loading) return <div className="py-10 text-center">Загрузка...</div>;
  if (error) return <div className="py-10 text-center text-red-400">Ошибка: {error}</div>;
  if (!place) return <div className="py-10 text-center">Площадка не найдена</div>;

  return (
    <div className="max-w-3xl mx-auto py-10">
      <h1 className="text-3xl font-bold">{place.title}</h1>
      {place.image && (
        <img
          src={place.image || "https://placehold.co/600x300?text=No+image"}
          className="rounded-xl my-4 w-full object-cover h-64"
          alt={place.title}
        />
      )}

      <h2 className="text-2xl font-bold mt-6 mb-4">Доступные слоты:</h2>
      {slots.length === 0 ? (
        <p className="text-zinc-400">Нет доступных слотов</p>
      ) : (
        <div className="space-y-4">
          {slots.map((slot) => (
            <SlotCard key={slot.id} slot={slot} onBook={() => book(slot.id)} />
          ))}
        </div>
      )}
    </div>
  );
}