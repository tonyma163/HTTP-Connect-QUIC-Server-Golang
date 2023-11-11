async function sendRequest() {
    while (true) {
      try {
        const response = await fetch('http://localhost:6969/', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ type: "FRAME", message: "from web" })
        });
  
        const data = await response.json();
        console.log(data);
      } catch (error) {
        console.error('Error:', error);
      }
    }
  }
  
  sendRequest();