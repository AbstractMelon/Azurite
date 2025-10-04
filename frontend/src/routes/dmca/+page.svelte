<script lang="ts">
	// No script needed for static content
</script>

<svelte:head>
	<title>DMCA Policy - Azurite</title>
	<meta
		name="description"
		content="Digital Millennium Copyright Act (DMCA) policy for Azurite mod hosting platform"
	/>
</svelte:head>

<div class="min-h-screen bg-background-primary">
	<div class="container mx-auto px-4 py-8 max-w-4xl">
		<div class="card">
			<div class="p-8">
				<h1 class="text-3xl font-bold text-text-primary mb-6">DMCA Policy</h1>
				<p class="text-text-secondary mb-6">Last updated: {new Date().toLocaleDateString()}</p>

				<div class="prose prose-invert max-w-none space-y-6">
					<section>
						<h2 class="text-2xl font-semibold text-text-primary mb-4">Overview</h2>
						<p class="text-text-secondary leading-relaxed">
							Azurite respects the intellectual property rights of others and expects its users to
							do the same. In accordance with the Digital Millennium Copyright Act of 1998 ("DMCA"),
							we will respond expeditiously to claims of copyright infringement committed using our
							service.
						</p>
					</section>

					<section>
						<h2 class="text-2xl font-semibold text-text-primary mb-4">Filing a DMCA Notice</h2>
						<p class="text-text-secondary leading-relaxed mb-4">
							If you believe that your copyrighted work has been copied in a way that constitutes
							copyright infringement and is accessible via our service, please notify our designated
							copyright agent by providing the following information:
						</p>
						<ul class="list-disc pl-6 space-y-2 text-text-secondary">
							<li>
								A physical or electronic signature of a person authorized to act on behalf of the
								owner of an exclusive right that is allegedly infringed;
							</li>
							<li>
								Identification of the copyrighted work claimed to have been infringed, or, if
								multiple copyrighted works are covered by a single notification, a representative
								list of such works;
							</li>
							<li>
								Identification of the material that is claimed to be infringing or to be the subject
								of infringing activity and that is to be removed or access to which is to be
								disabled, and information reasonably sufficient to permit us to locate the material;
							</li>
							<li>
								Information reasonably sufficient to permit us to contact the complaining party,
								such as an address, telephone number, and, if available, an electronic mail address;
							</li>
							<li>
								A statement that the complaining party has a good faith belief that use of the
								material in the manner complained of is not authorized by the copyright owner, its
								agent, or the law;
							</li>
							<li>
								A statement that the information in the notification is accurate, and under penalty
								of perjury, that the complaining party is authorized to act on behalf of the owner
								of an exclusive right that is allegedly infringed.
							</li>
						</ul>
					</section>

					<section>
						<h2 class="text-2xl font-semibold text-text-primary mb-4">
							Designated Copyright Agent
						</h2>
						<p class="text-text-secondary leading-relaxed mb-4">
							Our designated copyright agent for notice of claims of copyright infringement can be
							reached at:
						</p>
						<div class="bg-background-secondary p-4 rounded-lg border border-slate-600">
							<p class="text-text-primary"><strong>DMCA Agent</strong></p>
							<p class="text-text-secondary">Azurite Platform</p>
							<p class="text-text-secondary">Email: dmca@azurite.dev</p>
							<p class="text-text-secondary">Subject Line: "DMCA Takedown Notice"</p>
						</div>
					</section>

					<section>
						<h2 class="text-2xl font-semibold text-text-primary mb-4">Counter-Notification</h2>
						<p class="text-text-secondary leading-relaxed mb-4">
							If you believe that material you posted was removed or access to it was disabled by
							mistake or misidentification, you may file a counter-notification with us by providing
							the following information:
						</p>
						<ul class="list-disc pl-6 space-y-2 text-text-secondary">
							<li>Your physical or electronic signature;</li>
							<li>
								Identification of the material that has been removed or to which access has been
								disabled and the location at which the material appeared before it was removed or
								access to it was disabled;
							</li>
							<li>
								A statement under penalty of perjury that you have a good faith belief that the
								material was removed or disabled as a result of mistake or misidentification;
							</li>
							<li>
								Your name, address, and telephone number, and a statement that you consent to the
								jurisdiction of Federal District Court for the judicial district in which the
								address is located;
							</li>
							<li>
								A statement that you will accept service of process from the person who provided
								notification of the alleged infringement.
							</li>
						</ul>
					</section>

					<section>
						<h2 class="text-2xl font-semibold text-text-primary mb-4">Repeat Infringer Policy</h2>
						<p class="text-text-secondary leading-relaxed">
							In accordance with the DMCA and other applicable law, we have adopted a policy of
							terminating, in appropriate circumstances and at our sole discretion, users who are
							deemed to be repeat infringers. We may also at our sole discretion limit access to the
							service and/or terminate the accounts of any users who infringe any intellectual
							property rights of others, whether or not there is any repeat infringement.
						</p>
					</section>

					<section>
						<h2 class="text-2xl font-semibold text-text-primary mb-4">Processing Time</h2>
						<p class="text-text-secondary leading-relaxed">
							We will process all valid DMCA notices as quickly as possible, typically within 24-48
							hours of receipt. Upon receiving a valid takedown notice, we will:
						</p>
						<ul class="list-disc pl-6 space-y-2 text-text-secondary mt-4">
							<li>Remove or disable access to the allegedly infringing material</li>
							<li>Notify the user who posted the material of the removal</li>
							<li>Provide the user with information about filing a counter-notification</li>
						</ul>
					</section>

					<section>
						<h2 class="text-2xl font-semibold text-text-primary mb-4">False Claims</h2>
						<p class="text-text-secondary leading-relaxed">
							Please note that under Section 512(f) of the DMCA, any person who knowingly materially
							misrepresents that material or activity is infringing may be subject to liability for
							damages. We reserve the right to seek damages from any party that submits a false DMCA
							notice.
						</p>
					</section>

					<section>
						<h2 class="text-2xl font-semibold text-text-primary mb-4">Questions</h2>
						<p class="text-text-secondary leading-relaxed">
							If you have any questions about this DMCA policy or need assistance with filing a
							notice, please contact us at <a
								href="mailto:dmca@azurite.dev"
								class="text-primary-400 hover:text-primary-300 transition-colors"
								>dmca@azurite.dev</a
							>.
						</p>
					</section>
				</div>
			</div>
		</div>
	</div>
</div>

<style>
	/* Custom prose styles for better readability */
	.prose h2 {
		border-bottom: 1px solid rgb(71 85 105);
		padding-bottom: 0.5rem;
	}

	.prose ul li {
		margin-bottom: 0.5rem;
	}
</style>
