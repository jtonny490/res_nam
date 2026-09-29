import React from "react";
import { api } from "../api";
import { CommentThread } from "../components/CommentThread";
import { CategoryTag, SeverityBadge, StatusBadge } from "../components/Badges";
import { canEdit, isAuthenticated, userId } from "../state";
import { CATEGORIES, formatDate } from "../utils";

export function ReportDetail({ id, navigate }) {
  const [report, setReport] = React.useState(null);
  const [comments, setComments] = React.useState([]);
  const [body, setBody] = React.useState("");
  const [error, setError] = React.useState("");
  const [editing, setEditing] = React.useState(false);
  const [form, setForm] = React.useState({});
  const [assessment, setAssessment] = React.useState(null);

  const load = React.useCallback(() => {
    api
      .report(id)
      .then((data) => {
        setReport(data);
        setForm({
          title: data.title,
          description: data.description || "",
          category: data.category,
          severity: data.severity,
        });
      })
      .catch((e) => setError(e.message));
    api.comments(id).then((data) => setComments(data.comments || []));
  }, [id]);

  React.useEffect(load, [load]);

  if (error) return <p className="alert error">{error}</p>;
  if (!report) return <p className="muted">Loading…</p>;

  const liked = (report.likes || []).some((like) => like.user_id === userId());

  async function toggleLike() {
    try {
      if (liked) await api.unlike(id);
      else await api.like(id);
      load();
    } catch (e) {
      setError(e.message);
    }
  }

  async function addComment(event) {
    event.preventDefault();
    try {
      await api.comment(id, { body });
      setBody("");
      load();
    } catch (e) {
      setError(e.message);
    }
  }

  async function saveEdit(event) {
    event.preventDefault();
    try {
      await api.updateReport(id, { ...form, severity: Number(form.severity) });
      setEditing(false);
      load();
    } catch (e) {
      setError(e.message);
    }
  }

  async function remove() {
    if (!window.confirm("Delete this report permanently?")) return;
    try {
      await api.deleteReport(id);
      navigate("feed");
    } catch (e) {
      setError(e.message);
    }
  }

  async function loadAssessment() {
    try {
      setAssessment(await api.assessment(id));
    } catch (e) {
      setError(e.message);
    }
  }

  return (
    <section className="stack">
      <article className="card stack">
        <div className="row">
          <StatusBadge status={report.status} />
          <CategoryTag category={report.category} />
          <SeverityBadge severity={report.severity} />
          <span className="meta">{formatDate(report.created_at)}</span>
        </div>
        {editing ? (
          <form className="stack" onSubmit={saveEdit}>
            <div>
              <label htmlFor="title">Title</label>
              <input
                id="title"
                value={form.title}
                onChange={(e) => setForm({ ...form, title: e.target.value })}
                required
              />
            </div>
            <div>
              <label htmlFor="category">Category</label>
              <select
                id="category"
                value={form.category}
                onChange={(e) => setForm({ ...form, category: e.target.value })}
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
                type="number"
                min="1"
                max="5"
                value={form.severity}
                onChange={(e) => setForm({ ...form, severity: e.target.value })}
              />
            </div>
            <div>
              <label htmlFor="description">Description</label>
              <textarea
                id="description"
                rows="4"
                value={form.description}
                onChange={(e) => setForm({ ...form, description: e.target.value })}
              />
            </div>
            <div className="row">
              <button className="btn primary" type="submit">
                Save
              </button>
              <button className="btn" type="button" onClick={() => setEditing(false)}>
                Cancel
              </button>
            </div>
          </form>
        ) : (
          <>
            <h1>{report.title}</h1>
            <p>{report.description}</p>
            {report.photo_url && <img className="photo-preview" src={report.photo_url} alt="" />}
            <div className="row">
              {isAuthenticated() && (
                <button className="btn" type="button" onClick={toggleLike}>
                  {liked ? "Unlike" : "Like"} ({(report.likes || []).length})
                </button>
              )}
              <button className="btn" type="button" onClick={loadAssessment}>
                Environmental assessment
              </button>
              {canEdit(report) && (
                <button className="btn" type="button" onClick={() => setEditing(true)}>
                  Edit
                </button>
              )}
              {canEdit(report) && (
                <button className="btn danger" type="button" onClick={remove}>
                  Delete
                </button>
              )}
            </div>
            {assessment && (
              <p className="alert success">
                Risk: <strong>{assessment.risk}</strong> — {assessment.summary}
              </p>
            )}
          </>
        )}
      </article>

      <section className="card">
        <h2>Comments</h2>
        <CommentThread comments={comments} />
        {isAuthenticated() ? (
          <form className="row" onSubmit={addComment}>
            <input
              value={body}
              onChange={(e) => setBody(e.target.value)}
              placeholder="Add a comment"
              required
            />
            <button className="btn primary" type="submit">
              Comment
            </button>
          </form>
        ) : (
          <p className="muted">Log in to comment.</p>
        )}
      </section>
    </section>
  );
}
