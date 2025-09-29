<script lang="ts">
    import type { Chapter } from "../../types/dtos/chapter";
    import ChapterEditor from "./ChapterEditor.svelte";

    export let chapters: Chapter[] = [];

    // State for chapter editing
    let editingChapter: Chapter | null = null;

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
</script>

<div class="p-6">
    <div class="flex justify-between items-center mb-6">
        <h2 class="text-xl font-bold text-[#19183B]">Chapter Management</h2>
        <button
            class="cursor-pointer px-4 py-2 bg-[#19183B] text-white rounded-lg hover:bg-[#2a2852] transition-colors"
        >
            Add New Chapter
        </button>
    </div>

    <!-- Chapters Table -->
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
                        Novel ID
                    </th>
                    <th
                        class="px-6 py-3 text-left text-xs font-medium text-[#19183B] uppercase tracking-wider"
                    >
                        Title
                    </th>
                    <th
                        class="px-6 py-3 text-left text-xs font-medium text-[#19183B] uppercase tracking-wider"
                    >
                        Release Date
                    </th>
                    <th
                        class="px-6 py-3 text-left text-xs font-medium text-[#19183B] uppercase tracking-wider"
                    >
                        Actions
                    </th>
                </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
                {#each chapters as chapter}
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
                            <div class="text-sm font-medium text-[#19183B]">
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
                                    on:click={() => editChapter(chapter.id)}
                                >
                                    Edit
                                </button>
                                <button
                                    class="cursor-pointer px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors flex items-center w-20 justify-center"
                                    on:click={() => deleteChapter(chapter.id)}
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

    <!-- Chapter Editor -->
    {#if editingChapter}
        <ChapterEditor
            {editingChapter}
            onSave={handleChapterSave}
            onCancel={handleCancelEdit}
        />
    {/if}
</div>
