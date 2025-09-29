export interface AdminGetAllUsersRequest {
    token: string;
}

export interface AdminGetAllUsersResponse {
    users: User[];
    error: string;
}

export interface User {
    id: number;
    username: string;
    email: string;
    role: string;
    createdAt: string;
    image: string;
}
