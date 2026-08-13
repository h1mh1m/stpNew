import type { Metadata } from "next";
import "./globals.css";

import Sidebar from "./components/Sidebar";
import Header from "./components/Header";
import Footer from "./components/Footer";

export const metadata: Metadata = {
  title: "TrackIQ",
  description: "TrackIQ Dashboard",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>
        <Sidebar />

        <Header />

        <main className="ml-56 pt-16 min-h-screen bg-[#eef3ff]">
          {children}
        </main>

        <Footer />
      </body>
    </html>
  );
}
