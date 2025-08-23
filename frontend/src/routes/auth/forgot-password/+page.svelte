<script lang="ts">
	import { goto } from '$app/navigation';
	import { Mail, ArrowLeft, Check } from 'lucide-svelte';
	import { authApi } from '$lib/api/client';
	import { toast } from '$lib/stores/notifications';

	let email = '';
	let loading = false;
	let success = false;

	async function handleSubmit() {
		if (!email) {
			toast.error('Validation Error', 'Please enter your email address');
			return;
		}

		loading = true;

		try {
			const response = await authApi.forgotPassword(email);

			if (response.success) {
				success = true;
				toast.success('Reset link sent', 'Check your email for the password reset link');
			} else {
				toast.error('Reset failed', response.error || 'An error occurred');
			}
		} catch (e) {
			console.error('Forgot password error:', e);
			toast.error('Network error', 'Please try again later');
		} finally {
			loading = false;
		}
	}

	function goToLogin() {
		goto('/auth/login');
	}
</script>

<svelte:head>
	<title>Forgot Password - Azurite</title>
	<meta
		name="description"
		content="Reset your Azurite account password. Enter your email to receive a password reset link."
	/>
</svelte:head>

<div class="min-h-screen bg-background-primary">
	<!-- Header -->
	<div class="bg-gradient-to-r from-slate-800/50 to-slate-700/50 border-b border-slate-700">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
			<div class="text-center">
				<h1 class="text-3xl font-bold text-text-primary mb-2">
					<Mail class="w-8 h-8 inline-block mr-2 mb-1" />
					Reset Password
				</h1>
				<p class="text-text-secondary">
					Enter your email address and we'll send you a link to reset your password
				</p>
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
							class="w-16 h-16 bg-green-600 rounded-full flex items-center justify-center mx-auto mb-6"
						>
							<Check class="w-8 h-8 text-green-600" />
						</div>

						<h2 class="text-xl font-semibold text-text-primary mb-3">Check your email</h2>

						<div class="bg-green-900/20 border border-green-600 rounded-lg p-4 mb-6">
							<p class="text-green-200 text-sm">
								If an account with <strong>{email}</strong> exists, we've sent you a password reset link.
							</p>
						</div>

						<p class="text-text-secondary text-sm mb-8">
							Didn't receive an email? Check your spam folder or try again with a different email
							address.
						</p>

						<div class="space-y-3">
							<button on:click={goToLogin} class="btn btn-primary w-full">
								<ArrowLeft class="w-4 h-4 mr-2" />
								Back to Login
							</button>

							<button
								on:click={() => {
									success = false;
									email = '';
								}}
								class="btn btn-outline w-full"
							>
								Try different email
							</button>
						</div>
					</div>
				{:else}
					<!-- Form State -->
					<div class="mb-6">
						<h2 class="text-xl font-semibold text-text-primary mb-2 text-center">
							Forgot your password?
						</h2>
						<p class="text-text-secondary text-sm text-center">
							No worries! We'll help you get back into your account.
						</p>
					</div>

					<form on:submit|preventDefault={handleSubmit} class="space-y-6">
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
									bind:value={email}
									required
									class="input pl-10"
									placeholder="Enter your email address"
									disabled={loading}
								/>
							</div>
						</div>

						<button
							type="submit"
							disabled={loading || !email.trim()}
							class="btn btn-primary w-full"
						>
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
									Sending reset link...
								</div>
							{:else}
								<Mail class="w-4 h-4 mr-2" />
								Send reset link
							{/if}
						</button>

						<!-- Back to Login -->
						<div class="text-center pt-4 border-t border-slate-700">
							<button
								type="button"
								on:click={goToLogin}
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

		<!-- Help Text -->
		<div class="text-center mt-8">
			<p class="text-text-muted text-sm">
				Need help? <a
					href="/support"
					class="text-primary-400 hover:text-primary-300 transition-colors">Contact support</a
				>
			</p>
		</div>
	</div>
</div>
