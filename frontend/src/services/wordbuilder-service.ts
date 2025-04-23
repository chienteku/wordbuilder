import axios from 'axios';
import { ApiResponse } from '@/types/wordbuilder';

const API_BASE_URL = typeof window !== 'undefined'
    ? (window.location.hostname === 'localhost'
        ? 'http://localhost:8081/api/wordbuilder'
        : '/api/wordbuilder')
    : '/api/wordbuilder';

const DICTIONARY_API_BASE_URL = typeof window !== 'undefined'
    ? (window.location.hostname === 'localhost'
        ? 'http://localhost:8081/api/dictionary'
        : '/api/dictionary')
    : '/api/dictionary';

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
    // Updated to use our backend proxy for word details
    getWordDetails: async (word: string) => {
        try {
            // Call our backend endpoint that gets word details
            const response = await axios.get(`${DICTIONARY_API_BASE_URL}/details/${encodeURIComponent(word)}`);
            return response.data;
        } catch (err) {
            console.error('Failed to fetch word details:', err);
            // Return default empty values if the API call fails
            return {
                pronunciation: '',
                audio: '',
                meaning: 'Definition not available',
                example: ''
            };
        }
    },

    // Method to fetch a word image from our backend proxy
    getWordImage: async (word: string) => {
        try {
            const response = await axios.get(`${DICTIONARY_API_BASE_URL}/image/${encodeURIComponent(word)}`);
            return response.data.imageUrl || null;
        } catch (error) {
            console.error('Failed to fetch word image:', error);
            return null;
        }
    },

    // Optimized method that fetches both details and image in one request
    getCompleteWordDetails: async (word: string) => {
        try {
            const response = await axios.get(`${DICTIONARY_API_BASE_URL}/complete/${encodeURIComponent(word)}`);
            return response.data;
        } catch (err) {
            console.error('Failed to fetch complete word details:', err);
            // Return default empty values if the API call fails
            return {
                pronunciation: '',
                audio: '',
                meaning: 'Definition not available',
                example: '',
                imageUrl: null
            };
        }
    }
};