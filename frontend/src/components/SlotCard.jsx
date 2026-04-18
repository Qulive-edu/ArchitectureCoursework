export default function SlotCard({ slot, onBook }) {
  return (
    <div className="bg-zinc-800 p-4 rounded-xl flex justify-between items-center">
      <div>
        <p>{new Date(slot.start_time).toLocaleString()}</p>
        <p className="text-zinc-400">
          {new Date(slot.end_time).toLocaleTimeString()}
        </p>
      </div>

      {slot.is_available ? (
        <button className="btn-primary" onClick={onBook}>Бронь</button>
      ) : (
        <span className="text-red-400">Занято</span>
      )}
    </div>
  );
}
