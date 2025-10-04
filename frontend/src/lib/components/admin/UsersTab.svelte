<script lang="ts">
	import { onMount } from 'svelte';
	import { Users, Search } from 'lucide-svelte';
	import type { User } from '$lib/types';

	let users: User[] = [];
	let loading = true;
	let searchQuery = '';
	let filteredUsers: User[] = [];

	async function loadUsers() {
		loading = true;
		try {
			// TODO: Implement admin users list endpoint when backend is ready
			users = [];
			filteredUsers = [];
		} catch (error) {
			console.error('Error loading users:', error);
		} finally {
			loading = false;
		}
	}

	$: {
		if (searchQuery.trim()) {
			filteredUsers = users.filter(
				(user) =>
					user.username.toLowerCase().includes(searchQuery.toLowerCase()) ||
					user.display_name.toLowerCase().includes(searchQuery.toLowerCase()) ||
					user.email.toLowerCase().includes(searchQuery.toLowerCase())
			);
		} else {
			filteredUsers = users;
		}
	}

	onMount(() => {
		loadUsers();
	});
</script>

<div class="card">
	<div class="p-6">
		<div class="flex items-center justify-between mb-6">
			<h2 class="text-xl font-semibold text-text-primary flex items-center">
				<Users class="w-5 h-5 mr-2" />
				User Management
			</h2>
			<div class="flex items-center space-x-3">
				<div class="relative">
					<Search
						class="w-4 h-4 absolute left-3 top-1/2 transform -translate-y-1/2 text-text-muted"
					/>
					<input
						type="text"
						placeholder="Search users..."
						bind:value={searchQuery}
						class="input pl-10"
					/>
				</div>
			</div>
		</div>

		{#if loading}
			<div class="text-center py-8">Loading users...</div>
		{:else}
			<div class="text-center py-12">
				<Users class="w-16 h-16 mx-auto text-text-muted mb-4 opacity-50" />
				<h3 class="text-lg font-medium text-text-primary mb-2">
					{searchQuery ? 'No users found' : 'User Management'}
				</h3>
				<p class="text-text-secondary">
					{searchQuery
						? 'Try adjusting your search query'
						: 'User management endpoint is not yet implemented in the backend.'}
				</p>
				{#if !searchQuery}
					<p class="text-text-muted text-sm mt-2">
						This feature will be available once the backend API is complete.
					</p>
				{/if}
			</div>
		{/if}
	</div>
</div>
