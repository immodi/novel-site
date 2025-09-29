<script lang="ts">
    import { mapDbUsersToUsersDTO } from "../../lib/mappers/users_mapper";
    import { usersState } from "../../lib/states/users_state.svelte";
    import refreshIcon from "../../assets/refresh_icon.svg";

    const users = $derived(mapDbUsersToUsersDTO(usersState.data ?? []));

    function refresh(): void {
        usersState.refresh();
    }

    $effect(() => {
        if (!usersState.data && !usersState.loading) {
            refresh();
        }
    });
</script>

<div class="p-4 sm:p-6">
    <!-- Header -->
    <div
        class="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4 mb-6"
    >
        <h2 class="text-xl font-bold text-[#19183B] text-center sm:text-left">
            User Management
        </h2>
        <button
            class="flex items-center justify-center gap-2 cursor-pointer px-4 py-2 bg-[#19183B] text-white rounded-lg hover:bg-[#2a2852] transition-colors disabled:opacity-50 disabled:cursor-not-allowed w-full sm:w-auto"
            onclick={refresh}
            disabled={usersState.loading}
        >
            {#if usersState.loading}
                <!-- Loading Spinner -->
                <svg
                    class="animate-spin h-5 w-5 text-white"
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                >
                    <circle
                        class="opacity-25"
                        cx="12"
                        cy="12"
                        r="10"
                        stroke="currentColor"
                        stroke-width="4"
                    ></circle>
                    <path
                        class="opacity-75"
                        fill="currentColor"
                        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                    ></path>
                </svg>
                <span class="sm:hidden">Refreshing...</span>
            {:else}
                <img src={refreshIcon} alt="Refresh" class="w-5 h-5" />
                <span class="sm:hidden">Refresh Users</span>
            {/if}
        </button>
    </div>

    <!-- Error Display -->
    {#if usersState.error}
        <div class="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg">
            <div class="flex items-start">
                <svg
                    class="w-5 h-5 text-red-500 mr-2 mt-0.5 flex-shrink-0"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                >
                    <path
                        fill-rule="evenodd"
                        d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                        clip-rule="evenodd"
                    />
                </svg>
                <span class="text-red-700 text-sm font-medium break-words">
                    {usersState.error}
                </span>
            </div>
        </div>
    {/if}

    <!-- Mobile Card View -->
    <div class="block sm:hidden space-y-4">
        {#if usersState.loading}
            <!-- Loading State -->
            <div class="text-center py-8">
                <div class="flex justify-center items-center">
                    <svg
                        class="animate-spin h-8 w-8 text-[#19183B]"
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                    >
                        <circle
                            class="opacity-25"
                            cx="12"
                            cy="12"
                            r="10"
                            stroke="currentColor"
                            stroke-width="4"
                        ></circle>
                        <path
                            class="opacity-75"
                            fill="currentColor"
                            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                        ></path>
                    </svg>
                    <span class="ml-2 text-[#19183B]">Loading users...</span>
                </div>
            </div>
        {:else if users.length === 0}
            <!-- Empty State -->
            <div class="text-center py-8">
                <div class="text-[#708993]">
                    <svg
                        class="mx-auto h-12 w-12 text-gray-400"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z"
                        />
                    </svg>
                    <h3 class="mt-2 text-sm font-medium text-gray-900">
                        No users
                    </h3>
                    <p class="mt-1 text-sm text-gray-500">
                        Get started by creating a new user.
                    </p>
                </div>
            </div>
        {:else}
            <!-- Users Cards -->
            {#each users as user}
                <div
                    class="bg-white rounded-lg border border-gray-200 p-4 shadow-sm"
                >
                    <!-- User Header -->
                    <div class="flex items-start space-x-3 mb-3">
                        <img
                            class="h-12 w-12 rounded-full object-cover flex-shrink-0"
                            src={user.coverImage}
                            alt="User avatar"
                        />
                        <div class="min-w-0 flex-1">
                            <div class="flex items-center justify-between">
                                <div>
                                    <div
                                        class="text-base font-semibold text-[#19183B] truncate"
                                    >
                                        {user.username}
                                    </div>
                                    <div class="text-sm text-[#708993]">
                                        ID: {user.id}
                                    </div>
                                </div>
                                <span
                                    class={`px-2 py-1 text-xs font-semibold rounded-full ${user.role === "admin" ? "bg-purple-100 text-purple-800" : "bg-green-100 text-green-800"}`}
                                >
                                    {user.role}
                                </span>
                            </div>
                        </div>
                    </div>

                    <!-- User Details -->
                    <div class="space-y-2 text-sm">
                        <div class="flex justify-between">
                            <span class="text-[#708993] font-medium"
                                >Email:</span
                            >
                            <span class="text-[#19183B] text-right break-all"
                                >{user.email}</span
                            >
                        </div>
                        <div class="flex justify-between">
                            <span class="text-[#708993] font-medium"
                                >Created:</span
                            >
                            <span class="text-[#19183B]">{user.createdAt}</span>
                        </div>
                    </div>

                    <!-- Actions -->
                    <div class="mt-4 pt-3 border-t border-gray-100">
                        <div class="text-center">
                            <span class="text-sm text-[#708993] select-none">
                                No Actions Available
                            </span>
                        </div>
                    </div>
                </div>
            {/each}
        {/if}
    </div>

    <!-- Desktop Table View -->
    <div class="hidden sm:block overflow-x-auto w-full">
        <table class="w-full divide-y divide-gray-200">
            <thead>
                <tr class="bg-gray-50">
                    <th
                        class="px-4 py-3 text-left text-xs font-medium text-[#19183B] uppercase tracking-wider"
                    >
                        User
                    </th>
                    <th
                        class="px-4 py-3 text-left text-xs font-medium text-[#19183B] uppercase tracking-wider"
                    >
                        Email
                    </th>
                    <th
                        class="px-4 py-3 text-left text-xs font-medium text-[#19183B] uppercase tracking-wider"
                    >
                        Role
                    </th>
                    <th
                        class="px-4 py-3 text-left text-xs font-medium text-[#19183B] uppercase tracking-wider"
                    >
                        Created At
                    </th>
                    <th
                        class="px-4 py-3 text-left text-xs font-medium text-[#19183B] uppercase tracking-wider"
                    >
                        Actions
                    </th>
                </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
                {#if usersState.loading}
                    <!-- Loading State -->
                    <tr>
                        <td colspan="5" class="px-4 py-8 text-center">
                            <div class="flex justify-center items-center">
                                <svg
                                    class="animate-spin h-8 w-8 text-[#19183B]"
                                    xmlns="http://www.w3.org/2000/svg"
                                    fill="none"
                                    viewBox="0 0 24 24"
                                >
                                    <circle
                                        class="opacity-25"
                                        cx="12"
                                        cy="12"
                                        r="10"
                                        stroke="currentColor"
                                        stroke-width="4"
                                    ></circle>
                                    <path
                                        class="opacity-75"
                                        fill="currentColor"
                                        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                                    ></path>
                                </svg>
                                <span class="ml-2 text-[#19183B]"
                                    >Loading users...</span
                                >
                            </div>
                        </td>
                    </tr>
                {:else if users.length === 0}
                    <!-- Empty State -->
                    <tr>
                        <td colspan="5" class="px-4 py-8 text-center">
                            <div class="text-[#708993]">
                                <svg
                                    class="mx-auto h-12 w-12 text-gray-400"
                                    fill="none"
                                    viewBox="0 0 24 24"
                                    stroke="currentColor"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z"
                                    />
                                </svg>
                                <h3
                                    class="mt-2 text-sm font-medium text-gray-900"
                                >
                                    No users
                                </h3>
                                <p class="mt-1 text-sm text-gray-500">
                                    Get started by creating a new user.
                                </p>
                            </div>
                        </td>
                    </tr>
                {:else}
                    <!-- Users List -->
                    {#each users as user}
                        <tr class="hover:bg-gray-50 transition-colors">
                            <td class="px-4 py-4 whitespace-nowrap">
                                <div class="flex items-center">
                                    <div class="flex-shrink-0 h-10 w-10">
                                        <img
                                            class="h-10 w-10 rounded-full object-cover"
                                            src={user.coverImage}
                                            alt="User avatar"
                                        />
                                    </div>
                                    <div class="ml-3">
                                        <div
                                            class="text-sm font-medium text-[#19183B]"
                                        >
                                            {user.username}
                                        </div>
                                        <div class="text-sm text-[#708993]">
                                            ID: {user.id}
                                        </div>
                                    </div>
                                </div>
                            </td>
                            <td class="px-4 py-4 whitespace-nowrap">
                                <div
                                    class="text-sm text-[#19183B] max-w-xs truncate"
                                >
                                    {user.email}
                                </div>
                            </td>
                            <td class="px-4 py-4 whitespace-nowrap">
                                <span
                                    class={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${user.role === "admin" ? "bg-purple-100 text-purple-800" : "bg-green-100 text-green-800"}`}
                                >
                                    {user.role}
                                </span>
                            </td>
                            <td
                                class="px-4 py-4 whitespace-nowrap text-sm text-[#19183B]"
                            >
                                {user.createdAt}
                            </td>
                            <td
                                class="px-4 py-4 whitespace-nowrap text-sm font-medium"
                            >
                                <div class="flex space-x-2">
                                    <span
                                        class="text-sm text-[#19183B] select-none"
                                        >No Actions Available</span
                                    >
                                </div>
                            </td>
                        </tr>
                    {/each}
                {/if}
            </tbody>
        </table>
    </div>
</div>
