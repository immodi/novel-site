<script lang="ts">
    import Dashboard from "./components/Dashboard.svelte";
    import Login from "./components/Login.svelte";
    import { getUserData, type UserData } from "./lib/states/auth_state.svelte";

    const { token, coverImage, username }: UserData = $derived(
        getUserData() || { token: "", coverImage: "", username: "" },
    );
</script>

<main class="h-screen w-screen">
    {#if token}
        <Dashboard {username} {coverImage} />
    {:else}
        <Login />
    {/if}
</main>
