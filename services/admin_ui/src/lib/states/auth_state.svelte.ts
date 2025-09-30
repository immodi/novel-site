import { getCookie, removeCookie, setCookie } from "../../helpers/cookies";
import { AUTH_COOKIE_NAME } from "../constants";

export type UserToken = string | null
let userToken = $state<UserToken | null>(getCookie(AUTH_COOKIE_NAME));

export function getUserToken(): string | null {
    return userToken
}

export function clearUserToken() {
    userToken = null
    removeCookie(AUTH_COOKIE_NAME)
}

export function setUserToken(token: string) {
    userToken = token
    setCookie(AUTH_COOKIE_NAME, userToken);
}
