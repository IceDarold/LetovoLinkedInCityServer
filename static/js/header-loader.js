async function loadHeader() {
    const container = document.getElementById('header-container');
    if (!container) {
        console.warn("Элемент #header-container не найден. Шапка не будет загружена.");
        return;
    }
    try {
        const response = await fetch('/static/html/header.html');
        const html = await response.text();
        container.innerHTML = html;
        setupPlayerCount();
    } catch (e) {
        console.error("Ошибка загрузки шапки:", e);
    }
}


async function setupPlayerCount() {
    async function updateStatus() {
        try {
            const response = await fetch('/api/status');
            const data = await response.json();
            const playerCountElement = document.getElementById('playerCount');
            if (playerCountElement) {
                playerCountElement.innerText = "Игроков онлайн: " + data.count;
            }
        } catch {
            const playerCountElement = document.getElementById('playerCount');
            if (playerCountElement) {
                playerCountElement.innerText = "Статус недоступен";
            }
        }
    }
    await updateStatus();
    setInterval(updateStatus, 5000);
}

// Загружаем всё при старте страницы
document.addEventListener('DOMContentLoaded', () => {
    loadHeader();
    //loadPlayersTable();
});
