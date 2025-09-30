export function setCookie<T>(
    name: string,
    value: T,
    days: number = 7,
    sameSite: "Strict" | "Lax" | "None" = "Lax",
    secure: boolean = sameSite === "None"
) {
    const date = new Date();
    date.setTime(date.getTime() + days * 24 * 60 * 60 * 1000);
    const expires = `; expires=${date.toUTCString()}`;
    const encoded = encodeURIComponent(JSON.stringify(value));

    let cookie = `${encodeURIComponent(name)}=${encoded}${expires}; path=/; SameSite=${sameSite}`;
    if (secure) {
        cookie += "; Secure";
    }

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

export function removeCookie(name: string) {
    document.cookie = `${encodeURIComponent(name)}=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/; SameSite=Lax`;
}
