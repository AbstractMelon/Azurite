<script lang="ts">
	import { FileText, Calendar, Mail, TriangleAlert } from 'lucide-svelte';

	let lastUpdated = '2024-01-01';
</script>

<svelte:head>
	<title>Terms of Service - Azurite</title>
	<meta name="description" content="Terms of Service for using the Azurite modding platform" />
</svelte:head>

<div class="min-h-screen bg-background-primary">
	<div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
		<!-- Header -->
		<div class="text-center mb-12">
			<div
				class="inline-flex items-center justify-center w-16 h-16 bg-primary-500/20 rounded-full mb-6"
			>
				<FileText class="w-8 h-8 text-primary-400" />
			</div>
			<h1 class="text-4xl font-bold text-text-primary mb-4">Terms of Service</h1>
			<p class="text-text-secondary text-lg max-w-2xl mx-auto">
				Please read these terms carefully before using the Azurite platform
			</p>
			<div class="flex items-center justify-center mt-4 text-sm text-text-muted">
				<Calendar class="w-4 h-4 mr-2" />
				<span>Last updated: {lastUpdated}</span>
			</div>
		</div>

		<!-- Content -->
		<div class="card">
			<div class="p-8 prose prose-invert prose-blue max-w-none">
				<section class="mb-8">
					<h2 class="text-2xl font-semibold text-text-primary mb-4">1. Acceptance of Terms</h2>
					<p class="text-text-secondary">
						By accessing and using the Azurite platform ("Service"), you accept and agree to be
						bound by the terms and provision of this agreement. If you do not agree to abide by the
						above, please do not use this service.
					</p>
				</section>

				<section class="mb-8">
					<h2 class="text-2xl font-semibold text-text-primary mb-4">2. Description of Service</h2>
					<p class="text-text-secondary mb-4">Azurite is a platform that allows users to:</p>
					<ul class="list-disc list-inside text-text-secondary space-y-2 ml-4">
						<li>Upload, share, and download game modifications (mods)</li>
						<li>Create user accounts and profiles</li>
						<li>Comment and interact with other users</li>
						<li>Access documentation and community resources</li>
					</ul>
				</section>

				<section class="mb-8">
					<h2 class="text-2xl font-semibold text-text-primary mb-4">3. User Accounts</h2>
					<p class="text-text-secondary mb-4">
						To use certain features of the Service, you must register for an account. You agree to:
					</p>
					<ul class="list-disc list-inside text-text-secondary space-y-2 ml-4">
						<li>Provide accurate, current, and complete information</li>
						<li>Maintain and update your information to keep it accurate</li>
						<li>Maintain the security of your password</li>
						<li>Accept responsibility for all activities under your account</li>
						<li>Notify us immediately of unauthorized use of your account</li>
					</ul>
				</section>

				<section class="mb-8">
					<h2 class="text-2xl font-semibold text-text-primary mb-4">4. Content Guidelines</h2>
					<p class="text-text-secondary mb-4">
						When uploading mods or other content, you agree that your content will not:
					</p>
					<ul class="list-disc list-inside text-text-secondary space-y-2 ml-4">
						<li>Contain malware, viruses, or other harmful code</li>
						<li>Infringe on intellectual property rights</li>
						<li>Contain illegal, harmful, or offensive material</li>
						<li>Violate the terms of service of the games being modified</li>
						<li>Include personal information of others without consent</li>
					</ul>
				</section>

				<section class="mb-8">
					<h2 class="text-2xl font-semibold text-text-primary mb-4">5. User Conduct</h2>
					<p class="text-text-secondary mb-4">You agree not to:</p>
					<ul class="list-disc list-inside text-text-secondary space-y-2 ml-4">
						<li>Harass, abuse, or harm other users</li>
						<li>Use the service for any illegal purposes</li>
						<li>Attempt to gain unauthorized access to our systems</li>
						<li>Spam or flood the platform with unwanted content</li>
						<li>Create multiple accounts to circumvent restrictions</li>
						<li>Engage in any activity that disrupts the service</li>
					</ul>
				</section>

				<section class="mb-8">
					<h2 class="text-2xl font-semibold text-text-primary mb-4">6. Intellectual Property</h2>
					<p class="text-text-secondary mb-4">
						You retain ownership of the content you upload. By uploading content, you grant us a
						non-exclusive, worldwide, royalty-free license to use, display, and distribute your
						content on the platform.
					</p>
					<p class="text-text-secondary">
						The Azurite platform, including its design, functionality, and underlying code, is
						protected by intellectual property laws and is owned by Azurite or its licensors.
					</p>
				</section>

				<section class="mb-8">
					<h2 class="text-2xl font-semibold text-text-primary mb-4">7. Privacy</h2>
					<p class="text-text-secondary">
						Your privacy is important to us. Please review our <a
							href="/privacy"
							class="text-primary-400 hover:text-primary-300">Privacy Policy</a
						>
						to understand how we collect, use, and protect your information.
					</p>
				</section>

				<section class="mb-8">
					<h2 class="text-2xl font-semibold text-text-primary mb-4">
						8. Moderation and Enforcement
					</h2>
					<p class="text-text-secondary mb-4">We reserve the right to:</p>
					<ul class="list-disc list-inside text-text-secondary space-y-2 ml-4">
						<li>Review and moderate all content before or after publication</li>
						<li>Remove content that violates these terms</li>
						<li>Suspend or terminate accounts for violations</li>
						<li>Take legal action when necessary</li>
					</ul>
				</section>

				<section class="mb-8">
					<h2 class="text-2xl font-semibold text-text-primary mb-4">9. Disclaimers</h2>
					<p class="text-text-secondary mb-4">
						The service is provided "as is" without warranties of any kind. We do not guarantee:
					</p>
					<ul class="list-disc list-inside text-text-secondary space-y-2 ml-4">
						<li>Uninterrupted or error-free service</li>
						<li>The safety or functionality of user-uploaded mods</li>
						<li>The accuracy of information provided by users</li>
						<li>Compatibility with all games or systems</li>
					</ul>
				</section>

				<section class="mb-8">
					<h2 class="text-2xl font-semibold text-text-primary mb-4">10. Limitation of Liability</h2>
					<p class="text-text-secondary">
						To the maximum extent permitted by law, Azurite shall not be liable for any indirect,
						incidental, special, consequential, or punitive damages, including but not limited to
						loss of data, revenue, or profits arising from your use of the service.
					</p>
				</section>

				<section class="mb-8">
					<h2 class="text-2xl font-semibold text-text-primary mb-4">11. Termination</h2>
					<p class="text-text-secondary">
						Either party may terminate this agreement at any time. Upon termination, your right to
						use the service will cease, but these terms will continue to apply to your prior use of
						the service.
					</p>
				</section>

				<section class="mb-8">
					<h2 class="text-2xl font-semibold text-text-primary mb-4">12. Changes to Terms</h2>
					<p class="text-text-secondary">
						We may update these terms from time to time. Significant changes will be communicated to
						users via email or platform notifications. Continued use of the service after changes
						constitutes acceptance of the new terms.
					</p>
				</section>

				<section class="mb-8">
					<h2 class="text-2xl font-semibold text-text-primary mb-4">13. Contact Information</h2>
					<p class="text-text-secondary mb-4">
						If you have questions about these terms, please contact us:
					</p>
					<div class="bg-slate-800/50 rounded-lg p-4">
						<div class="flex items-center text-text-secondary">
							<Mail class="w-5 h-5 mr-3 text-primary-400" />
							<span>support@azurite.com</span>
						</div>
					</div>
				</section>

				<div class="mt-12 p-6 bg-yellow-900/20 border border-yellow-600/50 rounded-lg">
					<div class="flex items-start">
						<TriangleAlert class="w-6 h-6 text-yellow-400 mr-3 mt-0.5 flex-shrink-0" />
						<div>
							<h3 class="text-lg font-semibold text-yellow-300 mb-2">Important Notice</h3>
							<p class="text-yellow-200">
								By using Azurite, you acknowledge that you have read, understood, and agree to be
								bound by these Terms of Service. If you do not agree with any part of these terms,
								you must not use our service.
							</p>
						</div>
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
		@apply text-primary-400 hover:text-primary-300 no-underline;
	}
	.prose a:hover {
		@apply underline;
	}
</style>
