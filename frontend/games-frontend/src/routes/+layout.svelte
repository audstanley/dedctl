<script lang="ts">
  import favicon from '$lib/assets/favicon.svg';
  import { goto } from '$app/navigation';
  import { getUser, logout, type User } from '$lib/stores/auth';
  import { Navbar, Button } from 'flowbite-svelte';
  import { UserCircleSolid } from 'flowbite-svelte-icons';

  let { children } = $props();
  let user = $state<User | null>(getUser());

  $effect(() => {
    user = getUser();
  });

  function handleLogout() {
    logout();
    goto('/');
  }
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="min-h-screen bg-gray-900">
	<Navbar class="bg-gray-800 shadow-lg" fluid>
		<div class="flex w-full md:flex-row md:items-center md:justify-between">
			<a href="/dashboard" class="text-white text-xl font-bold hover:text-blue-400 transition">
				Game Server Control
			</a>
			<div class="flex items-center">
				{#if user}
					<div class="hidden md:block mr-4 text-gray-300 text-sm flex items-center">
						<UserCircleSolid class="h-5 w-5 mr-2" />
						<span>{user.username}</span>
						{#if user.is_admin}
							<span class="ml-2 px-2 py-1 bg-blue-600 text-white text-xs rounded">ADMIN</span>
						{/if}
					</div>
					<Button color="gray" size="xs" onclick={handleLogout}>
						Logout
					</Button>
				{/if}
			</div>
		</div>
	</Navbar>

	<main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
		{@render children()}
	</main>
</div>

