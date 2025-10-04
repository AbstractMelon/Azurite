<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { Lock, ArrowLeft, Eye, EyeOff, Check } from 'lucide-svelte';
	import { authApi } from '$lib/api/client';
	import { toast } from '$lib/stores/notifications';

	let password = '';
	let confirmPassword = '';
	let loading = false;
	let success = false;
	let token = '';
	let showPassword = false;
	let showConfirmPassword = false;

	// Password validation states
	$: passwordValid = {
		length: password.length >= 8,
		match: password === confirmPassword && confirmPassword !== ''
	};

	$: allValid = passwordValid.length && passwordValid.match && password && confirmPassword;

	onMount(() => {
		token = $page.url.searchParams.get('token') || '';
		if (!token) {
			toast.error('Invalid Reset Link', 'The reset token is missing or invalid');
		}
	});

	async function handleSubmit(event: Event) {
		// Prevent default form submission
		event.preventDefault();

		if (!password || !confirmPassword) {
			toast.error('Validation Error', 'Please fill in all fields');
			return;
		}

		if (password.length < 8) {
			toast.error('Invalid Password', 'Password must be at least 8 characters long');
			return;
		}

		if (password !== confirmPassword) {
			toast.error('Password Mismatch', 'Passwords do not match');
			return;
		}

		loading = true;

		try {
			const response = await authApi.resetPassword({ token, password });

			if (response.success) {
				success = true;
				toast.success('Password Reset Complete', 'Your password has been reset successfully');

				// Redirect after a delay
				setTimeout(() => {
					goto('/auth/login');
				}, 3000);
			} else {
				toast.error(
					'Reset Failed',
					response.error || 'An error occurred while resetting your password'
				);
			}
		} catch (error) {
			console.error('Reset password error:', error);
			toast.error('Network Error', 'Please try again later');
		} finally {
			loading = false;
		}
	}

	function goToLogin() {
		goto('/auth/login');
	}

	function togglePasswordVisibility() {
		showPassword = !showPassword;
	}

	function toggleConfirmPasswordVisibility() {
		showConfirmPassword = !showConfirmPassword;
	}
</script>

<svelte:head>
	<title>Reset Password - Azurite</title>
	<meta name="description" content="Create a new password for your Azurite account." />
</svelte:head>

<div class="min-h-screen bg-background-primary">
	<!-- Header -->
	<div class="bg-gradient-to-r from-slate-800/50 to-slate-700/50 border-b border-slate-700">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
			<div class="text-center">
				<h1 class="text-3xl font-bold text-text-primary mb-2">
					<Lock class="w-8 h-8 inline-block mr-2 mb-1" />
					Set New Password
				</h1>
				<p class="text-text-secondary">Create a new secure password for your account</p>
			</div>
		</div>
	</div>

	<!-- Content -->
	<div class="max-w-md mx-auto px-4 py-12">
		<div class="card">
			<div class="p-8">
				{#if success}
					<!-- Success State -->
					<div class="text-center">
						<div
							class="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-6"
						>
							<Check class="w-8 h-8 text-green-600" />
						</div>

						<h2 class="text-xl font-semibold text-text-primary mb-3">Password Reset Complete</h2>

						<div class="bg-green-900/20 border border-green-600 rounded-lg p-4 mb-6">
							<p class="text-green-200 text-sm">
								Your password has been reset successfully! You can now log in with your new
								password.
							</p>
						</div>

						<p class="text-text-secondary text-sm mb-8">
							You will be redirected to the login page shortly.
						</p>

						<button onclick={goToLogin} class="btn btn-primary w-full">
							<ArrowLeft class="w-4 h-4 mr-2" />
							Go to Login
						</button>
					</div>
				{:else if !token}
					<!-- Invalid Token State -->
					<div class="text-center">
						<div class="bg-red-900/20 border border-red-600 rounded-lg p-4 mb-6">
							<p class="text-red-200 text-sm">
								Invalid or missing reset token. This link may have expired or been used already.
							</p>
						</div>
						<button onclick={goToLogin} class="btn btn-primary w-full">
							<ArrowLeft class="w-4 h-4 mr-2" />
							Back to Login
						</button>
					</div>
				{:else}
					<!-- Form State -->
					<div class="mb-6">
						<h2 class="text-xl font-semibold text-text-primary mb-2 text-center">
							Create New Password
						</h2>
						<p class="text-text-secondary text-sm text-center">
							Choose a strong password to keep your account secure.
						</p>
					</div>

					<form onsubmit={handleSubmit} class="space-y-6">
						<!-- New Password Field -->
						<div>
							<label for="password" class="block text-sm font-medium text-text-primary mb-2">
								New Password
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
									bind:value={password}
									required
									minlength="8"
									class="input pl-10 pr-12"
									placeholder="Enter your new password"
									disabled={loading}
								/>
								<button
									type="button"
									onclick={togglePasswordVisibility}
									class="absolute inset-y-0 right-0 pr-3 flex items-center text-text-muted hover:text-text-secondary transition-colors"
									disabled={loading}
								>
									{#if showPassword}
										<EyeOff class="h-5 w-5" />
									{:else}
										<Eye class="h-5 w-5" />
									{/if}
								</button>
							</div>
						</div>

						<!-- Confirm Password Field -->
						<div>
							<label
								for="confirm-password"
								class="block text-sm font-medium text-text-primary mb-2"
							>
								Confirm New Password
								<span class="text-red-400">*</span>
							</label>
							<div class="relative">
								<div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
									<Lock class="h-5 w-5 text-text-muted" />
								</div>
								<input
									id="confirm-password"
									name="confirm-password"
									type={showConfirmPassword ? 'text' : 'password'}
									bind:value={confirmPassword}
									required
									minlength="8"
									class="input pl-10 pr-12"
									placeholder="Confirm your new password"
									disabled={loading}
								/>
								<button
									type="button"
									onclick={toggleConfirmPasswordVisibility}
									class="absolute inset-y-0 right-0 pr-3 flex items-center text-text-muted hover:text-text-secondary transition-colors"
									disabled={loading}
								>
									{#if showConfirmPassword}
										<EyeOff class="h-5 w-5" />
									{:else}
										<Eye class="h-5 w-5" />
									{/if}
								</button>
							</div>
						</div>

						<!-- Password Requirements -->
						<div class="bg-slate-800 border border-slate-700 rounded-lg p-4">
							<p class="text-sm font-medium text-text-primary mb-3">Password requirements:</p>
							<ul class="space-y-2">
								<li class="flex items-center text-sm">
									{#if passwordValid.length}
										<Check class="w-4 h-4 text-green-400 mr-2 flex-shrink-0" />
										<span class="text-green-400">At least 8 characters long</span>
									{:else}
										<div
											class="w-4 h-4 border border-text-muted rounded-full mr-2 flex-shrink-0"
										></div>
										<span class="text-text-muted">At least 8 characters long</span>
									{/if}
								</li>
								<li class="flex items-center text-sm">
									{#if passwordValid.match}
										<Check class="w-4 h-4 text-green-400 mr-2 flex-shrink-0" />
										<span class="text-green-400">Passwords match</span>
									{:else}
										<div
											class="w-4 h-4 border border-text-muted rounded-full mr-2 flex-shrink-0"
										></div>
										<span class="text-text-muted">Passwords must match</span>
									{/if}
								</li>
							</ul>
						</div>

						<!-- Submit Button -->
						<button type="submit" disabled={loading || !allValid} class="btn btn-primary w-full">
							{#if loading}
								<div class="flex items-center justify-center">
									<svg
										class="animate-spin -ml-1 mr-3 h-5 w-5"
										xmlns="http://www.w3.org/2000/svg"
										fill="none"
										viewBox="0 0 24 24"
									>
										<circle
											class="opacity-25"
											cx="12"
											cy="12"
											r="10"
											stroke="currentColor"
											stroke-width="4"
										></circle>
										<path
											class="opacity-75"
											fill="currentColor"
											d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
										></path>
									</svg>
									Resetting Password...
								</div>
							{:else}
								<Lock class="w-4 h-4 mr-2" />
								Reset Password
							{/if}
						</button>

						<!-- Back to Login -->
						<div class="text-center pt-4 border-t border-slate-700">
							<button
								type="button"
								onclick={goToLogin}
								class="text-sm text-primary-400 hover:text-primary-300 font-medium transition-colors"
							>
								<ArrowLeft class="w-4 h-4 inline mr-1" />
								Back to login
							</button>
						</div>
					</form>
				{/if}
			</div>
		</div>

		<!-- Security Note -->
		<div class="text-center mt-8">
			<p class="text-text-muted text-sm">
				Make sure you're on a secure connection and using a trusted device.
			</p>
		</div>
	</div>
</div>
