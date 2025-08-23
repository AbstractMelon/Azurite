<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import Header from '$lib/components/Header.svelte';
	import Footer from '$lib/components/Footer.svelte';
	import Toast from '$lib/components/Toast.svelte';
	import { auth } from '$lib/stores/auth';
	import { authApi } from '$lib/api/client';

	// Initialize auth state from localStorage
	onMount(async () => {
		// Try to restore user session
		const storedAuth = localStorage.getItem('azurite_auth');
		if (storedAuth) {
			try {
				const authData = JSON.parse(storedAuth);
				if (authData.token) {
					// Verify token is still valid by getting profile
					const response = await authApi.getProfile();
					if (response.success && response.data) {
						auth.login(response.data, authData.token);
					} else {
						// Token is invalid, clear stored auth
						auth.logout();
					}
				}
			} catch (error) {
				console.error('Failed to restore auth session:', error);
				auth.logout();
			}
		}
	});

	// Pages that should not show header/footer
	const noLayoutPages = ['/auth/login', '/auth/register', '/auth/reset-password'];

	$: showLayout = !noLayoutPages.some((path) => $page.url.pathname.startsWith(path));
</script>

<svelte:head>
	<meta name="theme-color" content="#1e3a8a" />
</svelte:head>

<div class="min-h-screen flex flex-col bg-background-primary text-text-primary">
	{#if showLayout}
		<Header />
	{/if}

	<main class="flex-1 w-full">
		<slot />
	</main>

	{#if showLayout}
		<Footer />
	{/if}
</div>

<!-- Toast Notifications -->
<Toast />

<style>
	:global(body) {
		margin: 0;
		padding: 0;
	}
</style>
