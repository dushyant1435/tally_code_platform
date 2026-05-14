import React, { createContext, useCallback, useContext, useEffect, useState } from 'react';
import { api, ApiError } from '../api';
import {
  getStoredUser,
  setStoredUser,
  setToken,
  clearToken,
  getToken,
} from '../config';

const AuthCtx = createContext(null);

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(() => getStoredUser());
  const [loading, setLoading] = useState(!!getToken() && !getStoredUser());

  // If we have a token but no cached user, fetch /me on mount.
  useEffect(() => {
    let cancelled = false;
    const token = getToken();
    if (!token || user) {
      setLoading(false);
      return;
    }
    api
      .get('/api/v1/auth/me')
      .then((u) => {
        if (cancelled) return;
        setUser(u);
        setStoredUser(u);
      })
      .catch(() => {
        clearToken();
        if (!cancelled) setUser(null);
      })
      .finally(() => !cancelled && setLoading(false));
    return () => {
      cancelled = true;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const login = useCallback(async (usernameOrEmail, password) => {
    const data = await api.post('/api/v1/auth/login', {
      username: usernameOrEmail,
      password,
    });
    setToken(data.token);
    setStoredUser(data.user);
    setUser(data.user);
    return data.user;
  }, []);

  const signup = useCallback(async (username, email, password) => {
    const data = await api.post('/api/v1/auth/signup', {
      username,
      email,
      password,
    });
    setToken(data.token);
    setStoredUser(data.user);
    setUser(data.user);
    return data.user;
  }, []);

  const logout = useCallback(() => {
    clearToken();
    setUser(null);
  }, []);

  const value = {
    user,
    loading,
    isAuthenticated: !!user,
    isAdmin: user?.role === 'admin',
    login,
    signup,
    logout,
  };
  return <AuthCtx.Provider value={value}>{children}</AuthCtx.Provider>;
};

export function useAuth() {
  const ctx = useContext(AuthCtx);
  if (!ctx) throw new Error('useAuth must be used inside <AuthProvider>');
  return ctx;
}

export { ApiError };
