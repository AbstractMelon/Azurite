<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { modsApi, commentsApi } from '$lib/api/client';
	import { toast } from '$lib/stores/notifications';
	import { user, isAuthenticated } from '$lib/stores/auth';
	import Loading from '$lib/components/Loading.svelte';
	import type { Mod, Comment } from '$lib/types';
	import {
		Download,
		Heart,
		Calendar,
		Package,
		User,
		ExternalLink,
		Mail,
		Share2,
		MessageCircle,
		Edit,
		Trash2,
		Reply,
		FileText
	} from 'lucide-svelte';
	import { marked } from 'marked';
	import DOMPurify from 'dompurify';

	// URL params
	$: gameSlug = $page.params.gameSlug;
	$: modSlug = $page.params.modSlug;

	// Data
	let mod: Mod | null = null;
	let comments: Comment[] = [];
	let isLoadingMod = true;
	let isLoadingComments = true;
	let isLiked = false;
	let isLiking = false;

	// Comments
	let newComment = '';
	let isPostingComment = false;
	let replyingTo: number | null = null;
	let editingComment: number | null = null;
	let editCommentText = '';

	// UI state
	let showAllFiles = false;

	onMount(async () => {
		await loadMod();
		await loadComments();
	});

	// Load mod data
	async function loadMod() {
		isLoadingMod = true;
		try {
			const response = await modsApi.getMod(gameSlug, modSlug);
			if (response.success && response.data) {
				mod = response.data;
				isLiked = mod.is_liked || false;
			} else {
				toast.error('Mod not found', 'The requested mod could not be found.');
				goto(`/games/${gameSlug}`);
			}
		} catch (error) {
			console.error('Error loading mod:', error);
			toast.error('Error', 'Failed to load mod information');
		} finally {
			isLoadingMod = false;
		}
	}

	// Load comments
	async function loadComments() {
		if (!mod) return;

		isLoadingComments = true;
		try {
			const response = await commentsApi.getModComments(mod.id);
			if (response.success && response.data) {
				comments = response.data.data || [];
			}
		} catch (error) {
			console.error('Error loading comments:', error);
		} finally {
			isLoadingComments = false;
		}
	}

	// Handle like/unlike
	async function toggleLike() {
		if (!$isAuthenticated) {
			goto('/auth/login?redirect=' + encodeURIComponent($page.url.pathname));
			return;
		}

		if (isLiking) return;
		isLiking = true;

		try {
			const response = isLiked ? await modsApi.unlikeMod(mod.id) : await modsApi.likeMod(mod.id);

			if (response.success) {
				isLiked = !isLiked;
				mod.likes += isLiked ? 1 : -1;
			} else {
				toast.error('Failed to update like', response.error);
			}
		} catch (error) {
			console.error('Error toggling like:', error);
			toast.error('Error', 'Failed to update like');
		} finally {
			isLiking = false;
		}
	}

	// Handle download
  function handleDownload() {
    if (mod?.files?.length > 0) {
      window.open(`/download/${gameSlug}/${modSlug}`, '_blank');
    }
  }

	// Handle comment submission
	async function postComment(parentId?: number) {
		if (!$isAuthenticated) {
			goto('/auth/login?redirect=' + encodeURIComponent($page.url.pathname));
			return;
		}

		const commentText = parentId ? newComment : newComment;
		if (!commentText.trim()) return;

		isPostingComment = true;
		try {
			const response = await commentsApi.createComment(mod.id, {
				content: commentText.trim(),
				parent_id: parentId || undefined
			});

			if (response.success && response.data) {
				await loadComments(); // Reload comments
				newComment = '';
				replyingTo = null;
				toast.success('Comment posted', 'Your comment has been posted successfully.');
			} else {
				toast.error('Failed to post comment', response.error);
			}
		} catch (error) {
			console.error('Error posting comment:', error);
			toast.error('Error', 'Failed to post comment');
		} finally {
			isPostingComment = false;
		}
	}

	// Handle comment editing
	async function updateComment(commentId: number) {
		if (!editCommentText.trim()) return;

		try {
			const response = await commentsApi.updateComment(commentId, {
				content: editCommentText.trim()
			});

			if (response.success) {
				await loadComments(); // Reload comments
				editingComment = null;
				editCommentText = '';
				toast.success('Comment updated', 'Your comment has been updated.');
			} else {
				toast.error('Failed to update comment', response.error);
			}
		} catch (error) {
			console.error('Error updating comment:', error);
			toast.error('Error', 'Failed to update comment');
		}
	}

	// Handle comment deletion
	async function deleteComment(commentId: number) {
		if (!confirm('Are you sure you want to delete this comment?')) return;

		try {
			const response = await commentsApi.deleteComment(commentId);
			if (response.success) {
				await loadComments(); // Reload comments
				toast.success('Comment deleted', 'Your comment has been deleted.');
			} else {
				toast.error('Failed to delete comment', response.error);
			}
		} catch (error) {
			console.error('Error deleting comment:', error);
			toast.error('Error', 'Failed to delete comment');
		}
	}

	// Start editing comment
	function startEditing(comment: any) {
		editingComment = comment.id;
		editCommentText = comment.content;
	}

	// Cancel editing
	function cancelEditing() {
		editingComment = null;
		editCommentText = '';
	}

	// Start replying
	function startReplying(commentId: number) {
		replyingTo = commentId;
		newComment = '';
	}

	// Format numbers
	function formatNumber(num: number): string {
		if (num >= 1000000) {
			return (num / 1000000).toFixed(1) + 'M';
		}
		if (num >= 1000) {
			return (num / 1000).toFixed(1) + 'K';
		}
		return num.toString();
	}

	// Format date
	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString();
	}

	// Format relative time
	function formatRelativeTime(dateString: string): string {
		const date = new Date(dateString);
		const now = new Date();
		const diffTime = Math.abs(now.getTime() - date.getTime());
		const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

		if (diffDays === 1) return '1 day ago';
		if (diffDays < 7) return `${diffDays} days ago`;
		if (diffDays < 30) return `${Math.ceil(diffDays / 7)} weeks ago`;
		if (diffDays < 365) return `${Math.ceil(diffDays / 30)} months ago`;
		return `${Math.ceil(diffDays / 365)} years ago`;
	}

	// Format file size
	function formatFileSize(bytes: number): string {
		if (bytes === 0) return '0 Bytes';
		const k = 1024;
		const sizes = ['Bytes', 'KB', 'MB', 'GB'];
		const i = Math.floor(Math.log(bytes) / Math.log(k));
		return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
	}

	// Render markdown
	function renderMarkdown(content: string): string {
		return DOMPurify.sanitize(marked(content));
	}

	// Check if user can edit comment
	function canEditComment(comment: any): boolean {
		return $user && ($user.id === comment.user_id || $user.role === 'admin');
	}

	// Share mod
	async function shareMod() {
		if (navigator.share) {
			try {
				await navigator.share({
					title: mod.name,
					text: mod.short_description,
					url: window.location.href
				});
			} catch {
				copyToClipboard();
			}
		} else {
			copyToClipboard();
		}
	}

	// Copy URL to clipboard
	function copyToClipboard() {
		navigator.clipboard.writeText(window.location.href).then(() => {
			toast.success('Link copied', 'Mod link copied to clipboard');
		});
	}
</script>

<svelte:head>
	<title>{mod?.name || 'Mod'} - {mod?.game?.name || 'Game'} - Azurite</title>
	<meta
		name="description"
		content={mod
			? `${mod.short_description} Download ${mod.name} for ${mod.game?.name} on Azurite.`
			: 'Browse game modifications on Azurite.'}
	/>
</svelte:head>

{#if isLoadingMod}
	<div class="min-h-screen flex items-center justify-center">
		<Loading size="lg" text="Loading mod..." />
	</div>
{:else if mod}
	<div class="min-h-screen bg-background-primary">
		<!-- Mod Header -->
		<div class="bg-gradient-to-r from-slate-800/50 to-slate-700/50 border-b border-slate-700">
			<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
				<div
					class="flex flex-col lg:flex-row items-start lg:items-center space-y-6 lg:space-y-0 lg:space-x-8"
				>
					<!-- Mod Icon -->
					<div class="flex-shrink-0">
						{#if mod.icon}
							<img
								src={mod.icon}
								alt={mod.name}
								class="w-24 h-24 md:w-32 md:h-32 rounded-2xl border border-slate-600 shadow-xl"
							/>
						{:else}
							<div
								class="w-24 h-24 md:w-32 md:h-32 bg-gradient-to-br from-slate-600 to-slate-700 rounded-2xl flex items-center justify-center shadow-xl"
							>
								<Package class="w-12 h-12 md:w-16 md:h-16 text-text-muted" />
							</div>
						{/if}
					</div>

					<!-- Mod Info -->
					<div class="flex-1 min-w-0">
						<div class="flex items-center space-x-2 text-sm text-text-muted mb-2">
							<a href="/games/{mod.game?.slug}" class="hover:text-primary-400 transition-colors">
								{mod.game?.name}
							</a>
							<span>•</span>
							<span>by</span>
							<a
								href="/profile/{mod.owner?.username}"
								class="hover:text-primary-400 transition-colors"
							>
								{mod.owner?.display_name || mod.owner?.username}
							</a>
						</div>

						<h1 class="text-3xl md:text-4xl font-bold text-text-primary mb-3">
							{mod.name}
						</h1>

						<p class="text-text-secondary text-lg mb-4 leading-relaxed">
							{mod.short_description}
						</p>

						<!-- Stats -->
						<div class="flex flex-wrap items-center gap-6 text-text-muted">
							<div class="flex items-center">
								<Download class="w-5 h-5 mr-2" />
								<span class="font-medium">{formatNumber(mod.downloads || 0)}</span>
								<span class="ml-1">downloads</span>
							</div>
							<div class="flex items-center">
								<Heart class="w-5 h-5 mr-2" />
								<span class="font-medium">{formatNumber(mod.likes || 0)}</span>
								<span class="ml-1">likes</span>
							</div>
							<div class="flex items-center">
								<Calendar class="w-5 h-5 mr-2" />
								<span>Updated {formatRelativeTime(mod.updated_at)}</span>
							</div>
						</div>
					</div>

					<!-- Actions -->
					<div class="flex flex-col sm:flex-row gap-3 w-full lg:w-auto">
						<button
							on:click={handleDownload}
							disabled={!mod.files || mod.files.length === 0}
							class="btn btn-primary btn-lg flex-1 sm:flex-initial"
						>
							<Download class="w-5 h-5 mr-2" />
							Download
						</button>

						<div class="flex gap-2">
							<button
								on:click={toggleLike}
								disabled={isLiking}
								class="btn {isLiked ? 'btn-primary' : 'btn-outline'} flex items-center"
								title={isLiked ? 'Unlike' : 'Like'}
							>
								<Heart class="w-5 h-5 {isLiked ? 'fill-current' : ''}" />
							</button>

							<button on:click={shareMod} class="btn btn-outline" title="Share">
								<Share2 class="w-5 h-5" />
							</button>
						</div>
					</div>
				</div>
			</div>
		</div>

		<!-- Content -->
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
			<div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
				<!-- Main Content -->
				<div class="lg:col-span-2 space-y-8">
					<!-- Description -->
					<div class="card">
						<div class="p-6">
							<h2 class="text-xl font-semibold text-text-primary mb-4">Description</h2>
							{#if mod.description}
								<div class="markdown-content">
									{@html renderMarkdown(mod.description)}
								</div>
							{:else}
								<p class="text-text-muted italic">No detailed description provided.</p>
							{/if}
						</div>
					</div>

					<!-- Comments Section -->
					<div class="card">
						<div class="p-6">
							<div class="flex items-center justify-between mb-6">
								<h2 class="text-xl font-semibold text-text-primary flex items-center">
									<MessageCircle class="w-5 h-5 mr-2" />
									Comments ({comments.length})
								</h2>
							</div>

							<!-- Post Comment -->
							{#if $isAuthenticated}
								<div class="mb-6">
									<div class="flex items-start space-x-3">
										{#if $user?.avatar}
											<img
												src={$user.avatar}
												alt={$user.display_name}
												class="w-10 h-10 rounded-full border border-slate-600"
											/>
										{:else}
											<div
												class="w-10 h-10 bg-slate-600 rounded-full flex items-center justify-center"
											>
												<User class="w-5 h-5 text-text-muted" />
											</div>
										{/if}
										<div class="flex-1">
											<textarea
												bind:value={newComment}
												placeholder="Add a comment..."
												class="textarea mb-3"
												rows="3"
												disabled={isPostingComment}
											></textarea>
											<div class="flex justify-end">
												<button
													on:click={() => postComment()}
													disabled={!newComment.trim() || isPostingComment}
													class="btn btn-primary btn-sm"
												>
													{#if isPostingComment}
														<Loading size="sm" inline />
													{:else}
														Post Comment
													{/if}
												</button>
											</div>
										</div>
									</div>
								</div>
							{:else}
								<div class="mb-6 p-4 bg-slate-800 rounded-lg text-center">
									<p class="text-text-secondary mb-3">Sign in to leave a comment</p>
									<a
										href="/auth/login?redirect={encodeURIComponent($page.url.pathname)}"
										class="btn btn-primary btn-sm"
									>
										Sign In
									</a>
								</div>
							{/if}

							<!-- Comments List -->
							{#if isLoadingComments}
								<div class="text-center py-8">
									<Loading size="md" text="Loading comments..." />
								</div>
							{:else if comments.length === 0}
								<div class="text-center py-8 text-text-muted">
									<MessageCircle class="w-12 h-12 mx-auto mb-3 opacity-50" />
									<p>No comments yet. Be the first to share your thoughts!</p>
								</div>
							{:else}
								<div class="space-y-4">
									{#each comments as comment (comment.id)}
										<div class="border-l-2 border-slate-700 pl-4">
											<div class="flex items-start space-x-3">
												<!-- Avatar -->
												{#if comment.user?.avatar}
													<img
														src={comment.user.avatar}
														alt={comment.user.display_name}
														class="w-8 h-8 rounded-full border border-slate-600 flex-shrink-0"
													/>
												{:else}
													<div
														class="w-8 h-8 bg-slate-600 rounded-full flex items-center justify-center flex-shrink-0"
													>
														<User class="w-4 h-4 text-text-muted" />
													</div>
												{/if}

												<div class="flex-1 min-w-0">
													<!-- Comment Header -->
													<div class="flex items-center space-x-2 mb-1">
														<span class="font-medium text-text-primary">
															{comment.user?.display_name || comment.user?.username || 'Anonymous'}
														</span>
														<span class="text-text-muted text-sm">
															{formatRelativeTime(comment.created_at)}
														</span>
														{#if comment.updated_at !== comment.created_at}
															<span class="text-text-muted text-sm italic">(edited)</span>
														{/if}
													</div>

													<!-- Comment Content -->
													{#if editingComment === comment.id}
														<div class="mb-3">
															<textarea bind:value={editCommentText} class="textarea mb-2" rows="3"
															></textarea>
															<div class="flex space-x-2">
																<button
																	on:click={() => updateComment(comment.id)}
																	class="btn btn-primary btn-sm"
																>
																	Save
																</button>
																<button on:click={cancelEditing} class="btn btn-outline btn-sm">
																	Cancel
																</button>
															</div>
														</div>
													{:else}
														<div class="text-text-secondary mb-3 leading-relaxed">
															{comment.content}
														</div>
													{/if}

													<!-- Comment Actions -->
													{#if editingComment !== comment.id}
														<div class="flex items-center space-x-4 text-sm">
															{#if $isAuthenticated}
																<button
																	on:click={() => startReplying(comment.id)}
																	class="text-text-muted hover:text-primary-400 transition-colors flex items-center"
																>
																	<Reply class="w-4 h-4 mr-1" />
																	Reply
																</button>
															{/if}

															{#if canEditComment(comment)}
																<button
																	on:click={() => startEditing(comment)}
																	class="text-text-muted hover:text-primary-400 transition-colors flex items-center"
																>
																	<Edit class="w-4 h-4 mr-1" />
																	Edit
																</button>
																<button
																	on:click={() => deleteComment(comment.id)}
																	class="text-text-muted hover:text-red-400 transition-colors flex items-center"
																>
																	<Trash2 class="w-4 h-4 mr-1" />
																	Delete
																</button>
															{/if}
														</div>
													{/if}

													<!-- Reply Form -->
													{#if replyingTo === comment.id}
														<div class="mt-3 ml-4">
															<textarea
																bind:value={newComment}
																placeholder="Write a reply..."
																class="textarea mb-2"
																rows="2"
															></textarea>
															<div class="flex space-x-2">
																<button
																	on:click={() => postComment(comment.id)}
																	disabled={!newComment.trim() || isPostingComment}
																	class="btn btn-primary btn-sm"
																>
																	{#if isPostingComment}
																		<Loading size="sm" inline />
																	{:else}
																		Post Reply
																	{/if}
																</button>
																<button
																	on:click={() => (replyingTo = null)}
																	class="btn btn-outline btn-sm"
																>
																	Cancel
																</button>
															</div>
														</div>
													{/if}
												</div>
											</div>
										</div>
									{/each}
								</div>
							{/if}
						</div>
					</div>
				</div>

				<!-- Sidebar -->
				<div class="space-y-6">
					<!-- Download Info -->
					<div class="card">
						<div class="p-4">
							<h3 class="text-lg font-semibold text-text-primary mb-4">Download</h3>

							{#if mod.files && mod.files.length > 0}
								<div class="space-y-3">
									{#each mod.files.slice(0, showAllFiles ? undefined : 3) as file (file.filename)}
										<div class="flex items-center justify-between p-3 bg-slate-800 rounded-lg">
											<div class="flex items-center min-w-0">
												<FileText class="w-5 h-5 text-text-muted mr-2 flex-shrink-0" />
												<div class="min-w-0">
													<p class="text-sm font-medium text-text-primary truncate">
														{file.filename}
													</p>
													<p class="text-xs text-text-muted">
														{formatFileSize(file.file_size)}
														{#if file.is_main}
															<span class="ml-2 text-primary-400">• Main</span>
														{/if}
													</p>
												</div>
											</div>
										</div>
									{/each}

									{#if mod.files.length > 3}
										<button
											on:click={() => (showAllFiles = !showAllFiles)}
											class="text-primary-400 hover:text-primary-300 text-sm font-medium transition-colors"
										>
											{showAllFiles ? 'Show Less' : `Show ${mod.files.length - 3} More`}
										</button>
									{/if}
								</div>

								<button on:click={handleDownload} class="btn btn-primary w-full mt-4">
									<Download class="w-4 h-4 mr-2" />
									Download Now
								</button>
							{:else}
								<p class="text-text-muted text-sm">No files available for download.</p>
							{/if}
						</div>
					</div>

					<!-- Mod Details -->
					<div class="card">
						<div class="p-4">
							<h3 class="text-lg font-semibold text-text-primary mb-4">Details</h3>
							<div class="space-y-3 text-sm">
								<div>
									<span class="text-text-muted">Version:</span>
									<span class="text-text-primary ml-2">{mod.version}</span>
								</div>
								<div>
									<span class="text-text-muted">Game Version:</span>
									<span class="text-text-primary ml-2">{mod.game_version}</span>
								</div>
								<div>
									<span class="text-text-muted">Created:</span>
									<span class="text-text-primary ml-2">{formatDate(mod.created_at)}</span>
								</div>
								<div>
									<span class="text-text-muted">Last Updated:</span>
									<span class="text-text-primary ml-2">{formatDate(mod.updated_at)}</span>
								</div>
							</div>
						</div>
					</div>

					<!-- Tags -->
					{#if mod.tags && mod.tags.length > 0}
						<div class="card">
							<div class="p-4">
								<h3 class="text-lg font-semibold text-text-primary mb-3">Tags</h3>
								<div class="flex flex-wrap gap-2">
									{#each mod.tags as tag (tag.id)}
										<a
											href="/games/{mod.game?.slug}?tags={tag.slug}"
											class="badge badge-secondary hover:bg-slate-600 transition-colors"
										>
											{tag.name}
										</a>
									{/each}
								</div>
							</div>
						</div>
					{/if}

					<!-- Author Info -->
					{#if mod.owner}
						<div class="card">
							<div class="p-4">
								<h3 class="text-lg font-semibold text-text-primary mb-3">Author</h3>
								<div class="flex items-center space-x-3">
									{#if mod.owner.avatar}
										<img
											src={mod.owner.avatar}
											alt={mod.owner.display_name}
											class="w-12 h-12 rounded-full border border-slate-600"
										/>
									{:else}
										<div
											class="w-12 h-12 bg-slate-600 rounded-full flex items-center justify-center"
										>
											<User class="w-6 h-6 text-text-muted" />
										</div>
									{/if}
									<div>
										<a
											href="/profile/{mod.owner.username}"
											class="text-text-primary font-medium hover:text-primary-400 transition-colors"
										>
											{mod.owner.display_name || mod.owner.username}
										</a>
										{#if mod.owner.bio}
											<p class="text-text-muted text-sm mt-1 line-clamp-2">
												{mod.owner.bio}
											</p>
										{/if}
									</div>
								</div>
							</div>
						</div>
					{/if}

					<!-- Links -->
					<div class="card">
						<div class="p-4">
							<h3 class="text-lg font-semibold text-text-primary mb-3">Links</h3>
							<div class="space-y-2">
								{#if mod.source_website}
									<a
										href={mod.source_website}
										target="_blank"
										rel="noopener noreferrer"
										class="flex items-center text-primary-400 hover:text-primary-300 transition-colors text-sm"
									>
										<ExternalLink class="w-4 h-4 mr-2" />
										Source Code
									</a>
								{/if}
								{#if mod.contact_info}
									<a
										href={mod.contact_info.startsWith('http')
											? mod.contact_info
											: `mailto:${mod.contact_info}`}
										target="_blank"
										rel="noopener noreferrer"
										class="flex items-center text-primary-400 hover:text-primary-300 transition-colors text-sm"
									>
										<Mail class="w-4 h-4 mr-2" />
										Contact Author
									</a>
								{/if}
								<button
									on:click={shareMod}
									class="flex items-center text-primary-400 hover:text-primary-300 transition-colors text-sm w-full text-left"
								>
									<Share2 class="w-4 h-4 mr-2" />
									Share Mod
								</button>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
{/if}

<style>
	.line-clamp-2 {
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
</style>
