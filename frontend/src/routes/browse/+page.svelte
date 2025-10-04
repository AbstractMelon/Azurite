<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { gamesApi } from '$lib/api/client';
	import { toast } from '$lib/stores/notifications';
	import Loading from '$lib/components/Loading.svelte';
	import { ModCard } from '$lib/components/cards';
	import type { Game, Mod } from '$lib/types';
	import {
		Search,
		Filter,
		Download,
		Heart,
		Calendar,
		Package,
		Grid,
		List,
		SortAsc,
		Clock,
		X
	} from 'lucide-svelte';
	import { SvelteURLSearchParams } from 'svelte/reactivity';

	// Data
	let mods: Mod[] = [];
	let games: Game[] = [];
	let isLoadingMods = true;

	// Search and filters
	let searchQuery = '';
	let selectedGame = '';
	let selectedTags: string[] = [];
	let sortBy = 'newest'; // newest, oldest, downloads, likes, name
	let sortOrder = 'desc'; // asc, desc
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
		{ value: 'newest', label: 'Newest First', icon: Clock },
		{ value: 'oldest', label: 'Oldest First', icon: Calendar },
		{ value: 'downloads', label: 'Most Downloads', icon: Download },
		{ value: 'likes', label: 'Most Liked', icon: Heart },
		{ value: 'name', label: 'Name A-Z', icon: SortAsc }
	];

	// Available tags (mock data)
	const availableTags = [
		'Technology',
		'Magic',
		'Adventure',
		'Tools',
		'Weapons',
		'Armor',
		'Building',
		'Decoration',
		'Food',
		'Animals',
		'Biomes',
		'Dimensions',
		'Performance',
		'Graphics',
		'UI',
		'Quality of Life',
		'Automation',
		'Transportation',
		'Storage',
		'Energy'
	];

	onMount(async () => {
		// Get initial search from URL params
		const urlParams = new URLSearchParams($page.url.search);
		searchQuery = urlParams.get('q') || '';
		selectedGame = urlParams.get('game') || '';
		const tagParam = urlParams.get('tags');
		if (tagParam) {
			selectedTags = tagParam.split(',').filter(Boolean);
		}

		await Promise.all([loadGames(), loadMods()]);
	});

	// Load available games
	async function loadGames() {
		try {
			const response = await gamesApi.getGames({ per_page: 50 });
			if (response.success && response.data) {
				games = response.data.data || [];
			}
		} catch (error) {
			console.error('Error loading games:', error);
		}
	}

	// Load mods with current filters
	async function loadMods() {
		isLoadingMods = true;
		try {
			// Mock data for demonstration
			const mockMods = [
				{
					id: 1,
					name: 'OptiFine HD',
					slug: 'optifine-hd',
					short_description: 'Boost FPS and add advanced graphics features to your game',
					game: { name: 'Minecraft', slug: 'minecraft' },
					owner: { username: 'sp614x', display_name: 'sp614x' },
					downloads: 125000,
					likes: 8500,
					created_at: '2023-06-01T10:00:00Z',
					updated_at: '2024-01-10T14:00:00Z',
					tags: ['Performance', 'Graphics'],
					icon: null,
					version: '1.20.4'
				},
				{
					id: 2,
					name: 'JEI - Just Enough Items',
					slug: 'jei',
					short_description: 'Recipe and item viewer with search functionality',
					game: { name: 'Minecraft', slug: 'minecraft' },
					owner: { username: 'mezz', display_name: 'mezz' },
					downloads: 98000,
					likes: 6200,
					created_at: '2023-05-15T08:00:00Z',
					updated_at: '2024-01-08T12:00:00Z',
					tags: ['UI', 'Tools'],
					icon: null,
					version: '15.2.0'
				},
				{
					id: 3,
					name: 'Biomes O Plenty',
					slug: 'biomes-o-plenty',
					short_description: 'Adds 80+ new biomes to explore and discover',
					game: { name: 'Minecraft', slug: 'minecraft' },
					owner: { username: 'glitchfiend', display_name: 'Glitchfiend' },
					downloads: 87000,
					likes: 5400,
					created_at: '2023-04-20T15:30:00Z',
					updated_at: '2024-01-05T16:20:00Z',
					tags: ['Biomes', 'Adventure'],
					icon: null,
					version: '18.0.0'
				},
				{
					id: 4,
					name: 'Tinkers Construct',
					slug: 'tinkers-construct',
					short_description: 'Tool and weapon crafting system with customization',
					game: { name: 'Minecraft', slug: 'minecraft' },
					owner: { username: 'mdiyo', display_name: 'mDiyo' },
					downloads: 76000,
					likes: 4800,
					created_at: '2023-03-10T11:00:00Z',
					updated_at: '2024-01-03T10:15:00Z',
					tags: ['Tools', 'Weapons'],
					icon: null,
					version: '3.7.1'
				},
				{
					id: 5,
					name: 'Applied Energistics 2',
					slug: 'ae2',
					short_description: 'Advanced storage and automation system',
					game: { name: 'Minecraft', slug: 'minecraft' },
					owner: { username: 'algester', display_name: 'AlgorithmX2' },
					downloads: 65000,
					likes: 4200,
					created_at: '2023-02-28T09:45:00Z',
					updated_at: '2023-12-28T14:30:00Z',
					tags: ['Technology', 'Storage', 'Automation'],
					icon: null,
					version: '15.0.16'
				}
			];

			// Apply filters
			let filteredMods = mockMods;

			if (searchQuery) {
				const query = searchQuery.toLowerCase();
				filteredMods = filteredMods.filter(
					(mod) =>
						mod.name.toLowerCase().includes(query) ||
						mod.short_description.toLowerCase().includes(query) ||
						mod.owner.display_name.toLowerCase().includes(query)
				);
			}

			if (selectedGame) {
				filteredMods = filteredMods.filter((mod) => mod.game.slug === selectedGame);
			}

			if (selectedTags.length > 0) {
				filteredMods = filteredMods.filter((mod) =>
					selectedTags.some((tag) => mod.tags.includes(tag))
				);
			}

			// Apply sorting
			filteredMods.sort((a, b) => {
				let comparison = 0;
				switch (sortBy) {
					case 'name':
						comparison = a.name.localeCompare(b.name);
						break;
					case 'downloads':
						comparison = b.downloads - a.downloads;
						break;
					case 'likes':
						comparison = b.likes - a.likes;
						break;
					case 'oldest':
						comparison = new Date(a.created_at).getTime() - new Date(b.created_at).getTime();
						break;
					case 'newest':
					default:
						comparison = new Date(b.created_at).getTime() - new Date(a.created_at).getTime();
						break;
				}
				return sortOrder === 'asc' ? comparison : -comparison;
			});

			// Apply pagination
			const startIndex = (currentPage - 1) * perPage;
			const endIndex = startIndex + perPage;
			mods = filteredMods.slice(startIndex, endIndex);

			totalMods = filteredMods.length;
			totalPages = Math.ceil(totalMods / perPage);
		} catch (error) {
			console.error('Error loading mods:', error);
			toast.error('Error', 'Failed to load mods');
		} finally {
			isLoadingMods = false;
		}
	}

	// Handle search
	function handleSearch(event: Event) {
		// Prevent default form submission
		event.preventDefault();

		currentPage = 1;
		updateURL();
		loadMods();
	}

	// Handle filter change
	function handleFilterChange() {
		currentPage = 1;
		updateURL();
		loadMods();
	}

	// Update URL with current filters
	function updateURL() {
		const params = new SvelteURLSearchParams();
		if (searchQuery) params.set('q', searchQuery);
		if (selectedGame) params.set('game', selectedGame);
		if (selectedTags.length > 0) params.set('tags', selectedTags.join(','));

		const newURL = `/browse${params.toString() ? `?${params.toString()}` : ''}`;
		goto(newURL, { replaceState: true });
	}

	// Add tag filter
	function addTag(tag: string) {
		if (!selectedTags.includes(tag)) {
			selectedTags = [...selectedTags, tag];
			handleFilterChange();
		}
	}

	// Remove tag filter
	function removeTag(tag: string) {
		selectedTags = selectedTags.filter((t) => t !== tag);
		handleFilterChange();
	}

	// Clear all filters
	function clearFilters() {
		searchQuery = '';
		selectedGame = '';
		selectedTags = [];
		sortBy = 'newest';
		currentPage = 1;
		updateURL();
		loadMods();
	}

	// Change page
	function changePage(newPage: number) {
		if (newPage >= 1 && newPage <= totalPages && newPage !== currentPage) {
			currentPage = newPage;
			loadMods();
			window.scrollTo({ top: 0, behavior: 'smooth' });
		}
	}

	// Watch for sort changes
	$: if (sortBy) {
		handleFilterChange();
	}
</script>

<svelte:head>
	<title>Browse Mods - Azurite</title>
	<meta
		name="description"
		content="Browse and discover thousands of game modifications on Azurite. Filter by game, category, popularity and more."
	/>
</svelte:head>

<div class="min-h-screen bg-background-primary">
	<!-- Header -->
	<div class="bg-gradient-to-r from-primary-600/20 to-primary-700/20 border-b border-slate-700">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
			<div class="text-center">
				<h1 class="text-4xl font-bold text-text-primary mb-4">
					<Package class="w-10 h-10 inline-block mr-2 mb-1" />
					Browse Mods
				</h1>
				<p class="text-text-secondary text-lg max-w-2xl mx-auto mb-8">
					Discover amazing modifications for your favorite games. Filter by category, game,
					popularity and more.
				</p>

				<!-- Main Search -->
				<div class="max-w-2xl mx-auto">
					<form onsubmit={handleSearch}>
						<div class="relative">
							<Search
								class="absolute left-4 top-1/2 transform -translate-y-1/2 text-text-muted w-6 h-6"
							/>
							<input
								type="search"
								placeholder="Search for mods, authors, or descriptions..."
								bind:value={searchQuery}
								class="input pl-12 pr-4 py-4 w-full text-lg"
							/>
						</div>
					</form>
				</div>
			</div>
		</div>
	</div>

	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
		<!-- Filters and Controls -->
		<div class="flex flex-col lg:flex-row gap-4 mb-8">
			<!-- Quick Filters -->
			<div class="flex-1">
				<div class="flex flex-wrap items-center gap-3">
					<!-- Game Filter -->
					<select
						bind:value={selectedGame}
						onchange={handleFilterChange}
						class="select min-w-[150px]"
					>
						<option value="">All Games</option>
						{#each games as game (game.id)}
							<option value={game.slug}>{game.name}</option>
						{/each}
					</select>

					<!-- Sort -->
					<select bind:value={sortBy} class="select min-w-[150px]">
						{#each sortOptions as option (option.label)}
							<option value={option.value}>{option.label}</option>
						{/each}
					</select>

					<!-- Filter Toggle -->
					<button onclick={() => (showFilters = !showFilters)} class="btn btn-outline">
						<Filter class="w-4 h-4 mr-2" />
						Filters
						{#if selectedTags.length > 0}
							<span class="ml-2 bg-primary-600 text-white text-xs px-2 py-0.5 rounded-full">
								{selectedTags.length}
							</span>
						{/if}
					</button>

					{#if searchQuery || selectedGame || selectedTags.length > 0}
						<button
							onclick={clearFilters}
							class="btn btn-outline text-red-400 hover:text-red-300 hover:border-red-500"
						>
							<X class="w-4 h-4 mr-2" />
							Clear All
						</button>
					{/if}
				</div>
			</div>

			<!-- View Controls -->
			<div class="flex items-center space-x-2">
				<!-- Results Count -->
				<span class="text-text-muted text-sm whitespace-nowrap">
					{totalMods} results
				</span>

				<!-- View Mode Toggle -->
				<div class="flex border border-slate-600 rounded-lg p-1">
					<button
						onclick={() => (viewMode = 'grid')}
						class="p-2 rounded {viewMode === 'grid'
							? 'bg-primary-600 text-white'
							: 'text-text-muted hover:text-text-primary'} transition-colors"
						title="Grid View"
					>
						<Grid class="w-4 h-4" />
					</button>
					<button
						onclick={() => (viewMode = 'list')}
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
						onclick={() => (showFilters = false)}
						class="text-text-muted hover:text-text-primary"
					>
						<X class="w-5 h-5" />
					</button>
				</div>

				<!-- Tags -->
				<div class="mb-4">
					<h4 class="text-sm font-medium text-text-primary mb-3">Categories</h4>
					<div class="flex flex-wrap gap-2">
						{#each availableTags as tag (tag)}
							<button
								onclick={() => (selectedTags.includes(tag) ? removeTag(tag) : addTag(tag))}
								class="px-3 py-1 text-sm rounded-full border transition-colors {selectedTags.includes(
									tag
								)
									? 'bg-primary-600 text-white border-primary-600'
									: 'border-slate-600 text-text-secondary hover:text-text-primary hover:border-slate-500'}"
							>
								{tag}
							</button>
						{/each}
					</div>
				</div>

				<!-- Active Filters Display -->
				{#if selectedTags.length > 0}
					<div class="flex flex-wrap items-center gap-2 pt-4 border-t border-slate-700">
						<span class="text-text-secondary text-sm">Active filters:</span>
						{#each selectedTags as tag (tag)}
							<span class="badge badge-secondary">
								{tag}
								<button
									onclick={() => removeTag(tag)}
									class="ml-1 text-text-muted hover:text-text-primary"
								>
									×
								</button>
							</span>
						{/each}
					</div>
				{/if}
			</div>
		{/if}

		<!-- Results -->
		{#if isLoadingMods}
			<div class="flex justify-center items-center py-20">
				<Loading size="lg" text="Loading mods..." />
			</div>
		{:else if mods.length === 0}
			<!-- Empty State -->
			<div class="text-center py-20">
				<Package class="w-16 h-16 mx-auto text-text-muted mb-6 opacity-50" />
				<h2 class="text-2xl font-semibold text-text-primary mb-2">No mods found</h2>
				<p class="text-text-secondary mb-6 max-w-md mx-auto">
					{searchQuery || selectedGame || selectedTags.length > 0
						? 'Try adjusting your search terms or clearing some filters.'
						: 'Be the first to upload a mod to get the community started!'}
				</p>
				{#if searchQuery || selectedGame || selectedTags.length > 0}
					<button onclick={clearFilters} class="btn btn-outline"> Clear All Filters </button>
				{:else}
					<a href="/dashboard" class="btn btn-primary">
						<Package class="w-4 h-4 mr-2" />
						Upload First Mod
					</a>
				{/if}
			</div>
		{:else if viewMode === 'grid'}
			<!-- Grid View -->
			<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-3 gap-6 mb-8">
				{#each mods as mod (mod.id)}
					<ModCard {mod} variant="default" />
				{/each}
			</div>
		{:else}
			<!-- List View -->
			<div class="space-y-4 mb-8">
				{#each mods as mod (mod.id)}
					<ModCard {mod} variant="list" />
				{/each}
			</div>
		{/if}

		<!-- Pagination -->
		{#if totalPages > 1}
			<div class="flex justify-center items-center space-x-2 mb-8">
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

			<p class="text-center text-text-muted text-sm">
				Showing {(currentPage - 1) * perPage + 1}–{Math.min(currentPage * perPage, totalMods)} of {totalMods}
				results
			</p>
		{/if}
	</div>
</div>