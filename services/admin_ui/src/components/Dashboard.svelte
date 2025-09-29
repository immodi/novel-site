<script lang="ts">
    import type { Novel } from "../types/dtos/novel";
    import type { Chapter } from "../types/dtos/chapter";
    import UsersTab from "./dashboard/UsersTab.svelte";
    import NovelsTab from "./dashboard/NovelsTab.svelte";
    import ChaptersTab from "./dashboard/ChaptersTab.svelte";
    type Tab = "users" | "novels" | "chapters";

    let activeTab: Tab = $state("users");
    let selectedNovel: Novel | null = $state(null);

    let allChapters: Chapter[] = [
        {
            id: 1,
            novelId: 1,
            title: "Chapter 1: The Beginning",
            content: "<p>This is the beginning of the story...</p>",
            releaseDate: "2023-05-01",
        },
        {
            id: 2,
            novelId: 1,
            title: "Chapter 2: The Journey",
            content: "<p>The journey continues with exciting adventures...</p>",
            releaseDate: "2023-05-08",
        },
        {
            id: 3,
            novelId: 1,
            title: "Chapter 3: The Discovery",
            content: "<p>They discover the ancient ruins...</p>",
            releaseDate: "2023-05-15",
        },
        {
            id: 4,
            novelId: 2,
            title: "Chapter 1: New World",
            content: "<p>Entering a new world full of possibilities...</p>",
            releaseDate: "2023-06-01",
        },
        {
            id: 5,
            novelId: 2,
            title: "Chapter 2: First Contact",
            content: "<p>Meeting the inhabitants of the new world...</p>",
            releaseDate: "2023-06-08",
        },
        {
            id: 6,
            novelId: 3,
            title: "Chapter 1: Darkness Falls",
            content: "<p>The darkness begins to spread...</p>",
            releaseDate: "2023-07-01",
        },
        // Add more chapters to test pagination...
    ];

    function handleTabChange(tab: Tab): void {
        activeTab = tab;
        // Clear selected novel when switching away from chapters
        if (tab !== "chapters") {
            selectedNovel = null;
        }
    }

    function handleNovelSelect(novel: Novel): void {
        selectedNovel = novel;
        activeTab = "chapters";
    }

    function handleClearSelectedNovel(): void {
        selectedNovel = null;
    }
</script>

<div class="min-h-screen bg-gradient-to-br from-[#19183B] to-[#708993] p-6">
    <!-- Header -->
    <header class="mb-8">
        <div class="flex justify-between items-center">
            <h1 class="text-2xl font-bold text-white">Admin Dashboard</h1>
            <div class="flex items-center space-x-4">
                <div class="text-white">Admin User</div>
                <div
                    class="cursor-pointer w-10 h-10 rounded-full bg-[#E7F2EF] flex items-center justify-center text-[#19183B] font-bold"
                >
                    A
                </div>
            </div>
        </div>
    </header>

    <!-- Main Content -->
    <main class="max-w-7xl mx-auto">
        <!-- Tab Navigation -->
        <div class="bg-[#E7F2EF] rounded-t-xl p-4">
            <div class="flex space-x-4">
                <button
                    class={`cursor-pointer px-4 py-2 rounded-lg font-medium transition-colors border-2 ${activeTab === "users" ? "bg-[#19183B] text-white border-[#19183B]" : "text-[#19183B] bg-white border-[#A1C2BD] hover:bg-[#A1C2BD]"}`}
                    onclick={() => handleTabChange("users")}
                >
                    Users
                </button>
                <button
                    class={`cursor-pointer px-4 py-2 rounded-lg font-medium transition-colors border-2 ${activeTab === "novels" ? "bg-[#19183B] text-white border-[#19183B]" : "text-[#19183B] bg-white border-[#A1C2BD] hover:bg-[#A1C2BD]"}`}
                    onclick={() => handleTabChange("novels")}
                >
                    Novels
                </button>
                <button
                    class={`cursor-pointer px-4 py-2 rounded-lg font-medium transition-colors border-2 ${activeTab === "chapters" ? "bg-[#19183B] text-white border-[#19183B]" : "text-[#19183B] bg-white border-[#A1C2BD] hover:bg-[#A1C2BD]"}`}
                    onclick={() => handleTabChange("chapters")}
                >
                    Chapters
                </button>
            </div>
        </div>

        <!-- Selected Novel Banner (only shows in chapters tab) -->
        {#if activeTab === "chapters" && selectedNovel}
            <div class="bg-[#E7F2EF] px-4 py-2 border-b border-[#A1C2BD]">
                <div class="flex justify-between items-center">
                    <div class="flex items-center space-x-2">
                        <span class="text-sm text-[#708993]"
                            >Viewing chapters for:</span
                        >
                        <span class="text-sm font-medium text-[#19183B]"
                            >{selectedNovel.title}</span
                        >
                    </div>
                    <button
                        class="cursor-pointer text-sm text-[#19183B] hover:text-[#2a2852] transition-colors"
                        onclick={handleClearSelectedNovel}
                    >
                        Clear Selection
                    </button>
                </div>
            </div>
        {/if}

        <!-- Tab Content -->
        <div class="bg-white rounded-b-xl shadow-lg overflow-hidden">
            {#if activeTab === "users"}
                <UsersTab />
            {:else if activeTab === "novels"}
                <NovelsTab onNovelSelect={handleNovelSelect} />
            {:else if activeTab === "chapters"}
                <ChaptersTab chapters={allChapters} {selectedNovel} />
            {/if}
        </div>
    </main>
</div>
