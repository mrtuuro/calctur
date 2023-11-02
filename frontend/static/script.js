let ws;
let sending = false;
let interval;

function initMap() {
    var map = new google.maps.Map(document.getElementById('map'), {
        center: {lat: -34.397, lng: 150.644},
        zoom: 8
    });

    map.addListener('click', function(event) {
        fetch('/coords', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({lat: event.latLng.lat(), lng: event.latLng.lng()})
        });
    });
}

document.getElementById("locationBtn").addEventListener("click", function() {
    if (sending) {
        clearInterval(interval);
        ws.close();
        this.textContent = "Koordinatları Başlat";
        sending = false;
    } else {
        if ("geolocation" in navigator) {
            ws = new WebSocket("ws://localhost:8080/ws");

            interval = setInterval(function() {
                navigator.geolocation.getCurrentPosition(function(position) {
                    ws.send(JSON.stringify({
                        lat: position.coords.latitude,
                        lng: position.coords.longitude
                    }));
                });
            }, 3000);

            this.textContent = "Koordinatları Durdur";
            sending = true;
        } else {
            alert("Tarayıcınız konum servislerini desteklemiyor.");
        }
    }
});
