export default function Header() {
  return (
    <header className="fixed top-0 right-0 left-56 h-16 bg-white flex items-center justify-end px-8 shadow-sm">
      <button className="bg-blue-900 text-white px-5 py-2 rounded-full text-sm font-medium hover:bg-blue-800 transition">
        ⇥ Log out
      </button>
    </header>
  );
}
