export default function PlaceCard({ place }) {
  return (
    <div className="bg-zinc-800 rounded-xl p-5 shadow hover:scale-[1.01] transition">
      {place.image && (
        <img
          src={place.image}
          className="rounded-xl mb-3 h-40 w-full object-cover"
        />
      )}

      <h2 className="text-xl font-bold">{place.title}</h2>
      <p className="text-zinc-400">{place.address}</p>
      <p className="mt-2 text-green-400">{place.price_per_hour} ₽ / час</p>

      <a
        className="btn-primary mt-4 inline-block"
        href={`/place/${place.id}`}
      >
        Подробнее
      </a>
    </div>
  );
}
