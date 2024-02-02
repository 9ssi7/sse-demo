const registerForm = document.getElementById('register-form');
const tryer = document.getElementById('tryer');
const tryerButton = document.getElementById('tryer-button');
const userInfo = document.getElementById('user-info');
const messages = document.getElementById('messages');

const createEventSource = (email) => {
    const eventSource = new EventSource(`http://localhost:8080/events?email=${email}`);
    eventSource.onmessage = (event) => {
        const element = document.createElement('div');
        element.innerHTML = `<p>${event.data} <span class="date">${new Date().toLocaleTimeString()}</span></p>`;
        messages.appendChild(element);
        console.log(event);
    };
}

tryerButton.addEventListener('click', () => {
    const email = document.querySelector('#user-email').innerText;
    fetch(`http://localhost:8080/api?email=${email}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({})
    })
})

registerForm.addEventListener('submit', (event) => {
    event.preventDefault();
    const email = event.target.email.value;
    const name = event.target.name.value;

    fetch('http://localhost:8080/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ email, name })
    }).then(() => {
        registerForm.style.display = 'none';
        createEventSource(email);
        tryer.style.display = 'block';
        userInfo.innerHTML = `<div><h4>User Info</h4><p>Email: <span id="user-email">${email}</span></p><p>Name: ${name}</p></div>`;
    })
});
