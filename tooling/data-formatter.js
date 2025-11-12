const fs = require('fs');
const path = require('path');

// Hardcoded file paths
const INPUT_FILE = path.join(__dirname, 'data','output.json');
const OUTPUT_FILE = path.join(__dirname, 'data','formatted-data.json');

/**
 * Converts a money string (e.g., "$1,020.00" or "$9.50") to a numeric value.
 * @param {string} moneyString - The money string to convert
 * @returns {number} The numeric value with 2 decimal places
 */
function parseMoneyString(moneyString) {
  if (typeof moneyString !== 'string') {
    return moneyString;
  }
  // Remove dollar sign and commas, then parse as float
  const numericValue = parseFloat(moneyString.replace(/[$,]/g, ''));
  // Round to 2 decimal places
  return numericValue;
}

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

    // Convert target_from and target_to to numeric values
    const processedItems = allItems.map(item => {
      const processedItem = { ...item };
      if (processedItem.target_from) {
        processedItem.target_from = parseMoneyString(processedItem.target_from);
      }
      if (processedItem.target_to) {
        processedItem.target_to = parseMoneyString(processedItem.target_to);
      }
      return processedItem;
    });

    // Write to output file
    fs.writeFileSync(OUTPUT_FILE, JSON.stringify(processedItems, null, 2), 'utf8');

    console.log(`✓ Successfully formatted data`);
    console.log(`  - Input: ${INPUT_FILE}`);
    console.log(`  - Output: ${OUTPUT_FILE}`);
    console.log(`  - Total items extracted: ${processedItems.length}`);

    return processedItems;
  } catch (error) {
    console.error('Error formatting data:', error.message);
    throw error;
  }
}

// Run the function if this script is executed directly
if (require.main === module) {
  formatData();
}

