import { apiRequest } from "../helpers/api";
import type { AdminGetAllNovelsResponse } from "../types/api/novels";

export async function getAllNovels(token: string) {
    return apiRequest<AdminGetAllNovelsResponse>("/admin/novels", {
        method: "GET",
        headers: {
            Authorization: `Bearer ${token}`,
        },
    });
}
