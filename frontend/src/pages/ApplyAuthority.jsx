import React from "react";
import { api } from "../api";

export function ApplyAuthority() {
  const [form, setForm] = React.useState({});
  const [error, setError] = React.useState("");
  const [done, setDone] = React.useState(false);

  async function submit(event) {
    event.preventDefault();
    setError("");
    try {
      await api.applyAuthority(form);
      setDone(true);
    } catch (e) {
      setError(e.message);
    }
  }

  if (done) {
    return (
      <section className="panel card">
        <h1>Application submitted</h1>
        <p className="alert success">
          Your authority request is pending review. An administrator will approve or reject it.
        </p>
      </section>
    );
  }

  return (
    <section className="panel card">
      <h1>Apply for authority access</h1>
      <p className="muted">
        Maritime authorities can review reports and update their status. Submit your organization
        details for administrator approval.
      </p>
      {error && <p className="alert error">{error}</p>}
      <form onSubmit={submit}>
        <div>
          <label htmlFor="organization_name">Organization name</label>
          <input
            id="organization_name"
            required
            onChange={(e) => setForm({ ...form, organization_name: e.target.value })}
          />
        </div>
        <div>
          <label htmlFor="justification">Justification</label>
          <textarea
            id="justification"
            rows="4"
            onChange={(e) => setForm({ ...form, justification: e.target.value })}
          />
        </div>
        <button className="btn primary" type="submit">
          Submit application
        </button>
      </form>
    </section>
  );
}
