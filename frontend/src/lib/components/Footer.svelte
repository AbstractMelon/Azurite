<script lang="ts">
	import { Github, Twitter, Heart, Mail } from 'lucide-svelte';

	const currentYear = new Date().getFullYear();

	const footerLinks = {
		platform: [
			{ label: 'Browse Games', href: '/games' },
			{ label: 'Popular Mods', href: '/browse?sort=popular' },
			{ label: 'New Releases', href: '/browse?sort=newest' },
			{ label: 'Categories', href: '/categories' }
		],
		community: [
			{ label: 'Discord Server', href: 'https://discord.gg/azurite' },
			{ label: 'Forums', href: '/forums' },
			{ label: 'Documentation', href: '/docs' },
			{ label: 'Developer API', href: '/api-docs' }
		],
		support: [
			{ label: 'Help Center', href: '/help' },
			{ label: 'Contact Us', href: '/contact' },
			{ label: 'Bug Reports', href: '/bug-report' },
			{ label: 'Feature Requests', href: '/feature-request' }
		],
		legal: [
			{ label: 'Terms of Service', href: '/terms' },
			{ label: 'Privacy Policy', href: '/privacy' },
			{ label: 'DMCA Policy', href: '/dmca' },
			{ label: 'Content Guidelines', href: '/guidelines' }
		]
	};

	const socialLinks = [
		{
			icon: Github,
			href: 'https://github.com/azurite-platform',
			label: 'GitHub',
			color: 'hover:text-gray-300'
		},
		{
			icon: Twitter,
			href: 'https://twitter.com/azurite_mods',
			label: 'Twitter',
			color: 'hover:text-blue-400'
		},
		{
			icon: Mail,
			href: 'mailto:hello@azurite.dev',
			label: 'Email',
			color: 'hover:text-primary-400'
		}
	];
</script>

<!-- Newsletter Signup Section -->
<div class="bg-slate-800/50 border-t border-slate-700">
	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
		<div class="text-center">
			<h3 class="text-lg font-semibold text-text-primary mb-2">Stay updated with Azurite</h3>
			<p class="text-text-muted mb-4 max-w-md mx-auto">
				Get notified about new features, popular mods, and platform updates.
			</p>
			<form class="flex flex-col sm:flex-row gap-2 max-w-md mx-auto">
				<input type="email" placeholder="Enter your email" class="input flex-1" required />
				<button type="submit" class="btn btn-primary whitespace-nowrap"> Subscribe </button>
			</form>
			<p class="text-xs text-text-muted mt-2">We respect your privacy. Unsubscribe at any time.</p>
		</div>
	</div>
</div>

<footer class="bg-background-secondary border-t border-slate-700 mt-auto">
	<!-- Main Footer Content -->
	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-8">
			<!-- Brand Section -->
			<div class="lg:col-span-1">
				<div class="flex items-center space-x-2 mb-4">
					<div
						class="w-8 h-8 bg-gradient-to-br from-primary-400 to-primary-600 rounded-lg flex items-center justify-center shadow-glow-blue"
					>
						<span class="text-white font-bold text-lg">A</span>
					</div>
					<span class="text-xl font-bold text-gradient">Azurite</span>
				</div>
				<p class="text-text-muted text-sm mb-6 max-w-xs">
					The ultimate platform for discovering, sharing, and managing game modifications. Join
					thousands of modders and gamers in our community.
				</p>

				<!-- Social Links -->
				<div class="flex space-x-4">
					{#each socialLinks as social (social.label)}
						<a
							href={social.href}
							target="_blank"
							rel="noopener noreferrer"
							class="text-text-muted {social.color} transition-colors p-2 hover:bg-slate-700 rounded-lg"
							title={social.label}
						>
							<svelte:component this={social.icon} class="w-5 h-5" />
						</a>
					{/each}
				</div>
			</div>

			<!-- Platform Links -->
			<div>
				<h3 class="text-text-primary font-semibold mb-4">Platform</h3>
				<ul class="space-y-2">
					{#each footerLinks.platform as link (link.label)}
						<li>
							<a
								href={link.href}
								class="text-text-muted hover:text-text-primary transition-colors text-sm"
							>
								{link.label}
							</a>
						</li>
					{/each}
				</ul>
			</div>

			<!-- Community Links -->
			<div>
				<h3 class="text-text-primary font-semibold mb-4">Community</h3>
				<ul class="space-y-2">
					{#each footerLinks.community as link (link.label)}
						<li>
							<a
								href={link.href}
								class="text-text-muted hover:text-text-primary transition-colors text-sm"
								target={link.href.startsWith('http') ? '_blank' : undefined}
								rel={link.href.startsWith('http') ? 'noopener noreferrer' : undefined}
							>
								{link.label}
							</a>
						</li>
					{/each}
				</ul>
			</div>

			<!-- Support Links -->
			<div>
				<h3 class="text-text-primary font-semibold mb-4">Support</h3>
				<ul class="space-y-2">
					{#each footerLinks.support as link (link.label)}
						<li>
							<a
								href={link.href}
								class="text-text-muted hover:text-text-primary transition-colors text-sm"
							>
								{link.label}
							</a>
						</li>
					{/each}
				</ul>
			</div>

			<!-- Legal Links -->
			<div>
				<h3 class="text-text-primary font-semibold mb-4">Legal</h3>
				<ul class="space-y-2">
					{#each footerLinks.legal as link (link.label)}
						<li>
							<a
								href={link.href}
								class="text-text-muted hover:text-text-primary transition-colors text-sm"
							>
								{link.label}
							</a>
						</li>
					{/each}
				</ul>
			</div>
		</div>
	</div>

	<!-- Bottom Bar -->
	<div class="border-t border-slate-700">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
			<div class="flex flex-col md:flex-row justify-between items-center space-y-4 md:space-y-0">
				<div class="flex items-center text-sm text-text-muted">
					<span>© {currentYear} Azurite Platform. All rights reserved.</span>
				</div>

				<div class="flex items-center space-x-1 text-sm text-text-muted">
					<span>Made with</span>
					<Heart class="w-4 h-4 text-red-500 fill-current" />
					<span>for the modding community</span>
				</div>

				<div class="flex items-center space-x-6 text-sm">
					<a href="/status" class="text-text-muted hover:text-text-primary transition-colors">
						Status
					</a>
					<a href="/changelog" class="text-text-muted hover:text-text-primary transition-colors">
						Changelog
					</a>
					<span class="text-text-muted">v1.0.0</span>
				</div>
			</div>
		</div>
	</div>
</footer>
