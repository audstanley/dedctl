import { auth } from '$lib/stores/auth';

class ApiClient {
  private get baseUrl(): string {
    if (typeof window !== 'undefined') {
      return `http://${window.location.hostname}:8080`;
    }
    return 'http://localhost:8080';
  }

  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    
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
    const response = await fetch(`${this.baseUrl}/auth/login`, {
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
    return this.request<GameResponse>('/games');
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
    let url = `${this.baseUrl}/games/${gameName}/logs`;
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
