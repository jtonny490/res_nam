import React from "react";
import { api } from "../api";
import { CATEGORIES, PILOT, isWithinPilot } from "../utils";

export function NewReport({ navigate }) {
  const [form, setForm] = React.useState({ category: "water", severity: 3 });
  const [photo, setPhoto] = React.useState(null);
  const [error, setError] = React.useState("");
  const [busy, setBusy] = React.useState(false);

  function update(field, value) {
    setForm((current) => ({ ...current, [field]: value }));
  }

  function useMyLocation() {
    if (!navigator.geolocation) {
      setError("Geolocation is not available in this browser");
      return;
    }
    navigator.geolocation.getCurrentPosition(
      (position) => {
        update("latitude", position.coords.latitude);
        update("longitude", position.coords.longitude);
      },
      () => setError("Could not read your location; enter coordinates manually")
    );
  }

  async function submit(event) {
    event.preventDefault();
    setError("");
    if (!isWithinPilot(form.latitude, form.longitude)) {
      setError(
        `Coordinates must fall within the pilot area (lat ${PILOT.minLat}..${PILOT.maxLat}, lng ${PILOT.minLng}..${PILOT.maxLng})`
      );
      return;
    }
    setBusy(true);
    try {
      await api.createReport(
        {
          title: form.title,
          description: form.description || "",
          category: form.category,
          severity: Number(form.severity),
          latitude: Number(form.latitude),
          longitude: Number(form.longitude),
        },
        photo
      );
      navigate("feed");
    } catch (e) {
      setError(e.message);
    } finally {
      setBusy(false);
    }
  }

  return (
    <section className="panel card">
      <h1>Submit a pollution report</h1>
      {error && <p className="alert error">{error}</p>}
      <form onSubmit={submit}>
        <div>
          <label htmlFor="title">Title</label>
          <input id="title" required onChange={(e) => update("title", e.target.value)} />
        </div>
        <div>
          <label htmlFor="category">Category</label>
          <select
            id="category"
            value={form.category}
            onChange={(e) => update("category", e.target.value)}
          >
            {CATEGORIES.map((c) => (
              <option key={c} value={c}>
                {c}
              </option>
            ))}
          </select>
        </div>
        <div>
          <label htmlFor="severity">Severity (1-5)</label>
          <input
            id="severity"
            type="range"
            min="1"
            max="5"
            value={form.severity}
            onChange={(e) => update("severity", e.target.value)}
          />
          <span className="meta">Level {form.severity}</span>
        </div>
        <div>
          <label htmlFor="description">Description</label>
          <textarea
            id="description"
            rows="4"
            onChange={(e) => update("description", e.target.value)}
          />
        </div>
        <div className="grid-2">
          <div>
            <label htmlFor="latitude">Latitude</label>
            <input
              id="latitude"
              type="number"
              step="any"
              required
              value={form.latitude ?? ""}
              onChange={(e) => update("latitude", e.target.value)}
            />
          </div>
          <div>
            <label htmlFor="longitude">Longitude</label>
            <input
              id="longitude"
              type="number"
              step="any"
              required
              value={form.longitude ?? ""}
              onChange={(e) => update("longitude", e.target.value)}
            />
          </div>
        </div>
        <button className="btn" type="button" onClick={useMyLocation}>
          Use my location
        </button>
        <div>
          <label htmlFor="photo">Photo</label>
          <input
            id="photo"
            type="file"
            accept="image/*"
            onChange={(e) => setPhoto(e.target.files?.[0] || null)}
          />
        </div>
        <button className="btn primary" type="submit" disabled={busy}>
          {busy ? "Submitting…" : "Submit report"}
        </button>
      </form>
    </section>
  );
}
