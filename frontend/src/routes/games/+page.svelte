<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { gamesApi } from '$lib/api/client';
	import { toast } from '$lib/stores/notifications';
	import { isAuthenticated } from '$lib/stores/auth';
	import Loading from '$lib/components/Loading.svelte';
	import { GameCard } from '$lib/components/cards';
	import type { Game } from '$lib/types';
	import { Search, Gamepad2, Plus } from 'lucide-svelte';

	let games: Game[] = [];
	let isLoading = true;
	let searchQuery = '';
	let filteredGames: Game[] = [];
	let currentPage = 1;
	let totalPages = 1;
	let totalCount = 0;
	let perPage = 20;

	// Load games data
	onMount(async () => {
		await loadGames();
	});

	// Load games from API
	async function loadGames() {
		isLoading = true;
		try {
			const response = await gamesApi.getGames({
				page: currentPage,
				per_page: perPage,
				search: searchQuery || undefined
			});

			if (response.success && response.data) {
				games = response.data.data || [];
				totalPages = response.data.total_pages || 1;
				totalCount = response.data.total || 0;
				currentPage = response.data.page || 1;
				filterGames();
			} else {
				toast.error('Failed to load games', response.error);
			}
		} catch (error) {
			console.error('Error loading games:', error);
			toast.error('Error', 'Failed to load games');
		} finally {
			isLoading = false;
		}
	}

	// Filter games based on search
	function filterGames() {
		if (!searchQuery.trim()) {
			filteredGames = games;
		} else {
			const query = searchQuery.toLowerCase().trim();
			filteredGames = games.filter(
				(game) =>
					game.name.toLowerCase().includes(query) ||
					(game.description && game.description.toLowerCase().includes(query))
			);
		}
	}

	// Handle search
	function handleSearch(event: Event) {
		// Prevent default form submission
		event.preventDefault();

		currentPage = 1;
		loadGames();
	}

	// Handle page change
	function changePage(newPage: number) {
		if (newPage >= 1 && newPage <= totalPages && newPage !== currentPage) {
			currentPage = newPage;
			loadGames();
		}
	}

	// Handle game request
	function handleGameRequest() {
		if ($isAuthenticated) {
			goto('/games/request');
		} else {
			goto('/auth/login?redirect=/games/request');
		}
	}

	// Reactive search
	$: if (searchQuery !== undefined) {
		filterGames();
	}
</script>

<svelte:head>
	<title>Games - Azurite</title>
	<meta
		name="description"
		content="Browse all games supported on Azurite. Discover mods for your favorite games and join thriving communities."
	/>
</svelte:head>

<div class="min-h-screen bg-background-primary">
	<!-- Header Section -->
	<div class="bg-gradient-to-r from-primary-600/20 to-primary-700/20 border-b border-slate-700">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
			<div class="text-center">
				<h1 class="text-4xl font-bold text-text-primary mb-4">
					<Gamepad2 class="w-10 h-10 inline-block mr-2 mb-1" />
					Browse Games
				</h1>
				<p class="text-text-secondary text-lg max-w-2xl mx-auto mb-8">
					Discover mods for your favorite games. Each game has its own dedicated community with
					thousands of modifications to enhance your gaming experience.
				</p>

				<!-- Search Bar -->
				<div class="max-w-md mx-auto">
					<form onsubmit={handleSearch}>
						<div class="relative">
							<Search
								class="absolute left-3 top-1/2 transform -translate-y-1/2 text-text-muted w-5 h-5"
							/>
							<input
								type="search"
								placeholder="Search games..."
								bind:value={searchQuery}
								oninput={filterGames}
								class="input pl-10 pr-4 py-3 w-full text-center"
							/>
						</div>
					</form>
				</div>
			</div>
		</div>
	</div>

	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
		<!-- Stats and Actions -->
		<div class="flex flex-col sm:flex-row justify-between items-center mb-8">
			<div class="mb-4 sm:mb-0">
				{#if !isLoading}
					<p class="text-text-secondary">
						{searchQuery ? `Found ${filteredGames.length} games` : `${totalCount} games available`}
					</p>
				{/if}
			</div>

			<button onclick={handleGameRequest} class="btn btn-primary">
				<Plus class="w-4 h-4 mr-2" />
				Request Game
			</button>
		</div>

		{#if isLoading}
			<div class="flex justify-center items-center py-20">
				<Loading size="lg" text="Loading games..." />
			</div>
		{:else if filteredGames.length === 0}
			<!-- Empty State -->
			<div class="text-center py-20">
				<Gamepad2 class="w-16 h-16 mx-auto text-text-muted mb-6 opacity-50" />
				<h2 class="text-2xl font-semibold text-text-primary mb-2">
					{searchQuery ? 'No games found' : 'No games available'}
				</h2>
				<p class="text-text-secondary mb-6 max-w-md mx-auto">
					{searchQuery
						? 'Try adjusting your search terms or browse all games.'
						: 'Be the first to request a game to get the community started!'}
				</p>
				{#if searchQuery}
					<button
						onclick={() => {
							searchQuery = '';
							filterGames();
						}}
						class="btn btn-outline"
					>
						Clear Search
					</button>
				{:else}
					<button onclick={handleGameRequest} class="btn btn-primary">
						<Plus class="w-4 h-4 mr-2" />
						Request Game
					</button>
				{/if}
			</div>
		{:else}
			<!-- Games Grid -->
			<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
				{#each filteredGames as game (game.id)}
					<GameCard {game} variant="default" />
				{/each}
			</div>

			<!-- Pagination -->
			{#if !searchQuery && totalPages > 1}
				<div class="flex justify-center items-center space-x-2 mt-12">
					<button
						onclick={() => changePage(currentPage - 1)}
						disabled={currentPage <= 1}
						class="btn btn-outline btn-sm disabled:opacity-50 disabled:cursor-not-allowed"
					>
						Previous
					</button>

					<!-- Page Numbers -->
					{#each Array.from({ length: Math.min(7, totalPages) }, (_, i) => {
						const start = Math.max(1, currentPage - 3);
						const end = Math.min(totalPages, start + 6);
						return start + i <= end ? start + i : null;
					}).filter(Boolean) as pageNum (pageNum)}
						<button
							onclick={() => changePage(pageNum)}
							class="btn btn-sm {currentPage === pageNum ? 'btn-primary' : 'btn-outline'}"
						>
							{pageNum}
						</button>
					{/each}

					<button
						onclick={() => changePage(currentPage + 1)}
						disabled={currentPage >= totalPages}
						class="btn btn-outline btn-sm disabled:opacity-50 disabled:cursor-not-allowed"
					>
						Next
					</button>
				</div>
			{/if}
		{/if}

		<!-- Game Request CTA -->
		{#if !searchQuery && filteredGames.length > 0}
			<div class="mt-16 text-center p-8 bg-slate-800/50 rounded-xl border border-slate-700">
				<h3 class="text-xl font-semibold text-text-primary mb-2">Don't see your game?</h3>
				<p class="text-text-secondary mb-4">
					Request support for your favorite game and help grow the community.
				</p>
				<button onclick={handleGameRequest} class="btn btn-primary m-auto">
					<Plus class="w-4 h-4 mr-2" />
					Request Game Support
				</button>
			</div>
		{/if}
	</div>
</div>
