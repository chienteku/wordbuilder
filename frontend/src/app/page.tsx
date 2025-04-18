"use client"
import { useState, useEffect } from 'react';
import axios from 'axios';
import { motion } from 'framer-motion';

interface WordBuilderState {
  answer: string;
  prefix_set: string[];
  suffix_set: string[];
  step: number;
}

interface ApiResponse {
  session_id?: string;
  state: WordBuilderState;
  success?: boolean;
  message?: string;
  error?: string;
}

// Update API URL to use the correct port (8081 instead of 8080)
// And add support for relative URLs when deployed
const API_BASE_URL = typeof window !== 'undefined' 
  ? (window.location.hostname === 'localhost' 
    ? 'http://localhost:8081/api/wordbuilder'
    : '/api/wordbuilder')
  : '/api/wordbuilder';

const HomePage: React.FC = () => {
  const [sessionId, setSessionId] = useState<string | null>(null);
  const [state, setState] = useState<WordBuilderState | null>(null);
  const [isValidWord, setIsValidWord] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // 初始化 WordBuilder
  const initializeWordBuilder = async () => {
    try {
      const response = await axios.post<ApiResponse>(`${API_BASE_URL}/init`, {});
      const newSessionId = response.data.session_id || null;
      setSessionId(newSessionId);
      setState(response.data.state);
      setIsValidWord(response.data.message?.includes('valid English word') || false);
      setError(null);
      if (newSessionId) {
        localStorage.setItem('sessionId', newSessionId);
      }
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to initialize WordBuilder.');
    }
  };

  // 添加字母
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
        setIsValidWord(response.data.message?.includes('valid English word') || false);
        setError(null);
      } else {
        setError(response.data.message || 'Failed to add letter.');
      }
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to add letter.');
    }
  };

  // 移除字母
  const removeLetter = async (index: number) => {
    if (!sessionId || !state) return;
    try {
      const response = await axios.post<ApiResponse>(`${API_BASE_URL}/remove`, {
        session_id: sessionId,
        index,
      });
      if (response.data.success) {
        setState(response.data.state);
        setIsValidWord(response.data.message?.includes('valid English word') || false);
        setError(null);
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

  // This function is no longer needed as we can click on any letter to remove it

  // 恢復狀態
  useEffect(() => {
    const savedSessionId = localStorage.getItem('sessionId');
    if (savedSessionId && !sessionId) {
      axios
        .get<ApiResponse>(`${API_BASE_URL}/state?session_id=${savedSessionId}`)
        .then((response) => {
          setSessionId(savedSessionId);
          setState(response.data.state);
          setIsValidWord(response.data.message?.includes('valid English word') || false);
        })
        .catch((err) => {
          setError(`Failed to restore session. Starting a new session.`);
          localStorage.removeItem('sessionId');
          initializeWordBuilder();
        });
    } else if (!sessionId) {
      // If no saved session, start a new one
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

  return (
    <div className="min-h-screen bg-gray-100 flex flex-col items-center justify-center p-4">
      <h1 className="text-3xl font-bold mb-6">Word Builder</h1>

      {/* 主畫面 */}
      {state && (
        <div className="flex w-full max-w-4xl">
          {/* 前綴字母集合 */}
          <div className="w-1/3 p-4 bg-white rounded-l-lg shadow">
            <h2 className="text-lg font-semibold mb-2">Prefix Letters</h2>

            {state.prefix_set.length > 0 ? (
              <div>
                {/* Vowels */}
                <h3 className="text-md font-medium mt-2 mb-1">Vowels</h3>
                <div className="flex flex-wrap gap-2 mb-3">
                  {sortAndSplitLetters(state.prefix_set).vowels.length > 0 ? (
                    sortAndSplitLetters(state.prefix_set).vowels.map((letter) => (
                      <motion.button
                        key={`vowel-${letter}`}
                        onClick={() => addLetter(letter, 'prefix')}
                        className="px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600"
                        initial={{ opacity: 0, scale: 0.8 }}
                        animate={{ opacity: 1, scale: 1 }}
                        transition={{ duration: 0.2 }}
                      >
                        {letter}
                      </motion.button>
                    ))
                  ) : (
                    <p className="text-gray-500">No vowels available.</p>
                  )}
                </div>

                {/* Consonants */}
                <h3 className="text-md font-medium mt-2 mb-1">Consonants</h3>
                <div className="flex flex-wrap gap-2">
                  {sortAndSplitLetters(state.prefix_set).consonants.length > 0 ? (
                    sortAndSplitLetters(state.prefix_set).consonants.map((letter) => (
                      <motion.button
                        key={`consonant-${letter}`}
                        onClick={() => addLetter(letter, 'prefix')}
                        className="px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600"
                        initial={{ opacity: 0, scale: 0.8 }}
                        animate={{ opacity: 1, scale: 1 }}
                        transition={{ duration: 0.2 }}
                      >
                        {letter}
                      </motion.button>
                    ))
                  ) : (
                    <p className="text-gray-500">No consonants available.</p>
                  )}
                </div>
              </div>
            ) : (
              <p className="text-gray-500">No prefix letters available.</p>
            )}
          </div>

          {/* 答案和提示 */}
          <div className="w-1/3 p-4 bg-white shadow flex flex-col items-center">
            <h2 className="text-lg font-semibold mb-2">Current Word</h2>
            <div className="flex gap-1 mb-4">
              {state.answer.split('').map((letter, index) => (
                <motion.span
                  key={`${letter}-${index}`}
                  onClick={() => removeLetter(index)}
                  className="px-2 py-1 bg-gray-200 rounded cursor-pointer hover:bg-gray-300"
                  initial={{ opacity: 0, y: -10 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: 10 }}
                  transition={{ duration: 0.2 }}
                >
                  {letter}
                </motion.span>
              ))}
            </div>

            {isValidWord && (
              <motion.p
                className="text-green-600 font-semibold mt-4"
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ duration: 0.3 }}
              >
                This is a valid English word!
              </motion.p>
            )}
          </div>

          {/* 後綴字母集合 */}
          <div className="w-1/3 p-4 bg-white rounded-r-lg shadow">
            <h2 className="text-lg font-semibold mb-2">Suffix Letters</h2>

            {state.suffix_set.length > 0 ? (
              <div>
                {/* Vowels */}
                <h3 className="text-md font-medium mt-2 mb-1">Vowels</h3>
                <div className="flex flex-wrap gap-2 mb-3">
                  {sortAndSplitLetters(state.suffix_set).vowels.length > 0 ? (
                    sortAndSplitLetters(state.suffix_set).vowels.map((letter) => (
                      <motion.button
                        key={`vowel-${letter}`}
                        onClick={() => addLetter(letter, 'suffix')}
                        className="px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600"
                        initial={{ opacity: 0, scale: 0.8 }}
                        animate={{ opacity: 1, scale: 1 }}
                        transition={{ duration: 0.2 }}
                      >
                        {letter}
                      </motion.button>
                    ))
                  ) : (
                    <p className="text-gray-500">No vowels available.</p>
                  )}
                </div>

                {/* Consonants */}
                <h3 className="text-md font-medium mt-2 mb-1">Consonants</h3>
                <div className="flex flex-wrap gap-2">
                  {sortAndSplitLetters(state.suffix_set).consonants.length > 0 ? (
                    sortAndSplitLetters(state.suffix_set).consonants.map((letter) => (
                      <motion.button
                        key={`consonant-${letter}`}
                        onClick={() => addLetter(letter, 'suffix')}
                        className="px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600"
                        initial={{ opacity: 0, scale: 0.8 }}
                        animate={{ opacity: 1, scale: 1 }}
                        transition={{ duration: 0.2 }}
                      >
                        {letter}
                      </motion.button>
                    ))
                  ) : (
                    <p className="text-gray-500">No consonants available.</p>
                  )}
                </div>
              </div>
            ) : (
              <p className="text-gray-500">No suffix letters available.</p>
            )}
          </div>
        </div>
      )}

      {/* 錯誤訊息 */}
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