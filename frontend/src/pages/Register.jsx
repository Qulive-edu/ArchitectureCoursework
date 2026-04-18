// frontend/src/pages/Register.jsx
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { registerRequest } from "../api/auth";
import { useUserStore } from "../store/user";

export default function Register() {
  const [form, setForm] = useState({ name: "", email: "", password: "" });
  const login = useUserStore((s) => s.login);
  const navigate = useNavigate();

  const submit = async (e) => {
    e.preventDefault();
    try {
      const res = await registerRequest(form.name, form.email, form.password);
      login(res.data.user, res.data.token);
      navigate("/", { replace: true });
    } catch (err) {
      alert("Ошибка регистрации: " + (err.response?.data?.error || "Попробуйте другое имя/email"));
    }
  };

  return (
    <div className="max-w-sm mx-auto py-10">
      <form className="bg-zinc-800 p-6 rounded-xl space-y-4" onSubmit={submit}>
        <h1 className="text-xl font-bold">Регистрация</h1>
        <input
          className="w-full p-2 rounded bg-zinc-900"
          placeholder="Имя"
          onChange={(e) => setForm({ ...form, name: e.target.value })}
        />
        <input
          className="w-full p-2 rounded bg-zinc-900"
          placeholder="Email"
          onChange={(e) => setForm({ ...form, email: e.target.value })}
        />
        <input
          type="password"
          className="w-full p-2 rounded bg-zinc-900"
          placeholder="Пароль"
          onChange={(e) => setForm({ ...form, password: e.target.value })}
        />
        <button className="btn-primary w-full">Создать аккаунт</button>
      </form>
    </div>
  );
}