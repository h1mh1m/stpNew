'use client';

import React from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Home, ShoppingCart, User } from 'lucide-react';

const menuItems = [
  { name: 'Beranda', href: '/', icon: Home },
  { name: 'Pemesanan', href: '/booking', icon: ShoppingCart },
  { name: 'Profil', href: '/profile', icon: User },
];

export const Sidebar = () => {
  const pathname = usePathname();

  return (
    <aside className="w-64 bg-white border-r border-gray-200 flex flex-col min-h-screen">
      <div className="p-6">
        <h1 className="text-2xl font-bold italic tracking-tighter">
          TRACKIQ
        </h1>
      </div>
      <nav className="px-4 space-y-2 flex-1 mt-4">
        {menuItems.map((item) => {
          const isActive = pathname === item.href || (item.href !== '/' && pathname.startsWith(item.href));
          const Icon = item.icon;
          return (
            <Link 
              key={item.name} 
              href={item.href}
              className={`flex items-center space-x-3 px-4 py-3 rounded-lg transition-colors ${
                isActive 
                  ? 'bg-[#1e1b4b] text-white font-medium' // Dark blue background for active
                  : 'text-gray-900 hover:bg-gray-100'
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
