const eventSource = new EventSource("http://localhost:8080/events");

eventSource.onmessage = (event) => {
    console.log(event);
};