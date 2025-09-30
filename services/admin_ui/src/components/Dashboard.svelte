<script lang="ts">
    import type { Novel } from "../types/dtos/novel";
    import UsersTab from "./dashboard/UsersTab.svelte";
    import NovelsTab from "./dashboard/NovelsTab.svelte";
    import ChaptersTab from "./dashboard/ChaptersTab.svelte";
    import logoutIcon from "../assets/logout_icon.svg";
    import {
        clearUserData,
        setUserData,
    } from "../lib/states/auth_state.svelte";
    import { AUTH_COOKIE_NAME } from "../lib/constants";

    type Tab = "users" | "novels" | "chapters";
    type Props = {
        username: string;
        coverImage: string;
    };

    const { username, coverImage }: Props = $props();
    let activeTab: Tab = $state("users");
    let selectedNovel: Novel | null = $state(null);

    function handleTabChange(tab: Tab): void {
        activeTab = tab;
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

    function handleLogout(): void {
        document.cookie = `${AUTH_COOKIE_NAME}=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/;`;
        clearUserData();
    }
</script>

<div class="min-h-screen bg-gradient-to-br from-[#19183B] to-[#708993] p-6">
    <!-- Header -->
    <header class="mb-8">
        <div class="flex justify-between items-center">
            <h1 class="text-2xl font-bold text-white">Dashboard</h1>
            <div class="flex items-center space-x-4">
                <div class="text-white">{username}</div>

                <!-- Logout Button -->
                <button
                    class="cursor-pointer p-2 bg-[#E7F2EF] rounded-lg hover:bg-[#A1C2BD] transition-colors"
                    onclick={handleLogout}
                    title="Logout"
                >
                    <img src={logoutIcon} alt="Logout" class="w-5 h-5" />
                </button>

                {#if coverImage}
                    <img
                        alt=""
                        loading="lazy"
                        src={coverImage}
                        class="cursor-pointer w-10 h-10 rounded-full bg-[#E7F2EF]"
                    />
                {:else}
                    <div
                        class="cursor-pointer w-10 h-10 rounded-full bg-[#E7F2EF] flex items-center justify-center text-[#19183B] font-bold"
                    >
                        {username ? username[0].toUpperCase() : "?"}
                    </div>
                {/if}
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
                <ChaptersTab {selectedNovel} />
            {/if}
        </div>
    </main>
</div>
