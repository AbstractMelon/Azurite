<script lang="ts">
	import { Server, Database, Activity, TrendingUp } from 'lucide-svelte';

	export let stats: any;
	export let activeBans: any[] = [];

	function formatNumber(num: number): string {
		if (num === null) return 'N/A';
		if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M';
		if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
		return num.toString();
	}
</script>

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
					<p class="text-lg font-semibold text-text-primary">
						{stats.system?.version || '1.0.0'}
					</p>
				</div>
				<div>
					<p class="text-text-muted text-sm">Uptime</p>
					<p class="text-lg font-semibold text-text-primary">
						{stats.system?.uptime || 'N/A'}
					</p>
				</div>
				<div>
					<p class="text-text-muted text-sm">Health Status</p>
					<p class="text-lg font-semibold text-green-400 capitalize">
						{stats.system?.health || 'Good'}
					</p>
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
					<p class="text-xl font-bold text-text-primary">{stats.games?.total ?? 0}</p>
				</div>
				<div>
					<p class="text-text-muted text-sm">Active Bans</p>
					<p class="text-xl font-bold text-text-primary">{activeBans.length}</p>
				</div>
			</div>
		</div>
	</div>

	<!-- Activity Metrics -->
	<div class="card">
		<div class="p-6">
			<h2 class="text-xl font-semibold text-text-primary mb-4 flex items-center">
				<Activity class="w-5 h-5 mr-2" />
				Activity Metrics
			</h2>
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
				<div>
					<p class="text-text-muted text-sm">Total Downloads</p>
					<p class="text-xl font-bold text-text-primary">
						{formatNumber(stats.system?.total_downloads ?? 0)}
					</p>
				</div>
				<div>
					<p class="text-text-muted text-sm">Total Likes</p>
					<p class="text-xl font-bold text-text-primary">
						{formatNumber(stats.system?.total_likes ?? 0)}
					</p>
				</div>
				<div>
					<p class="text-text-muted text-sm">Total Comments</p>
					<p class="text-xl font-bold text-text-primary">
						{formatNumber(stats.system?.total_comments ?? 0)}
					</p>
				</div>
			</div>
		</div>
	</div>

	<!-- Performance Indicators -->
	<div class="card">
		<div class="p-6">
			<h2 class="text-xl font-semibold text-text-primary mb-4 flex items-center">
				<TrendingUp class="w-5 h-5 mr-2" />
				Growth Indicators
			</h2>
			<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
				<div class="p-4 bg-slate-800 rounded-lg">
					<p class="text-text-muted text-sm mb-2">New Users (30 days)</p>
					<p class="text-2xl font-bold text-green-400">
						{formatNumber(stats.users?.new_this_month ?? 0)}
					</p>
				</div>
				<div class="p-4 bg-slate-800 rounded-lg">
					<p class="text-text-muted text-sm mb-2">Pending Mod Reviews</p>
					<p class="text-2xl font-bold text-yellow-400">
						{formatNumber(stats.mods?.pending ?? 0)}
					</p>
				</div>
			</div>
		</div>
	</div>
</div>
