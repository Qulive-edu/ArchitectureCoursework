// frontend/src/components/Navbar.jsx
import { Link } from "react-router-dom";
import { useUserStore } from "../store/user";

export default function Navbar() {
  const { user, logout } = useUserStore();

  return (
    <nav className="bg-zinc-900 p-4 flex justify-between items-center shadow">
      <Link to="/" className="text-xl font-bold">SportRent</Link>

      <div className="flex items-center gap-4">
        {!user ? (
          <>
            <Link to="/login">Войти</Link>
            <Link to="/register" className="btn-primary">Регистрация</Link>
          </>
        ) : (
          <>
            <span>Привет, {user.name}!</span>
            <Link to="/my-bookings">Мои брони</Link>
            <button onClick={logout} className="btn-secondary">Выйти</button>
          </>
        )}
      </div>
    </nav>
  );
}