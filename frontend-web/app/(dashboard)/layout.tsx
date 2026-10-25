'use client';

import { Sidebar } from '@/components/layout/Sidebar';
import { Footer } from '@/components/layout/Footer';
import { LogOut } from 'lucide-react';
import { useRouter } from 'next/navigation';
import { useEffect } from 'react';

import { getAuthToken, clearAuthToken } from '@/lib/auth';

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode
}) {
  const router = useRouter();

  useEffect(() => {
    const token = getAuthToken();
    if (!token) {
      router.push('/signin');
    }
  }, [router]);

  const handleLogout = () => {
    clearAuthToken();
    router.push('/signin');
  };

  return (
    <div className="min-h-screen bg-[#eef2f6] font-sans flex flex-col md:flex-row">
      <Sidebar />
      <div className="flex-1 flex flex-col min-w-0">
        <header className="flex justify-end p-4 bg-[#eef2f6]">
          <button 
            onClick={handleLogout}
            className="flex items-center space-x-2 bg-[#1e1b4b] text-white px-5 py-2 rounded-lg font-medium shadow-sm transition-opacity hover:opacity-90"
          >
            <LogOut className="w-4 h-4" />
            <span>Log out</span>
          </button>
        </header>
        <main className="flex-1 p-4 lg:p-8 overflow-x-hidden">
          {children}
        </main>
        <Footer />
      </div>
    </div>
  );
}
