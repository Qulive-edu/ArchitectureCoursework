import { useEffect, useState } from "react";
import { fetchPlaces } from "../api/places";
import PlaceCard from "../components/PlaceCard";

export default function Home() {
  const [places, setPlaces] = useState([]);

  useEffect(() => {
    fetchPlaces().then(res => setPlaces(res.data));
  }, []);

  return (
    <div className="container mx-auto grid grid-cols-1 md:grid-cols-3 gap-6 py-10">
      {places.map((p) => (
        <PlaceCard key={p.id} place={p} />
      ))}
    </div>
  );
}
