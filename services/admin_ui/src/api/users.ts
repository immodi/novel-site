import { apiRequest } from "../helpers/api";
import type { AdminGetAllUsersResponse } from "../types/api/users";

export async function getAllUsers(token: string) {
    return apiRequest<AdminGetAllUsersResponse>("/admin/users", {
        method: "GET",
        headers: {
            Authorization: `Bearer ${token}`,
        },
    });
}
