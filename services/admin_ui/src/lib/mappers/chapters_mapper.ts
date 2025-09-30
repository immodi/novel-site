import type { ChapterResponse } from "../../types/api/chapters";
import type { Chapter } from "../../types/dtos/chapter";

export function mapDbChaptersToChaptersDTO(dbChapters: ChapterResponse[]): Chapter[] {
    return dbChapters.map((dbChapter) => ({
        id: dbChapter.id,
        novelId: dbChapter.novelId,
        title: dbChapter.title,
        content: dbChapter.content,
        releaseDate: dbChapter.releaseDate
    }));
}



