<script lang="ts">
	import { goto } from '$app/navigation';
	import { adminApi } from '$lib/api/client';
	import { toast } from '$lib/stores/notifications';
	import { ArrowLeft, Save, Gamepad2 } from 'lucide-svelte';
	import Loading from '$lib/components/Loading.svelte';

	let formData = {
		name: '',
		description: '',
		icon: ''
	};

	let isSubmitting = false;
	let errors: { [key: string]: string } = {};

	function validateForm() {
		errors = {};

		if (!formData.name.trim()) {
			errors.name = 'Game name is required';
		} else if (formData.name.trim().length < 2) {
			errors.name = 'Game name must be at least 2 characters';
		}

		if (!formData.description.trim()) {
			errors.description = 'Description is required';
		}

		return Object.keys(errors).length === 0;
	}

	async function handleSubmit(event: Event) {
		// Prevent default form submission
		event.preventDefault();

		if (!validateForm()) {
			toast.error('Validation Error', 'Please fix the errors below');
			return;
		}

		isSubmitting = true;

		try {
			const response = await adminApi.createGame({
				name: formData.name.trim(),
				description: formData.description.trim(),
				icon: formData.icon.trim()
			});

			if (response.success) {
				toast.success('Game Created', `"${formData.name}" has been added successfully.`);
				goto('/admin?tab=games');
			} else {
				toast.error('Creation Failed', response.error || 'Failed to create game');
			}
		} catch (error) {
			console.error('Error creating game:', error);
			toast.error('Error', 'An unexpected error occurred. Please try again.');
		} finally {
			isSubmitting = false;
		}
	}
</script>

<svelte:head>
	<title>Add New Game - Admin - Azurite</title>
</svelte:head>

<div class="min-h-screen bg-background-primary">
	<div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
		<!-- Header -->
		<div class="mb-8">
			<button
				onclick={() => goto('/admin?tab=games')}
				class="flex items-center text-text-muted hover:text-text-secondary mb-6"
			>
				<ArrowLeft class="w-4 h-4 mr-2" />
				Back to Admin
			</button>

			<div class="flex items-center space-x-3">
				<div class="w-12 h-12 bg-primary-500/20 rounded-lg flex items-center justify-center">
					<Gamepad2 class="w-6 h-6 text-primary-400" />
				</div>
				<div>
					<h1 class="text-3xl font-bold text-text-primary">Add New Game</h1>
					<p class="text-text-secondary mt-1">Create a new game entry on the platform</p>
				</div>
			</div>
		</div>

		<!-- Form -->
		<div class="card">
			<div class="p-8">
				<form onsubmit={handleSubmit} class="space-y-6">
					<!-- Game Name -->
					<div>
						<label for="name" class="block text-sm font-medium text-text-primary mb-2">
							Game Name *
						</label>
						<input
							id="name"
							type="text"
							bind:value={formData.name}
							placeholder="e.g., Minecraft, Terraria"
							class="input w-full {errors.name ? 'border-red-500' : ''}"
							required
						/>
						{#if errors.name}
							<p class="mt-1 text-sm text-red-400">{errors.name}</p>
						{/if}
					</div>

					<!-- Description -->
					<div>
						<label for="description" class="block text-sm font-medium text-text-primary mb-2">
							Description *
						</label>
						<textarea
							id="description"
							bind:value={formData.description}
							placeholder="Brief description of the game..."
							class="textarea w-full h-32 {errors.description ? 'border-red-500' : ''}"
							required
						></textarea>
						{#if errors.description}
							<p class="mt-1 text-sm text-red-400">{errors.description}</p>
						{/if}
					</div>

					<!-- Icon URL -->
					<div>
						<label for="icon" class="block text-sm font-medium text-text-primary mb-2">
							Icon URL
						</label>
						<input
							id="icon"
							type="url"
							bind:value={formData.icon}
							placeholder="https://example.com/icon.png"
							class="input w-full"
						/>
						<p class="mt-1 text-sm text-text-muted">Optional: URL to the game's icon image</p>
					</div>

					<!-- Submit Button -->
					<div class="flex justify-end space-x-3">
						<button
							type="button"
							onclick={() => goto('/admin?tab=games')}
							class="btn btn-outline"
						>
							Cancel
						</button>
						<button type="submit" disabled={isSubmitting} class="btn btn-primary min-w-[140px]">
							{#if isSubmitting}
								<Loading size="sm" />
								Creating...
							{:else}
								<Save class="w-4 h-4 mr-2" />
								Create Game
							{/if}
						</button>
					</div>
				</form>
			</div>
		</div>
	</div>
</div>
