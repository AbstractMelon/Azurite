<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { gamesApi } from '$lib/api/client';
	import { toast } from '$lib/stores/notifications';
	import { ArrowLeft, UserCog, Plus, Trash2 } from 'lucide-svelte';
	import Loading from '$lib/components/Loading.svelte';
	import type { Game, User } from '$lib/types';

	let gameSlug: string;
	let game: Game | null = null;
	let moderators: User[] = [];
	let isLoading = true;
	let isAddingModerator = false;
	let newModeratorId = '';

	$: gameSlug = $page.params.slug;

	async function loadGame() {
		isLoading = true;
		try {
			const response = await gamesApi.getGame(gameSlug);
			if (response.success && response.data) {
				game = response.data as Game;
				await loadModerators();
			} else {
				toast.error('Error', 'Failed to load game');
				goto('/admin?tab=games');
			}
		} catch (error) {
			console.error('Error loading game:', error);
			toast.error('Error', 'Failed to load game');
			goto('/admin?tab=games');
		} finally {
			isLoading = false;
		}
	}

	async function loadModerators() {
		try {
			const response = await gamesApi.getGameModerators(gameSlug);
			if (response.success && response.data) {
				moderators = response.data as User[];
			}
		} catch (error) {
			console.error('Error loading moderators:', error);
		}
	}

	async function addModerator() {
		if (!newModeratorId.trim()) {
			toast.error('Error', 'Please enter a user ID');
			return;
		}

		const userId = parseInt(newModeratorId);
		if (isNaN(userId)) {
			toast.error('Error', 'Invalid user ID');
			return;
		}

		isAddingModerator = true;
		try {
			const response = await gamesApi.assignModerator(gameSlug, userId);
			if (response.success) {
				toast.success('Moderator Added', 'User has been assigned as moderator');
				newModeratorId = '';
				await loadModerators();
			} else {
				toast.error('Error', response.error || 'Failed to add moderator');
			}
		} catch (error) {
			console.error('Error adding moderator:', error);
			toast.error('Error', 'Failed to add moderator');
		} finally {
			isAddingModerator = false;
		}
	}

	async function removeModerator(userId: number, username: string) {
		if (!confirm(`Remove ${username} as moderator?`)) return;

		try {
			const response = await gamesApi.removeModerator(gameSlug, userId);
			if (response.success) {
				toast.success('Moderator Removed', `${username} has been removed as moderator`);
				await loadModerators();
			} else {
				toast.error('Error', response.error || 'Failed to remove moderator');
			}
		} catch (error) {
			console.error('Error removing moderator:', error);
			toast.error('Error', 'Failed to remove moderator');
		}
	}

	onMount(() => {
		loadGame();
	});
</script>

<svelte:head>
	<title>Manage Moderators - Admin - Azurite</title>
</svelte:head>

<div class="min-h-screen bg-background-primary">
	<div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
		{#if isLoading}
			<div class="flex items-center justify-center py-12">
				<Loading size="lg" text="Loading..." />
			</div>
		{:else if game}
			<!-- Header -->
			<div class="mb-8">
				<button
					onclick={() => goto('/admin?tab=games')}
					class="flex items-center text-text-muted hover:text-text-secondary mb-6"
				>
					<ArrowLeft class="w-4 h-4 mr-2" />
					Back to Admin
				</button>

				<div class="flex items-center space-x-3">
					<div class="w-12 h-12 bg-primary-500/20 rounded-lg flex items-center justify-center">
						<UserCog class="w-6 h-6 text-primary-400" />
					</div>
					<div>
						<h1 class="text-3xl font-bold text-text-primary">Manage Moderators</h1>
						<p class="text-text-secondary mt-1">{game.name}</p>
					</div>
				</div>
			</div>

			<!-- Add Moderator -->
			<div class="card mb-6">
				<div class="p-6">
					<h2 class="text-lg font-semibold text-text-primary mb-4">Add Moderator</h2>
					<div class="flex items-end space-x-3">
						<div class="flex-1">
							<label for="userId" class="block text-sm font-medium text-text-secondary mb-2">
								User ID
							</label>
							<input
								id="userId"
								type="number"
								bind:value={newModeratorId}
								placeholder="Enter user ID"
								class="input w-full"
							/>
						</div>
						<button onclick={addModerator} disabled={isAddingModerator} class="btn btn-primary">
							{#if isAddingModerator}
								<Loading size="sm" />
							{:else}
								<Plus class="w-4 h-4 mr-2" />
								Add
							{/if}
						</button>
					</div>
				</div>
			</div>

			<!-- Moderators List -->
			<div class="card">
				<div class="p-6">
					<h2 class="text-lg font-semibold text-text-primary mb-4">
						Current Moderators ({moderators.length})
					</h2>

					{#if moderators.length === 0}
						<div class="text-center py-8">
							<UserCog class="w-12 h-12 mx-auto text-text-muted mb-3 opacity-50" />
							<p class="text-text-secondary">No moderators assigned yet</p>
						</div>
					{:else}
						<div class="space-y-3">
							{#each moderators as moderator (moderator.id)}
								<div
									class="flex items-center justify-between p-4 bg-slate-800 rounded-lg hover:bg-slate-800/80"
								>
									<div class="flex items-center space-x-3">
										{#if moderator.avatar}
											<img
												src={moderator.avatar}
												alt={moderator.display_name}
												class="w-10 h-10 rounded-full object-cover"
											/>
										{:else}
											<div
												class="w-10 h-10 rounded-full bg-primary-600 flex items-center justify-center"
											>
												<span class="text-white font-semibold">
													{moderator.display_name.charAt(0).toUpperCase()}
												</span>
											</div>
										{/if}
										<div>
											<p class="font-medium text-text-primary">{moderator.display_name}</p>
											<p class="text-sm text-text-muted">@{moderator.username}</p>
										</div>
									</div>
									<button
										onclick={() => removeModerator(moderator.id, moderator.username)}
										class="btn btn-sm btn-outline text-red-400 hover:text-red-300"
									>
										<Trash2 class="w-4 h-4" />
									</button>
								</div>
							{/each}
						</div>
					{/if}
				</div>
			</div>
		{/if}
	</div>
</div>
