<!--  
    @component 
    This component represents a generic close button in the top-right corner of its container (must have "position: relative").

    #### Props
    - type: string - either "close", "minimise" or "maximise"
    - style: string - inline style
-->

<script lang="ts">
    export let style: string = "";
    export let type: string = "";
    let _additional_class: string = "";
	export { _additional_class as class };
    let _class = `window-control-button ${_additional_class}`;
    let content: string = "";

    switch (type) {
        case "close":
            content = '×'
            break
        case "minimise":
            content =  '-'
            break
        case "maximise":
            content = '+'
            break
    }
</script>

<button class={_class} {style} on:click >
{content}
<slot/>
</button>

<style>
    button {
        background: none;
    }
    .window-control-button {
        cursor: pointer;
        display: inline-block;
        padding: 0.5em;
        position: absolute;
        top: 0.7em;
        right: 0.7em;
        line-height: 1ex;
        font-weight: bolder;
        z-index: 800000000;
    }
</style>