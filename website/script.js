const chatbox = document.querySelector('.chatbox');
const inputbox = document.querySelector('.inputbox');
const container = document.querySelector('.container');
const input = document.querySelector('input[type="text"]');
const button = document.querySelector('button');
const clearButton = document.querySelector('#clear-button');
const text_stuff = document.querySelector('.text_stuff')
const body = document.body;

// Function to create and display a message in the chatbox
function displayMessage(message, sender) {
  const messageElem = document.createElement('div');
  messageElem.classList.add('message', sender);
  messageElem.textContent = message;
  chatbox.appendChild(messageElem);
  chatbox.scrollTop = chatbox.scrollHeight;
}

// Function to send a message
function sendMessage() {
  const message = input.value.trim();
  if (message !== '') {
    displayMessage(message, 'sent');
    input.value = '';
    // Add code here to send the message to the chatbot backend
    fetch('http://localhost:8000/api/chatbot',{
      method: 'POST',
      headers: {
        'Content-Type':'application/json'
      },
      body: JSON.stringify({
        message: message
      })
    })
    .then(response => response.json())
    .then(data => displayMessage(data.message, 'received'))
    .catch(error => console.error(error))
  }
}

// Event listener for the send button
button.addEventListener('click', sendMessage);

// Event listener for the Enter key in the input field
input.addEventListener('keydown', (e) => {
  if (e.key === 'Enter') {
    sendMessage();
  }
});

// Function to toggle light and dark mode
function toggleMode() {
  body.classList.toggle('dark');
  container.classList.toggle('dark');
  chatbox.classList.toggle('dark');
  inputbox.classList.toggle('dark');
  text_stuff.classList.toggle('dark');
}  

// Event listener for the mode toggle button
const modeToggle = document.getElementById('mode-toggle');
modeToggle.addEventListener('click', toggleMode);

// Function to clear chat
function clearChat() {
  chatbox.innerHTML = "";
}

// Event listener for the clear chat button
clearButton.addEventListener('click', clearChat);
