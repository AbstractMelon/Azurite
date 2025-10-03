<script lang="ts">
	import { onMount } from 'svelte';
	import { adminApi } from '$lib/api/client';
	import { toast } from '$lib/stores/notifications';
	import { Package, CheckCircle, AlertTriangle } from 'lucide-svelte';
	import type { Mod } from '$lib/types';

	let pendingMods: Mod[] = [];
	let loading = true;

	async function loadPendingMods() {
		loading = true;
		try {
			const response = await adminApi.getPendingMods({ per_page: 20 });
			if (response.success && response.data) {
				pendingMods = response.data.data || [];
			}
		} catch (error) {
			console.error('Error loading pending mods:', error);
		} finally {
			loading = false;
		}
	}

	async function approveMod(modId: number, modName: string) {
		try {
			const response = await adminApi.approveMod(modId);
			if (response.success) {
				toast.success('Mod approved', `"${modName}" has been approved.`);
				await loadPendingMods();
			}
		} catch (error) {
			toast.error('Error', 'Failed to approve mod');
		}
	}

	async function rejectMod(modId: number, modName: string) {
		const reason = prompt('Reason for rejection:');
		if (!reason) return;

		try {
			const response = await adminApi.rejectMod(modId, reason);
			if (response.success) {
				toast.success('Mod rejected', `"${modName}" has been rejected.`);
				await loadPendingMods();
			}
		} catch (error) {
			toast.error('Error', 'Failed to reject mod');
		}
	}

	onMount(() => {
		loadPendingMods();
	});
</script>

<div class="card">
	<div class="p-6">
		<h2 class="text-xl font-semibold text-text-primary mb-6 flex items-center">
			<Package class="w-5 h-5 mr-2" />
			Pending Mods ({pendingMods.length})
		</h2>

		{#if loading}
			<div class="text-center py-8">Loading...</div>
		{:else if pendingMods.length === 0}
			<div class="text-center py-12">
				<Package class="w-16 h-16 mx-auto text-text-muted mb-4 opacity-50" />
				<h3 class="text-lg font-medium text-text-primary mb-2">No pending mods</h3>
				<p class="text-text-secondary">All mods have been reviewed.</p>
			</div>
		{:else}
			<div class="space-y-3">
				{#each pendingMods as mod (mod.id)}
					<div class="flex items-center justify-between p-4 bg-slate-800 rounded-lg">
						<div>
							<p class="text-sm font-medium text-text-primary">{mod.name}</p>
							<p class="text-xs text-text-muted">by {mod.owner?.display_name}</p>
						</div>
						<div class="flex space-x-2">
							<button
								on:click={() => approveMod(mod.id, mod.name)}
								class="btn btn-sm bg-green-600 text-white hover:bg-green-700"
							>
								<CheckCircle class="w-4 h-4" />
							</button>
							<button
								on:click={() => rejectMod(mod.id, mod.name)}
								class="btn btn-sm bg-red-600 text-white hover:bg-red-700"
							>
								<AlertTriangle class="w-4 h-4" />
							</button>
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</div>
