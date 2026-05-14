// API base URL — overridden at build time via REACT_APP_API_URL so the same
// build can target localhost in dev and a real host in production.
export const API_BASE =
  process.env.REACT_APP_API_URL?.replace(/\/+$/, '') || 'http://localhost:8080';

// JWT storage helpers (localStorage — fine for a dev/demo project).
const TOKEN_KEY = 'tally.token';
const USER_KEY = 'tally.user';

export function getToken() {
  return localStorage.getItem(TOKEN_KEY);
}
export function setToken(t) {
  if (t) localStorage.setItem(TOKEN_KEY, t);
  else localStorage.removeItem(TOKEN_KEY);
}
export function clearToken() {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(USER_KEY);
}
export function getStoredUser() {
  const raw = localStorage.getItem(USER_KEY);
  if (!raw) return null;
  try {
    return JSON.parse(raw);
  } catch {
    return null;
  }
}
export function setStoredUser(u) {
  if (u) localStorage.setItem(USER_KEY, JSON.stringify(u));
  else localStorage.removeItem(USER_KEY);
}
