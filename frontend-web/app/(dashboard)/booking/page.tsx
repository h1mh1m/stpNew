'use client';

import React, { useState } from 'react';
import { ChevronLeft, ChevronRight, Calendar as CalendarIcon, X } from 'lucide-react';
import { useRouter } from 'next/navigation';

export default function BookingPage() {
  const router = useRouter();
  const [currentMonth, setCurrentMonth] = useState(new Date(2024, 7)); // August 2024 (0-indexed month)
  const [isTosModalOpen, setIsTosModalOpen] = useState(false);
  
  // Mock events
  const events = [
    { date: 30, month: 6, title: 'Drinks with Robin', color: 'bg-blue-500' },
    { date: 1, month: 7, title: 'Team updates', color: 'bg-blue-500' },
    { date: 2, month: 7, title: 'Robyn\'s birthday', color: 'bg-green-500' },
    { date: 2, month: 7, title: 'Surprise party', color: 'bg-blue-500' },
    { date: 4, month: 7, title: 'Learn to paint', color: 'bg-red-500' },
  ];

  const daysInMonth = new Date(currentMonth.getFullYear(), currentMonth.getMonth() + 1, 0).getDate();
  const firstDayOfMonth = new Date(currentMonth.getFullYear(), currentMonth.getMonth(), 1).getDay();
  
  const days = Array.from({ length: daysInMonth }, (_, i) => i + 1);
  const blanks = Array.from({ length: firstDayOfMonth }, (_, i) => i);

  const handleSubmit = () => {
    setIsTosModalOpen(true);
  };

  const handleAgreeTos = () => {
    setIsTosModalOpen(false);
    router.push('/payment');
  };

  return (
    <div className="max-w-7xl mx-auto pb-12">
      <div className="bg-white p-8 md:p-12 rounded-3xl shadow-sm border border-gray-100 relative">
        
        {/* Calendar Section */}
        <div className="mb-12">
          {/* Calendar Grid (Simplified) */}
          <div className="border border-gray-200 rounded-xl overflow-hidden mb-6 bg-white">
            <div className="grid grid-cols-7 border-b border-gray-200 bg-gray-50">
              {['Min', 'Sen', 'Sel', 'Rab', 'Kam', 'Jum', 'Sab'].map(day => (
                <div key={day} className="px-2 py-3 text-sm font-medium text-gray-500 text-center">{day}</div>
              ))}
            </div>
            <div className="grid grid-cols-7 min-h-[400px]">
              {blanks.map(blank => (
                <div key={`blank-${blank}`} className="border-b border-r border-gray-100 p-2 min-h-[100px] bg-gray-50/50"></div>
              ))}
              {days.map(day => {
                const dayEvents = events.filter(e => e.date === day && e.month === currentMonth.getMonth());
                return (
                  <div key={day} className="border-b border-r border-gray-100 p-2 min-h-[100px] relative">
                    <span className="text-sm text-gray-600 font-medium mb-1 inline-block">{day}</span>
                    <div className="space-y-1 mt-1">
                      {dayEvents.map((evt, idx) => (
                        <div key={idx} className={`${evt.color} text-white text-xs px-2 py-1 rounded truncate`}>
                          {evt.title}
                        </div>
                      ))}
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
          
          <div className="flex items-center justify-center gap-4">
            <button 
              onClick={() => setCurrentMonth(new Date(currentMonth.getFullYear(), currentMonth.getMonth() - 1, 1))}
              className="p-2 hover:bg-gray-100 rounded-full transition-colors text-gray-900"
              aria-label="Previous Month"
            >
              <ChevronLeft className="w-6 h-6" />
            </button>
            <h3 className="text-xl font-bold text-gray-900 min-w-36 text-center">
              {currentMonth.toLocaleString('id-ID', { month: 'long', year: 'numeric' })}
            </h3>
            <button 
              onClick={() => setCurrentMonth(new Date(currentMonth.getFullYear(), currentMonth.getMonth() + 1, 1))}
              className="p-2 hover:bg-gray-100 rounded-full transition-colors text-gray-900"
              aria-label="Next Month"
            >
              <ChevronRight className="w-6 h-6" />
            </button>
          </div>
        </div>

        <hr className="border-gray-200 my-8" />

        {/* Search by Date */}
        <div className="mb-10">
          <h4 className="text-sm font-bold text-gray-900 mb-4">Cari Berdasarkan Tanggal</h4>
          <div className="max-w-md">
            <label className="block text-xs text-gray-500 mb-2">Tanggal dan Waktu Mulai</label>
            <div className="relative">
              <CalendarIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
              <input 
                type="text" 
                placeholder="mm/dd/yyyy, --:-- --" 
                className="w-full pl-10 pr-4 py-2.5 border border-gray-300 rounded-lg text-sm text-gray-900 placeholder:text-gray-400 focus:outline-none focus:ring-2 focus:ring-[#512BD4]/20 focus:border-[#512BD4] bg-white"
              />
            </div>
          </div>
        </div>

        {/* Circuit Booking Request Form */}
        <div>
          <h3 className="text-2xl font-bold text-gray-900 mb-2">Formulir Permintaan Pemesanan Sirkuit</h3>
          <p className="text-gray-500 text-sm mb-8">Silakan lengkapi semua bidang di bawah ini. Tim kami akan meninjau dan mengonfirmasi permintaan pemesanan Anda.</p>
          
          <form className="space-y-6 max-w-4xl">
            <div>
              <label className="block text-sm font-medium text-gray-900 mb-2">
                Siapa Nama Anda? <span className="text-red-500">*</span>
              </label>
              <input 
                type="text" 
                placeholder="Misal, Nama Lengkap" 
                className="w-full px-4 py-2.5 border border-gray-300 rounded-lg text-sm text-gray-900 placeholder:text-gray-400 focus:outline-none focus:ring-2 focus:ring-[#512BD4]/20 focus:border-[#512BD4] bg-white"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-900 mb-2">
                Dari (Organisasi/Grup)? <span className="text-red-500">*</span>
              </label>
              <input 
                type="text" 
                placeholder="Misal, Nama Sekolah atau Klub" 
                className="w-full px-4 py-2.5 border border-gray-300 rounded-lg text-sm text-gray-900 placeholder:text-gray-400 focus:outline-none focus:ring-2 focus:ring-[#512BD4]/20 focus:border-[#512BD4] bg-white"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-900 mb-2">
                Berapa lama pemesanan? <span className="text-red-500">*</span>
              </label>
              <input 
                type="text" 
                placeholder="Misal, 2 Jam" 
                className="w-full px-4 py-2.5 border border-gray-300 rounded-lg text-sm text-gray-900 placeholder:text-gray-400 focus:outline-none focus:ring-2 focus:ring-[#512BD4]/20 focus:border-[#512BD4] bg-white"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-900 mb-2">
                Kapan ke Kapan? <span className="text-red-500">*</span>
              </label>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs text-gray-500 mb-2">Tanggal dan Waktu Mulai</label>
                  <div className="relative">
                    <CalendarIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                    <input 
                      type="text" 
                      placeholder="mm/dd/yyyy, --:-- --" 
                      className="w-full pl-10 pr-4 py-2.5 border border-gray-300 rounded-lg text-sm text-gray-900 placeholder:text-gray-400 focus:outline-none focus:ring-2 focus:ring-[#512BD4]/20 focus:border-[#512BD4] bg-white"
                    />
                  </div>
                </div>
                <div>
                  <label className="block text-xs text-gray-500 mb-2">Tanggal dan Waktu Berakhir</label>
                  <div className="relative">
                    <CalendarIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                    <input 
                      type="text" 
                      placeholder="mm/dd/yyyy, --:-- --" 
                      className="w-full pl-10 pr-4 py-2.5 border border-gray-300 rounded-lg text-sm text-gray-900 placeholder:text-gray-400 focus:outline-none focus:ring-2 focus:ring-[#512BD4]/20 focus:border-[#512BD4] bg-white"
                    />
                  </div>
                </div>
              </div>
            </div>

            <div className="pt-6 flex justify-center">
              <button 
                type="button"
                onClick={handleSubmit}
                className="bg-[#512BD4] hover:bg-[#4323af] text-white px-8 py-3 rounded-lg font-medium transition-colors"
              >
                Kirim Permintaan Pemesanan
              </button>
            </div>
          </form>
        </div>

      </div>

      {/* Terms of Service Modal */}
      {isTosModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
          <div className="bg-white rounded-3xl p-8 max-w-lg w-full shadow-xl relative">
            <button 
              onClick={() => setIsTosModalOpen(false)}
              className="absolute top-6 right-6 text-gray-400 hover:text-gray-600"
            >
              <X className="w-6 h-6" />
            </button>
            <h3 className="text-2xl font-bold text-gray-900 mb-4">Ketentuan Layanan</h3>
            <div className="text-gray-600 text-sm space-y-4 max-h-[60vh] overflow-y-auto mb-8 pr-2">
              <p>1. Pemesan wajib menjaga kebersihan dan fasilitas sirkuit selama penggunaan.</p>
              <p>2. Segala bentuk kerusakan yang diakibatkan oleh kelalaian pemesan akan menjadi tanggung jawab pemesan sepenuhnya.</p>
              <p>3. Dilarang keras membawa atau mengonsumsi minuman keras dan obat-obatan terlarang di area sirkuit.</p>
              <p>4. Pemesan harus mematuhi jadwal yang telah disepakati dan tidak melewati batas waktu pemesanan.</p>
              <p>5. Pembatalan pemesanan harus dilakukan minimal 24 jam sebelum waktu penggunaan yang dijadwalkan.</p>
              <p>6. Kami tidak bertanggung jawab atas kehilangan barang berharga milik pribadi di area sirkuit.</p>
            </div>
            <div className="flex justify-end">
              <button 
                onClick={handleAgreeTos}
                className="bg-[#512BD4] hover:bg-[#4323af] text-white px-8 py-3 rounded-lg font-medium transition-colors w-full"
              >
                Saya Setuju
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
