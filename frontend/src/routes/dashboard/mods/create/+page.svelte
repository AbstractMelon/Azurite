<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { isAuthenticated } from '$lib/stores/auth';
	import { modsApi, gamesApi } from '$lib/api/client';
	import { toast } from '$lib/stores/notifications';
	import Loading from '$lib/components/Loading.svelte';
	import {
		Package,
		ArrowLeft,
		Upload,
		Tag,
		Info,
		AlertTriangle,
		CheckCircle,
		FileText
	} from 'lucide-svelte';

	// Form data
	let formData = {
		name: '',
		short_description: '',
		description: '',
		version: '1.0.0',
		game_version: '',
		game_id: 0,
		source_website: '',
		contact_info: '',
		tags: [] as string[],
		dependencies: [] as number[]
	};

	// Available data
	let games: any[] = [];
	let tagInput = '';
	let selectedGameTags: { id: number; name: string; slug: string }[] = [];

	// UI state
	let isSubmitting = false;
	let isLoading = true;
	let errors: { [key: string]: string } = {};
	let currentStep = 1;
	let totalSteps = 3;

	// File upload
	let fileInput: HTMLInputElement;
	let selectedFile: File | null = null;
	let fileError = '';

	// Icon upload
	let iconInput: HTMLInputElement;
	let selectedIcon: File | null = null;
	let iconPreview: string | null = null;

	// Load initial data
	onMount(async () => {
		if (!$isAuthenticated) {
			goto('/auth/login?redirect=/dashboard/mods/create');
			return;
		}

		try {
			const gamesResponse = await gamesApi.getGames({ per_page: 100 });
			if (gamesResponse.success && gamesResponse.data) {
				games = gamesResponse.data.data || [];
			}
		} catch (error) {
			console.error('Error loading games:', error);
			toast.error('Error', 'Failed to load available games');
		}

		isLoading = false;
	});

	// Load tags when game is selected
	async function onGameChange() {
		if (formData.game_id) {
			try {
				const game = games.find((g) => g.id === formData.game_id);
				if (game) {
					const tagsResponse = await gamesApi.getGameTags(game.slug);
					if (tagsResponse.success && tagsResponse.data) {
						selectedGameTags = tagsResponse.data;
					}
				}
			} catch (error) {
				console.error('Error loading game tags:', error);
			}
		} else {
			selectedGameTags = [];
		}
		formData.tags = [];
	}

	// Handle icon selection
	function handleIconSelect(event: Event) {
		const input = event.target as HTMLInputElement;
		const file = input.files?.[0];

		if (!file) {
			selectedIcon = null;
			iconPreview = null;
			return;
		}

		// Validate file type
		if (!file.type.startsWith('image/')) {
			toast.error('Invalid File', 'Please select an image file');
			selectedIcon = null;
			iconPreview = null;
			return;
		}

		// Validate file size (max 2MB)
		if (file.size > 2 * 1024 * 1024) {
			toast.error('File Too Large', 'Icon must be less than 2MB');
			selectedIcon = null;
			iconPreview = null;
			return;
		}

		selectedIcon = file;

		// Create preview
		const reader = new FileReader();
		reader.onload = (e) => {
			iconPreview = e.target?.result as string;
		};
		reader.readAsDataURL(file);
	}

	// Handle file selection
	function handleFileSelect(event: Event) {
		const input = event.target as HTMLInputElement;
		const file = input.files?.[0];

		if (!file) {
			selectedFile = null;
			fileError = '';
			return;
		}

		// Validate file
		const maxSize = 100 * 1024 * 1024; // 100MB
		const allowedTypes = [
			'application/java-archive',
			'application/zip',
			'application/x-rar-compressed',
			'application/x-7z-compressed',
			'application/x-tar'
		];

		if (file.size > maxSize) {
			fileError = 'File size must be less than 100MB';
			selectedFile = null;
			input.value = '';
			return;
		}

		if (!allowedTypes.includes(file.type) && !file.name.match(/\.(jar|zip|rar|7z|tar\.gz)$/i)) {
			fileError = 'File must be a .jar, .zip, .rar, .7z, or .tar.gz file';
			selectedFile = null;
			input.value = '';
			return;
		}

		selectedFile = file;
		fileError = '';
	}

	// Add tag
	function addTag() {
		const tag = tagInput.trim().toLowerCase();
		if (tag && !formData.tags.includes(tag)) {
			formData.tags = [...formData.tags, tag];
			tagInput = '';
		}
	}

	// Remove tag
	function removeTag(tagToRemove: string) {
		formData.tags = formData.tags.filter((tag) => tag !== tagToRemove);
	}

	// Add predefined tag
	function addPredefinedTag(tag: { id: number; name: string; slug: string }) {
		if (!formData.tags.includes(tag.slug)) {
			formData.tags = [...formData.tags, tag.slug];
		}
	}

	// Form validation
	function validateStep(step: number): boolean {
		errors = {};

		if (step === 1) {
			if (!formData.name.trim()) {
				errors.name = 'Mod name is required';
			} else if (formData.name.trim().length < 3) {
				errors.name = 'Mod name must be at least 3 characters';
			} else if (formData.name.trim().length > 200) {
				errors.name = 'Mod name must be less than 200 characters';
			}

			if (!formData.short_description.trim()) {
				errors.short_description = 'Short description is required';
			} else if (formData.short_description.trim().length > 500) {
				errors.short_description = 'Short description must be less than 500 characters';
			}

			if (!formData.game_id) {
				errors.game_id = 'Please select a game';
			}

			if (!formData.version.trim()) {
				errors.version = 'Version is required';
			}

			if (!formData.game_version.trim()) {
				errors.game_version = 'Game version is required';
			}
		}

		if (step === 2) {
			if (!formData.description.trim()) {
				errors.description = 'Detailed description is required';
			} else if (formData.description.trim().length < 50) {
				errors.description = 'Description must be at least 50 characters';
			}
		}

		if (step === 3) {
			if (!selectedFile) {
				errors.file = 'Please select a mod file to upload';
			}
		}

		return Object.keys(errors).length === 0;
	}

	// Next step
	function nextStep() {
		if (validateStep(currentStep)) {
			currentStep++;
		}
	}

	// Previous step
	function prevStep() {
		currentStep--;
	}

	// Submit form
	async function handleSubmit(event: Event) {
		// Prevent default form submission
		event.preventDefault();

		if (!validateStep(3)) {
			toast.error('Validation Error', 'Please fix the errors below');
			return;
		}

		if (!selectedFile) {
			toast.error('File Required', 'Please select a mod file to upload');
			return;
		}

		isSubmitting = true;

		try {
			// First create the mod
			const modResponse = await modsApi.createMod({
				name: formData.name.trim(),
				description: formData.description.trim(),
				short_description: formData.short_description.trim(),
				version: formData.version.trim(),
				game_version: formData.game_version.trim(),
				game_id: formData.game_id,
				source_website: formData.source_website.trim() || undefined,
				contact_info: formData.contact_info.trim() || undefined,
				tags: formData.tags,
				dependencies: formData.dependencies
			});

			if (!modResponse.success) {
				toast.error('Creation Failed', modResponse.error || 'Failed to create mod');
				return;
			}

			const modId = modResponse.data.id;

			// Upload the icon if selected
			if (selectedIcon) {
				const iconResponse = await modsApi.uploadModIcon(modId, selectedIcon);
				if (!iconResponse.success) {
					console.warn('Icon upload failed:', iconResponse.error);
					// Don't fail the whole process if icon upload fails
				}
			}

			// Upload the file
			const fileResponse = await modsApi.uploadModFile(modId, selectedFile, true);

			if (!fileResponse.success) {
				toast.error('Upload Failed', fileResponse.error || 'Failed to upload mod file');
				// Note: The mod was created but file upload failed
				// We might want to handle this differently
				return;
			}

			toast.success(
				'Mod Created',
				"Your mod has been uploaded and is pending review. You will be notified when it's approved."
			);
			goto('/dashboard');
		} catch (error) {
			console.error('Error creating mod:', error);
			toast.error('Error', 'An unexpected error occurred. Please try again.');
		} finally {
			isSubmitting = false;
		}
	}

	// $: stepValid = validateStep(currentStep);
</script>

<svelte:head>
	<title>Upload New Mod - Azurite</title>
	<meta name="description" content="Upload a new mod to the Azurite platform" />
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
					onclick={() => goto('/dashboard')}
					class="flex items-center text-text-muted hover:text-text-secondary mb-6"
				>
					<ArrowLeft class="w-4 h-4 mr-2" />
					Back to Dashboard
				</button>

				<div class="text-center">
					<div
						class="inline-flex items-center justify-center w-16 h-16 bg-primary-500/20 rounded-full mb-6"
					>
						<Package class="w-8 h-8 text-primary-400" />
					</div>
					<h1 class="text-4xl font-bold text-text-primary mb-4">Upload New Mod</h1>
					<p class="text-text-secondary text-lg max-w-2xl mx-auto">
						Share your creation with the community! Fill out the information below to upload your
						mod.
					</p>
				</div>
			</div>

			<!-- Progress Indicator -->
			<div class="mb-8">
				<div class="flex items-center justify-center space-x-4">
					{#each Array(totalSteps) as _, i (i)}
						<div class="flex items-center">
							<div
								class="w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium
								{i + 1 < currentStep
									? 'bg-primary-500 text-white'
									: i + 1 === currentStep
										? 'bg-primary-500 text-white'
										: 'bg-slate-600 text-text-muted'}"
							>
								{i + 1 < currentStep ? '✓' : i + 1}
							</div>
							{#if i < totalSteps - 1}
								<div
									class="w-12 h-0.5 mx-2 {i + 1 < currentStep ? 'bg-primary-500' : 'bg-slate-600'}"
								></div>
							{/if}
						</div>
					{/each}
				</div>
				<div class="flex justify-center mt-2">
					<span class="text-sm text-text-muted">
						Step {currentStep} of {totalSteps}
					</span>
				</div>
			</div>

			<!-- Form -->
			<div class="card">
				<div class="p-8">
					<form onsubmit={handleSubmit}>
						<!-- Step 1: Basic Information -->
						{#if currentStep === 1}
							<div class="space-y-6">
								<div class="text-center mb-6">
									<h2 class="text-2xl font-semibold text-text-primary">Basic Information</h2>
									<p class="text-text-secondary">Tell us about your mod</p>
								</div>

								<!-- Mod Name -->
								<div>
									<label for="name" class="block text-sm font-medium text-text-primary mb-2">
										Mod Name *
									</label>
									<input
										id="name"
										type="text"
										bind:value={formData.name}
										placeholder="e.g., Awesome Feature Pack"
										class="input w-full {errors.name ? 'border-red-500' : ''}"
										required
										maxlength="200"
									/>
									{#if errors.name}
										<p class="mt-1 text-sm text-red-400">{errors.name}</p>
									{/if}
								</div>

								<!-- Short Description -->
								<div>
									<label
										for="short_description"
										class="block text-sm font-medium text-text-primary mb-2"
									>
										Short Description *
									</label>
									<textarea
										id="short_description"
										bind:value={formData.short_description}
										placeholder="A brief description of what your mod does (shown in mod lists)"
										class="textarea w-full h-20 {errors.short_description ? 'border-red-500' : ''}"
										required
										maxlength="500"
									></textarea>
									{#if errors.short_description}
										<p class="mt-1 text-sm text-red-400">{errors.short_description}</p>
									{/if}
									<p class="mt-1 text-sm text-text-muted">
										{formData.short_description.length}/500 characters
									</p>
								</div>

								<!-- Game Selection -->
								<div>
									<label for="game_id" class="block text-sm font-medium text-text-primary mb-2">
										Game *
									</label>
									<select
										id="game_id"
										bind:value={formData.game_id}
										onchange={onGameChange}
										class="select w-full {errors.game_id ? 'border-red-500' : ''}"
										required
									>
										<option value={0}>Select a game</option>
										{#each games as game (game.id)}
											<option value={game.id}>{game.name}</option>
										{/each}
									</select>
									{#if errors.game_id}
										<p class="mt-1 text-sm text-red-400">{errors.game_id}</p>
									{/if}
								</div>

								<!-- Mod Icon -->
								<div>
									<label class="block text-sm font-medium text-text-primary mb-2">
										Mod Icon (Optional)
									</label>
									<div class="flex items-center gap-4">
										{#if iconPreview}
											<div class="relative">
												<img
													src={iconPreview}
													alt="Icon preview"
													class="w-20 h-20 rounded-lg object-cover border-2 border-slate-600"
												/>
												<button
													type="button"
													onclick={() => {
														selectedIcon = null;
														iconPreview = null;
														if (iconInput) iconInput.value = '';
													}}
													class="absolute -top-2 -right-2 bg-red-500 text-white rounded-full w-6 h-6 flex items-center justify-center hover:bg-red-600"
												>
													×
												</button>
											</div>
										{/if}
										<label
											class="flex-1 flex flex-col items-center px-4 py-6 bg-background-secondary border-2 border-dashed border-slate-600 rounded-lg cursor-pointer hover:border-primary-500 transition-colors"
										>
											<Upload class="w-8 h-8 text-text-muted mb-2" />
											<span class="text-sm text-text-secondary">
												{selectedIcon ? 'Change Icon' : 'Upload Icon'}
											</span>
											<span class="text-xs text-text-muted mt-1">PNG, JPG up to 2MB</span>
											<input
												type="file"
												bind:this={iconInput}
												onchange={handleIconSelect}
												accept="image/*"
												class="hidden"
											/>
										</label>
									</div>
									<p class="mt-2 text-sm text-text-muted">
										Recommended: 256x256px square image
									</p>
								</div>

								<!-- Version Info -->
								<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
									<div>
										<label for="version" class="block text-sm font-medium text-text-primary mb-2">
											Mod Version *
										</label>
										<input
											id="version"
											type="text"
											bind:value={formData.version}
											placeholder="1.0.0"
											class="input w-full {errors.version ? 'border-red-500' : ''}"
											required
										/>
										{#if errors.version}
											<p class="mt-1 text-sm text-red-400">{errors.version}</p>
										{/if}
									</div>
									<div>
										<label
											for="game_version"
											class="block text-sm font-medium text-text-primary mb-2"
										>
											Game Version *
										</label>
										<input
											id="game_version"
											type="text"
											bind:value={formData.game_version}
											placeholder="1.19.2"
											class="input w-full {errors.game_version ? 'border-red-500' : ''}"
											required
										/>
										{#if errors.game_version}
											<p class="mt-1 text-sm text-red-400">{errors.game_version}</p>
										{/if}
									</div>
								</div>
							</div>
						{/if}

						<!-- Step 2: Detailed Information -->
						{#if currentStep === 2}
							<div class="space-y-6">
								<div class="text-center mb-6">
									<h2 class="text-2xl font-semibold text-text-primary">Detailed Information</h2>
									<p class="text-text-secondary">Provide more details about your mod</p>
								</div>

								<!-- Detailed Description -->
								<div>
									<label for="description" class="block text-sm font-medium text-text-primary mb-2">
										Detailed Description *
									</label>
									<textarea
										id="description"
										bind:value={formData.description}
										placeholder="Provide a comprehensive description of your mod, including features, installation instructions, and any other relevant information..."
										class="textarea w-full h-40 {errors.description ? 'border-red-500' : ''}"
										required
									></textarea>
									{#if errors.description}
										<p class="mt-1 text-sm text-red-400">{errors.description}</p>
									{/if}
									<p class="mt-1 text-sm text-text-muted">
										{formData.description.length} characters (minimum 50)
									</p>
								</div>

								<!-- Tags -->
								<div>
									<label for="tags" class="block text-sm font-medium text-text-primary mb-2">
										Tags
									</label>

									<!-- Predefined Tags -->
									{#if selectedGameTags.length > 0}
										<div class="mb-3">
											<p class="text-sm text-text-muted mb-2">Popular tags for this game:</p>
											<div class="flex flex-wrap gap-2">
												{#each selectedGameTags as tag (tag.id)}
													<button
														type="button"
														onclick={() => addPredefinedTag(tag)}
														class="px-3 py-1 bg-slate-700 hover:bg-slate-600 text-text-secondary text-sm rounded-full border border-slate-600 hover:border-primary-500 transition-colors"
														disabled={formData.tags.includes(tag.slug)}
													>
														{tag.name}
													</button>
												{/each}
											</div>
										</div>
									{/if}

									<!-- Custom Tag Input -->
									<div class="flex space-x-2">
										<input
											type="text"
											bind:value={tagInput}
											placeholder="Add a custom tag..."
											class="input flex-1"
											onkeypress={(e) => e.key === 'Enter' && (e.preventDefault(), addTag())}
										/>
										<button
											type="button"
											onclick={addTag}
											class="btn btn-outline px-4"
											disabled={!tagInput.trim()}
										>
											<Tag class="w-4 h-4" />
										</button>
									</div>

									<!-- Selected Tags -->
									{#if formData.tags.length > 0}
										<div class="mt-3 flex flex-wrap gap-2">
											{#each formData.tags as tag (tag)}
												<span
													class="inline-flex items-center px-3 py-1 bg-primary-500/20 text-primary-300 text-sm rounded-full"
												>
													{tag}
													<button
														type="button"
														onclick={() => removeTag(tag)}
														class="ml-2 hover:text-primary-200"
													>
														×
													</button>
												</span>
											{/each}
										</div>
									{/if}
								</div>

								<!-- Optional Fields -->
								<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
									<div>
										<label
											for="source_website"
											class="block text-sm font-medium text-text-primary mb-2"
										>
											Source Website
										</label>
										<input
											id="source_website"
											type="url"
											bind:value={formData.source_website}
											placeholder="https://github.com/username/mod"
											class="input w-full"
										/>
										<p class="mt-1 text-sm text-text-muted">
											Link to your mod's source code or development page
										</p>
									</div>
									<div>
										<label
											for="contact_info"
											class="block text-sm font-medium text-text-primary mb-2"
										>
											Contact Email
										</label>
										<input
											id="contact_info"
											type="email"
											bind:value={formData.contact_info}
											placeholder="your@email.com"
											class="input w-full"
										/>
										<p class="mt-1 text-sm text-text-muted">For support and bug reports</p>
									</div>
								</div>
							</div>
						{/if}

						<!-- Step 3: File Upload -->
						{#if currentStep === 3}
							<div class="space-y-6">
								<div class="text-center mb-6">
									<h2 class="text-2xl font-semibold text-text-primary">Upload Mod File</h2>
									<p class="text-text-secondary">Upload your mod file to complete the submission</p>
								</div>

								<!-- File Upload -->
								<div>
									<label class="block text-sm font-medium text-text-primary mb-2">
										Mod File *
									</label>
									<div
										class="border-2 border-dashed border-slate-600 rounded-lg p-8 text-center hover:border-primary-500 transition-colors"
									>
										<input
											type="file"
											bind:this={fileInput}
											onchange={handleFileSelect}
											accept=".jar,.zip,.rar,.7z,.tar.gz"
											class="hidden"
										/>

										{#if selectedFile}
											<div class="space-y-2">
												<CheckCircle class="w-12 h-12 text-green-400 mx-auto" />
												<p class="text-text-primary font-medium">{selectedFile.name}</p>
												<p class="text-text-muted text-sm">
													{(selectedFile.size / (1024 * 1024)).toFixed(2)} MB
												</p>
												<button
													type="button"
													onclick={() => fileInput.click()}
													class="btn btn-outline btn-sm"
												>
													Choose Different File
												</button>
											</div>
										{:else}
											<div class="space-y-2">
												<Upload class="w-12 h-12 text-text-muted mx-auto" />
												<p class="text-text-primary font-medium">Click to upload your mod file</p>
												<p class="text-text-muted text-sm">
													Supports .jar, .zip, .rar, .7z, .tar.gz files up to 100MB
												</p>
												<button
													type="button"
													onclick={() => fileInput.click()}
													class="btn btn-primary"
												>
													Choose File
												</button>
											</div>
										{/if}
									</div>

									{#if errors.file || fileError}
										<p class="mt-1 text-sm text-red-400">{errors.file || fileError}</p>
									{/if}
								</div>

								<!-- Upload Guidelines -->
								<div class="bg-blue-900/20 border border-blue-600/50 rounded-lg p-4">
									<div class="flex items-start">
										<Info class="w-5 h-5 text-blue-400 mr-3 mt-0.5 flex-shrink-0" />
										<div>
											<h3 class="text-blue-300 font-medium mb-2">Upload Guidelines</h3>
											<ul class="text-blue-200 text-sm space-y-1">
												<li>• Ensure your mod is compatible with the specified game version</li>
												<li>• Include installation instructions in your description</li>
												<li>• Test your mod thoroughly before uploading</li>
												<li>• Make sure you have permission to distribute any included assets</li>
											</ul>
										</div>
									</div>
								</div>
							</div>
						{/if}

						<!-- Navigation -->
						<div class="flex justify-between mt-8 pt-6 border-t border-slate-700">
							<div>
								{#if currentStep > 1}
									<button type="button" onclick={prevStep} class="btn btn-outline">
										<ArrowLeft class="w-4 h-4 mr-2" />
										Previous
									</button>
								{/if}
							</div>

							<div>
								{#if currentStep < totalSteps}
									<button type="button" onclick={nextStep} class="btn btn-primary">
										Next
										<ArrowLeft class="w-4 h-4 ml-2 rotate-180" />
									</button>
								{:else}
									<button
										type="submit"
										disabled={isSubmitting}
										class="btn btn-primary min-w-[140px]"
									>
										{#if isSubmitting}
											<Loading size="sm" />
											Uploading...
										{:else}
											<Upload class="w-4 h-4 mr-2" />
											Upload Mod
										{/if}
									</button>
								{/if}
							</div>
						</div>
					</form>
				</div>
			</div>

			<!-- Help Section -->
			<div class="mt-8 grid grid-cols-1 md:grid-cols-3 gap-4">
				<div class="card bg-yellow-900/20 border-yellow-600/50">
					<div class="p-4">
						<AlertTriangle class="w-6 h-6 text-yellow-400 mb-3" />
						<h3 class="font-semibold text-yellow-300 mb-2">Review Process</h3>
						<p class="text-yellow-200 text-sm">
							All mods go through a review process before being published. This ensures quality and
							safety for all users.
						</p>
					</div>
				</div>

				<div class="card bg-green-900/20 border-green-600/50">
					<div class="p-4">
						<CheckCircle class="w-6 h-6 text-green-400 mb-3" />
						<h3 class="font-semibold text-green-300 mb-2">Quality Tips</h3>
						<p class="text-green-200 text-sm">
							Clear descriptions, appropriate tags, and thorough testing will help your mod get
							approved faster.
						</p>
					</div>
				</div>

				<div class="card">
					<div class="p-4">
						<FileText class="w-6 h-6 text-primary-400 mb-3" />
						<h3 class="font-semibold text-text-primary mb-2">Need Help?</h3>
						<p class="text-text-secondary text-sm mb-3">
							Check our documentation for detailed uploading guidelines.
						</p>
						<a href="/docs" class="text-primary-400 hover:text-primary-300 text-sm">
							View Docs →
						</a>
					</div>
				</div>
			</div>
		</div>
	</div>
{/if}
