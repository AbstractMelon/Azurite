<script lang="ts">
	export let title: string;
	export let value: string | number;
	export let subtitle: string = '';
	export let icon: any;
	export let iconColor: string = 'text-blue-400';
	export let iconBg: string = 'bg-blue-500/20';
</script>

<div class="card">
	<div class="p-6">
		<div class="flex items-center">
			<div class="p-2 {iconBg} rounded-lg">
				<svelte:component this={icon} class="w-6 h-6 {iconColor}" />
			</div>
			<div class="ml-4">
				<p class="text-text-muted text-sm">{title}</p>
				<p class="text-2xl font-bold text-text-primary">{value}</p>
			</div>
		</div>
		{#if subtitle}
			<div class="mt-4 text-sm text-text-muted">{subtitle}</div>
		{/if}
	</div>
</div>
