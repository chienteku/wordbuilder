import { motion } from 'framer-motion';
import { WordDetails } from '@/types/wordbuilder';

interface CurrentWordProps {
    answer: string;
    step: number;
    onReset: () => void;
    onRemoveLetter: (index: number) => void;
}

export const CurrentWord: React.FC<CurrentWordProps> = ({
    answer,
    step,
    onReset,
    onRemoveLetter
}) => {
    return (
        <div className="bg-white rounded-t-lg shadow p-4">
            <div className="flex items-center justify-between mb-3">
                <div className="flex items-center gap-2">
                    <h2 className="text-lg font-semibold">Current Word</h2>
                    <motion.button
                        onClick={onReset}
                        className="p-1 bg-gray-200 text-gray-700 rounded-full hover:bg-gray-300 flex items-center justify-center"
                        whileHover={{ scale: 1.1 }}
                        whileTap={{ scale: 0.95 }}
                        disabled={answer.length === 0}
                        aria-label="Reset word"
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                            <path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"></path>
                            <path d="M3 3v5h5"></path>
                        </svg>
                    </motion.button>
                </div>
                <span className="text-xs bg-gray-100 px-2 py-1 rounded text-gray-600">
                    Step: {step}
                </span>
            </div>

            <div className="flex justify-center mb-4">
                {answer.length > 0 ? (
                    <div className="flex gap-1 flex-wrap justify-center">
                        {answer.split('').map((letter, index) => (
                            <motion.div
                                key={`${letter}-${index}`}
                                onClick={() => onRemoveLetter(index)}
                                className="w-9 h-9 flex items-center justify-center bg-gray-200 rounded-full cursor-pointer hover:bg-gray-300 text-lg font-medium"
                                whileHover={{ scale: 1.05 }}
                                initial={{ opacity: 0, y: -10 }}
                                animate={{ opacity: 1, y: 0 }}
                                transition={{ duration: 0.2 }}
                            >
                                {letter}
                            </motion.div>
                        ))}
                    </div>
                ) : (
                    <p className="text-gray-500 italic">Tap letters to build a word</p>
                )}
            </div>
        </div>
    );
};

interface WordDetailsDisplayProps {
    word: string;
    details: WordDetails;
}

export const WordDetailsDisplay: React.FC<WordDetailsDisplayProps> = ({ word, details }) => {
    return (
        <motion.div
            className="p-4 bg-gray-50 rounded-lg border border-gray-200"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 0.3 }}
        >
            <div className="flex items-center flex-wrap mb-2">
                <h3 className="font-semibold mr-2">{word}</h3>
                {details.pronunciation && (
                    <span className="text-gray-600 text-sm">/{details.pronunciation}/</span>
                )}
                {details.audio && (
                    <button
                        onClick={() => new Audio(details.audio).play()}
                        className="ml-2 text-blue-500"
                    >
                        🔊
                    </button>
                )}
            </div>

            {/* Updated image container with fixed dimensions of 640x480 */}
            {details.imageUrl && (
                <div className="mb-3 flex justify-center">
                    <img
                        src={details.imageUrl}
                        alt={`Image for ${word}`}
                        className="w-full max-w-md h-auto max-h-96 object-contain rounded-md mb-3"
                        style={{ maxWidth: "640px", maxHeight: "480px" }}
                        onError={(e) => {
                            // Hide the image on error
                            e.currentTarget.style.display = 'none';
                        }}
                    />
                </div>
            )}

            <p className="text-gray-800 mb-2 text-sm">{details.meaning}</p>
            {details.example && (
                <p className="text-gray-600 italic text-sm">"{details.example}"</p>
            )}
        </motion.div>
    );
};

interface SuggestedCompletionsProps {
    completions: string[];
    suggestion?: string;
}

export const SuggestedCompletions: React.FC<SuggestedCompletionsProps> = ({
    completions,
    suggestion
}) => {
    if (completions.length === 0) return null;

    return (
        <div className="mt-4 p-3 bg-gray-50 rounded-lg border border-gray-200">
            <h4 className="font-medium text-sm text-gray-700 mb-2">Possible words:</h4>
            <div className="flex flex-wrap gap-2">
                {completions.map(word => (
                    <span key={word} className="px-2 py-1 bg-blue-100 text-blue-800 rounded text-xs">
                        {word}
                    </span>
                ))}
            </div>
            {suggestion && (
                <p className="mt-2 text-sm text-gray-600 italic">{suggestion}</p>
            )}
        </div>
    );
};