<script lang="ts">
  import { goto } from '$app/navigation';
  import { gamesStore } from '$lib/stores/games';
  import { auth, getUser, type User } from '$lib/stores/auth';
  import { onMount } from 'svelte';
  import { Alert, Badge } from 'flowbite-svelte';
  import type { GameInfo } from '$lib/api/client';

  let games = $state<GameInfo[]>([]);
  let loading = $state(true);
  let currentUser = $state<User | null>(null);
  let errorMessage = $state('');
  let successMessage = $state('');
  let showDismissAlert = $state(false);
  let showDismissSuccess = $state(false);
  let updatingArt = $state<string | null>(null);
  let editingAppId = $state<Record<string, number>>({});
  let editingOrder = $state<Record<string, number>>({});

  onMount(async () => {
    currentUser = getUser();
    if (!currentUser?.is_admin) {
      goto('/dashboard');
      return;
    }
    const result = await gamesStore.fetchGames();
    if (result.success && result.games) {
      games = result.games;
      const initialAppIds: Record<string, number> = {};
      const initialOrders: Record<string, number> = {};
      games.forEach(g => {
        initialAppIds[g.name] = g.app_id;
        initialOrders[g.name] = g.order;
      });
      editingAppId = initialAppIds;
      editingOrder = initialOrders;
    } else {
      errorMessage = result.error || 'Failed to load games';
      showDismissAlert = true;
    }
    loading = false;
  });

  async function handleSaveMetadata(gameName: string) {
    const appId = editingAppId[gameName] || 0;
    const order = editingOrder[gameName] || 0;
    const result = await gamesStore.updateMetadata(gameName, appId, order);
    if (result.success) {
      games = games.map(g => {
        if (g.name === gameName) {
          return { ...g, app_id: appId, order };
        }
        return g;
      });
      successMessage = `Metadata for ${gameName} saved`;
      showDismissSuccess = true;
      errorMessage = '';
    } else {
      errorMessage = result.error || 'Failed to save metadata';
      showDismissAlert = true;
    }
  }

  async function handleUpdateArt(gameName: string, appId: number) {
    if (appId <= 0) {
      errorMessage = `No AppID set for ${gameName}. Set an AppID and save first.`;
      showDismissAlert = true;
      return;
    }
    updatingArt = gameName;
    errorMessage = '';
    const result = await gamesStore.updateArt(gameName);
    updatingArt = null;
    if (result.success) {
      games = games.map(g => {
        if (g.name === gameName) {
          return { ...g, has_image: true };
        }
        return g;
      });
      successMessage = `Game art for ${gameName} updated`;
      showDismissSuccess = true;
      errorMessage = '';
    } else {
      errorMessage = result.error || 'Failed to update game art';
      showDismissAlert = true;
    }
  }

  function getStatusBadge(game: GameInfo): 'gray' | 'green' | 'amber' {
    if (game.app_id <= 0) return 'gray';
    if (game.has_image) return 'green';
    return 'amber';
  }

  function getStatusText(game: GameInfo): string {
    if (game.app_id <= 0) return 'No AppID';
    if (game.has_image) return 'Image cached';
    return 'No image';
  }
</script>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-3xl font-bold text-white">Admin Settings</h1>
      <p class="text-sm text-gray-400 mt-1">Manage game metadata and download cover art</p>
    </div>
    <button
      onclick={() => goto('/dashboard')}
      class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-lg shadow-sm text-white bg-gray-600 hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500 transition"
    >
      Back to Dashboard
    </button>
  </div>

  {#if showDismissSuccess && successMessage}
    <Alert color="green" dismissable bind:alertStatus={showDismissSuccess}>
      {successMessage}
    </Alert>
  {/if}

  {#if showDismissAlert && errorMessage}
    <Alert color="red" dismissable bind:alertStatus={showDismissAlert}>
      {errorMessage}
    </Alert>
  {/if}

  {#if loading}
    <div class="flex items-center justify-center py-12">
      <svg class="animate-spin h-8 w-8 text-blue-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
      <span class="ml-3 text-gray-400">Loading games...</span>
    </div>
  {:else if games.length === 0}
    <div class="bg-gray-800 rounded-lg p-12 text-center border border-gray-700">
      <p class="text-gray-400 text-lg">No game servers found</p>
      <p class="text-gray-500 text-sm mt-2">Make sure Steam game services are installed and running</p>
    </div>
  {:else}
    <div class="bg-gray-800 rounded-lg shadow-lg overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-700">
          <thead class="bg-gray-900">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">Game</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">AppID</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">Order</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">Status</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody class="bg-gray-800 divide-y divide-gray-700">
            {#each games as game (game.name)}
              <tr class="hover:bg-gray-750 transition">
                <td class="px-6 py-4 whitespace-nowrap">
                  <span class="text-white font-medium capitalize">{game.name}</span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <input
                    type="number"
                    min="0"
                    value={editingAppId[game.name] ?? ''}
                    oninput={(e) => {
                      const target = e.target as HTMLInputElement;
                      editingAppId = { ...editingAppId, [game.name]: parseInt(target.value) || 0 };
                    }}
                    class="w-28 bg-gray-900 border border-gray-600 rounded px-3 py-1.5 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
                    placeholder="Enter AppID"
                  />
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <input
                    type="number"
                    min="0"
                    value={editingOrder[game.name] ?? ''}
                    oninput={(e) => {
                      const target = e.target as HTMLInputElement;
                      editingOrder = { ...editingOrder, [game.name]: parseInt(target.value) || 0 };
                    }}
                    class="w-20 bg-gray-900 border border-gray-600 rounded px-3 py-1.5 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
                    placeholder="Order"
                  />
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <Badge color={getStatusBadge(game)}>{getStatusText(game)}</Badge>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center space-x-2">
                    <button
                      onclick={() => handleSaveMetadata(game.name)}
                      class="inline-flex items-center px-3 py-1.5 bg-blue-700 hover:bg-blue-800 text-white text-xs font-medium rounded-lg transition"
                    >
                      Save
                    </button>
                    <button
                      onclick={() => handleUpdateArt(game.name, editingAppId[game.name] || 0)}
                      disabled={updatingArt !== null || editingAppId[game.name] <= 0}
                      class="inline-flex items-center px-3 py-1.5 bg-green-700 hover:bg-green-800 disabled:bg-gray-600 disabled:cursor-not-allowed text-white text-xs font-medium rounded-lg transition"
                    >
                      {#if updatingArt === game.name}
                        <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                        </svg>
                        Downloading...
                      {:else}
                        Update Art
                      {/if}
                    </button>
                  </div>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>
  {/if}
</div>
