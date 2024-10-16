<script lang="ts">
    export let height: string = "1em";
    export let width: string = "1em";
    export let title: string;
    let _class: string = "";
    export { _class as class };
</script>

<svg xmlns="http://www.w3.org/2000/svg" {height} {width} class={_class} viewBox="0 0 512 512"><title>{title}</title><circle cx="256" cy="192" r="32"/><path d="M256 32c-88.22 0-160 68.65-160 153 0 40.17 18.31 93.59 54.42 158.78 29 52.34 62.55 99.67 80 123.22a31.75 31.75 0 0051.22 0c17.42-23.55 51-70.88 80-123.22C397.69 278.61 416 225.19 416 185c0-84.35-71.78-153-160-153zm0 224a64 64 0 1164-64 64.07 64.07 0 01-64 64z"/></svg>
<!-- <svg xmlns="http://www.w3.org/2000/svg" {height} {width} class="ionicon" viewBox="0 0 512 512"><path d="M256 48c-79.5 0-144 61.39-144 137 0 87 96 224.87 131.25 272.49a15.77 15.77 0 0025.5 0C304 409.89 400 272.07 400 185c0-75.61-64.5-137-144-137z" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32"/><circle cx="256" cy="192" r="48" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32"/></svg> -->