import { motion } from 'framer-motion';

interface ErrorMessageProps {
    message: string;
}

export const ErrorMessage: React.FC<ErrorMessageProps> = ({ message }) => {
    if (!message) return null;

    return (
        <motion.div
            className="mt-4 p-3 bg-red-100 text-red-700 rounded-lg"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 0.3 }}
        >
            <p className="text-sm text-center">{message}</p>
        </motion.div>
    );
};