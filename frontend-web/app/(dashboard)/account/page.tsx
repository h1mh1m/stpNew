'use client';

import React from 'react';
import { User, Mail, Phone, Lock, Calendar, CreditCard, ChevronRight, LayoutGrid, Edit2, MoreHorizontal } from 'lucide-react';

export default function AccountPage() {
  return (
    <div className="max-w-4xl mx-auto pb-12">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Akun</h1>
      </div>

      <div className="space-y-6">
        {/* Main Profile Card */}
        <div className="bg-white rounded-2xl border border-gray-100 shadow-sm p-6 flex flex-col md:flex-row md:items-center gap-6">
          <div className="flex items-center gap-4 flex-1">
            <div className="relative">
              <div className="w-20 h-20 bg-[#512BD4] rounded-full flex items-center justify-center text-white shadow-inner">
                <User size={40} className="text-white/80" />
              </div>
              <button className="absolute bottom-0 right-0 bg-gray-100 p-1.5 rounded-full border-2 border-white text-gray-600 hover:bg-gray-200 transition-colors">
                <Edit2 size={14} />
              </button>
            </div>
            <div>
              <h2 className="text-xl font-bold text-gray-900">Nama</h2>
              <p className="text-sm text-gray-500">Status</p>
            </div>
          </div>
          
          <div className="flex-1 grid grid-cols-1 gap-4">
            <div className="flex items-center gap-3">
              <div className="w-8 h-8 rounded-full bg-gray-50 flex items-center justify-center text-gray-500">
                <CreditCard size={16} />
              </div>
              <div>
                <p className="text-sm font-semibold text-gray-900">0123456789</p>
                <p className="text-xs text-gray-500">ID</p>
              </div>
            </div>
            <div className="flex items-center gap-3">
              <div className="w-8 h-8 rounded-full bg-gray-50 flex items-center justify-center text-gray-500">
                <Phone size={16} />
              </div>
              <div>
                <p className="text-sm font-semibold text-gray-900">Belum diatur</p>
                <div className="flex items-center gap-1.5">
                  <p className="text-xs text-gray-500">Nomor Ponsel</p>
                  <div className="w-1.5 h-1.5 rounded-full bg-orange-500"></div>
                </div>
              </div>
            </div>
          </div>

          <div className="flex-1 flex items-start gap-3">
            <div className="w-8 h-8 rounded-full bg-gray-50 flex items-center justify-center text-gray-500">
              <Mail size={16} />
            </div>
            <div>
              <p className="text-sm font-semibold text-gray-900">Test@Test.com</p>
              <div className="flex items-center gap-1.5">
                <p className="text-xs text-gray-500">Email Utama</p>
                <div className="w-1.5 h-1.5 rounded-full bg-green-500"></div>
              </div>
            </div>
          </div>
        </div>

        {/* Microsoft Card */}
        <div className="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
          <div className="px-6 py-4 border-b border-gray-50 flex items-center gap-2">
            <LayoutGrid size={18} className="text-[#512BD4]" />
            <h3 className="font-semibold text-gray-900 text-sm">Microsoft</h3>
          </div>
          <div className="p-2">
            <div className="flex items-center justify-between p-4 hover:bg-gray-50 rounded-xl cursor-pointer transition-colors group">
              <div className="flex items-start gap-4">
                <div className="mt-0.5 text-gray-500"><User size={20} /></div>
                <div>
                  <p className="text-sm font-semibold text-gray-900 group-hover:text-[#512BD4] transition-colors">Pusat Akun Microsoft</p>
                  <p className="text-xs text-gray-500">Ubah foto profil, multi-factor authentication, dan lainnya</p>
                </div>
              </div>
              <ChevronRight size={20} className="text-gray-400" />
            </div>
          </div>
        </div>

        {/* Informasi Pribadi Card */}
        <div className="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
          <div className="px-6 py-5 border-b border-gray-50">
            <h3 className="font-semibold text-gray-900 text-base">Informasi pribadi</h3>
            <p className="text-sm text-gray-500 mt-1">Data yang umumnya ditampilkan di semua web myITS</p>
          </div>
          <div className="p-2 space-y-1">
            <div className="flex items-center justify-between p-4 hover:bg-gray-50 rounded-xl cursor-pointer transition-colors group">
              <div className="flex items-start gap-4">
                <div className="mt-0.5 text-gray-500"><User size={20} /></div>
                <div>
                  <p className="text-sm font-semibold text-gray-900 group-hover:text-[#512BD4] transition-colors">Nama</p>
                  <p className="text-xs text-gray-500">Ubah nama panggilan</p>
                </div>
              </div>
              <ChevronRight size={20} className="text-gray-400" />
            </div>
            
            <div className="flex items-center justify-between p-4 hover:bg-gray-50 rounded-xl cursor-pointer transition-colors group">
              <div className="flex items-start gap-4">
                <div className="mt-0.5 text-gray-500"><Calendar size={20} /></div>
                <div>
                  <p className="text-sm font-semibold text-gray-900 group-hover:text-[#512BD4] transition-colors">Tanggal Lahir</p>
                  <p className="text-xs text-gray-500">Melihat tanggal lahir</p>
                </div>
              </div>
              <ChevronRight size={20} className="text-gray-400" />
            </div>
          </div>
        </div>

        {/* Kontak Card */}
        <div className="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
          <div className="px-6 py-5 border-b border-gray-50">
            <h3 className="font-semibold text-gray-900 text-base">Kontak</h3>
            <p className="text-sm text-gray-500 mt-1">Ubah email dan nomor ponsel</p>
          </div>
          <div className="p-2 space-y-1">
            <div className="flex items-center justify-between p-4 hover:bg-gray-50 rounded-xl cursor-pointer transition-colors group">
              <div className="flex items-start gap-4">
                <div className="mt-0.5 text-gray-500"><Mail size={20} /></div>
                <div>
                  <p className="text-sm font-semibold text-gray-900 group-hover:text-[#512BD4] transition-colors">Email</p>
                  <p className="text-xs text-gray-500">Perbarui dan verifikasi email</p>
                </div>
              </div>
              <ChevronRight size={20} className="text-gray-400" />
            </div>
            
            <div className="flex items-center justify-between p-4 hover:bg-gray-50 rounded-xl cursor-pointer transition-colors group">
              <div className="flex items-start gap-4">
                <div className="mt-0.5 text-gray-500"><Phone size={20} /></div>
                <div>
                  <p className="text-sm font-semibold text-gray-900 group-hover:text-[#512BD4] transition-colors">Nomor Ponsel</p>
                  <p className="text-xs text-gray-500">Perbarui dan verifikasi nomor telepon</p>
                </div>
              </div>
              <ChevronRight size={20} className="text-gray-400" />
            </div>
          </div>
        </div>

        {/* Lainnya Card */}
        <div className="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
          <div className="px-6 py-5 border-b border-gray-50">
            <h3 className="font-semibold text-gray-900 text-base">Lainnya</h3>
            <p className="text-sm text-gray-500 mt-1">Kata sandi dan preferensi tampilan</p>
          </div>
          <div className="p-2 space-y-1">
            <div className="flex items-center justify-between p-4 hover:bg-gray-50 rounded-xl cursor-pointer transition-colors group">
              <div className="flex items-start gap-4">
                <div className="mt-0.5 text-gray-500"><MoreHorizontal size={20} /></div>
                <div>
                  <p className="text-sm font-semibold text-gray-900 group-hover:text-[#512BD4] transition-colors">Kata Sandi</p>
                  <p className="text-xs text-gray-500">Perbarui kata sandi</p>
                </div>
              </div>
              <ChevronRight size={20} className="text-gray-400" />
            </div>
            
            <div className="flex items-center justify-between p-4 hover:bg-gray-50 rounded-xl cursor-pointer transition-colors group">
              <div className="flex items-start gap-4">
                <div className="mt-0.5 text-gray-500"><Lock size={20} /></div>
                <div>
                  <p className="text-sm font-semibold text-gray-900 group-hover:text-[#512BD4] transition-colors">Multi-Factor Authentication</p>
                  <p className="text-xs text-gray-500">Kelola multi-factor authentication Anda</p>
                </div>
              </div>
              <ChevronRight size={20} className="text-gray-400" />
            </div>
          </div>
        </div>

      </div>
    </div>
  );
}
