import { DOMAIN } from "../lib/constants";

const isLocalhost = window.location.hostname === "localhost";

export function setCookie<T>(
    name: string,
    value: T,
    days: number = 7,
    sameSite: "Strict" | "Lax" | "None" = "Lax",
    secure: boolean = !isLocalhost,
    domain: string = DOMAIN
) {
    const date = new Date();
    date.setTime(date.getTime() + days * 24 * 60 * 60 * 1000);
    const expires = `; expires=${date.toUTCString()}`;
    const encoded = encodeURIComponent(JSON.stringify(value));

    let cookie = `${encodeURIComponent(name)}=${encoded}${expires}; path=/; SameSite=${sameSite}`;
    if (!isLocalhost) cookie += `; Domain=${domain}`;
    if (secure) cookie += "; Secure";

    document.cookie = cookie;
}

export function getCookie<T>(name: string): T | null {
    const match = document.cookie.match(
        new RegExp("(^| )" + encodeURIComponent(name) + "=([^;]+)")
    );
    if (!match) return null;
    try {
        return JSON.parse(decodeURIComponent(match[2])) as T;
    } catch {
        return null;
    }
}

export function removeCookie(name: string, domain: string = DOMAIN) {
    const isLocalhost = window.location.hostname === "localhost";
    let cookie = `${encodeURIComponent(name)}=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/; SameSite=Lax`;
    if (!isLocalhost) cookie += `; Domain=${domain}; Secure`;
    document.cookie = cookie;
}
