<script lang="ts">
  import { authStore } from '$lib/stores/auth';
  import { goto } from '$app/navigation';
  import { Alert, Button, Input, Label } from 'flowbite-svelte';

  let username = $state('');
  let password = $state('');
  let showDismissAlert = $state(false);
  let error = $state('');
  let loading = $state(false);

  async function handleLogin() {
    showDismissAlert = false;
    error = '';
    loading = true;

    const result = await authStore.login(username, password);

    if (result.success) {
      goto('/dashboard');
    } else {
      error = result.error || 'Login failed';
      showDismissAlert = true;
    }

    loading = false;
  }
</script>

<div class="min-h-screen bg-gray-900 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
  <div class="max-w-md w-full space-y-8">
    <div>
      <h2 class="mt-6 text-center text-3xl font-bold text-white">
        Game Server Control
      </h2>
      <p class="mt-2 text-center text-sm text-gray-400">
        Sign in to manage your servers
      </p>
    </div>

    {#if showDismissAlert}
      <Alert color="red" dismissable bind:alertStatus={showDismissAlert}>
        {error}
      </Alert>
    {/if}

    <form class="mt-8 space-y-6" onsubmit={(e) => { e.preventDefault(); handleLogin(); }}>
      <div class="space-y-4">
        <div>
          <Label for="username">Username</Label>
          <Input 
            id="username" 
            type="text" 
            bind:value={username}
            placeholder="Enter username"
            required
          />
        </div>
        <div>
          <Label for="password">Password</Label>
          <Input 
            id="password" 
            type="password" 
            bind:value={password}
            placeholder="Enter password"
            required
          />
        </div>
      </div>

      <Button 
        type="submit"
        disabled={loading || username.length < 3 || password.length < 6}
        color="blue"
      >
        {#if loading}
          <span class="flex items-center">
            <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Signing in...
          </span>
        {:else}
          Sign in
        {/if}
      </Button>
    </form>
  </div>
</div>
