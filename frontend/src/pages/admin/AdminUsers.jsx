import React from "react";
import { api } from "../../api";
import { ROLES, USER_STATUSES, formatDate } from "../../utils";

export function AdminUsers() {
  const [users, setUsers] = React.useState([]);
  const [error, setError] = React.useState("");
  const [notice, setNotice] = React.useState("");

  const load = React.useCallback(() => {
    api
      .users()
      .then((data) => setUsers(data.users || []))
      .catch((e) => setError(e.message));
  }, []);

  React.useEffect(load, [load]);

  function patch(id, changes) {
    setUsers((current) => current.map((u) => (u.id === id ? { ...u, ...changes } : u)));
  }

  async function save(user) {
    setError("");
    setNotice("");
    try {
      await api.updateUser(user.id, { role: user.role, status: user.status });
      setNotice(`Updated ${user.name || user.email}`);
    } catch (e) {
      setError(e.message);
    }
  }

  return (
    <section>
      <h1>User management</h1>
      {error && <p className="alert error">{error}</p>}
      {notice && <p className="alert success">{notice}</p>}
      <div className="table-wrap wide">
        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Email</th>
              <th>Role</th>
              <th>Status</th>
              <th>Joined</th>
              <th />
            </tr>
          </thead>
          <tbody>
            {users.map((user) => (
              <tr key={user.id}>
                <td>{user.name}</td>
                <td>{user.email}</td>
                <td>
                  <select value={user.role} onChange={(e) => patch(user.id, { role: e.target.value })}>
                    {ROLES.map((r) => (
                      <option key={r} value={r}>
                        {r}
                      </option>
                    ))}
                  </select>
                </td>
                <td>
                  <select
                    value={user.status}
                    onChange={(e) => patch(user.id, { status: e.target.value })}
                  >
                    {USER_STATUSES.map((s) => (
                      <option key={s} value={s}>
                        {s}
                      </option>
                    ))}
                  </select>
                </td>
                <td>{formatDate(user.created_at)}</td>
                <td>
                  <button className="btn" type="button" onClick={() => save(user)}>
                    Save
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
