export interface AdminGetAllUsersResponse {
    users: UserResponse[];
    error: string;
}

export interface UserResponse {
    id: number;
    username: string;
    email: string;
    role: string;
    createdAt: string;
    image: string;
}
