import { useState, useEffect } from "react";
import ImageGallery from "./ImageGallery";
import { useAuth } from "@/components/authProvider";
import { jwtDecode } from "jwt-decode";
import { createFolder, getFolders, getPhotos } from "@/pages/api/api";
import axios from "axios";
export default function Dashboard() {
  const token = useAuth();
  const { logout } = useAuth();

  const user = token && jwtDecode(JSON.stringify(token));
  const [folders, setFolders] = useState([]);
  const [activeMenu, setActiveMenu] = useState({
    name: null,
    id: null,
  });

  const [allPhotos, setAllPhotos] = useState([]);
  const getAllFolders = async () => {
    try {
      const res = await getFolders();
      if (!res) {
        setFolders([]);
        setActiveMenu({});
      }
      setFolders(res.data);
      setActiveMenu({ name: res.data[0]?.name, id: res.data[0]?.id });
    } catch (error) {}
  };

  const getAllPhotos = async () => {
    try {
      const res = await getPhotos();
      console.log(res.data);
      setAllPhotos(res.data);
    } catch (error) {}
  };

  const filteredPhotos =
    allPhotos && allPhotos?.filter((photo) => photo.folderId === activeMenu.id);

  useEffect(() => {
    getAllFolders();
    getAllPhotos();
  }, []);
  useEffect(() => {
    if (folders?.length > 0) {
      setActiveMenu({
        name: folders[0].name,
        id: folders[0].id,
      });
    }
  }, [folders]);
  // For upload modal
  const [showUploadModal, setShowUploadModal] = useState(false);
  const [uploadFiles, setUploadFiles] = useState([]);
  const uploadHandler = async () => {
    if (!uploadFiles.length) return alert("Please select files");

    try {
      const formData = new FormData();
      uploadFiles.forEach((file) => formData.append("files", file));
      formData.append("folderId", activeMenu.id); // folder ID

      // send to backend
      const res = await axios.post(
        process.env.NEXT_PUBLIC_BASE_URL + "/api/photos/upload",
        formData,
        {
          headers: {
            "Content-Type": "multipart/form-data",
            Authorization: `Bearer ${localStorage.getItem("token")}`,
          },
        }
      );
      if (res.status == 200) {
        setShowUploadModal(false);
        getAllPhotos();
      }
    } catch (error) {
      console.error(error);
      alert("Upload failed");
    }
  };
  // create new folder
  const [showNewFolderModal, setShowNewFolderModal] = useState(false);
  const [newFolderName, setNewFolderName] = useState("");
  const createNewFolder = async () => {
    if (!newFolderName) return alert("Please enter folder name");
    try {
      const res = await createFolder(newFolderName)
      if (res.status == 200) {
        setShowNewFolderModal(false);
        getAllFolders();
      }
    } catch (error) {
      console.error(error);
      alert("Folder creation failed");
    }
  };
  const [open, setOpen] = useState(false);
  return (
    <div className="flex h-screen bg-gray-100">
      {/* Sidebar */}
     <>
  {/* Mobile Hamburger */}
  <button
    onClick={() => setOpen(!open)}  // ✅ toggle instead of only true
    className={`md:hidden w-10 h-10 fixed top-4 ${open ? "right-4" : "left-4"} z-50 bg-gray-800 text-white p-2 rounded-lg`}
  >
    {open ? "✕" : "☰"} {/* ✅ show X when open */}
  </button>

  {/* Sidebar (desktop + mobile) */}
  <div
    className={`
      fixed top-0 left-0 h-full w-60 bg-white shadow-lg flex flex-col transform transition-transform duration-300
      ${open ? "translate-x-0 z-[999]" : "-translate-x-full"}
      md:translate-x-0 md:static md:block
    `}
  >
    {/* Profile Section */}
    <div className="flex items-center space-x-3 p-4 border-b">
      <div>
        <h3 className="text-lg font-semibold">{user?.username}</h3>
        <p className="text-sm text-gray-500">{user?.email}</p>
      </div>

      {/* Logout icon */}
      <button
        onClick={logout}
        className="ml-auto cursor-pointer w-8 h-8 flex items-center justify-center rounded-full hover:bg-gray-200 transition"
      >
        <img src="/logout.svg" alt="Logout" />
      </button>
    </div>

    {/* Menu Section */}
    <nav className="flex-1 p-4 space-y-2 overflow-y-auto">
      {folders?.length > 0 ? (
        folders.map((menu) => (
          <button
            key={menu.id}
            onClick={() => setActiveMenu(menu)}
            className={`flex w-full text-left px-1 py-2 rounded-lg transition ${
              activeMenu?.id === menu.id
                ? "bg-gray-500 text-white"
                : "text-gray-700 hover:bg-gray-200"
            }`}
          >
            <img src="/gallery.svg" className="w-6 h-6" alt="Gallery" />
            <span className="ml-2">{menu.name}</span>
          </button>
        ))
      ) : (
        <h2 className="text-gray-500 ">No Folders</h2>
      )}
    </nav>

    {/* Create Folder Button */}
    <div className="p-4 fixed bottom-0 w-full border-t">
      <button
        className="bg-blue-500 text-white px-4 py-2 w-full rounded-lg"
        onClick={() => setShowNewFolderModal(true)}
      >
        Create Folder
      </button>
    </div>
  </div>

  {/* Dark overlay (only on mobile) */}
  {open && (
    <div
      className="fixed inset-0  bg-opacity-50 z-40 md:hidden"
      onClick={() => setOpen(false)}
    />
  )}
</>


      {/* Main Content */}
      <div className="flex-1 p-6 overflow-y-auto ">
        <h2 className="text-2xl font-bold mb-4 text-center md:text-left ">
          {activeMenu && activeMenu.name}
        </h2>

        {/* Image Grid */}
        <ImageGallery filteredPhotos={filteredPhotos} />
      </div>
      {/* Upload Button */}
      <div className="fixed bottom-4 right-4">
        <button
          className="bg-blue-500 text-white px-4 py-2 rounded-lg"
          onClick={() => setShowUploadModal(true)}
        >
          Upload
        </button>
      </div>

      {/* Upload Modal */}
      {showUploadModal && (
        <div
          className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
          onClick={() => setShowUploadModal(false)}
        >
          <div
            className="bg-white rounded-lg p-6 w-100"
            onClick={(e) => e.stopPropagation()} // prevent closing modal when clicking inside
          >
            <h2 className="text-lg font-bold mb-4">Upload Image / Video</h2>
            <input
              type="file"
              multiple
              accept="image/*,video/*"
              className="mb-4 p-2 border border-gray-300 rounded-lg"
              placeholder="Upload Image / Video"
              onChange={(e) => setUploadFiles(Array.from(e.target.files))}
            />
            <div className="flex justify-end space-x-2">
              <button
                className="bg-gray-300 px-4 py-2 rounded-lg"
                onClick={() => setShowUploadModal(false)}
              >
                Cancel
              </button>
              <button
                onClick={() => uploadHandler()}
                className="bg-blue-500 text-white px-4 py-2 rounded-lg"
              >
                Upload
              </button>
            </div>
          </div>
        </div>
      )}

      {/* New Folder Modal */}
      {showNewFolderModal && (
        <div
          className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
          onClick={() => setShowNewFolderModal(false)}
        >
          <div
            className="bg-white rounded-lg p-6"
            onClick={(e) => e.stopPropagation()} // prevent closing modal when clicking inside
          >
            <h2 className="text-lg font-bold mb-4">Create New Folder</h2>
            <input
              type="text"
              className="mb-4 p-2 border border-gray-300 rounded-lg"
              placeholder="Folder Name"
              onChange={(e) => setNewFolderName(e.target.value)}
            />
            <div className="flex justify-end space-x-2">
              <button
                className="bg-gray-300 px-4 py-2 rounded-lg"
                onClick={() => setShowNewFolderModal(false)}
              >
                Cancel
              </button>
              <button
                onClick={() => createNewFolder()}
                className="bg-blue-500 text-white px-4 py-2 rounded-lg"
              >
                Create
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
