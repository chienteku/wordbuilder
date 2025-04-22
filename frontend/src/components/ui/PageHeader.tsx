import Link from 'next/link';

interface PageHeaderProps {
    title: string;
    showHomeLink?: boolean;
}

export const PageHeader: React.FC<PageHeaderProps> = ({
    title,
    showHomeLink = false
}) => {
    return (
        <div className="flex justify-between items-center mb-6">
            {showHomeLink ? (
                <Link href="/" className="text-blue-600 hover:text-blue-800 flex items-center">
                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                        <path d="M19 12H5M12 19l-7-7 7-7" />
                    </svg>
                    <span className="ml-1">Home</span>
                </Link>
            ) : (
                <div className="w-20">{/* Spacer for alignment */}</div>
            )}

            <h1 className="text-2xl md:text-3xl font-bold text-center">{title}</h1>

            <div className="w-20"></div> {/* Spacer for alignment */}
        </div>
    );
};