import Link from "next/link";

export default function Sidebar() {
  const menus = [
    { name: "Beranda", href: "/" },
    { name: "View", href: "/view" },
    { name: "Schedule Event", href: "/schedule" },
    { name: "Booking", href: "/booking" },
    { name: "Profile", href: "/profile" },
  ];

  return (
    <aside className="fixed left-0 top-0 h-screen w-56 bg-white shadow-sm px-4 py-6">
      <div className="mb-10 text-center">
        <h1 className="font-bold text-xl text-blue-900">
          TrackIQ
        </h1>
      </div>

      <nav className="flex flex-col gap-3">
        {menus.map((menu) => (
          <Link
            key={menu.name}
            href={menu.href}
            className="rounded-lg px-4 py-3 text-sm font-medium text-gray-700 hover:bg-blue-900 hover:text-white transition"
          >
            {menu.name}
          </Link>
        ))}
      </nav>
    </aside>
  );
}
