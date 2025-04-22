"use client"
import { useWordBuilder } from '@/hooks/useWordBuilder';
import { PageHeader } from '@/components/ui/PageHeader';
import { ErrorMessage } from '@/components/ui/ErrorMessage';
import { LetterSet } from '@/components/game/LetterSet';
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
    loading,
    addLetter,
    removeLetter,
    resetWordBuilder
  } = useWordBuilder();

  // Calculate min height to keep layout stable
  const detailsMinHeight = '150px';

  return (
    <div className="min-h-screen bg-gray-100 p-4 flex flex-col items-center">
      <div className="w-full max-w-lg">
        <PageHeader title="Word Builder" showHomeLink={true} />

        {state && (
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
        )}

        {/* Error Message */}
        <ErrorMessage message={error || ''} />
      </div>
    </div>
  );
};

export default GamePage;