<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { user, isAuthenticated } from '$lib/stores/auth';

	onMount(() => {
		if (!$isAuthenticated) {
			goto('/auth/login?redirect=/profile');
			return;
		}

		const username = $user?.username;
		if (username) {
			goto(`/profile/${username}`);
		} else {
			goto('/');
		}
	});
</script>