
<script>
let currentCard = 1; // Keeps track of the currently visible card
function showCard(cardNumber) {
    const card1 = document.getElementById('card1');
    const card2 = document.getElementById('card2');
    const card3 = document.getElementById('card3'); // Reference to the map card
    const dots = document.querySelectorAll('.dot');
    // Hide all cards
    card1.classList.remove('active');
    card2.classList.remove('active');
    card3.classList.remove('active');
    // Show the selected card
    if (cardNumber === 1) {
        card1.classList.add('active');
    } else if (cardNumber === 2) {
        card2.classList.add('active');
    } else if (cardNumber === 3) {
        card3.classList.add('active');
    }
    // Update the dots
    dots.forEach((dot, index) => {
        dot.classList.remove('active');
        if (index === cardNumber - 1) {
            dot.classList.add('active');
        }
    });
    currentCard = cardNumber; // Update the current card number
}

document.addEventListener("DOMContentLoaded", function () {
    // Initial setup to show the first card
    showCard(currentCard);

    // Get the current URL path
    const path = window.location.pathname;
    // Select the navigation links
    const homeLink = document.getElementById("homeLink");
    const artistsLink = document.getElementById("artistsLink");

    // Add 'active' class based on current URL
    if (path === "/") {
        homeLink.classList.add("active"); // Add active class to Home link
    } else if (path === "/Artists") {
        artistsLink.classList.add("active"); // Add active class to Artists link
    }
});

document.addEventListener('DOMContentLoaded', function () {
    const dateContainers = document.querySelectorAll('.content');

    dateContainers.forEach(container => {
        const links = container.querySelectorAll('.date-link');
        if (links.length > 1) {
            // Add 'multiple' class to links
            links.forEach(link => {
                link.classList.add('multiple');
            });

            // Set up auto-scrolling
            let scrollAmount = 0;
            let scrollStep = 1; // Number of pixels to scroll per frame
            const scrollInterval = 15; // Time (in ms) between each scroll step
            const maxScroll = container.scrollWidth - container.clientWidth; // Calculate max scroll value
            let scrolling; // To store the interval ID

            // Auto-scrolling function
            function autoScroll() {
                container.scrollLeft = scrollAmount;
                scrollAmount += scrollStep;

                // Reverse scrolling if we reach the end or the beginning
                if (scrollAmount >= maxScroll || scrollAmount <= 0) {
                    scrollStep = -scrollStep;
                }
            }

            // Start auto-scrolling
            function startScrolling() {
                scrolling = setInterval(autoScroll, scrollInterval);
            }

            // Stop auto-scrolling
            function stopScrolling() {
                clearInterval(scrolling);
            }

            // Start scrolling initially
            startScrolling();

            // Pause on hover
            container.addEventListener('mouseenter', stopScrolling);

            // Resume scrolling when hover ends
            container.addEventListener('mouseleave', startScrolling);
        }
    });
});

const map = L.map("map").setView([51.505, -0.09], 13);
L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution:  '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
}).addTo(map);

const customIcon = L.icon({
    iconUrl: '/static/image/pin (1).png',
    iconSize: [30, 30],
});

const clusterLayer = window.L.markerClusterGroup({
    iconCreateFunction: function (cluster) {
        return L.divIcon({
            html: '<div class="cluster-div">'+ cluster.getChildCount() +'</div>'
        })
    }
});

const coordinates = JSON.parse(`{{.Coordinates}}`);
coordinates.forEach((coordinate) => {
    if (coordinate.lat && coordinate.lon && coordinate.location) {
        const marker = L.marker([coordinate.lat, coordinate.lon], { icon: customIcon});
        const popup = coordinate.location;
        marker.bindPopup("<h4>"+ popup +"<h4>");
        clusterLayer.addLayer(marker);
    } else {
        console.warn("Invalid coordinate data:", coordinate);
    }
});
map.addLayer(clusterLayer);

const bounds = L.latLngBounds(coordinates.map(c => [c.lat, c.lon]));
map.fitBounds(bounds);

</script>