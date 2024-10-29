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
<!-- https://remixicon.com/icon/pencil-fill -->
<svg xmlns="http://www.w3.org/2000/svg" {onclick} {style} {height} {width} class={_class} viewBox="0 0 24 24" fill="currentColor"><title>{title}</title><path d="M12.8995 6.85453L17.1421 11.0972L7.24264 20.9967H3V16.754L12.8995 6.85453ZM14.3137 5.44032L16.435 3.319C16.8256 2.92848 17.4587 2.92848 17.8492 3.319L20.6777 6.14743C21.0682 6.53795 21.0682 7.17112 20.6777 7.56164L18.5563 9.68296L14.3137 5.44032Z"></path></svg>