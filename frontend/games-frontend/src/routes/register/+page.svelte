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
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-900 py-12 px-4 sm:px-6 lg:px-8">
  <div class="max-w-md w-full space-y-8">
    <div>
      <h2 class="mt-6 text-center text-3xl font-bold text-white">
        Create Account
      </h2>
      <p class="mt-2 text-center text-sm text-gray-400">
        Register to manage your servers
      </p>
    </div>

    {#if error}
      <div class="bg-red-500 bg-opacity-10 border border-red-500 text-red-400 px-4 py-3 rounded-lg text-sm">
        {error}
      </div>
    {/if}

    <form class="mt-8 space-y-6" onsubmit={(e) => { e.preventDefault(); handleRegister(); }}>
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
            class="w-full px-4 py-2.5 bg-gray-800 border border-gray-600 rounded-lg text-gray-100 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
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
            class="w-full px-4 py-2.5 bg-gray-800 border border-gray-600 rounded-lg text-gray-100 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
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
            class="w-full px-4 py-2.5 bg-gray-800 border border-gray-600 rounded-lg text-gray-100 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
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
          Admin account
        </label>
      </div>

      <div>
        <button
          type="submit"
          disabled={loading || username.length < 3 || password.length < 6 || password !== confirmPassword}
          class="w-full flex justify-center py-2.5 px-4 border border-transparent text-sm font-medium rounded-lg text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          {#if loading}
            Creating account...
          {:else}
            Register
          {/if}
        </button>
      </div>

      <div class="text-center">
        <a href="/" class="font-medium text-blue-400 hover:text-blue-300 text-sm">
          Already have an account? Sign in
        </a>
      </div>
    </form>
  </div>
</div>
