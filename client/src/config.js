// API base URL — overridden at build time with REACT_APP_API_URL so the same
// build can target localhost in dev and a real host in production.
export const API_BASE =
  process.env.REACT_APP_API_URL?.replace(/\/+$/, '') || 'http://localhost:8080';

// Hard-coded for now until auth lands.
export const CURRENT_USER_ID = 123;
