<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { auth } from '$lib/stores/auth';
	import { toast } from '$lib/stores/notifications';
	import { authApi } from '$lib/api/client';
	import type { AuthResponse } from '$lib/types';
	import Loading from '$lib/components/Loading.svelte';

	let isProcessing = true;
	let errorMessage = '';

	onMount(async () => {
		const code = $page.url.searchParams.get('code');
		const state = $page.url.searchParams.get('state');
		const error = $page.url.searchParams.get('error');

		// Handle OAuth error
		if (error) {
			errorMessage = 'Discord authentication was cancelled or failed.';
			isProcessing = false;
			toast.error('Authentication Failed', errorMessage);
			setTimeout(() => goto('/auth/login'), 3000);
			return;
		}

		// Validate required parameters
		if (!code || !state) {
			errorMessage = 'Invalid callback parameters. Please try again.';
			isProcessing = false;
			toast.error('Authentication Failed', errorMessage);
			setTimeout(() => goto('/auth/login'), 3000);
			return;
		}

		try {
			const response = await authApi.handleOAuthCallback('discord', code, state);

			if (response.success && response.data) {
				// Store the authentication data
				const authData = response.data as AuthResponse;
				auth.login(authData.user, authData.token);
				toast.success('Welcome!', 'You have been successfully authenticated with Discord.');

				// Redirect to home page
				goto('/');
			} else {
				errorMessage = response.error || 'Failed to authenticate with Discord.';
				isProcessing = false;
				toast.error('Authentication Failed', errorMessage);
				setTimeout(() => goto('/auth/login'), 3000);
			}
		} catch (error) {
			console.error('Discord OAuth callback error:', error);
			errorMessage = 'An unexpected error occurred. Please try again.';
			isProcessing = false;
			toast.error('Authentication Failed', errorMessage);
			setTimeout(() => goto('/auth/login'), 3000);
		}
	});
</script>

<svelte:head>
	<title>Discord Authentication - Azurite</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-background-primary">
	<div class="text-center">
		{#if isProcessing}
			<Loading size="lg" />
			<h2 class="mt-6 text-2xl font-bold text-text-primary">Authenticating with Discord...</h2>
			<p class="mt-2 text-text-secondary">Please wait while we complete your sign-in.</p>
		{:else}
			<div class="text-red-500 mb-4">
				<svg
					class="w-16 h-16 mx-auto"
					fill="none"
					stroke="currentColor"
					viewBox="0 0 24 24"
					xmlns="http://www.w3.org/2000/svg"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
					/>
				</svg>
			</div>
			<h2 class="text-2xl font-bold text-text-primary">Authentication Failed</h2>
			<p class="mt-2 text-text-secondary">{errorMessage}</p>
			<p class="mt-4 text-text-muted">Redirecting to login page...</p>
		{/if}
	</div>
</div>
