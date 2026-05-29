import { writable } from 'svelte/store';

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

export type { User };

export const authStore = writable<AuthState>(loadFromStorage());

export const authManager = {
  getToken(): string | null {
    let token: string | null = null;
    authStore.subscribe((value) => {
      token = value.token;
    })();
    return token;
  },

  getUser(): User | null {
    let user: User | null = null;
    authStore.subscribe((value) => {
      user = value.user;
    })();
    return user;
  },

  isAuthenticated(): boolean {
    return authManager.getToken() !== null;
  },

  setToken(token: string) {
    authStore.update((state) => ({ ...state, token }));
    localStorage.setItem(tokenKey, token);
  },

  setUser(user: User) {
    authStore.update((state) => ({ ...state, user }));
    localStorage.setItem(userKey, JSON.stringify(user));
  },

  async login(username: string, password: string) {
    try {
      const response = await api.login(username, password);
      authManager.setToken(response.token);
      authManager.setUser(response.user);
      return { success: true };
    } catch (error: any) {
      return { success: false, error: error.message };
    }
  },

  async register(username: string, password: string, isAdmin: boolean) {
    try {
      await api.register(username, password, isAdmin);
      return { success: true };
    } catch (error: any) {
      return { success: false, error: error.message };
    }
  },

  logout() {
    authStore.set({ token: null, user: null });
    localStorage.removeItem(tokenKey);
    localStorage.removeItem(userKey);
  },

  requireAuth(): boolean {
    if (!authManager.isAuthenticated()) {
      window.location.href = '/';
      return false;
    }
    return true;
  },
};
