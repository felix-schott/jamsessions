<script lang="ts">
    export let title: string;
    export let height: string = "1em";
    export let width: string = "1em";
    export let style: string = "";
    export let color: string = "black";
    let _class: string = "";
    export { _class as class };
</script>

<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" {style} {height} {width} class={_class} fill={color}><title>{title}</title><path d="M11 11V5H13V11H19V13H13V19H11V13H5V11H11Z"></path></svg>