<script lang="ts">
	import WindowControlButton from './WindowControlButton.svelte';

	export let isVisible: () => boolean;
	export let hide: () => void;

	// const dispatch = createEventDispatcher();

	const checkClick = (ev: Event) => {
		let clickedElement = ev.target as HTMLElement;
		if (clickedElement.id === 'outside-popup') {
			// dispatch('close');
			hide();
		}
	};
</script>

{#if isVisible()}
	<div class="blurry-background" />
	<div class="popup-pane" id="outside-popup" on:click={checkClick}>
		<div class="popup">
			<WindowControlButton type="close" on:click={hide} />
			<div class="top-left">
				<slot name="top-left" />
			</div>
			<div class="bottom-left">
				<slot name="bottom-left" />
			</div>
			<div class="modal-content">
				<slot />
			</div>
		</div>
	</div>
{/if}

<!-- on:click={() => dispatch('close')} /> -->

<style>
	.blurry-background {
		z-index: 800000001;
		/*
        z-index: 99999998;
        */
		position: absolute;
		top: 0;
		left: 0;
		background-color: rgba(128, 128, 128, 0.8);
		display: flex;
		justify-content: center;
		align-items: center;
		height: 100%;
		width: 100%;
	}

	.top-left {
		position: absolute;
		top: 0.7em;
		left: 0.7em;
		z-index: 800000000;
	}

	@media (max-width: 480px) {
		.top-left {
			top: 0.2em;
			left: 0.2em;
		}
	}

	.bottom-left {
		position: absolute;
		bottom: 0.7em;
		left: 0.7em;
		z-index: 800000000;
	}

	@media (max-width: 480px) {
		.bottom-left {
			bottom: 0.2em;
			left: 0.2em;
		}
	}


	.popup-pane {
		z-index: 800000002;
		/*
        z-index: 99999999;
        */
		position: absolute;
		top: 0;
		left: 0;
		height: 100%;
		width: 100%;
		display: flex;
		justify-content: center;
		align-items: center;
	}

	.popup {
		background-color: white;
		border-radius: 10px;
		padding: 2em;
		box-shadow: 2px 2px 2px grey;
		position: relative;
	}

	@media (max-width: 480px) {
		.popup {
			padding: 1em;
		}
	}

	.modal-content {
		overflow: auto;
		max-width: 70vw;
		max-height: 70vh;
	}

	@media (max-width: 480px) {
		.modal-content {
			max-width: 80vw;
			max-height: 80vh;
			font-size: smaller;
		}
	}
</style>
