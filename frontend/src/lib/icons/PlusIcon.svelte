<script lang="ts">
    interface Props {
        title: string;
        height?: string;
        width?: string;
        style?: string;
        color?: string;
        class?: string;
    }

    let {
        title,
        height = "1em",
        width = "1em",
        style = "",
        color = "black",
        class: _class = ""
    }: Props = $props();
    
</script>

<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" {style} {height} {width} class={_class} fill={color}><title>{title}</title><path d="M11 11V5H13V11H19V13H13V19H11V13H5V11H11Z"></path></svg>