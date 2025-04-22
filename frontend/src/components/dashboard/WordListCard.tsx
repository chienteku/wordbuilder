import { useState } from 'react';
import { motion } from 'framer-motion';
import { WordList } from '@/types/wordlist';

interface WordListCardProps {
    wordList: WordList;
    onEdit: (wordList: WordList) => void;
    onDelete: (wordList: WordList) => void;
    onDownload: (wordList: WordList) => void;
    onUse: (wordList: WordList) => void;
    onViewSample: (wordList: WordList) => void;
    isActive?: boolean;
}

export const WordListCard: React.FC<WordListCardProps> = ({
    wordList,
    onEdit,
    onDelete,
    onDownload,
    onUse,
    onViewSample,
    isActive = false,
}) => {
    const [isExpanded, setIsExpanded] = useState(false);

    // Format dates nicely
    const formatDate = (dateString: string) => {
        const date = new Date(dateString);
        return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
    };

    return (
        <motion.div
            className={`bg-white rounded-lg shadow-md overflow-hidden border ${isActive ? 'border-blue-500' : 'border-gray-200'}`}
            whileHover={{ y: -2 }}
            transition={{ duration: 0.2 }}
        >
            <div className="p-4 cursor-pointer" onClick={() => setIsExpanded(!isExpanded)}>
                <div className="flex justify-between items-center">
                    <h3 className="text-lg font-semibold flex items-center">
                        {isActive && (
                            <span className="inline-block w-3 h-3 bg-blue-500 rounded-full mr-2" title="Active dictionary"></span>
                        )}
                        {wordList.name}
                    </h3>
                    <div className="flex items-center gap-2">
                        <span className="text-sm text-gray-500">{wordList.word_count.toLocaleString()} words</span>
                        <motion.button
                            animate={{ rotate: isExpanded ? 180 : 0 }}
                            className="ml-2 text-gray-500 p-1"
                            aria-label="Toggle details"
                        >
                            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                <path d="M6 9l6 6 6-6" />
                            </svg>
                        </motion.button>
                    </div>
                </div>

                {!isExpanded && wordList.description && (
                    <p className="text-sm text-gray-600 mt-1 line-clamp-1">{wordList.description}</p>
                )}
            </div>

            {isExpanded && (
                <motion.div
                    initial={{ opacity: 0, height: 0 }}
                    animate={{ opacity: 1, height: 'auto' }}
                    exit={{ opacity: 0, height: 0 }}
                    className="px-4 pb-4"
                >
                    {wordList.description && (
                        <div className="mb-3">
                            <h4 className="text-sm font-medium text-gray-700">Description:</h4>
                            <p className="text-sm text-gray-600">{wordList.description}</p>
                        </div>
                    )}

                    {wordList.source && (
                        <div className="mb-3">
                            <h4 className="text-sm font-medium text-gray-700">Source:</h4>
                            <p className="text-sm text-gray-600">{wordList.source}</p>
                        </div>
                    )}

                    <div className="grid grid-cols-2 gap-x-4 gap-y-2 mb-3 text-sm">
                        <div>
                            <h4 className="font-medium text-gray-700">Created:</h4>
                            <p className="text-gray-600">{formatDate(wordList.created_at)}</p>
                        </div>
                        <div>
                            <h4 className="font-medium text-gray-700">Updated:</h4>
                            <p className="text-gray-600">{formatDate(wordList.updated_at)}</p>
                        </div>
                    </div>

                    <div className="flex flex-wrap gap-2 mt-4">
                        <button
                            onClick={() => onViewSample(wordList)}
                            className="px-3 py-1 text-xs bg-gray-200 text-gray-700 rounded hover:bg-gray-300"
                        >
                            View Sample
                        </button>
                        <button
                            onClick={() => onDownload(wordList)}
                            className="px-3 py-1 text-xs bg-gray-200 text-gray-700 rounded hover:bg-gray-300"
                        >
                            Download
                        </button>
                        {!isActive && (
                            <button
                                onClick={() => onUse(wordList)}
                                className="px-3 py-1 text-xs bg-blue-600 text-white rounded hover:bg-blue-700"
                            >
                                Use This List
                            </button>
                        )}
                        <button
                            onClick={() => onEdit(wordList)}
                            className="px-3 py-1 text-xs bg-yellow-100 text-yellow-700 rounded hover:bg-yellow-200"
                        >
                            Edit
                        </button>
                        <button
                            onClick={() => onDelete(wordList)}
                            className="px-3 py-1 text-xs bg-red-100 text-red-700 rounded hover:bg-red-200"
                        >
                            Delete
                        </button>
                    </div>
                </motion.div>
            )}
        </motion.div>
    );
};