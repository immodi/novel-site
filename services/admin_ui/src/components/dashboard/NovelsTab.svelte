<script lang="ts">
    import type { Novel } from "../../types/dtos/novel";

    export let novels: Novel[] = [];

    function editNovel(novelId: number): void {
        console.log(`Edit novel with ID: ${novelId}`);
        // Future implementation will open a form to edit the novel
    }

    function deleteNovel(novelId: number): void {
        console.log(`Delete novel with ID: ${novelId}`);
        // Future implementation will delete the novel
    }
</script>

<div class="p-6">
    <div class="flex justify-between items-center mb-6">
        <h2 class="text-xl font-bold text-[#19183B]">Novel Management</h2>
        <button
            class="cursor-pointer px-4 py-2 bg-[#19183B] text-white rounded-lg hover:bg-[#2a2852] transition-colors"
        >
            Add New Novel
        </button>
    </div>

    <!-- Novels Table -->
    <div class="overflow-x-auto w-full">
        <table class="w-full divide-y divide-gray-200">
            <thead>
                <tr class="bg-gray-50">
                    <th
                        class="px-6 py-3 text-left text-xs font-medium text-[#19183B] uppercase tracking-wider"
                    >
                        ID
                    </th>
                    <th
                        class="px-6 py-3 text-left text-xs font-medium text-[#19183B] uppercase tracking-wider"
                    >
                        Title
                    </th>
                    <th
                        class="px-6 py-3 text-left text-xs font-medium text-[#19183B] uppercase tracking-wider"
                    >
                        Views
                    </th>
                    <th
                        class="px-6 py-3 text-left text-xs font-medium text-[#19183B] uppercase tracking-wider"
                    >
                        Actions
                    </th>
                </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
                {#each novels as novel}
                    <tr class="hover:bg-gray-50 transition-colors">
                        <td
                            class="px-6 py-4 whitespace-nowrap text-sm text-[#19183B]"
                        >
                            {novel.id}
                        </td>
                        <td class="px-6 py-4 whitespace-nowrap">
                            <div class="text-sm font-medium text-[#19183B]">
                                {novel.title}
                            </div>
                        </td>
                        <td
                            class="px-6 py-4 whitespace-nowrap text-sm text-[#19183B]"
                        >
                            {novel.viewCount.toLocaleString()}
                        </td>
                        <td
                            class="px-6 py-4 whitespace-nowrap text-sm font-medium"
                        >
                            <div class="flex space-x-2">
                                <button
                                    class="cursor-pointer px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors flex items-center w-20 justify-center"
                                    on:click={() => editNovel(novel.id)}
                                >
                                    Edit
                                </button>
                                <button
                                    class="cursor-pointer px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors flex items-center w-20 justify-center"
                                    on:click={() => deleteNovel(novel.id)}
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
