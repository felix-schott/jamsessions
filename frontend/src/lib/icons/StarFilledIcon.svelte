<script lang="ts">
    interface Props {
        height?: string;
        width?: string;
        title?: string;
        onmouseover?: () => any;
        onclick?: () => any;
    }

    let { height = "1em", width = "1em", title = "", onmouseover, onclick }: Props = $props();
</script>

<svg xmlns="http://www.w3.org/2000/svg" {onmouseover} {onclick} {height} {width} viewBox="0 0 24 24" fill="currentColor">{title}<path d="M12.0006 18.26L4.94715 22.2082L6.52248 14.2799L0.587891 8.7918L8.61493 7.84006L12.0006 0.5L15.3862 7.84006L23.4132 8.7918L17.4787 14.2799L19.054 22.2082L12.0006 18.26Z"></path></svg>