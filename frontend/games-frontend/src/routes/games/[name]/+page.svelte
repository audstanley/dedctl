<script lang="ts">
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import { gamesStore } from '$lib/stores/games';
  import { auth } from '$lib/stores/auth';
  import { onMount } from 'svelte';
  import { Alert, Badge, Button } from 'flowbite-svelte';
  import type { GameInfo } from '$lib/api/client';

  $effect(() => {
    $page;
  });
  const name = $derived($page.params.name as string);
  let gameInfo = $state<GameInfo | null>(null);
  let gameStatus = $state('not-found');
  let loading = $state(true);
  let error = $state('');
  let adminUser = $derived(auth.getUser()?.is_admin === true);

  onMount(async () => {
    const game = gamesStore.getGameInfo(name);
    gameInfo = game ?? null;

    if (!game) {
      await gamesStore.fetchGames();
      const updated = gamesStore.getGameInfo(name);
      gameInfo = updated ?? null;
    }
    try {
      gameStatus = await gamesStore.updateGameStatus(name);
    } catch {
      gameStatus = 'not-found';
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

  async function handleEnable() {
    const result = await gamesStore.enableGame(name);
    if (result.success) {
      gameInfo = gameInfo ? { ...gameInfo, enabled: true } : null;
    } else {
      error = result.error || 'Failed to enable game';
    }
  }

  async function handleDisable() {
    const result = await gamesStore.disableGame(name);
    if (result.success) {
      gameInfo = gameInfo ? { ...gameInfo, enabled: false } : null;
    } else {
      error = result.error || 'Failed to disable game';
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

  <div class="bg-gray-800 rounded-lg shadow-lg overflow-hidden">
    {#if gameInfo?.has_image}
      <div class="h-48 bg-gray-900 overflow-hidden">
        <img
          src="/images/{gameInfo.app_id}.jpg"
          alt={name}
          class="w-full h-full object-cover"
        />
      </div>
    {/if}

    <div class="p-8">
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
        <Button color="red" onclick={handleStop} disabled={gameStatus === 'inactive' || gameStatus === 'not-found'}>Stop Server</Button>
        <Button color="blue" onclick={handleRestart}>Restart Server</Button>
      </div>

      {#if adminUser}
        <div class="mt-4">
          {#if gameInfo?.enabled}
            <Button color="red" onclick={handleDisable} disabled={loading}>Disable (Auto-start on login)</Button>
          {:else}
            <Button color="gray" onclick={handleEnable} disabled={loading}>Enable (Auto-start on login)</Button>
          {/if}
        </div>
      {/if}

      <div class="mt-8">
        <Button color="gray" onclick={() => goto(`/games/${name}/logs`)}>View Logs →</Button>
      </div>
    </div>
  </div>
</div>
