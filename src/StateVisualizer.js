// src/StateVisualizer.js
import React from 'react';
import { Tooltip } from 'react-tooltip'; // Correct import for react-tooltip v4+

const StateVisualizer = ({ stateData }) => {
  // Guard against undefined or null stateData
  if (!stateData) {
    return <p>No state data available. Please upload a Terraform state file.</p>;
  }

  const { terraform_version, version, resources, modules, outputs, variables, dataSources } = stateData;

  return (
    <div>
      <h2>Terraform State Visualization</h2>
      <p><strong>Terraform Version:</strong> {terraform_version}</p>
      <p><strong>State Version:</strong> {version}</p>

      <div>
        <h3>Resources</h3>
        {resources && resources.length === 0 ? (
          <p>No resources found in the state file.</p>
        ) : (
          <ul>
            {resources.map((resource, index) => (
              <li key={index} data-tip={`Type: ${resource.type}`}>
                <h4>{resource.name} (Type: {resource.type})</h4>
                <ul>
                  {resource.instances.map((instance, i) => (
                    <li key={i}>
                      <strong>Attributes:</strong>
                      <pre>{JSON.stringify(instance.attributes, null, 2)}</pre>
                    </li>
                  ))}
                </ul>
                <Tooltip />
              </li>
            ))}
          </ul>
        )}
      </div>

      <div>
        <h3>Modules</h3>
        {modules && modules.length === 0 ? (
          <p>No modules found in the state file.</p>
        ) : (
          <ul>
            {modules.map((module, index) => (
              <li key={index}>
                <h4>{module.path.join(' > ')}</h4>
                <ul>
                  {module.resources.map((resource, idx) => (
                    <li key={idx}>
                      {resource.name} (Type: {resource.type})
                    </li>
                  ))}
                </ul>
              </li>
            ))}
          </ul>
        )}
      </div>

      <div>
        <h3>Outputs</h3>
        {outputs ? (
          <pre>{JSON.stringify(outputs, null, 2)}</pre>
        ) : (
          <p>No outputs found in the state file.</p>
        )}
      </div>

      <div>
        <h3>Variables</h3>
        {variables ? (
          <pre>{JSON.stringify(variables, null, 2)}</pre>
        ) : (
          <p>No variables found in the state file.</p>
        )}
      </div>

      <div>
        <h3>Data Sources</h3>
        {dataSources ? (
          <pre>{JSON.stringify(dataSources, null, 2)}</pre>
        ) : (
          <p>No data sources found in the state file.</p>
        )}
      </div>
    </div>
  );
};

export default StateVisualizer;

