import { getAllUsers } from "../../api/users";
import type { UserResponse } from "../../types/api/users";
import { getUserToken } from "./auth_state.svelte";

class UsersState {
    data = $state<UserResponse[] | null>(null);
    loading = $state(false);
    error = $state<string | null>(null);

    async refresh() {
        this.loading = true;
        this.error = null;
        try {
            const token = getUserToken() as string;
            const { data, error: networkError } = await getAllUsers(token);
            if (networkError) {
                this.error = "Network error";
                this.data = [];
            } else if (data?.error) {
                this.error = data.error;
                this.data = [];
            } else {
                this.data = data?.users ?? [];
            }
        } catch (err) {
            this.error = (err as Error).message;
            this.data = [];
        } finally {
            this.loading = false;
        }
    }

    setUsers(users: UserResponse[] | null) {
        this.data = users;
        this.error = null;
        this.loading = false;
    }

    reset() {
        this.data = null;
        this.loading = false;
        this.error = null;
    }
}

export const usersState = new UsersState();
