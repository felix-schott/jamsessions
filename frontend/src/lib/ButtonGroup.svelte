<script lang="ts">
	interface Props {
		options: string[];
		activeIndex?: number;
		onchange?: (activeOption: string) => any;
	}

	let { options, activeIndex = $bindable(0), onchange }: Props = $props();
</script>

<div class="btn-group">
	{#each options as opt, idx}
		<button
			class:active={activeIndex === idx}
			onclick={() => {
				activeIndex = idx;
				if (onchange) onchange(opt);
			}}>{opt}</button
		>
	{/each}
</div>

<style>
	.btn-group button {
		background-color: lightgrey;
		padding: 10px 24px; /* Some padding */
		cursor: pointer; /* Pointer/hand icon */
		float: left; /* Float the buttons side by side */
		border-color: white;
	}

	/* .btn-group button.active {
		box-shadow: inset grey 0px 0px 10px -2px; 
        background-color: darkgrey;
	} */

	.btn-group button.active {
		box-shadow: 	inset black 0px 0px 10px -2px;
		background: darkgrey;
		/* border-color: grey; */
		color: white;

	}

	.btn-group button:not(:last-child) {
		border-right: none; /* Prevent double borders */
		border-top-right-radius: 0;
		border-bottom-right-radius: 0;
	}

	.btn-group button:last-child {
		border-top-left-radius: 0;
		border-bottom-left-radius: 0;
	}

	/* Clear floats (clearfix hack) */
	.btn-group:after {
		content: '';
		clear: both;
		display: table;
	}

	/* Add a background color on hover */
	.btn-group button:hover {
		background-color: grey;
	}
</style>
