import type { NovelResponse } from "../../types/api/novels";
import { getAllNovels } from "../../api/novels";
import { getUserToken } from "./auth_state.svelte";

class NovelsState {
    data = $state<NovelResponse[] | null>(null);
    loading = $state(false);
    error = $state<string | null>(null);

    async refresh() {
        this.loading = true;
        this.error = null;
        try {
            const token = getUserToken() as string;
            const { data, error: networkError } = await getAllNovels(token);
            if (networkError) {
                this.error = "Network error";
                this.data = [];
            } else if (data?.error) {
                this.error = data.error;
                this.data = [];
            } else {
                this.data = data?.novels ?? [];
            }
        } catch (err) {
            this.error = (err as Error).message;
            this.data = [];
        } finally {
            this.loading = false;
        }
    }

    setNovels(novels: NovelResponse[] | null) {
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

export const novelsState = new NovelsState();
