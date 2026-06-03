<script lang="ts">
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import { isAuthenticated, getUser, logout, type User } from '$lib/stores/auth';

  let { children } = $props();
  let currentUser = $state<User | null>(null);

  onMount(() => {
    if (!isAuthenticated()) {
      goto('/');
    }
    currentUser = getUser();
  });
</script>

{@render children()}
