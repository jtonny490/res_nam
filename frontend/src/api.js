import { getSession } from "./state.js";

function query(params = {}) {
  const entries = Object.entries(params).filter(
    ([, value]) => value !== undefined && value !== null && value !== ""
  );
  return entries.length ? `?${new URLSearchParams(entries)}` : "";
}

export async function request(path, { method = "GET", body, auth = true } = {}) {
  const headers = {};
  const token = getSession()?.token;
  if (auth && token) headers.Authorization = `Bearer ${token}`;
  let payload = body;
  if (body && !(body instanceof FormData)) {
    headers["Content-Type"] = "application/json";
    payload = JSON.stringify(body);
  }
  const response = await fetch(path, { method, headers, body: payload });
  const data = response.status === 204 ? null : await response.json().catch(() => ({}));
  if (!response.ok) {
    throw new Error((data && data.error) || `Request failed (${response.status})`);
  }
  return data;
}

export const api = {
  register: (body) => request("/api/auth/register", { method: "POST", body, auth: false }),
  login: (body) => request("/api/auth/login", { method: "POST", body, auth: false }),

  reports: (params) => request(`/api/reports${query(params)}`),
  report: (id) => request(`/api/reports/${id}`),
  createReport: (data, photo) => {
    if (!photo) return request("/api/reports", { method: "POST", body: data });
    const form = new FormData();
    Object.entries(data).forEach(([key, value]) => form.append(key, value));
    form.append("photo", photo);
    return request("/api/reports", { method: "POST", body: form });
  },
  updateReport: (id, body) => request(`/api/reports/${id}`, { method: "PATCH", body }),
  deleteReport: (id) => request(`/api/reports/${id}`, { method: "DELETE" }),

  comments: (id) => request(`/api/reports/${id}/comments`),
  comment: (id, body) => request(`/api/reports/${id}/comments`, { method: "POST", body }),
  like: (id) => request(`/api/reports/${id}/like`, { method: "POST" }),
  unlike: (id) => request(`/api/reports/${id}/like`, { method: "DELETE" }),

  applyAuthority: (body) => request("/api/authority-requests", { method: "POST", body }),
  authorityRequests: () => request("/api/authority-requests"),
  reviewAuthority: (id, status) =>
    request(`/api/authority-requests/${id}`, { method: "PATCH", body: { status } }),

  mapPins: () => request("/api/map/reports"),
  assessment: (id) => request(`/api/reports/${id}/assessment`),

  users: () => request("/api/admin/users"),
  updateUser: (id, body) => request(`/api/admin/users/${id}`, { method: "PATCH", body }),
  analytics: () => request("/api/admin/analytics"),
};
