import { auth } from '$lib/stores/auth';

const API_BASE = import.meta.env.VITE_API_BASE_URL;

class ApiClient {
  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${API_BASE}${endpoint}`;
    
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...(options.headers as Record<string, string> || {}),
    };

    const token = auth.getToken();
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(url, {
      ...options,
      headers,
    });

    if (response.status === 401) {
      auth.logout();
      window.location.href = '/';
      throw new Error('Authentication required');
    }

    const contentType = response.headers.get('content-type');
    if (contentType && contentType.includes('application/json')) {
      const data = await response.json();
      if (!response.ok) {
        throw new Error(data.message || 'API request failed');
      }
      return data;
    }

    if (!response.ok) {
      throw new Error('API request failed');
    }

    return response as unknown as T;
  }

  async login(username: string, password: string): Promise<AuthResponse> {
    const response = await fetch(`${API_BASE}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password }),
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.message || 'Login failed');
    }

    return data.data;
  }

  async listGames(): Promise<GameResponse> {
    const response = await this.request<CommonResponse & { data: GameResponse }>('/games');
    return response.data;
  }

  async startGame(gameName: string): Promise<ControlResponse> {
    const response = await this.request<CommonResponse & { data: ControlResponse }>(`/games/${gameName}/start`, {
      method: 'POST',
    });
    return response.data;
  }

  async stopGame(gameName: string): Promise<ControlResponse> {
    const response = await this.request<CommonResponse & { data: ControlResponse }>(`/games/${gameName}/stop`, {
      method: 'POST',
    });
    return response.data;
  }

  async restartGame(gameName: string): Promise<ControlResponse> {
    const response = await this.request<CommonResponse & { data: ControlResponse }>(`/games/${gameName}/restart`, {
      method: 'POST',
    });
    return response.data;
  }

  async getGameStatus(gameName: string): Promise<ControlResponse> {
    const response = await this.request<CommonResponse & { data: ControlResponse }>(`/games/${gameName}/status`, {
      method: 'GET',
    });
    return response.data;
  }

  streamLogs(gameName: string): EventSource {
    const token = auth.getToken();
    let url = `${API_BASE}/games/${gameName}/logs`;
    if (token) {
      url += `?token=${encodeURIComponent(token)}`;
    }

    return new EventSource(url);
  }
}

type CommonResponse = {
  success: boolean;
  message: string;
  data?: any;
};

type AuthResponse = {
  token: string;
  user: {
    username: string;
    is_admin: boolean;
  };
};

type GameResponse = string[];

type ControlResponse = {
  status: string;
  game: string;
};

export const api = new ApiClient();
