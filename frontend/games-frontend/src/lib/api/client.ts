type CommonResponse = {
  success: boolean;
  message: string;
  data?: any;
};

type LoginRequest = {
  username: string;
  password: string;
};

type User = {
  username: string;
  is_admin: boolean;
};

type AuthResponse = {
  token: string;
  user: User;
};

type GameResponse = string[];

type ControlResponse = {
  status: string;
  game: string;
};

type LogEntry = {
  timestamp: string;
  message: string;
};

class ApiClient {
  private get baseUrl(): string {
    if (typeof window !== 'undefined' && window.location.hostname !== 'localhost' && window.location.hostname !== '127.0.0.1') {
      return `http://${window.location.hostname}:8080`;
    }
    return '/api';
  }

  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    const token = localStorage.getItem('jwt_token');
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(url, {
      ...options,
      headers,
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.message || 'API request failed');
    }

    return data;
  }

  async login(username: string, password: string): Promise<AuthResponse> {
    return this.request<CommonResponse & { data: AuthResponse }>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    });
  }

  async listGames(): Promise<GameResponse> {
    return this.request<GameResponse>('/games');
  }

  async startGame(gameName: string): Promise<ControlResponse> {
    return this.request<CommonResponse & { data: ControlResponse }>(`/games/${gameName}/start`, {
      method: 'POST',
    });
  }

  async stopGame(gameName: string): Promise<ControlResponse> {
    return this.request<CommonResponse & { data: ControlResponse }>(`/games/${gameName}/stop`, {
      method: 'POST',
    });
  }

  async restartGame(gameName: string): Promise<ControlResponse> {
    return this.request<CommonResponse & { data: ControlResponse }>(`/games/${gameName}/restart`, {
      method: 'POST',
    });
  }

  async getGameStatus(gameName: string): Promise<ControlResponse> {
    return this.request<CommonResponse & { data: ControlResponse }>(`/games/${gameName}/status`, {
      method: 'GET',
    });
  }

  streamLogs(gameName: string): EventSource {
    const token = localStorage.getItem('jwt_token');
    let url = `${this.baseUrl}/games/${gameName}/logs`;
    if (token) {
      url += `?token=${encodeURIComponent(token)}`;
    }

    return new EventSource(url);
  }
}

export const api = new ApiClient();
