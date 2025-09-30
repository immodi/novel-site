import { apiRequest } from "../helpers/api";
import type { AdminDataResponse } from "../types/api/admin_data";

export async function getAdminData(token: string) {
    return apiRequest<AdminDataResponse>("/admin/data", {
        method: "GET",
        headers: {
            Authorization: `Bearer ${token}`,
        },
    });
}
