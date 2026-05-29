import { writable, get, read, type Writable } from 'svelte/store';

type User = {
  username: string;
  is_admin: boolean;
};

class AuthStore {
  private tokenKey = 'jwt_token';
  private userKey = 'user_data';

  private token: Writable<string | null> = writable(null);
  private user: Writable<User | null> = writable(null);

  constructor() {
    this.loadFromStorage();
  }

  private loadFromStorage() {
    const storedToken = localStorage.getItem(this.tokenKey);
    const storedUser = localStorage.getItem(this.userKey);

    if (storedToken) {
      this.token.set(storedToken);
    }

    if (storedUser) {
      try {
        const parsedUser = JSON.parse(storedUser);
        this.user.set(parsedUser);
      } catch (e) {
        localStorage.removeItem(this.userKey);
      }
    }
  }

  getToken(): string | null {
    return read(this.token);
  }

  getUser(): User | null {
    return read(this.user);
  }

  isAuthenticated(): boolean {
    const token = get(this.token);
    return token !== null;
  }

  setToken(token: string) {
    this.token.set(token);
    localStorage.setItem(this.tokenKey, token);
  }

  setUser(user: User) {
    this.user.set(user);
    localStorage.setItem(this.userKey, JSON.stringify(user));
  }

  async login(username: string, password: string) {
    try {
      const response = await api.login(username, password);
      this.setToken(response.token);
      this.setUser(response.user);
      return { success: true };
    } catch (error: any) {
      return { success: false, error: error.message };
    }
  }

  async register(username: string, password: string, isAdmin: boolean) {
    try {
      await api.register(username, password, isAdmin);
      return { success: true };
    } catch (error: any) {
      return { success: false, error: error.message };
    }
  }

  logout() {
    this.token.set(null);
    this.user.set(null);
    localStorage.removeItem(this.tokenKey);
    localStorage.removeItem(this.userKey);
  }

  requireAuth(): boolean {
    if (!this.isAuthenticated()) {
      window.location.href = '/';
      return false;
    }
    return true;
  }
}

export const authStore = new AuthStore();
