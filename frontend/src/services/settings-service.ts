import axios from 'axios';

export interface Settings {
    pixabay_api_key: string;
}

export interface SettingsResponse {
    settings: Settings;
}

export interface UpdateSettingsResponse {
    message: string;
    settings: Settings;
}

const API_BASE_URL = typeof window !== 'undefined'
    ? (window.location.hostname === 'localhost'
        ? 'http://localhost:8081/api/settings'
        : '/api/settings')
    : '/api/settings';

export const settingsService = {
    // Get current settings
    getSettings: async (): Promise<SettingsResponse> => {
        const response = await axios.get<SettingsResponse>(API_BASE_URL);
        return response.data;
    },

    // Update settings
    updateSettings: async (settings: Settings): Promise<UpdateSettingsResponse> => {
        const response = await axios.post<UpdateSettingsResponse>(API_BASE_URL, settings);
        return response.data;
    }
};