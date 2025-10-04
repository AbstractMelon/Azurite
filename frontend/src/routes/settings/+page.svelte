<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { auth, user, isAuthenticated } from '$lib/stores/auth';
	import { authApi } from '$lib/api/client';
	import { toast } from '$lib/stores/notifications';
	import Loading from '$lib/components/Loading.svelte';
	import {
		User,
		Lock,
		Bell,
		Eye,
		EyeOff,
		Save,
		Trash2,
		AlertTriangle,
		Check,
		X,
		Upload,
		Camera
	} from 'lucide-svelte';

	// Form data
	let profileData = {
		display_name: '',
		bio: '',
		notify_email: true,
		notify_in_site: true
	};

	let passwordData = {
		current_password: '',
		new_password: '',
		confirm_password: ''
	};

	// UI state
	let activeTab = 'profile';
	let isLoadingProfile = true;
	let isSavingProfile = false;
	let isChangingPassword = false;
	let showCurrentPassword = false;
	let showNewPassword = false;
	let showConfirmPassword = false;
	let profileErrors: { [key: string]: string } = {};
	let passwordErrors: { [key: string]: string } = {};

	// Avatar upload
	let avatarFile: File | null = null;
	let avatarPreview: string | null = null;
	let isUploadingAvatar = false;

	const tabs = [
		{ id: 'profile', label: 'Profile', icon: User },
		{ id: 'security', label: 'Security', icon: Lock },
		{ id: 'notifications', label: 'Notifications', icon: Bell }
	];

	onMount(async () => {
		if (!$isAuthenticated) {
			goto('/auth/login?redirect=/settings');
			return;
		}

		await loadProfile();
	});

	// Load current profile data
	async function loadProfile() {
		isLoadingProfile = true;
		try {
			if ($user) {
				profileData = {
					display_name: $user.display_name || '',
					bio: $user.bio || '',
					notify_email: $user.notify_email ?? true,
					notify_in_site: $user.notify_in_site ?? true
				};
			}
		} catch (error) {
			console.error('Error loading profile:', error);
			toast.error('Error', 'Failed to load profile settings');
		} finally {
			isLoadingProfile = false;
		}
	}

	// Save profile changes
	async function saveProfile(event: Event) {
		// Prevent default form submission
		event.preventDefault();

		// Clear previous errors
		profileErrors = {};

		// Validation
		if (!profileData.display_name.trim()) {
			profileErrors.display_name = 'Display name is required';
		} else if (profileData.display_name.length < 2) {
			profileErrors.display_name = 'Display name must be at least 2 characters';
		}

		if (profileData.bio.length > 500) {
			profileErrors.bio = 'Bio must be less than 500 characters';
		}

		if (Object.keys(profileErrors).length > 0) {
			return;
		}

		isSavingProfile = true;
		try {
			const response = await authApi.updateProfile({
				display_name: profileData.display_name.trim(),
				bio: profileData.bio.trim(),
				notify_email: profileData.notify_email,
				notify_in_site: profileData.notify_in_site
			});

			if (response.success && response.data) {
				auth.updateUser(response.data);
				toast.success('Profile updated', 'Your profile has been successfully updated.');
			} else {
				toast.error('Update failed', response.error || 'Failed to update profile');
			}
		} catch (error) {
			console.error('Error updating profile:', error);
			toast.error('Error', 'Failed to update profile');
		} finally {
			isSavingProfile = false;
		}
	}

	// Change password
	async function changePassword(event: Event) {
		// Prevent default form submission
		event.preventDefault();

		// Clear previous errors
		passwordErrors = {};

		// Validation
		if (!passwordData.current_password) {
			passwordErrors.current_password = 'Current password is required';
		}

		if (!passwordData.new_password) {
			passwordErrors.new_password = 'New password is required';
		} else if (passwordData.new_password.length < 8) {
			passwordErrors.new_password = 'New password must be at least 8 characters';
		} else if (!isStrongPassword(passwordData.new_password)) {
			passwordErrors.new_password =
				'Password must contain uppercase, lowercase, number, and special character';
		}

		if (!passwordData.confirm_password) {
			passwordErrors.confirm_password = 'Please confirm your new password';
		} else if (passwordData.new_password !== passwordData.confirm_password) {
			passwordErrors.confirm_password = 'Passwords do not match';
		}

		if (Object.keys(passwordErrors).length > 0) {
			return;
		}

		isChangingPassword = true;
		try {
			const response = await authApi.updatePassword({
				current_password: passwordData.current_password,
				new_password: passwordData.new_password
			});

			if (response.success) {
				passwordData = {
					current_password: '',
					new_password: '',
					confirm_password: ''
				};
				toast.success('Password changed', 'Your password has been successfully changed.');
			} else {
				toast.error('Failed to change password', response.error || 'Invalid current password');
			}
		} catch (error) {
			console.error('Error changing password:', error);
			toast.error('Error', 'Failed to change password');
		} finally {
			isChangingPassword = false;
		}
	}

	// Strong password validation
	function isStrongPassword(password: string): boolean {
		const hasUpperCase = /[A-Z]/.test(password);
		const hasLowerCase = /[a-z]/.test(password);
		const hasNumbers = /\d/.test(password);
		const hasSpecialChar = /[!@#$%^&*(),.?":{}|<>]/.test(password);
		return hasUpperCase && hasLowerCase && hasNumbers && hasSpecialChar;
	}

	// Handle avatar file selection
	function handleAvatarChange(event: Event) {
		const target = event.target as HTMLInputElement;
		const file = target.files?.[0];

		if (file) {
			// Validate file type
			if (!file.type.startsWith('image/')) {
				toast.error('Invalid file', 'Please select an image file');
				return;
			}

			// Validate file size (max 2MB)
			if (file.size > 2 * 1024 * 1024) {
				toast.error('File too large', 'Please select an image smaller than 2MB');
				return;
			}

			avatarFile = file;

			// Create preview
			const reader = new FileReader();
			reader.onload = (e) => {
				avatarPreview = e.target?.result as string;
			};
			reader.readAsDataURL(file);
		}
	}

	// Upload avatar
	async function uploadAvatar() {
		if (!avatarFile) return;

		isUploadingAvatar = true;
		try {
			// This would be implemented with a file upload API
			// const formData = new FormData();
			// formData.append('avatar', avatarFile);
			// const response = await authApi.uploadAvatar(formData);

			// Mock success for now
			await new Promise((resolve) => setTimeout(resolve, 2000));

			toast.success('Avatar updated', 'Your profile picture has been updated.');
			avatarFile = null;
			avatarPreview = null;
		} catch (error) {
			console.error('Error uploading avatar:', error);
			toast.error('Error', 'Failed to upload avatar');
		} finally {
			isUploadingAvatar = false;
		}
	}

	// Remove avatar preview
	function removeAvatarPreview() {
		avatarFile = null;
		avatarPreview = null;
	}

	// Toggle password visibility
	function togglePasswordVisibility(field: string) {
		switch (field) {
			case 'current':
				showCurrentPassword = !showCurrentPassword;
				break;
			case 'new':
				showNewPassword = !showNewPassword;
				break;
			case 'confirm':
				showConfirmPassword = !showConfirmPassword;
				break;
		}
	}

	// Real-time validation
	function validateField(field: string, value: string) {
		switch (field) {
			case 'display_name':
				if (value && value.length >= 2) {
					delete profileErrors.display_name;
					profileErrors = profileErrors;
				}
				break;
			case 'new_password':
				if (value && value.length >= 8 && isStrongPassword(value)) {
					delete passwordErrors.new_password;
					passwordErrors = passwordErrors;
				}
				break;
			case 'confirm_password':
				if (value && value === passwordData.new_password) {
					delete passwordErrors.confirm_password;
					passwordErrors = passwordErrors;
				}
				break;
		}
	}
</script>

<svelte:head>
	<title>Settings - Azurite</title>
	<meta
		name="description"
		content="Manage your Azurite account settings, profile, and preferences."
	/>
</svelte:head>

{#if isLoadingProfile}
	<div class="min-h-screen flex items-center justify-center">
		<Loading size="lg" text="Loading settings..." />
	</div>
{:else}
	<div class="min-h-screen bg-background-primary">
		<!-- Header -->
		<div class="bg-gradient-to-r from-slate-800/50 to-slate-700/50 border-b border-slate-700">
			<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
				<h1 class="text-3xl font-bold text-text-primary">Account Settings</h1>
				<p class="text-text-secondary mt-2">
					Manage your account preferences and security settings
				</p>
			</div>
		</div>

		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
			<div class="grid grid-cols-1 lg:grid-cols-4 gap-8">
				<!-- Sidebar Navigation -->
				<div class="lg:col-span-1">
					<nav class="space-y-1">
						{#each tabs as tab (tab.id)}
							<button
								onclick={() => (activeTab = tab.id)}
								class="w-full flex items-center px-4 py-3 text-left rounded-lg transition-colors {activeTab ===
								tab.id
									? 'bg-primary-600 text-white'
									: 'text-text-secondary hover:text-text-primary hover:bg-slate-800'}"
							>
								<svelte:component this={tab.icon} class="w-5 h-5 mr-3" />
								{tab.label}
							</button>
						{/each}
					</nav>
				</div>

				<!-- Content -->
				<div class="lg:col-span-3">
					{#if activeTab === 'profile'}
						<!-- Profile Settings -->
						<div class="card">
							<div class="p-6">
								<h2 class="text-xl font-semibold text-text-primary mb-6">Profile Information</h2>

								<form onsubmit={saveProfile} class="space-y-6">
									<!-- Avatar Section -->
									<div>
										<label class="block text-sm font-medium text-text-primary mb-3">
											Profile Picture
										</label>
										<div class="flex items-center space-x-6">
											<!-- Current/Preview Avatar -->
											<div class="relative">
												{#if avatarPreview}
													<img
														src={avatarPreview}
														alt="Avatar preview"
														class="w-20 h-20 rounded-full border border-slate-600"
													/>
													<button
														type="button"
														onclick={removeAvatarPreview}
														class="absolute -top-2 -right-2 bg-red-600 text-white rounded-full p-1 hover:bg-red-700 transition-colors"
													>
														<X class="w-4 h-4" />
													</button>
												{:else if $user?.avatar}
													<img
														src={$user.avatar}
														alt="Current avatar"
														class="w-20 h-20 rounded-full border border-slate-600"
													/>
												{:else}
													<div
														class="w-20 h-20 bg-slate-600 rounded-full flex items-center justify-center border border-slate-600"
													>
														<User class="w-10 h-10 text-text-muted" />
													</div>
												{/if}
											</div>

											<!-- Upload Controls -->
											<div class="space-y-2">
												<label class="btn btn-outline cursor-pointer">
													<Camera class="w-4 h-4 mr-2" />
													Choose Photo
													<input
														type="file"
														accept="image/*"
														onchange={handleAvatarChange}
														class="hidden"
													/>
												</label>

												{#if avatarFile}
													<button
														type="button"
														onclick={uploadAvatar}
														disabled={isUploadingAvatar}
														class="btn btn-primary btn-sm block"
													>
														{#if isUploadingAvatar}
															<Loading size="sm" inline />
														{:else}
															<Upload class="w-4 h-4 mr-1" />
															Upload
														{/if}
													</button>
												{/if}

												<p class="text-xs text-text-muted">JPG, PNG, or WebP. Max size 2MB.</p>
											</div>
										</div>
									</div>

									<!-- Display Name -->
									<div>
										<label
											for="display_name"
											class="block text-sm font-medium text-text-primary mb-2"
										>
											Display Name
											<span class="text-red-400">*</span>
										</label>
										<input
											id="display_name"
											type="text"
											bind:value={profileData.display_name}
											oninput={() => validateField('display_name', profileData.display_name)}
											class="input {profileErrors.display_name ? 'input-error' : ''}"
											placeholder="How others will see you"
											required
										/>
										{#if profileErrors.display_name}
											<p class="error-message">{profileErrors.display_name}</p>
										{/if}
									</div>

									<!-- Username Display -->
									<div>
										<label class="block text-sm font-medium text-text-primary mb-2">
											Username
										</label>
										<div class="input bg-slate-800 text-text-muted cursor-not-allowed">
											@{$user?.username}
										</div>
										<p class="text-xs text-text-muted mt-1">Your username cannot be changed</p>
									</div>

									<!-- Email Display -->
									<div>
										<label class="block text-sm font-medium text-text-primary mb-2">
											Email Address
										</label>
										<div class="input bg-slate-800 text-text-muted cursor-not-allowed">
											{$user?.email}
										</div>
										<p class="text-xs text-text-muted mt-1">
											Contact support to change your email address
										</p>
									</div>

									<!-- Bio -->
									<div>
										<label for="bio" class="block text-sm font-medium text-text-primary mb-2">
											Bio
										</label>
										<textarea
											id="bio"
											bind:value={profileData.bio}
											class="textarea {profileErrors.bio ? 'input-error' : ''}"
											placeholder="Tell others about yourself..."
											rows="4"
											maxlength="500"
										></textarea>
										<div class="flex justify-between items-center mt-1">
											{#if profileErrors.bio}
												<p class="error-message">{profileErrors.bio}</p>
											{:else}
												<p class="text-xs text-text-muted">Brief description for your profile</p>
											{/if}
											<p class="text-xs text-text-muted">
												{profileData.bio.length}/500
											</p>
										</div>
									</div>

									<!-- Save Button -->
									<div class="flex justify-end">
										<button type="submit" disabled={isSavingProfile} class="btn btn-primary">
											{#if isSavingProfile}
												<Loading size="sm" inline />
											{:else}
												<Save class="w-4 h-4 mr-2" />
												Save Changes
											{/if}
										</button>
									</div>
								</form>
							</div>
						</div>
					{:else if activeTab === 'security'}
						<!-- Security Settings -->
						<div class="card">
							<div class="p-6">
								<h2 class="text-xl font-semibold text-text-primary mb-6">Security Settings</h2>

								<form onsubmit={changePassword} class="space-y-6">
									<!-- Current Password -->
									<div>
										<label
											for="current_password"
											class="block text-sm font-medium text-text-primary mb-2"
										>
											Current Password
											<span class="text-red-400">*</span>
										</label>
										<div class="relative">
											<input
												id="current_password"
												type={showCurrentPassword ? 'text' : 'password'}
												bind:value={passwordData.current_password}
												class="input pr-10 {passwordErrors.current_password ? 'input-error' : ''}"
												placeholder="Enter your current password"
												required
											/>
											<button
												type="button"
												class="absolute inset-y-0 right-0 pr-3 flex items-center"
												onclick={() => togglePasswordVisibility('current')}
											>
												{#if showCurrentPassword}
													<EyeOff class="h-5 w-5 text-text-muted" />
												{:else}
													<Eye class="h-5 w-5 text-text-muted" />
												{/if}
											</button>
										</div>
										{#if passwordErrors.current_password}
											<p class="error-message">{passwordErrors.current_password}</p>
										{/if}
									</div>

									<!-- New Password -->
									<div>
										<label
											for="new_password"
											class="block text-sm font-medium text-text-primary mb-2"
										>
											New Password
											<span class="text-red-400">*</span>
										</label>
										<div class="relative">
											<input
												id="new_password"
												type={showNewPassword ? 'text' : 'password'}
												bind:value={passwordData.new_password}
												oninput={() => validateField('new_password', passwordData.new_password)}
												class="input pr-10 {passwordErrors.new_password ? 'input-error' : ''}"
												placeholder="Enter a new password"
												required
											/>
											<button
												type="button"
												class="absolute inset-y-0 right-0 pr-3 flex items-center"
												onclick={() => togglePasswordVisibility('new')}
											>
												{#if showNewPassword}
													<EyeOff class="h-5 w-5 text-text-muted" />
												{:else}
													<Eye class="h-5 w-5 text-text-muted" />
												{/if}
											</button>
										</div>
										{#if passwordErrors.new_password}
											<p class="error-message">{passwordErrors.new_password}</p>
										{:else}
											<p class="text-xs text-text-muted">
												Must be at least 8 characters with uppercase, lowercase, number, and special
												character
											</p>
										{/if}
									</div>

									<!-- Confirm New Password -->
									<div>
										<label
											for="confirm_password"
											class="block text-sm font-medium text-text-primary mb-2"
										>
											Confirm New Password
											<span class="text-red-400">*</span>
										</label>
										<div class="relative">
											<input
												id="confirm_password"
												type={showConfirmPassword ? 'text' : 'password'}
												bind:value={passwordData.confirm_password}
												oninput={() =>
													validateField('confirm_password', passwordData.confirm_password)}
												class="input pr-10 {passwordErrors.confirm_password ? 'input-error' : ''}"
												placeholder="Confirm your new password"
												required
											/>
											<button
												type="button"
												class="absolute inset-y-0 right-0 pr-3 flex items-center"
												onclick={() => togglePasswordVisibility('confirm')}
											>
												{#if showConfirmPassword}
													<EyeOff class="h-5 w-5 text-text-muted" />
												{:else}
													<Eye class="h-5 w-5 text-text-muted" />
												{/if}
											</button>
										</div>
										{#if passwordErrors.confirm_password}
											<p class="error-message">{passwordErrors.confirm_password}</p>
										{/if}
									</div>

									<!-- Change Password Button -->
									<div class="flex justify-end">
										<button type="submit" disabled={isChangingPassword} class="btn btn-primary">
											{#if isChangingPassword}
												<Loading size="sm" inline />
											{:else}
												<Lock class="w-4 h-4 mr-2" />
												Change Password
											{/if}
										</button>
									</div>
								</form>

								<!-- Account Actions -->
								<div class="mt-8 pt-8 border-t border-slate-700">
									<h3 class="text-lg font-medium text-text-primary mb-4">Account Actions</h3>

									<div class="bg-red-900/20 border border-red-600 rounded-lg p-4">
										<div class="flex items-center mb-2">
											<AlertTriangle class="w-5 h-5 text-red-400 mr-2" />
											<h4 class="text-red-400 font-medium">Danger Zone</h4>
										</div>
										<p class="text-text-secondary text-sm mb-4">
											This action cannot be undone. Your account and all associated data will be
											permanently deleted.
										</p>
										<button class="btn btn-danger btn-sm">
											<Trash2 class="w-4 h-4 mr-2" />
											Delete Account
										</button>
									</div>
								</div>
							</div>
						</div>
					{:else if activeTab === 'notifications'}
						<!-- Notification Settings -->
						<div class="card">
							<div class="p-6">
								<h2 class="text-xl font-semibold text-text-primary mb-6">
									Notification Preferences
								</h2>

								<form onsubmit={saveProfile} class="space-y-6">
									<div class="space-y-4">
										<!-- Email Notifications -->
										<div class="flex items-start">
											<div class="flex items-center h-5">
												<input
													id="notify_email"
													type="checkbox"
													bind:checked={profileData.notify_email}
													class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-slate-600 bg-slate-800 rounded"
												/>
											</div>
											<div class="ml-3">
												<label for="notify_email" class="text-sm font-medium text-text-primary">
													Email Notifications
												</label>
												<p class="text-sm text-text-secondary">
													Receive notifications via email about your mods, comments, and account
													activity.
												</p>
											</div>
										</div>

										<!-- In-Site Notifications -->
										<div class="flex items-start">
											<div class="flex items-center h-5">
												<input
													id="notify_in_site"
													type="checkbox"
													bind:checked={profileData.notify_in_site}
													class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-slate-600 bg-slate-800 rounded"
												/>
											</div>
											<div class="ml-3">
												<label for="notify_in_site" class="text-sm font-medium text-text-primary">
													In-Site Notifications
												</label>
												<p class="text-sm text-text-secondary">
													Show notifications in the site header when you're logged in.
												</p>
											</div>
										</div>
									</div>

									<!-- Notification Types -->
									<div class="border-t border-slate-700 pt-6">
										<h3 class="text-sm font-medium text-text-primary mb-4">
											You will receive notifications for:
										</h3>
										<div class="space-y-2 text-sm text-text-secondary">
											<div class="flex items-center">
												<Check class="w-4 h-4 text-green-400 mr-2" />
												New comments on your mods
											</div>
											<div class="flex items-center">
												<Check class="w-4 h-4 text-green-400 mr-2" />
												Mod approval or rejection status
											</div>
											<div class="flex items-center">
												<Check class="w-4 h-4 text-green-400 mr-2" />
												Download and like milestones
											</div>
											<div class="flex items-center">
												<Check class="w-4 h-4 text-green-400 mr-2" />
												Important account updates
											</div>
										</div>
									</div>

									<!-- Save Button -->
									<div class="flex justify-end pt-4">
										<button type="submit" disabled={isSavingProfile} class="btn btn-primary">
											{#if isSavingProfile}
												<Loading size="sm" inline />
											{:else}
												<Save class="w-4 h-4 mr-2" />
												Save Preferences
											{/if}
										</button>
									</div>
								</form>
							</div>
						</div>
					{/if}
				</div>
			</div>
		</div>
	</div>
{/if}
