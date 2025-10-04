<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { gamesApi } from '$lib/api/client';
	import { toast } from '$lib/stores/notifications';
	import { isAuthenticated } from '$lib/stores/auth';
	import Loading from '$lib/components/Loading.svelte';
	import { ModCard } from '$lib/components/cards';
	import type { Game, Mod, Tag } from '$lib/types';
	import {
		Search,
		Filter,
		Download,
		Heart,
		Calendar,
		Package,
		BookOpen,
		Grid,
		List
	} from 'lucide-svelte';

	// URL params
	$: gameSlug = $page.params.slug;

	// Data
	let game: Game | null = null;
	let mods: Mod[] = [];
	let tags: Tag[] = [];
	let isLoadingGame = true;
	let isLoadingMods = true;
	let isLoadingTags = true;

	// Filters and search
	let searchQuery = '';
	let selectedTags: string[] = [];
	let sortBy = 'newest'; // newest, oldest, downloads, likes
	let viewMode = 'grid'; // grid, list

	// Pagination
	let currentPage = 1;
	let totalPages = 1;
	let totalMods = 0;
	let perPage = 20;

	// UI state
	let showFilters = false;

	// Sort options
	const sortOptions = [
		{ value: 'newest', label: 'Newest First', icon: Calendar },
		{ value: 'oldest', label: 'Oldest First', icon: Calendar },
		{ value: 'downloads', label: 'Most Downloads', icon: Download },
		{ value: 'likes', label: 'Most Liked', icon: Heart }
	];

	onMount(async () => {
		await loadGameData();
		await loadMods();
		await loadTags();
	});

	// Load game information
	async function loadGameData() {
		isLoadingGame = true;
		try {
			const response = await gamesApi.getGame(gameSlug);
			if (response.success && response.data) {
				game = response.data;
			} else {
				toast.error('Game not found', 'The requested game could not be found.');
				goto('/games');
			}
		} catch (error) {
			console.error('Error loading game:', error);
			toast.error('Error', 'Failed to load game information');
		} finally {
			isLoadingGame = false;
		}
	}

	// Load mods for this game
	async function loadMods() {
		isLoadingMods = true;
		try {
			const response = await gamesApi.getGameMods(gameSlug, {
				page: currentPage,
				per_page: perPage,
				sort: sortBy,
				tags: selectedTags.length > 0 ? selectedTags.join(',') : undefined
			});

			if (response.success && response.data) {
				mods = response.data.data || [];
				totalPages = response.data.total_pages || 1;
				totalMods = response.data.total || 0;
				currentPage = response.data.page || 1;
			} else {
				console.error('Failed to load mods:', response.error);
			}
		} catch (error) {
			console.error('Error loading mods:', error);
		} finally {
			isLoadingMods = false;
		}
	}

	// Load available tags for this game
	async function loadTags() {
		isLoadingTags = true;
		try {
			const response = await gamesApi.getGameTags(gameSlug);
			if (response.success && response.data) {
				tags = response.data || [];
			}
		} catch (error) {
			console.error('Error loading tags:', error);
		} finally {
			isLoadingTags = false;
		}
	}

	// Handle search
	function handleSearch() {
		// Filter mods locally if search is active
		filteredMods = searchQuery.trim()
			? mods.filter(
					(mod) =>
						mod.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
						mod.short_description?.toLowerCase().includes(searchQuery.toLowerCase())
				)
			: mods;
	}

	// Handle tag selection
	function toggleTag(tagSlug: string) {
		if (selectedTags.includes(tagSlug)) {
			selectedTags = selectedTags.filter((t) => t !== tagSlug);
		} else {
			selectedTags = [...selectedTags, tagSlug];
		}
		currentPage = 1;
		loadMods();
	}

	// Handle sort change
	function handleSortChange(newSort: string) {
		sortBy = newSort;
		currentPage = 1;
		loadMods();
	}

	// Handle page change
	function changePage(newPage: number) {
		if (newPage >= 1 && newPage <= totalPages && newPage !== currentPage) {
			currentPage = newPage;
			loadMods();
			// Scroll to top
			window.scrollTo({ top: 0, behavior: 'smooth' });
		}
	}

	// Format numbers
	function formatNumber(num: number): string {
		if (num >= 1000000) {
			return (num / 1000000).toFixed(1) + 'M';
		}
		if (num >= 1000) {
			return (num / 1000).toFixed(1) + 'K';
		}
		return num.toString();
	}

	// Format date
	function formatDate(dateString: string): string {
		const date = new Date(dateString);
		const now = new Date();
		const diffTime = Math.abs(now.getTime() - date.getTime());
		const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

		if (diffDays === 1) return '1 day ago';
		if (diffDays < 7) return `${diffDays} days ago`;
		if (diffDays < 30) return `${Math.ceil(diffDays / 7)} weeks ago`;
		return date.toLocaleDateString();
	}

	// Filtered mods for search
	let filteredMods: any[] = [];
	$: filteredMods = searchQuery.trim()
		? mods.filter(
				(mod) =>
					mod.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
					mod.short_description?.toLowerCase().includes(searchQuery.toLowerCase())
			)
		: mods;

	// Clear all filters
	function clearFilters() {
		searchQuery = '';
		selectedTags = [];
		sortBy = 'newest';
		currentPage = 1;
		loadMods();
	}
</script>

<svelte:head>
	<title>{game?.name || 'Game'} - Azurite</title>
	<meta
		name="description"
		content={game
			? `Discover and download mods for ${game.name}. Browse ${game.mod_count || 0} modifications created by the community.`
			: 'Browse game modifications on Azurite.'}
	/>
</svelte:head>

{#if isLoadingGame}
	<div class="min-h-screen flex items-center justify-center">
		<Loading size="lg" text="Loading game..." />
	</div>
{:else if game}
	<div class="min-h-screen bg-background-primary">
		<!-- Game Header -->
		<div class="bg-gradient-to-r from-slate-800/50 to-slate-700/50 border-b border-slate-700">
			<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
				<div
					class="flex flex-col md:flex-row items-start md:items-center space-y-6 md:space-y-0 md:space-x-8"
				>
					<!-- Game Icon -->
					<div class="flex-shrink-0">
						{#if game.icon}
							<img
								src={game.icon}
								alt={game.name}
								class="w-24 h-24 md:w-32 md:h-32 rounded-2xl border border-slate-600 shadow-xl"
							/>
						{:else}
							<div
								class="w-24 h-24 md:w-32 md:h-32 bg-gradient-to-br from-slate-600 to-slate-700 rounded-2xl flex items-center justify-center shadow-xl"
							>
								<Package class="w-12 h-12 md:w-16 md:h-16 text-text-muted" />
							</div>
						{/if}
					</div>

					<!-- Game Info -->
					<div class="flex-1">
						<h1 class="text-3xl md:text-4xl font-bold text-text-primary mb-2">
							{game.name}
						</h1>

						{#if game.description}
							<p class="text-text-secondary text-lg mb-4 leading-relaxed">
								{game.description}
							</p>
						{/if}

						<!-- Stats -->
						<div class="flex flex-wrap items-center gap-6 text-text-muted">
							<div class="flex items-center">
								<Package class="w-5 h-5 mr-2" />
								<span class="font-medium">{formatNumber(game.mod_count || 0)}</span>
								<span class="ml-1">mods</span>
							</div>
							{#if game.created_at}
								<div class="flex items-center">
									<Calendar class="w-5 h-5 mr-2" />
									<span>Since {new Date(game.created_at).getFullYear()}</span>
								</div>
							{/if}
						</div>
					</div>

					<!-- Actions -->
					<div class="flex flex-col sm:flex-row gap-3">
						<a href="/docs/{game.slug}" class="btn btn-outline">
							<BookOpen class="w-4 h-4 mr-2" />
							Documentation
						</a>
						{#if $isAuthenticated}
							<a href="/dashboard/mods/create?game={game.slug}" class="btn btn-primary">
								<Package class="w-4 h-4 mr-2" />
								Upload Mod
							</a>
						{/if}
					</div>
				</div>
			</div>
		</div>

		<!-- Content -->
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
			<!-- Search and Filters -->
			<div class="mb-8">
				<div class="flex flex-col lg:flex-row gap-4 mb-4">
					<!-- Search -->
					<div class="flex-1">
						<div class="relative">
							<Search
								class="absolute left-3 top-1/2 transform -translate-y-1/2 text-text-muted w-5 h-5"
							/>
							<input
								type="search"
								placeholder="Search mods..."
								bind:value={searchQuery}
								on:input={handleSearch}
								class="input pl-10 pr-4"
							/>
						</div>
					</div>

					<!-- Sort -->
					<div class="flex items-center space-x-2">
						<label for="sort-select" class="text-text-secondary text-sm whitespace-nowrap"
							>Sort by:</label
						>
						<select
							id="sort-select"
							bind:value={sortBy}
							on:change={() => handleSortChange(sortBy)}
							class="select min-w-[150px]"
						>
							{#each sortOptions as option (option.value)}
								<option value={option.value}>{option.label}</option>
							{/each}
						</select>
					</div>

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

					<!-- Filter Toggle -->
					<button on:click={() => (showFilters = !showFilters)} class="btn btn-outline">
						<Filter class="w-4 h-4 mr-2" />
						Filters
						{#if selectedTags.length > 0}
							<span class="ml-2 bg-primary-600 text-white text-xs px-2 py-0.5 rounded-full">
								{selectedTags.length}
							</span>
						{/if}
					</button>
				</div>

				<!-- Filters Panel -->
				{#if showFilters}
					<div class="card p-6 mb-4">
						<div class="flex items-center justify-between mb-4">
							<h3 class="text-lg font-semibold text-text-primary">Filters</h3>
							<button
								on:click={clearFilters}
								class="text-primary-400 hover:text-primary-300 text-sm font-medium transition-colors"
							>
								Clear All
							</button>
						</div>

						<!-- Tags -->
						{#if !isLoadingTags && tags.length > 0}
							<div>
								<h4 class="text-sm font-medium text-text-primary mb-3">Tags</h4>
								<div class="flex flex-wrap gap-2">
									{#each tags as tag (tag.id)}
										<button
											on:click={() => toggleTag(tag.slug)}
											class="px-3 py-1 text-sm rounded-full border transition-colors {selectedTags.includes(
												tag.slug
											)
												? 'bg-primary-600 text-white border-primary-600'
												: 'border-slate-600 text-text-secondary hover:text-text-primary hover:border-slate-500'}"
										>
											{tag.name}
										</button>
									{/each}
								</div>
							</div>
						{/if}
					</div>
				{/if}

				<!-- Active Filters Display -->
				{#if searchQuery || selectedTags.length > 0}
					<div class="flex flex-wrap items-center gap-2 mb-4">
						<span class="text-text-secondary text-sm">Active filters:</span>
						{#if searchQuery}
							<span class="badge badge-secondary">
								Search: {searchQuery}
								<button
									on:click={() => {
										searchQuery = '';
										handleSearch();
									}}
									class="ml-1 text-text-muted hover:text-text-primary"
								>
									×
								</button>
							</span>
						{/if}
						{#each selectedTags as tagSlug (tagSlug)}
							{@const tag = tags.find((t) => t.slug === tagSlug)}
							{#if tag}
								<span class="badge badge-secondary">
									{tag.name}
									<button
										on:click={() => toggleTag(tagSlug)}
										class="ml-1 text-text-muted hover:text-text-primary"
									>
										×
									</button>
								</span>
							{/if}
						{/each}
					</div>
				{/if}
			</div>

			<!-- Mods Count -->
			<div class="flex items-center justify-between mb-6">
				<p class="text-text-secondary">
					{#if searchQuery && filteredMods.length !== mods.length}
						Showing {filteredMods.length} of {totalMods} mods
					{:else}
						{totalMods} mods available
					{/if}
				</p>
			</div>

			<!-- Mods Grid/List -->
			{#if isLoadingMods}
				<div class="flex justify-center items-center py-20">
					<Loading size="lg" text="Loading mods..." />
				</div>
			{:else if filteredMods.length === 0}
				<!-- Empty State -->
				<div class="text-center py-20">
					<Package class="w-16 h-16 mx-auto text-text-muted mb-6 opacity-50" />
					<h2 class="text-2xl font-semibold text-text-primary mb-2">
						{searchQuery || selectedTags.length > 0
							? 'No mods match your filters'
							: 'No mods available yet'}
					</h2>
					<p class="text-text-secondary mb-6 max-w-md mx-auto">
						{searchQuery || selectedTags.length > 0
							? 'Try adjusting your search terms or clearing some filters.'
							: `Be the first to create a mod for ${game.name}!`}
					</p>
					{#if searchQuery || selectedTags.length > 0}
						<button on:click={clearFilters} class="btn btn-outline"> Clear Filters </button>
					{:else if $isAuthenticated}
						<a href="/dashboard/mods/create?game={game.slug}" class="btn btn-primary">
							<Package class="w-4 h-4 mr-2" />
							Upload First Mod
						</a>
					{/if}
				</div>
			{:else if viewMode === 'grid'}
				<!-- Grid View -->
				<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
					{#each filteredMods as mod (mod.id)}
						<ModCard {mod} variant="default" showGame={false} />
					{/each}
				</div>
			{:else}
				<!-- List View -->
				<div class="space-y-4">
					{#each filteredMods as mod (mod.id)}
						<ModCard {mod} variant="list" showGame={false} />
					{/each}
				</div>
			{/if}

			<!-- Pagination -->
			{#if !searchQuery && totalPages > 1}
				<div class="flex justify-center items-center space-x-2 mt-12">
					<button
						on:click={() => changePage(currentPage - 1)}
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
							on:click={() => changePage(pageNum)}
							class="btn btn-sm {currentPage === pageNum ? 'btn-primary' : 'btn-outline'}"
						>
							{pageNum}
						</button>
					{/each}

					<button
						on:click={() => changePage(currentPage + 1)}
						disabled={currentPage >= totalPages}
						class="btn btn-outline btn-sm disabled:opacity-50 disabled:cursor-not-allowed"
					>
						Next
					</button>
				</div>

				<p class="text-center text-text-muted text-sm mt-4">
					Page {currentPage} of {totalPages} • {totalMods} total mods
				</p>
			{/if}
		</div>
	</div>
{/if}

<style>
	.line-clamp-1 {
		display: -webkit-box;
		-webkit-line-clamp: 1;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}

	.line-clamp-2 {
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
</style>
