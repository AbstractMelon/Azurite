<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { isAuthenticated, isAdmin } from '$lib/stores/auth';
	import { adminApi } from '$lib/api/client';
	import { toast } from '$lib/stores/notifications';
	import { type Mod, type GameRequest, type Ban, type ActivityItem, type Game } from '$lib/types';
	import { Shield, BarChart3, Users, Package, Gamepad2, Gavel, Settings } from 'lucide-svelte';
	
	// Import tab components
	import OverviewTab from '$lib/components/admin/OverviewTab.svelte';
	import UsersTab from '$lib/components/admin/UsersTab.svelte';
	import ModsTab from '$lib/components/admin/ModsTab.svelte';
	import GamesTab from '$lib/components/admin/GamesTab.svelte';
	import BansTab from '$lib/components/admin/BansTab.svelte';
	import SystemTab from '$lib/components/admin/SystemTab.svelte';

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
	let allGames: Game[] = [];
	let allGameRequests: GameRequest[] = [];

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
			}
		} catch (error) {
			console.error('Error loading admin stats:', error);
		}
	}

	// Load recent activity
	async function loadRecentActivity() {
		isLoadingActivity = true;
		try {
			const response = await adminApi.getActivity({ limit: 10 });
			if (response.success && response.data) {
				recentActivity = response.data.map((item: any, index: number) => ({
					...item,
					id: item.id ?? index,
					message: createActivityMessage(item)
				}));
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
			}
		} catch (error) {
			console.error('Error loading active bans:', error);
		}
	}

	// Load all games for admin
	async function loadAllGames() {
		try {
			const response = await adminApi.getAllGames?.();
			if (response?.success && response.data) {
				allGames = response.data.data || [];
			}
		} catch (error) {
			console.error('Error loading games:', error);
		}
	}

	// Load all game requests for admin
	async function loadAllGameRequests() {
		try {
			const response = await adminApi.getGameRequests();
			if (response.success && response.data) {
				allGameRequests = response.data.data || [];
			}
		} catch (error) {
			console.error('Error loading game requests:', error);
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

	// Combined load function for games tab
	async function loadGamesData() {
		await Promise.all([loadAllGames(), loadAllGameRequests()]);
	}

	// Create activity message
	function createActivityMessage(activity: any) {
		switch (activity.type) {
			case 'user':
			case 'user_register':
				return `${activity.user || 'Someone'} joined the platform`;
			case 'game_request':
				return `${activity.user || 'Someone'} requested a new game: "${activity.name}"`;
			case 'mod_upload':
				return `${activity.user || 'Someone'} uploaded a new mod: "${activity.name}"`;
			default:
				return activity.name || 'Activity occurred';
		}
	}

	// Load games when tab changes
	$: if (activeTab === 'games') {
		loadAllGames();
		loadAllGameRequests();
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
					onclick={() => (activeTab = tab.id)}
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

		<!-- Tab Content -->
		{#if activeTab === 'overview'}
			<OverviewTab
				{stats}
				{recentActivity}
				{pendingMods}
				{gameRequests}
				{isLoadingActivity}
				onApproveMod={approveMod}
				onRejectMod={rejectMod}
				onApproveGameRequest={approveGameRequest}
				onRejectGameRequest={rejectGameRequest}
			/>
		{:else if activeTab === 'users'}
			<UsersTab />
		{:else if activeTab === 'mods'}
			<ModsTab />
		{:else if activeTab === 'games'}
			<GamesTab {allGames} {allGameRequests} onLoad={loadGamesData} />
		{:else if activeTab === 'bans'}
			<BansTab />
		{:else if activeTab === 'system'}
			<SystemTab {stats} {activeBans} />
		{/if}
	</div>
</div>
