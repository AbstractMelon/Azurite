<script lang="ts">
	import { Gamepad2, Package, ArrowRight, Calendar } from 'lucide-svelte';
	import type { Game } from '$lib/types';

	export let game: Game;
	export let variant: 'default' | 'compact' | 'hero' = 'default';
	export let showStats = true;
	export let showDescription = true;
	export let interactive = true;

	function formatNumber(num: number): string {
		if (num >= 1000000) {
			return (num / 1000000).toFixed(1) + 'M';
		}
		if (num >= 1000) {
			return (num / 1000).toFixed(1) + 'K';
		}
		return num.toString();
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			month: 'short',
			year: 'numeric'
		});
	}

	$: cardClasses = `
		card relative overflow-hidden transition-all duration-300
		${interactive ? 'card-hover group cursor-pointer' : ''}
		${variant === 'hero' ? 'aspect-[16/9]' : variant === 'compact' ? 'h-20' : 'aspect-[4/3]'}
	`.trim();
</script>

<a
	href="/games/{game.slug}"
	class={cardClasses}
	role="button"
	tabindex="0"
	aria-label="View {game.name} details"
>
	{#if variant === 'hero'}
		<!-- Hero Card Layout - Full image background -->
		<div class="relative h-full">
			<!-- Large Background Image -->
			<div class="absolute inset-0">
				{#if game.icon}
					<img src={game.icon} alt={game.name} class="w-full h-full object-cover" />
					<!-- Dark overlay for readability -->
					<div
						class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/40 to-black/20"
					></div>
				{:else}
					<div
						class="w-full h-full bg-gradient-to-br from-slate-700 via-slate-600 to-slate-800 flex items-center justify-center"
					>
						<Gamepad2 class="w-20 h-20 text-white/40" />
					</div>
					<div
						class="absolute inset-0 bg-gradient-to-t from-black/60 via-black/30 to-transparent"
					></div>
				{/if}
			</div>

			<!-- Content positioned over image -->
			<div class="relative h-full flex flex-col justify-between p-6">
				<!-- Top corner - Status indicator -->
				<div class="flex justify-end">
					{#if interactive}
						<div
							class="opacity-0 group-hover:opacity-100 transition-opacity bg-black/30 backdrop-blur-sm rounded-full p-2"
						>
							<ArrowRight class="w-5 h-5 text-white" />
						</div>
					{/if}
				</div>

				<!-- Bottom content -->
				<div class="space-y-3">
					<div>
						<h3
							class="text-2xl font-bold text-white drop-shadow-lg group-hover:text-primary-300 transition-colors"
						>
							{game.name}
						</h3>
						{#if showDescription && game.description}
							<p class="text-white/90 text-sm line-clamp-2 mt-2 drop-shadow">
								{game.description}
							</p>
						{/if}
					</div>

					{#if showStats}
						<div class="flex items-center justify-between">
							<div class="flex items-center space-x-4 text-white/80 text-sm">
								<div class="flex items-center bg-black/30 backdrop-blur-sm rounded-full px-3 py-1">
									<Package class="w-4 h-4 mr-2" />
									<span class="font-medium">{formatNumber(game.mod_count || 0)} mods</span>
								</div>
							</div>
							<div class="flex items-center">
								{#if game.is_active}
									<div
										class="px-3 py-1 bg-green-500/90 backdrop-blur-sm text-white text-xs font-medium rounded-full"
									>
										Active
									</div>
								{:else}
									<div
										class="px-3 py-1 bg-yellow-500/90 backdrop-blur-sm text-white text-xs font-medium rounded-full"
									>
										Coming Soon
									</div>
								{/if}
							</div>
						</div>
					{/if}
				</div>
			</div>
		</div>
	{:else if variant === 'compact'}
		<!-- Compact Card Layout - Horizontal with prominent image -->
		<div class="flex h-full">
			<!-- Large Image -->
			<div class="w-20 h-20 flex-shrink-0">
				{#if game.icon}
					<img src={game.icon} alt={game.name} class="w-full h-full object-cover rounded-l-xl" />
				{:else}
					<div
						class="w-full h-full bg-gradient-to-br from-slate-600 to-slate-700 rounded-l-xl flex items-center justify-center"
					>
						<Gamepad2 class="w-8 h-8 text-white/70" />
					</div>
				{/if}
			</div>

			<!-- Content -->
			<div class="flex-1 flex items-center justify-between px-4 min-w-0">
				<div class="min-w-0 flex-1">
					<h3
						class="font-semibold text-text-primary group-hover:text-primary-400 transition-colors truncate"
					>
						{game.name}
					</h3>
					{#if showStats}
						<p class="text-text-muted text-sm">
							{formatNumber(game.mod_count || 0)} mods
						</p>
					{/if}
				</div>

				<div class="flex items-center space-x-2 flex-shrink-0">
					{#if game.is_active}
						<div class="w-2 h-2 bg-green-500 rounded-full"></div>
					{:else}
						<div class="w-2 h-2 bg-yellow-500 rounded-full"></div>
					{/if}
					{#if interactive}
						<ArrowRight
							class="w-4 h-4 text-text-muted group-hover:text-primary-400 transition-colors"
						/>
					{/if}
				</div>
			</div>
		</div>
	{:else}
		<!-- Default Card Layout - Image on top -->
		<div class="h-full flex flex-col">
			<!-- Large Image Section -->
			<div class="relative aspect-[16/10] overflow-hidden">
				{#if game.icon}
					<img
						src={game.icon}
						alt={game.name}
						class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
					/>
				{:else}
					<div
						class="w-full h-full bg-gradient-to-br from-slate-600 via-slate-700 to-slate-800 flex items-center justify-center"
					>
						<Gamepad2 class="w-16 h-16 text-white/50" />
					</div>
				{/if}

				<!-- Overlay elements -->
				<div
					class="absolute inset-0 bg-gradient-to-t from-black/30 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity"
				></div>

				<!-- Status badge in top right -->
				<div class="absolute top-3 right-3">
					{#if game.is_active}
						<div
							class="px-2 py-1 bg-green-500/90 backdrop-blur-sm text-white text-xs font-medium rounded-full"
						>
							Active
						</div>
					{:else}
						<div
							class="px-2 py-1 bg-yellow-500/90 backdrop-blur-sm text-white text-xs font-medium rounded-full"
						>
							Coming Soon
						</div>
					{/if}
				</div>

				<!-- Arrow indicator -->
				{#if interactive}
					<div class="absolute top-3 left-3 opacity-0 group-hover:opacity-100 transition-opacity">
						<div class="bg-black/50 backdrop-blur-sm rounded-full p-2">
							<ArrowRight class="w-4 h-4 text-white" />
						</div>
					</div>
				{/if}
			</div>

			<!-- Compact Content Section -->
			<div class="flex-1 p-4 flex flex-col">
				<div class="mb-3">
					<h3
						class="text-lg font-bold text-text-primary group-hover:text-primary-400 transition-colors line-clamp-1 mb-1"
					>
						{game.name}
					</h3>
					{#if showDescription && game.description}
						<p class="text-text-secondary text-sm line-clamp-2 leading-relaxed">
							{game.description}
						</p>
					{/if}
				</div>

				<!-- Stats at bottom -->
				{#if showStats}
					<div class="mt-auto flex items-center justify-between text-sm">
						<div class="flex items-center space-x-3 text-text-muted">
							<div class="flex items-center">
								<Package class="w-4 h-4 mr-1" />
								<span class="font-medium">{formatNumber(game.mod_count || 0)}</span>
							</div>
							<div class="flex items-center">
								<Calendar class="w-4 h-4 mr-1" />
								<span>{formatDate(game.created_at)}</span>
							</div>
						</div>
					</div>
				{/if}
			</div>
		</div>
	{/if}
</a>

<style>
	.line-clamp-1 {
		display: -webkit-box;
		-webkit-line-clamp: 1;
		line-clamp: 1;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}

	.line-clamp-2 {
		display: -webkit-box;
		-webkit-line-clamp: 2;
		line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
</style>
