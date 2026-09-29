import React from "react";
import { api } from "../../api";
import { CategoryTag, SeverityBadge, StatusBadge } from "../../components/Badges";
import { formatDate } from "../../utils";

export function AdminModeration() {
  const [reports, setReports] = React.useState([]);
  const [error, setError] = React.useState("");

  const load = React.useCallback(() => {
    api
      .reports({ limit: 100 })
      .then((data) => setReports(data.reports || []))
      .catch((e) => setError(e.message));
  }, []);

  React.useEffect(load, [load]);

  async function remove(id) {
    if (!window.confirm("Delete this report permanently?")) return;
    setError("");
    try {
      await api.deleteReport(id);
      load();
    } catch (e) {
      setError(e.message);
    }
  }

  return (
    <section>
      <h1>Content moderation</h1>
      {error && <p className="alert error">{error}</p>}
      <div className="table-wrap wide">
        <table>
          <thead>
            <tr>
              <th>Title</th>
              <th>Category</th>
              <th>Severity</th>
              <th>Status</th>
              <th>Created</th>
              <th />
            </tr>
          </thead>
          <tbody>
            {reports.map((report) => (
              <tr key={report.id}>
                <td>
                  <a href={`#/reports/${report.id}`}>{report.title}</a>
                </td>
                <td>
                  <CategoryTag category={report.category} />
                </td>
                <td>
                  <SeverityBadge severity={report.severity} />
                </td>
                <td>
                  <StatusBadge status={report.status} />
                </td>
                <td>{formatDate(report.created_at)}</td>
                <td>
                  <button className="btn danger" type="button" onClick={() => remove(report.id)}>
                    Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </section>
  );
}
