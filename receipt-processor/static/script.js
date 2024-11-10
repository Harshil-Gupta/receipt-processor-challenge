function getPoints() {
    const receiptData = document.getElementById('receiptData').value;
    const resultDiv = document.getElementById('result');
    const errorDiv = document.getElementById('error');
    resultDiv.style.display = 'none';
    errorDiv.style.display = 'none';

    try {
        const receipt = JSON.parse(receiptData);
        
        fetch('http://localhost:8080/receipts/process', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(receipt),
        })
        .then(response => response.json())
        .then(data => {
            const receiptId = data.id;
            return fetch(`http://localhost:8080/receipts/${receiptId}/points`);
        })
        .then(response => response.json())
        .then(pointsData => {
            resultDiv.style.display = 'block';
            resultDiv.innerHTML = `Points: ${pointsData.points}`;
        })
        .catch(error => {
            errorDiv.style.display = 'block';
            errorDiv.innerHTML = 'Error processing receipt or fetching points. Please try again.';
        });
    } catch (error) {
        errorDiv.style.display = 'block';
        errorDiv.innerHTML = 'Invalid JSON format.';
    }
}
function getPointsById() {
    // Get the value from the input field
    const receiptId = document.getElementById('transactionIdInput').value;
    const fetchResultDiv = document.getElementById('fetchResult');
    fetchResultDiv.style.display = 'none'; // Hide result initially

    if (!receiptId) {
        // If no transaction ID is entered, display an error message
        fetchResultDiv.style.display = 'block';
        fetchResultDiv.innerHTML = 'Please enter a transaction ID.';
        return;
    }

    // Make a fetch request to the server with the entered receipt ID
    fetch(`http://localhost:8080/receipts/${receiptId}/points`)
    .then(response => {
        if (!response.ok) {
            throw new Error('Transaction ID not found');
        }
        return response.json();
    })
    .then(data => {
        // Display points for the given transaction ID
        fetchResultDiv.style.display = 'block';
        fetchResultDiv.innerHTML = `Points for transaction ID ${receiptId}: ${data.points}`;
    })
    .catch(error => {
        // Display any errors if they occur
        fetchResultDiv.style.display = 'block';
        fetchResultDiv.innerHTML = error.message;
    });
}