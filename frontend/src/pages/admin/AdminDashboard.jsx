import React from "react";
import { api } from "../../api";
import { titleCase } from "../../utils";

export function AdminDashboard() {
  const [data, setData] = React.useState(null);
  const [error, setError] = React.useState("");

  React.useEffect(() => {
    api
      .analytics()
      .then(setData)
      .catch((e) => setError(e.message));
  }, []);

  if (error) return <p className="alert error">{error}</p>;
  if (!data) return <p className="muted">Loading…</p>;

  return (
    <section>
      <h1>Admin dashboard</h1>
      <div className="stats-grid">
        <div className="stat">
          <div className="value">{data.total_reports}</div>
          <div className="meta">Total reports</div>
        </div>
        <div className="stat">
          <div className="value">{data.active_users}</div>
          <div className="meta">Active users</div>
        </div>
      </div>

      <h2>Reports by category</h2>
      <div className="table-wrap">
        <table>
          <thead>
            <tr>
              <th>Category</th>
              <th>Count</th>
            </tr>
          </thead>
          <tbody>
            {(data.by_category || []).map((row) => (
              <tr key={row.category}>
                <td>{titleCase(row.category)}</td>
                <td>{row.count}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <h2>Reports by status</h2>
      <div className="table-wrap">
        <table>
          <thead>
            <tr>
              <th>Status</th>
              <th>Count</th>
            </tr>
          </thead>
          <tbody>
            {(data.by_status || []).map((row) => (
              <tr key={row.status}>
                <td>{titleCase(row.status)}</td>
                <td>{row.count}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <p className="meta">
        Manage <a href="#/admin/users">users</a>, <a href="#/admin/authority-requests">authority requests</a>, or{" "}
        <a href="#/admin/moderation">content</a>.
      </p>
    </section>
  );
}
