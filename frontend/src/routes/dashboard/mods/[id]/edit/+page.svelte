<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { isAuthenticated } from '$lib/stores/auth';
	import { modsApi, gamesApi } from '$lib/api/client';
	import { toast } from '$lib/stores/notifications';
	import Loading from '$lib/components/Loading.svelte';
	import {
		Package,
		ArrowLeft,
		Save,
		Tag,
		Info,
		AlertTriangle,
		CheckCircle,
		Trash2,
		Upload,
		Eye
	} from 'lucide-svelte';

	// URL parameter
	$: modId = parseInt($page.params.id);

	// Form data
	let formData = {
		name: '',
		short_description: '',
		description: '',
		version: '',
		game_version: '',
		game_id: 0,
		source_website: '',
		contact_info: '',
		tags: [] as string[],
		dependencies: [] as number[]
	};

	// Original mod data
	let originalMod: {
		id: number;
		name: string;
		short_description: string;
		description: string;
		version: string;
		game_version: string;
		game_id: number;
		source_website: string;
		contact_info: string;
		tags: { slug: string; name: string }[];
		dependencies: { id: number }[];
		files: any[];
		game: { slug: string; name: string };
		slug: string;
		is_rejected: boolean;
		is_scanned: boolean;
		downloads: number;
		likes: number;
	} | null = null;

	// Available data
	let games: { id: number; name: string; slug: string }[] = [];
	let tagInput = '';
	let selectedGameTags: { id: number; name: string; slug: string }[] = [];

	// UI state
	let isSubmitting = false;
	let isLoading = true;
	let errors: { [key: string]: string } = {};

	// File management
	let modFiles: {
		id: number;
		filename: string;
		file_size: number;
		created_at: string;
		is_main: boolean;
	}[] = [];
	let fileInput: HTMLInputElement;
	let selectedFile: File | null = null;

	// Icon upload
	let iconInput: HTMLInputElement;
	let selectedIcon: File | null = null;
	let iconPreview: string | null = null;
	let currentIcon: string = '';
	let fileError = '';
	let isUploadingFile = false;

	// Load mod data
	onMount(async () => {
		if (!$isAuthenticated) {
			goto('/auth/login?redirect=/dashboard/mods/' + modId + '/edit');
			return;
		}

		if (isNaN(modId)) {
			toast.error('Invalid Mod', 'The mod ID is not valid.');
			goto('/dashboard');
			return;
		}

		await loadModData();
		await loadGames();
		isLoading = false;
	});

	// Load mod data
	async function loadModData() {
		try {
			const response = await modsApi.getModById(modId);
			if (response.success && response.data) {
				originalMod = response.data;

				// Populate form data
				formData = {
					name: originalMod.name || '',
					short_description: originalMod.short_description || '',
					description: originalMod.description || '',
					version: originalMod.version || '',
					game_version: originalMod.game_version || '',
					game_id: originalMod.game_id || 0,
					source_website: originalMod.source_website || '',
					contact_info: originalMod.contact_info || '',
					tags:
						originalMod.tags?.map((tag: { slug: string; name: string }) => tag.slug || tag.name) ||
						[],
					dependencies: originalMod.dependencies?.map((dep: { id: number }) => dep.id) || []
				};

				modFiles = originalMod.files || [];
				currentIcon = (originalMod as any).icon || '';

				// Load game tags
				if (originalMod.game) {
					await loadGameTags(originalMod.game.slug);
				}
			} else {
				toast.error('Mod Not Found', response.error || 'The requested mod could not be found.');
				goto('/dashboard');
			}
		} catch (error) {
			console.error('Error loading mod data:', error);
			toast.error('Error', 'Failed to load mod data');
			goto('/dashboard');
		}
	}

	// Load games
	async function loadGames() {
		try {
			const response = await gamesApi.getGames({ per_page: 100 });
			if (response.success && response.data) {
				games = response.data.data || [];
			}
		} catch (error) {
			console.error('Error loading games:', error);
		}
	}

	// Load game tags
	async function loadGameTags(gameSlug: string) {
		try {
			const response = await gamesApi.getGameTags(gameSlug);
			if (response.success && response.data) {
				selectedGameTags = response.data;
			}
		} catch (error) {
			console.error('Error loading game tags:', error);
		}
	}

	// Handle game change
	async function onGameChange() {
		if (formData.game_id) {
			const game = games.find((g) => g.id === formData.game_id);
			if (game) {
				await loadGameTags(game.slug);
			}
		} else {
			selectedGameTags = [];
		}
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

	// Upload new file
	async function uploadNewFile() {
		if (!selectedFile) {
			toast.error('No File Selected', 'Please select a file to upload');
			return;
		}

		isUploadingFile = true;

		try {
			const response = await modsApi.uploadModFile(modId, selectedFile, false);
			if (response.success) {
				toast.success('File Uploaded', 'New file has been uploaded successfully');
				selectedFile = null;
				if (fileInput) fileInput.value = '';
				// Reload mod data to show new file
				await loadModData();
			} else {
				toast.error('Upload Failed', response.error || 'Failed to upload file');
			}
		} catch (error) {
			console.error('Error uploading file:', error);
			toast.error('Upload Error', 'An error occurred while uploading the file');
		} finally {
			isUploadingFile = false;
		}
	}

	// Delete file
	async function deleteFile(fileId: number, filename: string) {
		if (!confirm(`Are you sure you want to delete "${filename}"? This action cannot be undone.`)) {
			return;
		}

		try {
			// Note: Need to implement deleteModFile endpoint
			toast.error('Feature Not Available', 'File deletion is not yet implemented');
		} catch (error) {
			console.error('Error deleting file:', error);
			toast.error('Delete Error', 'Failed to delete file');
		}
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
	function validateForm(): boolean {
		errors = {};

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

		if (!formData.description.trim()) {
			errors.description = 'Detailed description is required';
		} else if (formData.description.trim().length < 50) {
			errors.description = 'Description must be at least 50 characters';
		}

		if (!formData.version.trim()) {
			errors.version = 'Version is required';
		}

		if (!formData.game_version.trim()) {
			errors.game_version = 'Game version is required';
		}

		if (!formData.game_id) {
			errors.game_id = 'Please select a game';
		}

		return Object.keys(errors).length === 0;
	}

	// Submit form
	async function handleSubmit(event: Event) {
		// Prevent default form submission
		event.preventDefault();

		if (!validateForm()) {
			toast.error('Validation Error', 'Please fix the errors below');
			return;
		}

		isSubmitting = true;

		try {
			// Upload icon if selected
			if (selectedIcon) {
				const iconResponse = await modsApi.uploadModIcon(modId, selectedIcon);
				if (!iconResponse.success) {
					console.warn('Icon upload failed:', iconResponse.error);
					// Don't fail the whole process if icon upload fails
				}
			}

			const response = await modsApi.updateMod(modId, {
				name: formData.name.trim(),
				description: formData.description.trim(),
				short_description: formData.short_description.trim(),
				version: formData.version.trim(),
				game_version: formData.game_version.trim(),
				source_website: formData.source_website.trim() || undefined,
				contact_info: formData.contact_info.trim() || undefined,
				tags: formData.tags,
				dependencies: formData.dependencies
			});

			if (response.success) {
				toast.success('Mod Updated', 'Your mod has been updated successfully');
				goto('/dashboard');
			} else {
				toast.error('Update Failed', response.error || 'Failed to update mod');
			}
		} catch (error) {
			console.error('Error updating mod:', error);
			toast.error('Error', 'An unexpected error occurred. Please try again.');
		} finally {
			isSubmitting = false;
		}
	}

	// Delete mod
	async function deleteMod() {
		if (
			!confirm(
				`Are you sure you want to delete "${originalMod?.name}"? This action cannot be undone and will remove all files and data associated with this mod.`
			)
		) {
			return;
		}

		try {
			const response = await modsApi.deleteMod(modId);
			if (response.success) {
				toast.success('Mod Deleted', 'Your mod has been deleted successfully');
				goto('/dashboard');
			} else {
				toast.error('Delete Failed', response.error || 'Failed to delete mod');
			}
		} catch (error) {
			console.error('Error deleting mod:', error);
			toast.error('Error', 'Failed to delete mod');
		}
	}

	// Format file size
	function formatFileSize(bytes: number): string {
		if (bytes === 0) return '0 Bytes';
		const k = 1024;
		const sizes = ['Bytes', 'KB', 'MB', 'GB'];
		const i = Math.floor(Math.log(bytes) / Math.log(k));
		return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
	}

	// Format date
	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString();
	}
</script>

<svelte:head>
	<title>Edit Mod - Azurite</title>
	<meta name="description" content="Edit your mod on the Azurite platform" />
</svelte:head>

{#if isLoading}
	<div class="min-h-screen bg-background-primary flex items-center justify-center">
		<Loading size="lg" text="Loading mod data..." />
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

				<div class="flex items-center justify-between">
					<div>
						<h1 class="text-3xl font-bold text-text-primary">Edit Mod</h1>
						<p class="text-text-secondary mt-1">Update your mod information and manage files</p>
					</div>

					<div class="flex items-center space-x-3">
						{#if originalMod?.game?.slug}
							<a
								href="/mods/{originalMod.game.slug}/{originalMod.slug}"
								class="btn btn-outline"
								title="View Mod"
							>
								<Eye class="w-4 h-4 mr-2" />
								View
							</a>
						{/if}
						<button onclick={deleteMod} class="btn btn-danger" title="Delete Mod">
							<Trash2 class="w-4 h-4 mr-2" />
							Delete
						</button>
					</div>
				</div>
			</div>

			<div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
				<!-- Main Form -->
				<div class="lg:col-span-2 space-y-6">
					<!-- Basic Information -->
					<div class="card">
						<div class="p-6">
							<h2 class="text-xl font-semibold text-text-primary mb-6">Basic Information</h2>

							<form onsubmit={handleSubmit} class="space-y-4">
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
										placeholder="A brief description of what your mod does"
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
										Mod Icon
									</label>
									<div class="flex items-center gap-4">
										{#if iconPreview || currentIcon}
											<div class="relative">
												<img
													src={iconPreview || currentIcon}
													alt="Icon preview"
													class="w-20 h-20 rounded-lg object-cover border-2 border-slate-600"
												/>
												{#if iconPreview}
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
												{/if}
											</div>
										{/if}
										<label
											class="flex-1 flex flex-col items-center px-4 py-6 bg-background-secondary border-2 border-dashed border-slate-600 rounded-lg cursor-pointer hover:border-primary-500 transition-colors"
										>
											<Upload class="w-8 h-8 text-text-muted mb-2" />
											<span class="text-sm text-text-secondary">
												{selectedIcon ? 'Change Icon' : currentIcon ? 'Update Icon' : 'Upload Icon'}
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

								<!-- Submit Button -->
								<div class="flex justify-end pt-4">
									<button
										type="submit"
										disabled={isSubmitting}
										class="btn btn-primary min-w-[140px]"
									>
										{#if isSubmitting}
											<Loading size="sm" />
											Updating...
										{:else}
											<Save class="w-4 h-4 mr-2" />
											Update Mod
										{/if}
									</button>
								</div>
							</form>
						</div>
					</div>

					<!-- Detailed Information -->
					<div class="card">
						<div class="p-6">
							<h2 class="text-xl font-semibold text-text-primary mb-6">Detailed Information</h2>

							<div class="space-y-4">
								<!-- Description -->
								<div>
									<label for="description" class="block text-sm font-medium text-text-primary mb-2">
										Detailed Description *
									</label>
									<textarea
										id="description"
										bind:value={formData.description}
										placeholder="Provide a comprehensive description of your mod..."
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
									</div>
								</div>
							</div>
						</div>
					</div>

					<!-- File Management -->
					<div class="card">
						<div class="p-6">
							<h2 class="text-xl font-semibold text-text-primary mb-6">File Management</h2>

							<!-- Current Files -->
							{#if modFiles.length > 0}
								<div class="mb-6">
									<h3 class="text-lg font-medium text-text-primary mb-3">Current Files</h3>
									<div class="space-y-2">
										{#each modFiles as file (file.id)}
											<div class="flex items-center justify-between p-3 bg-slate-800 rounded-lg">
												<div class="flex items-center">
													<Package class="w-5 h-5 text-text-muted mr-3" />
													<div>
														<p class="text-text-primary font-medium">{file.filename}</p>
														<p class="text-text-muted text-sm">
															{formatFileSize(file.file_size)} • Uploaded {formatDate(
																file.created_at
															)}
															{#if file.is_main}
																<span
																	class="ml-2 px-2 py-0.5 bg-primary-500/20 text-primary-300 text-xs rounded"
																	>Main</span
																>
															{/if}
														</p>
													</div>
												</div>
												<button
													onclick={() => deleteFile(file.id, file.filename)}
													class="btn btn-outline btn-sm text-red-400 hover:text-red-300 hover:border-red-500"
													title="Delete File"
												>
													<Trash2 class="w-4 h-4" />
												</button>
											</div>
										{/each}
									</div>
								</div>
							{/if}

							<!-- Upload New File -->
							<div>
								<h3 class="text-lg font-medium text-text-primary mb-3">Upload New File</h3>
								<div class="border-2 border-dashed border-slate-600 rounded-lg p-6 text-center">
									<input
										type="file"
										bind:this={fileInput}
										onchange={handleFileSelect}
										accept=".jar,.zip,.rar,.7z,.tar.gz"
										class="hidden"
									/>

									{#if selectedFile}
										<div class="space-y-2">
											<CheckCircle class="w-8 h-8 text-green-400 mx-auto" />
											<p class="text-text-primary font-medium">{selectedFile.name}</p>
											<p class="text-text-muted text-sm">
												{formatFileSize(selectedFile.size)}
											</p>
											<div class="flex justify-center space-x-2">
												<button
													type="button"
													onclick={uploadNewFile}
													disabled={isUploadingFile}
													class="btn btn-primary btn-sm"
												>
													{#if isUploadingFile}
														<Loading size="sm" />
														Uploading...
													{:else}
														<Upload class="w-4 h-4 mr-1" />
														Upload
													{/if}
												</button>
												<button
													type="button"
													onclick={() => fileInput.click()}
													class="btn btn-outline btn-sm"
												>
													Change File
												</button>
											</div>
										</div>
									{:else}
										<div class="space-y-2">
											<Upload class="w-8 h-8 text-text-muted mx-auto" />
											<p class="text-text-primary">Choose a file to upload</p>
											<p class="text-text-muted text-sm">
												Supports .jar, .zip, .rar, .7z, .tar.gz files up to 100MB
											</p>
											<button
												type="button"
												onclick={() => fileInput.click()}
												class="btn btn-outline"
											>
												Choose File
											</button>
										</div>
									{/if}
								</div>

								{#if fileError}
									<p class="mt-2 text-sm text-red-400">{fileError}</p>
								{/if}
							</div>
						</div>
					</div>
				</div>

				<!-- Sidebar -->
				<div class="space-y-6">
					<!-- Status -->
					{#if originalMod}
						<div class="card">
							<div class="p-4">
								<h3 class="text-lg font-semibold text-text-primary mb-3">Status</h3>
								<div class="space-y-2">
									<div class="flex items-center justify-between">
										<span class="text-text-muted">Current Status:</span>
										<span
											class="px-2 py-1 text-xs rounded
											{originalMod.is_rejected
												? 'bg-red-500/20 text-red-300'
												: originalMod.is_scanned
													? 'bg-green-500/20 text-green-300'
													: 'bg-yellow-500/20 text-yellow-300'}"
										>
											{originalMod.is_rejected
												? 'Rejected'
												: originalMod.is_scanned
													? 'Approved'
													: 'Pending'}
										</span>
									</div>
									<div class="flex items-center justify-between">
										<span class="text-text-muted">Downloads:</span>
										<span class="text-text-primary">{originalMod.downloads}</span>
									</div>
									<div class="flex items-center justify-between">
										<span class="text-text-muted">Likes:</span>
										<span class="text-text-primary">{originalMod.likes}</span>
									</div>
								</div>
							</div>
						</div>
					{/if}

					<!-- Guidelines -->
					<div class="card bg-blue-900/20 border-blue-600/50">
						<div class="p-4">
							<div class="flex items-center mb-3">
								<Info class="w-5 h-5 text-blue-400 mr-2" />
								<h3 class="text-lg font-semibold text-blue-300">Update Guidelines</h3>
							</div>
							<ul class="text-blue-200 text-sm space-y-1">
								<li>• Major changes may require re-review</li>
								<li>• Keep your description up to date</li>
								<li>• Update version numbers appropriately</li>
								<li>• Test changes before uploading new files</li>
							</ul>
						</div>
					</div>

					<!-- Danger Zone -->
					<div class="card bg-red-900/20 border-red-600/50">
						<div class="p-4">
							<div class="flex items-center mb-3">
								<AlertTriangle class="w-5 h-5 text-red-400 mr-2" />
								<h3 class="text-lg font-semibold text-red-300">Danger Zone</h3>
							</div>
							<p class="text-red-200 text-sm mb-3">
								Deleting your mod is permanent and cannot be undone.
							</p>
							<button onclick={deleteMod} class="btn btn-danger btn-sm w-full">
								<Trash2 class="w-4 h-4 mr-2" />
								Delete Mod
							</button>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
{/if}
