async function sendRequest() {
    while (true) {
      try {
        const res = await fetch('http://127.0.0.1:7000/', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ type: "FRAME", message: "from web2" })
        });
  
        const data = await res.json();
        console.log(data);
      } catch (error) {
        console.error('Error:', error);
      }
    }
  }
  
  sendRequest();