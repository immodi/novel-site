declare global {
    interface Window {
        __ENV__?: {
            VITE_API_URL?: string;
            VITE_WS_URL?: string;
            VITE_DOMAIN?: string;
        };
    }
}

export const DOMAIN: string =
    window.__ENV__?.VITE_DOMAIN ?? import.meta.env.VITE_DOMAIN;
export const API_URL: string =
    window.__ENV__?.VITE_API_URL ?? import.meta.env.VITE_API_URL;
export const WS_URL: string =
    window.__ENV__?.VITE_WS_URL ?? import.meta.env.VITE_WS_URL;

export const AUTH_COOKIE_NAME = 'admin_auth_cookie'
export const ITEMS_PER_PAGINATION_PAGE = 20
