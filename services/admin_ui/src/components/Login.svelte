<script>
    let email = "";
    let password = "";
    let rememberMe = false;
    let isLoading = false;

    async function handleSubmit() {
        isLoading = true;
        await new Promise((resolve) => setTimeout(resolve, 1000));
        console.log("Logging in with:", { email, password, rememberMe });
        isLoading = false;
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
                    class="w-full px-3 py-2 border border-[#A1C2BD] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#19183B]"
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
                    class="w-full px-3 py-2 border border-[#A1C2BD] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#19183B]"
                    placeholder="Enter your password"
                />
            </div>

            <button
                type="submit"
                disabled={isLoading}
                class="cursor-pointer w-full py-2 mt-4 bg-[#19183B] text-white rounded-lg hover:bg-[#2a2852] disabled:opacity-50"
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
