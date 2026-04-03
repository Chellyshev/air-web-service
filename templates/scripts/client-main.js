ymaps.ready(init);

function init() {
  var map = new ymaps.Map("map", {
    center: [57.153, 65.5343],
    zoom: 6,
  });

  var placemark = new ymaps.Placemark([57.153, 65.5343], {
    balloonContent: "Тюмень",
  });

  map.geoObjects.add(placemark);
}
