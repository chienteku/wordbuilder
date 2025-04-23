"use client"
import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { settingsService, Settings } from '@/services/settings-service';
import { PageHeader } from '@/components/ui/PageHeader';

export default function SettingsPage() {
    const [settings, setSettings] = useState<Settings>({ pixabay_api_key: '' });
    const [isLoading, setIsLoading] = useState(true);
    const [isSaving, setIsSaving] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [successMessage, setSuccessMessage] = useState<string | null>(null);

    // Fetch current settings
    useEffect(() => {
        const fetchSettings = async () => {
            try {
                setIsLoading(true);
                const response = await settingsService.getSettings();
                setSettings(response.settings);
                setError(null);
            } catch (err: any) {
                setError(err.response?.data?.error || 'Failed to load settings');
            } finally {
                setIsLoading(false);
            }
        };

        fetchSettings();
    }, []);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            setIsSaving(true);
            setError(null);
            setSuccessMessage(null);
            
            await settingsService.updateSettings(settings);
            setSuccessMessage('Settings saved successfully!');
            
            // Clear success message after 3 seconds
            setTimeout(() => {
                setSuccessMessage(null);
            }, 3000);
        } catch (err: any) {
            setError(err.response?.data?.error || 'Failed to save settings');
        } finally {
            setIsSaving(false);
        }
    };

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setSettings(prev => ({ ...prev, [name]: value }));
    };

    return (
        <div className="min-h-screen bg-gray-100 p-4">
            <div className="max-w-2xl mx-auto">
                <PageHeader title="Settings" showHomeLink={true} />

                <div className="bg-white rounded-lg shadow-md p-6">
                    <h2 className="text-xl font-semibold mb-4">Application Settings</h2>

                    {isLoading ? (
                        <div className="flex justify-center items-center py-8">
                            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
                        </div>
                    ) : (
                        <form onSubmit={handleSubmit}>
                            <div className="mb-4">
                                <label htmlFor="pixabay_api_key" className="block text-sm font-medium text-gray-700 mb-1">
                                    Pixabay API Key
                                </label>
                                <input
                                    type="text"
                                    id="pixabay_api_key"
                                    name="pixabay_api_key"
                                    value={settings.pixabay_api_key}
                                    onChange={handleInputChange}
                                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                                    placeholder="Enter your Pixabay API key"
                                />
                                <p className="mt-1 text-sm text-gray-500">
                                    This API key is used to fetch images for words in the game. 
                                    You can get a free API key from <a href="https://pixabay.com/api/docs/" target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:underline">Pixabay</a>.
                                </p>
                            </div>

                            {error && (
                                <motion.div
                                    initial={{ opacity: 0, y: -10 }}
                                    animate={{ opacity: 1, y: 0 }}
                                    className="mb-4 p-3 bg-red-100 text-red-700 rounded-md"
                                >
                                    {error}
                                </motion.div>
                            )}

                            {successMessage && (
                                <motion.div
                                    initial={{ opacity: 0, y: -10 }}
                                    animate={{ opacity: 1, y: 0 }}
                                    className="mb-4 p-3 bg-green-100 text-green-700 rounded-md"
                                >
                                    {successMessage}
                                </motion.div>
                            )}

                            <div className="flex justify-end">
                                <button
                                    type="submit"
                                    disabled={isSaving}
                                    className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors disabled:bg-blue-400"
                                >
                                    {isSaving ? 'Saving...' : 'Save Settings'}
                                </button>
                            </div>
                        </form>
                    )}
                </div>
                
                <div className="mt-6 bg-white rounded-lg shadow-md p-6">
                    <h2 className="text-xl font-semibold mb-4">About Pixabay API</h2>
                    <div className="prose">
                        <p className="text-gray-700 mb-3">
                            The Word Builder app uses the Pixabay API to display images for words, enhancing the learning experience.
                        </p>
                        <h3 className="text-lg font-medium mb-2">How to get a Pixabay API key:</h3>
                        <ol className="list-decimal list-inside space-y-2 text-gray-700">
                            <li>Visit the <a href="https://pixabay.com/api/docs/" target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:underline">Pixabay API documentation</a></li>
                            <li>Create a free Pixabay account if you don&apos;t have one</li>
                            <li>After logging in, you&apos;ll find your API key on the documentation page</li>
                            <li>Copy the API key and paste it in the field above</li>
                            <li>Click &quot;Save Settings&quot; to store your API key</li>
                        </ol>
                        <p className="mt-4 text-sm text-gray-600">
                            Note: The free Pixabay API has usage limits. For more information, please refer to their documentation.
                        </p>
                    </div>
                </div>
            </div>
        </div>
    );
}