<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { auth, isAuthenticated } from '$lib/stores/auth';
	import { toast } from '$lib/stores/notifications';
	import { authApi } from '$lib/api/client';
	import type { AuthResponse } from '$lib/types';
	import Loading from '$lib/components/Loading.svelte';
	import { Eye, EyeOff, Mail, Lock, Github, Chrome } from 'lucide-svelte';

	let email = '';
	let password = '';
	let showPassword = false;
	let isLoading = false;
	let errors: { [key: string]: string } = {};

	// Redirect if already authenticated
	onMount(() => {
		if ($isAuthenticated) {
			const redirectTo = $page.url.searchParams.get('redirect') || '/';
			goto(redirectTo);
		}
	});

	// Handle form submission
	async function handleLogin(event: Event) {
		event.preventDefault();

		// Clear previous errors
		errors = {};

		// Basic validation
		if (!email.trim()) {
			errors.email = 'Email is required';
		} else if (!isValidEmail(email)) {
			errors.email = 'Please enter a valid email address';
		}

		if (!password.trim()) {
			errors.password = 'Password is required';
		} else if (password.length < 6) {
			errors.password = 'Password must be at least 6 characters';
		}

		if (Object.keys(errors).length > 0) {
			return;
		}

		isLoading = true;
		auth.setLoading(true);

		try {
			const response = await authApi.login({ email: email.trim(), password });

			if (response.success && response.data) {
				const authData = response.data as AuthResponse;
				auth.login(authData.user, authData.token);
				toast.success('Welcome back!', 'You have been successfully logged in.');

				// Redirect to the intended page or homepage
				const redirectTo = $page.url.searchParams.get('redirect') || '/';
				goto(redirectTo);
			} else {
				toast.error('Login failed', response.error || 'Invalid email or password');
			}
		} catch (error) {
			console.error('Login error:', error);
			toast.error('Login failed', 'An unexpected error occurred. Please try again.');
		} finally {
			isLoading = false;
			auth.setLoading(false);
		}
	}

	// Email validation
	function isValidEmail(email: string): boolean {
		const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
		return emailRegex.test(email);
	}

	// Handle OAuth login
	async function handleOAuthLogin(provider: string) {
		try {
			const response = await authApi.getOAuthURL(provider);

			if (response.success && response.data) {
				const oauthData = response.data as { url: string };
				if (oauthData.url) {
					// Redirect to the OAuth provider's authorization page
					window.location.href = oauthData.url;
				}
			} else {
				toast.error('OAuth Error', response.error || `${provider} OAuth is not configured`);
			}
		} catch (error) {
			console.error(`${provider} OAuth error:`, error);
			toast.error('OAuth Error', 'Failed to initiate OAuth flow. Please try again.');
		}
	}

	// Toggle password visibility
	function togglePasswordVisibility() {
		showPassword = !showPassword;
	}

	// Handle forgot password
	function handleForgotPassword() {
		goto('/auth/forgot-password');
	}
</script>

<svelte:head>
	<title>Login - Azurite</title>
	<meta
		name="description"
		content="Sign in to your Azurite account to manage mods, leave comments, and connect with the community."
	/>
</svelte:head>

<div
	class="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8 bg-background-primary"
>
	<div class="max-w-md w-full space-y-8">
		<!-- Header -->
		<div class="text-center">
			<a href="/" class="inline-flex items-center space-x-2 mb-6">
				<div
					class="w-10 h-10 bg-gradient-to-br from-primary-400 to-primary-600 rounded-lg flex items-center justify-center shadow-glow-blue"
				>
					<span class="text-white font-bold text-xl">A</span>
				</div>
				<span class="text-2xl font-bold text-gradient">Azurite</span>
			</a>

			<h2 class="text-3xl font-bold text-text-primary">Welcome back</h2>
			<p class="mt-2 text-text-secondary">Sign in to your account to continue</p>
		</div>

		<!-- Login Form -->
		<div class="card">
			<div class="p-6 sm:p-8">
				<form on:submit={handleLogin} class="space-y-6">
					<!-- Email Field -->
					<div>
						<label for="email" class="block text-sm font-medium text-text-primary mb-2">
							Email address
						</label>
						<div class="relative">
							<div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
								<Mail class="h-5 w-5 text-text-muted" />
							</div>
							<input
								id="email"
								name="email"
								type="email"
								autocomplete="email"
								required
								bind:value={email}
								class="input pl-10 {errors.email ? 'input-error' : ''}"
								placeholder="Enter your email"
								disabled={isLoading}
							/>
						</div>
						{#if errors.email}
							<p class="error-message">{errors.email}</p>
						{/if}
					</div>

					<!-- Password Field -->
					<div>
						<label for="password" class="block text-sm font-medium text-text-primary mb-2">
							Password
						</label>
						<div class="relative">
							<div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
								<Lock class="h-5 w-5 text-text-muted" />
							</div>
							<input
								id="password"
								name="password"
								type={showPassword ? 'text' : 'password'}
								autocomplete="current-password"
								required
								bind:value={password}
								class="input pl-10 pr-10 {errors.password ? 'input-error' : ''}"
								placeholder="Enter your password"
								disabled={isLoading}
							/>
							<div class="absolute inset-y-0 right-0 pr-3 flex items-center">
								<button
									type="button"
									class="text-text-muted hover:text-text-primary transition-colors"
									onclick={togglePasswordVisibility}
									disabled={isLoading}
								>
									{#if showPassword}
										<EyeOff class="h-5 w-5" />
									{:else}
										<Eye class="h-5 w-5" />
									{/if}
								</button>
							</div>
						</div>
						{#if errors.password}
							<p class="error-message">{errors.password}</p>
						{/if}
					</div>

					<!-- Remember Me & Forgot Password -->
					<div class="flex items-center justify-between">
						<div class="flex items-center">
							<input
								id="remember-me"
								name="remember-me"
								type="checkbox"
								class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-slate-600 bg-slate-800 rounded"
								disabled={isLoading}
							/>
							<label for="remember-me" class="ml-2 block text-sm text-text-secondary">
								Remember me
							</label>
						</div>

						<div class="text-sm">
							<button
								type="button"
								onclick={handleForgotPassword}
								class="text-primary-400 hover:text-primary-300 transition-colors"
								disabled={isLoading}
							>
								Forgot password?
							</button>
						</div>
					</div>

					<!-- Submit Button -->
					<div>
						<button
							type="submit"
							disabled={isLoading}
							class="btn btn-primary w-full flex justify-center"
						>
							{#if isLoading}
								<Loading size="sm" inline />
							{:else}
								Sign in
							{/if}
						</button>
					</div>
				</form>

				<!-- Divider -->
				<div class="mt-6">
					<div class="relative">
						<div class="absolute inset-0 flex items-center">
							<div class="w-full border-t border-slate-600"></div>
						</div>
						<div class="relative flex justify-center text-sm">
							<span class="px-2 bg-background-secondary text-text-muted">Or continue with</span>
						</div>
					</div>
				</div>

				<!-- OAuth Buttons -->
				<div class="mt-6 grid grid-cols-1 gap-3">
					<!-- GitHub -->
					<button
						type="button"
						onclick={() => handleOAuthLogin('github')}
						disabled={isLoading}
						class="btn btn-outline w-full flex items-center justify-center"
					>
						<Github class="w-5 h-5 mr-2" />
						Sign in with GitHub
					</button>

					<!-- Google -->
					<button
						type="button"
						onclick={() => handleOAuthLogin('google')}
						disabled={isLoading}
						class="btn btn-outline w-full flex items-center justify-center"
					>
						<Chrome class="w-5 h-5 mr-2" />
						Sign in with Google
					</button>

					<!-- Discord -->
					<button
						type="button"
						onclick={() => handleOAuthLogin('discord')}
						disabled={isLoading}
						class="btn btn-outline w-full flex items-center justify-center"
					>
						<svg class="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 24 24">
							<path
								d="M20.317 4.37a19.791 19.791 0 0 0-4.885-1.515.074.074 0 0 0-.079.037c-.21.375-.444.864-.608 1.25a18.27 18.27 0 0 0-5.487 0 12.64 12.64 0 0 0-.617-1.25.077.077 0 0 0-.079-.037A19.736 19.736 0 0 0 3.677 4.37a.07.07 0 0 0-.032.027C.533 9.046-.32 13.58.099 18.057a.082.082 0 0 0 .031.057 19.9 19.9 0 0 0 5.993 3.03.078.078 0 0 0 .084-.028c.462-.63.874-1.295 1.226-1.994.021-.041.001-.09-.041-.106a13.107 13.107 0 0 1-1.872-.892.077.077 0 0 1-.008-.128 10.2 10.2 0 0 0 .372-.292.074.074 0 0 1 .077-.010c3.928 1.793 8.18 1.793 12.062 0a.074.074 0 0 1 .078.01c.12.098.246.198.373.292a.077.077 0 0 1-.006.127 12.299 12.299 0 0 1-1.873.892.077.077 0 0 0-.041.107c.36.698.772 1.362 1.225 1.993a.076.076 0 0 0 .084.028 19.839 19.839 0 0 0 6.002-3.03.077.077 0 0 0 .032-.054c.5-5.177-.838-9.674-3.549-13.66a.061.061 0 0 0-.031-.03zM8.02 15.33c-1.183 0-2.157-1.085-2.157-2.419 0-1.333.956-2.419 2.157-2.419 1.21 0 2.176 1.096 2.157 2.42 0 1.333-.956 2.418-2.157 2.418zm7.975 0c-1.183 0-2.157-1.085-2.157-2.419 0-1.333.955-2.419 2.157-2.419 1.21 0 2.176 1.096 2.157 2.42 0 1.333-.946 2.418-2.157 2.418z"
							/>
						</svg>
						Sign in with Discord
					</button>
				</div>
			</div>
		</div>

		<!-- Sign Up Link -->
		<div class="text-center">
			<p class="text-text-secondary">
				Don't have an account?
				<a
					href="/auth/register"
					class="text-primary-400 hover:text-primary-300 font-medium transition-colors"
				>
					Sign up for free
				</a>
			</p>
		</div>
	</div>
</div>
