<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { user, isAuthenticated } from '$lib/stores/auth';
	import { toast } from '$lib/stores/notifications';
	import { authApi } from '$lib/api/client';
	import Loading from '$lib/components/Loading.svelte';
	import { ModCard } from '$lib/components/cards';
	import {
		User,
		Calendar,
		Package,
		Download,
		Heart,
		Settings,
		Edit,
		Shield,
		Grid,
		List,
		Search
	} from 'lucide-svelte';

	// Data
	let profileUser: any = null;
	let userMods: any[] = [];
	let likedMods: any[] = [];
	let isLoadingProfile = true;
	let isLoadingMods = true;

	// UI state
	let activeTab = 'mods'; // mods, liked
	let viewMode = 'grid'; // grid, list
	let searchQuery = '';

	onMount(async () => {
		if (!$isAuthenticated) {
			goto('/auth/login?redirect=/profile');
			return;
		}

		profileUser = $user;
		isLoadingProfile = false;
		await loadUserMods();
	});

	// Load user's mods
	async function loadUserMods() {
		if (!profileUser) return;

		isLoadingMods = true;
		try {
			const res = await authApi.getUserMods({ page: 1, per_page: 20 });

			if (res.success && res.data) {
				userMods = Array.isArray(res.data) ? res.data : res.data.data || [];
			} else {
				toast.error('Error', res.error || 'Failed to load your mods');
				userMods = [];
			}

			// We need a real endpoint for this
			likedMods = [];
		} catch (error) {
			console.error('Error loading user mods:', error);
			toast.error('Error', 'Failed to load your mods');
		} finally {
			isLoadingMods = false;
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
		return new Date(dateString).toLocaleDateString();
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
	function filterMods(mods: any[]): any[] {
		if (!searchQuery.trim()) return mods;
		const query = searchQuery.toLowerCase();
		return mods.filter(
			(mod) =>
				mod.name.toLowerCase().includes(query) ||
				mod.short_description?.toLowerCase().includes(query)
		);
	}

	// Get status badge classes
	function getStatusBadge(status: string) {
		switch (status) {
			case 'approved':
				return 'badge-success';
			case 'pending':
				return 'badge-warning';
			case 'rejected':
				return 'badge-danger';
			default:
				return 'badge-secondary';
		}
	}

	// Get status text
	function getStatusText(status: string) {
		switch (status) {
			case 'approved':
				return 'Published';
			case 'pending':
				return 'Under Review';
			case 'rejected':
				return 'Rejected';
			default:
				return 'Unknown';
		}
	}

	$: filteredUserMods = filterMods(userMods);
	$: filteredLikedMods = filterMods(likedMods);
	$: currentMods = activeTab === 'mods' ? filteredUserMods : filteredLikedMods;
</script>

<svelte:head>
	<title>{profileUser?.display_name || profileUser?.username || 'Profile'} - Azurite</title>
	<meta name="description" content="View and manage your Azurite profile, mods, and activity." />
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
								alt={profileUser.display_name}
								class="w-24 h-24 md:w-32 md:h-32 rounded-full border-4 border-slate-600 shadow-xl"
							/>
						{:else}
							<div
								class="w-24 h-24 md:w-32 md:h-32 bg-gradient-to-br from-slate-600 to-slate-700 rounded-full flex items-center justify-center border-4 border-slate-600 shadow-xl"
							>
								<User class="w-12 h-12 md:w-16 md:h-16 text-text-muted" />
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
								<span class="font-medium"
									>{formatNumber(
										userMods.reduce((sum, mod) => sum + (mod.downloads || 0), 0)
									)}</span
								>
								<span class="ml-1">total downloads</span>
							</div>
							<div class="flex items-center">
								<Heart class="w-5 h-5 mr-2" />
								<span class="font-medium"
									>{formatNumber(userMods.reduce((sum, mod) => sum + (mod.likes || 0), 0))}</span
								>
								<span class="ml-1">total likes</span>
							</div>
							<div class="flex items-center">
								<Calendar class="w-5 h-5 mr-2" />
								<span>Joined {formatDate(profileUser.created_at)}</span>
							</div>
						</div>
					</div>

					<!-- Actions -->
					<div class="flex flex-col sm:flex-row gap-3">
						<a href="/settings" class="btn btn-outline">
							<Settings class="w-4 h-4 mr-2" />
							Edit Profile
						</a>
						<a href="/dashboard" class="btn btn-primary">
							<Package class="w-4 h-4 mr-2" />
							Dashboard
						</a>
					</div>
				</div>
			</div>
		</div>

		<!-- Content -->
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
			<!-- Tabs -->
			<div class="flex items-center justify-between mb-8">
				<div class="flex space-x-1 bg-slate-800 rounded-lg p-1">
					<button
						on:click={() => (activeTab = 'mods')}
						class="px-4 py-2 rounded-md text-sm font-medium transition-colors {activeTab === 'mods'
							? 'bg-primary-600 text-white'
							: 'text-text-secondary hover:text-text-primary'}"
					>
						<Package class="w-4 h-4 mr-2 inline" />
						My Mods ({userMods.length})
					</button>
					<button
						on:click={() => (activeTab = 'liked')}
						class="px-4 py-2 rounded-md text-sm font-medium transition-colors {activeTab === 'liked'
							? 'bg-primary-600 text-white'
							: 'text-text-secondary hover:text-text-primary'}"
					>
						<Heart class="w-4 h-4 mr-2 inline" />
						Liked ({likedMods.length})
					</button>
				</div>

				<!-- Controls -->
				<div class="flex items-center space-x-4">
					<!-- Search -->
					<div class="relative">
						<Search
							class="absolute left-3 top-1/2 transform -translate-y-1/2 text-text-muted w-4 h-4"
						/>
						<input
							type="search"
							placeholder="Search mods..."
							bind:value={searchQuery}
							class="input pl-9 pr-4 py-2 w-64"
						/>
					</div>

					<!-- View Toggle -->
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

			<!-- Mods Content -->
			{#if isLoadingMods}
				<div class="flex justify-center items-center py-20">
					<Loading size="lg" text="Loading mods..." />
				</div>
			{:else if currentMods.length === 0}
				<!-- Empty State -->
				<div class="text-center py-20">
					<Package class="w-16 h-16 mx-auto text-text-muted mb-6 opacity-50" />
					<h2 class="text-2xl font-semibold text-text-primary mb-2">
						{activeTab === 'mods'
							? searchQuery
								? 'No mods match your search'
								: 'No mods yet'
							: searchQuery
								? 'No liked mods match your search'
								: 'No liked mods yet'}
					</h2>
					<p class="text-text-secondary mb-6 max-w-md mx-auto">
						{activeTab === 'mods'
							? searchQuery
								? 'Try adjusting your search terms.'
								: 'Start creating and sharing your mods with the community!'
							: searchQuery
								? 'Try adjusting your search terms.'
								: 'Browse mods and like your favorites to see them here.'}
					</p>
					{#if !searchQuery}
						{#if activeTab === 'mods'}
							<a href="/dashboard" class="btn btn-primary">
								<Package class="w-4 h-4 mr-2" />
								Upload Your First Mod
							</a>
						{:else}
							<a href="/browse" class="btn btn-primary">
								<Search class="w-4 h-4 mr-2" />
								Browse Mods
							</a>
						{/if}
					{/if}
				</div>
			{:else if viewMode === 'grid'}
				<!-- Grid View -->
				<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
					{#each currentMods as mod (mod.id)}
						<div class="card card-hover group relative">
							<div class="p-4">
								<!-- Status Badge (for user's mods) -->
								{#if activeTab === 'mods' && mod.status}
									<div class="absolute top-2 right-2">
										<span class="badge {getStatusBadge(mod.status)} text-xs">
											{getStatusText(mod.status)}
										</span>
									</div>
								{/if}

								<!-- Mod Icon -->
								<div class="flex justify-center mb-3">
									{#if mod.icon}
										<img
											src={mod.icon}
											alt={mod.name}
											class="w-12 h-12 rounded-lg border border-slate-600"
										/>
									{:else}
										<div class="w-12 h-12 bg-slate-600 rounded-lg flex items-center justify-center">
											<Package class="w-6 h-6 text-text-muted" />
										</div>
									{/if}
								</div>

								<!-- Mod Info -->
								<div class="text-center">
									<h3
										class="text-lg font-semibold text-text-primary group-hover:text-primary-400 transition-colors mb-1 line-clamp-1"
									>
										<a href="/mods/{mod.game?.slug}/{mod.slug}">{mod.name}</a>
									</h3>

									<p class="text-text-muted text-sm mb-2">
										for {mod.game?.name}
									</p>

									{#if mod.short_description}
										<p class="text-text-secondary text-sm mb-3 line-clamp-2">
											{mod.short_description}
										</p>
									{/if}

									<!-- Stats -->
									<div
										class="flex items-center justify-center space-x-4 text-xs text-text-muted mb-2"
									>
										<div class="flex items-center">
											<Download class="w-3 h-3 mr-1" />
											{formatNumber(mod.downloads || 0)}
										</div>
										<div class="flex items-center">
											<Heart class="w-3 h-3 mr-1" />
											{formatNumber(mod.likes || 0)}
										</div>
									</div>

									<!-- Update Date -->
									<p class="text-xs text-text-muted">
										Updated {formatRelativeTime(mod.updated_at)}
									</p>

									<!-- Author (for liked mods) -->
									{#if activeTab === 'liked' && mod.owner}
										<p class="text-xs text-text-muted mt-1">
											by {mod.owner.display_name || mod.owner.username}
										</p>
									{/if}
								</div>
							</div>
						</div>
					{/each}
				</div>
			{:else}
				<!-- List View -->
				<div class="space-y-4">
					{#each currentMods as mod (mod.id)}
						<div class="card card-hover group relative">
							<div class="p-4 flex items-start space-x-4">
								<!-- Large Mod Image -->
								<div class="flex-shrink-0">
									{#if mod.icon}
										<img
											src={mod.icon}
											alt={mod.name}
											class="w-20 h-20 rounded-xl object-cover shadow-md"
										/>
									{:else}
										<div
											class="w-20 h-20 bg-gradient-to-br from-primary-600 to-primary-700 rounded-xl flex items-center justify-center shadow-md"
										>
											<Package class="w-10 h-10 text-white" />
										</div>
									{/if}
								</div>

								<!-- Content -->
								<div class="flex-grow min-w-0">
									<div class="flex items-start justify-between mb-2">
										<div class="min-w-0 flex-1">
											<div class="flex items-center space-x-2 mb-1">
												<h3
													class="text-lg font-bold text-text-primary group-hover:text-primary-400 transition-colors truncate"
												>
													<a href="/mods/{mod.game?.slug}/{mod.slug}">{mod.name}</a>
												</h3>
												{#if activeTab === 'mods' && mod.status}
													<span class="badge {getStatusBadge(mod.status)} text-xs">
														{getStatusText(mod.status)}
													</span>
												{/if}
											</div>

											<div class="flex items-center space-x-2 mb-1">
												<span class="text-text-muted text-sm">for {mod.game?.name}</span>
												{#if activeTab === 'liked' && mod.owner}
													<span class="text-text-muted text-sm"
														>• by {mod.owner.display_name || mod.owner.username}</span
													>
												{/if}
												{#if mod.version}
													<span
														class="px-2 py-0.5 bg-slate-700/50 text-text-muted text-xs rounded-full"
														>v{mod.version}</span
													>
												{/if}
											</div>
										</div>

										<!-- Actions -->
										{#if activeTab === 'mods'}
											<div class="flex items-center space-x-2 flex-shrink-0">
												<a
													href="/dashboard/mods/{mod.id}/edit"
													class="btn btn-outline btn-sm"
													title="Edit Mod"
												>
													<Edit class="w-4 h-4" />
												</a>
											</div>
										{/if}
									</div>

									{#if mod.short_description}
										<p class="text-text-secondary text-sm line-clamp-2 mb-3 leading-relaxed">
											{mod.short_description}
										</p>
									{/if}

									<div class="flex items-center justify-between">
										<div class="flex items-center space-x-4 text-sm text-text-muted">
											<div class="flex items-center">
												<Download class="w-4 h-4 mr-1" />
												<span class="font-medium">{formatNumber(mod.downloads || 0)}</span>
											</div>
											<div class="flex items-center">
												<Heart class="w-4 h-4 mr-1" />
												<span class="font-medium">{formatNumber(mod.likes || 0)}</span>
											</div>
											<div class="flex items-center">
												<Calendar class="w-4 h-4 mr-1" />
												<span>Updated {formatRelativeTime(mod.updated_at)}</span>
											</div>
										</div>
										{#if mod.is_scanned}
											<div
												class="flex items-center px-2 py-1 bg-green-500/10 text-green-400 text-xs font-medium rounded-full border border-green-500/20"
											>
												<Shield class="w-3 h-3 mr-1" />
												Verified
											</div>
										{/if}
									</div>
								</div>
							</div>
						</div>
					{/each}
				</div>
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
