import React from 'react'
import { createRoot } from 'react-dom/client'
import './styles.css'

function App() {
  const [reports, setReports] = React.useState([])
  const [error, setError] = React.useState('')
  React.useEffect(() => { fetch('/api/reports').then(r => r.json()).then(x => setReports(x.reports || [])).catch(() => setError('Unable to load reports')) }, [])
  return <main><header><h1>RES NAM</h1><p>Lake Victoria pollution reports</p></header><section><h2>Community feed</h2>{error && <p>{error}</p>}{reports.length ? reports.map(r => <article key={r.id}><h3>{r.title}</h3><p>{r.category} · severity {r.severity} · {r.status}</p><p>{r.description}</p></article>) : <p>No reports yet.</p>}</section></main>
}
createRoot(document.getElementById('root')).render(<App />)
