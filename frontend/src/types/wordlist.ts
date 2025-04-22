export interface WordList {
    id: number;
    name: string;
    description: string;
    source: string;
    file_path: string;
    word_count: number;
    created_at: string;
    updated_at: string;
}

export interface WordListResponse {
    word_lists: WordList[];
}

export interface WordListDetailResponse {
    word_list: WordList;
}

export interface WordSampleResponse {
    words: string[];
    count: number;
}

export interface CreateWordListResponse {
    message: string;
    word_list: WordList;
}

export interface UpdateWordListResponse {
    message: string;
    word_list: WordList;
}

export interface DeleteWordListResponse {
    message: string;
}

export interface UseWordListResponse {
    message: string;
}