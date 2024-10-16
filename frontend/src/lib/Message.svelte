<!-- Displays a message in the middle of the screen -->
<script lang="ts">
    import { onMount } from "svelte";

    let component: HTMLElement;
    export let message: string;
    
    onMount(async () => { 
        await new Promise(r => setTimeout(r, 3000)); // ms value must correspond to css animation duration (2s)
        component.parentNode!.removeChild(component); // destroy
    })
</script>

<div class="message-pane" bind:this={component}>
    <div class="message">
        {message}
    </div>
</div>

<style>
    .message-pane {
        z-index: 99999999;
        height: 100vh;
        width: 100vw;
        background: transparent;
        display: flex;
        justify-content: center;
        align-items: center;
        position: absolute;
        top: 0;
        right: 0;
    }

    .message {
        vertical-align: middle;
        animation: fade 3s linear;
        padding: 1em;
        background: rgba(0, 0, 0, 0.5);
        color: whitesmoke;
        border-radius: 10px;
        cursor: default;
        font-size: x-large;
    }

    @keyframes fade {
        0%, 100% {
            background: rgba(0, 0, 0, 0);
        }
        50% {
            background: rgba(0, 0, 0, 0.5);
        }
    }
</style>
