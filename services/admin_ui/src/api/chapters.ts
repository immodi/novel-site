import { apiRequest } from "../helpers/api";
import type { AdminGetAllNovelChaptersResponse } from "../types/api/chapters";

export async function getAllNovelChapters(token: string, novelId: number) {
    return apiRequest<AdminGetAllNovelChaptersResponse>(`/admin/novels/${novelId}/chapters`, {
        method: "GET",
        headers: {
            Authorization: `Bearer ${token}`,
        },
    });
}
