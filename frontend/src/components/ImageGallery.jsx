import { deleteImage } from "@/pages/api/api";
import { useEffect, useState } from "react";

export default function ImageGallery({ filteredPhotos }) {
  const [images, setImages] = useState(filteredPhotos);
  const [loading, setLoading] = useState(false);

  // For full-image modal
  const [selectedImage, setSelectedImage] = useState(null);

  useEffect(() => {
    setImages(filteredPhotos);
  }, [filteredPhotos]);

  // Load more images
  const loadMore = () => {
    setLoading(true);
    setTimeout(() => {
      setImages((prev) => [...filteredPhotos]);
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

  // delete image
  const handleDelete = async (id) => {
    // here ids add into array

    try {
      const res = await deleteImage([id]);
      console.log(res);
      loadMore()
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <div>
      {/* Image Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
        {images?.length > 0 ? (
          images.map((img) => (
            <div
              key={img.id}
              className="relative bg-white rounded-2xl shadow-md overflow-hidden cursor-pointer hover:scale-105 transition duration-300 ease-in-out"
            >
              {/* Delete Button */}
              <button
                onClick={(e) => {
                  e.stopPropagation(); // prevent opening preview when deleting
                  handleDelete(img.id); // your delete function
                }}
                className="absolute top-2 right-2 z-10 bg-red-500 text-white p-2 rounded-full shadow hover:bg-red-600 transition"
              >
                {/* Trash Icon (SVG) */}
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  className="w-4 h-4"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  strokeWidth={2}
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                  />
                </svg>
              </button>

              {/* Image / Video */}
              {img.filename.includes(".mp4") ? (
                <video
                  src={img.url}
                  className="w-full h-40 object-cover"
                  controls
                  onClick={() => setSelectedImage(img)}
                />
              ) : (
                <img
                  src={img.url}
                  alt={img.filename}
                  className="w-full h-40 object-cover"
                  onClick={() => setSelectedImage(img)}
                />
              )}
            </div>
          ))
        ) : (
          <h1 className="text-2xl font-bold text-center mt-6">
            No images found
          </h1>
        )}
      </div>

      {images?.length === 0 && (
        <h1 className="text-2xl font-bold text-center mt-6">No images found</h1>
      )}

      {/* Loader */}
      {loading && (
        <div className="flex justify-center items-center mt-6">
          <div className="w-8 h-8 border-4 border-blue-600 border-t-transparent rounded-full animate-spin"></div>
        </div>
      )}

      {/* Full Image Modal */}
      {selectedImage && (
        <div
          className="fixed inset-0 bg-black bg-opacity-70 flex items-center justify-center z-50"
          onClick={() => setSelectedImage(null)}
        >
          {selectedImage.filename.includes(".mp4") ? (
            <video
              src={selectedImage.url}
              autoPlay
              className="max-h-[100%] max-w-[100%] rounded-lg shadow-lg"
            />
          ) : (
            <img
              src={selectedImage.url}
              alt="Full"
              className="max-h-[100%] max-w-[100%] rounded-lg shadow-lg"
              onClick={(e) => e.stopPropagation()}
            />
          )}
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
