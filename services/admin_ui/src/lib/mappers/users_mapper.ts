import type { UserResponse } from "../../types/api/users";
import type { User, UserRole } from "../../types/dtos/user";

export function mapDbUsersToUsersDTO(dbUsers: UserResponse[]): User[] {
    return dbUsers.map((dbUser) => ({
        id: dbUser.id,
        username: dbUser.username,
        email: dbUser.email,
        role: dbUser.role as UserRole,
        createdAt: dbUser.createdAt,
        coverImage: dbUser.image
    }));
}



