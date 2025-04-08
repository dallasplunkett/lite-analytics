const SERVER_URL = 'http://localhost:8080/stats';

const stat = {
    website: 'example.com',
    data: {
        event: 'pageview',
        page: window.location.pathname,
        timestamp: new Date().toISOString(),
        userAgent: navigator.userAgent
    }
};

fetch(SERVER_URL, {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json'
    },
    body: JSON.stringify(stat)
})
    .then(response => response.json())
    .then(result => {
        console.log('Stat successfully sent:', result);
    })
    .catch(error => {
        console.error('Error sending stat:', error);
    });
