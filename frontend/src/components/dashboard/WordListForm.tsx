import { useState, useRef, ChangeEvent, FormEvent } from 'react';
import { motion } from 'framer-motion';
import { WordList } from '@/types/wordlist';

interface WordListFormProps {
    onSubmit: (formData: FormData) => Promise<void>;
    onCancel: () => void;
    initialData?: WordList;
    submitLabel: string;
}

export const WordListForm: React.FC<WordListFormProps> = ({
    onSubmit,
    onCancel,
    initialData,
    submitLabel
}) => {
    const [name, setName] = useState(initialData?.name || '');
    const [description, setDescription] = useState(initialData?.description || '');
    const [source, setSource] = useState(initialData?.source || '');
    const [fileName, setFileName] = useState('');
    const [error, setError] = useState('');
    const [isLoading, setIsLoading] = useState(false);

    const fileInputRef = useRef<HTMLInputElement>(null);

    const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
        if (e.target.files && e.target.files.length > 0) {
            const file = e.target.files[0];

            // Validate file type
            if (!file.name.endsWith('.txt')) {
                setError('Only .txt files are allowed');
                e.target.value = '';
                return;
            }

            // Check file size (max 10MB)
            if (file.size > 10 * 1024 * 1024) {
                setError('File size exceeds 10MB limit');
                e.target.value = '';
                return;
            }

            setFileName(file.name);
            setError('');
        }
    };

    const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault();

        if (!name.trim()) {
            setError('Name is required');
            return;
        }

        // If updating and no file is selected, don't require a file
        if (!initialData && (!fileInputRef.current?.files || fileInputRef.current.files.length === 0)) {
            setError('Please select a word list file');
            return;
        }

        setIsLoading(true);
        setError('');

        try {
            const formData = new FormData();
            formData.append('name', name);
            formData.append('description', description);
            formData.append('source', source);

            if (fileInputRef.current?.files && fileInputRef.current.files.length > 0) {
                formData.append('file', fileInputRef.current.files[0]);
            }

            await onSubmit(formData);
        } catch (err: any) {
            setError(err.response?.data?.error || 'Failed to submit form');
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="bg-white p-6 rounded-lg shadow-md"
        >
            <h2 className="text-xl font-semibold mb-4">
                {initialData ? 'Edit Word List' : 'Create New Word List'}
            </h2>

            <form onSubmit={handleSubmit}>
                <div className="mb-4">
                    <label htmlFor="name" className="block text-sm font-medium text-gray-700 mb-1">
                        Name *
                    </label>
                    <input
                        id="name"
                        type="text"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                        required
                    />
                </div>

                <div className="mb-4">
                    <label htmlFor="description" className="block text-sm font-medium text-gray-700 mb-1">
                        Description
                    </label>
                    <textarea
                        id="description"
                        value={description}
                        onChange={(e) => setDescription(e.target.value)}
                        className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 min-h-[100px]"
                    />
                </div>

                <div className="mb-4">
                    <label htmlFor="source" className="block text-sm font-medium text-gray-700 mb-1">
                        Source
                    </label>
                    <input
                        id="source"
                        type="text"
                        value={source}
                        onChange={(e) => setSource(e.target.value)}
                        className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                        placeholder="Where did this word list come from?"
                    />
                </div>

                <div className="mb-6">
                    <label htmlFor="file" className="block text-sm font-medium text-gray-700 mb-1">
                        Word List File {!initialData && '*'}
                    </label>
                    <input
                        id="file"
                        type="file"
                        ref={fileInputRef}
                        onChange={handleFileChange}
                        className="hidden"
                        accept=".txt"
                    />
                    <div className="flex items-center gap-2">
                        <button
                            type="button"
                            onClick={() => fileInputRef.current?.click()}
                            className="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 transition-colors"
                        >
                            Choose File
                        </button>
                        <span className="text-sm text-gray-500">
                            {fileName || (initialData ? 'Keep existing file' : 'No file selected')}
                        </span>
                    </div>
                    <p className="mt-1 text-xs text-gray-500">
                        Only .txt files up to 10MB are accepted
                    </p>
                </div>

                {error && (
                    <div className="mb-4 p-3 bg-red-100 text-red-700 rounded-md">
                        {error}
                    </div>
                )}

                <div className="flex justify-end gap-3">
                    <button
                        type="button"
                        onClick={onCancel}
                        className="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 transition-colors"
                        disabled={isLoading}
                    >
                        Cancel
                    </button>
                    <button
                        type="submit"
                        className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors disabled:bg-blue-400"
                        disabled={isLoading}
                    >
                        {isLoading ? 'Submitting...' : submitLabel}
                    </button>
                </div>
            </form>
        </motion.div>
    );
};