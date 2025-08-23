<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import {
		Search,
		Filter,
		ChevronDown,
		Download,
		Heart,
		Calendar,
		Package,
		User,
		X,
		Grid,
		List
	} from 'lucide-svelte';
	import Loading from '$lib/components/Loading.svelte';

	interface Mod {
		id: number;
		name: string;
		slug: string;
		description: string;
		short_description: string;
		icon: string;
		version: string;
		game_version: string;
		downloads: number;
		likes: number;
		created_at: string;
		updated_at: string;
		game: {
			id: number;
			name: string;
			slug: string;
			icon: string;
		};
		owner: {
			id: number;
			username: string;
			display_name: string;
			avatar: string;
		};
		is_liked: boolean;
	}

	interface Game {
		id: number;
		name: string;
		slug: string;
		icon: string;
	}

	let searchQuery = '';
	let selectedGameId = 0;
	let results: Mod[] = [];
	let games: Game[] = [];
	let loading = false;
	let searchPerformed = false;
	let currentPage = 1;
	let totalPages = 0;
	let totalResults = 0;
	let showFilters = false;
	let viewMode = 'grid'; // grid, list

	// Load search query from URL params
	onMount(async () => {
		const urlQuery = $page.url.searchParams.get('q');
		const urlGameId = $page.url.searchParams.get('game_id');
		const urlPage = $page.url.searchParams.get('page');

		if (urlQuery) {
			searchQuery = urlQuery;
			selectedGameId = urlGameId ? parseInt(urlGameId) : 0;
			currentPage = urlPage ? parseInt(urlPage) : 1;
			await performSearch();
		}

		// Load available games for filter
		await loadGames();
	});

	async function loadGames() {
		try {
			const response = await fetch('/api/games');
			const data = await response.json();

			if (data.success) {
				games = data.data;
			}
		} catch (error) {
			console.error('Failed to load games:', error);
		}
	}

	async function performSearch() {
		if (!searchQuery.trim()) {
			return;
		}

		loading = true;
		searchPerformed = true;

		try {
			const params = new URLSearchParams({
				q: searchQuery.trim(),
				page: currentPage.toString(),
				per_page: '20'
			});

			if (selectedGameId > 0) {
				params.set('game_id', selectedGameId.toString());
			}

			const response = await fetch(`/api/search/mods?${params}`);
			const data = await response.json();

			if (data.success) {
				results = data.data.data || [];
				totalPages = data.data.total_pages || 0;
				totalResults = data.data.total || 0;

				// Update URL without triggering a page reload
				const url = new URL(window.location.href);
				url.searchParams.set('q', searchQuery.trim());
				url.searchParams.set('page', currentPage.toString());
				if (selectedGameId > 0) {
					url.searchParams.set('game_id', selectedGameId.toString());
				} else {
					url.searchParams.delete('game_id');
				}
				window.history.replaceState({}, '', url.toString());
			} else {
				console.error('Search failed:', data.error);
			}
		} catch (error) {
			console.error('Search error:', error);
		} finally {
			loading = false;
		}
	}

	function handleSearch(event: Event) {
		event.preventDefault();
		currentPage = 1;
		performSearch();
	}

	function handlePageChange(newPage: number) {
		currentPage = newPage;
		performSearch();
		window.scrollTo({ top: 0, behavior: 'smooth' });
	}

	function handleGameFilter() {
		currentPage = 1;
		performSearch();
	}

	function clearFilters() {
		searchQuery = '';
		selectedGameId = 0;
		currentPage = 1;
		results = [];
		searchPerformed = false;
		updateURL();
	}

	function updateURL() {
		const params = new URLSearchParams();
		if (searchQuery.trim()) params.set('q', searchQuery.trim());
		if (selectedGameId > 0) params.set('game_id', selectedGameId.toString());
		if (currentPage > 1) params.set('page', currentPage.toString());

		const newURL = `/search${params.toString() ? `?${params.toString()}` : ''}`;
		goto(newURL, { replaceState: true });
	}

	function formatNumber(num: number): string {
		if (num >= 1000000) {
			return (num / 1000000).toFixed(1) + 'M';
		} else if (num >= 1000) {
			return (num / 1000).toFixed(1) + 'K';
		}
		return num.toString();
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString();
	}

	function getGameName(gameId: number): string {
		const game = games.find((g) => g.id === gameId);
		return game ? game.name : 'All Games';
	}
</script>

<svelte:head>
	<title>Search{searchQuery ? ` - ${searchQuery}` : ''} - Azurite</title>
	<meta name="description" content="Search for game mods on Azurite mod hosting platform" />
</svelte:head>

<div class="min-h-screen bg-background-primary">
	<!-- Header -->
	<div class="bg-gradient-to-r from-primary-600/20 to-primary-700/20 border-b border-slate-700">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
			<div class="text-center">
				<h1 class="text-4xl font-bold text-text-primary mb-4">
					<Search class="w-10 h-10 inline-block mr-2 mb-1" />
					Search Mods
				</h1>
				<p class="text-text-secondary text-lg max-w-2xl mx-auto mb-8">
					Discover amazing modifications for your favorite games. Search by name, description,
					author, or keywords.
				</p>

				<!-- Main Search -->
				<div class="max-w-2xl mx-auto">
					<form on:submit|preventDefault={handleSearch}>
						<div class="relative">
							<Search
								class="absolute left-4 top-1/2 transform -translate-y-1/2 text-text-muted w-6 h-6"
							/>
							<input
								type="search"
								placeholder="Search for mods, descriptions, or keywords..."
								bind:value={searchQuery}
								class="input pl-12 pr-4 py-4 w-full text-lg"
								disabled={loading}
							/>
							{#if searchQuery}
								<button
									type="button"
									on:click={() => {
										searchQuery = '';
										results = [];
										searchPerformed = false;
									}}
									class="absolute right-4 top-1/2 transform -translate-y-1/2 text-text-muted hover:text-text-primary"
								>
									<X class="w-5 h-5" />
								</button>
							{/if}
						</div>
						<button
							type="submit"
							disabled={loading || !searchQuery.trim()}
							class="btn btn-primary mt-4 px-8"
						>
							{#if loading}
								<Loading size="sm" inline />
							{:else}
								Search
							{/if}
						</button>
					</form>
				</div>
			</div>
		</div>
	</div>

	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
		<!-- Filters and Controls -->
		{#if searchPerformed || searchQuery}
			<div class="flex flex-col lg:flex-row gap-4 mb-8">
				<!-- Filters -->
				<div class="flex-1">
					<div class="flex flex-wrap items-center gap-3">
						<!-- Game Filter -->
						<select
							bind:value={selectedGameId}
							on:change={handleGameFilter}
							class="select min-w-[150px]"
							disabled={loading}
						>
							<option value={0}>All Games</option>
							{#each games as game (game.id)}
								<option value={game.id}>{game.name}</option>
							{/each}
						</select>

						<!-- Filter Toggle -->
						<button
							on:click={() => (showFilters = !showFilters)}
							class="btn btn-outline"
							disabled={loading}
						>
							<Filter class="w-4 h-4 mr-2" />
							Filters
							<ChevronDown
								class="w-4 h-4 ml-2 transition-transform {showFilters ? 'rotate-180' : ''}"
							/>
						</button>

						{#if selectedGameId > 0}
							<button
								on:click={clearFilters}
								class="btn btn-outline text-red-400 hover:text-red-300 hover:border-red-500"
							>
								<X class="w-4 h-4 mr-2" />
								Clear Filters
							</button>
						{/if}
					</div>
				</div>

				<!-- View Controls -->
				<div class="flex items-center space-x-4">
					<!-- Results Count -->
					{#if searchPerformed && !loading}
						<span class="text-text-muted text-sm whitespace-nowrap">
							{totalResults} result{totalResults !== 1 ? 's' : ''}
						</span>
					{/if}

					<!-- View Mode Toggle -->
					<div class="flex border border-slate-600 rounded-lg p-1">
						<button
							on:click={() => (viewMode = 'grid')}
							class="p-2 rounded {viewMode === 'grid'
								? 'bg-primary-600 text-white'
								: 'text-text-muted hover:text-text-primary'} transition-colors"
							title="Grid View"
						>
							<Grid class="w-4 h-4" />
						</button>
						<button
							on:click={() => (viewMode = 'list')}
							class="p-2 rounded {viewMode === 'list'
								? 'bg-primary-600 text-white'
								: 'text-text-muted hover:text-text-primary'} transition-colors"
							title="List View"
						>
							<List class="w-4 h-4" />
						</button>
					</div>
				</div>
			</div>

			<!-- Advanced Filters Panel -->
			{#if showFilters}
				<div class="card p-6 mb-8">
					<div class="flex items-center justify-between mb-4">
						<h3 class="text-lg font-semibold text-text-primary">Advanced Filters</h3>
						<button
							on:click={() => (showFilters = false)}
							class="text-text-muted hover:text-text-primary"
						>
							<X class="w-5 h-5" />
						</button>
					</div>

					<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
						<div>
							<label
								for="game-filter-advanced"
								class="block text-sm font-medium text-text-primary mb-2"
							>
								Game
							</label>
							<select
								id="game-filter-advanced"
								bind:value={selectedGameId}
								on:change={handleGameFilter}
								class="select w-full"
								disabled={loading}
							>
								<option value={0}>All Games</option>
								{#each games as game (game.id)}
									<option value={game.id}>{game.name}</option>
								{/each}
							</select>
						</div>
						<!-- Add more filter options here in the future -->
					</div>
				</div>
			{/if}
		{/if}

		<!-- Search Results -->
		{#if loading}
			<div class="flex justify-center items-center py-20">
				<Loading size="lg" text="Searching..." />
			</div>
		{:else if searchPerformed}
			{#if results.length > 0}
				<!-- Results Header -->
				<div class="mb-6">
					<p class="text-text-secondary">
						Found {totalResults} mod{totalResults !== 1 ? 's' : ''}
						{#if searchQuery}for "{searchQuery}"{/if}
						{#if selectedGameId > 0}in {getGameName(selectedGameId)}{/if}
					</p>
				</div>

				<!-- Results Display -->
				{#if viewMode === 'grid'}
					<!-- Grid View -->
					<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6 mb-8">
						{#each results as mod (mod.id)}
							<div class="card card-hover">
								<a href="/games/{mod.game.slug}/mods/{mod.slug}" class="block">
									<div class="p-4">
										<!-- Mod Icon -->
										<div
											class="aspect-video bg-background-secondary rounded-lg flex items-center justify-center overflow-hidden mb-4"
										>
											{#if mod.icon}
												<img src={mod.icon} alt={mod.name} class="w-full h-full object-cover" />
											{:else}
												<Package class="w-12 h-12 text-text-muted" />
											{/if}
										</div>

										<!-- Game Badge -->
										<div class="flex items-center gap-2 mb-2">
											{#if mod.game.icon}
												<img src={mod.game.icon} alt={mod.game.name} class="w-4 h-4" />
											{/if}
											<span class="text-sm text-primary-400 font-medium">{mod.game.name}</span>
										</div>

										<h3 class="font-semibold text-text-primary mb-2 line-clamp-2">{mod.name}</h3>
										<p class="text-sm text-text-secondary mb-4 line-clamp-3">
											{mod.short_description}
										</p>

										<!-- Stats -->
										<div class="flex items-center justify-between text-sm text-text-muted mb-3">
											<div class="flex items-center gap-4">
												<span class="flex items-center gap-1">
													<Download class="w-3 h-3" />
													{formatNumber(mod.downloads)}
												</span>
												<span class="flex items-center gap-1">
													<Heart class="w-3 h-3" />
													{formatNumber(mod.likes)}
												</span>
											</div>
											<span>{formatDate(mod.updated_at)}</span>
										</div>

										<!-- Author -->
										<div class="flex items-center gap-2">
											{#if mod.owner.avatar}
												<img
													src={mod.owner.avatar}
													alt={mod.owner.display_name}
													class="w-6 h-6 rounded-full"
												/>
											{:else}
												<div
													class="w-6 h-6 bg-primary-600 rounded-full flex items-center justify-center text-xs text-white"
												>
													{mod.owner.display_name.charAt(0).toUpperCase()}
												</div>
											{/if}
											<span class="text-sm text-text-secondary">by {mod.owner.display_name}</span>
										</div>
									</div>
								</a>
							</div>
						{/each}
					</div>
				{:else}
					<!-- List View -->
					<div class="space-y-4 mb-8">
						{#each results as mod (mod.id)}
							<div class="card card-hover">
								<a href="/games/{mod.game.slug}/mods/{mod.slug}" class="block">
									<div class="p-6">
										<div class="flex items-start space-x-4">
											<!-- Mod Icon -->
											<div class="flex-shrink-0">
												{#if mod.icon}
													<img
														src={mod.icon}
														alt={mod.name}
														class="w-16 h-16 rounded-lg border border-slate-600"
													/>
												{:else}
													<div
														class="w-16 h-16 bg-slate-600 rounded-lg flex items-center justify-center"
													>
														<Package class="w-8 h-8 text-text-muted" />
													</div>
												{/if}
											</div>

											<!-- Mod Info -->
											<div class="flex-1 min-w-0">
												<!-- Game Badge -->
												<div class="flex items-center gap-2 mb-1">
													{#if mod.game.icon}
														<img src={mod.game.icon} alt={mod.game.name} class="w-4 h-4" />
													{/if}
													<span class="text-sm text-primary-400 font-medium">{mod.game.name}</span>
												</div>

												<h3 class="text-xl font-semibold text-text-primary mb-2">{mod.name}</h3>
												<p class="text-text-secondary mb-4 line-clamp-2">{mod.short_description}</p>

												<div class="flex items-center justify-between">
													<div class="flex items-center space-x-6 text-sm text-text-muted">
														<span class="flex items-center gap-1">
															<Download class="w-4 h-4" />
															{formatNumber(mod.downloads)} downloads
														</span>
														<span class="flex items-center gap-1">
															<Heart class="w-4 h-4" />
															{formatNumber(mod.likes)} likes
														</span>
														<span class="flex items-center gap-1">
															<Calendar class="w-4 h-4" />
															{formatDate(mod.updated_at)}
														</span>
													</div>

													<!-- Author -->
													<div class="flex items-center gap-2">
														{#if mod.owner.avatar}
															<img
																src={mod.owner.avatar}
																alt={mod.owner.display_name}
																class="w-6 h-6 rounded-full"
															/>
														{:else}
															<div
																class="w-6 h-6 bg-primary-600 rounded-full flex items-center justify-center text-xs text-white"
															>
																{mod.owner.display_name.charAt(0).toUpperCase()}
															</div>
														{/if}
														<span class="text-sm text-text-secondary"
															>by {mod.owner.display_name}</span
														>
													</div>
												</div>
											</div>
										</div>
									</div>
								</a>
							</div>
						{/each}
					</div>
				{/if}

				<!-- Pagination -->
				{#if totalPages > 1}
					<div class="flex justify-center items-center space-x-2">
						<button
							on:click={() => handlePageChange(currentPage - 1)}
							disabled={currentPage <= 1 || loading}
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
								on:click={() => handlePageChange(pageNum)}
								disabled={loading}
								class="btn btn-sm {pageNum === currentPage ? 'btn-primary' : 'btn-outline'}"
							>
								{pageNum}
							</button>
						{/each}

						<button
							on:click={() => handlePageChange(currentPage + 1)}
							disabled={currentPage >= totalPages || loading}
							class="btn btn-outline btn-sm disabled:opacity-50 disabled:cursor-not-allowed"
						>
							Next
						</button>
					</div>
				{/if}
			{:else}
				<!-- No Results -->
				<div class="text-center py-20">
					<Search class="w-16 h-16 mx-auto text-text-muted mb-6 opacity-50" />
					<h2 class="text-2xl font-semibold text-text-primary mb-2">No mods found</h2>
					<p class="text-text-secondary mb-6 max-w-md mx-auto">
						{#if searchQuery}
							No results found for "{searchQuery}"
							{#if selectedGameId > 0}
								in {getGameName(selectedGameId)}
							{/if}
						{:else}
							Try searching for something
						{/if}
					</p>
					<div class="space-y-2 text-sm text-text-muted max-w-md mx-auto mb-8">
						<p class="font-medium">Try these tips:</p>
						<ul class="text-left space-y-1">
							<li>• Check spelling and try different keywords</li>
							<li>• Use more general search terms</li>
							<li>• Try searching in all games</li>
							<li>• Browse popular mods instead</li>
						</ul>
					</div>
					<div class="flex flex-col sm:flex-row gap-3 justify-center">
						<button on:click={clearFilters} class="btn btn-outline"> Clear Search </button>
						<a href="/browse" class="btn btn-primary"> Browse All Mods </a>
					</div>
				</div>
			{/if}
		{:else}
			<!-- Welcome State -->
			<div class="text-center py-20">
				<Search class="w-16 h-16 mx-auto text-text-muted mb-6 opacity-50" />
				<h2 class="text-2xl font-semibold text-text-primary mb-4">Find the perfect mod</h2>
				<p class="text-text-secondary mb-12 max-w-md mx-auto">
					Search through thousands of mods across all supported games. Use the search bar above to
					get started.
				</p>

				<!-- Search Suggestions -->
				<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 max-w-4xl mx-auto">
					<div class="card">
						<div class="p-6 text-center">
							<div
								class="w-12 h-12 bg-primary-600/20 rounded-lg flex items-center justify-center mx-auto mb-4"
							>
								<Package class="w-6 h-6 text-primary-400" />
							</div>
							<h3 class="font-semibold text-text-primary mb-2">Search by Name</h3>
							<p class="text-sm text-text-secondary">
								Find specific mods by their name or partial matches
							</p>
						</div>
					</div>

					<div class="card">
						<div class="p-6 text-center">
							<div
								class="w-12 h-12 bg-primary-600/20 rounded-lg flex items-center justify-center mx-auto mb-4"
							>
								<User class="w-6 h-6 text-primary-400" />
							</div>
							<h3 class="font-semibold text-text-primary mb-2">Find by Author</h3>
							<p class="text-sm text-text-secondary">
								Search for mods created by your favorite developers
							</p>
						</div>
					</div>

					<div class="card">
						<div class="p-6 text-center">
							<div
								class="w-12 h-12 bg-primary-600/20 rounded-lg flex items-center justify-center mx-auto mb-4"
							>
								<Filter class="w-6 h-6 text-primary-400" />
							</div>
							<h3 class="font-semibold text-text-primary mb-2">Filter by Game</h3>
							<p class="text-sm text-text-secondary">
								Narrow down results to specific games you're interested in
							</p>
						</div>
					</div>
				</div>

				<!-- Popular Searches (Mock Data) -->
				<div class="mt-12">
					<h3 class="text-lg font-semibold text-text-primary mb-4">Popular Searches</h3>
					<div class="flex flex-wrap justify-center gap-2">
						{#each ['OptiFine', 'JEI', 'Biomes O Plenty', 'Applied Energistics', 'Tinkers Construct', 'JourneyMap', 'Iron Chests', 'Thermal Expansion'] as term}
							<button
								on:click={() => {
									searchQuery = term;
									handleSearch(new Event('submit'));
								}}
								class="px-4 py-2 text-sm bg-slate-800 hover:bg-slate-700 text-text-secondary hover:text-text-primary rounded-full border border-slate-600 hover:border-slate-500 transition-colors"
							>
								{term}
							</button>
						{/each}
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>

<style>
	.line-clamp-2 {
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}

	.line-clamp-3 {
		display: -webkit-box;
		-webkit-line-clamp: 3;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
</style>
