async function updatePlayersTable() {
    const table = document.getElementById('players');
    if (!table) {
        console.warn("Элемент #players не найден на странице.");
        return;
    }

    try {
        const response = await fetch('/api/status');
        const data = await response.json();

        table.innerHTML = "<tr><th>Player ID</th><th>Position</th></tr>"; // очищаем таблицу

        for (const player of data.players) {
            const row = document.createElement('tr');
            row.innerHTML = `<td>${player.playerId}</td><td>(${player.position.x.toFixed(2)}, ${player.position.y.toFixed(2)}, ${player.position.z.toFixed(2)})</td>`;
            table.appendChild(row);
        }
    } catch (e) {
        console.error("Не удалось загрузить статус игроков:", e);
    }
}

document.addEventListener('DOMContentLoaded', () => {
    updatePlayersTable();
    setInterval(updatePlayersTable, 5000); // Обновлять таблицу каждые 5 секунд
});
