<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { auth, isAuthenticated } from '$lib/stores/auth';
	import { toast } from '$lib/stores/notifications';
	import { authApi } from '$lib/api/client';
	import type { AuthResponse } from '$lib/types';
	import Loading from '$lib/components/Loading.svelte';
	import { Eye, EyeOff, Mail, Lock, User, Github, Chrome } from 'lucide-svelte';

	let formData = {
		username: '',
		email: '',
		password: '',
		confirmPassword: '',
		display_name: ''
	};
	let showPassword = false;
	let showConfirmPassword = false;
	let isLoading = false;
	let errors: { [key: string]: string } = {};
	let acceptTerms = false;

	// Redirect if already authenticated
	onMount(() => {
		if ($isAuthenticated) {
			const redirectTo = $page.url.searchParams.get('redirect') || '/';
			goto(redirectTo);
		}
	});

	// Handle form submission
	async function handleRegister(event: Event) {
		event.preventDefault();

		// Clear previous errors
		errors = {};

		// Validation
		if (!formData.username.trim()) {
			errors.username = 'Username is required';
		} else if (formData.username.length < 3) {
			errors.username = 'Username must be at least 3 characters';
		} else if (!/^[a-zA-Z0-9_-]+$/.test(formData.username)) {
			errors.username = 'Username can only contain letters, numbers, hyphens, and underscores';
		}

		if (!formData.email.trim()) {
			errors.email = 'Email is required';
		} else if (!isValidEmail(formData.email)) {
			errors.email = 'Please enter a valid email address';
		}

		if (!formData.display_name.trim()) {
			errors.display_name = 'Display name is required';
		} else if (formData.display_name.length < 2) {
			errors.display_name = 'Display name must be at least 2 characters';
		}

		if (!formData.password.trim()) {
			errors.password = 'Password is required';
		} else if (formData.password.length < 8) {
			errors.password = 'Password must be at least 8 characters';
		} else if (!isStrongPassword(formData.password)) {
			errors.password =
				'Password must contain at least one uppercase letter, one lowercase letter, one number, and one special character';
		}

		if (!formData.confirmPassword.trim()) {
			errors.confirmPassword = 'Please confirm your password';
		} else if (formData.password !== formData.confirmPassword) {
			errors.confirmPassword = 'Passwords do not match';
		}

		if (!acceptTerms) {
			errors.terms = 'You must accept the Terms of Service and Privacy Policy';
		}

		if (Object.keys(errors).length > 0) {
			return;
		}

		isLoading = true;
		auth.setLoading(true);

		try {
			const response = await authApi.register({
				username: formData.username.trim().toLowerCase(),
				email: formData.email.trim().toLowerCase(),
				password: formData.password,
				display_name: formData.display_name.trim()
			});

			if (response.success && response.data) {
				const authData = response.data as AuthResponse;
				auth.login(authData.user, authData.token);
				toast.success('Welcome to Azurite!', 'Your account has been created successfully.');

				// Redirect to the intended page or homepage
				const redirectTo = $page.url.searchParams.get('redirect') || '/';
				goto(redirectTo);
			} else {
				toast.error('Registration failed', response.error || 'Unable to create account');
			}
		} catch (error) {
			console.error('Registration error:', error);
			toast.error('Registration failed', 'An unexpected error occurred. Please try again.');
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

	// Strong password validation
	function isStrongPassword(password: string): boolean {
		const hasUpperCase = /[A-Z]/.test(password);
		const hasLowerCase = /[a-z]/.test(password);
		const hasNumbers = /\d/.test(password);
		const hasSpecialChar = /[!@#$%^&*(),.?":{}|<>]/.test(password);
		return hasUpperCase && hasLowerCase && hasNumbers && hasSpecialChar;
	}

	// Handle OAuth registration
	async function handleOAuthRegister(provider: string) {
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

	function toggleConfirmPasswordVisibility() {
		showConfirmPassword = !showConfirmPassword;
	}

	// Auto-generate display name from username
	function generateDisplayName() {
		if (formData.username && !formData.display_name) {
			formData.display_name =
				formData.username.charAt(0).toUpperCase() + formData.username.slice(1);
		}
	}

	// Real-time validation feedback
	function validateField(field: string) {
		switch (field) {
			case 'username':
				if (
					formData.username &&
					formData.username.length >= 3 &&
					/^[a-zA-Z0-9_-]+$/.test(formData.username)
				) {
					delete errors.username;
					errors = errors;
				}
				break;
			case 'email':
				if (formData.email && isValidEmail(formData.email)) {
					delete errors.email;
					errors = errors;
				}
				break;
			case 'password':
				if (formData.password && formData.password.length >= 8) {
					delete errors.password;
					errors = errors;
				}
				// Also check confirm password
				if (formData.confirmPassword && formData.password === formData.confirmPassword) {
					delete errors.confirmPassword;
					errors = errors;
				}
				break;
			case 'confirmPassword':
				if (formData.confirmPassword && formData.password === formData.confirmPassword) {
					delete errors.confirmPassword;
					errors = errors;
				}
				break;
		}
	}
</script>

<svelte:head>
	<title>Sign Up - Azurite</title>
	<meta
		name="description"
		content="Create your Azurite account to start sharing mods, joining discussions, and connecting with the gaming community."
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

			<h2 class="text-3xl font-bold text-text-primary">Create your account</h2>
			<p class="mt-2 text-text-secondary">Join the community and start sharing your creations</p>
		</div>

		<!-- Registration Form -->
		<div class="card">
			<div class="p-6 sm:p-8">
				<form on:submit={handleRegister} class="space-y-6">
					<!-- Username Field -->
					<div>
						<label for="username" class="block text-sm font-medium text-text-primary mb-2">
							Username
							<span class="text-red-400">*</span>
						</label>
						<div class="relative">
							<div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
								<User class="h-5 w-5 text-text-muted" />
							</div>
							<input
								id="username"
								name="username"
								type="text"
								autocomplete="username"
								required
								bind:value={formData.username}
								on:input={() => {
									validateField('username');
									generateDisplayName();
								}}
								class="input pl-10 {errors.username ? 'input-error' : ''}"
								placeholder="Choose a unique username"
								disabled={isLoading}
							/>
						</div>
						{#if errors.username}
							<p class="error-message">{errors.username}</p>
						{:else}
							<p class="text-xs text-text-muted mt-1">
								Only letters, numbers, hyphens, and underscores allowed
							</p>
						{/if}
					</div>

					<!-- Display Name Field -->
					<div>
						<label for="display_name" class="block text-sm font-medium text-text-primary mb-2">
							Display Name
							<span class="text-red-400">*</span>
						</label>
						<input
							id="display_name"
							name="display_name"
							type="text"
							required
							bind:value={formData.display_name}
							on:input={() => validateField('display_name')}
							class="input {errors.display_name ? 'input-error' : ''}"
							placeholder="How others will see you"
							disabled={isLoading}
						/>
						{#if errors.display_name}
							<p class="error-message">{errors.display_name}</p>
						{:else}
							<p class="text-xs text-text-muted mt-1">
								This is how your name will appear to others
							</p>
						{/if}
					</div>

					<!-- Email Field -->
					<div>
						<label for="email" class="block text-sm font-medium text-text-primary mb-2">
							Email address
							<span class="text-red-400">*</span>
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
								bind:value={formData.email}
								on:input={() => validateField('email')}
								class="input pl-10 {errors.email ? 'input-error' : ''}"
								placeholder="Enter your email"
								disabled={isLoading}
							/>
						</div>
						{#if errors.email}
							<p class="error-message">{errors.email}</p>
						{:else}
							<p class="text-xs text-text-muted mt-1">
								We'll use this to send you important updates
							</p>
						{/if}
					</div>

					<!-- Password Field -->
					<div>
						<label for="password" class="block text-sm font-medium text-text-primary mb-2">
							Password
							<span class="text-red-400">*</span>
						</label>
						<div class="relative">
							<div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
								<Lock class="h-5 w-5 text-text-muted" />
							</div>
							<input
								id="password"
								name="password"
								type={showPassword ? 'text' : 'password'}
								autocomplete="new-password"
								required
								bind:value={formData.password}
								on:input={() => validateField('password')}
								class="input pl-10 pr-10 {errors.password ? 'input-error' : ''}"
								placeholder="Create a strong password"
								disabled={isLoading}
							/>
							<div class="absolute inset-y-0 right-0 pr-3 flex items-center">
								<button
									type="button"
									class="text-text-muted hover:text-text-primary transition-colors"
									on:click={togglePasswordVisibility}
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
						{:else}
							<p class="text-xs text-text-muted mt-1">
								At least 8 characters with uppercase, lowercase, number, and special character
							</p>
						{/if}
					</div>

					<!-- Confirm Password Field -->
					<div>
						<label for="confirmPassword" class="block text-sm font-medium text-text-primary mb-2">
							Confirm Password
							<span class="text-red-400">*</span>
						</label>
						<div class="relative">
							<div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
								<Lock class="h-5 w-5 text-text-muted" />
							</div>
							<input
								id="confirmPassword"
								name="confirmPassword"
								type={showConfirmPassword ? 'text' : 'password'}
								autocomplete="new-password"
								required
								bind:value={formData.confirmPassword}
								on:input={() => validateField('confirmPassword')}
								class="input pl-10 pr-10 {errors.confirmPassword ? 'input-error' : ''}"
								placeholder="Confirm your password"
								disabled={isLoading}
							/>
							<div class="absolute inset-y-0 right-0 pr-3 flex items-center">
								<button
									type="button"
									class="text-text-muted hover:text-text-primary transition-colors"
									on:click={toggleConfirmPasswordVisibility}
									disabled={isLoading}
								>
									{#if showConfirmPassword}
										<EyeOff class="h-5 w-5" />
									{:else}
										<Eye class="h-5 w-5" />
									{/if}
								</button>
							</div>
						</div>
						{#if errors.confirmPassword}
							<p class="error-message">{errors.confirmPassword}</p>
						{/if}
					</div>

					<!-- Terms Acceptance -->
					<div>
						<div class="flex items-start">
							<input
								id="accept-terms"
								name="accept-terms"
								type="checkbox"
								bind:checked={acceptTerms}
								class="h-4 w-4 mt-1 text-primary-600 focus:ring-primary-500 border-slate-600 bg-slate-800 rounded"
								disabled={isLoading}
							/>
							<label for="accept-terms" class="ml-2 block text-sm text-text-secondary">
								I agree to the
								<a
									href="/terms"
									target="_blank"
									class="text-primary-400 hover:text-primary-300 transition-colors"
								>
									Terms of Service
								</a>
								and
								<a
									href="/privacy"
									target="_blank"
									class="text-primary-400 hover:text-primary-300 transition-colors"
								>
									Privacy Policy
								</a>
							</label>
						</div>
						{#if errors.terms}
							<p class="error-message">{errors.terms}</p>
						{/if}
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
								Create Account
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
							<span class="px-2 bg-background-secondary text-text-muted">Or sign up with</span>
						</div>
					</div>
				</div>

				<!-- OAuth Buttons -->
				<div class="mt-6 grid grid-cols-1 gap-3">
					<!-- GitHub -->
					<button
						type="button"
						on:click={() => handleOAuthRegister('github')}
						disabled={isLoading}
						class="btn btn-outline w-full flex items-center justify-center"
					>
						<Github class="w-5 h-5 mr-2" />
						Sign up with GitHub
					</button>

					<!-- Google -->
					<button
						type="button"
						on:click={() => handleOAuthRegister('google')}
						disabled={isLoading}
						class="btn btn-outline w-full flex items-center justify-center"
					>
						<Chrome class="w-5 h-5 mr-2" />
						Sign up with Google
					</button>

					<!-- Discord -->
					<button
						type="button"
						on:click={() => handleOAuthRegister('discord')}
						disabled={isLoading}
						class="btn btn-outline w-full flex items-center justify-center"
					>
						<svg class="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 24 24">
							<path
								d="M20.317 4.37a19.791 19.791 0 0 0-4.885-1.515.074.074 0 0 0-.079.037c-.21.375-.444.864-.608 1.25a18.27 18.27 0 0 0-5.487 0 12.64 12.64 0 0 0-.617-1.25.077.077 0 0 0-.079-.037A19.736 19.736 0 0 0 3.677 4.37a.07.07 0 0 0-.032.027C.533 9.046-.32 13.58.099 18.057a.082.082 0 0 0 .031.057 19.9 19.9 0 0 0 5.993 3.03.078.078 0 0 0 .084-.028c.462-.63.874-1.295 1.226-1.994.021-.041.001-.09-.041-.106a13.107 13.107 0 0 1-1.872-.892.077.077 0 0 1-.008-.128 10.2 10.2 0 0 0 .372-.292.074.074 0 0 1 .077-.010c3.928 1.793 8.18 1.793 12.062 0a.074.074 0 0 1 .078.01c.12.098.246.198.373.292a.077.077 0 0 1-.006.127 12.299 12.299 0 0 1-1.873.892.077.077 0 0 0-.041.107c.36.698.772 1.362 1.225 1.993a.076.076 0 0 0 .084.028 19.839 19.839 0 0 0 6.002-3.03.077.077 0 0 0 .032-.054c.5-5.177-.838-9.674-3.549-13.66a.061.061 0 0 0-.031-.03zM8.02 15.33c-1.183 0-2.157-1.085-2.157-2.419 0-1.333.956-2.419 2.157-2.419 1.21 0 2.176 1.096 2.157 2.42 0 1.333-.956 2.418-2.157 2.418zm7.975 0c-1.183 0-2.157-1.085-2.157-2.419 0-1.333.955-2.419 2.157-2.419 1.21 0 2.176 1.096 2.157 2.42 0 1.333-.946 2.418-2.157 2.418z"
							/>
						</svg>
						Sign up with Discord
					</button>
				</div>
			</div>
		</div>

		<!-- Sign In Link -->
		<div class="text-center">
			<p class="text-text-secondary">
				Already have an account?
				<a
					href="/auth/login"
					class="text-primary-400 hover:text-primary-300 font-medium transition-colors"
				>
					Sign in here
				</a>
			</p>
		</div>
	</div>
</div>
