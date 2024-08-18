// yaMap.js
ymaps.ready(init);

const zoomVal = 16;

function init() {
    const modeSwitch = document.getElementById('modeSwitch');
    const coordinateVals = document.getElementById('coordinateVals');

    // Проверяем, что элементы существуют
    if (!modeSwitch || !coordinateVals) {
        console.error('Не удалось найти элементы modeSwitch или coordinateVals на странице.');
        return;
    }

    var map = new ymaps.Map("map", {
        center: [59.9343, 30.3351], 
        zoom: zoomVal 
    });

    let placemark = null;

    let editMode = modeSwitch.checked;

    modeSwitch.addEventListener('change', function() {
        editMode = this.checked;
    });

    fetch('/points')
        .then(response => {
            if (!response.ok) {
                throw new Error('Ошибка загрузки точек с сервера');
            }
            return response.json();
        })
        .then(points => {
            if (points.length > 0) {
                const point = points[0]; 
                placemark = new ymaps.Placemark([point.lat, point.lng], {}, {
                    iconLayout: 'default#image',
                    iconImageHref: '/public/image/gas_station_icon.png',
                    iconImageSize: [30, 30], 
                    iconImageOffset: [-15, -15] 
                });
                map.geoObjects.add(placemark);
                // Устанавливаем координаты в таблице
                coordinateVals.innerText = point.lat + ', ' + point.lng;
                map.setCenter([point.lat, point.lng], zoomVal);
            }
        })
        .catch(error => {
            console.error('Произошла ошибка:', error);
        });

    map.events.add('click', function (e) {
        if (!editMode) {
            return;
        }
        var coords = e.get('coords');

        if (placemark) {
            placemark.geometry.setCoordinates(coords);
        } else {
            placemark = new ymaps.Placemark(coords, {}, {
                iconLayout: 'default#image',
                iconImageHref: '/public/image/gas_station_icon.png',
                iconImageSize: [30, 30], 
                iconImageOffset: [-15, -15]
            });
            map.geoObjects.add(placemark);
        }

        coordinateVals.innerText = coords[0] + ', ' + coords[1];

        map.setCenter(coords, zoomVal);

        fetch('/save-point', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ lat: coords[0], lng: coords[1] })
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Ошибка сохранения точки на сервере');
            }
            return response.json();
        })
        .catch(error => {
            console.error('Произошла ошибка при сохранении точки:', error);
        });
    });
}
