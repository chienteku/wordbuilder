import { motion } from 'framer-motion';

interface ConfirmModalProps {
    title: string;
    message: string;
    confirmLabel: string;
    cancelLabel?: string;
    confirmVariant?: 'danger' | 'warning' | 'primary';
    onConfirm: () => void;
    onCancel: () => void;
    isLoading?: boolean;
}

export const ConfirmModal: React.FC<ConfirmModalProps> = ({
    title,
    message,
    confirmLabel,
    cancelLabel = 'Cancel',
    confirmVariant = 'danger',
    onConfirm,
    onCancel,
    isLoading = false,
}) => {
    const getButtonClass = () => {
        switch (confirmVariant) {
            case 'danger':
                return 'bg-red-600 hover:bg-red-700 text-white';
            case 'warning':
                return 'bg-yellow-600 hover:bg-yellow-700 text-white';
            case 'primary':
            default:
                return 'bg-blue-600 hover:bg-blue-700 text-white';
        }
    };

    return (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
            <motion.div
                initial={{ opacity: 0, scale: 0.9 }}
                animate={{ opacity: 1, scale: 1 }}
                exit={{ opacity: 0, scale: 0.9 }}
                className="bg-white rounded-lg shadow-lg max-w-md w-full"
            >
                <div className="p-4 border-b">
                    <h2 className="text-xl font-semibold">{title}</h2>
                </div>

                <div className="p-4">
                    <p className="text-gray-700">{message}</p>
                </div>

                <div className="p-4 border-t flex justify-end space-x-2">
                    <button
                        onClick={onCancel}
                        disabled={isLoading}
                        className="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 transition-colors disabled:bg-gray-100 disabled:text-gray-500"
                    >
                        {cancelLabel}
                    </button>
                    <button
                        onClick={onConfirm}
                        disabled={isLoading}
                        className={`px-4 py-2 rounded-md transition-colors disabled:opacity-70 ${getButtonClass()}`}
                    >
                        {isLoading ? 'Processing...' : confirmLabel}
                    </button>
                </div>
            </motion.div>
        </div>
    );
};