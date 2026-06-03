import { writable } from 'svelte/store';
import { api } from '$lib/api/client';

export type User = {
  username: string;
  is_admin: boolean;
};

type AuthState = {
  token: string | null;
  user: User | null;
};

const tokenKey = 'jwt_token';
const userKey = 'user_data';

function loadFromStorage(): AuthState {
  if (typeof localStorage === 'undefined') {
    return { token: null, user: null };
  }
  
  const storedToken = localStorage.getItem(tokenKey);
  const storedUser = localStorage.getItem(userKey);

  if (storedToken && storedUser) {
    try {
      return {
        token: storedToken,
        user: JSON.parse(storedUser),
      };
    } catch (e) {
      localStorage.removeItem(tokenKey);
      localStorage.removeItem(userKey);
    }
  }

  return { token: null, user: null };
}

// Create the store
const store = writable<AuthState>(loadFromStorage());

// Get the store value
function getStoreValue() {
  let value: AuthState | null = null;
  store.subscribe((v) => {
    value = v;
  })();
  return value;
}

// Auth functions
export const auth = {
  getToken(): string | null {
    return getStoreValue()?.token || null;
  },

  getUser(): User | null {
    return getStoreValue()?.user || null;
  },

  isAuthenticated(): boolean {
    return auth.getToken() !== null;
  },

  setToken(token: string) {
    store.update((state) => ({ ...state, token }));
    localStorage.setItem(tokenKey, token);
  },

  setUser(user: User) {
    store.update((state) => ({ ...state, user }));
    localStorage.setItem(userKey, JSON.stringify(user));
  },

  async login(username: string, password: string) {
    try {
      const response = await api.login(username, password);
      auth.setToken(response.token);
      auth.setUser(response.user);
      return { success: true };
    } catch (error: any) {
      return { success: false, error: error.message };
    }
  },

  logout() {
    store.set({ token: null, user: null });
    localStorage.removeItem(tokenKey);
    localStorage.removeItem(userKey);
  },

  requireAuth(): boolean {
    if (!auth.isAuthenticated()) {
      window.location.href = '/';
      return false;
    }
    return true;
  },
};

// Export both the store and the auth functions
export { store as authStore };
export const { getToken, getUser, isAuthenticated, setToken, setUser, login, logout, requireAuth } = auth;
