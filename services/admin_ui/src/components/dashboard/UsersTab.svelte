<script lang="ts">
    import type { User } from "../../types/dtos/user";

    export let users: User[] = [];

    function editUser(userId: number): void {
        console.log(`Edit user with ID: ${userId}`);
        // Future implementation will open a form to edit the user
    }

    function deleteUser(userId: number): void {
        console.log(`Delete user with ID: ${userId}`);
        // Future implementation will delete the user
    }
</script>

<div class="p-6">
    <div class="flex justify-between items-center mb-6">
        <h2 class="text-xl font-bold text-[#19183B]">User Management</h2>
        <button
            class="cursor-pointer px-4 py-2 bg-[#19183B] text-white rounded-lg hover:bg-[#2a2852] transition-colors"
        >
            Add New User
        </button>
    </div>

    <!-- Users Table -->
    <div class="overflow-x-auto w-full">
        <table class="w-full divide-y divide-gray-200">
            <thead>
                <tr class="bg-gray-50">
                    <th
                        class="px-6 py-3 text-left text-xs font-medium text-[#19183B] uppercase tracking-wider"
                    >
                        User
                    </th>
                    <th
                        class="px-6 py-3 text-left text-xs font-medium text-[#19183B] uppercase tracking-wider"
                    >
                        Email
                    </th>
                    <th
                        class="px-6 py-3 text-left text-xs font-medium text-[#19183B] uppercase tracking-wider"
                    >
                        Role
                    </th>
                    <th
                        class="px-6 py-3 text-left text-xs font-medium text-[#19183B] uppercase tracking-wider"
                    >
                        Created At
                    </th>
                    <th
                        class="px-6 py-3 text-left text-xs font-medium text-[#19183B] uppercase tracking-wider"
                    >
                        Actions
                    </th>
                </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
                {#each users as user}
                    <tr class="hover:bg-gray-50 transition-colors">
                        <td class="px-6 py-4 whitespace-nowrap">
                            <div class="flex items-center">
                                <div class="flex-shrink-0 h-10 w-10">
                                    <img
                                        class="h-10 w-10 rounded-full object-cover"
                                        src={user.coverImage}
                                        alt="User avatar"
                                    />
                                </div>
                                <div class="ml-4">
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
                        <td class="px-6 py-4 whitespace-nowrap">
                            <div class="text-sm text-[#19183B]">
                                {user.email}
                            </div>
                        </td>
                        <td class="px-6 py-4 whitespace-nowrap">
                            <span
                                class={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${user.role === "admin" ? "bg-purple-100 text-purple-800" : "bg-green-100 text-green-800"}`}
                            >
                                {user.role}
                            </span>
                        </td>
                        <td
                            class="px-6 py-4 whitespace-nowrap text-sm text-[#19183B]"
                        >
                            {user.createdAt}
                        </td>
                        <td
                            class="px-6 py-4 whitespace-nowrap text-sm font-medium"
                        >
                            <div class="flex space-x-2">
                                <button
                                    class="cursor-pointer px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors flex items-center w-20 justify-center"
                                    on:click={() => editUser(user.id)}
                                >
                                    Edit
                                </button>
                                <button
                                    class="cursor-pointer px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors flex items-center w-20 justify-center"
                                    on:click={() => deleteUser(user.id)}
                                >
                                    Delete
                                </button>
                            </div>
                        </td>
                    </tr>
                {/each}
            </tbody>
        </table>
    </div>
</div>
