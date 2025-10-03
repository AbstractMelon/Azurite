<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { isAuthenticated } from '$lib/stores/auth';
	import { gamesApi } from '$lib/api/client';
	import { toast } from '$lib/stores/notifications';
	import Loading from '$lib/components/Loading.svelte';
	import {
		Gamepad2,
		ArrowLeft,
		Send,
		Info,
		CheckCircle,
		AlertTriangle,
		FileText
	} from 'lucide-svelte';

	// Form data
	let formData = {
		name: '',
		reason: '',
		description: '',
		existing_community: '',
		mod_loader: '',
		contact: ''
	};

	// UI state
	let isSubmitting = false;
	let isLoading = true;
	let errors: { [key: string]: string } = {};
	// Guidelines are always visible

	// Form validation
	function validateForm() {
		errors = {};

		if (!formData.name.trim()) {
			errors.name = 'Game name is required';
		} else if (formData.name.trim().length < 2) {
			errors.name = 'Game name must be at least 2 characters';
		} else if (formData.name.trim().length > 100) {
			errors.name = 'Game name must be less than 100 characters';
		}

		if (!formData.reason.trim()) {
			errors.reason = 'Reason is required';
		} else if (formData.reason.trim().length < 20) {
			errors.reason = 'Reason must be at least 20 characters';
		} else if (formData.reason.trim().length > 500) {
			errors.reason = 'Reason must be less than 500 characters';
		}

		if (!formData.description.trim()) {
			errors.description = 'Description is required';
		} else if (formData.description.trim().length < 50) {
			errors.description = 'Description must be at least 50 characters';
		} else if (formData.description.trim().length > 1000) {
			errors.description = 'Description must be less than 1000 characters';
		}

		return Object.keys(errors).length === 0;
	}

	// Submit form
	async function handleSubmit() {
		if (!validateForm()) {
			toast.error('Validation Error', 'Please fix the errors below');
			return;
		}

		isSubmitting = true;

		try {
			const response = await gamesApi.createGameRequest({
				name: formData.name.trim(),
				reason: formData.reason.trim(),
				description: formData.description.trim(),
				existing_community: formData.existing_community.trim(),
				mod_loader: formData.mod_loader.trim(),
				contact: formData.contact.trim()
			});

			if (response.success) {
				toast.success(
					'Request Submitted',
					"Your game request has been submitted for review. We'll notify you when it's processed."
				);
				goto('/games');
			} else {
				toast.error('Submission Failed', response.error || 'Failed to submit game request');
			}
		} catch (error) {
			console.error('Error submitting game request:', error);
			toast.error('Error', 'An unexpected error occurred. Please try again.');
		} finally {
			isSubmitting = false;
		}
	}

	// Check authentication on mount
	onMount(() => {
		if (!$isAuthenticated) {
			goto('/auth/login?redirect=/games/request');
			return;
		}
		isLoading = false;
	});
</script>

<svelte:head>
	<title>Request New Game - Azurite</title>
	<meta
		name="description"
		content="Request support for a new game on the Azurite modding platform"
	/>
</svelte:head>

{#if isLoading}
	<div class="min-h-screen bg-background-primary flex items-center justify-center">
		<Loading size="lg" text="Loading..." />
	</div>
{:else}
	<div class="min-h-screen bg-background-primary">
		<div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
			<!-- Header -->
			<div class="mb-8">
				<button
					on:click={() => goto('/games')}
					class="flex items-center text-text-muted hover:text-text-secondary mb-6"
				>
					<ArrowLeft class="w-4 h-4 mr-2" />
					Back to Games
				</button>

				<div class="text-center">
					<div
						class="inline-flex items-center justify-center w-16 h-16 bg-primary-500/20 rounded-full mb-6"
					>
						<Gamepad2 class="w-8 h-8 text-primary-400" />
					</div>
					<h1 class="text-4xl font-bold text-text-primary mb-4">Request New Game</h1>
					<p class="text-text-secondary text-lg max-w-2xl mx-auto">
						Help expand the Azurite community by requesting support for a new game. Our team will
						review your request and consider adding it to the platform.
					</p>
				</div>
			</div>

			<div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
				<!-- Form -->
				<div class="lg:col-span-2">
					<div class="card">
						<div class="p-8">
							<form on:submit|preventDefault={handleSubmit} class="space-y-6">
								<!-- Game Name -->
								<div>
									<label for="name" class="block text-sm font-medium text-text-primary mb-2">
										Game Name *
									</label>
									<input
										id="name"
										type="text"
										bind:value={formData.name}
										placeholder="e.g., Terraria, Stardew Valley, Factorio"
										class="input w-full {errors.name ? 'border-red-500' : ''}"
										required
										maxlength="100"
									/>
									{#if errors.name}
										<p class="mt-1 text-sm text-red-400">{errors.name}</p>
									{/if}
									<p class="mt-1 text-sm text-text-muted">
										The official name of the game you'd like to see on Azurite
									</p>
								</div>

								<!-- Reason -->
								<div>
									<label for="reason" class="block text-sm font-medium text-text-primary mb-2">
										Why should this game be added? *
									</label>
									<textarea
										id="reason"
										bind:value={formData.reason}
										placeholder="Briefly explain why this game would be a great addition to Azurite..."
										class="textarea w-full h-24 {errors.reason ? 'border-red-500' : ''}"
										required
										maxlength="500"
									></textarea>
									{#if errors.reason}
										<p class="mt-1 text-sm text-red-400">{errors.reason}</p>
									{/if}
									<div class="flex justify-between mt-1">
										<p class="text-sm text-text-muted">
											Provide a concise reason for adding this game
										</p>
										<span class="text-sm text-text-muted">
											{formData.reason.length}/500
										</span>
									</div>
								</div>

								<!-- Description -->
								<div>
									<label for="description" class="block text-sm font-medium text-text-primary mb-2">
										Game Description *
									</label>
									<textarea
										id="description"
										bind:value={formData.description}
										placeholder="Provide detailed information about the game, its popularity, modding community, and potential for user-generated content..."
										class="textarea w-full h-32 {errors.description ? 'border-red-500' : ''}"
										required
										maxlength="1000"
									></textarea>
									{#if errors.description}
										<p class="mt-1 text-sm text-red-400">{errors.description}</p>
									{/if}
									<div class="flex justify-between mt-1">
										<p class="text-sm text-text-muted">
											Detailed description of the game and its modding potential
										</p>
										<span class="text-sm text-text-muted">
											{formData.description.length}/1000
										</span>
									</div>
								</div>

								<!-- Existing Community -->
								<div>
									<label
										for="existing_community"
										class="block text-sm font-medium text-text-primary mb-2"
									>
										Existing Community
									</label>
									<input
										id="existing_community"
										type="text"
										bind:value={formData.existing_community}
										placeholder="e.g., Discord server, subreddit, forum URL"
										class="input w-full"
										maxlength="200"
									/>
									<p class="mt-1 text-sm text-text-muted">
										Link to existing modding community or forums (optional)
									</p>
								</div>

								<!-- Mod Loader -->
								<div>
									<label for="mod_loader" class="block text-sm font-medium text-text-primary mb-2">
										Mod Loader / Framework
									</label>
									<input
										id="mod_loader"
										type="text"
										bind:value={formData.mod_loader}
										placeholder="e.g., Forge, Fabric, BepInEx, MelonLoader"
										class="input w-full"
										maxlength="100"
									/>
									<p class="mt-1 text-sm text-text-muted">
										Popular mod loader or framework used for this game (optional)
									</p>
								</div>

								<!-- Contact -->
								<div>
									<label for="contact" class="block text-sm font-medium text-text-primary mb-2">
										Contact Information
									</label>
									<input
										id="contact"
										type="text"
										bind:value={formData.contact}
										placeholder="e.g., Discord username, email"
										class="input w-full"
										maxlength="100"
									/>
									<p class="mt-1 text-sm text-text-muted">
										How we can reach you for follow-up questions (optional)
									</p>
								</div>

								<!-- Submit Button -->
								<div class="flex justify-end">
									<button
										type="submit"
										disabled={isSubmitting}
										class="btn btn-primary min-w-[140px]"
									>
										{#if isSubmitting}
											<Loading size="sm" />
											Submitting...
										{:else}
											<Send class="w-4 h-4 mr-2" />
											Submit Request
										{/if}
									</button>
								</div>
							</form>
						</div>
					</div>
				</div>

				<!-- Sidebar -->
				<div class="space-y-6">
					<!-- Guidelines -->
					<div class="card">
						<div class="p-6">
							<div class="flex items-center mb-4">
								<FileText class="w-5 h-5 text-primary-400 mr-2" />
								<h3 class="text-lg font-semibold text-text-primary">Request Guidelines</h3>
							</div>
							<div class="space-y-3 text-sm text-text-secondary">
								<div class="flex items-start">
									<CheckCircle class="w-4 h-4 text-green-400 mr-2 mt-0.5 flex-shrink-0" />
									<span>Use the official game name</span>
								</div>
								<div class="flex items-start">
									<CheckCircle class="w-4 h-4 text-green-400 mr-2 mt-0.5 flex-shrink-0" />
									<span>Explain the game's modding potential</span>
								</div>
								<div class="flex items-start">
									<CheckCircle class="w-4 h-4 text-green-400 mr-2 mt-0.5 flex-shrink-0" />
									<span>Mention existing community size</span>
								</div>
								<div class="flex items-start">
									<CheckCircle class="w-4 h-4 text-green-400 mr-2 mt-0.5 flex-shrink-0" />
									<span>Be specific about why it fits Azurite</span>
								</div>
							</div>
						</div>
					</div>

					<!-- Process Info -->
					<div class="card bg-blue-900/20 border-blue-600/50">
						<div class="p-6">
							<div class="flex items-center mb-4">
								<Info class="w-5 h-5 text-blue-400 mr-2" />
								<h3 class="text-lg font-semibold text-blue-300">Review Process</h3>
							</div>
							<div class="space-y-3 text-sm text-blue-200">
								<p><strong>1. Submission:</strong> Your request is received and queued</p>
								<p><strong>2. Review:</strong> Our team evaluates the game's suitability</p>
								<p><strong>3. Decision:</strong> We'll notify you of approval or feedback</p>
								<p><strong>4. Implementation:</strong> Approved games are added to the platform</p>
							</div>
							<div class="mt-4 p-3 bg-blue-800/30 rounded-lg">
								<p class="text-xs text-blue-200">
									<strong>Timeline:</strong> Most requests are reviewed within 1-2 weeks. Popular games
									with active communities are prioritized.
								</p>
							</div>
						</div>
					</div>

					<!-- Requirements -->
					<div class="card bg-yellow-900/20 border-yellow-600/50">
						<div class="p-6">
							<div class="flex items-center mb-4">
								<AlertTriangle class="w-5 h-5 text-yellow-400 mr-2" />
								<h3 class="text-lg font-semibold text-yellow-300">Requirements</h3>
							</div>
							<div class="space-y-2 text-sm text-yellow-200">
								<p>• Game must support user modifications</p>
								<p>• Active player and modding community</p>
								<p>• Not already available on Azurite</p>
								<p>• Appropriate content (no adult-only games)</p>
								<p>• Legal to modify and distribute mods</p>
							</div>
						</div>
					</div>

					<!-- Help -->
					<div class="card">
						<div class="p-6">
							<h3 class="text-lg font-semibold text-text-primary mb-3">Need Help?</h3>
							<p class="text-sm text-text-secondary mb-4">
								Have questions about submitting a game request?
							</p>
							<div class="space-y-2">
								<a href="/docs" class="btn btn-outline btn-sm w-full">
									<FileText class="w-4 h-4 mr-2" />
									View Documentation
								</a>
								<a href="/contact" class="btn btn-outline btn-sm w-full">
									<Send class="w-4 h-4 mr-2" />
									Contact Support
								</a>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
{/if}
