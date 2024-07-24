
document.addEventListener('DOMContentLoaded', function () {
    let leasesData = [];

    const fetchLeases = async () => {
        try {
            const response = await fetch('/api/leases');
            const result = await response.json();

            if (result.code !== 0) {
                throw new Error(`Error fetching data: ${result.message}`);
            }

            leasesData = result.data;
            displayLeases(leasesData);
        } catch (error) {
            console.error('Error fetching data:', error);
            alert('Failed to fetch leases data. Please try again later.');
        }
    };

    const displayLeases = (leases) => {
        const cardGrid = document.querySelector('.card-grid');
        cardGrid.innerHTML = '';

        leases.forEach(lease => {
            const card = document.createElement('div');
            card.className = 'card';
            
            const ip = document.createElement('h2');
            ip.textContent = `IP: ${lease.ip}`;
            card.appendChild(ip);
            
            const bindingState = document.createElement('p');
            bindingState.innerHTML = `<strong>Binding State:</strong> ${lease.state}`;
            card.appendChild(bindingState);
            
            const hardwareEthernet = document.createElement('p');
            hardwareEthernet.innerHTML = `<strong>Hardware Ethernet:</strong> ${lease.hardware}`;
            card.appendChild(hardwareEthernet);
            
            if (lease.vendor) {
                const vendorClassIdentifier = document.createElement('p');
                vendorClassIdentifier.innerHTML = `<strong>Vendor:</strong> ${lease.vendor}`;
                card.appendChild(vendorClassIdentifier);
            }

            if (lease.host) {
                const hostname = document.createElement('p');
                hostname.innerHTML = `<strong>Hostname:</strong> ${lease.host}`;
                card.appendChild(hostname);
            }
            
            cardGrid.appendChild(card);
        });

        document.getElementById('total-ips').textContent = leases.length;
    };

    const applyFilter = () => {
        const stateFilter = document.getElementById('state-filter').value;
        const filteredLeases = leasesData.filter(lease => {
            return stateFilter === '' || lease.state === stateFilter;
        });
        displayLeases(filteredLeases);
    };

    document.getElementById('apply-filter').addEventListener('click', applyFilter);

    fetchLeases();
});
