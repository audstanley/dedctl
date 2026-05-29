<script lang="ts">
  import { goto } from '$app/navigation';
  import { authManager } from '$lib/stores/auth';
  import { onMount } from 'svelte';
  import type { User } from '$lib/stores/auth';

  let { children } = $props();
  let currentUser = $state<User | null>(null);

  onMount(() => {
    if (!authManager.isAuthenticated()) {
      goto('/');
    }
    currentUser = authManager.getUser();
  });
</script>

<div class="min-h-screen bg-gray-900">
	<nav class="bg-gray-800 shadow-lg">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="flex items-center justify-between h-16">
				<div class="flex items-center">
					<a href="/dashboard" class="text-white text-xl font-bold hover:text-blue-400 transition">
						Game Server Control
					</a>
				</div>
				<div class="flex items-center space-x-4">
					{#if currentUser}
						<span class="text-gray-300 text-sm">
							{currentUser.username}
							{#if currentUser.is_admin}
								<span class="ml-2 px-2 py-1 bg-blue-600 text-white text-xs rounded">ADMIN</span>
							{/if}
						</span>
						<button
							onclick={authManager.logout}
							class="bg-gray-700 hover:bg-gray-600 text-white px-4 py-2 rounded transition"
						>
							Logout
						</button>
					{/if}
				</div>
			</div>
		</div>
	</nav>

	<main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
		{@render children()}
	</main>
</div>
