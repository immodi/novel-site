<script lang="ts">
    import Dashboard from "./components/Dashboard.svelte";
    import Login from "./components/Login.svelte";
    import { getUserToken } from "./lib/states/auth_state.svelte";

    const token = $derived(getUserToken());
</script>

<main class="h-screen w-screen">
    {#if token}
        <Dashboard />
    {:else}
        <Login />
    {/if}
</main>
