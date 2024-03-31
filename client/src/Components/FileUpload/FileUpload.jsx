import React, { useEffect, useState } from 'react';
import Map from '../Map/Map';

const FileUpload = () => {
  const [geojsonData, setGeojsonData] = useState(null);
  useEffect(() => {
    console.log("Updated geojsonData:", geojsonData);
  }, [geojsonData]);

  const handleFileUpload = (event) => {
    const file = event.target.files[0];
    const reader = new FileReader();
    reader.onload = (e) => {
      const parsedData = JSON.parse(e.target.result);
      setGeojsonData(parsedData);
    };
    reader.readAsText(file);
  };
  


  return (
    <div>
      <input type="file" accept=".geojson,.kml" onChange={handleFileUpload} />
      <Map geojsonData={geojsonData} />
    </div>
  );
};

export default FileUpload;


// import React, { useState } from 'react';
// import axios from 'axios';

// const FileUpload = () => {
//   const [file, setFile] = useState(null);

//   const handleFileSelect = async (event) => {
//     const selectedFile = event.target.files[0];

//     const fileName = selectedFile.name;
//     const fileSize = (selectedFile.size / 1024).toFixed(2) + ' KB';
//     const fileType = selectedFile.type;

//     setFile({ name: fileName, size: fileSize, type: fileType });

//     // Upload GeoJSON file
//     if (fileType === 'application/vnd.geo+json') {
//       try {
//         const formData = new FormData();
//         formData.append('file', selectedFile);
//         await axios.post('http://localhost:8080/api/v1/geospatial-data', formData, {
//           headers: {
//             'Content-Type': 'multipart/form-data'
//           }
//         });
//       } catch (error) {
//         console.error('Error uploading GeoJSON file:', error);
//       }
//     }
//   };

//   return (
//     <div>
//       <h1>GeoJSON & KML File Manager</h1>
//       <input type="file" onChange={handleFileSelect} accept=".geojson,.kml" />
//       {file && (
//         <div>
//           <h2>Uploaded File:</h2>
//           <p>Name: {file.name}</p>
//           <p>Size: {file.size}</p>
//           <p>Type: {file.type}</p>
//         </div>
//       )}
//     </div>
//   );
// };

// export default FileUpload;
