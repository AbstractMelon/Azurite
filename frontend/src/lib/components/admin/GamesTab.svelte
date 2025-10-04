<script lang="ts">
	import { goto } from '$app/navigation';
	import { adminApi } from '$lib/api/client';
	import { toast } from '$lib/stores/notifications';
	import {
		Gamepad2,
		Plus,
		Eye,
		Edit,
		Trash2,
		UserCog,
		Flag,
		CheckCircle,
		AlertTriangle
	} from 'lucide-svelte';
	import type { Game, GameRequest } from '$lib/types';

	export let allGames: Game[] = [];
	export let allGameRequests: GameRequest[] = [];
	export let onLoad: () => Promise<void>;

	async function deleteGame(gameId: number, gameName: string) {
		if (!confirm(`Are you sure you want to delete "${gameName}"?`)) return;

		try {
			const response = await adminApi.deleteGame(gameId);
			if (response.success) {
				toast.success('Game deleted', `"${gameName}" has been deleted.`);
				await onLoad();
			}
		} catch (error) {
			toast.error('Error', 'Failed to delete game');
		}
	}

	async function approveGameRequest(requestId: number, gameName: string) {
		try {
			const response = await adminApi.approveGameRequest(requestId);
			if (response.success) {
				toast.success('Game approved', `"${gameName}" has been added.`);
				await onLoad();
			}
		} catch (error) {
			toast.error('Error', 'Failed to approve game request');
		}
	}

	async function rejectGameRequest(requestId: number, gameName: string) {
		const reason = prompt('Reason for rejection:');
		if (!reason) return;

		try {
			const response = await adminApi.rejectGameRequest(requestId, reason);
			if (response.success) {
				toast.success('Request rejected', `"${gameName}" request has been rejected.`);
				await onLoad();
			}
		} catch (error) {
			toast.error('Error', 'Failed to reject game request');
		}
	}

	function formatRelativeTime(dateString: string): string {
		const date = new Date(dateString);
		const now = new Date();
		const diffTime = Math.abs(now.getTime() - date.getTime());
		const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));
		const diffHours = Math.floor(diffTime / (1000 * 60 * 60));
		const diffMinutes = Math.floor(diffTime / (1000 * 60));

		if (diffMinutes < 60) return `${diffMinutes}m ago`;
		if (diffHours < 24) return `${diffHours}h ago`;
		return `${diffDays}d ago`;
	}
</script>

<div class="space-y-6">
	<!-- Games List -->
	<div class="card">
		<div class="p-6">
			<div class="flex items-center justify-between mb-6">
				<h2 class="text-xl font-semibold text-text-primary flex items-center">
					<Gamepad2 class="w-5 h-5 mr-2" />
					Games ({allGames.length})
				</h2>
				<button class="btn btn-primary" onclick={() => goto('/admin/games/new')}>
					<Plus class="w-4 h-4 mr-2" />
					Add Game
				</button>
			</div>

			{#if allGames.length === 0}
				<div class="text-center py-12">
					<Gamepad2 class="w-16 h-16 mx-auto text-text-muted mb-4 opacity-50" />
					<h3 class="text-lg font-medium text-text-primary mb-2">No games yet</h3>
					<p class="text-text-secondary">Add your first game to get started.</p>
				</div>
			{:else}
				<div class="overflow-x-auto">
					<table class="w-full">
						<thead>
							<tr class="border-b border-slate-700">
								<th class="text-left py-3 px-4 text-text-muted font-medium">Game</th>
								<th class="text-left py-3 px-4 text-text-muted font-medium">Slug</th>
								<th class="text-left py-3 px-4 text-text-muted font-medium">Mods</th>
								<th class="text-left py-3 px-4 text-text-muted font-medium">Status</th>
								<th class="text-right py-3 px-4 text-text-muted font-medium">Actions</th>
							</tr>
						</thead>
						<tbody>
							{#each allGames as game (game.id)}
								<tr class="border-b border-slate-800 hover:bg-slate-800/50">
									<td class="py-3 px-4">
										<div class="flex items-center space-x-3">
											{#if game.icon}
												<img
													src={game.icon}
													alt={game.name}
													class="w-10 h-10 rounded object-cover"
												/>
											{:else}
												<div
													class="w-10 h-10 rounded bg-slate-700 flex items-center justify-center"
												>
													<Gamepad2 class="w-5 h-5 text-text-muted" />
												</div>
											{/if}
											<div>
												<p class="font-medium text-text-primary">{game.name}</p>
												<p class="text-xs text-text-muted">ID: {game.id}</p>
											</div>
										</div>
									</td>
									<td class="py-3 px-4 text-text-secondary">{game.slug}</td>
									<td class="py-3 px-4 text-text-secondary">{game.mod_count || 0}</td>
									<td class="py-3 px-4">
										{#if game.is_active}
											<span
												class="px-2 py-1 text-xs font-medium bg-green-500/20 text-green-400 rounded"
											>
												Active
											</span>
										{:else}
											<span
												class="px-2 py-1 text-xs font-medium bg-red-500/20 text-red-400 rounded"
											>
												Inactive
											</span>
										{/if}
									</td>
									<td class="py-3 px-4">
										<div class="flex items-center justify-end space-x-2">
											<button
												onclick={() => goto(`/games/${game.slug}`)}
												class="btn btn-sm btn-outline"
												title="View"
											>
												<Eye class="w-4 h-4" />
											</button>
											<button
												onclick={() => goto(`/admin/games/${game.id}/edit`)}
												class="btn btn-sm btn-outline"
												title="Edit"
											>
												<Edit class="w-4 h-4" />
											</button>
											<button
												onclick={() => goto(`/admin/games/${game.slug}/moderators`)}
												class="btn btn-sm btn-outline"
												title="Manage Moderators"
											>
												<UserCog class="w-4 h-4" />
											</button>
											<button
												onclick={() => deleteGame(game.id, game.name)}
												class="btn btn-sm btn-outline text-red-400 hover:text-red-300"
												title="Delete"
											>
												<Trash2 class="w-4 h-4" />
											</button>
										</div>
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			{/if}
		</div>
	</div>

	<!-- Game Requests -->
	<div class="card">
		<div class="p-6">
			<h2 class="text-xl font-semibold text-text-primary mb-6 flex items-center">
				<Flag class="w-5 h-5 mr-2" />
				Game Requests ({allGameRequests.filter((r) => r.status === 'pending').length} pending)
			</h2>

			{#if allGameRequests.length === 0}
				<div class="text-center py-12">
					<Flag class="w-16 h-16 mx-auto text-text-muted mb-4 opacity-50" />
					<h3 class="text-lg font-medium text-text-primary mb-2">No game requests</h3>
					<p class="text-text-secondary">Game requests will appear here.</p>
				</div>
			{:else}
				<div class="space-y-3">
					{#each allGameRequests as request (request.id)}
						<div
							class="flex items-center justify-between p-4 bg-slate-800 rounded-lg hover:bg-slate-800/80"
						>
							<div class="flex-1">
								<div class="flex items-center space-x-3">
									<h3 class="font-medium text-text-primary">{request.name}</h3>
									{#if request.status === 'pending'}
										<span
											class="px-2 py-1 text-xs font-medium bg-yellow-500/20 text-yellow-400 rounded"
										>
											Pending
										</span>
									{:else if request.status === 'approved'}
										<span
											class="px-2 py-1 text-xs font-medium bg-green-500/20 text-green-400 rounded"
										>
											Approved
										</span>
									{:else if request.status === 'denied'}
										<span class="px-2 py-1 text-xs font-medium bg-red-500/20 text-red-400 rounded">
											Denied
										</span>
									{/if}
								</div>
								<p class="text-sm text-text-muted mt-1">
									Requested by {request.requested_by?.display_name || 'Unknown'} •
									{formatRelativeTime(request.created_at)}
								</p>
								{#if request.reason}
									<p class="text-sm text-text-secondary mt-2">
										<strong>Reason:</strong>
										{request.reason}
									</p>
								{/if}
								{#if request.description}
									<p class="text-sm text-text-secondary mt-1">
										<strong>Description:</strong>
										{request.description.length > 150
											? request.description.substring(0, 150) + '...'
											: request.description}
									</p>
								{/if}
								<div class="flex flex-wrap gap-3 mt-2 text-xs text-text-muted">
									{#if request.existing_community}
										<span>
											<strong>Community:</strong>
											{request.existing_community}
										</span>
									{/if}
									{#if request.mod_loader}
										<span>
											<strong>Mod Loader:</strong>
											{request.mod_loader}
										</span>
									{/if}
									{#if request.contact}
										<span>
											<strong>Contact:</strong>
											{request.contact}
										</span>
									{/if}
								</div>
								{#if request.admin_notes}
									<p class="text-sm text-yellow-400 mt-2">
										<strong>Admin Notes:</strong>
										{request.admin_notes}
									</p>
								{/if}
							</div>
							<div class="flex items-center space-x-2">
								{#if request.status === 'pending'}
									<button
										onclick={() => approveGameRequest(request.id, request.name)}
										class="btn btn-sm bg-green-600 text-white hover:bg-green-700"
										title="Approve"
									>
										<CheckCircle class="w-4 h-4" />
									</button>
									<button
										onclick={() => rejectGameRequest(request.id, request.name)}
										class="btn btn-sm bg-red-600 text-white hover:bg-red-700"
										title="Reject"
									>
										<AlertTriangle class="w-4 h-4" />
									</button>
								{/if}
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</div>
	</div>
</div>
