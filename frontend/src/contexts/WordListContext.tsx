import React, { createContext, useContext, useEffect, useState, ReactNode } from 'react';
import { wordListService } from '@/services/wordlist-service';

type WordList = {
  id: number;
  name: string;
  // ...other fields as needed
};

type WordListContextType = {
  wordLists: WordList[];
  loading: boolean;
  error: string | null;
  refresh: () => Promise<void>;
};

const WordListContext = createContext<WordListContextType | undefined>(undefined);

export const WordListProvider = ({ children }: { children: ReactNode }) => {
  const [wordLists, setWordLists] = useState<WordList[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchWordLists = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await wordListService.getAllWordLists();
      setWordLists(response.word_lists || []);
    } catch (err: any) {
      setError('Failed to fetch word lists');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchWordLists();
  }, []);

  return (
    <WordListContext.Provider value={{ wordLists, loading, error, refresh: fetchWordLists }}>
      {children}
    </WordListContext.Provider>
  );
};

export const useWordListContext = () => {
  const context = useContext(WordListContext);
  if (!context) throw new Error('useWordListContext must be used within a WordListProvider');
  return context;
};