function getOrder() {
    const orderId = document.getElementById('orderId').value.trim();
    if (!orderId) {
        showError('Please enter an Order ID');
        return;
    }

    showLoading(true);
    hideError();
    hideOrderInfo();

    fetch(`/order/${orderId}`)
        .then(response => {
            if (!response.ok) {
                throw new Error('Order not found');
            }
            return response.json();
        })
        .then(order => {
            displayOrder(order);
            showLoading(false);
        })
        .catch(error => {
            showError(error.message);
            showLoading(false);
        });
}

function displayOrder(order) {
    // Basic information
    document.getElementById('basicInfo').innerHTML = `
        <p><strong>Order UID:</strong> ${order.order_uid}</p>
        <p><strong>Track Number:</strong> ${order.track_number}</p>
        <p><strong>Entry:</strong> ${order.entry}</p>
        <p><strong>Customer ID:</strong> ${order.customer_id}</p>
        <p><strong>Date Created:</strong> ${new Date(order.date_created).toLocaleString()}</p>
    `;

    // Delivery information
    document.getElementById('deliveryInfo').innerHTML = `
        <p><strong>Name:</strong> ${order.delivery.name}</p>
        <p><strong>Phone:</strong> ${order.delivery.phone}</p>
        <p><strong>Address:</strong> ${order.delivery.address}, ${order.delivery.city}</p>
        <p><strong>Region:</strong> ${order.delivery.region}</p>
        <p><strong>ZIP:</strong> ${order.delivery.zip}</p>
        <p><strong>Email:</strong> ${order.delivery.email}</p>
    `;

    // Payment information
    document.getElementById('paymentInfo').innerHTML = `
        <p><strong>Transaction:</strong> ${order.payment.transaction}</p>
        <p><strong>Amount:</strong> $${(order.payment.amount / 100).toFixed(2)}</p>
        <p><strong>Currency:</strong> ${order.payment.currency}</p>
        <p><strong>Provider:</strong> ${order.payment.provider}</p>
        <p><strong>Bank:</strong> ${order.payment.bank}</p>
    `;

    // Items
    const itemsHtml = order.items.map(item => `
        <div class="item">
            <p><strong>Name:</strong> ${item.name}</p>
            <p><strong>Brand:</strong> ${item.brand}</p>
            <p><strong>Price:</strong> $${(item.price / 100).toFixed(2)}</p>
            <p><strong>Quantity:</strong> ${item.total_price / item.price}</p>
            <p><strong>Total:</strong> $${(item.total_price / 100).toFixed(2)}</p>
        </div>
    `).join('');

    document.getElementById('itemsInfo').innerHTML = itemsHtml;
    document.getElementById('orderInfo').style.display = 'block';
}

function showLoading(show) {
    document.getElementById('loading').style.display = show ? 'block' : 'none';
}

function showError(message) {
    const errorDiv = document.getElementById('error');
    errorDiv.textContent = message;
    errorDiv.style.display = 'block';
}

function hideError() {
    document.getElementById('error').style.display = 'none';
}

function hideOrderInfo() {
    document.getElementById('orderInfo').style.display = 'none';
}

// Handle Enter key press
document.getElementById('orderId').addEventListener('keypress', function(e) {
    if (e.key === 'Enter') {
        getOrder();
    }
});