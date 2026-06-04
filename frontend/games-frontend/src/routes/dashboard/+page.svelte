<script lang="ts">
  import { goto } from '$app/navigation';
  import { gamesStore } from '$lib/stores/games';
  import { onMount } from 'svelte';
  import { Alert } from 'flowbite-svelte';
  import type { GameInfo } from '$lib/api/client';

  let games = $state<GameInfo[]>([]);
  let loading = $state(true);
  let showDismissAlert = $state(false);
  let errorMessage = $state('');

  onMount(async () => {
    const { isAuthenticated } = await import('$lib/stores/auth');
    if (!isAuthenticated()) {
      loading = false;
      return;
    }
    const result = await gamesStore.fetchGames();
    if (result.success && result.games) {
      games = result.games;
    } else {
      errorMessage = result.error || 'Failed to load games';
      showDismissAlert = true;
    }
    loading = false;
  });

  function handleGameClick(game: GameInfo) {
    goto(`/games/${game.name}`);
  }
</script>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-3xl font-bold text-white">Game Servers</h1>
      <p class="text-sm text-gray-400 mt-1">Manage and monitor your game servers</p>
    </div>
    <button
      onclick={() => gamesStore.fetchGames()}
      disabled={loading}
      class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-lg shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:bg-gray-600 disabled:cursor-not-allowed transition"
    >
      {#if loading}
        <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        Refreshing...
      {:else}
        Refresh
      {/if}
    </button>
  </div>

  {#if showDismissAlert && errorMessage}
    <Alert color="red" dismissable bind:alertStatus={showDismissAlert}>
      {errorMessage}
    </Alert>
  {/if}

  {#if games.length === 0 && !loading}
    <div class="bg-gray-800 rounded-lg p-12 text-center border border-gray-700">
      <svg class="mx-auto h-12 w-12 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01" />
      </svg>
      <p class="text-gray-400 text-lg mt-4">No game servers found</p>
      <p class="text-gray-500 text-sm mt-2">Make sure Steam game services are installed and running</p>
    </div>
  {:else}
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {#each games as game (game.name)}
        <button
          onclick={() => handleGameClick(game)}
          class="text-left bg-gray-800 hover:bg-gray-750 border border-gray-700 rounded-lg overflow-hidden transition hover:border-blue-500 hover:shadow-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          {#if game.has_image}
            <div class="h-32 bg-gray-900 overflow-hidden">
              <img
                src="/images/{game.app_id}.jpg"
                alt={game.name}
                class="w-full h-full object-cover"
                loading="lazy"
              />
            </div>
          {/if}
          <div class="p-6">
            <div class="flex items-center justify-between mb-4">
              <h2 class="text-xl font-bold text-white capitalize">{#if game.has_image}{game.name}{:else}{game.name}{/if}</h2>
              <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-900 text-green-300">
                Ready
              </span>
            </div>
            <p class="text-gray-400 text-sm">Click to manage server</p>
            <p class="text-[10px] text-gray-600 mt-2 lowercase">{game.name}</p>
          </div>
        </button>
      {/each}
    </div>
  {/if}
</div>
