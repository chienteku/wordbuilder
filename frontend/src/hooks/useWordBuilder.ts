import { useState, useEffect } from 'react';
import { WordBuilderState, WordDetails } from '@/types/wordbuilder';
import { wordBuilderService, dictionaryService } from '@/services/wordbuilder-service';

export const useWordBuilder = () => {
    const [sessionId, setSessionId] = useState<string | null>(null);
    const [state, setState] = useState<WordBuilderState | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [wordDetails, setWordDetails] = useState<WordDetails | null>(null);
    const [loading, setLoading] = useState<boolean>(false);
    const [isInitialized, setIsInitialized] = useState<boolean>(false);

    // Initialize WordBuilder
    const initializeWordBuilder = async () => {
        try {
            setLoading(true);
            const response = await wordBuilderService.initSession();
            const newSessionId = response.session_id || null;
            setSessionId(newSessionId);
            setState(response.state);
            setError(null);
            if (newSessionId) {
                localStorage.setItem('sessionId', newSessionId);
            }
        } catch (err: any) {
            setError(err.response?.data?.error || 'Failed to initialize WordBuilder.');
        } finally {
            setLoading(false);
        }
    };

    // Add letter
    const addLetter = async (letter: string, position: 'prefix' | 'suffix') => {
        if (!sessionId || !state) return;
        try {
            setLoading(true);
            const response = await wordBuilderService.addLetter(sessionId, letter, position);
            if (response.success) {
                setState(response.state);
                setError(null);

                // Only fetch word details if it's valid AND has more than one letter
                if (response.state.is_valid_word && response.state.answer.length > 1) {
                    fetchWordDetails(response.state.answer);
                } else {
                    setWordDetails(null);
                }
            } else {
                setError(response.message || 'Failed to add letter.');
            }
        } catch (err: any) {
            setError(err.response?.data?.error || 'Failed to add letter.');
        } finally {
            setLoading(false);
        }
    };

    // Remove letter
    const removeLetter = async (index: number) => {
        if (!sessionId || !state) return;
        try {
            setLoading(true);
            const response = await wordBuilderService.removeLetter(sessionId, index);
            if (response.success) {
                setState(response.state);
                setError(null);

                // Only fetch word details if it's valid AND has more than one letter
                if (response.state.is_valid_word && response.state.answer.length > 1) {
                    fetchWordDetails(response.state.answer);
                } else {
                    setWordDetails(null);
                }
            } else {
                setError(response.message || 'Failed to remove letter.');
            }
        } catch (err: any) {
            setError(err.response?.data?.error || 'Failed to remove letter.');
            if (err.response?.status === 404) {
                localStorage.removeItem('sessionId');
                setSessionId(null);
                setState(null);
            }
        } finally {
            setLoading(false);
        }
    };

    // Reset WordBuilder
    const resetWordBuilder = async () => {
        if (!sessionId) {
            // If no session, create a new one
            initializeWordBuilder();
            return;
        }

        try {
            setLoading(true);
            const response = await wordBuilderService.resetSession(sessionId);
            if (response.success) {
                setState(response.state);
                setWordDetails(null);
                setError(null);
            } else {
                setError(response.message || 'Failed to reset word builder.');
            }
        } catch (err: any) {
            setError(err.response?.data?.error || 'Failed to reset word builder.');
            if (err.response?.status === 404) {
                localStorage.removeItem('sessionId');
                setSessionId(null);
                setState(null);
                initializeWordBuilder();
            }
        } finally {
            setLoading(false);
        }
    };

    // Fetch word details from dictionary API
    const fetchWordDetails = async (word: string) => {
        try {
            const details = await dictionaryService.getWordDetails(word);
            setWordDetails(details);
        } catch (err) {
            console.error('Failed to fetch word details:', err);
        }
    };

    // Restore state from localStorage or initialize new session
    useEffect(() => {
        const savedSessionId = localStorage.getItem('sessionId');
        if (savedSessionId && !sessionId) {
            wordBuilderService.getState(savedSessionId)
                .then((response) => {
                    setSessionId(savedSessionId);
                    setState(response.state);

                    // If the word is valid, fetch its details
                    if (response.state.is_valid_word && response.state.answer.length > 1) {
                        fetchWordDetails(response.state.answer);
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

    // Update the useEffect to track initialization
    useEffect(() => {
        const initSession = async () => {
            const savedSessionId = localStorage.getItem('sessionId');
            if (savedSessionId) {
                try {
                    const response = await wordBuilderService.getState(savedSessionId);
                    setSessionId(savedSessionId);
                    setState(response.state);

                    // If the word is valid, fetch its details
                    if (response.state.is_valid_word && response.state.answer.length > 1) {
                        fetchWordDetails(response.state.answer);
                    }
                    setIsInitialized(true);
                } catch (err) {
                    console.error("Failed to restore session:", err);
                    localStorage.removeItem('sessionId');
                    await createNewSession();
                }
            } else {
                await createNewSession();
            }
        };

        const createNewSession = async () => {
            try {
                await initializeWordBuilder();
                setIsInitialized(true);
            } catch (err) {
                console.error("Failed to create new session:", err);
            }
        };

        if (!isInitialized) {
            initSession();
        }
    }, [isInitialized]);

    return {
        state,
        error,
        wordDetails,
        loading,
        addLetter,
        removeLetter,
        resetWordBuilder,
        initializeWordBuilder
    };
};