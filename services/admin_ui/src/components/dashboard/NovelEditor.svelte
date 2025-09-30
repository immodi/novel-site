<script lang="ts">
    import type { Novel } from "../../types/dtos/novel";

    type Props = {
        editingNovel: Novel;
        onSave: (novel: Novel) => void;
        onCancel: () => void;
    };
    const { editingNovel, onSave, onCancel }: Props = $props();

    let editedTitle = $derived(editingNovel.title);
    let editedAuthor = $derived(editingNovel.author);
    let editedStatus = $derived(editingNovel.status.toLowerCase());
    let editedViews = $derived(editingNovel.views);
    const statusOptions = [
        { value: "ongoing", label: "Ongoing" },
        { value: "completed", label: "Completed" },
    ];

    function handleSave(): void {
        const updatedNovel: Novel = {
            ...editingNovel,
            title: editedTitle,
            author: editedAuthor,
            status: editedStatus,
            views: editedViews,
        };
        onSave(updatedNovel);
    }

    function handleCancel(): void {
        onCancel();
    }
</script>

<div class="mt-8 p-6 border-2 border-[#19183B] rounded-lg bg-[#E7F2EF]">
    <h3 class="text-lg font-medium text-[#19183B] mb-4">
        Edit Novel: {editingNovel.title}
    </h3>
    <div class="space-y-4">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
                <label
                    for="title"
                    class="block text-sm font-medium text-[#19183B] mb-1"
                    >Novel Title</label
                >
                <input
                    type="text"
                    id="title"
                    class="w-full px-3 py-2 border border-[#A1C2BD] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#19183B]"
                    bind:value={editedTitle}
                    placeholder="Enter novel title"
                />
            </div>
            <div>
                <label
                    for="author"
                    class="block text-sm font-medium text-[#19183B] mb-1"
                    >Author</label
                >
                <input
                    type="text"
                    id="author"
                    class="w-full px-3 py-2 border border-[#A1C2BD] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#19183B]"
                    bind:value={editedAuthor}
                    placeholder="Enter author name"
                />
            </div>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
                <label
                    for="status"
                    class="block text-sm font-medium text-[#19183B] mb-1"
                    >Status</label
                >
                <select
                    id="status"
                    class="w-full px-3 py-2 border border-[#A1C2BD] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#19183B] bg-white"
                    bind:value={editedStatus}
                >
                    {#each statusOptions as option}
                        <option value={option.value}>
                            {option.label}
                        </option>
                    {/each}
                </select>
            </div>
            <div>
                <label
                    for="views"
                    class="block text-sm font-medium text-[#19183B] mb-1"
                    >Views</label
                >
                <input
                    type="number"
                    id="views"
                    class="w-full px-3 py-2 border border-[#A1C2BD] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#19183B]"
                    bind:value={editedViews}
                    min="0"
                />
            </div>
        </div>

        <!-- Additional fields can be added here -->
        <div class="bg-blue-50 border border-blue-200 rounded-lg p-3">
            <div class="flex items-center">
                <svg
                    class="w-5 h-5 text-blue-500 mr-2"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                    ></path>
                </svg>
                <span class="text-sm text-blue-700">
                    Novel ID: <strong>{editingNovel.id}</strong>
                </span>
            </div>
        </div>

        <div class="flex space-x-2 pt-4">
            <button
                class="cursor-pointer px-6 py-2 bg-[#19183B] text-white rounded-lg hover:bg-[#2a2852] transition-colors font-medium"
                onclick={handleSave}
            >
                Save Changes
            </button>
            <button
                class="cursor-pointer px-6 py-2 bg-gray-500 text-white rounded-lg hover:bg-gray-600 transition-colors font-medium"
                onclick={handleCancel}
            >
                Cancel
            </button>
        </div>
    </div>
</div>
