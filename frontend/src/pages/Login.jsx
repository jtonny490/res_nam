import React from "react";
import { api } from "../api";
import { setSession } from "../state";

export function Login({ navigate }) {
  const [form, setForm] = React.useState({});
  const [error, setError] = React.useState("");

  async function submit(event) {
    event.preventDefault();
    setError("");
    try {
      setSession(await api.login(form));
      navigate("feed");
    } catch (e) {
      setError(e.message);
    }
  }

  return (
    <section className="panel card">
      <h1>Log in</h1>
      {error && <p className="alert error">{error}</p>}
      <form onSubmit={submit}>
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
          <label htmlFor="password">Password</label>
          <input
            id="password"
            type="password"
            required
            onChange={(e) => setForm({ ...form, password: e.target.value })}
          />
        </div>
        <button className="btn primary" type="submit">
          Log in
        </button>
      </form>
      <p className="meta">
        No account? <a href="#/register">Register</a>
      </p>
    </section>
  );
}
