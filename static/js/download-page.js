async function updateDownloadLink() {
    try {
        const response = await fetch('/api/latest-version');
        const data = await response.json();
        
        const downloadButton = document.querySelector('.download-button');
        downloadButton.href = '/static/downloads/' + data.fileName;
        downloadButton.download = data.fileName;
    } catch (error) {
        console.error('Ошибка получения версии клиента:', error);
    }
}

updateDownloadLink();