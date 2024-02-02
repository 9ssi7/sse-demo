const eventSource = new EventSource("http://localhost:8080/events");

const yLira = document.getElementById('y-lira');
const xLira = document.getElementById('x-lira');
const zLira = document.getElementById('z-lira');

const setNewPrice = (elementId, oldPrice, newPrice) => {
    const rates = [...document.querySelectorAll(`#${elementId} span#rate`)];
    rates.forEach(rate => {
        rate.innerHTML = newPrice;
    })
    const arrowProvider = document.querySelector(`#${elementId} #rate-with-arrow`);
    const arrow = document.querySelector(`#${elementId} #rate-with-arrow span#arrow`);
    if (oldPrice === newPrice) {
        arrow.innerHTML = '';
        arrowProvider.classList.remove('up');
        arrowProvider.classList.remove('down');
    }else if (oldPrice > newPrice) {
        arrow.innerHTML = '&#8595;';
        arrowProvider.classList.remove('up');
        arrowProvider.classList.add('down');
    }else {
        arrow.innerHTML = '&#8593;';
        arrowProvider.classList.remove('down');
        arrowProvider.classList.add('up');
    }
}

eventSource.onmessage = (event) => {
    try {
        const data = JSON.parse(event.data);
        if(!data.old || !data.new) return;
        const { old, new: newPrice } = data;
        if (old.yLira !== newPrice.yLira) {
            setNewPrice('y-lira', old.yLira, newPrice.yLira);
        }
        if (old.xLira !== newPrice.xLira) {
            setNewPrice('x-lira', old.xLira, newPrice.xLira);
        }
        if (old.zLira !== newPrice.zLira) {
            setNewPrice('z-lira', old.zLira, newPrice.zLira);
        }
        document.querySelector("#last-updated").innerHTML = new Date().toLocaleTimeString();
    }catch(e) {}
};