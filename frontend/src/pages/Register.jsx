import React from "react";
import { api } from "../api";
import { setSession } from "../state";

export function Register({ navigate }) {
  const [form, setForm] = React.useState({});
  const [error, setError] = React.useState("");

  async function submit(event) {
    event.preventDefault();
    setError("");
    try {
      setSession(await api.register(form));
      navigate("feed");
    } catch (e) {
      setError(e.message);
    }
  }

  return (
    <section className="panel card">
      <h1>Create account</h1>
      {error && <p className="alert error">{error}</p>}
      <form onSubmit={submit}>
        <div>
          <label htmlFor="name">Name</label>
          <input
            id="name"
            required
            onChange={(e) => setForm({ ...form, name: e.target.value })}
          />
        </div>
        <div>
          <label htmlFor="email">Email</label>
          <input
            id="email"
            type="email"
            required
            onChange={(e) => setForm({ ...form, email: e.target.value })}
          />
        </div>
        <div>
          <label htmlFor="password">Password (min 8 characters)</label>
          <input
            id="password"
            type="password"
            minLength="8"
            required
            onChange={(e) => setForm({ ...form, password: e.target.value })}
          />
        </div>
        <button className="btn primary" type="submit">
          Register
        </button>
      </form>
      <p className="meta">
        Already registered? <a href="#/login">Log in</a>
      </p>
    </section>
  );
}
