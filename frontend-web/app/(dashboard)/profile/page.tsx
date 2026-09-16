'use client';

import React from 'react';
import { User, Edit2, IdCard, Mail, Phone, ChevronRight, LayoutGrid, Lock, ShieldCheck, User2, CalendarDays } from 'lucide-react';

export default function ProfilePage() {
  return (
    <div className="max-w-7xl mx-auto pb-12">
      <div className="space-y-6">
        <h1 className="text-3xl font-bold text-gray-900 mb-8 px-2">Akun</h1>

        {/* User Card */}
        <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-8 flex flex-col md:flex-row md:items-center gap-8">
          <div className="flex items-center gap-6">
            <div className="relative">
              <div className="w-24 h-24 bg-[#512BD4] rounded-full flex items-center justify-center text-white">
                <User className="w-10 h-10" />
              </div>
              <button className="absolute bottom-0 right-0 p-1.5 bg-white border border-gray-200 rounded-full shadow-sm hover:bg-gray-50">
                <Edit2 className="w-4 h-4 text-gray-500" />
              </button>
            </div>
            <div>
              <h2 className="text-2xl font-bold text-gray-900">Nama</h2>
              <p className="text-gray-500">Status</p>
            </div>
          </div>
          
          <div className="flex-1 grid grid-cols-1 md:grid-cols-2 gap-6 pl-0 md:pl-8 md:border-l border-gray-100">
            <div className="flex items-center gap-4">
              <div className="w-10 h-10 bg-gray-100 rounded-full flex items-center justify-center shrink-0">
                <IdCard className="w-5 h-5 text-gray-600" />
              </div>
              <div>
                <p className="font-medium text-gray-900">0123456789</p>
                <p className="text-xs text-gray-500">ID</p>
              </div>
            </div>
            <div className="flex items-center gap-4">
              <div className="w-10 h-10 bg-gray-100 rounded-full flex items-center justify-center shrink-0">
                <Mail className="w-5 h-5 text-gray-600" />
              </div>
              <div>
                <p className="font-medium text-gray-900">Test@Test.com</p>
                <p className="text-xs text-gray-500 flex items-center gap-1">
                  Email Utama <span className="w-1.5 h-1.5 rounded-full bg-green-500 inline-block"></span>
                </p>
              </div>
            </div>
            <div className="flex items-center gap-4">
              <div className="w-10 h-10 bg-gray-100 rounded-full flex items-center justify-center shrink-0">
                <Phone className="w-5 h-5 text-gray-600" />
              </div>
              <div>
                <p className="font-medium text-gray-900">Belum diatur</p>
                <p className="text-xs text-gray-500 flex items-center gap-1">
                  Nomor Ponsel <span className="w-1.5 h-1.5 rounded-full bg-orange-500 inline-block"></span>
                </p>
              </div>
            </div>
          </div>
        </div>

        {/* Microsoft Section */}
        <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
          <div className="px-6 py-4 border-b border-gray-100 flex items-center gap-2">
            <LayoutGrid className="w-5 h-5 text-[#512BD4]" />
            <h3 className="font-semibold text-gray-900">Microsoft</h3>
          </div>
          <button className="w-full px-6 py-5 flex items-center justify-between hover:bg-gray-50 transition-colors text-left">
            <div className="flex items-center gap-4">
              <div className="w-10 h-10 rounded-full border border-gray-200 flex items-center justify-center text-gray-500">
                <User2 className="w-5 h-5" />
              </div>
              <div>
                <p className="font-medium text-gray-900 text-sm">Pusat Akun Microsoft</p>
                <p className="text-xs text-gray-500 mt-0.5">Ubah foto profil, multi-factor authentication, dan lainnya</p>
              </div>
            </div>
            <ChevronRight className="w-5 h-5 text-gray-400" />
          </button>
        </div>

        {/* Informasi Pribadi */}
        <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
          <div className="px-6 py-5 border-b border-gray-100">
            <h3 className="font-medium text-gray-900">Informasi pribadi</h3>
            <p className="text-sm text-gray-500 mt-1">Data yang umumnya ditampilkan di semua web myITS</p>
          </div>
          <button className="w-full px-6 py-4 flex items-center justify-between hover:bg-gray-50 transition-colors border-b border-gray-100 text-left">
            <div className="flex items-center gap-4">
              <User2 className="w-5 h-5 text-gray-500" />
              <div>
                <p className="font-medium text-gray-900 text-sm">Nama</p>
                <p className="text-xs text-gray-500">Ubah nama panggilan</p>
              </div>
            </div>
            <ChevronRight className="w-5 h-5 text-gray-400" />
          </button>
          <button className="w-full px-6 py-4 flex items-center justify-between hover:bg-gray-50 transition-colors text-left">
            <div className="flex items-center gap-4">
              <CalendarDays className="w-5 h-5 text-gray-500" />
              <div>
                <p className="font-medium text-gray-900 text-sm">Tanggal Lahir</p>
                <p className="text-xs text-gray-500">Melihat tanggal lahir</p>
              </div>
            </div>
            <ChevronRight className="w-5 h-5 text-gray-400" />
          </button>
        </div>

        {/* Kontak */}
        <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
          <div className="px-6 py-5 border-b border-gray-100">
            <h3 className="font-medium text-gray-900">Kontak</h3>
            <p className="text-sm text-gray-500 mt-1">Ubah email dan nomor ponsel</p>
          </div>
          <button className="w-full px-6 py-4 flex items-center justify-between hover:bg-gray-50 transition-colors border-b border-gray-100 text-left">
            <div className="flex items-center gap-4">
              <Mail className="w-5 h-5 text-gray-500" />
              <div>
                <p className="font-medium text-gray-900 text-sm">Email</p>
                <p className="text-xs text-gray-500">Perbarui dan verifikasi email</p>
              </div>
            </div>
            <ChevronRight className="w-5 h-5 text-gray-400" />
          </button>
          <button className="w-full px-6 py-4 flex items-center justify-between hover:bg-gray-50 transition-colors text-left">
            <div className="flex items-center gap-4">
              <Phone className="w-5 h-5 text-gray-500" />
              <div>
                <p className="font-medium text-gray-900 text-sm">Nomor Ponsel</p>
                <p className="text-xs text-gray-500">Perbarui dan verifikasi nomor telepon</p>
              </div>
            </div>
            <ChevronRight className="w-5 h-5 text-gray-400" />
          </button>
        </div>

        {/* Lainnya */}
        <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
          <div className="px-6 py-5 border-b border-gray-100">
            <h3 className="font-medium text-gray-900">Lainnya</h3>
            <p className="text-sm text-gray-500 mt-1">Kata sandi dan preferensi tampilan</p>
          </div>
          <button className="w-full px-6 py-4 flex items-center justify-between hover:bg-gray-50 transition-colors border-b border-gray-100 text-left">
            <div className="flex items-center gap-4">
              <Lock className="w-5 h-5 text-gray-500" />
              <div>
                <p className="font-medium text-gray-900 text-sm">Kata Sandi</p>
                <p className="text-xs text-gray-500">Perbarui kata sandi</p>
              </div>
            </div>
            <ChevronRight className="w-5 h-5 text-gray-400" />
          </button>
          <button className="w-full px-6 py-4 flex items-center justify-between hover:bg-gray-50 transition-colors text-left">
            <div className="flex items-center gap-4">
              <ShieldCheck className="w-5 h-5 text-gray-500" />
              <div>
                <p className="font-medium text-gray-900 text-sm">Multi-Factor Authentication</p>
                <p className="text-xs text-gray-500">Kelola multi-factor authentication Anda</p>
              </div>
            </div>
            <ChevronRight className="w-5 h-5 text-gray-400" />
          </button>
        </div>

      </div>
    </div>
  );
}
