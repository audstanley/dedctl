type CommonResponse = {
  success: boolean;
  message: string;
  data?: any;
};

type LoginRequest = {
  username: string;
  password: string;
};

type RegisterRequest = {
  username: string;
  password: string;
  is_admin: boolean;
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
  private baseUrl: string = '/api';

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

  async register(username: string, password: string, isAdmin: boolean): Promise<CommonResponse> {
    return this.request<CommonResponse>('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ username, password, is_admin: isAdmin }),
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

  streamLogs(gameName: string): EventSource {
    const token = localStorage.getItem('jwt_token');
    const url = `${this.baseUrl}/games/${gameName}/logs`;
    
    const headers: Record<string, string> = {};
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const eventSource = new EventSource(url, {
      withCredentials: false,
      headers,
    });

    return eventSource;
  }
}

export const api = new ApiClient();
