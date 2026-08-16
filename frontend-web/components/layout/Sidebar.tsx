'use client';

import React from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Home, List, FileText, Heart, User } from 'lucide-react';

const menuItems = [
  { name: 'Beranda', href: '/', icon: Home },
  { name: 'Daftar Ruangan', href: '/rooms', icon: List },
  { name: 'Pesanan Saya', href: '/orders', icon: FileText },
  { name: 'Wishlist', href: '/wishlist', icon: Heart },
  { name: 'Akun', href: '/account', icon: User },
];

export const Sidebar = () => {
  const pathname = usePathname();

  return (
    <aside className="w-64 bg-[#f8f9fa] border-r border-gray-200 min-h-[calc(100vh-73px)] hidden md:block">
      <nav className="p-4 space-y-2">
        {menuItems.map((item) => {
          const isActive = pathname === item.href;
          const Icon = item.icon;
          return (
            <Link 
              key={item.name} 
              href={item.href}
              className={`flex items-center space-x-3 px-4 py-3 rounded-lg transition-colors ${
                isActive 
                  ? 'bg-white text-[#0055b8] shadow-sm font-semibold' 
                  : 'text-gray-700 hover:bg-white hover:text-[#0055b8]'
              }`}
            >
              <Icon className="w-5 h-5" />
              <span>{item.name}</span>
            </Link>
          );
        })}
      </nav>
    </aside>
  );
};
