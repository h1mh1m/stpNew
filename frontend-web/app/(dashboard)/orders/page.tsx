'use client';

import React from 'react';
import { FileText, MapPin, Calendar, Clock, ChevronRight } from 'lucide-react';
import Link from 'next/link';

const myOrders = [
  {
    id: 'ORD-12345',
    roomTitle: 'Lab Pemrograman 1',
    department: 'Departemen Teknik Informatika ITS',
    location: 'Zona A - Gedung Teknik Informatika Lt. 2',
    date: '15 Agustus 2026',
    time: '08:00 - 12:00',
    status: 'Menunggu Pembayaran',
    totalPrice: 'Rp 360.000',
  },
  {
    id: 'ORD-12344',
    roomTitle: 'Ruang Rapat Utama',
    department: 'Rektorat ITS',
    location: 'Zona B',
    date: '10 Agustus 2026',
    time: '13:00 - 15:00',
    status: 'Selesai',
    totalPrice: 'Rp 400.000',
  },
  {
    id: 'ORD-12340',
    roomTitle: 'Auditorium Pascasarjana',
    department: 'Sekolah Pascasarjana ITS',
    location: 'Zona B',
    date: '1 Agustus 2026',
    time: '08:00 - 17:00',
    status: 'Dibatalkan',
    totalPrice: 'Rp 15.000.000',
  }
];

export default function OrdersPage() {
  return (
    <div className="max-w-6xl mx-auto pb-12">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Pesanan Saya</h1>
        <p className="text-gray-500">Kelola dan pantau status pemesanan ruangan Anda</p>
      </div>

      <div className="space-y-6">
        {myOrders.map((order) => (
          <div key={order.id} className="bg-white rounded-xl border border-gray-200 overflow-hidden shadow-sm hover:shadow-md transition-shadow">
            <div className="p-6">
              <div className="flex flex-col md:flex-row justify-between md:items-center mb-6 pb-6 border-b border-gray-100 gap-4">
                <div className="flex items-center text-sm font-medium text-gray-500">
                  <FileText className="w-4 h-4 mr-2" />
                  Order ID: <span className="text-gray-900 ml-1">{order.id}</span>
                </div>
                <div>
                  <span className={`px-3 py-1 rounded-full text-xs font-semibold ${
                    order.status === 'Menunggu Pembayaran' ? 'bg-amber-100 text-amber-700' :
                    order.status === 'Selesai' ? 'bg-green-100 text-green-700' :
                    'bg-red-100 text-red-700'
                  }`}>
                    {order.status}
                  </span>
                </div>
              </div>

              <div className="flex flex-col md:flex-row justify-between gap-6">
                <div className="flex-1 flex gap-6">
                  {/* Thumbnail Placeholder */}
                  <div className="w-24 h-24 bg-[#512BD4] rounded-lg flex-shrink-0"></div>
                  
                  <div>
                    <h3 className="text-xl font-bold text-gray-900 mb-1">{order.roomTitle}</h3>
                    <p className="text-sm text-gray-500 mb-3">{order.department}</p>
                    <div className="flex flex-col sm:flex-row sm:items-center gap-2 sm:gap-6 text-sm text-gray-600">
                      <div className="flex items-center">
                        <MapPin className="w-4 h-4 mr-1 text-gray-400" />
                        {order.location}
                      </div>
                      <div className="flex items-center">
                        <Calendar className="w-4 h-4 mr-1 text-gray-400" />
                        {order.date}
                      </div>
                      <div className="flex items-center">
                        <Clock className="w-4 h-4 mr-1 text-gray-400" />
                        {order.time}
                      </div>
                    </div>
                  </div>
                </div>
                
                <div className="flex flex-col justify-between items-end min-w-[200px] border-t md:border-t-0 md:border-l border-gray-100 pt-4 md:pt-0 md:pl-6">
                  <div className="text-right w-full mb-4 md:mb-0">
                    <p className="text-sm text-gray-500 mb-1">Total Pembayaran</p>
                    <p className="text-xl font-bold text-[#0055b8]">{order.totalPrice}</p>
                  </div>
                  <button className="w-full md:w-auto inline-flex items-center justify-center px-4 py-2 bg-gray-50 hover:bg-gray-100 text-gray-700 border border-gray-200 rounded-lg text-sm font-semibold transition-colors">
                    Lihat Detail
                    <ChevronRight className="w-4 h-4 ml-1" />
                  </button>
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
