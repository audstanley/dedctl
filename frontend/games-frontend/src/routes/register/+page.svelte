<script lang="ts">
  import { authStore } from '$lib/stores/auth';
  import { goto } from '$app/navigation';

  let username = $state('');
  let password = $state('');
  let confirmPassword = $state('');
  let isAdmin = $state(false);
  let error = $state('');
  let loading = $state(false);

  async function handleRegister() {
    error = '';
    loading = true;

    if (password !== confirmPassword) {
      error = 'Passwords do not match';
      loading = false;
      return;
    }

    if (password.length < 6) {
      error = 'Password must be at least 6 characters';
      loading = false;
      return;
    }

    const result = await authStore.register(username, password, isAdmin);

    if (result.success) {
      goto('/');
    } else {
      error = result.error || 'Registration failed';
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
        Create Account
      </h2>
      <p class="mt-2 text-center text-sm text-gray-400">
        Register to manage game servers
      </p>
    </div>

    {#if error}
      <div class="bg-red-500 bg-opacity-20 border border-red-500 text-red-400 px-4 py-3 rounded">
        {error}
      </div>
    {/if}

    <form class="mt-8 space-y-6" on:submit|preventDefault={handleRegister}>
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
            placeholder="Choose username (min 3 chars)"
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
            placeholder="Min 6 characters"
          />
        </div>
        <div>
          <label for="confirmPassword" class="block text-sm font-medium text-gray-300 mb-1">
            Confirm Password
          </label>
          <input
            id="confirmPassword"
            name="confirmPassword"
            type="password"
            bind:value={confirmPassword}
            required
            class="appearance-none rounded relative block w-full px-3 py-2 border border-gray-600 placeholder-gray-500 text-gray-100 bg-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            placeholder="Re-enter password"
          />
        </div>
      </div>

      <div class="flex items-center">
        <input
          id="isAdmin"
          name="isAdmin"
          type="checkbox"
          bind:checked={isAdmin}
          class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-600 rounded bg-gray-700"
        />
        <label for="isAdmin" class="ml-2 block text-sm text-gray-300">
          Admin account (can create other admins)
        </label>
      </div>

      <div>
        <button
          type="submit"
          disabled={loading || !validateUsername() || !validatePassword() || password !== confirmPassword}
          class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {#if loading}
            <span class="flex items-center">
              <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Creating account...
            </span>
          {:else}
            Register
          {/if}
        </button>
      </div>

      <div class="text-center">
        <a href="/" class="font-medium text-blue-400 hover:text-blue-300">
          Already have an account? Sign in
        </a>
      </div>
    </form>
  </div>
</div>
