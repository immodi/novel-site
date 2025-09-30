import type { NovelResponse } from "../../types/api/novels";
import type { Novel } from "../../types/dtos/novel";

export function mapDbNovelsToNovelsDTO(dbNovels: NovelResponse[]): Novel[] {
    return dbNovels.map((dbNovel) => ({
        id: dbNovel.id,
        title: dbNovel.title,
        views: dbNovel.views,
        author: dbNovel.author,
        status: dbNovel.status
    }));
}



