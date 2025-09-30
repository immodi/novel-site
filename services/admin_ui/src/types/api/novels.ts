export interface NovelResponse {
    id: number;
    title: string;
    views: number;
    author: string;
    status: string;
}

export interface AdminGetAllNovelsResponse {
    novels: NovelResponse[];
    error: string;
}
