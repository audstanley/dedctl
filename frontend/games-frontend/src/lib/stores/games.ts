import { writable, type Writable } from 'svelte/store';
import { api } from '$lib/api/client';

type GameStatus = 'active' | 'inactive' | 'not-found' | string;

type GameInfo = {
  name: string;
  status: GameStatus;
};

class GamesStore {
  private games: Writable<string[]> = writable([]);
  private statuses: Writable<Record<string, GameStatus>> = writable({});

  async fetchGames() {
    try {
      const games = await api.listGames();
      this.games.set(games);
      
      const statuses: Record<string, GameStatus> = {};
      games.forEach(game => {
        statuses[game] = 'not-found';
      });
      this.statuses.set(statuses);
      
      return { success: true, games };
    } catch (error: any) {
      return { success: false, error: error.message };
    }
  }

  async updateGameStatus(gameName: string): Promise<GameStatus> {
    try {
      const status = await this.getGameStatus(gameName);
      const current = get(this.statuses);
      current[gameName] = status;
      this.statuses.set(current);
      return status;
    } catch (error) {
      return 'error';
    }
  }

  async getGameStatus(gameName: string): Promise<GameStatus> {
    const games = get(this.games);
    if (!games.includes(gameName)) {
      return 'not-found';
    }
    try {
      const res = await api.getGameStatus(gameName);
      const status = res.status;
      return status === 'active' ? 'active' : 'inactive';
    } catch {
      return 'inactive';
    }
  }

  async startGame(gameName: string): Promise<{ success: boolean; error?: string }> {
    try {
      await api.startGame(gameName);
      await this.updateGameStatus(gameName);
      return { success: true };
    } catch (error: any) {
      return { success: false, error: error.message };
    }
  }

  async stopGame(gameName: string): Promise<{ success: boolean; error?: string }> {
    try {
      await api.stopGame(gameName);
      await this.updateGameStatus(gameName);
      return { success: true };
    } catch (error: any) {
      return { success: false, error: error.message };
    }
  }

  async restartGame(gameName: string): Promise<{ success: boolean; error?: string }> {
    try {
      await api.restartGame(gameName);
      await this.updateGameStatus(gameName);
      return { success: true };
    } catch (error: any) {
      return { success: false, error: error.message };
    }
  }

  getGames() {
    return get(this.games);
  }

  getStatus(gameName: string): GameStatus {
    const statuses = get(this.statuses);
    return statuses[gameName] || 'not-found';
  }

  getStatuses() {
    return get(this.statuses);
  }

  refreshStatuses() {
    const games = this.getGames();
    games.forEach(game => this.updateGameStatus(game));
  }
}

export const gamesStore = new GamesStore();
