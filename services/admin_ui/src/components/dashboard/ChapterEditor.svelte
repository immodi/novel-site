<script lang="ts">
    import type { Chapter } from "../../types/dtos/chapter";

    export let editingChapter: Chapter;
    export let onSave: (chapter: Chapter) => void;
    export let onCancel: () => void;

    let editedTitle = editingChapter.title;
    let editedReleaseDate = editingChapter.releaseDate;
    let editedContent = editingChapter.content;

    function handleSave(): void {
        const updatedChapter: Chapter = {
            ...editingChapter,
            title: editedTitle,
            releaseDate: editedReleaseDate,
            content: editedContent,
        };
        onSave(updatedChapter);
    }

    function handleCancel(): void {
        onCancel();
    }
</script>

<div class="mt-8 p-6 border-2 border-[#19183B] rounded-lg bg-[#E7F2EF]">
    <h3 class="text-lg font-medium text-[#19183B] mb-4">
        Edit Chapter: {editingChapter.title}
    </h3>
    <div class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
            <div>
                <label class="block text-sm font-medium text-[#19183B] mb-1"
                    >Chapter Title</label
                >
                <input
                    type="text"
                    class="w-full px-3 py-2 border border-[#A1C2BD] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#19183B]"
                    bind:value={editedTitle}
                />
            </div>
            <div>
                <label class="block text-sm font-medium text-[#19183B] mb-1"
                    >Release Date</label
                >
                <input
                    type="date"
                    class="w-full px-3 py-2 border border-[#A1C2BD] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#19183B]"
                    bind:value={editedReleaseDate}
                />
            </div>
        </div>
        <div>
            <label class="block text-sm font-medium text-[#19183B] mb-1"
                >Content (HTML)</label
            >
            <div class="mb-2 text-xs text-[#708993]">
                This will be replaced with a rich text editor in the final
                implementation
            </div>
            <textarea
                class="w-full h-64 px-3 py-2 border border-[#A1C2BD] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#19183B] font-mono text-sm"
                bind:value={editedContent}
            ></textarea>
        </div>
        <div class="flex space-x-2">
            <button
                class="cursor-pointer px-4 py-2 bg-[#19183B] text-white rounded-lg hover:bg-[#2a2852] transition-colors"
                on:click={handleSave}
            >
                Save Changes
            </button>
            <button
                class="cursor-pointer px-4 py-2 bg-gray-500 text-white rounded-lg hover:bg-gray-600 transition-colors"
                on:click={handleCancel}
            >
                Cancel
            </button>
        </div>
    </div>
</div>
