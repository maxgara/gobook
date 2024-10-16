package main

var homepageStr = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Simple Text Input Form</title>
    <style>
        #messageList {
            list-style-type: none;
            padding: 0;
            margin: 20px 0;
        }
        #messageList li {
            padding: 8px;
            background-color: #f0f0f0;
            margin-bottom: 5px;
            border-radius: 4px;
        }
    </style>
    <script>
        let socket;

        // Function to handle form submission
        function handleSubmit(event) {
            event.preventDefault(); // Prevent the default form submission (no redirect)
            
            const userInput = document.getElementById('userInput').value;
            
            // Send the input text to the server via WebSocket
            if (socket && socket.readyState === WebSocket.OPEN) {
                socket.send(JSON.stringify({ text: userInput }));
                document.getElementById('userInput').value = ''; // Clear the input
            } else {
                console.error('WebSocket is not open.');
            }
        }

        // Function to handle incoming messages and update the UI
        function addMessage(message) {
            const messageList = document.getElementById('messageList');
            const listItem = document.createElement('li');
            listItem.textContent = message.text; // Assuming the message is an object with a "text" field
            messageList.appendChild(listItem);
        }

        // Initialize WebSocket connection
        function initWebSocket() {
            socket = new WebSocket('ws://' + window.location.host + '/soc'); // Adjusting WebSocket URL based on host

            socket.onopen = function() {
                console.log('WebSocket connection established.');
            };

            socket.onmessage = function(event) {
                const message = JSON.parse(event.data);
                addMessage(message);
            };

            socket.onclose = function() {
                console.log('WebSocket connection closed. Reconnecting...');
                setTimeout(initWebSocket, 5000); // Attempt to reconnect after 5 seconds
            };

            socket.onerror = function(error) {
                console.error('WebSocket error:', error);
            };
        }

        // Initialize WebSocket when the page loads
        window.onload = initWebSocket;
    </script>
</head>
<body style="background-color:gray">
    <h1>Messages</h1>
    <ul id="messageList">
        <!-- Messages will be dynamically populated here via WebSocket -->
    </ul>

    <form id="textForm" onsubmit="handleSubmit(event)">
        <label for="userInput">Enter your text:</label><br><br>
        <input type="text" id="userInput" name="userInput" required><br><br>
        <input type="submit" value="Submit">
    </form>
</body>
</html>`
