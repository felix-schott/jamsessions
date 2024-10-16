<script lang="ts">
    export let height: string = "1em";
    export let width: string = "1em";
    export let title: string;
    let _class: string = "";
    export { _class as class };
</script>

<svg xmlns="http://www.w3.org/2000/svg" {height} {width} class={_class} viewBox="0 0 512 512"><title>{title}</title><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="48" d="M112 184l144 144 144-144"/></svg>
<!-- <svg xmlns="http://www.w3.org/2000/svg" {height} {width} viewBox="0 0 24 24" fill="currentColor"><title>{title}</title><path d="M6.41421 5H10V3H3V10H5V6.41421L9.29289 10.7071L10.7071 9.29289L6.41421 5ZM21 14H19V17.5858L14.7071 13.2929L13.2929 14.7071L17.5858 19H14V21H21V14Z"></path></svg> -->