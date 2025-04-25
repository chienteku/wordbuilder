import { motion } from 'framer-motion';
import { sortAndSplitLetters } from '@/lib/utils';

interface LetterButtonsProps {
    allLetters: string[];
    enabledLetters: string[];
    position: 'prefix' | 'suffix';
    type: 'vowels' | 'consonants';
    onAddLetter: (letter: string, position: 'prefix' | 'suffix') => void;
}

export const LetterButtons: React.FC<LetterButtonsProps> = ({
    allLetters,
    enabledLetters,
    position,
    type,
    onAddLetter,
}) => {
    if (allLetters.length === 0) {
        return <p className="text-sm text-gray-500">No {type} available</p>;
    }

    return (
        <div className="flex flex-wrap gap-1 justify-center">
            {allLetters.map((letter) => {
                const enabled = enabledLetters.includes(letter);
                const baseColor = position === 'prefix' ? 'bg-blue-500' : 'bg-green-500';
                const hoverColor = position === 'prefix' ? 'hover:bg-blue-600' : 'hover:bg-green-600';

                return (
                    <motion.button
                        key={`${type}-${letter}`}
                        onClick={() => enabled && onAddLetter(letter, position)}
                        className={`w-9 h-9 flex items-center justify-center rounded-full text-white text-lg font-medium
                ${baseColor}
                ${enabled ? `${hoverColor} cursor-pointer` : 'opacity-40 cursor-default pointer-events-none'}
              `}
                        style={{
                            opacity: enabled ? 1 : 0.3,
                            transition: 'opacity 0.2s',
                        }}
                        whileHover={enabled ? { scale: 1.05 } : {}}
                        whileTap={enabled ? { scale: 0.95 } : {}}
                        transition={{ duration: 0.1 }}
                        disabled={!enabled}
                    >
                        {letter}
                    </motion.button>
                );
            })}
        </div>
    );
};

// Replace in LetterSet.tsx
interface LetterSetProps {
    position: 'prefix' | 'suffix';
    allLetters: string[];
    enabledLetters: string[];
    onAddLetter: (letter: string, position: 'prefix' | 'suffix') => void;
}

export const LetterSet: React.FC<LetterSetProps> = ({
    position,
    allLetters,
    enabledLetters,
    onAddLetter,
}) => {
    const { vowels, consonants } = sortAndSplitLetters(allLetters);
    const enabledVowels = vowels.filter((l) => enabledLetters.includes(l));
    const enabledConsonants = consonants.filter((l) => enabledLetters.includes(l));
    const title = position === 'prefix' ? 'Prefix Letters' : 'Suffix Letters';

    return (
        <div className={`w-full ${position === 'prefix' ? 'border-r border-gray-200' : ''}`}>
            <h3 className={`text-center font-semibold mb-2 ${position === 'prefix' ? 'text-blue-600' : 'text-green-600'}`}>
                {title}
            </h3>

            <div className="mb-3">
                <p className="text-xs text-gray-500 mb-1 text-center">Vowels</p>
                <LetterButtons
                    allLetters={vowels}
                    enabledLetters={enabledVowels}
                    position={position}
                    type="vowels"
                    onAddLetter={onAddLetter}
                />
            </div>

            <div>
                <p className="text-xs text-gray-500 mb-1 text-center">Consonants</p>
                <LetterButtons
                    allLetters={consonants}
                    enabledLetters={enabledConsonants}
                    position={position}
                    type="consonants"
                    onAddLetter={onAddLetter}
                />
            </div>
        </div>
    );
};