import React from "react";
import { api } from "../../api";
import { formatDate } from "../../utils";

export function AdminAuthorityRequests() {
  const [requests, setRequests] = React.useState([]);
  const [error, setError] = React.useState("");

  const load = React.useCallback(() => {
    api
      .authorityRequests()
      .then((data) => setRequests(data.requests || []))
      .catch((e) => setError(e.message));
  }, []);

  React.useEffect(load, [load]);

  async function review(id, status) {
    setError("");
    try {
      await api.reviewAuthority(id, status);
      load();
    } catch (e) {
      setError(e.message);
    }
  }

  return (
    <section>
      <h1>Authority requests</h1>
      {error && <p className="alert error">{error}</p>}
      {requests.length === 0 ? (
        <p className="muted">No pending requests.</p>
      ) : (
        <div className="stack">
          {requests.map((request) => (
            <article key={request.id} className="card stack">
              <header className="row">
                <h2>{request.organization_name}</h2>
                <span className="meta">{formatDate(request.created_at)}</span>
              </header>
              {request.justification && <p>{request.justification}</p>}
              <div className="row">
                <button className="btn primary" type="button" onClick={() => review(request.id, "approved")}>
                  Approve
                </button>
                <button className="btn danger" type="button" onClick={() => review(request.id, "rejected")}>
                  Reject
                </button>
              </div>
            </article>
          ))}
        </div>
      )}
    </section>
  );
}
