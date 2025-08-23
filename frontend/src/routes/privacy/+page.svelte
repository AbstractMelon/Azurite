<script lang="ts">
	import { Shield, Calendar, Mail, Eye, Lock, Globe, Users } from 'lucide-svelte';
	let lastUpdated = '2024-01-01';
</script>

<svelte:head>
	<title>Privacy Policy - Azurite</title>
	<meta
		name="description"
		content="Privacy Policy for the Azurite modding platform - Learn how we protect your data"
	/>
</svelte:head>

<div class="min-h-screen bg-background-primary">
	<div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
		<!-- Header -->
		<div class="text-center mb-12">
			<div
				class="inline-flex items-center justify-center w-16 h-16 bg-green-500/20 rounded-full mb-6"
			>
				<Shield class="w-8 h-8 text-green-400" />
			</div>
			<h1 class="text-4xl font-bold text-text-primary mb-4">Privacy Policy</h1>
			<p class="text-text-secondary text-lg max-w-2xl mx-auto">
				Your privacy matters to us. Learn how we collect, use, and protect your information on
				Azurite.
			</p>
			<div class="flex items-center justify-center mt-4 text-sm text-text-muted">
				<Calendar class="w-4 h-4 mr-2" />
				<span>Last updated: {lastUpdated}</span>
			</div>
		</div>

		<!-- Content -->
		<div class="card">
			<div class="p-8 prose prose-invert prose-green max-w-none">
				<section class="mb-8">
					<h2 class="text-2xl font-semibold text-text-primary mb-4">1. Information We Collect</h2>

					<h3 class="text-xl font-medium text-text-primary mb-3">Account Information</h3>
					<p class="text-text-secondary mb-4">When you create an account, we collect:</p>
					<ul class="list-disc list-inside text-text-secondary space-y-2 ml-4">
						<li>Username and display name</li>
						<li>Email address</li>
						<li>Password (encrypted and hashed)</li>
						<li>Profile information you choose to provide (bio, avatar)</li>
					</ul>

					<h3 class="text-xl font-medium text-text-primary mb-3 mt-6">Content and Activity</h3>
					<ul class="list-disc list-inside text-text-secondary space-y-2 ml-4">
						<li>Mods you upload and their metadata</li>
						<li>Comments and reviews you post</li>
						<li>Likes and downloads you make</li>
						<li>Messages you send through our platform</li>
					</ul>

					<h3 class="text-xl font-medium text-text-primary mb-3 mt-6">Technical Information</h3>
					<ul class="list-disc list-inside text-text-secondary space-y-2 ml-4">
						<li>IP address and location data</li>
						<li>Device and browser information</li>
						<li>Usage patterns and preferences</li>
						<li>Error logs and performance data</li>
					</ul>
				</section>

				<section class="mb-8">
					<h2 class="text-2xl font-semibold text-text-primary mb-4">
						2. How We Use Your Information
					</h2>
					<p class="text-text-secondary mb-4">We use your information to:</p>
					<ul class="list-disc list-inside text-text-secondary space-y-2 ml-4">
						<li>Provide and maintain the Azurite service</li>
						<li>Process your uploads and manage your content</li>
						<li>Send you notifications about your account and activity</li>
						<li>Improve our platform and develop new features</li>
						<li>Prevent fraud, spam, and abuse</li>
						<li>Comply with legal obligations</li>
						<li>Communicate with you about service updates</li>
					</ul>
				</section>

				<section class="mb-8">
					<h2 class="text-2xl font-semibold text-text-primary mb-4">3. Information Sharing</h2>
					<p class="text-text-secondary mb-4">
						We do not sell your personal information. We may share information in these situations:
					</p>

					<h3 class="text-xl font-medium text-text-primary mb-3">Public Information</h3>
					<ul class="list-disc list-inside text-text-secondary space-y-2 ml-4">
						<li>Your username, display name, and public profile information</li>
						<li>Mods you upload and their descriptions</li>
						<li>Comments and reviews you post publicly</li>
						<li>Your activity statistics (downloads, likes)</li>
					</ul>

					<h3 class="text-xl font-medium text-text-primary mb-3 mt-6">Third-Party Services</h3>
					<ul class="list-disc list-inside text-text-secondary space-y-2 ml-4">
						<li>OAuth providers (GitHub, Google, Discord) for authentication</li>
						<li>Cloud storage providers for file hosting</li>
						<li>Analytics services to understand platform usage</li>
						<li>Email services for notifications</li>
					</ul>

					<h3 class="text-xl font-medium text-text-primary mb-3 mt-6">Legal Requirements</h3>
					<p class="text-text-secondary">
						We may disclose your information when required by law, court order, or to protect the
						rights, property, or safety of Azurite, our users, or others.
					</p>
				</section>

				<section class="mb-8">
					<h2 class="text-2xl font-semibold text-text-primary mb-4">4. Data Security</h2>
					<p class="text-text-secondary mb-4">
						We implement security measures to protect your information:
					</p>
					<ul class="list-disc list-inside text-text-secondary space-y-2 ml-4">
						<li>Encryption of data in transit and at rest</li>
						<li>Secure password hashing using industry standards</li>
						<li>Regular security audits and updates</li>
						<li>Access controls and authentication requirements</li>
						<li>Monitoring for suspicious activity</li>
						<li>Incident response procedures</li>
					</ul>
				</section>

				<section class="mb-8">
					<h2 class="text-2xl font-semibold text-text-primary mb-4">5. Cookies and Tracking</h2>
					<p class="text-text-secondary mb-4">We use cookies and similar technologies to:</p>
					<ul class="list-disc list-inside text-text-secondary space-y-2 ml-4">
						<li>Keep you logged in to your account</li>
						<li>Remember your preferences and settings</li>
						<li>Analyze how you use our platform</li>
						<li>Provide personalized content and recommendations</li>
					</ul>
					<p class="text-text-secondary mt-4">
						You can control cookies through your browser settings, but some features may not work
						properly if cookies are disabled.
					</p>
				</section>

				<section class="mb-8">
					<h2 class="text-2xl font-semibold text-text-primary mb-4">6. Your Rights and Choices</h2>
					<p class="text-text-secondary mb-4">You have the right to:</p>
					<ul class="list-disc list-inside text-text-secondary space-y-2 ml-4">
						<li>Access and download your personal data</li>
						<li>Correct inaccurate or outdated information</li>
						<li>Delete your account and associated data</li>
						<li>Control email notification preferences</li>
						<li>Opt out of certain data processing activities</li>
						<li>Request information about how your data is used</li>
					</ul>
					<p class="text-text-secondary mt-4">
						To exercise these rights, please contact us using the information provided at the end of
						this policy.
					</p>
				</section>

				<section class="mb-8">
					<h2 class="text-2xl font-semibold text-text-primary mb-4">7. Data Retention</h2>
					<p class="text-text-secondary mb-4">We retain your information for as long as:</p>
					<ul class="list-disc list-inside text-text-secondary space-y-2 ml-4">
						<li>Your account remains active</li>
						<li>Required to provide our services</li>
						<li>Necessary for legal compliance</li>
						<li>Needed to resolve disputes</li>
					</ul>
					<p class="text-text-secondary mt-4">
						When you delete your account, we will delete or anonymize your personal information
						within 30 days, except where retention is required by law.
					</p>
				</section>

				<section class="mb-8">
					<h2 class="text-2xl font-semibold text-text-primary mb-4">8. Children's Privacy</h2>
					<p class="text-text-secondary">
						Azurite is not intended for children under 13 years of age. We do not knowingly collect
						personal information from children under 13. If you believe we have collected
						information from a child under 13, please contact us immediately.
					</p>
				</section>

				<section class="mb-8">
					<h2 class="text-2xl font-semibold text-text-primary mb-4">9. International Users</h2>
					<p class="text-text-secondary">
						Azurite is operated from the United States. If you are accessing our service from
						outside the US, please be aware that your information may be transferred to, stored, and
						processed in the US where our servers are located.
					</p>
				</section>

				<section class="mb-8">
					<h2 class="text-2xl font-semibold text-text-primary mb-4">10. Changes to This Policy</h2>
					<p class="text-text-secondary">
						We may update this privacy policy from time to time. We will notify you of significant
						changes by email or through prominent notices on our platform. Your continued use of
						Azurite after changes become effective constitutes acceptance of the updated policy.
					</p>
				</section>

				<section class="mb-8">
					<h2 class="text-2xl font-semibold text-text-primary mb-4">11. Contact Us</h2>
					<p class="text-text-secondary mb-4">
						If you have questions or concerns about this privacy policy, please contact us:
					</p>
					<div class="bg-slate-800/50 rounded-lg p-4 space-y-3">
						<div class="flex items-center text-text-secondary">
							<Mail class="w-5 h-5 mr-3 text-green-400" />
							<span>privacy@azurite.com</span>
						</div>
						<div class="flex items-center text-text-secondary">
							<Globe class="w-5 h-5 mr-3 text-green-400" />
							<span>https://azurite.com/contact</span>
						</div>
					</div>
				</section>

				<div class="grid grid-cols-1 md:grid-cols-3 gap-4 mt-12">
					<div class="p-4 bg-blue-900/20 border border-blue-600/50 rounded-lg">
						<div class="flex items-center mb-3">
							<Eye class="w-5 h-5 text-blue-400 mr-2" />
							<h3 class="font-semibold text-blue-300">Transparency</h3>
						</div>
						<p class="text-blue-200 text-sm">
							We are committed to being transparent about how we collect and use your data.
						</p>
					</div>
					<div class="p-4 bg-green-900/20 border border-green-600/50 rounded-lg">
						<div class="flex items-center mb-3">
							<Lock class="w-5 h-5 text-green-400 mr-2" />
							<h3 class="font-semibold text-green-300">Security</h3>
						</div>
						<p class="text-green-200 text-sm">
							Your data is protected with industry-standard security measures and encryption.
						</p>
					</div>
					<div class="p-4 bg-purple-900/20 border border-purple-600/50 rounded-lg">
						<div class="flex items-center mb-3">
							<Users class="w-5 h-5 text-purple-400 mr-2" />
							<h3 class="font-semibold text-purple-300">Control</h3>
						</div>
						<p class="text-purple-200 text-sm">
							You have full control over your data and can access, modify, or delete it at any time.
						</p>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>

<style>
	.prose ul li {
		@apply text-text-secondary;
	}
	.prose a {
		@apply text-green-400 hover:text-green-300 no-underline;
	}
	.prose a:hover {
		@apply underline;
	}
</style>
