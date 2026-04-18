// frontend/src/pages/Login.jsx
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { loginRequest } from "../api/auth";
import { useUserStore } from "../store/user";

export default function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const login = useUserStore((s) => s.login);
  const navigate = useNavigate();

  const submit = async (e) => {
    e.preventDefault();
    try {
      const res = await loginRequest(email, password);
      login(res.data.user, res.data.token);
      navigate("/", { replace: true });
    } catch (err) {
      alert("Ошибка входа: " + (err.response?.data?.error || "Проверьте email и пароль"));
    }
  };

  return (
    <div className="max-w-sm mx-auto py-10">
      <form onSubmit={submit} className="bg-zinc-800 p-6 rounded-xl space-y-4">
        <h1 className="text-xl font-bold">Вход</h1>
        <input
          className="w-full p-2 rounded bg-zinc-900"
          placeholder="Email"
          onChange={(e) => setEmail(e.target.value)}
        />
        <input
          type="password"
          className="w-full p-2 rounded bg-zinc-900"
          placeholder="Пароль"
          onChange={(e) => setPassword(e.target.value)}
        />
        <button className="btn-primary w-full">Войти</button>
      </form>
    </div>
  );
}