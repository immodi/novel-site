<script lang="ts">
    import type { User } from "../types/dtos/user";
    import type { Novel } from "../types/dtos/novel";
    import type { Chapter } from "../types/dtos/chapter";
    import UsersTab from "./dashboard/UsersTab.svelte";
    import NovelsTab from "./dashboard/NovelsTab.svelte";
    import ChaptersTab from "./dashboard/ChaptersTab.svelte";
    import { getUsers } from "../lib/states/user_data.svelte";
    type Tab = "users" | "novels" | "chapters";

    let activeTab: Tab = $state("users");
    let users = $state(getUsers());

    let novels: Novel[] = [
        { id: 1, title: "The Lost Kingdom", viewCount: 15432 },
        { id: 2, title: "Echoes of Tomorrow", viewCount: 9876 },
        { id: 3, title: "Shadow and Light", viewCount: 23451 },
    ];

    let chapters: Chapter[] = [
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
            novelId: 2,
            title: "Chapter 1: New World",
            content: "<p>Entering a new world full of possibilities...</p>",
            releaseDate: "2023-06-01",
        },
    ];

    // Event handlers
    function handleTabChange(tab: "users" | "novels" | "chapters") {
        activeTab = tab;
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
                    on:click={() => handleTabChange("users")}
                >
                    Users
                </button>
                <button
                    class={`cursor-pointer px-4 py-2 rounded-lg font-medium transition-colors border-2 ${activeTab === "novels" ? "bg-[#19183B] text-white border-[#19183B]" : "text-[#19183B] bg-white border-[#A1C2BD] hover:bg-[#A1C2BD]"}`}
                    on:click={() => handleTabChange("novels")}
                >
                    Novels
                </button>
                <button
                    class={`cursor-pointer px-4 py-2 rounded-lg font-medium transition-colors border-2 ${activeTab === "chapters" ? "bg-[#19183B] text-white border-[#19183B]" : "text-[#19183B] bg-white border-[#A1C2BD] hover:bg-[#A1C2BD]"}`}
                    on:click={() => handleTabChange("chapters")}
                >
                    Chapters
                </button>
            </div>
        </div>

        <!-- Tab Content -->
        <div class="bg-white rounded-b-xl shadow-lg overflow-hidden">
            {#if activeTab === "users"}
                <UsersTab {users} />
            {:else if activeTab === "novels"}
                <NovelsTab {novels} />
            {:else if activeTab === "chapters"}
                <ChaptersTab {chapters} />
            {/if}
        </div>
    </main>
</div>
