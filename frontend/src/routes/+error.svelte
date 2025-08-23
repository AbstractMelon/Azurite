<script lang="ts">
	export let data: { error?: Error & { status?: number }; status?: number };
	$: error = data?.error;
	$: status = data?.status || error?.status || 500;
</script>

<div class="flex flex-col items-center justify-center min-h-screen text-center px-6">
	<div class="glass rounded-2xl p-10 max-w-lg w-full fade-in">
		<h1 class="text-6xl font-extrabold text-gradient mb-4">
			{status || 500}
		</h1>
		<p class="text-2xl font-semibold text-text-primary mb-2">
			{error?.message || 'Something went wrong'}
		</p>
		<p class="text-text-secondary mb-6">
			{status === 404
				? 'We couldn’t find the page you’re looking for.'
				: 'An unexpected error has occurred. Please try again later.'}
		</p>
		<div class="flex gap-4 justify-center">
			<a href="/" class="btn btn-primary">Go Home</a>
			<button class="btn btn-secondary" on:click={() => location.reload()}> Retry </button>
		</div>
	</div>
</div>
