"use client"
import { useState, useEffect } from 'react';
import { useWordBuilder } from '@/hooks/useWordBuilder';
import { PageHeader } from '@/components/ui/PageHeader';
import { ErrorMessage } from '@/components/ui/ErrorMessage';
import { LetterSet } from '@/components/game/LetterSet';
import { WordListSelector } from '@/components/game/WordListSelector';
import { wordListService } from '@/services/wordlist-service';
import {
  CurrentWord,
  WordDetailsDisplay,
  SuggestedCompletions
} from '@/components/game/WordDisplay';
import { WordListProvider } from '@/contexts/WordListContext';

const ALL_LETTERS = [...'abcdefghijklmnopqrstuvwxyz'];

const GamePage = () => {
  const {
    state,
    error,
    wordDetails,
    resetWordBuilder,
    addLetter,
    removeLetter,
    loading
  } = useWordBuilder();

  const [activeWordListId, setActiveWordListId] = useState<number | null>(null);
  const [changeMessage, setChangeMessage] = useState<string | null>(null);

  // // Fetch word lists and validate activeWordListId on mount
  // useEffect(() => {
  //   async function validateWordList() {
  //     const listsResponse = await wordListService.getAllWordLists();
  //     const availableLists = listsResponse.word_lists || [];

  //     const storedWordListId = localStorage.getItem('activeWordListId');
  //     let wordListId = storedWordListId ? parseInt(storedWordListId, 10) : null;

  //     // If no word list selected or it was deleted, pick the first available
  //     if (!wordListId || !availableLists.some(wl => wl.id === wordListId)) {
  //       if (availableLists.length > 0) {
  //         wordListId = availableLists[0].id;
  //         setActiveWordListId(wordListId);
  //         localStorage.setItem('activeWordListId', wordListId.toString());
  //         await handleSelectWordList(wordListId); // This resets the builder
  //       }
  //     } else {
  //       setActiveWordListId(wordListId);
  //       // Always reset the builder to sync with the word list
  //       resetWordBuilder();
  //     }
  //   }
  //   validateWordList();
  //   // eslint-disable-next-line
  // }, []);

  // Load active word list ID from localStorage on component mount
  useEffect(() => {
    const storedWordListId = localStorage.getItem('activeWordListId');
    if (storedWordListId) {
      const wordListId = parseInt(storedWordListId, 10);
      if (!isNaN(wordListId)) {
        setActiveWordListId(wordListId);
      }
    }
  }, []);

  // Handle word list selection
  const handleSelectWordList = async (wordListId: number) => {
    try {
      await wordListService.useWordList(wordListId);
      setActiveWordListId(wordListId);

      // Save to localStorage
      localStorage.setItem('activeWordListId', wordListId.toString());

      // Show message
      setChangeMessage('Dictionary changed! Please start a new word.');

      // Reset the game
      resetWordBuilder();

      // Clear message after 5 seconds
      setTimeout(() => {
        setChangeMessage(null);
      }, 5000);
    } catch (err) {
      throw err;
    }
  };

  return (
    <div className="min-h-screen bg-gray-100 p-4 flex flex-col items-center">
      <div className="w-full max-w-lg">
        <div className="flex flex-col gap-4 mb-6">
          <PageHeader title="Word Builder" showHomeLink={true} />
          <WordListSelector
            activeWordListId={activeWordListId}
            onSelectWordList={handleSelectWordList}
          />
        </div>

        {changeMessage && (
          <div className="mb-4 p-3 bg-blue-100 text-blue-800 rounded-md">
            {changeMessage}
          </div>
        )}

        {loading && !state ? (
          <div className="bg-white rounded-lg shadow-md p-8 text-center">
            <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mb-4"></div>
            <p>Loading dictionary...</p>
          </div>
        ) : state ? (
          <>
            {/* Current Word Section */}
            <CurrentWord
              answer={state.answer}
              step={state.step}
              onReset={resetWordBuilder}
              onRemoveLetter={removeLetter}
            />
            {/* Letter Sets - Split into Left/Right */}
            <div className="flex bg-white rounded-b-lg border-t border-gray-200 shadow p-4">
              <div className="w-1/2">
                <LetterSet
                  position="prefix"
                  allLetters={ALL_LETTERS}
                  enabledLetters={state.prefix_set}
                  onAddLetter={addLetter}
                />
              </div>
              <div className="w-1/2">
                <LetterSet
                  position="suffix"
                  allLetters={ALL_LETTERS}
                  enabledLetters={state.suffix_set}
                  onAddLetter={addLetter}
                />
              </div>
            </div>

            {/* Word Details with fixed min-height */}
            <div
              className="transition-all duration-300 bg-white px-4 pb-4 min-h-[600px] max-h-[600px] overflow-y-auto flex flex-col justify-start"
            >
              {state.is_valid_word && state.answer.length > 1 && wordDetails && (
                <WordDetailsDisplay
                  word={state.answer}
                  details={wordDetails}
                />
              )}
              {state.valid_completions && state.valid_completions.length > 0 && (
                <SuggestedCompletions
                  completions={state.valid_completions}
                  suggestion={state.suggestion}
                />
              )}
            </div>

            {/* Mobile Layout Info */}
            <div className="mt-4 text-center text-sm text-gray-500">
              <p>Use left side for prefix letters, right side for suffix letters</p>
            </div>
          </>
        ) : (
          <div className="bg-white rounded-lg shadow-md p-8 text-center">
            <p className="text-red-600 mb-4">Failed to load word builder. Please try again.</p>
            <button
              onClick={() => resetWordBuilder()}
              className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
            >
              Retry
            </button>
          </div>
        )}

        {/* Error Message */}
        <ErrorMessage message={error || ''} />
      </div>
    </div>
  );
};

const GamePageWrapper = () => (
  <WordListProvider>
    <GamePage />
  </WordListProvider>
);

export default GamePageWrapper;