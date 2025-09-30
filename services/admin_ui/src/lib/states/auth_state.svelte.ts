import { getCookie, setCookie } from "../../helpers/cookies";
import { AUTH_COOKIE_NAME } from "../constants";

export type UserData = {
    token: string
    coverImage: string
    username: string
};
let userData = $state<UserData | null>(null);

export function getUserData(): UserData | null {
    return userData ?? getCookie(AUTH_COOKIE_NAME);
}

export function getUserToken(): string | null {
    return userData?.token ?? (getCookie(AUTH_COOKIE_NAME) as UserData | null)?.token ?? null;
}

export function clearUserData() {
    userData = null
    setCookie(AUTH_COOKIE_NAME, null)
}

export function setUserData(data: UserData) {
    userData = data;
    setCookie(AUTH_COOKIE_NAME, data);
}

export function setUserToken(token: string) {
    userData = userData
        ? { ...userData, token }
        : { token, coverImage: "", username: "" };

    setCookie(AUTH_COOKIE_NAME, userData);
}
