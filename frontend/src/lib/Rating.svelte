<script lang="ts">
	import StarEmptyIcon from './icons/StarEmptyIcon.svelte';
	import StarHalfIcon from './icons/StarHalfIcon.svelte';
	import StarFilledIcon from './icons/StarFilledIcon.svelte';

	export let n: number;
    export let style: string = "";
    export let size: string = "1em";

    console.log("n", n)

	if (n < 0 || n > 5) {
		throw Error('Rating must be between 0 and 5');
	}
</script>

<div class="rating" {style}>
	{#each Array.from(Array(Math.floor(n)).keys()) as _}
		<StarFilledIcon height={size} width={size}/>
	{/each}
	{#if n % 1 !== 0}
		<StarHalfIcon height={size} width={size} />
	{/if}
	{#each Array.from(Array(Math.floor(5 - n)).keys()) as _}
		<StarEmptyIcon height={size} width={size}/>
	{/each}
</div>

<style>
    .rating {
        display: flex;
    }
</style>