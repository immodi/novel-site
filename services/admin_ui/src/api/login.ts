import { apiRequest } from "../helpers/api";
import type { AdminLoginRequest, AdminLoginResponse } from "../types/api/login";

export async function login(body: AdminLoginRequest) {
    return apiRequest<AdminLoginResponse>("/admin/login", {
        method: "POST",
        body: JSON.stringify(body),
    });
}
