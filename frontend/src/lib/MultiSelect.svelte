<!-- Changes <select multiple> behaviour to also accept a single click instead of the default shift+click -->
<script lang="ts">
    import { createEventDispatcher } from "svelte";
    export let name: string;
    export let id: string;
    export let title: string;

    let dispatch = createEventDispatcher();

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
        dispatch("change")
    };
</script>

<select {id} {title} {name} on:change multiple on:mousedown={onMousedown}>
    <slot />
</select>
