import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { WordList } from '@/types/wordlist';
import { wordListService } from '@/services/wordlist-service';

interface WordSampleModalProps {
    wordList: WordList;
    onClose: () => void;
}

export const WordSampleModal: React.FC<WordSampleModalProps> = ({
    wordList,
    onClose
}) => {
    const [words, setWords] = useState<string[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [sampleSize, setSampleSize] = useState(100);

    useEffect(() => {
        const fetchSample = async () => {
            try {
                setLoading(true);
                const response = await wordListService.getWordListSample(wordList.id, sampleSize);
                setWords(response.words);
                setError('');
            } catch (err: any) {
                setError(err.response?.data?.error || 'Failed to load word sample');
            } finally {
                setLoading(false);
            }
        };

        fetchSample();
    }, [wordList.id, sampleSize]);

    const handleSampleSizeChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
        setSampleSize(Number(e.target.value));
    };

    return (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
            <motion.div
                initial={{ opacity: 0, scale: 0.9 }}
                animate={{ opacity: 1, scale: 1 }}
                exit={{ opacity: 0, scale: 0.9 }}
                className="bg-white rounded-lg shadow-lg max-w-2xl w-full max-h-[80vh] flex flex-col"
            >
                <div className="p-4 border-b flex justify-between items-center">
                    <h2 className="text-xl font-semibold">Sample Words from "{wordList.name}"</h2>
                    <button
                        onClick={onClose}
                        className="text-gray-500 hover:text-gray-700"
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                            <path d="M18 6L6 18M6 6l12 12" />
                        </svg>
                    </button>
                </div>

                <div className="p-4 border-b flex justify-between items-center">
                    <div className="flex items-center gap-2">
                        <label htmlFor="sampleSize" className="text-sm font-medium">Sample size:</label>
                        <select
                            id="sampleSize"
                            value={sampleSize}
                            onChange={handleSampleSizeChange}
                            className="border rounded px-2 py-1 text-sm"
                        >
                            <option value={50}>50 words</option>
                            <option value={100}>100 words</option>
                            <option value={250}>250 words</option>
                            <option value={500}>500 words</option>
                        </select>
                    </div>

                    <div className="text-sm text-gray-600">
                        Total: {wordList.word_count.toLocaleString()} words
                    </div>
                </div>

                <div className="p-4 overflow-y-auto flex-grow">
                    {loading ? (
                        <div className="flex items-center justify-center py-8">
                            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
                        </div>
                    ) : error ? (
                        <div className="p-4 bg-red-100 text-red-700 rounded">
                            {error}
                        </div>
                    ) : words.length === 0 ? (
                        <div className="text-center text-gray-500 py-8">
                            No words found in this word list.
                        </div>
                    ) : (
                        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-2">
                            {words.map((word, index) => (
                                <div
                                    key={index}
                                    className="bg-gray-100 px-3 py-2 rounded text-sm hover:bg-gray-200 transition-colors"
                                >
                                    {word}
                                </div>
                            ))}
                        </div>
                    )}
                </div>

                <div className="p-4 border-t flex justify-end">
                    <button
                        onClick={onClose}
                        className="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 transition-colors"
                    >
                        Close
                    </button>
                </div>
            </motion.div>
        </div>
    );
};