<!-- Changes <select multiple> behaviour to also accept a single click instead of the default shift+click -->
<script lang="ts">
    interface Props {
        name: string;
        id: string;
        title: string;
        children?: import('svelte').Snippet;
        onchange?: () => any;
    }

    let {
        name,
        id,
        title,
        children,
        onchange
    }: Props = $props();

    // adapted from https://stackoverflow.com/a/59084958
    const onMousedown = (e: MouseEvent) => {
        // native mobile select already has checkboxes on both ios and android
        if (
            navigator.userAgent.match(/Android/i) ||
            navigator.userAgent.match(/iPhone|iPad|iPod/i)
        )
            return;
        if (e.shiftKey) return;
        e.preventDefault();
        (e.target! as HTMLOptionElement).focus();
        var scroll = (e.target! as HTMLOptionElement).scrollTop;
        (e.target! as HTMLOptionElement).selected = !(
            e.target! as HTMLOptionElement
        ).selected;
        (e.target! as HTMLOptionElement).scrollTop = scroll;
        // make sure that other listeners pick up on the change
        // even though we prevented the default
        if (onchange) onchange()
    };
</script>

<select {id} {title} {name} {onchange} multiple onmousedown={onMousedown}>
    {@render children?.()}
</select>
