<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { isAuthenticated, isAdmin } from '$lib/stores/auth';
	import { adminApi } from '$lib/api/client';
	import { toast } from '$lib/stores/notifications';
	import Loading from '$lib/components/Loading.svelte';
	import { type Mod, type GameRequest, type Ban, type ActivityItem } from '$lib/types';
	import {
		Shield,
		Users,
		Package,
		Gamepad2,
		TrendingUp,
		Activity,
		AlertTriangle,
		CheckCircle,
		Clock,
		Gavel,
		Edit,
		Plus,
		BarChart3,
		Settings,
		Flag,
		MessageCircle,
		Globe,
		Database,
		Server,
		UserPlus
	} from 'lucide-svelte';

	// Data
	let stats = {
		users: { total: 0, active: 0, new_this_month: 0 },
		mods: { total: 0, approved: 0, pending: 0, rejected: 0 },
		games: { total: 0, active: 0 },
		system: { uptime: '', version: '1.0.0', health: 'good' }
	};
	let recentActivity: ActivityItem[] = [];
	let pendingMods: Mod[] = [];
	let gameRequests: GameRequest[] = [];
	let activeBans: Ban[] = [];

	// UI state
	let isLoadingActivity = true;
	let activeTab = 'overview';

	// Tabs
	const tabs = [
		{ id: 'overview', label: 'Overview', icon: BarChart3 },
		{ id: 'users', label: 'Users', icon: Users },
		{ id: 'mods', label: 'Mods', icon: Package },
		{ id: 'games', label: 'Games', icon: Gamepad2 },
		{ id: 'bans', label: 'Bans', icon: Gavel },
		{ id: 'system', label: 'System', icon: Settings }
	];

	onMount(async () => {
		if (!$isAuthenticated) {
			goto('/auth/login?redirect=/admin');
			return;
		}

		if (!$isAdmin) {
			toast.error('Access denied', 'You do not have permission to access the admin panel.');
			goto('/');
			return;
		}

		await Promise.all([
			loadStats(),
			loadRecentActivity(),
			loadPendingMods(),
			loadGameRequests(),
			loadActiveBans()
		]);
	});

	// Load admin statistics
	async function loadStats() {
		try {
			const response = await adminApi.getStats();
			if (response.success && response.data) {
				stats = {
					...response.data,
					users: response.data.users ?? { total: 0, active: 0, new_this_month: 0 },
					mods: response.data.mods ?? { total: 0, approved: 0, pending: 0, rejected: 0 },
					games: response.data.games ?? { total: 0, active: 0 },
					system: response.data.system ?? { uptime: '', version: '1.0.0', health: 'good' }
				};
			} else {
				// Mock data for demonstration
				stats = {
					users: { total: 15420, active: 12850, new_this_month: 342 },
					mods: { total: 8750, approved: 8200, pending: 450, rejected: 100 },
					games: { total: 12, active: 11 },
					system: { uptime: '72h 30m', version: '1.0.0', health: 'good' }
				};
			}
		} catch (error) {
			console.error('Error loading admin stats:', error);
			// Use mock data on error
			stats = {
				users: { total: 15420, active: 12850, new_this_month: 342 },
				mods: { total: 8750, approved: 8200, pending: 450, rejected: 100 },
				games: { total: 12, active: 11 },
				system: { uptime: '72h 30m', version: '1.0.0', health: 'good' }
			};
		}
	}

	// Load recent activity
	async function loadRecentActivity() {
		isLoadingActivity = true;
		try {
			const response = await adminApi.getActivity({ limit: 10 });
			if (response.success && response.data) {
				// Make sure every activity has a unique id for Svelte
				recentActivity = response.data.map((item, index) => ({
					...item,
					id: item.id ?? index,
					message: createActivityMessage(item)
				}));

				console.log('Recent activity loaded:', recentActivity);
			} else {
				recentActivity = [];
			}
		} catch (error) {
			console.error('Error loading recent activity:', error);
			recentActivity = [];
		} finally {
			isLoadingActivity = false;
		}
	}

	// Load pending mods
	async function loadPendingMods() {
		try {
			const response = await adminApi.getPendingMods({ per_page: 10 });
			if (response.success && response.data) {
				pendingMods = response.data.data || [];
			} else {
				pendingMods = [];
				if (response.error) {
					console.error('Failed to load pending mods:', response.error);
				}
			}
		} catch (error) {
			console.error('Error loading pending mods:', error);
			pendingMods = [];
		}
	}

	// Load game requests
	async function loadGameRequests() {
		try {
			const response = await adminApi.getGameRequests();
			if (response.success && response.data) {
				gameRequests = response.data.data || [];
			} else {
				console.error('Failed to load game requests:', response.error);
			}
		} catch (error) {
			console.error('Error loading game requests:', error);
		}
	}

	// Load active bans
	async function loadActiveBans() {
		try {
			const response = await adminApi.getBans({ active: true });
			if (response.success && response.data) {
				activeBans = response.data.data || [];
			} else {
				console.error('Failed to load active bans:', response.error);
			}
		} catch (error) {
			console.error('Error loading active bans:', error);
		}
	}

	// Format numbers
	function formatNumber(num: number): string {
		// Filter null
		if (num === null) return 'N/A';

		if (num >= 1000000) {
			return (num / 1000000).toFixed(1) + 'M';
		}
		if (num >= 1000) {
			return (num / 1000).toFixed(1) + 'K';
		}
		return num.toString();
	}

	// Format relative time
	function formatRelativeTime(dateString: string): string {
		const date = new Date(dateString);
		const now = new Date();
		const diffTime = Math.abs(now.getTime() - date.getTime());
		const diffMinutes = Math.floor(diffTime / (1000 * 60));
		const diffHours = Math.floor(diffTime / (1000 * 60 * 60));
		const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));

		if (diffMinutes < 60) return `${diffMinutes}m ago`;
		if (diffHours < 24) return `${diffHours}h ago`;
		return `${diffDays}d ago`;
	}

	// Get activity icon
	function getActivityIcon(type: string) {
		switch (type) {
			case 'mod_upload':
				return Package;
			case 'user_register':
				return UserPlus;
			case 'user':
				return Users;
			case 'game_request':
				return Gamepad2;
			case 'mod_report':
				return Flag;
			case 'comment':
				return MessageCircle;
			default:
				return Activity;
		}
	}

	// Get activity color
	function getActivityColor(type: string, severity?: string) {
		if (severity === 'high') return 'text-red-400';
		switch (type) {
			case 'mod_upload':
				return 'text-green-400';
			case 'user_register':
			case 'user':
				return 'text-blue-400';
			case 'game_request':
				return 'text-purple-400';
			case 'mod_report':
				return 'text-yellow-400';
			default:
				return 'text-text-muted';
		}
	}

	// Create activity message
	function createActivityMessage(activity: ActivityItem) {
		console.log('Activity:', activity);
		switch (activity.type) {
			case 'user':
			case 'user_register':
				return `${activity.user || 'Someone'} joined the platform`;
			case 'game_request':
				return `${activity.user || 'Someone'} requested a new game: "${activity.name}"`;
			case 'mod_upload':
				return `${activity.user || 'Someone'} uploaded a new mod: "${activity.name}"`;
			case 'mod_report':
				return `${activity.user || 'Someone'} reported a mod: "${activity.name}"`;
			case 'comment':
				return `${activity.user || 'Someone'} commented: "${activity.name}"`;
			default:
				return activity.name || 'Activity occurred';
		}
	}

	// Approve mod
	async function approveMod(modId: number, modName: string) {
		try {
			const response = await adminApi.approveMod?.(modId);
			if (response?.success) {
				toast.success('Mod approved', `"${modName}" has been approved and published.`);
				await loadPendingMods();
				await loadStats();
			} else {
				toast.error('Failed to approve mod', response?.error || 'Unknown error');
			}
		} catch (error) {
			console.error('Error approving mod:', error);
			toast.error('Error', 'Failed to approve mod');
		}
	}

	// Reject mod
	async function rejectMod(modId: number, modName: string) {
		const reason = prompt('Please provide a reason for rejection:');
		if (!reason) return;

		try {
			const response = await adminApi.rejectMod?.(modId, reason);
			if (response?.success) {
				toast.success('Mod rejected', `"${modName}" has been rejected.`);
				await loadPendingMods();
				await loadStats();
			} else {
				toast.error('Failed to reject mod', response?.error || 'Unknown error');
			}
		} catch (error) {
			console.error('Error rejecting mod:', error);
			toast.error('Error', 'Failed to reject mod');
		}
	}

	// Approve game request
	async function approveGameRequest(requestId: number, gameName: string) {
		try {
			const response = await adminApi.approveGameRequest(requestId);
			if (response.success) {
				toast.success('Game approved', `"${gameName}" has been added to the platform.`);
				await loadGameRequests();
			} else {
				toast.error('Error', response.error || 'Failed to approve game request');
			}
		} catch (error) {
			console.error('Error approving game request:', error);
			toast.error('Error', 'Failed to approve game request');
		}
	}

	// Reject game request
	async function rejectGameRequest(requestId: number, gameName: string) {
		const reason = prompt('Please provide a reason for rejection:');
		if (!reason) return;

		try {
			const response = await adminApi.rejectGameRequest(requestId, reason);
			if (response.success) {
				toast.success('Game request rejected', `"${gameName}" request has been rejected.`);
				await loadGameRequests();
			} else {
				toast.error('Error', response.error || 'Failed to reject game request');
			}
		} catch (error) {
			console.error('Error rejecting game request:', error);
			toast.error('Error', 'Failed to reject game request');
		}
	}

	// Unban user
	async function unbanUser(banId: number, username: string) {
		try {
			const response = await adminApi.unbanUser(banId);
			if (response.success) {
				toast.success('User unbanned', `${username} has been unbanned.`);
				await loadActiveBans();
			} else {
				toast.error('Failed to unban user', response.error);
			}
		} catch (error) {
			console.error('Error unbanning user:', error);
			toast.error('Error', 'Failed to unban user');
		}
	}
</script>

<svelte:head>
	<title>Admin Dashboard - Azurite</title>
	<meta name="description" content="Azurite admin panel for site management and moderation." />
</svelte:head>

<div class="min-h-screen bg-background-primary">
	<!-- Header -->
	<div class="bg-gradient-to-r from-red-600/20 to-red-700/20 border-b border-slate-700">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
			<div class="flex items-center space-x-3">
				<Shield class="w-8 h-8 text-red-400" />
				<div>
					<h1 class="text-3xl font-bold text-text-primary">Admin Dashboard</h1>
					<p class="text-text-secondary mt-1">Site administration and moderation tools</p>
				</div>
			</div>
		</div>
	</div>

	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
		<!-- Tabs -->
		<div class="flex space-x-1 bg-slate-800 rounded-lg p-1 mb-8 overflow-x-auto">
			{#each tabs as tab (tab.id)}
				<button
					on:click={() => (activeTab = tab.id)}
					class="flex items-center px-4 py-2 rounded-md text-sm font-medium whitespace-nowrap transition-colors {activeTab ===
					tab.id
						? 'bg-primary-600 text-white'
						: 'text-text-secondary hover:text-text-primary hover:bg-slate-700'}"
				>
					<svelte:component this={tab.icon} class="w-4 h-4 mr-2" />
					{tab.label}
				</button>
			{/each}
		</div>

		{#if activeTab === 'overview'}
			<!-- Overview Tab -->
			<div class="space-y-8">
				<!-- Stats Grid -->
				<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
					<!-- Users -->
					<div class="card">
						<div class="p-6">
							<div class="flex items-center">
								<div class="p-2 bg-blue-500/20 rounded-lg">
									<Users class="w-6 h-6 text-blue-400" />
								</div>
								<div class="ml-4">
									<p class="text-text-muted text-sm">Total Users</p>
									<p class="text-2xl font-bold text-text-primary">
										{formatNumber(stats.users?.total ?? 0)}
									</p>
								</div>
							</div>
							<div class="mt-4 flex items-center text-sm text-text-muted">
								<TrendingUp class="w-4 h-4 mr-1 text-green-400" />
								{stats.users.new_this_month} new this month
							</div>
						</div>
					</div>

					<!-- Mods -->
					<div class="card">
						<div class="p-6">
							<div class="flex items-center">
								<div class="p-2 bg-green-500/20 rounded-lg">
									<Package class="w-6 h-6 text-green-400" />
								</div>
								<div class="ml-4">
									<p class="text-text-muted text-sm">Total Mods</p>
									<p class="text-2xl font-bold text-text-primary">
										{formatNumber(stats.mods?.total ?? 0)}
									</p>
								</div>
							</div>
							<div class="mt-4 flex items-center text-sm text-text-muted">
								<Clock class="w-4 h-4 mr-1 text-yellow-400" />
								{stats.mods.pending} pending review
							</div>
						</div>
					</div>

					<!-- Games -->
					<div class="card">
						<div class="p-6">
							<div class="flex items-center">
								<div class="p-2 bg-purple-500/20 rounded-lg">
									<Gamepad2 class="w-6 h-6 text-purple-400" />
								</div>
								<div class="ml-4">
									<p class="text-text-muted text-sm">Games</p>
									<p class="text-2xl font-bold text-text-primary">{stats.games.total}</p>
								</div>
							</div>
							<div class="mt-4 flex items-center text-sm text-text-muted">
								<CheckCircle class="w-4 h-4 mr-1 text-green-400" />
								{stats.games.active} active
							</div>
						</div>
					</div>

					<!-- System Health -->
					<div class="card">
						<div class="p-6">
							<div class="flex items-center">
								<div class="p-2 bg-emerald-500/20 rounded-lg">
									<Server class="w-6 h-6 text-emerald-400" />
								</div>
								<div class="ml-4">
									<p class="text-text-muted text-sm">System Status</p>
									<p class="text-lg font-semibold text-green-400">Healthy</p>
								</div>
							</div>
							<div class="mt-4 flex items-center text-sm text-text-muted">
								<Globe class="w-4 h-4 mr-1" />
								Uptime: {stats.system.uptime}
							</div>
						</div>
					</div>
				</div>

				<!-- Content Grid -->
				<div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
					<!-- Recent Activity -->
					<div class="card">
						<div class="p-6">
							<h2 class="text-xl font-semibold text-text-primary mb-4 flex items-center">
								<Activity class="w-5 h-5 mr-2" />
								Recent Activity
							</h2>

							{#if isLoadingActivity}
								<div class="text-center py-4">
									<Loading size="md" text="Loading activity..." />
								</div>
							{:else if recentActivity.length === 0}
								<p class="text-text-muted text-center py-8">No recent activity</p>
							{:else}
								<div class="space-y-4">
									{#each recentActivity as activity (activity.id)}
										<div class="flex items-start space-x-3">
											<div class="p-1 bg-slate-700 rounded-full">
												<svelte:component
													this={getActivityIcon(activity.type)}
													class="w-4 h-4 {getActivityColor(activity.type, activity.severity)}"
												/>
											</div>
											<div class="flex-1 min-w-0">
												<p class="text-sm text-text-secondary">
													{activity.message}
												</p>
												<p class="text-xs text-text-muted mt-1">
													{formatRelativeTime(activity.created_at)}
												</p>
											</div>
										</div>
									{/each}
								</div>
							{/if}
						</div>
					</div>

					<!-- Pending Actions -->
					<div class="card">
						<div class="p-6">
							<h2 class="text-xl font-semibold text-text-primary mb-4 flex items-center">
								<Clock class="w-5 h-5 mr-2" />
								Pending Actions
							</h2>

							<div class="space-y-4">
								{#if pendingMods.length > 0}
									<div>
										<h3 class="text-sm font-medium text-text-primary mb-2">
											Mod Reviews ({pendingMods.length})
										</h3>
										<div class="space-y-2">
											{#each pendingMods.slice(0, 3) as mod (mod.id)}
												<div class="flex items-center justify-between p-3 bg-slate-800 rounded-lg">
													<div>
														<p class="text-sm font-medium text-text-primary">{mod.name}</p>
														<p class="text-xs text-text-muted">by {mod.owner.display_name}</p>
													</div>
													<div class="flex space-x-1">
														<button
															on:click={() => approveMod(mod.id, mod.name)}
															class="btn btn-sm bg-green-600 text-white hover:bg-green-700"
															title="Approve"
														>
															<CheckCircle class="w-4 h-4" />
														</button>
														<button
															on:click={() => rejectMod(mod.id, mod.name)}
															class="btn btn-sm bg-red-600 text-white hover:bg-red-700"
															title="Reject"
														>
															<AlertTriangle class="w-4 h-4" />
														</button>
													</div>
												</div>
											{/each}
										</div>
									</div>
								{/if}

								{#if gameRequests.length > 0}
									<div>
										<h3 class="text-sm font-medium text-text-primary mb-2">
											Game Requests ({gameRequests.length})
										</h3>
										<div class="space-y-2">
											{#each gameRequests.slice(0, 2) as request (request.id)}
												<div class="flex items-center justify-between p-3 bg-slate-800 rounded-lg">
													<div>
														<p class="text-sm font-medium text-text-primary">{request.name}</p>
														<p class="text-xs text-text-muted">
															by {request.requested_by.display_name}
														</p>
													</div>
													<div class="flex space-x-1">
														<button
															on:click={() => approveGameRequest(request.id, request.name)}
															class="btn btn-sm bg-green-600 text-white hover:bg-green-700"
															title="Approve"
														>
															<CheckCircle class="w-4 h-4" />
														</button>
														<button
															on:click={() => rejectGameRequest(request.id, request.name)}
															class="btn btn-sm bg-red-600 text-white hover:bg-red-700"
															title="Reject"
														>
															<AlertTriangle class="w-4 h-4" />
														</button>
													</div>
												</div>
											{/each}
										</div>
									</div>
								{/if}

								{#if pendingMods.length === 0 && gameRequests.length === 0}
									<p class="text-text-muted text-center py-4">No pending actions</p>
								{/if}
							</div>
						</div>
					</div>
				</div>
			</div>
		{:else if activeTab === 'bans'}
			<!-- Bans Tab -->
			<div class="card">
				<div class="p-6">
					<div class="flex items-center justify-between mb-6">
						<h2 class="text-xl font-semibold text-text-primary flex items-center">
							<Gavel class="w-5 h-5 mr-2" />
							Active Bans ({activeBans.length})
						</h2>
						<button class="btn btn-primary">
							<Plus class="w-4 h-4 mr-2" />
							Create Ban
						</button>
					</div>

					{#if activeBans.length === 0}
						<div class="text-center py-12">
							<Gavel class="w-16 h-16 mx-auto text-text-muted mb-4 opacity-50" />
							<h3 class="text-lg font-medium text-text-primary mb-2">No active bans</h3>
							<p class="text-text-secondary">All users are currently in good standing.</p>
						</div>
					{:else}
						<div class="space-y-4">
							{#each activeBans as ban (ban.id)}
								<div class="flex items-center justify-between p-4 bg-slate-800 rounded-lg">
									<div class="flex items-center space-x-4">
										<div class="w-10 h-10 bg-red-600 rounded-full flex items-center justify-center">
											<Gavel class="w-5 h-5 text-white" />
										</div>
										<div>
											<h3 class="font-medium text-text-primary">
												{ban.user?.display_name || ban.user?.username}
											</h3>
											<p class="text-sm text-text-muted">
												Banned by {ban.banned_by_user?.display_name}
											</p>
											<p class="text-sm text-red-400">{ban.reason}</p>
											<p class="text-xs text-text-muted mt-1">
												Expires: {new Date(ban.expires_at).toLocaleDateString()}
											</p>
										</div>
									</div>
									<div class="flex items-center space-x-2">
										<button
											on:click={() =>
												unbanUser(ban.id, ban.user?.display_name || ban.user?.username)}
											class="btn btn-outline btn-sm text-green-400 hover:text-green-300"
										>
											Unban
										</button>
										<button class="btn btn-outline btn-sm">
											<Edit class="w-4 h-4" />
										</button>
									</div>
								</div>
							{/each}
						</div>
					{/if}
				</div>
			</div>
		{:else if activeTab === 'system'}
			<!-- System Tab -->
			<div class="space-y-6">
				<!-- System Information -->
				<div class="card">
					<div class="p-6">
						<h2 class="text-xl font-semibold text-text-primary mb-4 flex items-center">
							<Server class="w-5 h-5 mr-2" />
							System Information
						</h2>
						<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
							<div>
								<p class="text-text-muted text-sm">Version</p>
								<p class="text-lg font-semibold text-text-primary">{stats.system.version}</p>
							</div>
							<div>
								<p class="text-text-muted text-sm">Uptime</p>
								<p class="text-lg font-semibold text-text-primary">{stats.system.uptime}</p>
							</div>
							<div>
								<p class="text-text-muted text-sm">Health Status</p>
								<p class="text-lg font-semibold text-green-400 capitalize">{stats.system.health}</p>
							</div>
						</div>
					</div>
				</div>

				<!-- Database Statistics -->
				<div class="card">
					<div class="p-6">
						<h2 class="text-xl font-semibold text-text-primary mb-4 flex items-center">
							<Database class="w-5 h-5 mr-2" />
							Database Statistics
						</h2>
						<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
							<div>
								<p class="text-text-muted text-sm">Total Users</p>
								<p class="text-xl font-bold text-text-primary">
									{formatNumber(stats.users?.total ?? 0)}
								</p>
							</div>
							<div>
								<p class="text-text-muted text-sm">Total Mods</p>
								<p class="text-xl font-bold text-text-primary">
									{formatNumber(stats.mods?.total ?? 0)}
								</p>
							</div>
							<div>
								<p class="text-text-muted text-sm">Total Games</p>
								<p class="text-xl font-bold text-text-primary">{stats.games.total}</p>
							</div>
							<div>
								<p class="text-text-muted text-sm">Active Bans</p>
								<p class="text-xl font-bold text-text-primary">{activeBans.length}</p>
							</div>
						</div>
					</div>
				</div>
			</div>
		{:else}
			<!-- Other tabs placeholder -->
			<div class="card">
				<div class="p-6 text-center">
					<h2 class="text-xl font-semibold text-text-primary mb-2">
						{tabs.find((t) => t.id === activeTab)?.label} Management
					</h2>
					<p class="text-text-secondary">This section is under development.</p>
				</div>
			</div>
		{/if}
	</div>
</div>
