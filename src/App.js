// src/App.js
import React, { useState } from 'react';
import UploadForm from './UploadForm';
import StateVisualizer from './StateVisualizer';

function App() {
  const [stateData, setStateData] = useState(null); // Default to null instead of undefined

  const handleUploadSuccess = (data) => {
    console.log('Upload success:', data);
    setStateData(data); // Set the state with the uploaded file data
  };

  return (
    <div>
      <h1>Terraform State Parser and Visualizer</h1>
      <UploadForm onUploadSuccess={handleUploadSuccess} />
      {stateData ? (
        <StateVisualizer stateData={stateData} /> // Only render this if stateData is not null
      ) : (
        <p>Please upload a Terraform state file to get started.</p>
      )}
    </div>
  );
}

export default App;

