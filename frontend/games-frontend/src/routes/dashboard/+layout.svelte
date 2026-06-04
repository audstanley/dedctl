<script lang="ts">
  import { goto } from '$app/navigation';
  import { isAuthenticated } from '$lib/stores/auth';

  let { children } = $props();

  $effect(() => {
    if (!isAuthenticated()) {
      goto('/');
    }
  });
</script>

<main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
	{@render children()}
</main>
