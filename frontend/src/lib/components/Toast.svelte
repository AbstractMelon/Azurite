<script lang="ts">
	import { fly } from 'svelte/transition';
	import { notifications, type Notification } from '$lib/stores/notifications';
	import { CheckCircle, AlertCircle, AlertTriangle, Info, X } from 'lucide-svelte';

	let notificationList: Notification[] = [];

	// Subscribe to notifications store
	notifications.subscribe((state) => {
		notificationList = state.notifications;
	});

	// Get icon for notification type
	function getIcon(type: string) {
		switch (type) {
			case 'success':
				return CheckCircle;
			case 'error':
				return AlertCircle;
			case 'warning':
				return AlertTriangle;
			case 'info':
			default:
				return Info;
		}
	}

	// Get styles for notification type
	function getTypeStyles(type: string) {
		switch (type) {
			case 'success':
				return {
					container: 'bg-green-900/90 border-green-600 text-green-100',
					icon: 'text-green-400',
					button: 'text-green-300 hover:text-green-100'
				};
			case 'error':
				return {
					container: 'bg-red-900/90 border-red-600 text-red-100',
					icon: 'text-red-400',
					button: 'text-red-300 hover:text-red-100'
				};
			case 'warning':
				return {
					container: 'bg-yellow-900/90 border-yellow-600 text-yellow-100',
					icon: 'text-yellow-400',
					button: 'text-yellow-300 hover:text-yellow-100'
				};
			case 'info':
			default:
				return {
					container: 'bg-blue-900/90 border-blue-600 text-blue-100',
					icon: 'text-blue-400',
					button: 'text-blue-300 hover:text-blue-100'
				};
		}
	}

	// Remove notification
	function removeNotification(id: string) {
		notifications.remove(id);
	}

	// Handle notification action
	function handleAction(action: () => void, notificationId: string) {
		action();
		removeNotification(notificationId);
	}
</script>

<!-- Toast Container -->
<div class="fixed top-4 right-4 z-50 space-y-2 max-w-sm w-full pointer-events-none">
	{#each notificationList as notification (notification.id)}
		{@const styles = getTypeStyles(notification.type)}
		{@const IconComponent = getIcon(notification.type)}

		<div
			class="pointer-events-auto backdrop-blur-sm border rounded-lg shadow-lg overflow-hidden {styles.container}"
			transition:fly={{ x: 300, duration: 300 }}
		>
			<div class="p-4">
				<div class="flex items-start space-x-3">
					<!-- Icon -->
					<div class="flex-shrink-0 {styles.icon}">
						<svelte:component this={IconComponent} class="w-5 h-5" />
					</div>

					<!-- Content -->
					<div class="flex-1 min-w-0">
						<div class="flex items-start justify-between">
							<div class="flex-1">
								<h4 class="text-sm font-medium leading-5">
									{notification.title}
								</h4>
								{#if notification.message}
									<p class="mt-1 text-sm opacity-90 leading-5">
										{notification.message}
									</p>
								{/if}
							</div>

							<!-- Close Button -->
							<button
								onclick={() => removeNotification(notification.id)}
								class="ml-2 flex-shrink-0 p-1 rounded-md transition-colors {styles.button}"
								title="Dismiss"
							>
								<X class="w-4 h-4" />
							</button>
						</div>

						<!-- Actions -->
						{#if notification.actions && notification.actions.length > 0}
							<div class="mt-3 flex space-x-2">
								{#each notification.actions as action, index (index)}
									<button
										onclick={() => handleAction(action.action, notification.id)}
										class="text-xs font-medium px-3 py-1 rounded-md transition-colors bg-white/10 hover:bg-white/20"
									>
										{action.label}
									</button>
								{/each}
							</div>
						{/if}
					</div>
				</div>
			</div>

			<!-- Progress Bar for Auto-dismiss -->
			{#if !notification.persistent && notification.duration}
				<div class="h-1 bg-black/20 relative overflow-hidden">
					<div
						class="h-full bg-white/30 absolute left-0 top-0 animate-[shrink_{notification.duration}ms_linear_forwards]"
						style="width: 100%;"
					></div>
				</div>
			{/if}
		</div>
	{/each}
</div>

<style>
	@keyframes shrink {
		from {
			width: 100%;
		}
		to {
			width: 0%;
		}
	}
</style>
