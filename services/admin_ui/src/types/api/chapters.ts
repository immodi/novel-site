export interface AdminGetAllNovelChaptersResponse {
    chapters: ChapterResponse[];
    error: string;
}

export interface ChapterResponse {
    id: number;
    novelId: number;
    title: string;
    content: string;
    releaseDate: string;
}
