export default function Home() {
  return (
    <div className="space-y-8">

      {/* Hero */}
      <section className="bg-white rounded-xl overflow-hidden shadow">
        <div className="h-80 bg-gray-300 flex items-center justify-center">
          <h1 className="text-3xl font-bold text-white">
            Welcome To TrackIQ
          </h1>
        </div>

        <div className="p-6 text-center">
          <p className="text-gray-600">
            Get your adrenaline pumping on the best track!
            The ITS Circuit is ready to witness your speed today.
          </p>
        </div>
      </section>


      {/* Directors */}
      <section className="bg-white rounded-xl p-8 shadow">
        <h2 className="text-2xl font-bold text-center mb-6">
          Our Directors
        </h2>

        <div className="grid md:grid-cols-2 gap-8">
          <div className="h-56 bg-purple-600 rounded-xl flex items-end p-5">
            <p className="font-bold">
              Pak Julendra <br />
              Director 1
            </p>
          </div>

          <div className="h-56 bg-purple-600 rounded-xl flex items-end p-5">
            <p className="font-bold">
              Pak Julendra <br />
              Director 1
            </p>
          </div>
        </div>
      </section>


      {/* Vision Mission */}
      <section className="bg-white rounded-xl p-8 shadow space-y-8">

        <div className="grid md:grid-cols-2 gap-8 items-center">
          <div className="h-48 bg-gray-300 rounded-xl"></div>

          <div>
            <h2 className="text-2xl font-bold">
              Our Vision
            </h2>
            <p className="text-gray-600 mt-3">
              Creating the best racing experience and
              innovative track management system.
            </p>
          </div>
        </div>


        <div className="grid md:grid-cols-2 gap-8 items-center">

          <div>
            <h2 className="text-2xl font-bold">
              Our Mission
            </h2>
            <p className="text-gray-600 mt-3">
              Providing safe, modern, and professional
              circuit facilities.
            </p>
          </div>

          <div className="h-48 bg-gray-300 rounded-xl"></div>

        </div>

      </section>

    </div>
  );
}
