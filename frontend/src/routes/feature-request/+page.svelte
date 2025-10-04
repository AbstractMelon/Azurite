<script lang="ts">
	import { goto } from '$app/navigation';
	import { isAuthenticated, user } from '$lib/stores/auth';
	import { toast } from '$lib/stores/notifications';
	import { Lightbulb, Send, AlertCircle, Heart, Users, Zap } from 'lucide-svelte';

	let title = '';
	let description = '';
	let useCase = '';
	let currentSolution = '';
	let priority = 'medium';
	let category = 'general';
	let implementationIdeas = '';
	let loading = false;
	let submitted = false;

	const priorities = [
		{ value: 'low', label: 'Low', description: 'Nice to have, but not essential' },
		{ value: 'medium', label: 'Medium', description: 'Would improve user experience' },
		{ value: 'high', label: 'High', description: 'Important for functionality' },
		{ value: 'critical', label: 'Critical', description: 'Essential for platform success' }
	];

	const categories = [
		{ value: 'general', label: 'General' },
		{ value: 'ui-ux', label: 'User Interface / UX' },
		{ value: 'mod-management', label: 'Mod Management' },
		{ value: 'search-discovery', label: 'Search & Discovery' },
		{ value: 'community', label: 'Community Features' },
		{ value: 'user-profile', label: 'User Profile' },
		{ value: 'notifications', label: 'Notifications' },
		{ value: 'api', label: 'API / Integration' },
		{ value: 'mobile', label: 'Mobile Experience' },
		{ value: 'performance', label: 'Performance' },
		{ value: 'security', label: 'Security' },
		{ value: 'analytics', label: 'Analytics / Insights' },
		{ value: 'other', label: 'Other' }
	];

	async function handleSubmit(event: Event) {
		event.preventDefault();

		if (!$isAuthenticated) {
			toast.error('Authentication required', 'Please log in to submit a feature request');
			goto('/auth/login?redirect=/feature-request');
			return;
		}

		if (!title.trim() || !description.trim() || !useCase.trim()) {
			toast.error('Missing information', 'Please fill in all required fields');
			return;
		}

		loading = true;

		try {
			// Since we don't have a dedicated feature request API endpoint,
			// we'll simulate submitting it (in a real app, this would go to a product management system)
			await new Promise((resolve) => setTimeout(resolve, 1000)); // Simulate API call

			const featureRequest = {
				title: title.trim(),
				description: description.trim(),
				useCase: useCase.trim(),
				currentSolution: currentSolution.trim(),
				implementationIdeas: implementationIdeas.trim(),
				priority,
				category,
				submittedBy: $user?.username,
				submittedAt: new Date().toISOString()
			};

			console.log('Feature request submitted:', featureRequest);

			submitted = true;
			toast.success(
				'Feature request submitted',
				'Thank you for helping shape the future of Azurite!'
			);
		} catch (error) {
			console.error('Failed to submit feature request:', error);
			toast.error('Submission failed', 'Please try again or contact support directly');
		} finally {
			loading = false;
		}
	}

	function resetForm() {
		title = '';
		description = '';
		useCase = '';
		currentSolution = '';
		implementationIdeas = '';
		priority = 'medium';
		category = 'general';
		submitted = false;
	}

	function getPriorityColor(priorityValue: string): string {
		switch (priorityValue) {
			case 'low':
				return 'text-green-400';
			case 'medium':
				return 'text-yellow-400';
			case 'high':
				return 'text-orange-400';
			case 'critical':
				return 'text-red-400';
			default:
				return 'text-gray-400';
		}
	}
</script>

<svelte:head>
	<title>Feature Request - Azurite</title>
	<meta name="description" content="Suggest new features and improvements for Azurite" />
</svelte:head>

<div class="min-h-screen bg-background-primary">
	<div class="container mx-auto px-4 py-8 max-w-4xl">
		{#if submitted}
			<!-- Success Message -->
			<div class="card">
				<div class="p-8 text-center">
					<div
						class="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4"
					>
						<Lightbulb class="w-8 h-8 text-blue-600" />
					</div>
					<h1 class="text-3xl font-bold text-text-primary mb-4">Feature Request Submitted</h1>
					<p class="text-text-secondary mb-6">
						Thank you for sharing your ideas! Your suggestion helps us understand what our community
						needs and guides our development priorities.
					</p>
					<div class="bg-background-secondary p-4 rounded-lg border border-slate-600 mb-6">
						<p class="text-text-primary font-medium mb-2">What happens next?</p>
						<ul class="text-sm text-text-secondary space-y-1 text-left">
							<li>• Our product team will review your suggestion</li>
							<li>• Popular requests may be prioritized for development</li>
							<li>• We'll consider implementation feasibility and impact</li>
							<li>• You may be contacted for additional input</li>
							<li>• Updates will be posted on our roadmap when available</li>
						</ul>
					</div>
					<div class="flex gap-4 justify-center">
						<button onclick={resetForm} class="btn btn-outline"> Submit Another Request </button>
						<a href="/" class="btn btn-primary"> Back to Home </a>
					</div>
				</div>
			</div>
		{:else}
			<!-- Feature Request Form -->
			<div class="card">
				<div class="p-8">
					<div class="flex items-center gap-3 mb-6">
						<Lightbulb class="h-8 w-8 text-primary-400" />
						<div>
							<h1 class="text-3xl font-bold text-text-primary">Request a Feature</h1>
							<p class="text-text-secondary">Share your ideas to help improve Azurite</p>
						</div>
					</div>

					{#if !$isAuthenticated}
						<div
							class="bg-yellow-100 border border-yellow-400 text-yellow-700 px-4 py-3 rounded mb-6"
						>
							<div class="flex items-center gap-2">
								<AlertCircle class="h-5 w-5" />
								<span>You need to be logged in to submit a feature request.</span>
							</div>
							<a href="/auth/login?redirect=/feature-request" class="underline hover:no-underline">
								Click here to log in
							</a>
						</div>
					{/if}

					<!-- Popular Feature Categories -->
					<div class="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-8">
						<div class="bg-background-secondary p-4 rounded-lg border border-slate-600">
							<Heart class="h-6 w-6 text-primary-400 mb-2" />
							<h3 class="font-medium text-text-primary mb-1">User Experience</h3>
							<p class="text-sm text-text-secondary">
								Interface improvements, usability enhancements
							</p>
						</div>
						<div class="bg-background-secondary p-4 rounded-lg border border-slate-600">
							<Users class="h-6 w-6 text-primary-400 mb-2" />
							<h3 class="font-medium text-text-primary mb-1">Community</h3>
							<p class="text-sm text-text-secondary">Social features, collaboration tools</p>
						</div>
						<div class="bg-background-secondary p-4 rounded-lg border border-slate-600">
							<Zap class="h-6 w-6 text-primary-400 mb-2" />
							<h3 class="font-medium text-text-primary mb-1">Functionality</h3>
							<p class="text-sm text-text-secondary">New capabilities, integrations, tools</p>
						</div>
					</div>

					<form onsubmit={handleSubmit} class="space-y-6">
						<!-- Title -->
						<div>
							<label for="title" class="block text-sm font-medium text-text-primary mb-2">
								Feature Title <span class="text-red-400">*</span>
							</label>
							<input
								id="title"
								type="text"
								bind:value={title}
								placeholder="Brief, descriptive title for your feature idea"
								class="input w-full"
								required
								disabled={loading || !$isAuthenticated}
								maxlength="200"
							/>
							<p class="text-xs text-text-muted mt-1">
								Keep it clear and concise (max 200 characters)
							</p>
						</div>

						<!-- Category and Priority -->
						<div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
							<div>
								<label for="category" class="block text-sm font-medium text-text-primary mb-2">
									Category
								</label>
								<select
									id="category"
									bind:value={category}
									class="select w-full"
									disabled={loading || !$isAuthenticated}
								>
									{#each categories as cat}
										<option value={cat.value}>{cat.label}</option>
									{/each}
								</select>
							</div>

							<div>
								<label for="priority" class="block text-sm font-medium text-text-primary mb-2">
									Priority
								</label>
								<select
									id="priority"
									bind:value={priority}
									class="select w-full"
									disabled={loading || !$isAuthenticated}
								>
									{#each priorities as prio}
										<option value={prio.value} class={getPriorityColor(prio.value)}>
											{prio.label} - {prio.description}
										</option>
									{/each}
								</select>
							</div>
						</div>

						<!-- Description -->
						<div>
							<label for="description" class="block text-sm font-medium text-text-primary mb-2">
								Feature Description <span class="text-red-400">*</span>
							</label>
							<textarea
								id="description"
								bind:value={description}
								placeholder="Detailed description of the feature you'd like to see..."
								class="textarea w-full h-32"
								required
								disabled={loading || !$isAuthenticated}
								maxlength="2000"
							></textarea>
							<p class="text-xs text-text-muted mt-1">
								Explain what the feature should do and how it should work
							</p>
						</div>

						<!-- Use Case -->
						<div>
							<label for="usecase" class="block text-sm font-medium text-text-primary mb-2">
								Use Case / Problem Statement <span class="text-red-400">*</span>
							</label>
							<textarea
								id="usecase"
								bind:value={useCase}
								placeholder="What problem does this feature solve? Who would benefit from it?"
								class="textarea w-full h-24"
								required
								disabled={loading || !$isAuthenticated}
								maxlength="1000"
							></textarea>
							<p class="text-xs text-text-muted mt-1">
								Help us understand the problem this feature would solve
							</p>
						</div>

						<!-- Current Solution -->
						<div>
							<label
								for="current-solution"
								class="block text-sm font-medium text-text-primary mb-2"
							>
								Current Workaround (if any)
							</label>
							<textarea
								id="current-solution"
								bind:value={currentSolution}
								placeholder="How do you currently solve this problem, if at all?"
								class="textarea w-full h-20"
								disabled={loading || !$isAuthenticated}
								maxlength="500"
							></textarea>
							<p class="text-xs text-text-muted mt-1">
								Describe any existing workarounds or limitations
							</p>
						</div>

						<!-- Implementation Ideas -->
						<div>
							<label for="implementation" class="block text-sm font-medium text-text-primary mb-2">
								Implementation Ideas (Optional)
							</label>
							<textarea
								id="implementation"
								bind:value={implementationIdeas}
								placeholder="Any thoughts on how this feature could be implemented?"
								class="textarea w-full h-24"
								disabled={loading || !$isAuthenticated}
								maxlength="1000"
							></textarea>
							<p class="text-xs text-text-muted mt-1">
								Share any technical ideas or similar features you've seen elsewhere
							</p>
						</div>

						<!-- Guidelines -->
						<div class="bg-background-secondary p-4 rounded-lg border border-slate-600">
							<h3 class="font-medium text-text-primary mb-2">Tips for Great Feature Requests</h3>
							<ul class="text-sm text-text-secondary space-y-1">
								<li>• Focus on the problem, not just the solution</li>
								<li>• Explain who would benefit from this feature</li>
								<li>• Include examples or mockups if you have them</li>
								<li>• Consider the impact on different types of users</li>
								<li>• Search existing requests to avoid duplicates</li>
								<li>• Be specific about your needs and constraints</li>
							</ul>
						</div>

						<!-- Submit Button -->
						<div class="flex gap-4">
							<button
								type="submit"
								disabled={loading ||
									!$isAuthenticated ||
									!title.trim() ||
									!description.trim() ||
									!useCase.trim()}
								class="btn btn-primary flex items-center gap-2"
							>
								{#if loading}
									<div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
									Submitting...
								{:else}
									<Send class="h-4 w-4" />
									Submit Feature Request
								{/if}
							</button>

							<a href="/help" class="btn btn-outline"> View Help Center </a>
						</div>
					</form>
				</div>
			</div>

			<!-- Community Voting Notice -->
			<div class="mt-6 text-center">
				<p class="text-text-muted text-sm">
					Feature requests are reviewed by our team and prioritized based on community feedback and
					technical feasibility.
				</p>
			</div>
		{/if}
	</div>
</div>
