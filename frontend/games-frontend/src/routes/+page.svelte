<script lang="ts">
  import { goto } from '$app/navigation';
  import { login } from '$lib/stores/auth';
  import { gamesStore } from '$lib/stores/games';
  import { onMount } from 'svelte';
  import { Alert } from 'flowbite-svelte';
  import type { ServerInfo } from '$lib/api/client';

  let username = $state('');
  let password = $state('');
  let showDismissAlert = $state(false);
  let error = $state('');
  let loading = $state(false);
  let bannerImage = $state<string>('');

  onMount(async () => {
    const info = await gamesStore.getServerInfo();
    if (info?.main_image) {
      bannerImage = info.main_image;
    }
  });

  async function handleLogin() {
    showDismissAlert = false;
    error = '';
    loading = true;

    const result = await login(username, password);

    if (result.success) {
      goto('/dashboard');
    } else {
      error = result.error || 'Login failed';
      showDismissAlert = true;
    }

    loading = false;
  }
</script>

<div class="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
  {#if bannerImage}
    <div class="absolute inset-0 overflow-hidden">
      <img src="/images/{bannerImage}" alt="banner" class="w-full h-full object-cover opacity-10" />
    </div>
  {/if}
  <div class="max-w-md w-full space-y-8 relative">
    <div class="text-center">
      <div class="flex justify-center mb-4">
        <div class="w-16 h-16 bg-blue-600 rounded-full flex items-center justify-center">
          <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01" />
          </svg>
        </div>
      </div>
      <h2 class="text-3xl font-bold text-white">Game Server Control</h2>
      <p class="mt-2 text-sm text-gray-400">Sign in to manage your game servers</p>
    </div>

    {#if showDismissAlert}
      <Alert color="red" dismissable bind:alertStatus={showDismissAlert}>
        {error}
      </Alert>
    {/if}

    <form class="mt-8 space-y-6" onsubmit={(e) => { e.preventDefault(); handleLogin(); }}>
      <div class="space-y-4">
        <div>
          <label for="username" class="block text-sm font-medium text-gray-300 mb-1">Username</label>
          <input
            id="username"
            type="text"
            bind:value={username}
            placeholder="Enter username"
            class="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition"
            required
          />
        </div>
        <div>
          <label for="password" class="block text-sm font-medium text-gray-300 mb-1">Password</label>
          <input
            id="password"
            type="password"
            bind:value={password}
            placeholder="Enter password"
            class="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition"
            required
          />
        </div>
      </div>

      <button
        type="submit"
        disabled={loading || username.length < 3 || password.length < 6}
        class="w-full flex justify-center py-3 px-4 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:bg-gray-600 disabled:cursor-not-allowed transition"
      >
        {#if loading}
          <span class="flex items-center">
            <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Signing in...
          </span>
        {:else}
          Sign in
        {/if}
      </button>
    </form>
  </div>
</div>
