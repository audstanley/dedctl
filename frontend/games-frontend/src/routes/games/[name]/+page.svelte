<script lang="ts">
  import { goto } from '$app/navigation';
  import { gamesStore } from '$lib/stores/games';
  import { onMount } from 'svelte';
  import { Alert, Badge, Button } from 'flowbite-svelte';

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

  async function refreshStatus() {
    gameStatus = await gamesStore.getGameStatus(name);
  }
</script>

<div class="space-y-6">
  <div class="flex items-center space-x-4 mb-6">
    <Button color="gray" onclick={handleGoBack}>Back to Dashboard</Button>
  </div>

  <div class="bg-gray-800 rounded-lg p-8 shadow-lg">
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-3xl font-bold text-white capitalize">{name}</h1>
      {#if !loading}
        <Badge color={gameStatus === 'active' ? 'green' : gameStatus === 'inactive' ? 'red' : 'gray'}>{gameStatus}</Badge>
      {/if}
    </div>

    {#if error}
      <Alert color="red" dismissable>
        {error}
      </Alert>
    {/if}

    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <Button color="green" onclick={handleStart} disabled={loading || gameStatus === 'active'}>Start Server</Button>
      <Button color="red" onclick={handleStop} disabled={loading || gameStatus === 'inactive'}>Stop Server</Button>
      <Button color="blue" onclick={handleRestart} disabled={loading}>Restart Server</Button>
    </div>

    <div class="mt-8">
      <Button color="gray" onclick={() => goto(`/games/${name}/logs`)}>View Logs →</Button>
    </div>
  </div>
</div>
