<script lang="ts">
  import { goto } from '$app/navigation';
  import { gamesStore } from '$lib/stores/games';
  import { onMount } from 'svelte';
  import { Alert, Card, Badge } from 'flowbite-svelte';

  let games = $state<string[]>([]);
  let loading = $state(true);
  let showDismissAlert = $state(false);
  let errorMessage = $state('');

  onMount(async () => {
    const result = await gamesStore.fetchGames();
    if (result.success) {
      games = result.games;
    } else {
      errorMessage = result.error || 'Failed to load games';
      showDismissAlert = true;
    }
    loading = false;
  });

  function handleGameClick(gameName: string) {
    goto(`/games/${gameName}`);
  }

  async function refreshGames() {
    loading = true;
    showDismissAlert = false;
    errorMessage = '';
    const result = await gamesStore.fetchGames();
    if (result.success) {
      games = result.games;
    } else {
      errorMessage = result.error || 'Failed to refresh games';
      showDismissAlert = true;
    }
    loading = false;
  }
</script>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <h1 class="text-3xl font-bold text-white">Game Servers</h1>
    <button
      onclick={refreshGames}
      disabled={loading}
      class="bg-blue-600 hover:bg-blue-700 disabled:bg-gray-600 text-white px-4 py-2 rounded transition"
    >
      {#if loading}
        <span class="flex items-center">
          <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          Refreshing...
        </span>
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
    <div class="bg-gray-800 rounded-lg p-8 text-center">
      <p class="text-gray-400 text-lg">No game servers found</p>
      <p class="text-gray-500 text-sm mt-2">Make sure Steam game services are installed and running</p>
    </div>
  {:else}
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {#each games as game (game)}
        <Card class="p-6 bg-gray-800 hover:bg-gray-750 cursor-pointer transition" onclick={() => handleGameClick(game)}>
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-xl font-bold text-white capitalize">{game}</h2>
            <Badge color="green">Ready</Badge>
          </div>
          <p class="text-gray-400 text-sm">Click to manage server</p>
        </Card>
      {/each}
    </div>
  {/if}
</div>
