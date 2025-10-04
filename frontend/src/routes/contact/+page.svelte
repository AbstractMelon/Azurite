<script lang="ts">
	import { onMount } from 'svelte';
	import { isAuthenticated, user } from '$lib/stores/auth';
	import { toast } from '$lib/stores/notifications';
	import { Mail, Send, Clock, MessageSquare } from 'lucide-svelte';

	let name = '';
	let email = '';
	let subject = '';
	let message = '';
	let category = 'general';
	let loading = false;
	let submitted = false;

	const categories = [
		{ value: 'general', label: 'General Inquiry' },
		{ value: 'technical', label: 'Technical Support' },
		{ value: 'account', label: 'Account Issues' },
		{ value: 'mod-review', label: 'Mod Review' },
		{ value: 'dmca', label: 'DMCA / Copyright' },
		{ value: 'partnership', label: 'Partnership / Business' },
		{ value: 'press', label: 'Press / Media' },
		{ value: 'other', label: 'Other' }
	];

	onMount(() => {
		// Pre-fill user information if logged in
		if ($isAuthenticated && $user) {
			name = $user.display_name || $user.username;
			email = $user.email || '';
		}
	});

	async function handleSubmit(event: Event) {
		event.preventDefault();

		if (!name.trim() || !email.trim() || !subject.trim() || !message.trim()) {
			toast.error('Missing information', 'Please fill in all required fields');
			return;
		}

		if (!isValidEmail(email)) {
			toast.error('Invalid email', 'Please enter a valid email address');
			return;
		}

		loading = true;

		try {
			// Since we don't have a dedicated contact API endpoint,
			// we'll simulate submitting it (in a real app, this would go to a support system)
			await new Promise(resolve => setTimeout(resolve, 1000)); // Simulate API call

			const contactForm = {
				name: name.trim(),
				email: email.trim(),
				subject: subject.trim(),
				message: message.trim(),
				category,
				submittedAt: new Date().toISOString(),
				userAgent: navigator.userAgent,
				userId: $user?.id || null
			};

			console.log('Contact form submitted:', contactForm);

			submitted = true;
			toast.success('Message sent', 'Thank you for contacting us! We\'ll get back to you soon.');
		} catch (error) {
			console.error('Failed to submit contact form:', error);
			toast.error('Submission failed', 'Please try again or contact us directly via email');
		} finally {
			loading = false;
		}
	}

	function resetForm() {
		name = $user?.display_name || $user?.username || '';
		email = $user?.email || '';
		subject = '';
		message = '';
		category = 'general';
		submitted = false;
	}

	function isValidEmail(email: string): boolean {
		const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
		return emailRegex.test(email);
	}

	function getCategoryDescription(categoryValue: string): string {
		const descriptions = {
			'general': 'Questions, feedback, or other general inquiries',
			'technical': 'Bug reports, technical issues, or platform problems',
			'account': 'Login issues, account recovery, or profile problems',
			'mod-review': 'Questions about mod approval process or review status',
			'dmca': 'Copyright infringement reports or DMCA takedown requests',
			'partnership': 'Business partnerships, sponsorships, or collaborations',
			'press': 'Media inquiries, interviews, or press releases',
			'other': 'Anything else not covered by the above categories'
		};
		return descriptions[categoryValue] || '';
	}
</script>

<svelte:head>
	<title>Contact Us - Azurite</title>
	<meta name="description" content="Get in touch with the Azurite team for support, feedback, or business inquiries" />
</svelte:head>

<div class="min-h-screen bg-background-primary">
	<div class="container mx-auto px-4 py-8 max-w-6xl">
		{#if submitted}
			<!-- Success Message -->
			<div class="max-w-2xl mx-auto">
				<div class="card">
					<div class="p-8 text-center">
						<div class="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
							<svg class="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
							</svg>
						</div>
						<h1 class="text-3xl font-bold text-text-primary mb-4">Message Sent Successfully</h1>
						<p class="text-text-secondary mb-6">
							Thank you for reaching out! We've received your message and will get back to you
							as soon as possible, typically within 24-48 hours.
						</p>
						<div class="bg-background-secondary p-4 rounded-lg border border-slate-600 mb-6">
							<p class="text-text-primary font-medium mb-2">What happens next?</p>
							<ul class="text-sm text-text-secondary space-y-1 text-left">
								<li>• Our team will review your message</li>
								<li>• You'll receive a confirmation email shortly</li>
								<li>• We'll respond to your inquiry within 1-2 business days</li>
								<li>• For urgent issues, we may contact you directly</li>
							</ul>
						</div>
						<div class="flex gap-4 justify-center">
							<button onclick={resetForm} class="btn btn-outline">
								Send Another Message
							</button>
							<a href="/" class="btn btn-primary">
								Back to Home
							</a>
						</div>
					</div>
				</div>
			</div>
		{:else}
			<!-- Contact Page Content -->
			<div class="text-center mb-8">
				<h1 class="text-4xl font-bold text-text-primary mb-4">Get in Touch</h1>
				<p class="text-text-secondary text-lg max-w-2xl mx-auto">
					Have a question, suggestion, or need support? We'd love to hear from you.
					Our team is here to help make your Azurite experience better.
				</p>
			</div>

			<div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
				<!-- Contact Information -->
				<div class="lg:col-span-1 space-y-6">
					<!-- Contact Methods -->
					<div class="card">
						<div class="p-6">
							<h2 class="text-xl font-bold text-text-primary mb-4">Contact Information</h2>
							<div class="space-y-4">
								<div class="flex items-start gap-3">
									<Mail class="h-5 w-5 text-primary-400 mt-0.5" />
									<div>
										<p class="font-medium text-text-primary">Email Support</p>
										<a href="mailto:support@azurite.dev" class="text-primary-400 hover:text-primary-300 text-sm">
											support@azurite.dev
										</a>
									</div>
								</div>

								<div class="flex items-start gap-3">
									<MessageSquare class="h-5 w-5 text-primary-400 mt-0.5" />
									<div>
										<p class="font-medium text-text-primary">Live Chat</p>
										<p class="text-text-secondary text-sm">Available during business hours</p>
									</div>
								</div>

								<div class="flex items-start gap-3">
									<Clock class="h-5 w-5 text-primary-400 mt-0.5" />
									<div>
										<p class="font-medium text-text-primary">Response Time</p>
										<p class="text-text-secondary text-sm">Typically within 24-48 hours</p>
									</div>
								</div>
							</div>
						</div>
					</div>

					<!-- Quick Links -->
					<div class="card">
						<div class="p-6">
							<h3 class="text-lg font-semibold text-text-primary mb-4">Quick Links</h3>
							<div class="space-y-3">
								<a href="/help" class="flex items-center gap-2 text-text-secondary hover:text-primary-400 transition-colors">
									<span>📚</span>
									<span>Help Center</span>
								</a>
								<a href="/bug-report" class="flex items-center gap-2 text-text-secondary hover:text-primary-400 transition-colors">
									<span>🐛</span>
									<span>Report a Bug</span>
								</a>
								<a href="/feature-request" class="flex items-center gap-2 text-text-secondary hover:text-primary-400 transition-colors">
									<span>💡</span>
									<span>Request a Feature</span>
								</a>
								<a href="/guidelines" class="flex items-center gap-2 text-text-secondary hover:text-primary-400 transition-colors">
									<span>📋</span>
									<span>Community Guidelines</span>
								</a>
								<a href="/dmca" class="flex items-center gap-2 text-text-secondary hover:text-primary-400 transition-colors">
									<span>⚖️</span>
									<span>DMCA Policy</span>
								</a>
							</div>
						</div>
					</div>

					<!-- Office Hours -->
					<div class="card">
						<div class="p-6">
							<h3 class="text-lg font-semibold text-text-primary mb-4">Support Hours</h3>
							<div class="space-y-2 text-sm">
								<div class="flex justify-between">
									<span class="text-text-secondary">Monday - Friday</span>
									<span class="text-text-primary">9:00 AM - 6:00 PM PST</span>
								</div>
								<div class="flex justify-between">
									<span class="text-text-secondary">Saturday</span>
									<span class="text-text-primary">10:00 AM - 4:00 PM PST</span>
								</div>
								<div class="flex justify-between">
									<span class="text-text-secondary">Sunday</span>
									<span class="text-text-primary">Closed</span>
								</div>
							</div>
							<p class="text-xs text-text-muted mt-3">
								Emergency issues are handled 24/7. We'll respond as quickly as possible.
							</p>
						</div>
					</div>
				</div>

				<!-- Contact Form -->
				<div class="lg:col-span-2">
					<div class="card">
						<div class="p-8">
							<div class="flex items-center gap-3 mb-6">
								<Mail class="h-6 w-6 text-primary-400" />
								<div>
									<h2 class="text-2xl font-bold text-text-primary">Send us a Message</h2>
									<p class="text-text-secondary">We'll get back to you as soon as possible</p>
								</div>
							</div>

							<form onsubmit={handleSubmit} class="space-y-6">
								<!-- Name and Email -->
								<div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
									<div>
										<label for="name" class="block text-sm font-medium text-text-primary mb-2">
											Your Name <span class="text-red-400">*</span>
										</label>
										<input
											id="name"
											type="text"
											bind:value={name}
											placeholder="Enter your full name"
											class="input w-full"
											required
											disabled={loading}
											maxlength="100"
										/>
									</div>

									<div>
										<label for="email" class="block text-sm font-medium text-text-primary mb-2">
											Email Address <span class="text-red-400">*</span>
										</label>
										<input
											id="email"
											type="email"
											bind:value={email}
											placeholder="your@email.com"
											class="input w-full"
											required
											disabled={loading}
										/>
									</div>
								</div>

								<!-- Category -->
								<div>
									<label for="category" class="block text-sm font-medium text-text-primary mb-2">
										Category
									</label>
									<select
										id="category"
										bind:value={category}
										class="select w-full"
										disabled={loading}
									>
										{#each categories as cat}
											<option value={cat.value}>{cat.label}</option>
										{/each}
									</select>
									<p class="text-xs text-text-muted mt-1">
										{getCategoryDescription(category)}
									</p>
								</div>

								<!-- Subject -->
								<div>
									<label for="subject" class="block text-sm font-medium text-text-primary mb-2">
										Subject <span class="text-red-400">*</span>
									</label>
									<input
										id="subject"
										type="text"
										bind:value={subject}
										placeholder="Brief description of your inquiry"
										class="input w-full"
										required
										disabled={loading}
										maxlength="200"
									/>
								</div>

								<!-- Message -->
								<div>
									<label for="message" class="block text-sm font-medium text-text-primary mb-2">
										Message <span class="text-red-400">*</span>
									</label>
									<textarea
										id="message"
										bind:value={message}
										placeholder="Please provide as much detail as possible about your inquiry..."
										class="textarea w-full h-32"
										required
										disabled={loading}
										maxlength="2000"
									></textarea>
									<p class="text-xs text-text-muted mt-1">
										{message.length}/2000 characters
									</p>
								</div>

								<!-- Additional Info -->
								<div class="bg-background-secondary p-4 rounded-lg border border-slate-600">
									<h3 class="font-medium text-text-primary mb-2">Tips for Better Support</h3>
									<ul class="text-sm text-text-secondary space-y-1">
										<li>• Be as specific as possible about your issue</li>
										<li>• Include error messages or screenshots if applicable</li>
										<li>• Mention your browser and operating system for technical issues</li>
										<li>• Provide steps to reproduce the problem</li>
										<li>• Include your username or account details (but never passwords)</li>
									</ul>
								</div>

								<!-- Submit Button -->
								<div class="flex gap-4">
									<button
										type="submit"
										disabled={loading || !name.trim() || !email.trim() || !subject.trim() || !message.trim()}
										class="btn btn-primary flex items-center gap-2"
									>
										{#if loading}
											<div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
											Sending...
										{:else}
											<Send class="h-4 w-4" />
											Send Message
										{/if}
									</button>

									<a href="/help" class="btn btn-outline">
										Check Help Center First
									</a>
								</div>
							</form>
						</div>
					</div>
				</div>
			</div>

			<!-- FAQ Section -->
			<div class="mt-12">
				<div class="card">
					<div class="p-8">
						<h2 class="text-2xl font-bold text-text-primary mb-6 text-center">Frequently Asked Questions</h2>
						<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
							<div>
								<h3 class="font-semibold text-text-primary mb-2">How long does mod approval take?</h3>
								<p class="text-text-secondary text-sm">
									Most mods are reviewed within 24-48 hours. Complex mods or those requiring additional
									review may take longer.
								</p>
							</div>

							<div>
								<h3 class="font-semibold text-text-primary mb-2">Can I update my mod after approval?</h3>
								<p class="text-text-secondary text-sm">
									Yes! You can update your mods anytime. Updates go through a quick review process
									to ensure they meet our guidelines.
								</p>
							</div>

							<div>
								<h3 class="font-semibold text-text-primary mb-2">How do I report inappropriate content?</h3>
								<p class="text-text-secondary text-sm">
									Use the "Report" button on any mod or comment, or contact us directly with details
									about the content that violates our guidelines.
								</p>
							</div>

							<div>
								<h3 class="font-semibold text-text-primary mb-2">Is there an API for developers?</h3>
								<p class="text-text-secondary text-sm">
									Yes! We provide a REST API for developers. Check our documentation for details
									on authentication and available endpoints.
								</p>
							</div>
						</div>
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>
