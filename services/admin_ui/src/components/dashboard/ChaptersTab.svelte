<script lang="ts">
    import type { Chapter } from "../../types/dtos/chapter";
    import type { Novel } from "../../types/dtos/novel";
    import ChapterEditor from "./ChapterEditor.svelte";

    type ChapterProps = {
        chapters: Chapter[];
        selectedNovel: Novel | null;
    };

    let { chapters = [], selectedNovel = null }: ChapterProps = $props();

    // State for chapter editing
    let editingChapter: Chapter | null = $state(null);

    // Filter state
    let filterText = $state("");

    // Pagination state
    let currentPage = $state(1);
    const itemsPerPage = 20;

    // Computed values
    const filteredChapters = $derived(
        chapters.filter((chapter) => {
            if (!selectedNovel) return false;
            if (filterText.trim() === "")
                return chapter.novelId === selectedNovel.id;

            return (
                chapter.novelId === selectedNovel.id &&
                chapter.title.toLowerCase().includes(filterText.toLowerCase())
            );
        }),
    );

    const totalPages = $derived(
        Math.ceil(filteredChapters.length / itemsPerPage),
    );
    const paginatedChapters = $derived(
        filteredChapters.slice(
            (currentPage - 1) * itemsPerPage,
            currentPage * itemsPerPage,
        ),
    );

    function editChapter(chapterId: number): void {
        const chapter = chapters.find((c) => c.id === chapterId);
        if (chapter) {
            editingChapter = chapter;
        }
    }

    function deleteChapter(chapterId: number): void {
        console.log(`Delete chapter with ID: ${chapterId}`);
        // Future implementation will delete the chapter
        chapters = chapters.filter((c) => c.id !== chapterId);
    }

    function handleChapterSave(updatedChapter: Chapter): void {
        const chapterIndex = chapters.findIndex(
            (c) => c.id === updatedChapter.id,
        );
        if (chapterIndex !== -1) {
            chapters[chapterIndex] = updatedChapter;
        }
        editingChapter = null;
    }

    function handleCancelEdit(): void {
        editingChapter = null;
    }

    function handleFilterChange(event: Event): void {
        const target = event.target as HTMLInputElement;
        filterText = target.value;
        currentPage = 1; // Reset to first page when filtering
    }

    function clearFilter(): void {
        filterText = "";
        currentPage = 1;
    }

    function goToPage(page: number): void {
        if (page >= 1 && page <= totalPages) {
            currentPage = page;
        }
    }

    function nextPage(): void {
        if (currentPage < totalPages) {
            currentPage++;
        }
    }

    function previousPage(): void {
        if (currentPage > 1) {
            currentPage--;
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
            Chapter Management
        </h2>
        <button
            class="cursor-pointer px-4 py-2 bg-[#19183B] text-white rounded-lg hover:bg-[#2a2852] transition-colors w-full sm:w-auto"
        >
            Add New Chapter
        </button>
    </div>

    {#if !selectedNovel}
        <!-- Empty state when no novel is selected -->
        <div class="text-center py-8">
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
                No novel selected
            </h3>
            <p class="mt-1 text-sm text-gray-500">
                Select a novel from the Novels tab to view its chapters.
            </p>
        </div>
    {:else}
        <!-- Selected Novel Info -->
        <div class="bg-[#E7F2EF] border border-[#A1C2BD] rounded-lg p-4 mb-6">
            <div class="flex items-center justify-between">
                <div class="flex items-center space-x-3">
                    <div class="bg-[#19183B] text-white rounded-full p-2">
                        <svg
                            class="w-4 h-4"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"
                            />
                        </svg>
                    </div>
                    <div>
                        <p class="text-sm text-[#708993]">
                            Viewing chapters for:
                        </p>
                        <p class="text-base font-semibold text-[#19183B]">
                            {selectedNovel.title}
                        </p>
                    </div>
                </div>
                <div
                    class="text-sm text-[#708993] bg-white px-3 py-1 rounded-full"
                >
                    {filteredChapters.length} chapter{filteredChapters.length !==
                    1
                        ? "s"
                        : ""}
                </div>
            </div>
        </div>

        <!-- Filter Section -->
        <div class="mb-6">
            <div class="flex flex-col sm:flex-row sm:items-center gap-4">
                <div class="flex-1">
                    <label for="chapter-filter" class="sr-only"
                        >Filter chapters</label
                    >
                    <div class="relative rounded-md shadow-sm">
                        <div
                            class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"
                        >
                            <svg
                                class="h-5 w-5 text-gray-400"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                                />
                            </svg>
                        </div>
                        <input
                            id="chapter-filter"
                            type="text"
                            class="block w-full rounded-md border border-[#A1C2BD] py-2 pl-10 pr-10 focus:border-[#19183B] focus:outline-none focus:ring-2 focus:ring-[#19183B] sm:text-sm"
                            placeholder="Filter chapters by title..."
                            value={filterText}
                            oninput={handleFilterChange}
                        />
                        {#if filterText}
                            <button
                                class="absolute inset-y-0 right-0 flex items-center pr-3"
                                onclick={clearFilter}
                            >
                                <svg
                                    class="h-5 w-5 text-gray-400"
                                    viewBox="0 0 20 20"
                                    fill="currentColor"
                                >
                                    <path
                                        d="M6.28 5.22a.75.75 0 00-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 101.06 1.06L10 11.06l3.72 3.72a.75.75 0 101.06-1.06L11.06 10l3.72-3.72a.75.75 0 00-1.06-1.06L10 8.94 6.28 5.22z"
                                    />
                                </svg>
                            </button>
                        {/if}
                    </div>
                </div>
            </div>
        </div>

        <!-- Mobile Card View -->
        <div class="block sm:hidden space-y-4">
            {#if paginatedChapters.length === 0}
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
                                d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                            />
                        </svg>
                        <h3 class="mt-2 text-sm font-medium text-gray-900">
                            {#if filterText}
                                No chapters found
                            {:else}
                                No chapters
                            {/if}
                        </h3>
                        <p class="mt-1 text-sm text-gray-500">
                            {#if filterText}
                                No chapters found matching "{filterText}"
                            {:else}
                                This novel doesn't have any chapters yet.
                            {/if}
                        </p>
                    </div>
                </div>
            {:else}
                <!-- Chapters Cards -->
                {#each paginatedChapters as chapter}
                    <div
                        class="bg-white rounded-lg border border-gray-200 p-4 shadow-sm"
                    >
                        <!-- Chapter Header -->
                        <div class="flex items-start justify-between mb-3">
                            <div class="min-w-0 flex-1">
                                <div class="flex items-center space-x-2 mb-2">
                                    <span
                                        class="text-xs font-medium text-[#708993] bg-gray-100 px-2 py-1 rounded"
                                    >
                                        ID: {chapter.id}
                                    </span>
                                    <span
                                        class="text-xs font-medium text-[#19183B] bg-blue-50 px-2 py-1 rounded"
                                    >
                                        Novel: {chapter.novelId}
                                    </span>
                                </div>
                                <h3
                                    class="text-base font-semibold text-[#19183B] mb-1"
                                >
                                    {chapter.title}
                                </h3>
                                <p class="text-sm text-[#708993]">
                                    Released: {chapter.releaseDate}
                                </p>
                            </div>
                        </div>

                        <!-- Actions -->
                        <div
                            class="flex flex-wrap gap-2 pt-3 border-t border-gray-100"
                        >
                            <button
                                class="cursor-pointer flex-1 min-w-[70px] px-3 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors text-xs font-medium text-center"
                                onclick={() => editChapter(chapter.id)}
                            >
                                Edit
                            </button>
                            <button
                                class="cursor-pointer flex-1 min-w-[70px] px-3 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors text-xs font-medium text-center"
                                onclick={() => deleteChapter(chapter.id)}
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
                            class="w-[10%] px-6 py-3 text-left text-xs font-medium text-[#19183B] uppercase tracking-wider"
                        >
                            ID
                        </th>
                        <th
                            class="w-[10%] px-6 py-3 text-left text-xs font-medium text-[#19183B] uppercase tracking-wider"
                        >
                            Novel ID
                        </th>
                        <th
                            class="w-[40%] px-6 py-3 text-left text-xs font-medium text-[#19183B] uppercase tracking-wider"
                        >
                            Title
                        </th>
                        <th
                            class="w-[20%] px-6 py-3 text-left text-xs font-medium text-[#19183B] uppercase tracking-wider"
                        >
                            Release Date
                        </th>
                        <th
                            class="w-[20%] px-6 py-3 text-left text-xs font-medium text-[#19183B] uppercase tracking-wider"
                        >
                            Actions
                        </th>
                    </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                    {#if paginatedChapters.length === 0}
                        <tr>
                            <td
                                colspan="5"
                                class="px-6 py-8 text-center text-sm text-[#708993]"
                            >
                                {#if filterText}
                                    No chapters found matching "{filterText}"
                                {:else}
                                    No chapters found for this novel
                                {/if}
                            </td>
                        </tr>
                    {:else}
                        {#each paginatedChapters as chapter}
                            <tr class="hover:bg-gray-50 transition-colors">
                                <td
                                    class="px-6 py-4 whitespace-nowrap text-sm text-[#19183B]"
                                >
                                    {chapter.id}
                                </td>
                                <td
                                    class="px-6 py-4 whitespace-nowrap text-sm text-[#19183B]"
                                >
                                    {chapter.novelId}
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap">
                                    <div
                                        class="text-sm font-medium text-[#19183B]"
                                    >
                                        {chapter.title}
                                    </div>
                                </td>
                                <td
                                    class="px-6 py-4 whitespace-nowrap text-sm text-[#19183B]"
                                >
                                    {chapter.releaseDate}
                                </td>
                                <td
                                    class="px-6 py-4 whitespace-nowrap text-sm font-medium"
                                >
                                    <div class="flex space-x-2">
                                        <button
                                            class="cursor-pointer px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors flex items-center w-20 justify-center"
                                            onclick={() =>
                                                editChapter(chapter.id)}
                                        >
                                            Edit
                                        </button>
                                        <button
                                            class="cursor-pointer px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors flex items-center w-20 justify-center"
                                            onclick={() =>
                                                deleteChapter(chapter.id)}
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
                                    filteredChapters.length,
                                )}</span
                            >
                            of
                            <span class="font-medium"
                                >{filteredChapters.length}</span
                            >
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
    {/if}

    <!-- Chapter Editor -->
    {#if editingChapter}
        <ChapterEditor
            {editingChapter}
            onSave={handleChapterSave}
            onCancel={handleCancelEdit}
        />
    {/if}
</div>
