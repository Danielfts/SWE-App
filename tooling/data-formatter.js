const fs = require('fs');
const path = require('path');

// Hardcoded file paths
const INPUT_FILE = path.join(__dirname, 'output.json');
const OUTPUT_FILE = path.join(__dirname, 'formatted-data.json');

/**
 * Reads a JSON file containing an array of objects with 'next_page' and 'items' keys,
 * extracts all items, flattens them into a single array, and writes to an output file.
 */
function formatData() {
  try {
    // Read the input file
    const rawData = fs.readFileSync(INPUT_FILE, 'utf8');
    const data = JSON.parse(rawData);

    // Flatten all items using flatMap
    const allItems = data.flatMap(obj => obj.items || []);

    // Write to output file
    fs.writeFileSync(OUTPUT_FILE, JSON.stringify(allItems, null, 2), 'utf8');

    console.log(`✓ Successfully formatted data`);
    console.log(`  - Input: ${INPUT_FILE}`);
    console.log(`  - Output: ${OUTPUT_FILE}`);
    console.log(`  - Total items extracted: ${allItems.length}`);

    return allItems;
  } catch (error) {
    console.error('Error formatting data:', error.message);
    throw error;
  }
}

// Run the function if this script is executed directly
if (require.main === module) {
  formatData();
}

