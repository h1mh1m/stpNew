import Link from 'next/link';
import { ChevronDown, LogIn } from 'lucide-react';

export default function Navbar() {
  return (
    <nav className="bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between sticky top-0 z-50">
      <div className="flex items-center space-x-8">
        <Link href="/" className="text-2xl font-bold text-[#001021]">
          <span className="text-[#0055b8]">Wira</span>nala
        </Link>
      </div>
      <div className="flex items-center space-x-6">
        <button className="flex items-center text-sm font-medium text-gray-700 hover:text-gray-900">
          ID <ChevronDown className="ml-1 w-4 h-4" />
        </button>
        <Link href="/signin" className="flex items-center text-sm font-semibold text-gray-800 hover:text-[#0055b8] transition-colors">
          <LogIn className="w-5 h-5 mr-2" />
          Masuk
        </Link>
      </div>
    </nav>
  );
}
