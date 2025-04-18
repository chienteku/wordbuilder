"use client"
import { useState, useEffect } from 'react';
import axios from 'axios';
import { motion } from 'framer-motion';

interface WordBuilderState {
  answer: string;
  prefix_set: string[];
  suffix_set: string[];
  step: number;
  is_valid_word: boolean; // Added this field to explicitly track if the word is valid
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

// Update API URL to use the correct port (8081) and support relative URLs when deployed
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

        // If the updated word is valid, fetch its details
        if (response.data.state.is_valid_word) {
          fetchWordDetails(response.data.state.answer);
        } else {
          // Clear word details if it's no longer a valid word
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
          if (response.data.state.is_valid_word) {
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
      return <p className="text-gray-500">No {type} available.</p>;
    }

    return (
      <div className="flex flex-wrap gap-2 mb-3">
        {letters.map((letter) => (
          <motion.button
            key={`${type}-${letter}`}
            onClick={() => addLetter(letter, position)}
            className="px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600"
            initial={{ opacity: 0, scale: 0.8 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.2 }}
          >
            {letter}
          </motion.button>
        ))}
      </div>
    );
  };

  // Component to render letter set (prefix or suffix)
  const LetterSet = ({
    title,
    letters,
    position
  }: {
    title: string,
    letters: string[],
    position: 'prefix' | 'suffix'
  }) => {
    const { vowels, consonants } = sortAndSplitLetters(letters);

    return (
      <div className={`w-1/3 p-4 bg-white ${position === 'prefix' ? 'rounded-l-lg' : 'rounded-r-lg'} shadow`}>
        <h2 className="text-lg font-semibold mb-2">{title}</h2>

        {letters.length > 0 ? (
          <div>
            <h3 className="text-md font-medium mt-2 mb-1">Vowels</h3>
            <LetterButtons letters={vowels} position={position} type="vowels" />

            <h3 className="text-md font-medium mt-2 mb-1">Consonants</h3>
            <LetterButtons letters={consonants} position={position} type="consonants" />
          </div>
        ) : (
          <p className="text-gray-500">No {title.toLowerCase()} letters available.</p>
        )}
      </div>
    );
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

  return (
    <div className="min-h-screen bg-gray-100 flex flex-col items-center justify-center p-4">
      <h1 className="text-3xl font-bold mb-6">Word Builder</h1>

      {state && (
        <div className="flex w-full max-w-4xl">
          {/* Prefix Letters */}
          <LetterSet
            title="Prefix Letters"
            letters={state.prefix_set}
            position="prefix"
          />

          {/* Current Word */}
          <div className="w-1/3 p-4 bg-white shadow flex flex-col items-center">
            <div className="flex items-center gap-2 mb-2">
              <h2 className="text-lg font-semibold">Current Word</h2>
              <motion.button
                onClick={resetWordBuilder}
                className="p-2 bg-gray-200 text-gray-700 rounded hover:bg-gray-300 flex items-center justify-center"
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
                disabled={state.answer.length === 0}
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"></path>
                  <path d="M3 3v5h5"></path>
                </svg>
              </motion.button>
            </div>
            <div className="flex gap-1 mb-4">
              {state.answer.split('').map((letter, index) => (
                <motion.span
                  key={`${letter}-${index}`}
                  onClick={() => removeLetter(index)}
                  className="px-2 py-1 bg-gray-200 rounded cursor-pointer hover:bg-gray-300"
                  initial={{ opacity: 0, y: -10 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ duration: 0.2 }}
                >
                  {letter}
                </motion.span>
              ))}
            </div>
            {state.is_valid_word && state.answer.length > 1 && wordDetails && (
              <motion.div
                className="mt-4 p-4 bg-gray-50 rounded-lg border border-gray-200"
                initial={{ opacity: 0, height: 0 }}
                animate={{ opacity: 1, height: "auto" }}
                transition={{ duration: 0.3 }}
              >
                <div className="flex items-center mb-2">
                  <h3 className="font-semibold mr-2">{state.answer}</h3>
                  <span className="text-gray-600">/{wordDetails.pronunciation}/</span>
                  {wordDetails.audio && (
                    <button
                      onClick={() => new Audio(wordDetails.audio).play()}
                      className="ml-2 text-blue-500"
                    >
                      🔊
                    </button>
                  )}
                </div>
                <p className="text-gray-800 mb-2">{wordDetails.meaning}</p>
                {wordDetails.example && (
                  <p className="text-gray-600 italic">"{wordDetails.example}"</p>
                )}
              </motion.div>
            )}
          </div>

          {/* Suffix Letters */}
          <LetterSet
            title="Suffix Letters"
            letters={state.suffix_set}
            position="suffix"
          />
        </div>
      )}

      {/* Error Message */}
      {error && (
        <motion.p
          className="mt-4 text-red-600"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.3 }}
        >
          {error}
        </motion.p>
      )}
    </div>
  );
};

export default HomePage;