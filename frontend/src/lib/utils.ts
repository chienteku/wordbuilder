/**
 * Check if a letter is a vowel
 */
export const isVowel = (letter: string): boolean => {
    return ['a', 'e', 'i', 'o', 'u'].includes(letter.toLowerCase());
};

/**
 * Sort and split letters into vowels and consonants
 */
export const sortAndSplitLetters = (letters: string[]) => {
    const vowels = letters.filter(letter => isVowel(letter)).sort();
    const consonants = letters.filter(letter => !isVowel(letter)).sort();
    return { vowels, consonants };
};