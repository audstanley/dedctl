<script lang="ts">
  import { goto } from '$app/navigation';
  import { auth, isAuthenticated, getUser, type User } from '$lib/stores/auth';

  let { children } = $props();
  let currentUser = $state<User | null>(null);

  $effect(() => {
    return auth.subscribe(() => {
      currentUser = getUser();
      if (!isAuthenticated() || !currentUser?.is_admin) {
        goto(currentUser?.is_admin ? '/dashboard' : '/');
      }
    });
  });
</script>

{@render children()}
