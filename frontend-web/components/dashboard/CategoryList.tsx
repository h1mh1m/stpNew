import React from 'react';

const categories = [
  { id: 1, name: 'Rekayasa Vibrasi dan Sistem Otomotif', icon: 'https://images.unsplash.com/photo-1518770660439-4636190af475?q=80&w=200&auto=format&fit=crop' },
  { id: 2, name: 'Rekayasa Sistem dan Kontrol', icon: 'https://images.unsplash.com/photo-1581092160562-40aa08e78837?q=80&w=200&auto=format&fit=crop' },
  { id: 3, name: 'Transportasi dan Geoteknik (LTG)', icon: 'https://images.unsplash.com/photo-1503387762-592deb58ef4e?q=80&w=200&auto=format&fit=crop' },
  { id: 4, name: 'Mekanika Benda Padat', icon: 'https://images.unsplash.com/photo-1537462715879-360eeb61a0ad?q=80&w=200&auto=format&fit=crop' },
  { id: 5, name: 'Material Fungsional Maju', icon: 'https://images.unsplash.com/photo-1614935151651-0bea6508abb0?q=80&w=200&auto=format&fit=crop' },
  { id: 6, name: 'Material dan Struktur Gedung', icon: 'https://images.unsplash.com/photo-1503387762-592deb58ef4e?q=80&w=200&auto=format&fit=crop' },
];

export default function CategoryList() {
  return (
    <div className="py-12 mt-4 border-t border-gray-200">
      <h2 className="text-2xl font-bold text-gray-900 mb-8 px-2 border-l-4 border-[#512BD4]">Laboratorium</h2>
      <div className="flex overflow-x-auto space-x-6 pb-6 scrollbar-hide px-2">
        {categories.map((cat) => (
          <div key={cat.id} className="flex flex-col items-center flex-shrink-0 w-32 cursor-pointer group">
            <div className="w-24 h-24 rounded-full overflow-hidden mb-4 border-2 border-transparent group-hover:border-[#512BD4] transition-all duration-300 shadow-md group-hover:shadow-lg transform group-hover:-translate-y-1">
              <img src={cat.icon} alt={cat.name} className="w-full h-full object-cover" />
            </div>
            <p className="text-sm text-center text-gray-700 font-medium group-hover:text-[#512BD4] transition-colors line-clamp-3">
              {cat.name}
            </p>
          </div>
        ))}
      </div>
    </div>
  );
}
