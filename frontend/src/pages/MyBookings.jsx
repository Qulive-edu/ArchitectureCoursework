import { useEffect, useState } from "react";
import { fetchMyBookings } from "../api/bookings";

export default function MyBookings() {
  const [list, setList] = useState([]);

  useEffect(() => {
    fetchMyBookings().then((res) => setList(res.data));
  }, []);

  return (
    <div className="max-w-2xl mx-auto py-10">
      <h1 className="text-2xl font-bold mb-4">Мои брони</h1>

      {list.map((b) => (
        <div key={b.id} className="bg-zinc-800 p-4 rounded-xl mb-4">
          <p><strong>Бронь #{b.id}</strong></p>
          <p>Площадка: {b.place_id}</p>
          <p>Слот: {b.slot_id}</p>
          <p className="text-zinc-400">{new Date(b.created_at).toLocaleString()}</p>
        </div>
      ))}
    </div>
  );
}
