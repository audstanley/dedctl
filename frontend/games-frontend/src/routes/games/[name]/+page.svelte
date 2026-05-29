<script lang="ts">
  import { goto } from '$app/navigation';
  import { gamesStore } from '$lib/stores/games';
  import { onMount } from 'svelte';
  import { page } from '$app/stores';

  let { name } = $props();
  let gameStatus = $state('not-found');
  let loading = $state(true);
  let error = $state('');

  onMount(async () => {
    await gamesStore.fetchGames();
    const games = gamesStore.getGames();
    if (games.includes(name)) {
      gameStatus = await gamesStore.updateGameStatus(name);
    }
    loading = false;
  });

  async function handleStart() {
    const result = await gamesStore.startGame(name);
    if (result.success) {
      gameStatus = 'active';
    } else {
      error = result.error || 'Failed to start game';
    }
  }

  async function handleStop() {
    const result = await gamesStore.stopGame(name);
    if (result.success) {
      gameStatus = 'inactive';
    } else {
      error = result.error || 'Failed to stop game';
    }
  }

  async function handleRestart() {
    const result = await gamesStore.restartGame(name);
    if (result.success) {
      gameStatus = 'active';
    } else {
      error = result.error || 'Failed to restart game';
    }
  }

  function handleGoBack() {
    goto('/dashboard');
  }

  function getStatusColor(status: string): string {
    switch (status) {
      case 'active':
        return 'bg-green-600';
      case 'inactive':
        return 'bg-red-600';
      case 'not-found':
        return 'bg-gray-600';
      default:
        return 'bg-blue-600';
    }
  }
</script>

<div class="space-y-6">
  <div class="flex items-center space-x-4 mb-6">
    <button
      onclick={handleGoBack}
      class="bg-gray-700 hover:bg-gray-600 text-white px-4 py-2 rounded transition"
    >
      ← Back to Dashboard
    </button>
  </div>

  <div class="bg-gray-800 rounded-lg p-8 shadow-lg">
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-3xl font-bold text-white capitalize">{name}</h1>
      {#if !loading}
        <span class={`px-4 py-2 ${getStatusColor(gameStatus)} text-white text-sm rounded-full`}>
          {gameStatus}
        </span>
      {/if}
    </div>

    {#if error}
      <div class="bg-red-500 bg-opacity-20 border border-red-500 text-red-400 px-4 py-3 rounded mb-6">
        {error}
      </div>
    {/if}

    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <button
        onclick={handleStart}
        disabled={loading || gameStatus === 'active'}
        class="bg-green-600 hover:bg-green-700 disabled:bg-gray-600 disabled:cursor-not-allowed text-white px-6 py-3 rounded transition"
      >
        Start Server
      </button>
      <button
        onclick={handleStop}
        disabled={loading || gameStatus === 'inactive'}
        class="bg-red-600 hover:bg-red-700 disabled:bg-gray-600 disabled:cursor-not-allowed text-white px-6 py-3 rounded transition"
      >
        Stop Server
      </button>
      <button
        onclick={handleRestart}
        disabled={loading}
        class="bg-blue-600 hover:bg-blue-700 disabled:bg-gray-600 disabled:cursor-not-allowed text-white px-6 py-3 rounded transition"
      >
        Restart Server
      </button>
    </div>

    <div class="mt-8">
      <a
        href={`/games/${name}/logs`}
        class="inline-block bg-gray-700 hover:bg-gray-600 text-white px-6 py-3 rounded transition"
      >
        View Logs →
      </a>
    </div>
  </div>
</div>
