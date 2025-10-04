<script lang="ts">
	import { onMount } from 'svelte';
	import { adminApi } from '$lib/api/client';
	import { toast } from '$lib/stores/notifications';
	import { Gavel, Plus, X } from 'lucide-svelte';
	import type { Ban } from '$lib/types';

	interface Props {
		data?: any;
	}

	let { data }: Props = $props();

	let activeBans = $state<Ban[]>([]);
	let loading = $state(true);
	let showCreateModal = $state(false);
	let newBan = $state({
		user_id: null as number | null,
		ip_address: '',
		reason: '',
		duration: undefined as number | undefined
	});

	async function loadBans() {
		loading = true;
		try {
			const response = await adminApi.getBans({ active: true });
			if (response.success && response.data) {
				activeBans = response.data.data || [];
			}
		} catch (error) {
			console.error('Error loading bans:', error);
		} finally {
			loading = false;
		}
	}

	async function createBan() {
		if (!newBan.reason) {
			toast.error('Error', 'Reason is required');
			return;
		}

		try {
			const response = await adminApi.createBan({
				user_id: newBan.user_id ?? undefined,
				ip_address: newBan.ip_address,
				reason: newBan.reason,
				duration: newBan.duration
			});

			if (response.success) {
				toast.success('Ban created', 'User has been banned successfully');
				showCreateModal = false;
				newBan = { user_id: null, ip_address: '', reason: '', duration: undefined };
				await loadBans();
			}
		} catch (error) {
			toast.error('Error', 'Failed to create ban');
		}
	}

	async function unbanUser(banId: number, username: string) {
		try {
			const response = await adminApi.unbanUser(banId);
			if (response.success) {
				toast.success('User unbanned', `${username} has been unbanned.`);
				await loadBans();
			}
		} catch (error) {
			toast.error('Error', 'Failed to unban user');
		}
	}

	onMount(() => {
		loadBans();
	});
</script>

<div class="card">
	<div class="p-6">
		<div class="flex items-center justify-between mb-6">
			<h2 class="text-xl font-semibold text-text-primary flex items-center">
				<Gavel class="w-5 h-5 mr-2" />
				Active Bans ({activeBans.length})
			</h2>
			<button class="btn btn-primary" onclick={() => (showCreateModal = true)}>
				<Plus class="w-4 h-4 mr-2" />
				Create Ban
			</button>
		</div>

		{#if loading}
			<div class="text-center py-8">Loading...</div>
		{:else if activeBans.length === 0}
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
									{ban.user?.display_name || ban.user?.username || ban.ip_address}
								</h3>
								<p class="text-sm text-text-muted">
									Banned by {ban.banned_by_user?.display_name}
								</p>
								<p class="text-sm text-red-400">{ban.reason}</p>
								{#if ban.expires_at}
									<p class="text-xs text-text-muted mt-1">
										Expires: {new Date(ban.expires_at).toLocaleDateString()}
									</p>
								{/if}
							</div>
						</div>
						<div class="flex items-center space-x-2">
							<button
								onclick={() =>
									unbanUser(ban.id, ban.user?.display_name || ban.user?.username || 'User')}
								class="btn btn-outline btn-sm text-green-400 hover:text-green-300"
							>
								Unban
							</button>
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</div>

<!-- Create Ban Modal -->
{#if showCreateModal}
	<div
		class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
		onclick={() => (showCreateModal = false)}
	>
		<div
			class="bg-slate-900 rounded-lg max-w-md w-full p-6"
			onclick={(event) => event.stopPropagation()}
		>
			<div class="flex items-center justify-between mb-6">
				<h2 class="text-2xl font-bold text-text-primary">Create Ban</h2>
				<button
					onclick={() => (showCreateModal = false)}
					class="text-text-muted hover:text-text-primary"
				>
					<X class="w-6 h-6" />
				</button>
			</div>

			<div class="space-y-4">
				<div>
					<label class="block text-sm font-medium text-text-secondary mb-2"
						>User ID (optional)</label
					>
					<input
						type="number"
						bind:value={newBan.user_id}
						class="input w-full"
						placeholder="Leave empty for IP ban"
					/>
				</div>

				<div>
					<label class="block text-sm font-medium text-text-secondary mb-2"
						>IP Address (optional)</label
					>
					<input
						type="text"
						bind:value={newBan.ip_address}
						class="input w-full"
						placeholder="e.g., 192.168.1.1"
					/>
				</div>

				<div>
					<label class="block text-sm font-medium text-text-secondary mb-2">Reason *</label>
					<textarea
						bind:value={newBan.reason}
						rows="3"
						class="input w-full"
						placeholder="Reason for ban..."
						required
					/>
				</div>

				<div>
					<label class="block text-sm font-medium text-text-secondary mb-2">Duration (days)</label>
					<input
						type="number"
						bind:value={newBan.duration}
						class="input w-full"
						placeholder="0 for permanent"
					/>
				</div>
			</div>

			<div class="flex items-center justify-end space-x-3 mt-6">
				<button onclick={() => (showCreateModal = false)} class="btn btn-outline"> Cancel </button>
				<button onclick={createBan} class="btn btn-primary"> Create Ban </button>
			</div>
		</div>
	</div>
{/if}
