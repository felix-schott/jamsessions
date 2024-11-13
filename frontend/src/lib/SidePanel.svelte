<script lang="ts">
	import type { Snippet } from 'svelte';
	import WindowControlButton from './WindowControlButton.svelte';

	interface Props {
		background?: boolean;
		children: Snippet;
		hide: () => void;
	}

	let { background = false, children, hide }: Props = $props();
</script>

<div class="side-panel" class:background class:padding={!background}>
	{#if !background}
		<WindowControlButton type="close" onclick={hide} />
		<div>
			{@render children()}
		</div>
	{/if}
</div>

<style>
	.side-panel {
		background: white;
		width: auto;
		position: relative;
		min-width: 0;
		max-height: 100%;
	}

	.padding {
		padding: 1em 3em;
	}

	.side-panel > div {
		overflow: auto;
		height: 100%;
	}

	@media (max-width: 480px) {
		.padding {
			padding: 1em;
		}
	}

	.background {
		width: 0;
	}
</style>
