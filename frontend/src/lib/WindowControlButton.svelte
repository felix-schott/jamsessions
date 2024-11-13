<!--  
    @component 
    This component represents a generic close button in the top-right corner of its container (must have "position: relative").

    #### Props
    - type: string - either "close", "minimise" or "maximise"
    - style: string - inline style
-->

<script lang="ts">
	interface Props {
		style?: string;
		type?: string;
		class?: string;
		onclick: () => void;
	}
	let content: string = $state('');

	let { style = '', type = '', class: _class = '', onclick }: Props = $props();
	let allClasses = `window-control-button ${_class}`;

	switch (type) {
		case 'close':
			content = '×';
			break;
		case 'minimise':
			content = '-';
			break;
		case 'maximise':
			content = '+';
			break;
	}
</script>

<button class={allClasses} {style} {onclick}>
	{content}
</button>

<style>
	button {
		background: none;
	}
	.window-control-button {
		cursor: pointer;
		display: inline-block;
		padding: 0.5em;
		position: absolute;
		top: 0.7em;
		right: 0.7em;
		line-height: 1ex;
		font-weight: bolder;
		z-index: 800000000;
	}

	@media (max-width: 480px) {
		.window-control-button {
			top: 0.2em;
			right: 0.2em;
		}
	}
</style>
