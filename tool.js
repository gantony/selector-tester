document.getElementById('selectorForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    const url = document.getElementById('url').value;
    const token = document.getElementById('token').value;
    const selector = document.getElementById('selector').value;
    const timeFrom = document.getElementById('timeFrom').value;
    const timeTo = document.getElementById('timeTo').value;
    const resultEl = document.getElementById('result');
    const payloadPreview = document.getElementById('payloadPreview');

    // Build and show formatted payload preview
    const payloadObj = {
        "page_size": 100,
        "page_num": 0,
        "sort_by": [],
        "time_range": { "from": timeFrom, "to": timeTo },
        "selector": selector
    };
    payloadPreview.textContent = JSON.stringify(payloadObj, null, 2);

    resultEl.textContent = 'Loading...';
    let text = '';
    try {
        const response = await fetch('/proxy', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                url: url,
                payload: JSON.stringify(payloadObj),
                token: token
            })
        });
        text = await response.text();
        const data = JSON.parse(text)
        resultEl.textContent = JSON.stringify(data, null, 2);
    } catch (err) {
        resultEl.textContent = text;
    }
});

// Update payload preview live as user types
['selector', 'timeFrom', 'timeTo'].forEach(id => {
    document.getElementById(id).addEventListener('input', function() {
        const selector = document.getElementById('selector').value;
        const timeFrom = document.getElementById('timeFrom').value;
        const timeTo = document.getElementById('timeTo').value;
        const payloadObj = {
            "page_size": 100,
            "page_num": 0,
            "sort_by": [],
            "time_range": { "from": timeFrom, "to": timeTo },
            "selector": selector
        };
        document.getElementById('payloadPreview').textContent = JSON.stringify(payloadObj, null, 2);
    });
});
