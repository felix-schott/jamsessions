<script lang="ts">
    interface Props {
        title: string;
        height?: string;
        width?: string;
        style?: string;
        class?: string;
        onclick?: () => any;
    }

    let {
        title,
        height = "1em",
        width = "1em",
        style = "",
        class: _class = "",
        onclick
    }: Props = $props();
    
</script>

<svg xmlns="http://www.w3.org/2000/svg" {onclick} {style} {height} {width} class={_class} viewBox="0 0 512 512"><title>{title}</title><path d="M384 336a63.78 63.78 0 00-46.12 19.7l-148-83.27a63.85 63.85 0 000-32.86l148-83.27a63.8 63.8 0 10-15.73-27.87l-148 83.27a64 64 0 100 88.6l148 83.27A64 64 0 10384 336z"/></svg>