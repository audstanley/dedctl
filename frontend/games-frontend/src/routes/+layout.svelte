<script lang="ts">
  import '../app.css';
  import { goto } from '$app/navigation';
  import { auth, logout, type User } from '$lib/stores/auth';
  import { Navbar, NavContainer, Button } from 'flowbite-svelte';

  let { children } = $props();
  let currentUser = $state<User | null>(auth.getUser());

  $effect(() => {
    return auth.subscribe(() => {
      currentUser = auth.getUser();
    });
  });

  function handleLogout() {
    logout();
    goto('/');
  }
</script>

<div class="min-h-screen bg-gray-900">
	<Navbar class="bg-gray-800 shadow-lg">
		<NavContainer>
			<a href="/dashboard" class="text-white text-xl font-bold hover:text-blue-400 transition">
				Game Server Control
			</a>
		</NavContainer>
		<div class="flex items-center">
			{#if currentUser}
				<div class="hidden md:block mr-4 text-gray-300 text-sm flex items-center">
					<span>{currentUser.username}</span>
					{#if currentUser.is_admin}
						<span class="ml-2 px-2 py-1 bg-blue-600 text-white text-xs rounded">ADMIN</span>
					{/if}
				</div>
				<Button color="gray" size="xs" onclick={handleLogout}>
					Logout
				</Button>
			{/if}
		</div>
	</Navbar>

	<main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
		{@render children()}
	</main>
</div>
