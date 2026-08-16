import React from 'react';

const services = [
  { id: 1, title: 'PENGUJIAN KAYU', price: 'Rp 300.000', image: 'https://images.unsplash.com/photo-1589939705384-5185137a7f0f?q=80&w=500&auto=format&fit=crop' },
  { id: 2, title: 'Pengukuran persentase kadar air dan berat jenis', price: 'Rp 350.000 - Rp 450.000', image: 'https://images.unsplash.com/photo-1581092580497-e0d23cbdf1dc?q=80&w=500&auto=format&fit=crop' },
  { id: 3, title: 'PENGUJIAN BETON', price: 'Rp 0 - Rp 900.000', image: 'https://images.unsplash.com/photo-1504328345606-18bbc8c9d7d1?q=80&w=500&auto=format&fit=crop' },
  { id: 4, title: 'PENGUJIAN BAJA', price: 'Rp 0', image: 'https://images.unsplash.com/photo-1504328345606-18bbc8c9d7d1?q=80&w=500&auto=format&fit=crop' },
  { id: 5, title: 'Pengujian Lanjutan', price: 'Rp 600.000 - Rp 1.500.000', image: 'https://images.unsplash.com/photo-1581092580497-e0d23cbdf1dc?q=80&w=500&auto=format&fit=crop' },
];

export default function ServiceGrid() {
  return (
    <div className="py-12 border-t border-gray-200">
      <h2 className="text-2xl font-bold text-gray-900 mb-8 px-2 border-l-4 border-[#512BD4]">Layanan Terbaru</h2>
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6 px-2">
        {services.map((service) => (
          <div key={service.id} className="bg-white rounded-xl border border-gray-200 overflow-hidden hover:shadow-xl transition-shadow cursor-pointer flex flex-col h-full transform hover:-translate-y-1 transition-all duration-300">
            <div className="h-48 w-full bg-gray-100 relative">
              <img src={service.image} alt={service.title} className="w-full h-full object-cover" />
            </div>
            <div className="p-4 flex flex-col flex-grow">
              <h3 className="font-semibold text-gray-900 mb-3 line-clamp-2 leading-tight">{service.title}</h3>
              <p className="text-[#512BD4] font-bold mt-auto text-lg">{service.price}</p>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
