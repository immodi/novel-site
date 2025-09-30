import { getAdminData } from "../../api/admin_data";
import type { AdminDataResponse } from "../../types/api/admin_data";
import { getUserToken } from "./auth_state.svelte";

class AdminDataState {
    emptyObject: AdminDataResponse = {
        username: "",
        coverImage: "",
        error: "",
    }

    data = $state<AdminDataResponse | null>(null);
    loading = $state(false);
    error = $state<string | null>(null);

    async refresh() {
        this.loading = true;
        this.error = null;
        try {
            const token = getUserToken() as string;
            const { data, error: networkError } = await getAdminData(token);
            if (networkError) {
                this.error = "Network error";
                this.data = {
                    ...this.emptyObject
                };
            } else if (data?.error) {
                this.error = data.error;
                this.data = {
                    ...this.emptyObject
                };
            } else {
                this.data = data ?? {
                    ...this.emptyObject
                };;
            }
        } catch (err) {
            this.error = (err as Error).message;
            this.data = {
                ...this.emptyObject
            };
        } finally {
            this.loading = false;
        }
    }

    setAdminData(adminData: AdminDataResponse | null) {
        this.data = adminData;
        this.error = null;
        this.loading = false;
    }

    reset() {
        this.data = null;
        this.loading = false;
        this.error = null;
    }
}

export const adminDataState = new AdminDataState();
