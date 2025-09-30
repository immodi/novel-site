import type { ChapterResponse } from "../../types/api/chapters";
import { getAllNovelChapters } from "../../api/chapters";
import { getUserToken } from "./auth_state.svelte";

class ChaptersState {
    data = $state<ChapterResponse[] | null>(null);
    loading = $state(false);
    error = $state<string | null>(null);

    async refresh(novelID: number) {
        this.loading = true;
        this.error = null;
        try {
            const token = getUserToken() as string;
            const { data, error: networkError } = await getAllNovelChapters(token, novelID);
            if (networkError) {
                this.error = "Network error";
                this.data = [];
            } else if (data?.error) {
                this.error = data.error;
                this.data = [];
            } else {
                this.data = data?.chapters ?? [];
            }
        } catch (err) {
            this.error = (err as Error).message;
            this.data = [];
        } finally {
            this.loading = false;
        }
    }

    setNovels(novels: ChapterResponse[] | null) {
        this.data = novels;
        this.error = null;
        this.loading = false;
    }

    reset() {
        this.data = null;
        this.loading = false;
        this.error = null;
    }
}

export const chaptersState = new ChaptersState();
