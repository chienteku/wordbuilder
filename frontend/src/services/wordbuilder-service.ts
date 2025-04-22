import axios from 'axios';
import { ApiResponse } from '@/types/wordbuilder';

const API_BASE_URL = typeof window !== 'undefined'
    ? (window.location.hostname === 'localhost'
        ? 'http://localhost:8081/api/wordbuilder'
        : '/api/wordbuilder')
    : '/api/wordbuilder';

export const wordBuilderService = {
    // Initialize a new WordBuilder session
    initSession: async (): Promise<ApiResponse> => {
        const response = await axios.post<ApiResponse>(`${API_BASE_URL}/init`, {});
        return response.data;
    },

    // Get current state
    getState: async (sessionId: string): Promise<ApiResponse> => {
        const response = await axios.get<ApiResponse>(`${API_BASE_URL}/state?session_id=${sessionId}`);
        return response.data;
    },

    // Add a letter to the word
    addLetter: async (
        sessionId: string,
        letter: string,
        position: 'prefix' | 'suffix'
    ): Promise<ApiResponse> => {
        const response = await axios.post<ApiResponse>(`${API_BASE_URL}/add`, {
            session_id: sessionId,
            letter,
            position,
        });
        return response.data;
    },

    // Remove a letter from the word
    removeLetter: async (sessionId: string, index: number): Promise<ApiResponse> => {
        const response = await axios.post<ApiResponse>(`${API_BASE_URL}/remove`, {
            session_id: sessionId,
            index,
        });
        return response.data;
    },

    // Reset the WordBuilder
    resetSession: async (sessionId: string): Promise<ApiResponse> => {
        const response = await axios.post<ApiResponse>(`${API_BASE_URL}/reset`, {
            session_id: sessionId
        });
        return response.data;
    }
};

// Dictionary API service
export const dictionaryService = {
    // Fetch word details from an external dictionary API
    getWordDetails: async (word: string) => {
        const response = await axios.get(`https://api.dictionaryapi.dev/api/v2/entries/en/${word}`);
        const data = response.data[0];

        return {
            pronunciation: data.phonetics[0]?.text || '',
            audio: data.phonetics[0]?.audio || '',
            meaning: data.meanings[0]?.definitions[0]?.definition || '',
            example: data.meanings[0]?.definitions[0]?.example || ''
        };
    }
};