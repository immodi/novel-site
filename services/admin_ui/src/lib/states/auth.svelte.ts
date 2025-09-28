import { getCookie, setCookie } from "../../helpers/cookies";
import { AUTH_COOKIE_NAME } from "../constants";

export type UserToken = string | null;
let userToken = $state<UserToken>(null);

export function getUserToken() {
    return userToken || getCookie(AUTH_COOKIE_NAME);
}

export function setUserToken(token: string) {
    userToken = token;
    setCookie(AUTH_COOKIE_NAME, token)
}
