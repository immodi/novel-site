export function setCookie<T>(name: string, value: T, days: number = 7) {
    const date = new Date();
    date.setTime(date.getTime() + days * 24 * 60 * 60 * 1000);
    const expires = `; expires=${date.toUTCString()}`;
    const encoded = encodeURIComponent(JSON.stringify(value));
    document.cookie = `${encodeURIComponent(name)}=${encoded}${expires}`;
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
    document.cookie = `${encodeURIComponent(name)}=; expires=Thu, 01 Jan 1970 00:00:00 GMT`;
}
