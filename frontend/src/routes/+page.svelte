<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { gamesApi } from '$lib/api/client';
	import Loading from '$lib/components/Loading.svelte';
	import { GameCard, ModCard } from '$lib/components/cards';
	import type { Game, Mod } from '$lib/types';
	import {
		Users,
		TrendingUp,
		ArrowRight,
		Gamepad2,
		Package,
		Shield,
		Zap,
		Globe
	} from 'lucide-svelte';

	let isLoading = true;
	let featuredGames: Game[] = [];
	let popularMods: Mod[] = [];
	let recentMods: Mod[] = [];
	let stats = {
		totalMods: 0,
		totalDownloads: 0,
		totalUsers: 0,
		totalGames: 0
	};

	// Load homepage data
	onMount(async () => {
		try {
			// Load featured games
			const gamesResponse = await gamesApi.getGames({ per_page: 6 });
			if (gamesResponse.success) {
				featuredGames = gamesResponse.data?.data || [];
			}

			// Mock popular mods (should be replaced with actual API call)
			popularMods = [
				{
					id: 1,
					name: 'OptiFine',
					slug: 'optifine',
					game: { name: 'Minecraft', slug: 'minecraft' },
					downloads: 125000,
					likes: 8500,
					short_description: 'Boost FPS and add advanced graphics features',
					icon: '/api/files/images/optifine-icon.png'
				},
				{
					id: 2,
					name: 'JEI',
					slug: 'jei',
					game: { name: 'Minecraft', slug: 'minecraft' },
					downloads: 98000,
					likes: 6200,
					short_description: 'Just Enough Items - Recipe and item viewer',
					icon: '/api/files/images/jei-icon.png'
				},
				{
					id: 3,
					name: 'Biomes O Plenty',
					slug: 'biomes-o-plenty',
					game: { name: 'Minecraft', slug: 'minecraft' },
					downloads: 87000,
					likes: 5400,
					short_description: 'Adds 80+ new biomes to explore',
					icon: '/api/files/images/bop-icon.png'
				}
			];

			// Mock recent mods
			recentMods = [
				{
					id: 4,
					name: 'Create Mod',
					slug: 'create',
					game: { name: 'Minecraft', slug: 'minecraft' },
					downloads: 45000,
					likes: 3200,
					short_description: 'Contraptions and automation',
					icon: '/api/files/images/create-icon.png',
					created_at: new Date(Date.now() - 86400000).toISOString()
				},
				{
					id: 5,
					name: 'Twilight Forest',
					slug: 'twilight-forest',
					game: { name: 'Minecraft', slug: 'minecraft' },
					downloads: 32000,
					likes: 2100,
					short_description: 'Explore a mysterious twilight dimension',
					icon: '/api/files/images/tf-icon.png',
					created_at: new Date(Date.now() - 172800000).toISOString()
				}
			];

			// Mock stats
			stats = {
				totalMods: 15420,
				totalDownloads: 2500000,
				totalUsers: 45000,
				totalGames: 12
			};
		} catch (error) {
			console.error('Failed to load homepage data:', error);
		} finally {
			isLoading = false;
		}
	});

	function formatNumber(num: number): string {
		if (num >= 1000000) {
			return (num / 1000000).toFixed(1) + 'M';
		}
		if (num >= 1000) {
			return (num / 1000).toFixed(1) + 'K';
		}
		return num.toString();
	}
</script>

<svelte:head>
	<title>Azurite - Discover, Share, and Manage Game Mods</title>
	<meta
		name="description"
		content="The ultimate platform for discovering, sharing, and managing game modifications. Join thousands of modders and gamers in the Azurite community."
	/>
</svelte:head>

{#if isLoading}
	<div class="min-h-[60vh] flex items-center justify-center">
		<Loading size="lg" text="Loading Azurite..." />
	</div>
{:else}
	<!-- Hero Section -->
	<section
		class="relative overflow-hidden bg-gradient-to-br from-background-primary via-background-secondary to-slate-800 py-20 lg:py-32"
	>
		<!-- Background Elements -->
		<div class="absolute inset-0">
			<div
				class="absolute top-1/4 left-1/4 w-64 h-64 bg-primary-500/10 rounded-full blur-3xl"
			></div>
			<div
				class="absolute bottom-1/4 right-1/4 w-48 h-48 bg-blue-500/10 rounded-full blur-3xl"
			></div>
		</div>

		<div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="text-center">
				<!-- Main Hero Content -->
				<h1 class="text-4xl sm:text-5xl lg:text-6xl font-bold mb-6">
					<span class="text-text-primary">Welcome to</span>
					<span class="text-gradient block mt-2">Azurite</span>
				</h1>

				<p class="text-xl text-text-secondary max-w-3xl mx-auto mb-8 leading-relaxed">
					The ultimate platform for discovering, sharing, and managing game modifications. Join
					thousands of modders and gamers in our thriving community.
				</p>

				<!-- CTA Buttons -->
				<div class="flex flex-col sm:flex-row gap-4 justify-center mb-12">
					<button on:click={() => goto('/games')} class="btn btn-primary btn-lg">
						<Gamepad2 class="w-5 h-5 mr-2" />
						Browse Games
					</button>
					<button on:click={() => goto('/browse')} class="btn btn-outline btn-lg">
						<Package class="w-5 h-5 mr-2" />
						Discover Mods
					</button>
				</div>

				<!-- Stats -->
				<div class="grid grid-cols-2 md:grid-cols-4 gap-6 max-w-4xl mx-auto">
					<div class="text-center p-4">
						<div class="text-2xl lg:text-3xl font-bold text-primary-400 mb-1">
							{formatNumber(stats.totalMods)}
						</div>
						<div class="text-text-muted text-sm">Total Mods</div>
					</div>
					<div class="text-center p-4">
						<div class="text-2xl lg:text-3xl font-bold text-primary-400 mb-1">
							{formatNumber(stats.totalDownloads)}
						</div>
						<div class="text-text-muted text-sm">Downloads</div>
					</div>
					<div class="text-center p-4">
						<div class="text-2xl lg:text-3xl font-bold text-primary-400 mb-1">
							{formatNumber(stats.totalUsers)}
						</div>
						<div class="text-text-muted text-sm">Users</div>
					</div>
					<div class="text-center p-4">
						<div class="text-2xl lg:text-3xl font-bold text-primary-400 mb-1">
							{stats.totalGames}
						</div>
						<div class="text-text-muted text-sm">Games</div>
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- Features Section -->
	<section class="py-16 bg-background-secondary">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="text-center mb-12">
				<h2 class="text-3xl font-bold text-text-primary mb-4">Why Choose Azurite?</h2>
				<p class="text-text-secondary max-w-2xl mx-auto">
					Built by modders, for modders. Experience the next generation of mod discovery and
					sharing.
				</p>
			</div>

			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
				<div class="text-center p-6">
					<div
						class="w-12 h-12 bg-primary-500/20 rounded-lg flex items-center justify-center mx-auto mb-4"
					>
						<Zap class="w-6 h-6 text-primary-400" />
					</div>
					<h3 class="text-lg font-semibold text-text-primary mb-2">Lightning Fast</h3>
					<p class="text-text-secondary">
						Optimized for speed with instant search, fast downloads, and responsive design.
					</p>
				</div>

				<div class="text-center p-6">
					<div
						class="w-12 h-12 bg-primary-500/20 rounded-lg flex items-center justify-center mx-auto mb-4"
					>
						<Shield class="w-6 h-6 text-primary-400" />
					</div>
					<h3 class="text-lg font-semibold text-text-primary mb-2">Safe & Secure</h3>
					<p class="text-text-secondary">
						All mods are scanned for malware and reviewed by our community moderators.
					</p>
				</div>

				<div class="text-center p-6">
					<div
						class="w-12 h-12 bg-primary-500/20 rounded-lg flex items-center justify-center mx-auto mb-4"
					>
						<Globe class="w-6 h-6 text-primary-400" />
					</div>
					<h3 class="text-lg font-semibold text-text-primary mb-2">Multi-Game Support</h3>
					<p class="text-text-secondary">
						Supporting multiple games with dedicated communities and documentation.
					</p>
				</div>
			</div>
		</div>
	</section>

	<!-- Featured Games -->
	{#if featuredGames.length > 0}
		<section class="py-16">
			<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
				<div class="flex items-center justify-between mb-8">
					<div>
						<h2 class="text-2xl font-bold text-text-primary">Featured Games</h2>
						<p class="text-text-secondary mt-2">Explore mods for your favorite games</p>
					</div>
					<a href="/games" class="btn btn-outline">
						View All
						<ArrowRight class="w-4 h-4 ml-2" />
					</a>
				</div>

				<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
					{#each featuredGames.slice(0, 6) as game (game.id)}
						<GameCard {game} variant="default" />
					{/each}
				</div>
			</div>
		</section>
	{/if}

	<!-- Popular Mods -->
	{#if popularMods.length > 0}
		<section class="py-16 bg-background-secondary">
			<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
				<div class="flex items-center justify-between mb-8">
					<div>
						<h2 class="text-2xl font-bold text-text-primary">Popular Mods</h2>
						<p class="text-text-secondary mt-2">Most downloaded and loved by the community</p>
					</div>
					<a href="/browse?sort=popular" class="btn btn-outline">
						<TrendingUp class="w-4 h-4 mr-2" />
						View All
					</a>
				</div>

				<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
					{#each popularMods as mod (mod.id)}
						<ModCard {mod} variant="default" />
					{/each}
				</div>
			</div>
		</section>
	{/if}

	<!-- Recent Mods -->
	{#if recentMods.length > 0}
		<section class="py-16">
			<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
				<div class="flex items-center justify-between mb-8">
					<div>
						<h2 class="text-2xl font-bold text-text-primary">Recently Added</h2>
						<p class="text-text-secondary mt-2">Fresh mods from our creators</p>
					</div>
					<a href="/browse?sort=newest" class="btn btn-outline">
						View All
						<ArrowRight class="w-4 h-4 ml-2" />
					</a>
				</div>

				<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
					{#each recentMods as mod (mod.id)}
						<ModCard {mod} variant="list" />
					{/each}
				</div>
			</div>
		</section>
	{/if}

	<!-- CTA Section -->
	<section class="py-20 bg-gradient-to-r from-primary-600 to-primary-700">
		<div class="max-w-4xl mx-auto text-center px-4 sm:px-6 lg:px-8">
			<h2 class="text-3xl font-bold text-white mb-4">Ready to Share Your Creations?</h2>
			<p class="text-primary-100 text-lg mb-8 max-w-2xl mx-auto">
				Join our community of creators and share your mods with thousands of players worldwide.
				Upload, manage, and grow your modding projects with ease.
			</p>
			<div class="flex flex-col sm:flex-row gap-4 justify-center">
				<a href="/auth/register" class="btn bg-white text-primary-600 hover:bg-gray-100 btn-lg">
					<Users class="w-5 h-5 mr-2" />
					Join Community
				</a>
				<a href="/dashboard" class="btn border-2 border-white text-white hover:bg-white/10 btn-lg">
					<Package class="w-5 h-5 mr-2" />
					Upload Mod
				</a>
			</div>
		</div>
	</section>
{/if}

<style>
	.line-clamp-2 {
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
</style>
