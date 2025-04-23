export interface WordBuilderState {
    answer: string;
    prefix_set: string[];
    suffix_set: string[];
    step: number;
    is_valid_word: boolean;
    valid_completions?: string[];
    suggestion?: string;
}

export interface WordDetails {
    pronunciation: string;
    audio: string;
    meaning: string;
    example: string;
    imageUrl?: string;
}

export interface ApiResponse {
    session_id?: string;
    state: WordBuilderState;
    success?: boolean;
    message?: string;
    error?: string;
}