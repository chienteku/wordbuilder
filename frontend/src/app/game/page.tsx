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

  // Calculate min height to keep layout stable
  const detailsMinHeight = '150px';

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

            {/* Word Details with fixed min-height */}
            <div style={{ minHeight: detailsMinHeight }} className="transition-all duration-300 bg-white px-4 pb-4">
              {state.is_valid_word && state.answer.length > 1 && wordDetails ? (
                <WordDetailsDisplay
                  word={state.answer}
                  details={wordDetails}
                />
              ) : !state.is_valid_word && state.valid_completions && state.valid_completions.length > 0 ? (
                <SuggestedCompletions
                  completions={state.valid_completions}
                  suggestion={state.suggestion}
                />
              ) : null}
            </div>

            {/* Letter Sets - Split into Left/Right */}
            <div className="flex bg-white rounded-b-lg border-t border-gray-200 shadow p-4">
              <div className="w-1/2">
                <LetterSet
                  position="prefix"
                  letters={state.prefix_set}
                  onAddLetter={addLetter}
                />
              </div>
              <div className="w-1/2">
                <LetterSet
                  position="suffix"
                  letters={state.suffix_set}
                  onAddLetter={addLetter}
                />
              </div>
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

export default GamePage;