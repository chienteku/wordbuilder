"use client"
import { useState } from 'react';
import { motion } from 'framer-motion';
import Link from 'next/link';
import { wordListService } from '@/services/wordlist-service';
import { WordList } from '@/types/wordlist';
import { PageHeader } from '@/components/ui/PageHeader';
import { WordListCard } from '@/components/dashboard/WordListCard';
import { WordListForm } from '@/components/dashboard/WordListForm';
import { WordSampleModal } from '@/components/dashboard/WordSampleModal';
import { ConfirmModal } from '@/components/ui/ConfirmModal';
import { useWordListContext } from '@/contexts/WordListContext';

export default function DashboardPage() {
    // Use the context instead of local state for wordLists
    const { wordLists, loading: isLoading, error: contextError, refresh } = useWordListContext();

    const [error, setError] = useState('');
    const [showCreateForm, setShowCreateForm] = useState(false);
    const [editingWordList, setEditingWordList] = useState<WordList | null>(null);
    const [activeWordListId, setActiveWordListId] = useState<number | null>(null);
    const [viewingSample, setViewingSample] = useState<WordList | null>(null);
    const [confirmDeleteWordList, setConfirmDeleteWordList] = useState<WordList | null>(null);
    const [isProcessing, setIsProcessing] = useState(false);
    const [notification, setNotification] = useState({ message: '', type: '' });

    // Set active word list when wordLists change
    if (wordLists.length > 0 && activeWordListId === null) {
        setActiveWordListId(wordLists[0].id);
    }

    const handleCreateWordList = async (formData: FormData) => {
        try {
            setIsProcessing(true);
            await wordListService.createWordList(formData);
            setShowCreateForm(false);
            showNotification('Word list created successfully', 'success');
            // Use the refresh function from context instead of local fetchWordLists
            await refresh();
        } catch (err: any) {
            setError(err.response?.data?.error || 'Failed to create word list');
            showNotification('Failed to create word list', 'error');
        } finally {
            setIsProcessing(false);
        }
    };

    const handleUpdateWordList = async (formData: FormData) => {
        if (!editingWordList) return;

        try {
            setIsProcessing(true);
            await wordListService.updateWordList(editingWordList.id, formData);
            setEditingWordList(null);
            showNotification('Word list updated successfully', 'success');
            await refresh();
        } catch (err: any) {
            setError(err.response?.data?.error || 'Failed to update word list');
            showNotification('Failed to update word list', 'error');
        } finally {
            setIsProcessing(false);
        }
    };

    const handleDeleteWordList = async () => {
        if (!confirmDeleteWordList) return;

        try {
            setIsProcessing(true);
            await wordListService.deleteWordList(confirmDeleteWordList.id);
            setConfirmDeleteWordList(null);
            showNotification('Word list deleted successfully', 'success');
            await refresh();
        } catch (err: any) {
            setError(err.response?.data?.error || 'Failed to delete word list');
            showNotification('Failed to delete word list', 'error');
        } finally {
            setIsProcessing(false);
        }
    };

    const handleUseWordList = async (wordList: WordList) => {
        try {
            setIsProcessing(true);
            await wordListService.useWordList(wordList.id);
            setActiveWordListId(wordList.id);
            showNotification(`Now using "${wordList.name}" as the active dictionary`, 'success');
        } catch (err: any) {
            setError(err.response?.data?.error || 'Failed to set active word list');
            showNotification('Failed to set active word list', 'error');
        } finally {
            setIsProcessing(false);
        }
    };

    const showNotification = (message: string, type: 'success' | 'error') => {
        setNotification({ message, type });
        setTimeout(() => {
            setNotification({ message: '', type: '' });
        }, 5000);
    };

    return (
        <div className="min-h-screen bg-gray-100 p-4">
            <div className="max-w-6xl mx-auto">
                <PageHeader title="Word List Dashboard" showHomeLink={true} />

                <div className="mb-6 flex justify-between items-center">
                    <div>
                        <h2 className="text-xl font-semibold">Manage Word Lists</h2>
                        <p className="text-gray-600 text-sm">
                            Create, edit, and manage your word lists for the Word Builder game
                        </p>
                    </div>

                    <div className="flex gap-3">
                        <Link href="/game" className="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition-colors flex items-center">
                            <span>Play Game</span>
                        </Link>
                        <button
                            onClick={() => setShowCreateForm(true)}
                            className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
                        >
                            New Word List
                        </button>
                    </div>
                </div>

                {notification.message && (
                    <motion.div
                        initial={{ opacity: 0, y: -20 }}
                        animate={{ opacity: 1, y: 0 }}
                        exit={{ opacity: 0, y: -20 }}
                        className={`mb-4 p-3 rounded-md ${notification.type === 'success' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'
                            }`}
                    >
                        {notification.message}
                    </motion.div>
                )}

                {(error || contextError) && (
                    <div className="mb-4 p-3 bg-red-100 text-red-700 rounded-md">
                        {error || contextError}
                    </div>
                )}

                {isLoading ? (
                    <div className="bg-white rounded-lg shadow-md p-8 text-center">
                        <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mb-4"></div>
                        <p>Loading word lists...</p>
                    </div>
                ) : showCreateForm ? (
                    <WordListForm
                        onSubmit={handleCreateWordList}
                        onCancel={() => setShowCreateForm(false)}
                        submitLabel="Create Word List"
                    />
                ) : editingWordList ? (
                    <WordListForm
                        onSubmit={handleUpdateWordList}
                        onCancel={() => setEditingWordList(null)}
                        initialData={editingWordList}
                        submitLabel="Update Word List"
                    />
                ) : wordLists.length === 0 ? (
                    <div className="bg-white rounded-lg shadow-md p-8 text-center">
                        <p className="text-gray-600 mb-4">No word lists found. Create your first word list to get started.</p>
                        <button
                            onClick={() => setShowCreateForm(true)}
                            className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
                        >
                            Create Word List
                        </button>
                    </div>
                ) : (
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        {wordLists.map((wordList) => (
                            <WordListCard
                                key={wordList.id}
                                wordList={wordList}
                                isActive={wordList.id === activeWordListId}
                                onEdit={() => setEditingWordList(wordList)}
                                onDelete={() => setConfirmDeleteWordList(wordList)}
                                onDownload={() => wordListService.downloadWordList(wordList.id)}
                                onUse={() => handleUseWordList(wordList)}
                                onViewSample={() => setViewingSample(wordList)}
                            />
                        ))}
                    </div>
                )}
            </div>

            {viewingSample && (
                <WordSampleModal
                    wordList={viewingSample}
                    onClose={() => setViewingSample(null)}
                />
            )}

            {confirmDeleteWordList && (
                <ConfirmModal
                    title="Delete Word List"
                    message={`Are you sure you want to delete "${confirmDeleteWordList.name}"? This action cannot be undone.`}
                    confirmLabel="Delete"
                    confirmVariant="danger"
                    onConfirm={handleDeleteWordList}
                    onCancel={() => setConfirmDeleteWordList(null)}
                    isLoading={isProcessing}
                />
            )}
        </div>
    );
}