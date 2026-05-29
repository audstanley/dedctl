<script lang="ts">
  import { authStore } from '$lib/stores/auth';
  import { goto } from '$app/navigation';

  let username = $state('');
  let password = $state('');
  let error = $state('');
  let loading = $state(false);

  async function handleLogin() {
    error = '';
    loading = true;

    const result = await authStore.login(username, password);

    if (result.success) {
      goto('/dashboard');
    } else {
      error = result.error || 'Login failed';
    }

    loading = false;
  }

  function validateUsername(): boolean {
    return username.trim().length >= 3;
  }

  function validatePassword(): boolean {
    return password.length >= 6;
  }
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-900">
  <div class="max-w-md w-full space-y-8 p-10 bg-gray-800 rounded-lg shadow-2xl">
    <div>
      <h2 class="mt-6 text-center text-3xl font-extrabold text-white">
        Game Server Control
      </h2>
      <p class="mt-2 text-center text-sm text-gray-400">
        Sign in to manage your game servers
      </p>
    </div>

    {#if error}
      <div class="bg-red-500 bg-opacity-20 border border-red-500 text-red-400 px-4 py-3 rounded">
        {error}
      </div>
    {/if}

    <form class="mt-8 space-y-6" onsubmit={(e) => { e.preventDefault(); handleLogin(); }}>
      <div class="rounded-md shadow-sm space-y-4">
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
            class="appearance-none rounded relative block w-full px-3 py-2 border border-gray-600 placeholder-gray-500 text-gray-100 bg-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
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
            class="appearance-none rounded relative block w-full px-3 py-2 border border-gray-600 placeholder-gray-500 text-gray-100 bg-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            placeholder="Enter password"
          />
        </div>
      </div>

      <div>
        <button
          type="submit"
          disabled={loading || !validateUsername() || !validatePassword()}
          class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
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
      </div>

      <div class="text-center">
        <a href="/register" class="font-medium text-blue-400 hover:text-blue-300">
          Don't have an account? Register
        </a>
      </div>
    </form>
  </div>
</div>

