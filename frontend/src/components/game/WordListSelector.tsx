import { useState } from 'react';
// import { WordList } from '@/types/wordlist';
// import { wordListService } from '@/services/wordlist-service';
import { useWordListContext } from '@/contexts/WordListContext';

interface WordListSelectorProps {
    activeWordListId: number | null;
    onSelectWordList: (wordListId: number) => Promise<void>;
}

export const WordListSelector: React.FC<WordListSelectorProps> = ({
    activeWordListId,
    onSelectWordList
}) => {
    // const [wordLists, setWordLists] = useState<WordList[]>([]);
    // const [loading, setLoading] = useState(true);
    // const [error, setError] = useState<string | null>(null);
    const [isChanging, setIsChanging] = useState(false);
    const { wordLists, loading, error } = useWordListContext();

    // useEffect(() => {
    //     fetchWordLists();
    // }, []);

    // const fetchWordLists = async () => {
    //     try {
    //         setLoading(true);
    //         const response = await wordListService.getAllWordLists();
    //         setWordLists(response.word_lists || []);
    //         setError(null);
    //     } catch (err: any) {
    //         setError('Failed to load word lists');
    //         console.error(err);
    //     } finally {
    //         setLoading(false);
    //     }
    // };

    const handleChange = async (e: React.ChangeEvent<HTMLSelectElement>) => {
        const wordListId = parseInt(e.target.value, 10);
        if (isNaN(wordListId)) return;

        try {
            setIsChanging(true);
            await onSelectWordList(wordListId);
        } catch (err: any) {
            console.error(err);
        } finally {
            setIsChanging(false);
        }
    };

    if (loading) {
        return (
            <div className="flex items-center text-sm text-gray-500">
                <span className="inline-block animate-spin rounded-full h-4 w-4 border-b-2 border-blue-500 mr-2"></span>
                Loading word lists...
            </div>
        );
    }

    if (wordLists.length === 0) {
        return null;
    }

    return (
        <div className="flex flex-col sm:flex-row items-center gap-2">
            <label htmlFor="word-list-select" className="text-sm font-medium whitespace-nowrap">
                Dictionary:
            </label>
            <div className="relative flex-grow min-w-[200px]">
                <select
                    id="word-list-select"
                    value={activeWordListId || ''}
                    onChange={handleChange}
                    disabled={isChanging}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm pr-10 appearance-none"
                >
                    {wordLists.map((wordList) => (
                        <option key={wordList.id} value={wordList.id}>
                            {wordList.name} ({wordList.word_count.toLocaleString()} words)
                        </option>
                    ))}
                </select>
                <div className="absolute inset-y-0 right-0 flex items-center pr-2 pointer-events-none">
                    {isChanging ? (
                        <span className="inline-block animate-spin rounded-full h-4 w-4 border-b-2 border-blue-500"></span>
                    ) : (
                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                            <path d="M6 9l6 6 6-6" />
                        </svg>
                    )}
                </div>
            </div>
            {error && (
                <span className="text-xs text-red-600">{error}</span>
            )}
        </div>
    );
};