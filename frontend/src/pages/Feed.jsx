import React from "react";
import { api } from "../api";
import { ReportCard } from "../components/ReportCard";
import { CATEGORIES, STATUSES } from "../utils";

export function Feed() {
  const [filters, setFilters] = React.useState({ category: "", status: "" });
  const [page, setPage] = React.useState(1);
  const [data, setData] = React.useState({ reports: [], total: 0, limit: 20 });
  const [error, setError] = React.useState("");

  React.useEffect(() => {
    setError("");
    api
      .reports({ ...filters, page })
      .then((result) =>
        setData((current) => ({
          ...current,
          reports: Array.isArray(result?.reports) ? result.reports : [],
          total: result?.total ?? 0,
          limit: result?.limit ?? current.limit,
        }))
      )
      .catch((e) => setError(e.message));
  }, [filters, page]);

  const reports = Array.isArray(data.reports) ? data.reports : [];
  const limit = data.limit || 20;
  const totalPages = Math.max(1, Math.ceil((data.total || 0) / limit));

  function updateFilter(event) {
    setPage(1);
    setFilters({ ...filters, [event.target.name]: event.target.value });
  }

  return (
    <section>
      <h1>Community reports</h1>
      <div className="filters">
        <div>
          <label htmlFor="category">Category</label>
          <select id="category" name="category" value={filters.category} onChange={updateFilter}>
            <option value="">All categories</option>
            {CATEGORIES.map((c) => (
              <option key={c} value={c}>
                {c}
              </option>
            ))}
          </select>
        </div>
        <div>
          <label htmlFor="status">Status</label>
          <select id="status" name="status" value={filters.status} onChange={updateFilter}>
            <option value="">All statuses</option>
            {STATUSES.map((s) => (
              <option key={s} value={s}>
                {s}
              </option>
            ))}
          </select>
        </div>
        <span className="meta">{data.total || 0} report(s)</span>
      </div>

      {error && <p className="alert error">{error}</p>}
      {reports.length === 0 && !error ? (
        <p className="muted">No reports match these filters.</p>
      ) : (
        <div className="stack">
          {reports.map((report) => (
            <ReportCard key={report.id} report={report} />
          ))}
        </div>
      )}

      <div className="pagination">
        <button
          className="btn"
          type="button"
          disabled={page <= 1}
          onClick={() => setPage(page - 1)}
        >
          Previous
        </button>
        <span className="meta">
          Page {page} of {totalPages}
        </span>
        <button
          className="btn"
          type="button"
          disabled={page >= totalPages}
          onClick={() => setPage(page + 1)}
        >
          Next
        </button>
      </div>
    </section>
  );
}
