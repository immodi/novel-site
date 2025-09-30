export interface AdminLoginRequest {
    email: string;
    password: string;
}

export interface AdminLoginResponse {
    token: string;
    username: string;
    coverImage: string;
    error: string;
}
