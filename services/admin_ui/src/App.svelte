<script lang="ts">
    import Dashboard from "./components/Dashboard.svelte";
    import Login from "./components/Login.svelte";
    import {
        getUserToken,
        type UserToken,
    } from "./lib/states/auth_state.svelte";
    const token: UserToken = $derived(getUserToken());
</script>

<main class="h-screen w-screen">
    {#if token}
        <Dashboard />
    {:else}
        <Login />
    {/if}
</main>
