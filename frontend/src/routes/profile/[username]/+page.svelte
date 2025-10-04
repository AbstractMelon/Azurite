<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { user as currentUser, isAuthenticated } from '$lib/stores/auth';
	import { toast } from '$lib/stores/notifications';
	import { usersApi } from '$lib/api/client';
	import Loading from '$lib/components/Loading.svelte';
	import { ModCard } from '$lib/components/cards';
	import type { User, Mod, PaginatedResponse } from '$lib/types';
	import {
		User as UserIcon,
		Calendar,
		Package,
		Download,
		Heart,
		Settings,
		Grid,
		List,
		Search
	} from 'lucide-svelte';

	// Get username from URL params
	$: username = $page.params.username;

	// Data
	let profileUser: User | null = null;
	let userMods: Mod[] = [];
	let isLoadingProfile = true;
	let isLoadingMods = true;
	let modsResponse: PaginatedResponse<Mod> | null = null;

	// UI state
	let viewMode: 'grid' | 'list' = 'grid';
	let searchQuery = '';
	let currentPage = 1;
	let isLoadingMore = false;

	// Check if this is the current user's own profile
	$: isOwnProfile = $currentUser && profileUser && $currentUser.username === profileUser.username;

	onMount(async () => {
		await loadUserProfile();
		await loadUserMods();
	});

	// Load user profile by username
	async function loadUserProfile() {
		if (!username) return;

		isLoadingProfile = true;
		try {
			const res = await usersApi.getUserByUsername(username);

			if (res.success && res.data) {
				profileUser = res.data as User;
			} else {
				toast.error('Error', res.error || 'User not found');
				goto('/404');
				return;
			}
		} catch (error) {
			console.error('Error loading user profile:', error);
			toast.error('Error', 'Failed to load user profile');
			goto('/404');
			return;
		} finally {
			isLoadingProfile = false;
		}
	}

	// Load user's mods
	async function loadUserMods(page: number = 1, append: boolean = false) {
		if (!profileUser) return;

		if (append) {
			isLoadingMore = true;
		} else {
			isLoadingMods = true;
		}

		try {
			const res = await usersApi.getUserMods(profileUser.id, {
				page,
				per_page: 20
			});

			if (res.success && res.data) {
				modsResponse = res.data as PaginatedResponse<Mod>;

				if (append) {
					userMods = [...userMods, ...modsResponse.data];
				} else {
					userMods = modsResponse.data;
				}

				currentPage = page;
			} else {
				toast.error('Error', res.error || 'Failed to load user mods');
				if (!append) userMods = [];
			}
		} catch (error) {
			console.error('Error loading user mods:', error);
			toast.error('Error', 'Failed to load user mods');
			if (!append) userMods = [];
		} finally {
			isLoadingMods = false;
			isLoadingMore = false;
		}
	}

	// Load more mods (pagination)
	async function loadMore() {
		if (!modsResponse || currentPage >= modsResponse.total_pages || isLoadingMore) return;
		await loadUserMods(currentPage + 1, true);
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
		return new Date(dateString).toLocaleDateString('en-US', {
			month: 'long',
			day: 'numeric',
			year: 'numeric'
		});
	}

	// Format relative time
	function formatRelativeTime(dateString: string): string {
		const date = new Date(dateString);
		const now = new Date();
		const diffTime = Math.abs(now.getTime() - date.getTime());
		const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

		if (diffDays === 1) return '1 day ago';
		if (diffDays < 7) return `${diffDays} days ago`;
		if (diffDays < 30) return `${Math.ceil(diffDays / 7)} weeks ago`;
		return formatDate(dateString);
	}

	// Filter mods based on search
	function filterMods(mods: Mod[]): Mod[] {
		if (!searchQuery.trim()) return mods;
		const query = searchQuery.toLowerCase();
		return mods.filter(
			(mod) =>
				mod.name.toLowerCase().includes(query) ||
				mod.short_description?.toLowerCase().includes(query) ||
				mod.game?.name.toLowerCase().includes(query)
		);
	}

	$: filteredMods = filterMods(userMods);
	$: totalDownloads = userMods.reduce((sum, mod) => sum + (mod.downloads || 0), 0);
	$: totalLikes = userMods.reduce((sum, mod) => sum + (mod.likes || 0), 0);
</script>

<svelte:head>
	<title>{profileUser?.display_name || profileUser?.username || 'User Profile'} - Azurite</title>
	<meta
		name="description"
		content="View {profileUser?.display_name ||
			profileUser?.username ||
			'user'}'s profile, mods, and activity on Azurite."
	/>
</svelte:head>

{#if isLoadingProfile}
	<div class="min-h-screen flex items-center justify-center">
		<Loading size="lg" text="Loading profile..." />
	</div>
{:else if profileUser}
	<div class="min-h-screen bg-background-primary">
		<!-- Profile Header -->
		<div class="bg-gradient-to-r from-slate-800/50 to-slate-700/50 border-b border-slate-700">
			<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
				<div
					class="flex flex-col md:flex-row items-start md:items-center space-y-6 md:space-y-0 md:space-x-8"
				>
					<!-- Profile Picture -->
					<div class="flex-shrink-0">
						{#if profileUser.avatar}
							<img
								src={profileUser.avatar}
								alt={profileUser.display_name || profileUser.username}
								class="w-24 h-24 md:w-32 md:h-32 rounded-full border-4 border-slate-600 shadow-xl"
							/>
						{:else}
							<div
								class="w-24 h-24 md:w-32 md:h-32 bg-gradient-to-br from-slate-600 to-slate-700 rounded-full flex items-center justify-center border-4 border-slate-600 shadow-xl"
							>
								<UserIcon class="w-12 h-12 md:w-16 md:h-16 text-text-muted" />
							</div>
						{/if}
					</div>

					<!-- Profile Info -->
					<div class="flex-1">
						<div class="flex items-center space-x-3 mb-2">
							<h1 class="text-3xl md:text-4xl font-bold text-text-primary">
								{profileUser.display_name || profileUser.username}
							</h1>
							{#if profileUser.role !== 'user'}
								<span
									class="badge {profileUser.role === 'admin' ? 'badge-danger' : 'badge-primary'}"
								>
									{profileUser.role.replace('_', ' ').toUpperCase()}
								</span>
							{/if}
							{#if !profileUser.is_active}
								<span class="badge badge-secondary"> INACTIVE </span>
							{/if}
						</div>

						{#if profileUser.display_name && profileUser.username !== profileUser.display_name}
							<p class="text-text-muted text-lg mb-3">@{profileUser.username}</p>
						{/if}

						{#if profileUser.bio}
							<p class="text-text-secondary text-lg mb-4 leading-relaxed max-w-2xl">
								{profileUser.bio}
							</p>
						{/if}

						<!-- Profile Stats -->
						<div class="flex flex-wrap items-center gap-6 text-text-muted mb-4">
							<div class="flex items-center">
								<Package class="w-5 h-5 mr-2" />
								<span class="font-medium">{userMods.length}</span>
								<span class="ml-1">mods</span>
							</div>
							<div class="flex items-center">
								<Download class="w-5 h-5 mr-2" />
								<span class="font-medium">{formatNumber(totalDownloads)}</span>
								<span class="ml-1">total downloads</span>
							</div>
							<div class="flex items-center">
								<Heart class="w-5 h-5 mr-2" />
								<span class="font-medium">{formatNumber(totalLikes)}</span>
								<span class="ml-1">total likes</span>
							</div>
							<div class="flex items-center">
								<Calendar class="w-5 h-5 mr-2" />
								<span>Joined {formatDate(profileUser.created_at)}</span>
							</div>
							{#if profileUser.last_login_at}
								<div class="flex items-center">
									<span class="w-2 h-2 bg-green-500 rounded-full mr-2"></span>
									<span>Last seen {formatRelativeTime(profileUser.last_login_at)}</span>
								</div>
							{/if}
						</div>
					</div>

					<!-- Actions -->
					<div class="flex flex-col sm:flex-row gap-3">
						{#if isOwnProfile}
							<a href="/settings" class="btn btn-outline">
								<Settings class="w-4 h-4 mr-2" />
								Edit Profile
							</a>
							<a href="/dashboard" class="btn btn-primary">
								<Package class="w-4 h-4 mr-2" />
								Dashboard
							</a>
						{:else if $isAuthenticated}
							<!-- Future: Add follow/message functionality -->
							<button class="btn btn-outline" disabled>
								<UserIcon class="w-4 h-4 mr-2" />
								Follow
							</button>
						{/if}
					</div>
				</div>
			</div>
		</div>

		<!-- Content -->
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
			<!-- Header with search and controls -->
			<div class="flex flex-col lg:flex-row lg:items-center justify-between gap-4 mb-8">
				<div>
					<h2 class="text-2xl font-bold text-text-primary mb-2">
						{isOwnProfile
							? 'Your Mods'
							: `${profileUser.display_name || profileUser.username}'s Mods`}
					</h2>
					<p class="text-text-muted">
						{modsResponse?.total || userMods.length} mod{userMods.length !== 1 ? 's' : ''} published
					</p>
				</div>

				<!-- Controls -->
				<div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-4">
					<!-- Search -->
					<div class="relative">
						<Search
							class="absolute left-3 top-1/2 transform -translate-y-1/2 text-text-muted w-4 h-4"
						/>
						<input
							type="search"
							placeholder="Search mods..."
							bind:value={searchQuery}
							class="input pl-9 pr-4 py-2 w-full sm:w-64"
						/>
					</div>

					<!-- View Toggle -->
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

			<!-- Mods Content -->
			{#if isLoadingMods}
				<div class="flex justify-center items-center py-20">
					<Loading size="lg" text="Loading mods..." />
				</div>
			{:else if filteredMods.length === 0}
				<!-- Empty State -->
				<div class="text-center py-20">
					<Package class="w-16 h-16 mx-auto text-text-muted mb-6 opacity-50" />
					<h2 class="text-2xl font-semibold text-text-primary mb-2">
						{searchQuery
							? 'No mods match your search'
							: isOwnProfile
								? 'No mods yet'
								: 'No mods published'}
					</h2>
					<p class="text-text-secondary mb-6 max-w-md mx-auto">
						{searchQuery
							? 'Try adjusting your search terms.'
							: isOwnProfile
								? 'Start creating and sharing your mods with the community!'
								: `${profileUser.display_name || profileUser.username} hasn't published any mods yet.`}
					</p>
					{#if !searchQuery && isOwnProfile}
						<a href="/dashboard" class="btn btn-primary">
							<Package class="w-4 h-4 mr-2" />
							Upload Your First Mod
						</a>
					{:else if !searchQuery}
						<a href="/browse" class="btn btn-primary">
							<Search class="w-4 h-4 mr-2" />
							Browse Other Mods
						</a>
					{/if}
				</div>
			{:else if viewMode === 'grid'}
				<!-- Grid View -->
				<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
					{#each filteredMods as mod (mod.id)}
						<ModCard {mod} variant="default" />
					{/each}
				</div>

				<!-- Load More Button -->
				{#if modsResponse && currentPage < modsResponse.total_pages && !searchQuery}
					<div class="text-center">
						<button onclick={loadMore} disabled={isLoadingMore} class="btn btn-outline">
							{#if isLoadingMore}
								<Loading size="sm" inline />
								Loading...
							{:else}
								Load More
							{/if}
						</button>
					</div>
				{/if}
			{:else}
				<!-- List View -->
				<div class="space-y-4 mb-8">
					{#each filteredMods as mod (mod.id)}
						<ModCard {mod} variant="list" />
					{/each}
				</div>

				<!-- Load More Button -->
				{#if modsResponse && currentPage < modsResponse.total_pages && !searchQuery}
					<div class="text-center">
						<button onclick={loadMore} disabled={isLoadingMore} class="btn btn-outline">
							{#if isLoadingMore}
								<Loading size="sm" inline />
								Loading...
							{:else}
								Load More
							{/if}
						</button>
					</div>
				{/if}
			{/if}
		</div>
	</div>
{:else}
	<!-- User not found -->
	<div class="min-h-screen flex items-center justify-center">
		<div class="text-center py-20">
			<UserIcon class="w-16 h-16 mx-auto text-text-muted mb-6 opacity-50" />
			<h1 class="text-2xl font-semibold text-text-primary mb-2">User not found</h1>
			<p class="text-text-secondary mb-6">
				The user you're looking for doesn't exist or has been deactivated.
			</p>
			<a href="/browse" class="btn btn-primary">
				<Search class="w-4 h-4 mr-2" />
				Browse Community
			</a>
		</div>
	</div>
{/if}
