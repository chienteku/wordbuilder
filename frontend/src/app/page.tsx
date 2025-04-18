"use client"
import { useState, useEffect } from 'react';
import axios from 'axios';
import { motion } from 'framer-motion';

interface WordBuilderState {
  answer: string;
  prefix_set: string[];
  suffix_set: string[];
  step: number;
  is_valid_word: boolean;
  valid_completions?: string[]; // New field
  suggestion?: string;          // New field
}

interface WordDetails {
  pronunciation: string;
  audio: string;
  meaning: string;
  example: string;
}

interface ApiResponse {
  session_id?: string;
  state: WordBuilderState;
  success?: boolean;
  message?: string;
  error?: string;
}

const API_BASE_URL = typeof window !== 'undefined'
  ? (window.location.hostname === 'localhost'
    ? 'http://localhost:8081/api/wordbuilder'
    : '/api/wordbuilder')
  : '/api/wordbuilder';

const HomePage: React.FC = () => {
  const [sessionId, setSessionId] = useState<string | null>(null);
  const [state, setState] = useState<WordBuilderState | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [wordDetails, setWordDetails] = useState<WordDetails | null>(null);

  // Initialize WordBuilder
  const initializeWordBuilder = async () => {
    try {
      const response = await axios.post<ApiResponse>(`${API_BASE_URL}/init`, {});
      const newSessionId = response.data.session_id || null;
      setSessionId(newSessionId);
      setState(response.data.state);
      setError(null);
      if (newSessionId) {
        localStorage.setItem('sessionId', newSessionId);
      }
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to initialize WordBuilder.');
    }
  };

  // Add letter
  const addLetter = async (letter: string, position: 'prefix' | 'suffix') => {
    if (!sessionId || !state) return;
    try {
      const response = await axios.post<ApiResponse>(`${API_BASE_URL}/add`, {
        session_id: sessionId,
        letter,
        position,
      });
      if (response.data.success) {
        setState(response.data.state);
        setError(null);

        // Only fetch word details if it's valid AND has more than one letter
        if (response.data.state.is_valid_word && response.data.state.answer.length > 1) {
          fetchWordDetails(response.data.state.answer);
        } else {
          setWordDetails(null);
        }
      } else {
        setError(response.data.message || 'Failed to add letter.');
      }
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to add letter.');
    }
  };

  // Remove letter
  const removeLetter = async (index: number) => {
    if (!sessionId || !state) return;
    try {
      const response = await axios.post<ApiResponse>(`${API_BASE_URL}/remove`, {
        session_id: sessionId,
        index,
      });
      if (response.data.success) {
        setState(response.data.state);
        setError(null);

        // Only fetch word details if it's valid AND has more than one letter
        if (response.data.state.is_valid_word && response.data.state.answer.length > 1) {
          fetchWordDetails(response.data.state.answer);
        } else {
          setWordDetails(null);
        }
      } else {
        setError(response.data.message || 'Failed to remove letter.');
      }
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to remove letter.');
      if (err.response?.status === 404) {
        localStorage.removeItem('sessionId');
        setSessionId(null);
        setState(null);
      }
    }
  };

  // Reset WordBuilder
  const resetWordBuilder = async () => {
    if (!sessionId) return;
    try {
      const response = await axios.post<ApiResponse>(`${API_BASE_URL}/reset`, {
        session_id: sessionId
      });
      if (response.data.success) {
        setState(response.data.state);
        setWordDetails(null);
        setError(null);
      } else {
        setError(response.data.message || 'Failed to reset word builder.');
      }
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to reset word builder.');
      if (err.response?.status === 404) {
        localStorage.removeItem('sessionId');
        setSessionId(null);
        setState(null);
        initializeWordBuilder();
      }
    }
  };

  // Restore state from localStorage or initialize new session
  useEffect(() => {
    const savedSessionId = localStorage.getItem('sessionId');
    if (savedSessionId && !sessionId) {
      axios
        .get<ApiResponse>(`${API_BASE_URL}/state?session_id=${savedSessionId}`)
        .then((response) => {
          setSessionId(savedSessionId);
          setState(response.data.state);

          // If the word is valid, fetch its details
          if (response.data.state.is_valid_word && response.data.state.answer.length > 1) {
            fetchWordDetails(response.data.state.answer);
          }
        })
        .catch(() => {
          localStorage.removeItem('sessionId');
          initializeWordBuilder();
        });
    } else if (!sessionId) {
      initializeWordBuilder();
    }
  }, [sessionId]);

  // Helper function to check if a letter is a vowel
  const isVowel = (letter: string): boolean => {
    return ['a', 'e', 'i', 'o', 'u'].includes(letter.toLowerCase());
  };

  // Function to sort and split letters into vowels and consonants
  const sortAndSplitLetters = (letters: string[]) => {
    const vowels = letters.filter(letter => isVowel(letter)).sort();
    const consonants = letters.filter(letter => !isVowel(letter)).sort();
    return { vowels, consonants };
  };

  // When a valid word is detected
  const fetchWordDetails = async (word: string) => {
    try {
      const response = await axios.get(`https://api.dictionaryapi.dev/api/v2/entries/en/${word}`);
      const data = response.data[0];
      setWordDetails({
        pronunciation: data.phonetics[0]?.text || '',
        audio: data.phonetics[0]?.audio || '',
        meaning: data.meanings[0]?.definitions[0]?.definition || '',
        example: data.meanings[0]?.definitions[0]?.example || ''
      });
    } catch (err) {
      console.error('Failed to fetch word details:', err);
    }
  };

  // Component to render letter buttons
  const LetterButtons = ({
    letters,
    position,
    type
  }: {
    letters: string[],
    position: 'prefix' | 'suffix',
    type: 'vowels' | 'consonants'
  }) => {
    if (letters.length === 0) {
      return <p className="text-sm text-gray-500">No {type} available</p>;
    }

    return (
      <div className="flex flex-wrap gap-1 justify-center">
        {letters.map((letter) => (
          <motion.button
            key={`${type}-${letter}`}
            onClick={() => addLetter(letter, position)}
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

  // LetterSet component for displaying a set of letters (prefix or suffix)
  const LetterSet = ({ position, letters }: { position: 'prefix' | 'suffix', letters: string[] }) => {
    const { vowels, consonants } = sortAndSplitLetters(letters);
    const title = position === 'prefix' ? 'Prefix Letters' : 'Suffix Letters';
    const alignment = position === 'prefix' ? 'items-end pr-2' : 'items-start pl-2';

    return (
      <div className={`w-full ${position === 'prefix' ? 'border-r border-gray-200' : ''}`}>
        <h3 className={`text-center font-semibold mb-2 ${position === 'prefix' ? 'text-blue-600' : 'text-green-600'}`}>
          {title}
        </h3>

        <div className="mb-3">
          <p className="text-xs text-gray-500 mb-1 text-center">Vowels</p>
          <LetterButtons letters={vowels} position={position} type="vowels" />
        </div>

        <div>
          <p className="text-xs text-gray-500 mb-1 text-center">Consonants</p>
          <LetterButtons letters={consonants} position={position} type="consonants" />
        </div>
      </div>
    );
  };

  // Component to display suggested completions
  const SuggestedCompletions = () => {
    if (!state || !state.valid_completions || state.valid_completions.length === 0) return null;
    
    return (
      <div className="mt-4 p-3 bg-gray-50 rounded-lg border border-gray-200">
        <h4 className="font-medium text-sm text-gray-700 mb-2">Possible words:</h4>
        <div className="flex flex-wrap gap-2">
          {state.valid_completions.map(word => (
            <span key={word} className="px-2 py-1 bg-blue-100 text-blue-800 rounded text-xs">
              {word}
            </span>
          ))}
        </div>
        {state.suggestion && (
          <p className="mt-2 text-sm text-gray-600 italic">{state.suggestion}</p>
        )}
      </div>
    );
  };

  // Calculate min height to keep layout stable
  const detailsMinHeight = '150px';

  return (
    <div className="min-h-screen bg-gray-100 p-4 flex flex-col items-center">
      <div className="w-full max-w-lg">
        <h1 className="text-2xl md:text-3xl font-bold text-center mb-6">Word Builder</h1>

        {state && (
          <>
            {/* Current Word Section with fixed height */}
            <div className="bg-white rounded-t-lg shadow p-4">
              <div className="flex items-center justify-between mb-3">
                <div className="flex items-center gap-2">
                  <h2 className="text-lg font-semibold">Current Word</h2>
                  <motion.button
                    onClick={resetWordBuilder}
                    className="p-1 bg-gray-200 text-gray-700 rounded-full hover:bg-gray-300 flex items-center justify-center"
                    whileHover={{ scale: 1.1 }}
                    whileTap={{ scale: 0.95 }}
                    disabled={state.answer.length === 0}
                    aria-label="Reset word"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                      <path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"></path>
                      <path d="M3 3v5h5"></path>
                    </svg>
                  </motion.button>
                </div>
                <span className="text-xs bg-gray-100 px-2 py-1 rounded text-gray-600">
                  Step: {state.step}
                </span>
              </div>

              <div className="flex justify-center mb-4">
                {state.answer.length > 0 ? (
                  <div className="flex gap-1 flex-wrap justify-center">
                    {state.answer.split('').map((letter, index) => (
                      <motion.div
                        key={`${letter}-${index}`}
                        onClick={() => removeLetter(index)}
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

              {/* Word Details with fixed min-height */}
              <div style={{ minHeight: detailsMinHeight }} className="transition-all duration-300">
                {state.is_valid_word && state.answer.length > 1 && wordDetails ? (
                  <motion.div
                    className="p-4 bg-gray-50 rounded-lg border border-gray-200"
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    transition={{ duration: 0.3 }}
                  >
                    <div className="flex items-center flex-wrap mb-2">
                      <h3 className="font-semibold mr-2">{state.answer}</h3>
                      {wordDetails.pronunciation && (
                        <span className="text-gray-600 text-sm">/{wordDetails.pronunciation}/</span>
                      )}
                      {wordDetails.audio && (
                        <button
                          onClick={() => new Audio(wordDetails.audio).play()}
                          className="ml-2 text-blue-500"
                        >
                          🔊
                        </button>
                      )}
                    </div>
                    <p className="text-gray-800 mb-2 text-sm">{wordDetails.meaning}</p>
                    {wordDetails.example && (
                      <p className="text-gray-600 italic text-sm">"{wordDetails.example}"</p>
                    )}
                  </motion.div>
                ) : !state.is_valid_word && state.valid_completions && state.valid_completions.length > 0 ? (
                  <SuggestedCompletions />
                ) : null}
              </div>
            </div>

            {/* Letter Sets - Split into Left/Right */}
            <div className="flex bg-white rounded-b-lg border-t border-gray-200 shadow p-4">
              <div className="w-1/2">
                <LetterSet position="prefix" letters={state.prefix_set} />
              </div>
              <div className="w-1/2">
                <LetterSet position="suffix" letters={state.suffix_set} />
              </div>
            </div>

            {/* Mobile Layout Info */}
            <div className="mt-4 text-center text-sm text-gray-500">
              <p>Use left side for prefix letters, right side for suffix letters</p>
            </div>
          </>
        )}

        {/* Error Message */}
        {error && (
          <motion.div
            className="mt-4 p-3 bg-red-100 text-red-700 rounded-lg"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 0.3 }}
          >
            <p className="text-sm text-center">{error}</p>
          </motion.div>
        )}
      </div>
    </div>
  );
};

export default HomePage;