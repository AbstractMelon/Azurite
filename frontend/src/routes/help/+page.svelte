<script lang="ts">
	import { onMount } from 'svelte';
	import {
		Search,
		Book,
		MessageCircle,
		Bug,
		Lightbulb,
		Mail,
		ExternalLink,
		ChevronRight
	} from 'lucide-svelte';

	let searchQuery = '';
	let selectedCategory = 'all';

	const categories = [
		{ value: 'all', label: 'All Topics' },
		{ value: 'getting-started', label: 'Getting Started' },
		{ value: 'mod-management', label: 'Mod Management' },
		{ value: 'account', label: 'Account & Profile' },
		{ value: 'community', label: 'Community' },
		{ value: 'troubleshooting', label: 'Troubleshooting' },
		{ value: 'policies', label: 'Policies' }
	];

	const helpTopics = [
		{
			id: 1,
			title: 'Getting Started with Azurite',
			description: 'Learn the basics of using Azurite to discover and manage mods',
			category: 'getting-started',
			popular: true,
			content: `
				<h3>Welcome to Azurite!</h3>
				<p>Azurite is a community-driven platform for discovering, downloading, and sharing game modifications (mods). Here's how to get started:</p>

				<h4>1. Create Your Account</h4>
				<ul>
					<li>Click "Sign Up" in the top right corner</li>
					<li>Fill in your details or use social login (GitHub, Google, Discord)</li>
					<li>Verify your email address</li>
					<li>Customize your profile</li>
				</ul>

				<h4>2. Discover Mods</h4>
				<ul>
					<li>Browse games from the main navigation</li>
					<li>Use the search function to find specific mods</li>
					<li>Filter by categories, popularity, and update date</li>
					<li>Read descriptions and user reviews</li>
				</ul>

				<h4>3. Download and Install</h4>
				<ul>
					<li>Click on any mod to view details</li>
					<li>Check compatibility with your game version</li>
					<li>Read installation instructions carefully</li>
					<li>Download the mod files</li>
				</ul>
			`
		},
		{
			id: 2,
			title: 'How to Upload Your First Mod',
			description: 'Step-by-step guide to sharing your creations with the community',
			category: 'mod-management',
			popular: true,
			content: `
				<h3>Uploading Your Mod</h3>
				<p>Ready to share your creation? Follow these steps to upload your mod:</p>

				<h4>Prerequisites</h4>
				<ul>
					<li>Have a verified account</li>
					<li>Your mod files ready and tested</li>
					<li>Screenshots or images of your mod</li>
					<li>Clear description and installation instructions</li>
				</ul>

				<h4>Upload Process</h4>
				<ol>
					<li>Go to your Dashboard and click "Upload Mod"</li>
					<li>Select the game your mod is for</li>
					<li>Fill in the mod details (name, description, tags)</li>
					<li>Upload your mod files (ZIP format recommended)</li>
					<li>Add screenshots and preview images</li>
					<li>Set the appropriate tags and categories</li>
					<li>Review and submit for approval</li>
				</ol>

				<h4>Approval Process</h4>
				<p>All mods go through a review process that typically takes 24-48 hours. We check for:</p>
				<ul>
					<li>Malware and security threats</li>
					<li>Content policy compliance</li>
					<li>Basic functionality</li>
					<li>Appropriate categorization</li>
				</ul>
			`
		},
		{
			id: 3,
			title: 'Managing Your Account Settings',
			description: 'Customize your profile, notifications, and privacy settings',
			category: 'account',
			popular: false,
			content: `
				<h3>Account Settings</h3>
				<p>Personalize your Azurite experience through your account settings:</p>

				<h4>Profile Settings</h4>
				<ul>
					<li>Update your display name and bio</li>
					<li>Upload a profile picture</li>
					<li>Link your social media accounts</li>
					<li>Set your timezone and language preferences</li>
				</ul>

				<h4>Notification Settings</h4>
				<ul>
					<li>Choose email notification preferences</li>
					<li>Set in-app notification types</li>
					<li>Configure mod update alerts</li>
					<li>Manage comment and mention notifications</li>
				</ul>

				<h4>Privacy Settings</h4>
				<ul>
					<li>Control who can see your profile</li>
					<li>Manage your mod visibility</li>
					<li>Set download history privacy</li>
					<li>Configure data sharing preferences</li>
				</ul>
			`
		},
		{
			id: 4,
			title: 'Community Guidelines and Best Practices',
			description: 'Learn how to be a positive member of the Azurite community',
			category: 'community',
			popular: true,
			content: `
				<h3>Community Guidelines</h3>
				<p>Help us maintain a welcoming and supportive community:</p>

				<h4>Respectful Interaction</h4>
				<ul>
					<li>Be kind and respectful to all users</li>
					<li>Provide constructive feedback on mods</li>
					<li>Help new users learn the platform</li>
					<li>Report inappropriate behavior</li>
				</ul>

				<h4>Quality Content</h4>
				<ul>
					<li>Upload original, working mods</li>
					<li>Provide clear descriptions and instructions</li>
					<li>Use appropriate tags and categories</li>
					<li>Keep mods updated and maintained</li>
				</ul>

				<h4>Prohibited Content</h4>
				<ul>
					<li>Malicious or harmful software</li>
					<li>Copyrighted material without permission</li>
					<li>Inappropriate or offensive content</li>
					<li>Spam or duplicate uploads</li>
				</ul>
			`
		},
		{
			id: 5,
			title: 'Troubleshooting Download Issues',
			description: 'Common problems and solutions for mod downloads',
			category: 'troubleshooting',
			popular: false,
			content: `
				<h3>Download Troubleshooting</h3>
				<p>Having trouble downloading mods? Try these solutions:</p>

				<h4>Common Issues</h4>
				<ul>
					<li><strong>Download fails or stops:</strong> Check your internet connection and try again</li>
					<li><strong>File appears corrupted:</strong> Clear browser cache and re-download</li>
					<li><strong>Download button not working:</strong> Disable ad blockers temporarily</li>
					<li><strong>Slow download speeds:</strong> Try downloading during off-peak hours</li>
				</ul>

				<h4>Browser-Specific Fixes</h4>
				<ul>
					<li><strong>Chrome:</strong> Check if downloads are blocked in settings</li>
					<li><strong>Firefox:</strong> Verify file type permissions</li>
					<li><strong>Safari:</strong> Allow downloads from unknown developers</li>
					<li><strong>Edge:</strong> Check SmartScreen filter settings</li>
				</ul>

				<h4>Still Having Issues?</h4>
				<p>If problems persist, try:</p>
				<ul>
					<li>Using a different browser</li>
					<li>Disabling VPN or proxy</li>
					<li>Contacting the mod author</li>
					<li>Reporting the issue to our support team</li>
				</ul>
			`
		},
		{
			id: 6,
			title: 'Understanding Our Policies',
			description: 'Learn about our content, privacy, and community policies',
			category: 'policies',
			popular: false,
			content: `
				<h3>Azurite Policies</h3>
				<p>Understanding our policies helps ensure a safe and legal platform:</p>

				<h4>Content Policy</h4>
				<ul>
					<li>All content must comply with applicable laws</li>
					<li>No malicious software or harmful code</li>
					<li>Respect intellectual property rights</li>
					<li>Age-appropriate content only</li>
				</ul>

				<h4>Privacy Policy</h4>
				<ul>
					<li>We protect your personal information</li>
					<li>Data is used only for platform functionality</li>
					<li>You control your privacy settings</li>
					<li>We don't sell your data to third parties</li>
				</ul>

				<h4>DMCA Policy</h4>
				<ul>
					<li>We respond to valid DMCA takedown notices</li>
					<li>Copyright holders can report infringement</li>
					<li>Users can file counter-notifications</li>
					<li>Repeat offenders may be banned</li>
				</ul>
			`
		}
	];

	let filteredTopics = helpTopics;
	let selectedTopic = null;

	function filterTopics() {
		filteredTopics = helpTopics.filter((topic) => {
			const matchesCategory = selectedCategory === 'all' || topic.category === selectedCategory;
			const matchesSearch =
				!searchQuery ||
				topic.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
				topic.description.toLowerCase().includes(searchQuery.toLowerCase());
			return matchesCategory && matchesSearch;
		});
	}

	function selectTopic(topic) {
		selectedTopic = topic;
	}

	function clearTopic() {
		selectedTopic = null;
	}

	$: if (searchQuery || selectedCategory) {
		filterTopics();
	}
</script>

<svelte:head>
	<title>Help Center - Azurite</title>
	<meta
		name="description"
		content="Get help with Azurite mod hosting platform - guides, tutorials, and support"
	/>
</svelte:head>

<div class="min-h-screen bg-background-primary">
	<div class="container mx-auto px-4 py-8">
		{#if selectedTopic}
			<!-- Topic Detail View -->
			<div class="mb-6">
				<button onclick={clearTopic} class="btn btn-outline flex items-center gap-2 mb-4">
					← Back to Help Center
				</button>
			</div>

			<div class="card">
				<div class="p-8">
					<h1 class="text-3xl font-bold text-text-primary mb-4">{selectedTopic.title}</h1>
					<p class="text-text-secondary mb-6">{selectedTopic.description}</p>
					<div class="prose prose-invert max-w-none">
						{@html selectedTopic.content}
					</div>
				</div>
			</div>
		{:else}
			<!-- Help Center Main View -->
			<div class="text-center mb-8">
				<h1 class="text-4xl font-bold text-text-primary mb-4">How can we help you?</h1>
				<p class="text-text-secondary text-lg">
					Find answers, guides, and support for using Azurite
				</p>
			</div>

			<!-- Search and Filter -->
			<div class="card mb-8">
				<div class="p-6">
					<div class="flex flex-col sm:flex-row gap-4 mb-6">
						<div class="flex-1 relative">
							<div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
								<Search class="h-5 w-5 text-text-muted" />
							</div>
							<input
								type="text"
								bind:value={searchQuery}
								placeholder="Search help topics..."
								class="input pl-10 w-full"
							/>
						</div>
						<select bind:value={selectedCategory} class="select w-full sm:w-auto">
							{#each categories as category}
								<option value={category.value}>{category.label}</option>
							{/each}
						</select>
					</div>

					<!-- Popular Topics -->
					<div class="mb-4">
						<h3 class="text-lg font-semibold text-text-primary mb-3">Popular Topics</h3>
						<div class="flex flex-wrap gap-2">
							{#each helpTopics.filter((t) => t.popular) as topic}
								<button onclick={() => selectTopic(topic)} class="btn btn-outline btn-sm">
									{topic.title}
								</button>
							{/each}
						</div>
					</div>
				</div>
			</div>

			<!-- Quick Actions -->
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
				<a href="/bug-report" class="card hover:shadow-lg transition-shadow">
					<div class="p-6 text-center">
						<Bug class="h-8 w-8 text-red-400 mx-auto mb-3" />
						<h3 class="font-semibold text-text-primary mb-2">Report a Bug</h3>
						<p class="text-sm text-text-secondary">Found an issue? Let us know so we can fix it.</p>
					</div>
				</a>

				<a href="/feature-request" class="card hover:shadow-lg transition-shadow">
					<div class="p-6 text-center">
						<Lightbulb class="h-8 w-8 text-yellow-400 mx-auto mb-3" />
						<h3 class="font-semibold text-text-primary mb-2">Request a Feature</h3>
						<p class="text-sm text-text-secondary">Have an idea? Share it with our team.</p>
					</div>
				</a>

				<a href="/contact" class="card hover:shadow-lg transition-shadow">
					<div class="p-6 text-center">
						<Mail class="h-8 w-8 text-blue-400 mx-auto mb-3" />
						<h3 class="font-semibold text-text-primary mb-2">Contact Support</h3>
						<p class="text-sm text-text-secondary">Need personal assistance? Get in touch.</p>
					</div>
				</a>

				<a href="/guidelines" class="card hover:shadow-lg transition-shadow">
					<div class="p-6 text-center">
						<Book class="h-8 w-8 text-green-400 mx-auto mb-3" />
						<h3 class="font-semibold text-text-primary mb-2">Guidelines</h3>
						<p class="text-sm text-text-secondary">Learn about our community standards.</p>
					</div>
				</a>
			</div>

			<!-- Help Topics -->
			<div class="card">
				<div class="p-6">
					<h2 class="text-2xl font-bold text-text-primary mb-6">Help Topics</h2>

					{#if filteredTopics.length > 0}
						<div class="space-y-4">
							{#each filteredTopics as topic}
								<button
									onclick={() => selectTopic(topic)}
									class="w-full text-left p-4 rounded-lg border border-slate-600 hover:border-primary-500 hover:bg-background-secondary transition-colors"
								>
									<div class="flex items-center justify-between">
										<div class="flex-1">
											<h3 class="font-semibold text-text-primary mb-1">
												{topic.title}
												{#if topic.popular}
													<span
														class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-primary-100 text-primary-800 ml-2"
													>
														Popular
													</span>
												{/if}
											</h3>
											<p class="text-text-secondary text-sm">{topic.description}</p>
										</div>
										<ChevronRight class="h-5 w-5 text-text-muted ml-4" />
									</div>
								</button>
							{/each}
						</div>
					{:else}
						<div class="text-center py-8">
							<div class="text-4xl text-text-muted mb-4">🔍</div>
							<h3 class="text-lg font-semibold text-text-primary mb-2">No topics found</h3>
							<p class="text-text-secondary">Try adjusting your search or category filter.</p>
						</div>
					{/if}
				</div>
			</div>

			<!-- Additional Resources -->
			<div class="mt-8">
				<h2 class="text-2xl font-bold text-text-primary mb-6 text-center">Additional Resources</h2>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
					<div class="card">
						<div class="p-6">
							<MessageCircle class="h-6 w-6 text-primary-400 mb-3" />
							<h3 class="font-semibold text-text-primary mb-2">Community Forum</h3>
							<p class="text-text-secondary mb-4">
								Connect with other users, ask questions, and share knowledge.
							</p>
							<a href="#" class="text-primary-400 hover:text-primary-300 flex items-center gap-1">
								Visit Forum <ExternalLink class="h-4 w-4" />
							</a>
						</div>
					</div>

					<div class="card">
						<div class="p-6">
							<Book class="h-6 w-6 text-primary-400 mb-3" />
							<h3 class="font-semibold text-text-primary mb-2">API Documentation</h3>
							<p class="text-text-secondary mb-4">
								Technical documentation for developers and integrators.
							</p>
							<a
								href="/api/docs"
								class="text-primary-400 hover:text-primary-300 flex items-center gap-1"
							>
								View Docs <ExternalLink class="h-4 w-4" />
							</a>
						</div>
					</div>
				</div>
			</div>

			<!-- Still Need Help -->
			<div class="mt-8 text-center">
				<div class="card">
					<div class="p-8">
						<h3 class="text-xl font-semibold text-text-primary mb-4">Still need help?</h3>
						<p class="text-text-secondary mb-6">
							Can't find what you're looking for? Our support team is here to help.
						</p>
						<a href="/contact" class="btn btn-primary"> Contact Support </a>
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>

<style>
	/* Custom prose styles for help content */
	.prose h3 {
		color: rgb(226 232 240);
		font-size: 1.25rem;
		font-weight: 600;
		margin-top: 1.5rem;
		margin-bottom: 0.75rem;
	}

	.prose h4 {
		color: rgb(203 213 225);
		font-size: 1.125rem;
		font-weight: 500;
		margin-top: 1.25rem;
		margin-bottom: 0.5rem;
	}

	.prose p {
		color: rgb(148 163 184);
		margin-bottom: 1rem;
		line-height: 1.6;
	}

	.prose ul,
	.prose ol {
		color: rgb(148 163 184);
		margin-bottom: 1rem;
		padding-left: 1.5rem;
	}

	.prose li {
		margin-bottom: 0.5rem;
	}

	.prose strong {
		color: rgb(226 232 240);
	}
</style>
