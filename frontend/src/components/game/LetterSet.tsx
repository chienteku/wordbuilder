import { motion } from 'framer-motion';
import { sortAndSplitLetters } from '@/lib/utils';

interface LetterButtonsProps {
    letters: string[];
    position: 'prefix' | 'suffix';
    type: 'vowels' | 'consonants';
    onAddLetter: (letter: string, position: 'prefix' | 'suffix') => void;
}

export const LetterButtons: React.FC<LetterButtonsProps> = ({
    letters,
    position,
    type,
    onAddLetter
}) => {
    if (letters.length === 0) {
        return <p className="text-sm text-gray-500">No {type} available</p>;
    }

    return (
        <div className="flex flex-wrap gap-1 justify-center">
            {letters.map((letter) => (
                <motion.button
                    key={`${type}-${letter}`}
                    onClick={() => onAddLetter(letter, position)}
                    className={`w-9 h-9 flex items-center justify-center rounded-full text-white text-lg font-medium
                   ${position === 'prefix' ? 'bg-blue-500 hover:bg-blue-600' : 'bg-green-500 hover:bg-green-600'}`}
                    whileHover={{ scale: 1.05 }}
                    whileTap={{ scale: 0.95 }}
                    transition={{ duration: 0.1 }}
                >
                    {letter}
                </motion.button>
            ))}
        </div>
    );
};

interface LetterSetProps {
    position: 'prefix' | 'suffix';
    letters: string[];
    onAddLetter: (letter: string, position: 'prefix' | 'suffix') => void;
}

export const LetterSet: React.FC<LetterSetProps> = ({ position, letters, onAddLetter }) => {
    const { vowels, consonants } = sortAndSplitLetters(letters);
    const title = position === 'prefix' ? 'Prefix Letters' : 'Suffix Letters';

    return (
        <div className={`w-full ${position === 'prefix' ? 'border-r border-gray-200' : ''}`}>
            <h3 className={`text-center font-semibold mb-2 ${position === 'prefix' ? 'text-blue-600' : 'text-green-600'}`}>
                {title}
            </h3>

            <div className="mb-3">
                <p className="text-xs text-gray-500 mb-1 text-center">Vowels</p>
                <LetterButtons
                    letters={vowels}
                    position={position}
                    type="vowels"
                    onAddLetter={onAddLetter}
                />
            </div>

            <div>
                <p className="text-xs text-gray-500 mb-1 text-center">Consonants</p>
                <LetterButtons
                    letters={consonants}
                    position={position}
                    type="consonants"
                    onAddLetter={onAddLetter}
                />
            </div>
        </div>
    );
};