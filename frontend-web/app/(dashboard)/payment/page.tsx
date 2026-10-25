'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { ArrowLeft, ChevronDown, ChevronUp, X, CheckCircle2 } from 'lucide-react';
import { useRouter } from 'next/navigation';

export default function PaymentPage() {
  const router = useRouter();
  const [isAccordionOpen, setIsAccordionOpen] = useState(true);
  const [isSummaryModalOpen, setIsSummaryModalOpen] = useState(false);

  const handleNext = () => {
    // Default behavior for UI testing
    setIsSummaryModalOpen(true);
  };

  const handleCompleteOrder = () => {
    setIsSummaryModalOpen(false);
    router.push('/');
  };

  return (
    <div className="max-w-7xl mx-auto pb-12">
      <div className="bg-white p-8 md:p-12 rounded-3xl shadow-sm border border-gray-100 min-h-[800px] relative">
        
        {/* Header */}
        <div className="flex items-center gap-4 mb-10">
          <Link href="/booking" className="p-2 hover:bg-gray-100 rounded-full transition-colors">
            <ArrowLeft className="w-6 h-6 text-gray-900" />
          </Link>
          <h1 className="text-3xl font-bold text-gray-900">Payment</h1>
        </div>

        <div className="max-w-3xl space-y-10">
          
          {/* BNI Logo Placeholder */}
          <div className="flex items-center text-4xl font-black text-[#F15A23] tracking-tighter">
            <span className="text-[#005E6A] mr-1">BNI</span>
          </div>

          {/* Informasi penting */}
          <div>
            <h3 className="text-xl font-bold text-gray-900 mb-4">Informasi penting</h3>
            <ul className="space-y-3 text-gray-600 list-disc pl-5">
              <li>Pembayaran ini hanya berlaku dengan menggunakan akun BNI</li>
              <li>Pastikan nomor Virtual Account Anda sesuai dengan cara pembayaran</li>
              <li>Pastikan bahwa nominal pembayaran sudah sesuai dengan pemesanan yang Anda lakukan</li>
              <li>Untuk pembayaran dengan menggunakan bank lain, mohon menggunakan metode transfer bank online.</li>
              <li>Untuk transfer antar bank, mohon tidak menggunakan metode pembayaran BI Fast</li>
            </ul>
          </div>

          {/* Silahkan Transfer Ke */}
          <div>
            <h3 className="text-xl font-bold text-gray-900 mb-4">Silahkan Transfer Ke</h3>
            <div className="border border-gray-200 rounded-2xl p-6 bg-white">
              <div className="flex justify-between items-center mb-6 pb-6 border-b border-gray-100">
                <div>
                  <p className="text-sm text-gray-500 mb-1">Virtual Account</p>
                  <p className="text-lg font-bold text-gray-900">BNI</p>
                </div>
                <div className="text-3xl font-black text-[#F15A23] tracking-tighter">
                  <span className="text-[#005E6A] mr-1">BNI</span>
                </div>
              </div>

              <div className="flex justify-between items-end mb-6 pb-6 border-b border-gray-100">
                <div>
                  <p className="text-sm text-gray-500 mb-1">Nomor Virtual Account</p>
                  <p className="text-2xl font-bold text-gray-900 tracking-wider">XXX XXX XXX XXX</p>
                </div>
                <button className="text-[#0055b8] text-sm font-bold tracking-wide hover:underline uppercase">
                  SALIN
                </button>
              </div>

              <div>
                <p className="text-sm text-gray-500 mb-1">Total Pembayaran</p>
                <p className="text-3xl font-bold text-[#007bff]">Rp XX.XXX.XXX</p>
              </div>
            </div>
          </div>

          {/* Tata Cara Pembayaran */}
          <div>
            <h3 className="text-xl font-bold text-gray-900 mb-4">Tata Cara Pembayaran Melalui BNI</h3>
            <div className="border border-gray-200 rounded-2xl overflow-hidden">
              <button 
                onClick={() => setIsAccordionOpen(!isAccordionOpen)}
                className="w-full flex justify-between items-center p-6 bg-white hover:bg-gray-50 transition-colors"
              >
                <span className="font-medium text-gray-900">Transfer Melalui ATM</span>
                {isAccordionOpen ? <ChevronUp className="w-5 h-5 text-gray-500" /> : <ChevronDown className="w-5 h-5 text-gray-500" />}
              </button>
              
              {isAccordionOpen && (
                <div className="px-6 pb-6 pt-2 bg-white">
                  <ul className="space-y-3 text-gray-600 text-sm">
                    <li>Masukkan Kartu ATM</li>
                    <li>Pilih <span className="font-medium">Bahasa</span></li>
                    <li>Masukkan <span className="font-medium">PIN ATM</span></li>
                    <li>Pilih <span className="font-medium">Menu Lainnya</span></li>
                    <li>Pilih <span className="font-medium">Transfer</span></li>
                    <li>Pilih <span className="font-medium">Jenis rekening</span> yang akan Anda gunakan (Contoh: <span className="font-medium">Rekening Tabungan</span>)</li>
                  </ul>
                </div>
              )}
            </div>
          </div>

          <div className="pt-8 flex justify-center">
            <button 
              onClick={handleNext}
              className="bg-[#512BD4] hover:bg-[#4323af] text-white px-12 py-3 rounded-lg font-medium transition-colors w-64"
            >
              Next
            </button>
          </div>

        </div>
      </div>

      {/* Order Summary Modal */}
      {isSummaryModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
          <div className="bg-white rounded-3xl p-8 max-w-sm w-full shadow-xl relative text-center">
            <button 
              onClick={() => setIsSummaryModalOpen(false)}
              className="absolute top-6 right-6 text-gray-400 hover:text-gray-600"
            >
              <X className="w-6 h-6" />
            </button>
            <div className="flex justify-center mb-6">
              <div className="w-20 h-20 bg-green-100 rounded-full flex items-center justify-center">
                <CheckCircle2 className="w-10 h-10 text-green-500" />
              </div>
            </div>
            <h3 className="text-2xl font-bold text-gray-900 mb-2">Pemesanan Berhasil</h3>
            <p className="text-gray-500 text-sm mb-8">
              Terima kasih, pembayaran Anda telah berhasil kami terima.
            </p>
            <div className="bg-gray-50 rounded-xl p-4 text-left mb-8 border border-gray-100">
              <p className="text-sm text-gray-500 mb-1">Total yang dibayar</p>
              <p className="text-xl font-bold text-gray-900">Rp XX.XXX.XXX</p>
            </div>
            <button 
              onClick={handleCompleteOrder}
              className="bg-[#512BD4] hover:bg-[#4323af] text-white px-8 py-3 rounded-lg font-medium transition-colors w-full"
            >
              Selesai
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
