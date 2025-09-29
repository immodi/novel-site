import { getAllUsers } from "../../api/get_all_users";
import type { User } from "../../types/api/users_data";
import { getUserToken } from "./auth.svelte";

type UsersState = {
    data: User[] | null;
    loading: boolean;
    error: string | null;
};

let users = $state<UsersState>({
    data: null,
    loading: false,
    error: null,
});

export function getUsers(): UsersState {
    if (!users.data && !users.loading) {
        requestUsers();
    }
    return users;
}

export function setUsers(_users: User[]) {
    users.data = _users;
    users.error = null;
    users.loading = false;
}

async function requestUsers() {
    users.loading = true;
    users.error = null;

    try {
        const token = getUserToken() as string;
        const { data, error: networkError } = await getAllUsers(token);
        console.log({ data, error: networkError });

        if (networkError) {
            users.error = "Network error";
            users.data = [];
        } else if (data?.error) {
            users.error = data.error;
            users.data = [];
        } else if (data?.users) {
            users.data = data.users;
        }
    } catch (err) {
        users.error = (err as Error).message;
        users.data = [];
    } finally {
        users.loading = false;
    }
}
