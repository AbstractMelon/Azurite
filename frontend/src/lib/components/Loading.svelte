<script lang="ts">
	export let size: 'sm' | 'md' | 'lg' | 'xl' = 'md';
	export let text: string | null = null;
	export let inline = false;
	export let variant: 'default' | 'primary' | 'white' = 'default';

	// Size classes
	const sizeClasses = {
		sm: 'w-4 h-4',
		md: 'w-6 h-6',
		lg: 'w-8 h-8',
		xl: 'w-12 h-12'
	};

	// Variant classes
	const variantClasses = {
		default: 'border-gray-300 border-t-primary-600',
		primary: 'border-primary-200 border-t-primary-600',
		white: 'border-white/30 border-t-white'
	};

	// Text size classes
	const textSizeClasses = {
		sm: 'text-xs',
		md: 'text-sm',
		lg: 'text-base',
		xl: 'text-lg'
	};
</script>

{#if inline}
	<div class="inline-flex items-center space-x-2">
		<div
			class="loading-spinner {sizeClasses[size]} border-2 {variantClasses[variant]}"
			role="status"
			aria-label="Loading"
		></div>
		{#if text}
			<span class="text-text-secondary {textSizeClasses[size]}">{text}</span>
		{/if}
	</div>
{:else}
	<div class="flex flex-col items-center justify-center space-y-3 p-8">
		<div
			class="loading-spinner {sizeClasses[size]} border-2 {variantClasses[variant]}"
			role="status"
			aria-label="Loading"
		></div>
		{#if text}
			<p class="text-text-secondary {textSizeClasses[size]} text-center">{text}</p>
		{/if}
	</div>
{/if}
