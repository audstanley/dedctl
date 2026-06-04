import { writable, get } from 'svelte/store';
import { goto } from '$app/navigation';
import { auth } from '$lib/stores/auth';
import { api } from '$lib/api/client';

type GameInfo = {
  name: string;
  app_id: number;
  order: number;
  has_image: boolean;
};

type ServerInfo = {
  main_image: string;
  icon: string;
};

type GameStatus = 'active' | 'inactive' | 'not-found' | string;

class GamesStore {
  private games = writable<GameInfo[]>([]);
  private statuses = writable<Record<string, GameStatus>>({});

  async fetchGames(): Promise<{ success: boolean; games?: GameInfo[]; error?: string }> {
    try {
      const games = await api.listGames();
      this.games.set(games);
      const statusMap: Record<string, GameStatus> = {};
      games.forEach(g => {
        statusMap[g.name] = 'not-found';
      });
      this.statuses.set(statusMap);
      return { success: true, games };
    } catch (err) {
      this.handleAuthError(err);
      return { success: false, error: err instanceof Error ? err.message : 'Failed to load games' };
    }
  }

  async updateGameStatus(gameName: string): Promise<string> {
    try {
      const status = await this.getGameStatus(gameName);
      return status;
    } catch (err) {
      this.handleAuthError(err);
      return 'not-found';
    }
  }

  async getGameStatus(gameName: string): Promise<string> {
    try {
      const response = await api.getGameStatus(gameName);
      this.statuses.update(s => {
        s[gameName] = response.status;
        return { ...s };
      });
      return response.status;
    } catch (err) {
      this.handleAuthError(err);
      return 'not-found';
    }
  }

  async startGame(gameName: string): Promise<{ success: boolean; error?: string }> {
    try {
      await api.startGame(gameName);
      await this.getGameStatus(gameName);
      return { success: true };
    } catch (err) {
      this.handleAuthError(err);
      return { success: false, error: err instanceof Error ? err.message : 'Failed to start game' };
    }
  }

  async stopGame(gameName: string): Promise<{ success: boolean; error?: string }> {
    try {
      await api.stopGame(gameName);
      await this.getGameStatus(gameName);
      return { success: true };
    } catch (err) {
      this.handleAuthError(err);
      return { success: false, error: err instanceof Error ? err.message : 'Failed to stop game' };
    }
  }

  async restartGame(gameName: string): Promise<{ success: boolean; error?: string }> {
    try {
      await api.restartGame(gameName);
      await this.getGameStatus(gameName);
      return { success: true };
    } catch (err) {
      this.handleAuthError(err);
      return { success: false, error: err instanceof Error ? err.message : 'Failed to restart game' };
    }
  }

  getGames(): GameInfo[] {
    return get(this.games);
  }

  getStatus(gameName: string): GameStatus {
    const statuses = get(this.statuses);
    return statuses[gameName] || 'not-found';
  }

  getStatuses(): Record<string, GameStatus> {
    return get(this.statuses);
  }

  refreshStatuses(): void {
    const games = this.getGames();
    games.forEach(g => this.updateGameStatus(g.name));
  }

  getGameInfo(name: string): GameInfo | undefined {
    const games = get(this.games);
    return games.find(g => g.name === name);
  }

  getGameNameList(): string[] {
    return get(this.games).map(g => g.name);
  }

  private handleAuthError(err: unknown): void {
    if (err instanceof Error && err.message === 'Authentication required') {
      auth.logout();
      goto('/');
    }
  }

  async updateMetadata(gameName: string, appId: number, order: number): Promise<{ success: boolean; error?: string }> {
    try {
      await api.updateMetadata(gameName, appId, order);
      await this.fetchGames();
      return { success: true };
    } catch (err) {
      this.handleAuthError(err);
      return { success: false, error: err instanceof Error ? err.message : 'Failed to update metadata' };
    }
  }

  async updateArt(gameName: string): Promise<{ success: boolean; error?: string }> {
    try {
      await api.updateArt(gameName);
      await this.fetchGames();
      return { success: true };
    } catch (err) {
      this.handleAuthError(err);
      return { success: false, error: err instanceof Error ? err.message : 'Failed to update game art' };
    }
  }

  async updateGlobalSettings(mainImage: string, icon: string): Promise<{ success: boolean; error?: string }> {
    try {
      await api.updateGlobalSettings(mainImage, icon);
      return { success: true };
    } catch (err) {
      this.handleAuthError(err);
      return { success: false, error: err instanceof Error ? err.message : 'Failed to update settings' };
    }
  }

  async getServerInfo(): Promise<ServerInfo | null> {
    try {
      return await api.getServerInfo();
    } catch {
      return null;
    }
  }
}

export const gamesStore = new GamesStore();
