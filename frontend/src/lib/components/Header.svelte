<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { auth, user, isAuthenticated } from '$lib/stores/auth';
	import { notificationsApi } from '$lib/api/client';
	import { Bell, Menu, Search, User, Settings, LogOut, Shield, Plus } from 'lucide-svelte';

	let showMobileMenu = false;
	let showUserMenu = false;
	let showNotifications = false;
	let searchQuery = '';
	let unreadCount = 0;

	// Toggle mobile menu
	function toggleMobileMenu() {
		showMobileMenu = !showMobileMenu;
	}

	// Toggle user menu
	function toggleUserMenu() {
		showUserMenu = !showUserMenu;
	}

	// Toggle notifications
	function toggleNotifications() {
		showNotifications = !showNotifications;
	}

	// Handle search
	function handleSearch() {
		if (searchQuery.trim()) {
			goto(`/search?q=${encodeURIComponent(searchQuery.trim())}`);
			searchQuery = '';
		}
	}

	// Handle logout
	async function handleLogout() {
		auth.logout();
		showUserMenu = false;
		goto('/');
	}

	// Load unread notifications count
	async function loadUnreadCount() {
		if ($isAuthenticated) {
			const response = await notificationsApi.getUnreadCount();
			if (response.success) {
				unreadCount = response.data?.count || 0;
			}
		}
	}

	// Close menus when clicking outside
	function handleClickOutside(event: MouseEvent) {
		const target = event.target as Element;

		// Close user menu
		if (showUserMenu && !target.closest('.user-menu-container')) {
			showUserMenu = false;
		}

		// Close notifications
		if (showNotifications && !target.closest('.notifications-container')) {
			showNotifications = false;
		}

		// Close mobile menu
		if (
			showMobileMenu &&
			!target.closest('.mobile-menu-container') &&
			!target.closest('.mobile-menu-button')
		) {
			showMobileMenu = false;
		}
	}

	onMount(() => {
		loadUnreadCount();
		document.addEventListener('click', handleClickOutside);

		return () => {
			document.removeEventListener('click', handleClickOutside);
		};
	});

	// Reactive statement to update unread count when auth changes
	$: if ($isAuthenticated) {
		loadUnreadCount();
	}

	// Navigation links
	const navLinks = [
		{ href: '/', label: 'Home' },
		{ href: '/games', label: 'Games' },
		{ href: '/browse', label: 'Browse Mods' }
	];

	// Check if current route is active
	function isActive(href: string): boolean {
		if (href === '/') {
			return $page.url.pathname === '/';
		}
		return $page.url.pathname.startsWith(href);
	}
</script>

<svelte:window
	on:keydown={(e) =>
		e.key === 'Escape' && (showMobileMenu = showUserMenu = showNotifications = false)}
/>

<header
	class="sticky top-0 z-40 w-full bg-background-primary/95 backdrop-blur-sm border-b border-slate-700"
>
	<nav class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
		<div class="flex justify-between items-center h-16">
			<!-- Logo and Brand -->
			<div class="flex items-center">
				<a href="/" class="flex items-center space-x-2 group">
					<img
						src="/logo-nobg.png"
						alt="Azurite Logo"
						class="w-8 h-8 rounded-lg group-hover:shadow-glow-blue-lg transition-all duration-200"
					/>
					<span class="text-xl font-bold text-gradient">Azurite</span>
				</a>
			</div>

			<!-- Desktop Navigation -->
			<div class="hidden md:flex items-center space-x-8">
				{#each navLinks as link (link.href)}
					<a href={link.href} class="nav-link {isActive(link.href) ? 'nav-link-active' : ''}">
						{link.label}
					</a>
				{/each}
			</div>

			<!-- Search Bar (Desktop) -->
			<div class="hidden md:flex flex-1 max-w-lg mx-8">
				<form on:submit|preventDefault={handleSearch} class="w-full">
					<div class="relative">
						<Search
							class="absolute left-3 top-1/2 transform -translate-y-1/2 text-text-muted w-4 h-4"
						/>
						<input
							type="search"
							placeholder="Search mods, games..."
							bind:value={searchQuery}
							class="input pl-10 pr-4 py-2"
						/>
					</div>
				</form>
			</div>

			<!-- Right Side Actions -->
			<div class="flex items-center space-x-4">
				<!-- Mobile Search Toggle -->
				<button class="md:hidden p-2 text-text-secondary hover:text-text-primary transition-colors">
					<Search class="w-5 h-5" />
				</button>

				{#if $isAuthenticated}
					<!-- Notifications -->
					<div class="relative notifications-container">
						<button
							on:click={toggleNotifications}
							class="p-2 text-text-secondary hover:text-text-primary transition-colors relative"
							title="Notifications"
						>
							<Bell class="w-5 h-5" />
							{#if unreadCount > 0}
								<span
									class="absolute -top-1 -right-1 bg-primary-600 text-white text-xs rounded-full min-w-[1.25rem] h-5 flex items-center justify-center"
								>
									{unreadCount > 99 ? '99+' : unreadCount}
								</span>
							{/if}
						</button>

						<!-- Notifications Dropdown -->
						{#if showNotifications}
							<div class="dropdown right-0 w-80 max-h-96 overflow-y-auto">
								<div class="p-3 border-b border-slate-600">
									<h3 class="font-medium text-text-primary">Notifications</h3>
								</div>

								<div class="p-2">
									<!-- Notification items would go here -->
									<div class="text-center py-8 text-text-muted">
										<Bell class="w-8 h-8 mx-auto mb-2 opacity-50" />
										<p>No new notifications</p>
									</div>
								</div>

								<div class="p-2 border-t border-slate-600">
									<a href="/notifications" class="dropdown-item text-center text-primary-400">
										View all notifications
									</a>
								</div>
							</div>
						{/if}
					</div>

					<!-- User Menu -->
					<div class="relative user-menu-container">
						<button
							on:click={toggleUserMenu}
							class="flex items-center space-x-2 p-2 rounded-lg hover:bg-slate-700 transition-colors"
						>
							{#if $user?.avatar}
								<img
									src={$user.avatar}
									alt={$user.display_name}
									class="w-8 h-8 rounded-full border border-slate-600"
								/>
							{:else}
								<div class="w-8 h-8 bg-slate-600 rounded-full flex items-center justify-center">
									<User class="w-4 h-4" />
								</div>
							{/if}
							<span class="hidden sm:block text-sm font-medium text-text-primary">
								{$user?.display_name || $user?.username}
							</span>
						</button>

						<!-- User Dropdown -->
						{#if showUserMenu}
							<div class="dropdown right-0 w-48">
								<a href="/profile" class="dropdown-item">
									<User class="w-4 h-4 mr-2" />
									Profile
								</a>
								<a href="/dashboard" class="dropdown-item">
									<Plus class="w-4 h-4 mr-2" />
									Creator Dashboard
								</a>
								<a href="/settings" class="dropdown-item">
									<Settings class="w-4 h-4 mr-2" />
									Settings
								</a>

								{#if $user?.role === 'admin' || $user?.role === 'community_moderator'}
									<div class="border-t border-slate-600 my-1"></div>
									<a href="/admin" class="dropdown-item">
										<Shield class="w-4 h-4 mr-2" />
										Admin Panel
									</a>
								{/if}

								<div class="border-t border-slate-600 my-1"></div>
								<button
									on:click={handleLogout}
									class="dropdown-item text-red-400 hover:bg-red-500/10"
								>
									<LogOut class="w-4 h-4 mr-2" />
									Logout
								</button>
							</div>
						{/if}
					</div>
				{:else}
					<!-- Auth Actions -->
					<div class="flex items-center space-x-2">
						<a href="/auth/login" class="btn btn-outline btn-sm"> Login </a>
						<a href="/auth/register" class="btn btn-primary btn-sm"> Sign Up </a>
					</div>
				{/if}

				<!-- Mobile Menu Button -->
				<button
					on:click={toggleMobileMenu}
					class="md:hidden p-2 text-text-secondary hover:text-text-primary transition-colors mobile-menu-button"
				>
					<Menu class="w-5 h-5" />
				</button>
			</div>
		</div>

		<!-- Mobile Menu -->
		{#if showMobileMenu}
			<div class="md:hidden mobile-menu-container">
				<div class="px-2 pt-2 pb-3 space-y-1 bg-background-secondary border-t border-slate-700">
					<!-- Mobile Search -->
					<form on:submit|preventDefault={handleSearch} class="mb-4">
						<div class="relative">
							<Search
								class="absolute left-3 top-1/2 transform -translate-y-1/2 text-text-muted w-4 h-4"
							/>
							<input
								type="search"
								placeholder="Search mods, games..."
								bind:value={searchQuery}
								class="input pl-10 pr-4 py-2 w-full"
							/>
						</div>
					</form>

					<!-- Mobile Navigation Links -->
					{#each navLinks as link (link.label)}
						<a
							href={link.href}
							class="block px-3 py-2 text-base font-medium {isActive(link.href)
								? 'text-primary-400 bg-slate-700'
								: 'text-text-secondary hover:text-text-primary hover:bg-slate-700'} rounded-md transition-colors"
							on:click={() => (showMobileMenu = false)}
						>
							{link.label}
						</a>
					{/each}

					{#if $isAuthenticated}
						<div class="border-t border-slate-600 mt-4 pt-4">
							<a
								href="/profile"
								class="block px-3 py-2 text-base font-medium text-text-secondary hover:text-text-primary hover:bg-slate-700 rounded-md transition-colors"
								on:click={() => (showMobileMenu = false)}
							>
								<User class="w-4 h-4 inline mr-2" />
								Profile
							</a>
							<a
								href="/dashboard"
								class="block px-3 py-2 text-base font-medium text-text-secondary hover:text-text-primary hover:bg-slate-700 rounded-md transition-colors"
								on:click={() => (showMobileMenu = false)}
							>
								<Plus class="w-4 h-4 inline mr-2" />
								Creator Dashboard
							</a>
							<a
								href="/settings"
								class="block px-3 py-2 text-base font-medium text-text-secondary hover:text-text-primary hover:bg-slate-700 rounded-md transition-colors"
								on:click={() => (showMobileMenu = false)}
							>
								<Settings class="w-4 h-4 inline mr-2" />
								Settings
							</a>

							{#if $user?.role === 'admin' || $user?.role === 'community_moderator'}
								<a
									href="/admin"
									class="block px-3 py-2 text-base font-medium text-text-secondary hover:text-text-primary hover:bg-slate-700 rounded-md transition-colors"
									on:click={() => (showMobileMenu = false)}
								>
									<Shield class="w-4 h-4 inline mr-2" />
									Admin Panel
								</a>
							{/if}

							<button
								on:click={() => {
									handleLogout();
									showMobileMenu = false;
								}}
								class="block w-full text-left px-3 py-2 text-base font-medium text-red-400 hover:text-red-300 hover:bg-red-500/10 rounded-md transition-colors"
							>
								<LogOut class="w-4 h-4 inline mr-2" />
								Logout
							</button>
						</div>
					{:else}
						<div class="border-t border-slate-600 mt-4 pt-4 flex space-x-2">
							<a
								href="/auth/login"
								class="btn btn-outline flex-1 text-center"
								on:click={() => (showMobileMenu = false)}
							>
								Login
							</a>
							<a
								href="/auth/register"
								class="btn btn-primary flex-1 text-center"
								on:click={() => (showMobileMenu = false)}
							>
								Sign Up
							</a>
						</div>
					{/if}
				</div>
			</div>
		{/if}
	</nav>
</header>

<style>
	.user-menu-container,
	.notifications-container,
	.mobile-menu-container {
		position: relative;
	}
</style>
