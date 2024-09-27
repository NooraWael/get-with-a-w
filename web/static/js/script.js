document.addEventListener('DOMContentLoaded', function() {
    const form = document.querySelector('form');
    
    form.addEventListener('submit', function(event) {
        event.preventDefault();
        
        const url = document.getElementById('url').value;
        fetch('/download', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: `url=${encodeURIComponent(url)}`
        })
        .then(response => response.text())
        .then(data => {
            alert('Download initiated successfully!');
            console.log(data); // Optionally output the server response to the console
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Failed to initiate download. See console for more details.');
        });
    });
});
