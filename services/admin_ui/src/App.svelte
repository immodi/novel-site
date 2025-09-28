<script lang="ts">
    import Login from "./components/Login.svelte";
    import { getUserToken, type UserToken } from "./lib/states/auth.svelte";
    let token: UserToken = $derived(getUserToken());
</script>

<main class="h-screen w-screen">
    {#if token}
        <div>authed user</div>
    {:else}
        <Login />
    {/if}
</main>
