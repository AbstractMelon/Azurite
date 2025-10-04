<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { user, isAuthenticated } from '$lib/stores/auth';
	import { modsApi, authApi, notificationsApi } from '$lib/api/client';
	import { toast } from '$lib/stores/notifications';
	import Loading from '$lib/components/Loading.svelte';
	import type { Mod, ActivityItem, DashboardStats } from '$lib/types';
	import { SvelteSet } from 'svelte/reactivity';
	import {
		Package,
		Plus,
		Download,
		Heart,
		Eye,
		Clock,
		CheckCircle,
		XCircle,
		Edit,
		Trash2,
		Settings,
		MessageCircle,
		User
	} from 'lucide-svelte';

	// Data
	let userMods: Mod[] = [];
	let stats: DashboardStats = {
		totalMods: 0,
		totalDownloads: 0,
		totalLikes: 0,
		totalComments: 0,
		pendingMods: 0,
		approvedMods: 0,
		rejectedMods: 0
	};
	let recentActivity: ActivityItem[] = [];
	let isLoadingMods = true;

	// UI state
	let selectedMods = new SvelteSet<number>();
	let showBulkActions = false;

	onMount(async () => {
		if (!$isAuthenticated) {
			goto('/auth/login?redirect=/dashboard');
			return;
		}

		await Promise.all([loadUserMods(), loadStats(), loadRecentActivity()]);
	});

	// Load user's mods
	async function loadUserMods() {
		isLoadingMods = true;
		try {
			const response = await authApi.getUserMods({ per_page: 100 });
			if (response.success && response.data) {
				userMods = response.data.data || [];
			} else {
				toast.error('Error', response.error || 'Failed to load your mods');
				userMods = [];
			}
		} catch (error) {
			console.error('Error loading user mods:', error);
			toast.error('Error', 'Failed to load your mods');
			userMods = [];
		} finally {
			isLoadingMods = false;
		}
	}

	// Helper function to get mod status
	function getModStatus(mod: Mod): string {
		if (mod.is_rejected) return 'rejected';
		if (mod.is_scanned) return 'approved';
		return 'pending';
	}

	// Load dashboard stats
	async function loadStats() {
		try {
			// Calculate stats from mods
			stats = {
				totalMods: userMods.length,
				totalDownloads: userMods.reduce((sum, mod) => sum + (mod.downloads || 0), 0),
				totalLikes: userMods.reduce((sum, mod) => sum + (mod.likes || 0), 0),
				totalComments: userMods.reduce((sum, mod) => sum + (mod.comments || 0), 0),
				pendingMods: userMods.filter((mod) => getModStatus(mod) === 'pending').length,
				approvedMods: userMods.filter((mod) => getModStatus(mod) === 'approved').length,
				rejectedMods: userMods.filter((mod) => getModStatus(mod) === 'rejected').length
			};
		} catch (error) {
			console.error('Error loading stats:', error);
		}
	}

	// Load recent activity
	async function loadRecentActivity() {
		try {
			// Load recent notifications for activity feed
			const response = await notificationsApi.getNotifications({ per_page: 10 });
			if (response.success && response.data) {
				recentActivity = (response.data.data || []).map((notification: unknown) => ({
					type: notification.type || 'info',
					message: notification.title || notification.message,
					timestamp: notification.created_at,
					mod: notification.data ? JSON.parse(notification.data) : null
				}));
			} else {
				// Fallback to empty array if API fails
				recentActivity = [];
			}
		} catch (error) {
			console.error('Error loading recent activity:', error);
			recentActivity = [];
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
	function formatRelativeTime(dateString: string): string {
		const date = new Date(dateString);
		const now = new Date();
		const diffTime = Math.abs(now.getTime() - date.getTime());
		const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

		if (diffDays === 1) return '1 day ago';
		if (diffDays < 7) return `${diffDays} days ago`;
		if (diffDays < 30) return `${Math.ceil(diffDays / 7)} weeks ago`;
		return date.toLocaleDateString();
	}

	// Get status badge class
	function getStatusBadge(mod: any): string {
		const status = getModStatus(mod);
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
	function getStatusText(mod: any): string {
		const status = getModStatus(mod);
		switch (status) {
			case 'approved':
				return 'Approved';
			case 'pending':
				return 'Pending Review';
			case 'rejected':
				return 'Rejected';
			default:
				return 'Unknown';
		}
	}

	// Handle mod selection
	function toggleModSelection(modId: number) {
		if (selectedMods.has(modId)) {
			selectedMods.delete(modId);
		} else {
			selectedMods.add(modId);
		}
		selectedMods = selectedMods; // Trigger reactivity
		showBulkActions = selectedMods.size > 0;
	}

	// Select all mods
	function selectAllMods() {
		selectedMods = new SvelteSet(userMods.map((mod) => mod.id));
		showBulkActions = true;
	}

	// Clear selection
	function clearSelection() {
		selectedMods.clear();
		selectedMods = selectedMods;
		showBulkActions = false;
	}

	// Delete mod
	async function deleteMod(modId: number, modName: string) {
		if (!confirm(`Are you sure you want to delete "${modName}"? This action cannot be undone.`)) {
			return;
		}

		try {
			const response = await modsApi.deleteMod(modId);
			if (response.success) {
				toast.success('Mod deleted', `"${modName}" has been deleted successfully.`);
				await loadUserMods();
				await loadStats();
			} else {
				toast.error('Failed to delete mod', response.error);
			}
		} catch (error) {
			console.error('Error deleting mod:', error);
			toast.error('Error', 'Failed to delete mod');
		}
	}

	// Bulk delete
	async function bulkDelete() {
		const selectedModsArray = Array.from(selectedMods);
		const modNames = selectedModsArray
			.map((id) => userMods.find((m) => m.id === id)?.name)
			.filter(Boolean);

		if (
			!confirm(
				`Are you sure you want to delete ${selectedModsArray.length} mods? This action cannot be undone.\n\nMods to delete:\n${modNames.join('\n')}`
			)
		) {
			return;
		}

		try {
			const deletePromises = selectedModsArray.map((id) => modsApi.deleteMod(id));
			await Promise.all(deletePromises);

			toast.success(
				'Mods deleted',
				`${selectedModsArray.length} mods have been deleted successfully.`
			);
			clearSelection();
			await loadUserMods();
			await loadStats();
		} catch (error) {
			console.error('Error bulk deleting mods:', error);
			toast.error('Error', 'Failed to delete some mods');
		}
	}

	// Navigate to mod edit
	function editMod(modId: number) {
		goto(`/dashboard/mods/${modId}/edit`);
	}

	// Navigate to create mod
	function createMod() {
		goto('/dashboard/mods/create');
	}

	// Get activity icon
	function getActivityIcon(type: string) {
		switch (type) {
			case 'comment':
				return MessageCircle;
			case 'like':
				return Heart;
			case 'download':
				return Download;
			default:
				return Eye;
		}
	}

	$: selectedCount = selectedMods.size;
</script>

<svelte:head>
	<title>Creator Dashboard - Azurite</title>
	<meta
		name="description"
		content="Manage your mods, track performance, and grow your creator profile on Azurite."
	/>
</svelte:head>

<div class="min-h-screen bg-background-primary">
	<!-- Header -->
	<div class="bg-gradient-to-r from-primary-600/20 to-primary-700/20 border-b border-slate-700">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
			<div
				class="flex flex-col sm:flex-row justify-between items-start sm:items-center space-y-4 sm:space-y-0"
			>
				<div>
					<h1 class="text-3xl font-bold text-text-primary">Creator Dashboard</h1>
					<p class="text-text-secondary mt-1">
						Welcome back, {$user?.display_name || $user?.username}! Here's how your mods are
						performing.
					</p>
				</div>
				<button on:click={createMod} class="btn btn-primary">
					<Plus class="w-5 h-5 mr-2" />
					Upload New Mod
				</button>
			</div>
		</div>
	</div>

	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
		<!-- Stats Grid -->
		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
			<div class="card">
				<div class="p-6">
					<div class="flex items-center">
						<div class="p-2 bg-primary-500/20 rounded-lg">
							<Package class="w-6 h-6 text-primary-400" />
						</div>
						<div class="ml-4">
							<p class="text-text-muted text-sm">Total Mods</p>
							<p class="text-2xl font-bold text-text-primary">{stats.totalMods}</p>
						</div>
					</div>
				</div>
			</div>

			<div class="card">
				<div class="p-6">
					<div class="flex items-center">
						<div class="p-2 bg-green-500/20 rounded-lg">
							<Download class="w-6 h-6 text-green-400" />
						</div>
						<div class="ml-4">
							<p class="text-text-muted text-sm">Total Downloads</p>
							<p class="text-2xl font-bold text-text-primary">
								{formatNumber(stats.totalDownloads)}
							</p>
						</div>
					</div>
				</div>
			</div>

			<div class="card">
				<div class="p-6">
					<div class="flex items-center">
						<div class="p-2 bg-red-500/20 rounded-lg">
							<Heart class="w-6 h-6 text-red-400" />
						</div>
						<div class="ml-4">
							<p class="text-text-muted text-sm">Total Likes</p>
							<p class="text-2xl font-bold text-text-primary">{formatNumber(stats.totalLikes)}</p>
						</div>
					</div>
				</div>
			</div>

			<div class="card">
				<div class="p-6">
					<div class="flex items-center">
						<div class="p-2 bg-blue-500/20 rounded-lg">
							<MessageCircle class="w-6 h-6 text-blue-400" />
						</div>
						<div class="ml-4">
							<p class="text-text-muted text-sm">Total Comments</p>
							<p class="text-2xl font-bold text-text-primary">
								{formatNumber(stats.totalComments)}
							</p>
						</div>
					</div>
				</div>
			</div>
		</div>

		<!-- Status Overview -->
		{#if stats.pendingMods > 0 || stats.rejectedMods > 0}
			<div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
				<div class="card bg-green-900/20 border-green-600/50">
					<div class="p-4">
						<div class="flex items-center">
							<CheckCircle class="w-5 h-5 text-green-400 mr-2" />
							<span class="text-green-400 font-medium">Published</span>
						</div>
						<p class="text-2xl font-bold text-green-300 mt-2">{stats.approvedMods}</p>
					</div>
				</div>

				{#if stats.pendingMods > 0}
					<div class="card bg-yellow-900/20 border-yellow-600/50">
						<div class="p-4">
							<div class="flex items-center">
								<Clock class="w-5 h-5 text-yellow-400 mr-2" />
								<span class="text-yellow-400 font-medium">Under Review</span>
							</div>
							<p class="text-2xl font-bold text-yellow-300 mt-2">{stats.pendingMods}</p>
						</div>
					</div>
				{/if}

				{#if stats.rejectedMods > 0}
					<div class="card bg-red-900/20 border-red-600/50">
						<div class="p-4">
							<div class="flex items-center">
								<XCircle class="w-5 h-5 text-red-400 mr-2" />
								<span class="text-red-400 font-medium">Rejected</span>
							</div>
							<p class="text-2xl font-bold text-red-300 mt-2">{stats.rejectedMods}</p>
						</div>
					</div>
				{/if}
			</div>
		{/if}

		<div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
			<!-- Mods List -->
			<div class="lg:col-span-2">
				<div class="card">
					<div class="p-6">
						<div class="flex items-center justify-between mb-6">
							<h2 class="text-xl font-semibold text-text-primary">Your Mods</h2>
							<div class="flex items-center space-x-2">
								{#if showBulkActions}
									<span class="text-sm text-text-muted">{selectedCount} selected</span>
									<button on:click={bulkDelete} class="btn btn-danger btn-sm">
										<Trash2 class="w-4 h-4 mr-1" />
										Delete Selected
									</button>
									<button on:click={clearSelection} class="btn btn-outline btn-sm"> Clear </button>
								{:else}
									<button on:click={selectAllMods} class="btn btn-outline btn-sm">
										Select All
									</button>
								{/if}
							</div>
						</div>

						{#if isLoadingMods}
							<div class="text-center py-8">
								<Loading size="md" text="Loading your mods..." />
							</div>
						{:else if userMods.length === 0}
							<div class="text-center py-12">
								<Package class="w-16 h-16 mx-auto text-text-muted mb-4 opacity-50" />
								<h3 class="text-lg font-medium text-text-primary mb-2">No mods yet</h3>
								<p class="text-text-secondary mb-4">
									Upload your first mod to start sharing with the community!
								</p>
								<button on:click={createMod} class="btn btn-primary">
									<Plus class="w-4 h-4 mr-2" />
									Upload Your First Mod
								</button>
							</div>
						{:else}
							<div class="space-y-4">
								{#each userMods as mod (mod.id)}
									<div
										class="flex items-center p-4 bg-slate-800/50 rounded-lg hover:bg-slate-700/50 transition-colors"
									>
										<input
											type="checkbox"
											checked={selectedMods.has(mod.id)}
											on:change={() => toggleModSelection(mod.id)}
											class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-slate-600 bg-slate-800 rounded mr-4"
										/>

										<div class="flex items-center flex-1 min-w-0">
											<!-- Mod Icon -->
											{#if mod.icon}
												<img
													src={mod.icon}
													alt={mod.name}
													class="w-12 h-12 rounded-lg border border-slate-600 mr-4 flex-shrink-0"
												/>
											{:else}
												<div
													class="w-12 h-12 bg-slate-600 rounded-lg flex items-center justify-center mr-4 flex-shrink-0"
												>
													<Package class="w-6 h-6 text-text-muted" />
												</div>
											{/if}

											<!-- Mod Info -->
											<div class="flex-1 min-w-0">
												<div class="flex items-center space-x-2 mb-1">
													<h3 class="text-lg font-medium text-text-primary truncate">
														{mod.name}
													</h3>
													<span class="badge {getStatusBadge(mod)} text-xs">
														{getStatusText(mod)}
													</span>
												</div>
												<p class="text-text-muted text-sm mb-1">
													v{mod.version} • {mod.game?.name || 'Unknown Game'}
												</p>
												<p class="text-text-secondary text-sm truncate">
													{mod.short_description}
												</p>
												{#if mod.is_rejected && mod.rejection_reason}
													<p class="text-red-400 text-sm mt-1">
														Reason: {mod.rejection_reason}
													</p>
												{/if}
											</div>

											<!-- Stats -->
											<div
												class="hidden sm:flex items-center space-x-6 text-sm text-text-muted mx-6"
											>
												<div class="flex items-center">
													<Download class="w-4 h-4 mr-1" />
													{formatNumber(mod.downloads)}
												</div>
												<div class="flex items-center">
													<Heart class="w-4 h-4 mr-1" />
													{formatNumber(mod.likes)}
												</div>
												<div class="flex items-center">
													<MessageCircle class="w-4 h-4 mr-1" />
													{mod.comments}
												</div>
											</div>

											<!-- Actions -->
											<div class="flex items-center space-x-2">
												<a
													href="/mods/{mod.game.slug}/{mod.slug}"
													class="btn btn-outline btn-sm"
													title="View Mod"
												>
													<Eye class="w-4 h-4" />
												</a>
												<button
													on:click={() => editMod(mod.id)}
													class="btn btn-outline btn-sm"
													title="Edit Mod"
												>
													<Edit class="w-4 h-4" />
												</button>
												<button
													on:click={() => deleteMod(mod.id, mod.name)}
													class="btn btn-outline btn-sm text-red-400 hover:text-red-300 hover:border-red-500"
													title="Delete Mod"
												>
													<Trash2 class="w-4 h-4" />
												</button>
											</div>
										</div>
									</div>
								{/each}
							</div>
						{/if}
					</div>
				</div>
			</div>

			<!-- Sidebar -->
			<div class="space-y-6">
				<!-- Quick Actions -->
				<div class="card">
					<div class="p-4">
						<h3 class="text-lg font-semibold text-text-primary mb-4">Quick Actions</h3>
						<div class="space-y-3">
							<button on:click={createMod} class="btn btn-primary w-full">
								<Plus class="w-4 h-4 mr-2" />
								Upload New Mod
							</button>
							<a href="/profile" class="btn btn-outline w-full">
								<User class="w-4 h-4 mr-2" />
								View Profile
							</a>
							<a href="/settings" class="btn btn-outline w-full">
								<Settings class="w-4 h-4 mr-2" />
								Account Settings
							</a>
						</div>
					</div>
				</div>

				<!-- Recent Activity -->
				<div class="card">
					<div class="p-4">
						<h3 class="text-lg font-semibold text-text-primary mb-4">Recent Activity</h3>
						{#if recentActivity.length === 0}
							<p class="text-text-muted text-sm">No recent activity</p>
						{:else}
							<div class="space-y-3">
								{#each recentActivity.slice(0, 5) as activity (activity)}
									<div class="flex items-start space-x-3">
										<div class="p-1 bg-primary-500/20 rounded-full mt-1">
											<svelte:component
												this={getActivityIcon(activity.type)}
												class="w-3 h-3 text-primary-400"
											/>
										</div>
										<div class="flex-1 min-w-0">
											<p class="text-sm text-text-secondary">
												{activity.message}
											</p>
											<p class="text-xs text-text-muted mt-1">
												{formatRelativeTime(activity.timestamp)}
											</p>
										</div>
									</div>
								{/each}
							</div>
						{/if}
					</div>
				</div>

				<!-- Tips -->
				<div class="card bg-blue-900/20 border-blue-600/50">
					<div class="p-4">
						<h3 class="text-lg font-semibold text-blue-300 mb-3">Creator Tips</h3>
						<div class="space-y-2 text-sm text-blue-200">
							<p>• Write clear, detailed descriptions</p>
							<p>• Include screenshots or videos</p>
							<p>• Respond to comments quickly</p>
							<p>• Keep your mods updated</p>
							<p>• Use relevant tags</p>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
