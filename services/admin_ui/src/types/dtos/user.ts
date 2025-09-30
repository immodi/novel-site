export type UserRole = 'admin' | 'user';
export interface User {
    id: number;
    username: string;
    email: string;
    role: UserRole;
    createdAt: string;
    coverImage: string;
}
