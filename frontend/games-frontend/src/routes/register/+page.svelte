<script lang="ts">
  import { authStore } from '$lib/stores/auth';
  import { goto } from '$app/navigation';
  import { Alert, Button, Input, Label, Checkbox, Helper } from 'flowbite-svelte';

  let username = $state('');
  let password = $state('');
  let confirmPassword = $state('');
  let isAdmin = $state(false);
  let showDismissAlert = $state(false);
  let registerError = $state('');

  async function handleRegister() {
    showDismissAlert = false;
    registerError = '';

    if (username.length < 3) {
      registerError = 'Username must be at least 3 characters';
      showDismissAlert = true;
      return;
    }

    if (password.length < 6) {
      registerError = 'Password must be at least 6 characters';
      showDismissAlert = true;
      return;
    }

    if (password !== confirmPassword) {
      registerError = 'Passwords do not match';
      showDismissAlert = true;
      return;
    }

    const result = await authStore.register(username, password, isAdmin);

    if (result.success) {
      goto('/');
    } else {
      registerError = result.error || 'Registration failed';
      showDismissAlert = true;
    }
  }
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-900">
  <div class="w-full max-w-md p-8 space-y-8 bg-gray-800 rounded-xl shadow-2xl">
    <div class="text-center">
      <h2 class="text-4xl font-extrabold text-white">
        Create Account
      </h2>
      <p class="mt-3 text-base text-gray-400">
        Register to manage game servers
      </p>
    </div>

    {#if showDismissAlert && registerError}
      <Alert color="red" dismissable bind:alertStatus={showDismissAlert}>
        {registerError}
      </Alert>
    {/if}

    <form class="mt-8 space-y-6" onsubmit={(e) => { e.preventDefault(); handleRegister(); }}>
      <div class="space-y-4">
        <div>
          <Label for="username" class="mb-2">Username</Label>
          <Input 
            id="username" 
            type="text" 
            bind:value={username}
            placeholder="Choose username (min 3 chars)"
            required
          />
        </div>
        <div>
          <Label for="password" class="mb-2">Password</Label>
          <Input 
            id="password" 
            type="password" 
            bind:value={password}
            placeholder="Min 6 characters"
            required
          />
        </div>
        <div>
          <Label for="confirmPassword" class="mb-2">Confirm Password</Label>
          <Input 
            id="confirmPassword" 
            type="password" 
            bind:value={confirmPassword}
            placeholder="Re-enter password"
            required
          />
        </div>
      </div>

      <div class="flex items-start">
        <Checkbox id="isAdmin" bind:checked={isAdmin} />
        <label for="isAdmin" class="ml-2 block text-sm text-gray-300">
          Admin account (can create other admins)
        </label>
      </div>

      <Button type="submit" color="blue">
        Register
      </Button>

      <div class="text-center">
        <a href="/" class="font-medium text-blue-400 hover:text-blue-300 transition">
          Already have an account? Sign in
        </a>
      </div>
    </form>
  </div>
</div>
