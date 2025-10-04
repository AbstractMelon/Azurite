<script lang="ts">
	import { Package, Download, Heart, Calendar, User, ArrowRight, Shield } from 'lucide-svelte';
	import type { Mod } from '$lib/types';

	export let mod: Mod;
	export let variant: 'default' | 'compact' | 'hero' | 'list' = 'default';
	export let showStats = true;
	export let showDescription = true;
	export let showGame = true;
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
			day: 'numeric',
			year: 'numeric'
		});
	}

	$: cardClasses = `
		card relative overflow-hidden transition-all duration-300
		${interactive ? 'card-hover group cursor-pointer' : ''}
		${variant === 'hero' ? 'aspect-[16/9]' : variant === 'compact' ? 'h-24' : variant === 'list' ? 'h-auto' : 'aspect-[4/3]'}
	`.trim();
</script>

<a
	href="/mods/{mod.game?.slug || 'unknown'}/{mod.slug}"
	class={cardClasses}
	role="button"
	tabindex="0"
	aria-label="View {mod.name} mod details"
>
	{#if variant === 'hero'}
		<!-- Hero Card Layout - Full image background -->
		<div class="relative h-full">
			<!-- Large Background Image -->
			<div class="absolute inset-0">
				{#if mod.icon}
					<img src={mod.icon} alt={mod.name} class="w-full h-full object-cover" />
					<!-- Dark overlay for readability -->
					<div
						class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/40 to-black/20"
					></div>
				{:else}
					<div
						class="w-full h-full bg-gradient-to-br from-primary-600 via-primary-700 to-primary-800 flex items-center justify-center"
					>
						<Package class="w-20 h-20 text-white/40" />
					</div>
					<div
						class="absolute inset-0 bg-gradient-to-t from-black/60 via-black/30 to-transparent"
					></div>
				{/if}
			</div>

			<!-- Content positioned over image -->
			<div class="relative h-full flex flex-col justify-between p-6">
				<!-- Top section with game badge -->
				<div class="flex items-start justify-between">
					{#if showGame && mod.game}
						<div class="px-3 py-1 bg-black/50 backdrop-blur-sm rounded-full border border-white/20">
							<span class="text-white text-sm font-medium">{mod.game.name}</span>
						</div>
					{/if}

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
							{mod.name}
						</h3>
						{#if mod.version}
							<p class="text-white/70 text-sm">Version {mod.version}</p>
						{/if}
						{#if showDescription && mod.short_description}
							<p class="text-white/90 text-sm line-clamp-2 mt-2 drop-shadow">
								{mod.short_description}
							</p>
						{/if}
					</div>

					{#if showStats}
						<div class="flex items-center justify-between">
							<div class="flex items-center space-x-3 text-white/80 text-sm">
								<div class="flex items-center bg-black/30 backdrop-blur-sm rounded-full px-3 py-1">
									<Download class="w-4 h-4 mr-2" />
									<span class="font-medium">{formatNumber(mod.downloads || 0)}</span>
								</div>
								<div class="flex items-center bg-black/30 backdrop-blur-sm rounded-full px-3 py-1">
									<Heart class="w-4 h-4 mr-2" />
									<span class="font-medium">{formatNumber(mod.likes || 0)}</span>
								</div>
							</div>
							{#if mod.is_scanned}
								<div
									class="flex items-center px-3 py-1 bg-green-500/90 backdrop-blur-sm text-white text-xs font-medium rounded-full"
								>
									<Shield class="w-3 h-3 mr-1" />
									Verified
								</div>
							{/if}
						</div>
					{/if}
				</div>
			</div>
		</div>
	{:else if variant === 'compact'}
		<!-- Compact Card Layout - Horizontal with prominent image -->
		<div class="flex h-full">
			<!-- Large Image -->
			<div class="w-24 h-24 flex-shrink-0">
				{#if mod.icon}
					<img src={mod.icon} alt={mod.name} class="w-full h-full object-cover rounded-l-xl" />
				{:else}
					<div
						class="w-full h-full bg-gradient-to-br from-primary-600 to-primary-700 rounded-l-xl flex items-center justify-center"
					>
						<Package class="w-8 h-8 text-white" />
					</div>
				{/if}
			</div>

			<!-- Content -->
			<div class="flex-1 flex items-center justify-between px-4 py-2 min-w-0">
				<div class="min-w-0 flex-1">
					<h3
						class="font-semibold text-text-primary group-hover:text-primary-400 transition-colors truncate mb-1"
					>
						{mod.name}
					</h3>
					{#if showGame && mod.game}
						<p class="text-text-muted text-xs truncate mb-1">for {mod.game.name}</p>
					{/if}
					{#if showStats}
						<div class="flex items-center space-x-3 text-xs text-text-muted">
							<div class="flex items-center">
								<Download class="w-3 h-3 mr-1" />
								<span>{formatNumber(mod.downloads || 0)}</span>
							</div>
							<div class="flex items-center">
								<Heart class="w-3 h-3 mr-1" />
								<span>{formatNumber(mod.likes || 0)}</span>
							</div>
						</div>
					{/if}
				</div>

				<div class="flex items-center space-x-2 flex-shrink-0">
					{#if mod.is_scanned}
						<div class="w-2 h-2 bg-green-500 rounded-full"></div>
					{/if}
					{#if interactive}
						<ArrowRight
							class="w-4 h-4 text-text-muted group-hover:text-primary-400 transition-colors"
						/>
					{/if}
				</div>
			</div>
		</div>
	{:else if variant === 'list'}
		<!-- List Card Layout - Horizontal with medium image -->
		<div class="p-4 flex items-start space-x-4">
			<!-- Medium Image -->
			<div class="flex-shrink-0">
				{#if mod.icon}
					<img src={mod.icon} alt={mod.name} class="w-20 h-20 rounded-xl object-cover shadow-md" />
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
						<h3
							class="text-lg font-bold text-text-primary group-hover:text-primary-400 transition-colors truncate mb-1"
						>
							{mod.name}
						</h3>
						<div class="flex items-center space-x-2 mb-1">
							{#if showGame && mod.game}
								<span class="text-text-muted text-sm">for {mod.game.name}</span>
							{/if}
							{#if mod.version}
								<span class="px-2 py-0.5 bg-slate-700/50 text-text-muted text-xs rounded-full"
									>v{mod.version}</span
								>
							{/if}
						</div>
					</div>
					{#if interactive}
						<div class="opacity-0 group-hover:opacity-100 transition-opacity">
							<ArrowRight class="w-5 h-5 text-primary-400" />
						</div>
					{/if}
				</div>

				{#if showDescription && mod.short_description}
					<p class="text-text-secondary text-sm line-clamp-2 mb-3 leading-relaxed">
						{mod.short_description}
					</p>
				{/if}

				{#if showStats}
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
								<span>{formatDate(mod.updated_at)}</span>
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
				{/if}
			</div>
		</div>
	{:else}
		<!-- Default Card Layout - Image on top -->
		<div class="h-full flex flex-col">
			<!-- Large Image Section -->
			<div class="relative aspect-[16/10] overflow-hidden">
				{#if mod.icon}
					<img
						src={mod.icon}
						alt={mod.name}
						class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
					/>
				{:else}
					<div
						class="w-full h-full bg-gradient-to-br from-primary-600 via-primary-700 to-primary-800 flex items-center justify-center"
					>
						<Package class="w-16 h-16 text-white/50" />
					</div>
				{/if}

				<!-- Overlay elements -->
				<div
					class="absolute inset-0 bg-gradient-to-t from-black/20 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity"
				></div>

				<!-- Game badge -->
				{#if showGame && mod.game}
					<div class="absolute top-3 left-3">
						<div
							class="px-2 py-1 bg-black/70 backdrop-blur-sm text-white text-xs font-medium rounded-full"
						>
							{mod.game.name}
						</div>
					</div>
				{/if}

				<!-- Status badges in top right -->
				<div class="absolute top-3 right-3 flex items-center space-x-2">
					{#if mod.is_scanned}
						<div
							class="px-2 py-1 bg-green-500/90 backdrop-blur-sm text-white text-xs font-medium rounded-full"
						>
							<Shield class="w-3 h-3" />
						</div>
					{/if}
					{#if mod.version}
						<div
							class="px-2 py-1 bg-primary-500/90 backdrop-blur-sm text-white text-xs font-medium rounded-full"
						>
							v{mod.version}
						</div>
					{/if}
				</div>

				<!-- Arrow indicator -->
				{#if interactive}
					<div
						class="absolute bottom-3 right-3 opacity-0 group-hover:opacity-100 transition-opacity"
					>
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
						class="text-lg font-bold text-text-primary group-hover:text-primary-400 transition-colors line-clamp-2 mb-2"
					>
						{mod.name}
					</h3>
					{#if showDescription && mod.short_description}
						<p class="text-text-secondary text-sm line-clamp-2 leading-relaxed">
							{mod.short_description}
						</p>
					{/if}
				</div>

				<!-- Stats at bottom -->
				{#if showStats}
					<div class="mt-auto flex items-center justify-between text-sm">
						<div class="flex items-center space-x-3 text-text-muted">
							<div class="flex items-center">
								<Download class="w-4 h-4 mr-1" />
								<span class="font-medium">{formatNumber(mod.downloads || 0)}</span>
							</div>
							<div class="flex items-center">
								<Heart class="w-4 h-4 mr-1" />
								<span class="font-medium">{formatNumber(mod.likes || 0)}</span>
							</div>
						</div>
						{#if mod.owner}
							<div class="flex items-center text-text-muted text-xs truncate">
								<User class="w-3 h-3 mr-1 flex-shrink-0" />
								<span class="truncate">{mod.owner.display_name || mod.owner.username}</span>
							</div>
						{/if}
					</div>
				{/if}
			</div>
		</div>
	{/if}
</a>

<style>
	.line-clamp-2 {
		display: -webkit-box;
		-webkit-line-clamp: 2;
		line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
</style>
