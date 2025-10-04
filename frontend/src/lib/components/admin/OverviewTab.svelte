<script lang="ts">
	import {
		Users,
		Package,
		Gamepad2,
		Server,
		Clock,
		CheckCircle,
		Activity,
		AlertTriangle
	} from 'lucide-svelte';
	import StatsCard from './StatsCard.svelte';
	import Loading from '$lib/components/Loading.svelte';
	import type { Mod, GameRequest, ActivityItem } from '$lib/types';

	export let stats: any;
	export let recentActivity: ActivityItem[];
	export let pendingMods: Mod[];
	export let gameRequests: GameRequest[];
	export let isLoadingActivity: boolean;
	export let onApproveMod: (modId: number, modName: string) => Promise<void>;
	export let onRejectMod: (modId: number, modName: string) => Promise<void>;
	export let onApproveGameRequest: (requestId: number, gameName: string) => Promise<void>;
	export let onRejectGameRequest: (requestId: number, gameName: string) => Promise<void>;

	function formatNumber(num: number): string {
		if (num === null) return 'N/A';
		if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M';
		if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
		return num.toString();
	}

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

	function getActivityIcon(type: string) {
		switch (type) {
			case 'mod_upload':
				return Package;
			case 'user_register':
			case 'user':
				return Users;
			case 'game_request':
				return Gamepad2;
			default:
				return Activity;
		}
	}

	function getActivityColor(type: string) {
		switch (type) {
			case 'mod_upload':
				return 'text-green-400';
			case 'user_register':
			case 'user':
				return 'text-blue-400';
			case 'game_request':
				return 'text-purple-400';
			default:
				return 'text-text-muted';
		}
	}
</script>

<div class="space-y-8">
	<!-- Stats Grid -->
	<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
		<StatsCard
			title="Total Users"
			value={formatNumber(stats.users?.total ?? 0)}
			subtitle="{stats.users?.new_this_month ?? 0} new this month"
			icon={Users}
			iconColor="text-blue-400"
			iconBg="bg-blue-500/20"
		/>

		<StatsCard
			title="Total Mods"
			value={formatNumber(stats.mods?.total ?? 0)}
			subtitle="{stats.mods?.pending ?? 0} pending review"
			icon={Package}
			iconColor="text-green-400"
			iconBg="bg-green-500/20"
		/>

		<StatsCard
			title="Games"
			value={stats.games?.total ?? 0}
			subtitle="{stats.games?.active ?? 0} active"
			icon={Gamepad2}
			iconColor="text-purple-400"
			iconBg="bg-purple-500/20"
		/>

		<StatsCard
			title="System Status"
			value="Healthy"
			subtitle="Uptime: {stats.system?.uptime || 'N/A'}"
			icon={Server}
			iconColor="text-emerald-400"
			iconBg="bg-emerald-500/20"
		/>
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
										class="w-4 h-4 {getActivityColor(activity.type)}"
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
											<p class="text-xs text-text-muted">by {mod.owner?.display_name}</p>
										</div>
										<div class="flex space-x-1">
											<button
												on:click={() => onApproveMod(mod.id, mod.name)}
												class="btn btn-sm bg-green-600 text-white hover:bg-green-700"
												title="Approve"
											>
												<CheckCircle class="w-4 h-4" />
											</button>
											<button
												on:click={() => onRejectMod(mod.id, mod.name)}
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
												by {request.requested_by?.display_name}
											</p>
										</div>
										<div class="flex space-x-1">
											<button
												on:click={() => onApproveGameRequest(request.id, request.name)}
												class="btn btn-sm bg-green-600 text-white hover:bg-green-700"
												title="Approve"
											>
												<CheckCircle class="w-4 h-4" />
											</button>
											<button
												on:click={() => onRejectGameRequest(request.id, request.name)}
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
