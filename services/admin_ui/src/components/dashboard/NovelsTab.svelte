<script lang="ts">
    import { mapDbNovelsToNovelsDTO } from "../../lib/mappers/novels_mapper";
    import { novelsState } from "../../lib/states/novels_state.svelte";
    import type { Novel } from "../../types/dtos/novel";

    type NovelProps = {
        novels: Novel[];
        onNovelSelect: (novel: Novel) => void;
    };

    const novels = $derived(mapDbNovelsToNovelsDTO(novelsState.data ?? []));

    function refresh(): void {
        novelsState.refresh();
    }

    $effect(() => {
        if (!novelsState.data && !novelsState.loading) {
            refresh();
        }
    });
    const { onNovelSelect }: NovelProps = $props();

    let currentPage = $state(1);
    const itemsPerPage = 20;

    const totalPages = $derived(Math.ceil(novels.length / itemsPerPage));
    const paginatedNovels = $derived(
        novels.slice(
            (currentPage - 1) * itemsPerPage,
            currentPage * itemsPerPage,
        ),
    );

    function editNovel(id: number) {
        console.log("Edit novel", id);
    }

    function deleteNovel(id: number) {
        console.log("Delete novel", id);
    }

    function viewChapters(novel: Novel) {
        onNovelSelect(novel);
    }

    function goToPage(page: number) {
        if (page >= 1 && page <= totalPages) currentPage = page;
    }

    function nextPage() {
        if (currentPage < totalPages) currentPage++;
    }

    function previousPage() {
        if (currentPage > 1) currentPage--;
    }

    // Helper function to get status badge color
    function getStatusColor(status: string): string {
        switch (status.toLowerCase()) {
            case "published":
                return "bg-green-100 text-green-800";
            case "draft":
                return "bg-yellow-100 text-yellow-800";
            case "archived":
                return "bg-gray-100 text-gray-800";
            case "suspended":
                return "bg-red-100 text-red-800";
            default:
                return "bg-blue-100 text-blue-800";
        }
    }
</script>

<div class="p-4 sm:p-6">
    <!-- Header -->
    <div
        class="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4 mb-6"
    >
        <h2
            class="text-lg sm:text-xl font-bold text-[#19183B] text-center sm:text-left"
        >
            Novel Management
        </h2>
        <button
            class="cursor-pointer px-4 py-2 bg-[#19183B] text-white rounded-lg hover:bg-[#2a2852] transition-colors w-full sm:w-auto"
        >
            Add New Novel
        </button>
    </div>

    <!-- Mobile Card View -->
    <div class="block sm:hidden space-y-4">
        {#if novelsState.loading}
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
                    <span class="ml-2 text-[#19183B]">Loading novels...</span>
                </div>
            </div>
        {:else if paginatedNovels.length === 0}
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
                            d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"
                        />
                    </svg>
                    <h3 class="mt-2 text-sm font-medium text-gray-900">
                        No novels
                    </h3>
                    <p class="mt-1 text-sm text-gray-500">
                        Get started by creating a new novel.
                    </p>
                </div>
            </div>
        {:else}
            <!-- Novels Cards -->
            {#each paginatedNovels as novel}
                <div
                    class="bg-white rounded-lg border border-gray-200 p-4 shadow-sm"
                >
                    <!-- Novel Header -->
                    <div class="flex items-start justify-between mb-3">
                        <div class="min-w-0 flex-1">
                            <div class="flex items-center space-x-2 mb-2">
                                <span
                                    class="text-xs font-medium text-[#708993] bg-gray-100 px-2 py-1 rounded"
                                >
                                    ID: {novel.id}
                                </span>
                                <span
                                    class="text-xs font-medium text-[#19183B] bg-blue-50 px-2 py-1 rounded"
                                >
                                    {novel.views.toLocaleString()} views
                                </span>
                                <span
                                    class={`text-xs font-medium px-2 py-1 rounded ${getStatusColor(novel.status)}`}
                                >
                                    {novel.status}
                                </span>
                            </div>
                            <h3
                                class="text-base font-semibold text-[#19183B] mb-1"
                            >
                                {novel.title}
                            </h3>
                            <p class="text-sm text-[#708993]">
                                by {novel.author}
                            </p>
                        </div>
                    </div>

                    <!-- Actions -->
                    <div
                        class="flex flex-wrap gap-2 pt-3 border-t border-gray-100"
                    >
                        <button
                            class="cursor-pointer flex-1 min-w-[80px] px-3 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors text-xs font-medium text-center"
                            onclick={() => viewChapters(novel)}
                        >
                            View Chapters
                        </button>
                        <button
                            class="cursor-pointer flex-1 min-w-[60px] px-3 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors text-xs font-medium text-center"
                            onclick={() => editNovel(novel.id)}
                        >
                            Edit
                        </button>
                        <button
                            class="cursor-pointer flex-1 min-w-[70px] px-3 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors text-xs font-medium text-center"
                            onclick={() => deleteNovel(novel.id)}
                        >
                            Delete
                        </button>
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
                        class="w-[8%] px-4 py-3 text-left text-xs font-medium text-[#19183B] uppercase tracking-wider"
                    >
                        ID
                    </th>
                    <th
                        class="w-[30%] px-4 py-3 text-left text-xs font-medium text-[#19183B] uppercase tracking-wider"
                    >
                        Title
                    </th>
                    <th
                        class="w-[15%] px-4 py-3 text-left text-xs font-medium text-[#19183B] uppercase tracking-wider"
                    >
                        Author
                    </th>
                    <th
                        class="w-[12%] px-4 py-3 text-left text-xs font-medium text-[#19183B] uppercase tracking-wider"
                    >
                        Status
                    </th>
                    <th
                        class="w-[10%] px-4 py-3 text-left text-xs font-medium text-[#19183B] uppercase tracking-wider"
                    >
                        Views
                    </th>
                    <th
                        class="w-[25%] px-4 py-3 text-left text-xs font-medium text-[#19183B] uppercase tracking-wider"
                    >
                        Actions
                    </th>
                </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
                {#if novelsState.loading}
                    <!-- Loading State -->
                    <tr>
                        <td colspan="6" class="px-4 py-8 text-center">
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
                                    >Loading novels...</span
                                >
                            </div>
                        </td>
                    </tr>
                {:else if paginatedNovels.length === 0}
                    <!-- Empty State -->
                    <tr>
                        <td colspan="6" class="px-4 py-8 text-center">
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
                                        d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"
                                    />
                                </svg>
                                <h3
                                    class="mt-2 text-sm font-medium text-gray-900"
                                >
                                    No novels
                                </h3>
                                <p class="mt-1 text-sm text-gray-500">
                                    Get started by creating a new novel.
                                </p>
                            </div>
                        </td>
                    </tr>
                {:else}
                    {#each paginatedNovels as novel}
                        <tr class="hover:bg-gray-50 transition-colors">
                            <td
                                class="px-4 py-4 whitespace-nowrap text-sm text-[#19183B]"
                            >
                                {novel.id}
                            </td>
                            <td class="px-4 py-4">
                                <div class="text-sm font-medium text-[#19183B]">
                                    {novel.title}
                                </div>
                            </td>
                            <td class="px-4 py-4 whitespace-nowrap">
                                <div class="text-sm text-[#19183B]">
                                    {novel.author}
                                </div>
                            </td>
                            <td class="px-4 py-4 whitespace-nowrap">
                                <span
                                    class={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${getStatusColor(novel.status)}`}
                                >
                                    {novel.status}
                                </span>
                            </td>
                            <td
                                class="px-4 py-4 whitespace-nowrap text-sm text-[#19183B]"
                            >
                                {novel.views.toLocaleString()}
                            </td>
                            <td
                                class="px-4 py-4 whitespace-nowrap text-sm font-medium"
                            >
                                <div class="flex space-x-2">
                                    <button
                                        class="cursor-pointer px-3 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors flex items-center justify-center text-xs"
                                        onclick={() => viewChapters(novel)}
                                    >
                                        View Chapters
                                    </button>
                                    <button
                                        class="cursor-pointer px-3 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors flex items-center justify-center text-xs"
                                        onclick={() => editNovel(novel.id)}
                                    >
                                        Edit
                                    </button>
                                    <button
                                        class="cursor-pointer px-3 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors flex items-center justify-center text-xs"
                                        onclick={() => deleteNovel(novel.id)}
                                    >
                                        Delete
                                    </button>
                                </div>
                            </td>
                        </tr>
                    {/each}
                {/if}
            </tbody>
        </table>
    </div>

    <!-- Pagination -->
    {#if totalPages > 1}
        <div
            class="flex flex-col sm:flex-row items-center justify-between mt-6 gap-4 px-4 py-3 bg-white border-t border-gray-200 sm:px-6"
        >
            <!-- Mobile Pagination -->
            <div class="flex justify-between w-full sm:hidden">
                <button
                    class="cursor-pointer relative inline-flex items-center px-4 py-2 text-sm font-medium text-[#19183B] bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                    onclick={previousPage}
                    disabled={currentPage === 1}
                >
                    Previous
                </button>

                <!-- Current page indicator -->
                <div class="flex items-center">
                    <span class="text-sm text-[#19183B]">
                        Page {currentPage} of {totalPages}
                    </span>
                </div>

                <button
                    class="cursor-pointer relative inline-flex items-center px-4 py-2 text-sm font-medium text-[#19183B] bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                    onclick={nextPage}
                    disabled={currentPage === totalPages}
                >
                    Next
                </button>
            </div>

            <!-- Desktop Pagination -->
            <div
                class="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between"
            >
                <div>
                    <p class="text-sm text-[#708993]">
                        Showing
                        <span class="font-medium"
                            >{(currentPage - 1) * itemsPerPage + 1}</span
                        >
                        to
                        <span class="font-medium"
                            >{Math.min(
                                currentPage * itemsPerPage,
                                novels.length,
                            )}</span
                        >
                        of
                        <span class="font-medium">{novels.length}</span>
                        results
                    </p>
                </div>
                <div>
                    <nav
                        class="isolate inline-flex -space-x-px rounded-md shadow-sm"
                        aria-label="Pagination"
                    >
                        <button
                            class="cursor-pointer relative inline-flex items-center rounded-l-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0 disabled:opacity-50 disabled:cursor-not-allowed"
                            onclick={previousPage}
                            disabled={currentPage === 1}
                        >
                            <span class="sr-only">Previous</span>
                            <svg
                                class="h-5 w-5"
                                viewBox="0 0 20 20"
                                fill="currentColor"
                                aria-hidden="true"
                            >
                                <path
                                    fill-rule="evenodd"
                                    d="M12.79 5.23a.75.75 0 01-.02 1.06L8.832 10l3.938 3.71a.75.75 0 11-1.04 1.08l-4.5-4.25a.75.75 0 010-1.08l4.5-4.25a.75.75 0 011.06.02z"
                                    clip-rule="evenodd"
                                />
                            </svg>
                        </button>
                        {#each Array.from({ length: totalPages }, (_, i) => i + 1) as page}
                            <button
                                class={`cursor-pointer relative inline-flex items-center px-4 py-2 text-sm font-semibold ${
                                    page === currentPage
                                        ? "z-10 bg-[#19183B] text-white focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#19183B]"
                                        : "text-[#19183B] ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0"
                                }`}
                                onclick={() => goToPage(page)}
                            >
                                {page}
                            </button>
                        {/each}
                        <button
                            class="cursor-pointer relative inline-flex items-center rounded-r-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0 disabled:opacity-50 disabled:cursor-not-allowed"
                            onclick={nextPage}
                            disabled={currentPage === totalPages}
                        >
                            <span class="sr-only">Next</span>
                            <svg
                                class="h-5 w-5"
                                viewBox="0 0 20 20"
                                fill="currentColor"
                                aria-hidden="true"
                            >
                                <path
                                    fill-rule="evenodd"
                                    d="M7.21 14.77a.75.75 0 01.02-1.06L11.168 10 7.23 6.29a.75.75 0 111.04-1.08l4.5 4.25a.75.75 0 010 1.08l-4.5 4.25a.75.75 0 01-1.06-.02z"
                                    clip-rule="evenodd"
                                />
                            </svg>
                        </button>
                    </nav>
                </div>
            </div>
        </div>
    {/if}
</div>
