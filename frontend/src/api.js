const token = () => JSON.parse(localStorage.getItem('session') || 'null')?.token
export async function request(path, options = {}) {
  const headers = { 'Content-Type': 'application/json', ...(options.headers || {}) }
  if (token()) headers.Authorization = `Bearer ${token()}`
  const response = await fetch(path, { ...options, headers })
  const data = await response.json().catch(() => ({}))
  if (!response.ok) throw new Error(data.error || 'Request failed')
  return data
}
export const api = {
  reports: () => request('/api/reports'),
  report: id => request(`/api/reports/${id}`),
  login: body => request('/api/auth/login', { method: 'POST', body: JSON.stringify(body) }),
  register: body => request('/api/auth/register', { method: 'POST', body: JSON.stringify(body) }),
  createReport: body => request('/api/reports', { method: 'POST', body: JSON.stringify(body) }),
  comment: (id, body) => request(`/api/reports/${id}/comments`, { method: 'POST', body: JSON.stringify(body) }),
  like: id => request(`/api/reports/${id}/like`, { method: 'POST' })
}
