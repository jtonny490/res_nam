export const PILOT = { minLat: -1.6, maxLat: 0.7, minLng: 31.5, maxLng: 35.0 };

export const CATEGORIES = ["water", "air", "waste", "deforestation", "other"];
export const STATUSES = ["open", "investigating", "stale", "resolved"];
export const ROLES = ["public", "authority", "admin"];
export const USER_STATUSES = ["active", "pending", "banned"];

const SEVERITY_LABELS = {
  1: "Minor",
  2: "Low",
  3: "Moderate",
  4: "High",
  5: "Severe",
};

export function titleCase(value = "") {
  return value
    .split(/[\s_-]+/)
    .filter(Boolean)
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(" ");
}

export function severityLabel(level) {
  return SEVERITY_LABELS[level] ?? "Unknown";
}

export function formatDate(value) {
  if (!value) return "";
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? "" : date.toLocaleString();
}

export function isWithinPilot(latitude, longitude) {
  const lat = Number(latitude);
  const lng = Number(longitude);
  if (Number.isNaN(lat) || Number.isNaN(lng)) return false;
  return lat >= PILOT.minLat && lat <= PILOT.maxLat && lng >= PILOT.minLng && lng <= PILOT.maxLng;
}

export function classNames(...values) {
  return values.filter(Boolean).join(" ");
}
