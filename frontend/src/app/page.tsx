"use client"
import { motion } from 'framer-motion';
import Link from 'next/link';

export default function HomePage() {
    return (
        <div className="min-h-screen bg-gray-100 p-4 flex flex-col items-center justify-center">
            <motion.div
                className="max-w-2xl w-full bg-white rounded-lg shadow-lg p-8"
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5 }}
            >
                <h1 className="text-3xl font-bold text-center mb-6">Word Builder Game</h1>

                <p className="text-gray-700 mb-6 text-center">
                    Challenge yourself with this word-building game! Add letters to
                    the beginning or end of your word to create valid English words.
                </p>

                <div className="mb-8 bg-blue-50 p-4 rounded-lg">
                    <h2 className="text-xl font-semibold mb-3 text-blue-800">How to Play:</h2>
                    <ul className="space-y-2 text-gray-700">
                        <li className="flex items-start">
                            <span className="mr-2 text-blue-500 font-bold">1.</span>
                            <span>Start by selecting a letter from either the prefix (left) or suffix (right) section.</span>
                        </li>
                        <li className="flex items-start">
                            <span className="mr-2 text-blue-500 font-bold">2.</span>
                            <span>Continue adding letters to form a valid word. The game will show you available options.</span>
                        </li>
                        <li className="flex items-start">
                            <span className="mr-2 text-blue-500 font-bold">3.</span>
                            <span>When you form a valid word, you'll see its definition and pronunciation.</span>
                        </li>
                        <li className="flex items-start">
                            <span className="mr-2 text-blue-500 font-bold">4.</span>
                            <span>Remove letters by clicking on them if you want to try a different path.</span>
                        </li>
                        <li className="flex items-start">
                            <span className="mr-2 text-blue-500 font-bold">5.</span>
                            <span>Try to form the longest word possible or discover new words!</span>
                        </li>
                    </ul>
                </div>

                <div className="flex justify-center">
                    <motion.div
                        whileHover={{ scale: 1.05 }}
                        whileTap={{ scale: 0.95 }}
                    >
                        <Link href="/game" className="inline-block bg-blue-600 hover:bg-blue-700 text-white font-bold py-3 px-6 rounded-lg transition-colors duration-200">
                            Start Playing
                        </Link>
                    </motion.div>
                </div>
            </motion.div>
        </div>
    );
}