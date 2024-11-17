function getPoints() {
    var receiptData = document.getElementById('receiptData').value;
    var resultDiv = document.getElementById('result');
    var errorDiv = document.getElementById('error');
    resultDiv.style.display = 'none';
    errorDiv.style.display = 'none';
    try {
        // Parse the JSON input
        var receipt = JSON.parse(receiptData);
        
        // Send the receipt data to the backend
        fetch('http://localhost:8080/receipts/process', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(receipt),
        })
        .then(response => response.json())
        .then(data => {
            var receiptId = data.id;
            fetch(`http://localhost:8080/receipts/${receiptId}/points`)
            .then(response => response.json())
            .then(data => {
                resultDiv.style.display = 'block';
                resultDiv.innerHTML = `Points: ${data.points}`;
            })
            .catch(error => {
                errorDiv.style.display = 'block';
                errorDiv.innerHTML = 'Error fetching points. Please try again later.';
            });
        })
        .catch(error => {
            errorDiv.style.display = 'block';
            errorDiv.innerHTML = 'Error processing receipt. Please check the JSON format and try again.';
        });
    } catch (error) {
        errorDiv.style.display = 'block';
        errorDiv.innerHTML = 'Invalid JSON format.';
    }
}