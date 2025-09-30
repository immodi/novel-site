import { API_URL } from "../lib/constants";

export interface ApiResult<T> {
    data: T | null;
    error: string | null;
}


export async function apiRequest<T>(
    url: string,
    options: RequestInit = {}
): Promise<ApiResult<T>> {
    try {
        const res = await fetch(`${API_URL}${url}`, {
            headers: {
                "Content-Type": "application/json",
                ...(options.headers || {}),
            },
            ...options,
        });

        let data: any = null;
        try {
            data = await res.json();
        } catch { console.log("Failed to parse JSON"); }

        if (!res.ok) {
            const backendError = data?.error as string | undefined;
            return { data: null, error: backendError || `HTTP ${res.status} ${res.statusText}` };
        }

        return { data: data as T, error: null };
    } catch (err: any) {
        return { data: null, error: err.message || "Unknown error" };
    }
}

