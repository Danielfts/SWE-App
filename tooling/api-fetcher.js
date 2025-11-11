/**
 * Fetches paginated data from an API endpoint iteratively.
 * 
 * @param {string} url - The base URL to fetch from
 * @param {string} token - The authentication token (will be prefixed with "Bearer ")
 * @param {string} filePath - The file path to write results to continuously
 * @returns {Promise<Array>} An array containing all fetched response objects
 * 
 * @example
 * const results = await fetchAllPages('https://api.example.com/data', 'my-token-123', './output.json');
 * console.log(`Fetched ${results.length} pages`);
 */
async function fetchAllPages(url, token, filePath) {
  const fs = require('fs');
  const results = [];
  let nextPage = null;
  
  // Initialize the file with an empty array
  fs.writeFileSync(filePath, '[\n', 'utf8');
  
  try {
    while (true) {
      // Build the URL with next_page query parameter if available
      console.log(`Requesting Page: ${nextPage || '0'} - ${results.length}`);
      const requestUrl = new URL(url);
      if (nextPage) {
        requestUrl.searchParams.set('next_page', nextPage);
      }
      
      // Make the request with proper headers
      const response = await fetch(requestUrl.toString(), {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      });
      
      // Check if the response is OK
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      
      // Parse the JSON response
      const data = await response.json();
      
      // Append the response object to results
      results.push(data);
      
      // Write to file continuously
      const jsonString = JSON.stringify(data, null, 2);
      const prefix = results.length > 1 ? ',\n' : '';
      fs.appendFileSync(filePath, prefix + jsonString, 'utf8');
      
      // Check if there's a next_page key to continue pagination
      if (data.next_page && typeof data.next_page === 'string') {
        nextPage = data.next_page;
      } else {
        // No more pages, break the loop
        break;
      }
    }
    
    // Close the JSON array
    fs.appendFileSync(filePath, '\n]', 'utf8');
    
  } catch (error) {
    // Log the error and close the JSON array properly
    console.error('Error during pagination:', error.message);
    // Close the JSON array even on error
    fs.appendFileSync(filePath, '\n]', 'utf8');
    // Return the results collected before the error
    return results;
  }
  
  return results;
}

/**
 * CLI function to fetch data with user-provided token
 */
async function runCLI() {
  const readline = require('readline');
  
  const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
  });
  
  // Prompt for token
  rl.question('Enter your API token: ', async (token) => {
    rl.close();
    
    if (!token || token.trim() === '') {
      console.error('Error: Token is required');
      process.exit(1);
    }
    
    // TODO: Replace with your actual API URL
    const url = 'https://api.karenai.click/swechallenge/list';
    
    // TODO: Replace with your desired output file path
    const outputFile = './output.json';
    
    console.log(`\nFetching data from: ${url}`);
    console.log(`Writing results to: ${outputFile}`);
    console.log('Please wait...\n');
    
    try {
      const results = await fetchAllPages(url, token.trim(), outputFile);
      
      console.log(`✓ Successfully fetched ${results.length} pages`);
      console.log(`✓ Results written to: ${outputFile}`);
    } catch (error) {
      console.error('✗ Failed to fetch data:', error.message);
      process.exit(1);
    }
  });
}

// Run CLI if executed directly
if (require.main === module) {
  runCLI();
}

