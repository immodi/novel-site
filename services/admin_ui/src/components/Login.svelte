<script lang="ts">
    import { login } from "../api/login";
    import { setUserData } from "../lib/states/auth_state.svelte";

    let email = "";
    let password = "";
    let isLoading = false;
    let error = "";

    async function handleSubmit() {
        isLoading = true;
        error = "";
        const { data, error: networkError } = await login({ email, password });

        if (networkError) error = networkError;
        else if (data?.error) error = data.error;
        else if (data?.token) setUserData(data);

        isLoading = false;
    }

    function clearError() {
        error = "";
    }
</script>

<div
    class="h-screen flex items-center justify-center bg-gradient-to-br from-[#19183B] to-[#708993] p-4 overflow-hidden"
>
    <div
        class="max-w-md w-full bg-[#E7F2EF] p-6 rounded-xl shadow-lg mx-auto my-auto"
    >
        <div class="text-center mb-6">
            <h2 class="text-2xl font-bold text-[#19183B] mb-2">
                Sign in to your account
            </h2>
            <p class="text-sm text-[#708993]">
                Enter your credentials to continue
            </p>
        </div>

        <!-- Error Display -->
        {#if error}
            <div class="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg">
                <div class="flex items-center">
                    <svg
                        class="w-5 h-5 text-red-500 mr-2"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                    >
                        <path
                            fill-rule="evenodd"
                            d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                            clip-rule="evenodd"
                        />
                    </svg>
                    <span class="text-red-700 text-sm font-medium">{error}</span
                    >
                </div>
            </div>
        {/if}

        <form on:submit|preventDefault={handleSubmit} class="space-y-4">
            <div>
                <label
                    for="email"
                    class="block text-sm font-medium text-[#19183B] mb-1"
                    >Email</label
                >
                <input
                    id="email"
                    type="email"
                    required
                    bind:value={email}
                    on:input={clearError}
                    class={`w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[#19183B] ${!error ? "border-[#A1C2BD]" : "border-red-300"} `}
                    placeholder="your@email.com"
                />
            </div>

            <div>
                <label
                    for="password"
                    class="block text-sm font-medium text-[#19183B] mb-1"
                    >Password</label
                >
                <input
                    id="password"
                    type="password"
                    required
                    bind:value={password}
                    on:input={clearError}
                    class={`w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[#19183B] ${!error ? "border-[#A1C2BD]" : "border-red-300"}`}
                    placeholder="Enter your password"
                />
            </div>

            <button
                type="submit"
                disabled={isLoading}
                class="cursor-pointer w-full py-2 mt-4 bg-[#19183B] text-white rounded-lg hover:bg-[#2a2852] disabled:opacity-50 transition-colors"
            >
                {#if isLoading}
                    Signing in...
                {:else}
                    Sign in
                {/if}
            </button>
        </form>
    </div>
</div>
