import { api } from '$lib/api/client';

export type User = {
  username: string;
  is_admin: boolean;
};

const tokenKey = 'jwt_token';
const userKey = 'user_data';

let token: string | null = null;
let user: User | null = null;

const listeners: Set<() => void> = new Set();

function loadFromStorage() {
  if (typeof localStorage === 'undefined') {
    return;
  }
  const storedToken = localStorage.getItem(tokenKey);
  const storedUser = localStorage.getItem(userKey);
  if (storedToken && storedUser) {
    try {
      token = storedToken;
      user = JSON.parse(storedUser);
    } catch {
      localStorage.removeItem(tokenKey);
      localStorage.removeItem(userKey);
      token = null;
      user = null;
    }
  } else {
    token = null;
    user = null;
  }
  notifyListeners();
}

function notifyListeners() {
  for (const listener of listeners) {
    listener();
  }
}

function subscribe(fn: () => void) {
  listeners.add(fn);
  return () => {
    listeners.delete(fn);
  };
}

if (typeof window !== 'undefined') {
  loadFromStorage();
  window.addEventListener('storage', () => loadFromStorage());
}

function setAuth(t: string, u: User) {
  token = t;
  user = u;
  if (typeof localStorage !== 'undefined') {
    localStorage.setItem(tokenKey, t);
    localStorage.setItem(userKey, JSON.stringify(u));
  }
  notifyListeners();
}

function clearAuth() {
  token = null;
  user = null;
  if (typeof localStorage !== 'undefined') {
    localStorage.removeItem(tokenKey);
    localStorage.removeItem(userKey);
  }
  notifyListeners();
}

export const auth = {
  getToken: () => token,
  getUser: () => user,
  isAuthenticated: () => token !== null,
  setToken(t: string) { setAuth(t, user!); },
  setUser(u: User) { setAuth(token!, u); },
  async login(username: string, password: string) {
    try {
      const response = await api.login(username, password);
      setAuth(response.token, response.user);
      return { success: true };
    } catch (error: any) {
      return { success: false, error: error.message };
    }
  },
  logout() { clearAuth(); },
  requireAuth() {
    if (!auth.isAuthenticated()) {
      window.location.href = '/';
      return false;
    }
    return true;
  },
  subscribe,
};

export const { getToken, getUser, isAuthenticated, setToken, setUser, login, logout, requireAuth } = auth;
export { subscribe };
