<script lang="ts">
    export let title: string;
    export let height: string = "1em";
    export let width: string = "1em";
    export let style: string = "";
    export let colour: string = "black";
    let _class: string = "";
    export { _class as class };
</script>

<svg xmlns="http://www.w3.org/2000/svg" on:click {style} {height} {width} class={_class} viewBox="0 0 512 512"><title>{title}</title><path fill={colour} d="M256 56C145.72 56 56 145.72 56 256s89.72 200 200 200 200-89.72 200-200S366.28 56 256 56zm0 82a26 26 0 11-26 26 26 26 0 0126-26zm48 226h-88a16 16 0 010-32h28v-88h-16a16 16 0 010-32h32a16 16 0 0116 16v104h28a16 16 0 010 32z"/></svg>
<!-- <svg xmlns="http://www.w3.org/2000/svg" on:click {style} {height} {width} class="ionicon" viewBox="0 0 512 512"><title>{title}</title><path d="M248 64C146.39 64 64 146.39 64 248s82.39 184 184 184 184-82.39 184-184S349.61 64 248 64z" fill="none" stroke="currentColor" stroke-miterlimit="10" stroke-width="32"/><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32" d="M220 220h32v116"/><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-miterlimit="10" stroke-width="32" d="M208 340h88"/><path d="M248 130a26 26 0 1026 26 26 26 0 00-26-26z"/></svg> -->