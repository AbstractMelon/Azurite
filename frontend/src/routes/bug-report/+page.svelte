<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { isAuthenticated, user } from '$lib/stores/auth';
	import { toast } from '$lib/stores/notifications';
	import { Bug, Send, AlertCircle } from 'lucide-svelte';

	let title = '';
	let description = '';
	let stepsToReproduce = '';
	let expectedBehavior = '';
	let actualBehavior = '';
	let browserInfo = '';
	let operatingSystem = '';
	let priority = 'medium';
	let category = 'general';
	let loading = false;
	let submitted = false;

	const priorities = [
		{ value: 'low', label: 'Low', description: 'Minor issue, doesn\'t affect core functionality' },
		{ value: 'medium', label: 'Medium', description: 'Noticeable issue, affects some functionality' },
		{ value: 'high', label: 'High', description: 'Significant issue, affects important functionality' },
		{ value: 'critical', label: 'Critical', description: 'Severe issue, breaks core functionality' }
	];

	const categories = [
		{ value: 'general', label: 'General' },
		{ value: 'authentication', label: 'Authentication' },
		{ value: 'mod-upload', label: 'Mod Upload' },
		{ value: 'mod-download', label: 'Mod Download' },
		{ value: 'search', label: 'Search' },
		{ value: 'user-profile', label: 'User Profile' },
		{ value: 'comments', label: 'Comments' },
		{ value: 'notifications', label: 'Notifications' },
		{ value: 'ui-ux', label: 'UI/UX' },
		{ value: 'performance', label: 'Performance' },
		{ value: 'security', label: 'Security' },
		{ value: 'other', label: 'Other' }
	];

	onMount(() => {
		// Auto-detect browser and OS information
		browserInfo = navigator.userAgent;
		operatingSystem = navigator.platform;
	});

	async function handleSubmit(event: Event) {
		event.preventDefault();

		if (!$isAuthenticated) {
			toast.error('Authentication required', 'Please log in to submit a bug report');
			goto('/auth/login?redirect=/bug-report');
			return;
		}

		if (!title.trim() || !description.trim()) {
			toast.error('Missing information', 'Please fill in all required fields');
			return;
		}

		loading = true;

		try {
			// Since we don't have a dedicated bug report API endpoint,
			// we'll simulate submitting it (in a real app, this would go to a support system)
			await new Promise(resolve => setTimeout(resolve, 1000)); // Simulate API call

			const bugReport = {
				title: title.trim(),
				description: description.trim(),
				stepsToReproduce: stepsToReproduce.trim(),
				expectedBehavior: expectedBehavior.trim(),
				actualBehavior: actualBehavior.trim(),
				browserInfo,
				operatingSystem,
				priority,
				category,
				submittedBy: $user?.username,
				submittedAt: new Date().toISOString()
			};

			console.log('Bug report submitted:', bugReport);

			submitted = true;
			toast.success('Bug report submitted', 'Thank you for helping us improve Azurite!');
		} catch (error) {
			console.error('Failed to submit bug report:', error);
			toast.error('Submission failed', 'Please try again or contact support directly');
		} finally {
			loading = false;
		}
	}

	function resetForm() {
		title = '';
		description = '';
		stepsToReproduce = '';
		expectedBehavior = '';
		actualBehavior = '';
		priority = 'medium';
		category = 'general';
		submitted = false;
	}

	function getPriorityColor(priorityValue: string): string {
		switch (priorityValue) {
			case 'low': return 'text-green-400';
			case 'medium': return 'text-yellow-400';
			case 'high': return 'text-orange-400';
			case 'critical': return 'text-red-400';
			default: return 'text-gray-400';
		}
	}
</script>

<svelte:head>
	<title>Bug Report - Azurite</title>
	<meta name="description" content="Report bugs and technical issues to help improve Azurite" />
</svelte:head>

<div class="min-h-screen bg-background-primary">
	<div class="container mx-auto px-4 py-8 max-w-4xl">
		{#if submitted}
			<!-- Success Message -->
			<div class="card">
				<div class="p-8 text-center">
					<div class="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
						<svg class="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
						</svg>
					</div>
					<h1 class="text-3xl font-bold text-text-primary mb-4">Bug Report Submitted</h1>
					<p class="text-text-secondary mb-6">
						Thank you for taking the time to report this issue. Our team will review your report
						and work on a fix as soon as possible.
					</p>
					<div class="bg-background-secondary p-4 rounded-lg border border-slate-600 mb-6">
						<p class="text-text-primary font-medium mb-2">What happens next?</p>
						<ul class="text-sm text-text-secondary space-y-1 text-left">
							<li>• Our development team will review your report</li>
							<li>• We may reach out if we need additional information</li>
							<li>• You'll be notified when the issue is resolved</li>
							<li>• Critical issues are prioritized and addressed first</li>
						</ul>
					</div>
					<div class="flex gap-4 justify-center">
						<button on:click={resetForm} class="btn btn-outline">
							Submit Another Report
						</button>
						<a href="/" class="btn btn-primary">
							Back to Home
						</a>
					</div>
				</div>
			</div>
		{:else}
			<!-- Bug Report Form -->
			<div class="card">
				<div class="p-8">
					<div class="flex items-center gap-3 mb-6">
						<Bug class="h-8 w-8 text-primary-400" />
						<div>
							<h1 class="text-3xl font-bold text-text-primary">Report a Bug</h1>
							<p class="text-text-secondary">Help us improve Azurite by reporting issues you encounter</p>
						</div>
					</div>

					{#if !$isAuthenticated}
						<div class="bg-yellow-100 border border-yellow-400 text-yellow-700 px-4 py-3 rounded mb-6">
							<div class="flex items-center gap-2">
								<AlertCircle class="h-5 w-5" />
								<span>You need to be logged in to submit a bug report.</span>
							</div>
							<a href="/auth/login?redirect=/bug-report" class="underline hover:no-underline">
								Click here to log in
							</a>
						</div>
					{/if}

					<form on:submit={handleSubmit} class="space-y-6">
						<!-- Title -->
						<div>
							<label for="title" class="block text-sm font-medium text-text-primary mb-2">
								Bug Title <span class="text-red-400">*</span>
							</label>
							<input
								id="title"
								type="text"
								bind:value={title}
								placeholder="Brief description of the issue"
								class="input w-full"
								required
								disabled={loading || !$isAuthenticated}
								maxlength="200"
							/>
							<p class="text-xs text-text-muted mt-1">
								Keep it concise and descriptive (max 200 characters)
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
								Description <span class="text-red-400">*</span>
							</label>
							<textarea
								id="description"
								bind:value={description}
								placeholder="Detailed description of the bug..."
								class="textarea w-full h-32"
								required
								disabled={loading || !$isAuthenticated}
								maxlength="2000"
							></textarea>
							<p class="text-xs text-text-muted mt-1">
								Provide as much detail as possible about the issue
							</p>
						</div>

						<!-- Steps to Reproduce -->
						<div>
							<label for="steps" class="block text-sm font-medium text-text-primary mb-2">
								Steps to Reproduce
							</label>
							<textarea
								id="steps"
								bind:value={stepsToReproduce}
								placeholder="1. Go to...&#10;2. Click on...&#10;3. Observe..."
								class="textarea w-full h-24"
								disabled={loading || !$isAuthenticated}
								maxlength="1000"
							></textarea>
							<p class="text-xs text-text-muted mt-1">
								List the exact steps to reproduce the issue
							</p>
						</div>

						<!-- Expected vs Actual Behavior -->
						<div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
							<div>
								<label for="expected" class="block text-sm font-medium text-text-primary mb-2">
									Expected Behavior
								</label>
								<textarea
									id="expected"
									bind:value={expectedBehavior}
									placeholder="What should happen..."
									class="textarea w-full h-20"
									disabled={loading || !$isAuthenticated}
									maxlength="500"
								></textarea>
							</div>

							<div>
								<label for="actual" class="block text-sm font-medium text-text-primary mb-2">
									Actual Behavior
								</label>
								<textarea
									id="actual"
									bind:value={actualBehavior}
									placeholder="What actually happens..."
									class="textarea w-full h-20"
									disabled={loading || !$isAuthenticated}
									maxlength="500"
								></textarea>
							</div>
						</div>

						<!-- System Information -->
						<div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
							<div>
								<label for="browser" class="block text-sm font-medium text-text-primary mb-2">
									Browser Information
								</label>
								<textarea
									id="browser"
									bind:value={browserInfo}
									class="textarea w-full h-16 text-xs"
									readonly
									disabled={loading || !$isAuthenticated}
								></textarea>
							</div>

							<div>
								<label for="os" class="block text-sm font-medium text-text-primary mb-2">
									Operating System
								</label>
								<input
									id="os"
									type="text"
									bind:value={operatingSystem}
									class="input w-full"
									readonly
									disabled={loading || !$isAuthenticated}
								/>
							</div>
						</div>

						<!-- Additional Information -->
						<div class="bg-background-secondary p-4 rounded-lg border border-slate-600">
							<h3 class="font-medium text-text-primary mb-2">Tips for Better Bug Reports</h3>
							<ul class="text-sm text-text-secondary space-y-1">
								<li>• Include screenshots or screen recordings if possible</li>
								<li>• Mention if the issue is consistently reproducible</li>
								<li>• Include any error messages you see</li>
								<li>• Test in different browsers if applicable</li>
								<li>• Clear your browser cache and try again before reporting</li>
							</ul>
						</div>

						<!-- Submit Button -->
						<div class="flex gap-4">
							<button
								type="submit"
								disabled={loading || !$isAuthenticated || !title.trim() || !description.trim()}
								class="btn btn-primary flex items-center gap-2"
							>
								{#if loading}
									<div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
									Submitting...
								{:else}
									<Send class="h-4 w-4" />
									Submit Bug Report
								{/if}
							</button>

							<a href="/help" class="btn btn-outline">
								View Help Center
							</a>
						</div>
					</form>
				</div>
			</div>
		{/if}
	</div>
</div>
