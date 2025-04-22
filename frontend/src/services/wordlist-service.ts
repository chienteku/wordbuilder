import axios from 'axios';
import {
  WordListResponse,
  WordListDetailResponse,
  WordSampleResponse,
  CreateWordListResponse,
  UpdateWordListResponse,
  DeleteWordListResponse,
  UseWordListResponse,
} from '@/types/wordlist';

const API_BASE_URL = typeof window !== 'undefined'
  ? (window.location.hostname === 'localhost'
    ? 'http://localhost:8081/api/wordlists'
    : '/api/wordlists')
  : '/api/wordlists';

export const wordListService = {
  // Get all word lists
  getAllWordLists: async (): Promise<WordListResponse> => {
    const response = await axios.get<WordListResponse>(API_BASE_URL);
    return response.data;
  },

  // Get a specific word list
  getWordList: async (id: number): Promise<WordListDetailResponse> => {
    const response = await axios.get<WordListDetailResponse>(`${API_BASE_URL}/${id}`);
    return response.data;
  },

  // Get a sample of words from a word list
  getWordListSample: async (id: number, limit = 100): Promise<WordSampleResponse> => {
    const response = await axios.get<WordSampleResponse>(`${API_BASE_URL}/${id}/sample?limit=${limit}`);
    return response.data;
  },

  // Create a new word list
  createWordList: async (formData: FormData): Promise<CreateWordListResponse> => {
    const response = await axios.post<CreateWordListResponse>(API_BASE_URL, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },

  // Update an existing word list
  updateWordList: async (id: number, formData: FormData): Promise<UpdateWordListResponse> => {
    const response = await axios.put<UpdateWordListResponse>(`${API_BASE_URL}/${id}`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },

  // Delete a word list
  deleteWordList: async (id: number): Promise<DeleteWordListResponse> => {
    const response = await axios.delete<DeleteWordListResponse>(`${API_BASE_URL}/${id}`);
    return response.data;
  },

  // Use a word list as the active dictionary
  useWordList: async (id: number): Promise<UseWordListResponse> => {
    const response = await axios.post<UseWordListResponse>(`${API_BASE_URL}/${id}/use`);
    return response.data;
  },

  // Download a word list
  downloadWordList: (id: number): void => {
    // Redirect the browser to the download URL
    window.location.href = `${API_BASE_URL}/${id}/download`;
  },
};