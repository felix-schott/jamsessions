<script lang="ts">
	import type { Snippet } from 'svelte';
	import ExpandIcon from './icons/ExpandIcon.svelte';

	interface Props {
		id: string;
		title?: string;
		blur?: boolean;
		onclick?: () => any;
		onclose?: () => void;
		heading?: Snippet;
		content: Snippet;
	}
	let { id, title = '', blur = false, onclick, onclose, heading, content }: Props = $props();
</script>

<div class="popup-container">
	<div {id} class:popup={!blur} class:popup-blur={blur} {title}>
		<span {onclick}>
			{#if heading}
				{@render heading()}
			{/if}
			{@render content()}
		</span>
		<!-- <div class="popup-controls"> -->
		<div
			id={id + 'close'}
			onclick={() => {
				if (onclose) onclose();
			}}
			class="popup-close"
		>
			×
		</div>
		{#if blur}
			<div class="expand"><ExpandIcon title="Expand" /></div>
		{/if}
	</div>
</div>

<style>
	.popup-container {
		max-width: 40vw;
	}

	.popup,
	.popup-blur {
		background-color: white;
		border-radius: 10px;
		padding: 2em;
		display: inline-flex;
		font-size: smaller;
		cursor: pointer;
		overflow: hidden;
	}

	.popup-blur {
		padding-bottom: 3em;
	}

	.popup-blur:after {
		content: '';
		border-radius: 10px;
		position: absolute;
		left: 0px;
		right: 0px;
		height: 90%;
		bottom: 0px;
		background: linear-gradient(180deg, rgba(0, 0, 0, 0) 0%, rgba(255, 255, 255, 1) 75%);
		pointer-events: none;
	}

	@media (max-width: 480px) {
		.popup,
		.popup-blur {
			padding: 1em;
		}
		.popup-container {
			max-width: 60vw;
		}
	}

	.popup-close {
		position: absolute;
		top: 0;
		right: 0;
		cursor: pointer;
		color: red;
		padding: 0.5em;
	}

	.popup > span,
	.popup-blur > span {
		font-size: larger;
		/* margin-bottom: 1.5em; */
	}

	.expand {
		position: absolute;
		bottom: 0;
		right: 0;
		left: 0;
		z-index: 99999999999;
		display: flex;
		justify-content: center;
		padding: 0.7em;
		box-shadow: rgba(255, 255, 255, 1);
		pointer-events: none;
	}
</style>
