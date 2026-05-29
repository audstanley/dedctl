<script lang="ts">
  import { authStore } from '$lib/stores/auth';
  import { goto } from '$app/navigation';
  import { Alert, Button } from 'flowbite-svelte';

  let username = $state('');
  let password = $state('');
  let showDismissAlert = $state(true);
  let loading = $state(false);

  function validateUsername(): boolean {
    return username.trim().length >= 3;
  }

  function validatePassword(): boolean {
    return password.length >= 6;
  }

  async function handleLogin() {
    showDismissAlert = true;
    loading = true;

    const result = await authStore.login(username, password);

    if (result.success) {
      goto('/dashboard');
    } else {
      showDismissAlert = true;
    }

    loading = false;
  }
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-900">
  <div class="w-full max-w-md p-8 space-y-8 bg-gray-800 rounded-xl shadow-2xl">
    <div class="text-center">
      <h1 class="text-4xl font-extrabold text-white tracking-tight">
        Game Server Control
      </h1>
      <p class="mt-3 text-base text-gray-400">
        Sign in to manage your game servers
      </p>
    </div>

    {#if showDismissAlert}
      <Alert color="red" dismissable bind:alertStatus={showDismissAlert}>
        Login failed - please try again
      </Alert>
    {/if}

    <form class="mt-8 space-y-6" onsubmit={(e) => { e.preventDefault(); handleLogin(); }}>
      <div class="space-y-4">
        <div>
          <label for="username" class="block text-sm font-medium text-gray-300 mb-1">
            Username
          </label>
          <input
            id="username"
            name="username"
            type="text"
            bind:value={username}
            required
            class="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-gray-100 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition"
            placeholder="Enter username"
          />
        </div>
        <div>
          <label for="password" class="block text-sm font-medium text-gray-300 mb-1">
            Password
          </label>
          <input
            id="password"
            name="password"
            type="password"
            bind:value={password}
            required
            class="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-gray-100 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition"
            placeholder="Enter password"
          />
        </div>
      </div>

      <div>
        <Button 
          type="submit"
          disabled={loading || !validateUsername() || !validatePassword()}
          color="blue"
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
        </Button>
      </div>

      <div class="text-center">
        <a href="/register" class="font-medium text-blue-400 hover:text-blue-300 transition">
          Don't have an account? Register
        </a>
      </div>
    </form>
  </div>
</div>

