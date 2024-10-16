<script lang="ts">
    export let height: string = "1em";
    export let width: string = "1em";
    export let title: string;
    let _class: string = "";
    export { _class as class };
</script>

<svg xmlns="http://www.w3.org/2000/svg" {height} {width} class={_class} viewBox="0 0 512 512"><title>{title}</title><path d="M256 48C141.13 48 48 141.13 48 256s93.13 208 208 208 208-93.13 208-208S370.87 48 256 48zm96 240h-96a16 16 0 01-16-16V128a16 16 0 0132 0v128h80a16 16 0 010 32z"/></svg>
<!-- <svg xmlns="http://www.w3.org/2000/svg" {width} {height} class="ionicon" viewBox="0 0 512 512"><path d="M256 64C150 64 64 150 64 256s86 192 192 192 192-86 192-192S362 64 256 64z" fill="none" stroke="currentColor" stroke-miterlimit="10" stroke-width="32"/><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32" d="M256 128v144h96"/></svg> -->