const SESSION_KEY = "session";

function storage() {
  return typeof localStorage !== "undefined" ? localStorage : null;
}

export function getSession() {
  const store = storage();
  if (!store) return null;
  try {
    const raw = store.getItem(SESSION_KEY);
    return raw ? JSON.parse(raw) : null;
  } catch {
    return null;
  }
}

export function setSession(session) {
  storage()?.setItem(SESSION_KEY, JSON.stringify(session));
}

export function clearSession() {
  storage()?.removeItem(SESSION_KEY);
}

export function currentUser() {
  return getSession()?.user ?? null;
}

export function currentRole() {
  return currentUser()?.role ?? "guest";
}

export function isAuthenticated() {
  return Boolean(getSession()?.token);
}

export function isAdmin() {
  return currentRole() === "admin";
}

export function isAuthority() {
  return currentRole() === "authority";
}

export function isPublic() {
  return currentRole() === "public";
}

export function canAccessMap() {
  return isAdmin() || isAuthority();
}

export function canModerate() {
  return isAdmin();
}

export function canManageAuthority() {
  return isAdmin();
}

export function userId() {
  return currentUser()?.id ?? 0;
}

export function isOwner(resource) {
  return Boolean(resource) && resource.user_id === userId();
}

export function canEdit(resource) {
  return isOwner(resource) || isAdmin();
}
