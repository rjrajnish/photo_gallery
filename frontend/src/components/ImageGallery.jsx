import { useEffect, useState } from "react";

export default function ImageGallery({ filteredPhotos }) {
  const [images, setImages] = useState(filteredPhotos);
  const [loading, setLoading] = useState(false);

  // For modal
  const [selectedImage, setSelectedImage] = useState(null);

  useEffect(() => {
    setImages(filteredPhotos);
  }, [filteredPhotos]);

  // Load more images
  const loadMore = () => {
    setLoading(true);
    setTimeout(() => {
      setImages((prev) => [...prev, ...filteredPhotos]);
      setLoading(false);
    }, 500);
  };

  // Infinite scroll
  useEffect(() => {
    const handleScroll = () => {
      if (
        window.innerHeight + document.documentElement.scrollTop + 50 >=
        document.documentElement.scrollHeight
      ) {
        if (!loading) loadMore();
      }
    };
    window.addEventListener("scroll", handleScroll);
    return () => window.removeEventListener("scroll", handleScroll);
  }, [loading]);

  return (
    <div>
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
        {images?.length > 0 &&
          images.map((img) => (
            <div
              key={img.id}
              className="bg-white rounded-2xl shadow-md overflow-hidden cursor-pointer hover:scale-105 transition duration-300 ease-in-out"
              onClick={() => setSelectedImage(img.url)}
            >
              <img
                src={img.url}
                alt={img.filename}
                className="w-full h-40 object-cover "
              />
            </div>
          ))}
      </div>
      {images.length === 0 && (
        <h1 className="text-2xl font-bold text-center mt-6">No images found</h1>
      )}

      {/* Loader */}
      {loading && (
        <div className="flex justify-center items-center mt-6">
          <div className="w-8 h-8 border-4 border-blue-600 border-t-transparent rounded-full animate-spin"></div>
        </div>
      )}

      {/* Modal */}
      {selectedImage && (
        <div
          className="fixed inset-0 bg-black bg-opacity-70 flex items-center justify-center z-50"
          onClick={() => setSelectedImage(null)}
        >
          <img
            src={selectedImage}
            alt="Full"
            className="max-h-[100%] max-w-[100%] rounded-lg shadow-lg"
            onClick={(e) => e.stopPropagation()} // prevent closing modal when clicking image
          />
          <button
            className="absolute top-4 right-4 text-white text-2xl font-bold"
            onClick={() => setSelectedImage(null)}
          >
            &times;
          </button>
        </div>
      )}
    </div>
  );
}
