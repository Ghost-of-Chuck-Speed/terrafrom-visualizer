#!/bin/bash

# Terraform State Viewer Generator - Multi-File Support
# Usage: 
#   ./tfstate-viewer.sh file1.tfstate file2.tfstate file3.tfstate
#   ./tfstate-viewer.sh /path/to/states/*.tfstate
#   ./tfstate-viewer.sh /path/to/states/

if [ $# -eq 0 ]; then
    echo "Usage: $0 <tfstate-files-or-directory>"
    echo ""
    echo "Examples:"
    echo "  $0 terraform.tfstate"
    echo "  $0 file1.tfstate file2.tfstate file3.tfstate"
    echo "  $0 /path/to/states/*.tfstate"
    echo "  $0 /path/to/states/"
    exit 1
fi

OUTPUT_FILE="terraform-viewer.html"

# Check if jq is installed
if ! command -v jq &> /dev/null; then
    echo "Error: jq is required but not installed. Install with: apt-get install jq / brew install jq"
    exit 1
fi

# Collect all tfstate files
TFSTATE_FILES=()

for arg in "$@"; do
    if [ -d "$arg" ]; then
        # If it's a directory, find all .tfstate files in it
        while IFS= read -r file; do
            TFSTATE_FILES+=("$file")
        done < <(find "$arg" -maxdepth 1 -name "*.tfstate" -type f)
    elif [ -f "$arg" ]; then
        # If it's a file, add it
        TFSTATE_FILES+=("$arg")
    else
        echo "Warning: '$arg' not found, skipping..."
    fi
done

if [ ${#TFSTATE_FILES[@]} -eq 0 ]; then
    echo "Error: No tfstate files found"
    exit 1
fi

echo "Found ${#TFSTATE_FILES[@]} tfstate file(s)"

# Build JSON object with all state files
ALL_STATES_JSON="{"

for i in "${!TFSTATE_FILES[@]}"; do
    TFSTATE_FILE="${TFSTATE_FILES[$i]}"
    FILENAME=$(basename "$TFSTATE_FILE")
    
    echo "Parsing $FILENAME..."
    
    # Extract resources - capture parent type properly
    RESOURCES_JSON=$(jq -c '[.resources[] | select(.mode == "managed") | . as $parent | .instances[] | {
      type: $parent.type,
      name: (.attributes.name // .attributes.id // .attributes.cluster_identifier // .attributes.function_name // .attributes.bucket // "N/A"),
      id: (.attributes.id // "N/A"),
      arn: (.attributes.arn // "N/A"),
      tags: (.attributes.tags // {}),
      name_tag: (.attributes.tags.Name // "N/A"),
      port: (.attributes.port // .attributes.listener_port // "N/A"),
      protocol: (.attributes.protocol // "N/A"),
      target_type: (.attributes.target_type // "N/A"),
      min_size: (.attributes.min_size // "N/A"),
      max_size: (.attributes.max_size // "N/A"),
      desired_capacity: (.attributes.desired_capacity // "N/A"),
      instance_type: (.attributes.instance_type // .attributes.node_type // "N/A"),
      ami: (.attributes.ami // .attributes.image_id // "N/A"),
      vpc_id: (.attributes.vpc_id // "N/A"),
      engine: (.attributes.engine // "N/A"),
      engine_version: (.attributes.engine_version // "N/A")
    }]' "$TFSTATE_FILE" 2>/dev/null)
    
    if [ -z "$RESOURCES_JSON" ] || [ "$RESOURCES_JSON" == "null" ] || [ "$RESOURCES_JSON" == "[]" ]; then
        echo "Warning: Failed to extract resources from $FILENAME, skipping..."
        continue
    fi
    
    NUM_RESOURCES=$(echo "$RESOURCES_JSON" | jq 'length')
    echo "  → Found $NUM_RESOURCES resources"
    
    # Add to the all states object
    if [ $i -gt 0 ]; then
        ALL_STATES_JSON="$ALL_STATES_JSON,"
    fi
    
    # Escape the filename for JSON key (replace special chars)
    SAFE_FILENAME=$(echo "$FILENAME" | sed 's/[^a-zA-Z0-9._-]/_/g')
    
    ALL_STATES_JSON="$ALL_STATES_JSON\"$SAFE_FILENAME\":$RESOURCES_JSON"
done

ALL_STATES_JSON="$ALL_STATES_JSON}"

# Build the complete HTML in one go
cat > "$OUTPUT_FILE" <<EOF
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Terraform State Viewer</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: #0f172a;
            color: #e2e8f0;
            padding: 20px;
        }
        .container {
            max-width: 1400px;
            margin: 0 auto;
        }
        header {
            margin-bottom: 30px;
        }
        h1 {
            font-size: 2em;
            margin-bottom: 10px;
            color: #60a5fa;
        }
        .stats {
            color: #94a3b8;
            font-size: 0.9em;
        }
        .file-selector {
            background: #1e293b;
            padding: 15px 20px;
            border-radius: 8px;
            margin-bottom: 15px;
            display: flex;
            align-items: center;
            gap: 15px;
        }
        .file-selector label {
            font-weight: 600;
            color: #cbd5e1;
        }
        .file-select {
            flex: 1;
            padding: 10px;
            font-size: 1em;
            border: 2px solid #334155;
            border-radius: 6px;
            background: #0f172a;
            color: #e2e8f0;
            cursor: pointer;
        }
        .file-select:focus {
            outline: none;
            border-color: #60a5fa;
        }
        .controls {
            background: #1e293b;
            padding: 20px;
            border-radius: 8px;
            margin-bottom: 20px;
        }
        .search-box {
            width: 100%;
            padding: 12px;
            font-size: 1em;
            border: 2px solid #334155;
            border-radius: 6px;
            background: #0f172a;
            color: #e2e8f0;
            margin-bottom: 15px;
        }
        .search-box:focus {
            outline: none;
            border-color: #60a5fa;
        }
        .filters {
            display: flex;
            flex-wrap: wrap;
            gap: 10px;
        }
        .filter-btn {
            padding: 8px 16px;
            border: 2px solid #334155;
            border-radius: 6px;
            background: #1e293b;
            color: #e2e8f0;
            cursor: pointer;
            font-size: 0.9em;
            transition: all 0.2s;
        }
        .filter-btn:hover {
            border-color: #60a5fa;
            background: #334155;
        }
        .filter-btn.active {
            background: #60a5fa;
            border-color: #60a5fa;
            color: #0f172a;
            font-weight: 600;
        }
        .resources {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
            gap: 15px;
        }
        .resource-card {
            background: #1e293b;
            border: 1px solid #334155;
            border-radius: 8px;
            padding: 15px;
            transition: all 0.2s;
        }
        .resource-card:hover {
            border-color: #60a5fa;
            transform: translateY(-2px);
        }
        .resource-type {
            display: inline-block;
            background: #334155;
            color: #60a5fa;
            padding: 4px 10px;
            border-radius: 4px;
            font-size: 0.8em;
            font-weight: 600;
            margin-bottom: 10px;
        }
        .resource-name {
            font-size: 1.1em;
            font-weight: 600;
            margin-bottom: 10px;
            color: #f1f5f9;
            word-break: break-word;
        }
        .resource-details {
            font-size: 0.85em;
            color: #94a3b8;
        }
        .resource-details div {
            margin: 5px 0;
            word-break: break-all;
        }
        .resource-details strong {
            color: #cbd5e1;
        }
        .no-results {
            text-align: center;
            padding: 40px;
            color: #64748b;
            font-size: 1.1em;
        }
        .tag {
            display: inline-block;
            background: #0f172a;
            padding: 2px 8px;
            border-radius: 3px;
            margin: 2px;
            font-size: 0.8em;
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>🏗️ Terraform State Viewer</h1>
            <div class="stats">
                <span id="resource-count">Loading resources...</span>
            </div>
        </header>

        <div class="file-selector">
            <label for="file-select">📁 State File:</label>
            <select id="file-select" class="file-select">
                <!-- Options will be populated by JavaScript -->
            </select>
        </div>

        <div class="controls">
            <input type="text" id="search" class="search-box" placeholder="Search resources by name, type, ARN, or tags...">
            <div class="filters">
                <button class="filter-btn active" data-filter="all">All Resources</button>
                <button class="filter-btn" data-filter="compute">Compute (EC2/EKS/Lambda)</button>
                <button class="filter-btn" data-filter="autoscaling">Auto Scaling</button>
                <button class="filter-btn" data-filter="launch">Launch Templates</button>
                <button class="filter-btn" data-filter="load_balancer">Load Balancers</button>
                <button class="filter-btn" data-filter="target">Target Groups</button>
                <button class="filter-btn" data-filter="database">Databases (RDS/DynamoDB)</button>
                <button class="filter-btn" data-filter="storage">Storage (S3/EBS/EFS)</button>
                <button class="filter-btn" data-filter="vpc">Networking (VPC/Route53)</button>
                <button class="filter-btn" data-filter="security_group">Security Groups</button>
                <button class="filter-btn" data-filter="iam">IAM</button>
                <button class="filter-btn" data-filter="monitoring">Monitoring (CloudWatch/SNS)</button>
            </div>
        </div>

        <div id="resources" class="resources"></div>
    </div>

    <script>
        const allStatesData = $ALL_STATES_JSON;
        
        let currentStateFile = '';
        let currentFilter = 'all';
        let searchTerm = '';

        // Populate file selector dropdown
        const fileSelect = document.getElementById('file-select');
        const stateFiles = Object.keys(allStatesData);
        
        stateFiles.forEach((filename, index) => {
            const option = document.createElement('option');
            option.value = filename;
            option.textContent = filename + ' (' + allStatesData[filename].length + ' resources)';
            fileSelect.appendChild(option);
            
            // Set first file as current
            if (index === 0) {
                currentStateFile = filename;
            }
        });

        console.log('Loaded', stateFiles.length, 'state file(s)');

        function getResourceCategory(type) {
            const t = type.toLowerCase();
            
            // Compute - EKS, ECS, EC2, Lambda
            if (t.includes('eks') || t.includes('ecs') || t.includes('fargate')) return 'compute';
            if (t.includes('lambda')) return 'compute';
            if (t.includes('instance') && !t.includes('db')) return 'compute';
            
            // Auto Scaling
            if (t.includes('autoscaling') || t.includes('asg') || t.includes('_asg')) return 'autoscaling';
            
            // Launch Templates
            if (t.includes('launch_template') || t.includes('launch_configuration')) return 'launch';
            
            // Load Balancing
            if (t.includes('target_group') || t.includes('lb_target_group')) return 'target';
            if (t.includes('lb_listener') || t.includes('listener')) return 'load_balancer';
            if (t.includes('load_balancer') || t.includes('_lb') || t.includes('alb') || t.includes('elb') || t.includes('nlb')) return 'load_balancer';
            
            // Security
            if (t.includes('security_group') || t.includes('_sg')) return 'security_group';
            if (t.includes('iam') || t.includes('role') || t.includes('policy') || t.includes('_user')) return 'iam';
            
            // Networking
            if (t.includes('vpc') || t.includes('subnet') || t.includes('route_table') || t.includes('route53') || t.includes('dns')) return 'vpc';
            if (t.includes('gateway') || t.includes('nat') || t.includes('internet_gateway') || t.includes('vpn')) return 'vpc';
            if (t.includes('network_interface') || t.includes('eni') || t.includes('eip') || t.includes('elastic_ip')) return 'vpc';
            
            // Databases
            if (t.includes('rds') || t.includes('db_instance') || t.includes('db_cluster')) return 'database';
            if (t.includes('dynamodb') || t.includes('documentdb') || t.includes('docdb')) return 'database';
            if (t.includes('elasticache') || t.includes('memorydb')) return 'database';
            if (t.includes('neptune') || t.includes('redshift')) return 'database';
            
            // Storage
            if (t.includes('s3') || t.includes('bucket')) return 'storage';
            if (t.includes('efs') || t.includes('fsx')) return 'storage';
            if (t.includes('ebs') || t.includes('volume') || t.includes('snapshot')) return 'storage';
            
            // Monitoring & Logs
            if (t.includes('cloudwatch') || t.includes('log_group') || t.includes('metric') || t.includes('alarm')) return 'monitoring';
            if (t.includes('sns') || t.includes('sqs') || t.includes('eventbridge')) return 'monitoring';
            
            return 'other';
        }

        function matchesSearch(resource, term) {
            if (!term) return true;
            const searchable = [
                resource.name,
                resource.type,
                resource.id,
                resource.arn,
                resource.name_tag,
                JSON.stringify(resource.tags)
            ].join(' ').toLowerCase();
            return searchable.includes(term.toLowerCase());
        }

        function renderResources() {
            const container = document.getElementById('resources');
            const resourcesData = allStatesData[currentStateFile] || [];
            
            const filtered = resourcesData.filter(r => {
                const category = getResourceCategory(r.type);
                const matchesFilter = currentFilter === 'all' || category === currentFilter;
                const matchesSearchTerm = matchesSearch(r, searchTerm);
                return matchesFilter && matchesSearchTerm;
            });

            document.getElementById('resource-count').textContent = 
                \`\${filtered.length} of \${resourcesData.length} resources in \${currentStateFile}\`;

            if (filtered.length === 0) {
                container.innerHTML = '<div class="no-results">No resources found</div>';
                return;
            }

            container.innerHTML = filtered.map(r => {
                const details = [];
                
                if (r.arn && r.arn !== 'N/A') details.push(\`<strong>ARN:</strong> \${r.arn}\`);
                if (r.id && r.id !== 'N/A' && r.id !== r.arn) details.push(\`<strong>ID:</strong> \${r.id}\`);
                if (r.name_tag && r.name_tag !== 'N/A') details.push(\`<strong>Name Tag:</strong> \${r.name_tag}\`);
                
                // Type-specific details
                if (r.port && r.port !== 'N/A') details.push(\`<strong>Port:</strong> \${r.port}\`);
                if (r.protocol && r.protocol !== 'N/A') details.push(\`<strong>Protocol:</strong> \${r.protocol}\`);
                if (r.target_type && r.target_type !== 'N/A') details.push(\`<strong>Target Type:</strong> \${r.target_type}\`);
                if (r.min_size && r.min_size !== 'N/A') details.push(\`<strong>ASG Capacity:</strong> \${r.min_size} / \${r.desired_capacity} / \${r.max_size} (min/desired/max)\`);
                if (r.instance_type && r.instance_type !== 'N/A') details.push(\`<strong>Instance Type:</strong> \${r.instance_type}\`);
                if (r.ami && r.ami !== 'N/A') details.push(\`<strong>AMI:</strong> \${r.ami}\`);
                if (r.vpc_id && r.vpc_id !== 'N/A') details.push(\`<strong>VPC:</strong> \${r.vpc_id}\`);
                if (r.engine && r.engine !== 'N/A') details.push(\`<strong>Engine:</strong> \${r.engine} \${r.engine_version !== 'N/A' ? r.engine_version : ''}\`);

                // Tags
                if (r.tags && Object.keys(r.tags).length > 0) {
                    const tagHtml = Object.entries(r.tags)
                        .map(([k, v]) => \`<span class="tag">\${k}: \${v}</span>\`)
                        .join('');
                    details.push(\`<strong>Tags:</strong><br>\${tagHtml}\`);
                }

                return \`
                    <div class="resource-card">
                        <div class="resource-type">\${r.type}</div>
                        <div class="resource-name">\${r.name}</div>
                        <div class="resource-details">
                            \${details.map(d => \`<div>\${d}</div>\`).join('')}
                        </div>
                    </div>
                \`;
            }).join('');
        }

        // File selector change event
        fileSelect.addEventListener('change', (e) => {
            currentStateFile = e.target.value;
            // Reset filter and search when switching files
            currentFilter = 'all';
            searchTerm = '';
            document.getElementById('search').value = '';
            document.querySelectorAll('.filter-btn').forEach(b => b.classList.remove('active'));
            document.querySelector('[data-filter="all"]').classList.add('active');
            renderResources();
        });

        // Filter button event listeners
        document.querySelectorAll('.filter-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                document.querySelectorAll('.filter-btn').forEach(b => b.classList.remove('active'));
                e.target.classList.add('active');
                currentFilter = e.target.dataset.filter;
                renderResources();
            });
        });

        // Search event listener
        document.getElementById('search').addEventListener('input', (e) => {
            searchTerm = e.target.value;
            renderResources();
        });

        // Initial render
        if (stateFiles.length > 0) {
            renderResources();
        } else {
            document.getElementById('resources').innerHTML = '<div class="no-results">No state files loaded</div>';
            document.getElementById('resource-count').textContent = 'No state files';
        }
    </script>
</body>
</html>
EOF

echo "✅ Generated $OUTPUT_FILE with ${#TFSTATE_FILES[@]} state file(s)"
echo "📂 Open it in your browser: file://$(pwd)/$OUTPUT_FILE"

# Try to open in browser
if command -v open &> /dev/null; then
    open "$OUTPUT_FILE"
elif command -v xdg-open &> /dev/null; then
    xdg-open "$OUTPUT_FILE"
fi