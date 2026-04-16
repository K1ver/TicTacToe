const API_BASE = 'http://localhost:8080/api';

let currentUser = null;
let currentGame = null;
let gameInterval = null;

// Управление экранами
function showScreen(screenId) {
    document.querySelectorAll('.screen').forEach(screen => {
        screen.classList.remove('active');
    });
    document.getElementById(screenId).classList.add('active');
    
    if (screenId === 'browser-screen') {
        loadGames();
    }
}

// Управление табами авторизации
function showAuthTab(tab) {
    document.querySelectorAll('.auth-form').forEach(form => {
        form.classList.remove('active');
    });
    document.querySelectorAll('.tab-btn').forEach(btn => {
        btn.classList.remove('active');
    });

    document.getElementById(`${tab}-form`).classList.add('active');
    event.target.classList.add('active');
}

// Регистрация с отладкой
document.getElementById('signup-form').addEventListener('submit', async (e) => {
    e.preventDefault();

    const username = document.getElementById('signup-username').value;
    const password = document.getElementById('signup-password').value;

    try {
        const response = await fetch(`${API_BASE}/auth/signUp`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                login: username,
                password: password
            })
        });

        const responseText = await response.text();

        if (response.ok) {
            showAuthMessage('Регистрация успешна! Теперь войдите.', 'success');
        } else {
            showAuthMessage(`Ошибка регистрации: ${responseText}`, 'error');
        }
    } catch (error) {
        showAuthMessage('Ошибка соединения с сервером', 'error');
    }
});

// Вход (JWT) — ОБНОВЛЯЕМ ЛОГИКУ получения userId, чтобы видеть причину ошибок
document.getElementById('signin-form').addEventListener('submit', async (e) => {
    e.preventDefault();

    const username = document.getElementById('signin-username').value;
    const password = document.getElementById('signin-password').value;

    try {
        const response = await fetch(`${API_BASE}/auth/signIn`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                login: username,
                password: password
            })
        });

        if (response.ok) {
            const data = await response.json();
            const accessToken = data.accessToken;
            const refreshToken = data.refreshToken;

            // Получаем userId (и остальные поля) по accessToken через дополнительный запрос
            const userRes = await fetch(`${API_BASE}/user/me`, {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${accessToken}`
                }
            });

            let user = { username, accessToken, refreshToken };
            if (userRes.ok) {
                const userJson = await userRes.json();
                // Логируем полученный ответ для отладки
                console.log('[DEBUG] /user/me response:', userJson);
                user.userId = userJson.id || userJson.userId || null;
                if (!user.userId) {
                    alert('Ошибка входа: сервер не вернул userId. Смотрите ответ в консоли (F12), отправьте его разработчику.');
                }
            } else {
                const errText = await userRes.text();
                console.error('[DEBUG] /user/me error:', errText);
                user.userId = null;
                alert('Ошибка получения информации о пользователе после входа. Смотрите ответ в консоли.');
            }
            currentUser = user;
            localStorage.setItem('user', JSON.stringify(currentUser));
            document.getElementById('current-user').textContent = username;
            showScreen('main-menu');
        } else {
            const responseText = await response.text();
            showAuthMessage(`Неверное имя пользователя или пароль: ${responseText}`, 'error');
        }
    } catch (error) {
        showAuthMessage('Ошибка соединения с сервером', 'error');
    }
});

function showAuthMessage(message, type) {
    const messageEl = document.getElementById('auth-message');
    messageEl.textContent = message;
    messageEl.className = type;
}

// Создание игры
async function createGame() {
    const side = parseInt(document.getElementById('side-select').value);
    const vsCpu = document.getElementById('cpu-game').checked;
    
    try {
        const response = await fetch(`${API_BASE}/game`, {
            method: 'POST',
            headers: getAuthHeaders(),
            body: JSON.stringify({
                side: side,
                gameWithCPU: vsCpu
            })
        });
        
        if (response.ok) {
            const data = await response.json();
            currentGame = data.id;
            initializeGameBoard();
            showScreen('game-screen');
            startGamePolling();
        } else {
            alert('Ошибка создания игры');
        }
    } catch (error) {
        alert('Ошибка соединения с сервером');
    }
}

// Загрузка списка игр
async function loadGames() {
    try {
        const response = await fetch(`${API_BASE}/game/browse`, {
            headers: getAuthHeaders()
        });
        
        const gamesList = document.getElementById('games-list');
        if (response.ok) {
            const games = await response.json();
            
            if (games.length === 0) {
                gamesList.innerHTML = '<div class="empty-state">Нет доступных игр</div>';
                return;
            }
            
            gamesList.innerHTML = '';
            games.forEach(game => {
                const gameItem = document.createElement('div');
                gameItem.className = 'game-item';
                gameItem.innerHTML = `
                    <div class="game-id">Игра ${game.id.slice(0, 8)}...</div>
                    <div class="game-side">Свободная сторона: ${game.side === 1 ? 'Крестики (X)' : 'Нолики (O)'}</div>
                `;
                gameItem.onclick = () => joinGame(game.id);
                gamesList.appendChild(gameItem);
            });
        } else {
            gamesList.innerHTML = '<div class="empty-state">Ошибка загрузки игр</div>';
        }
    } catch (error) {
        document.getElementById('games-list').innerHTML = '<div class="empty-state">Ошибка соединения</div>';
    }
}

// Присоединение к игре
async function joinGame(gameId) {
    try {
        const response = await fetch(`${API_BASE}/game/${gameId}/join`, {
            method: 'POST',
            headers: getAuthHeaders()
        });
        
        if (response.ok) {
            const gameData = await response.json();
            currentGame = gameId;
            initializeGameBoard();
            showScreen('game-screen');
            startGamePolling();
        } else {
            alert('Не удалось присоединиться к игре');
        }
    } catch (error) {
        alert('Ошибка соединения с сервером');
    }
}

// Инициализация игрового поля
function initializeGameBoard() {
    const board = document.getElementById('game-board');
    board.innerHTML = '';
    
    for (let i = 0; i < 9; i++) {
        const cell = document.createElement('button');
        cell.className = 'cell';
        cell.dataset.index = i;
        cell.onclick = () => makeMove(i);
        board.appendChild(cell);
    }
}

// Ход игрока
async function makeMove(cellIndex) {
    if (!currentGame) return;
    
    try {
        // Получаем текущее состояние игры
        const gameState = await getGameState();
        if (!gameState || gameState.status !== 1) return;
        
        // Проверяем, что это ход текущего игрока
        if (gameState.currentTurn !== currentUser.userId) {
            alert('Сейчас не ваш ход!');
            return;
        }
        
        // Создаем новое поле с ходом
        const row = Math.floor(cellIndex / 3);
        const col = cellIndex % 3;
        const newField = JSON.parse(JSON.stringify(gameState.field));
        
        // Проверяем, что клетка свободна
        if (newField[row][col] !== 0) {
            alert('Клетка уже занята!');
            return;
        }
        
        // Устанавливаем ход
        if (gameState.playerO === currentUser.userId) {
            newField[row][col] = 2
        } else {
            newField[row][col] = 1
        }

        
        // Отправляем ход
        const response = await fetch(`${API_BASE}/game/${currentGame}`, {
            method: 'POST',
            headers: getAuthHeaders(),
            body: JSON.stringify({
                field: newField
            })
        });
        
        if (!response.ok) {
            const error = await response.text();
            alert(`Неверный ход: ${error}`);
        }
        
    } catch (error) {
        console.error('Ошибка хода:', error);
        alert('Ошибка при выполнении хода');
    }
}

// Опрос состояния игры
function startGamePolling() {
    if (gameInterval) clearInterval(gameInterval);

    gameInterval = setInterval(async () => {
        if (!currentGame) return;

        const gameState = await getGameState();
        if (gameState) {
            updateGameBoard(gameState);
            updateGameInfo(gameState);
        }
    }, 2000);
}

// Определение стороны текущего игрока
function getCurrentPlayerSide(gameState) {
    if (!currentUser || !currentUser.userId) {
        console.error('No user ID found');
        return 0;
    }

    // Очищаем userId на случай если в localStorage сохранился с кавычками
    const cleanUserId = currentUser.userId.replace(/["\n]/g, '').trim();

    console.log('Current user ID:', cleanUserId);
    console.log('Player X ID:', gameState.playerX);
    console.log('Player O ID:', gameState.playerO);
    console.log('Current player:', gameState.currentTurn);

    if (gameState.currentTurn === gameState.playerX  && gameState.playerX.toString() === cleanUserId) {
        console.log('User is playing as X');
        return 1;
    }
    if (gameState.currentTurn === gameState.playerO && gameState.playerO.toString() === cleanUserId) {
        console.log('User is playing as O');
        return 2;
    }

    console.log('User not found in game players');
    return 0;
}

// Получение состояния игры
async function getGameState() {
    try {
        const response = await fetch(`${API_BASE}/game/${currentGame}`, {
            headers: getAuthHeaders()
        });
        
        if (response.ok) {
            return await response.json();
        }
    } catch (error) {
        console.error('Ошибка получения состояния игры:', error);
    }
    return null;
}

// Обновление игрового поля
function updateGameBoard(gameState) {
    const cells = document.querySelectorAll('.cell');
    const field = gameState.field.flat();

    // Определяем сторону текущего игрока
    const playerSide = getCurrentPlayerSide(gameState);
    const isMyTurn = gameState.currentTurn === currentUser.userId;
    const isGameActive = gameState.status === 1;

    console.log('Board update:');
    console.log('Player side:', playerSide);
    console.log('Current player in game:', gameState.currentTurn);
    console.log('Is my turn:', isMyTurn);
    console.log('Is game active:', isGameActive);

    cells.forEach((cell, index) => {
        const value = field[index];
        cell.textContent = value === 1 ? 'X' : value === 2 ? 'O' : '';
        cell.className = `cell ${value === 1 ? 'x' : value === 2 ? 'o' : ''}`;

        // Блокируем клетку если:
        // - игра не активна
        // - не наш ход
        // - клетка уже занята
        const isCellOccupied = value !== 0;
        cell.disabled = !isGameActive || !isMyTurn || isCellOccupied;

        // Визуальная индикация
        if (cell.disabled) {
            cell.style.cursor = 'not-allowed';
            cell.style.opacity = isCellOccupied ? '1' : '0.5';
        } else {
            cell.style.cursor = 'pointer';
            cell.style.opacity = '1';
        }
    });
}

// Обновление информации об игре
function updateGameInfo(gameState) {
    const statusElement = document.getElementById('game-status');
    const playerXElement = document.getElementById('player-x');
    const playerOElement = document.getElementById('player-o');
    
    // Обновление статуса
    let statusText = '';
    let statusClass = '';
    
    if (gameState.status === 0) {
        statusText = '⏳ Ожидание второго игрока...';
        statusClass = 'waiting';
    } else if (gameState.status === 1) {
        const currentSymbol = gameState.currentTurn === gameState.playerX ? 'X' : 'O';
        statusText = `🎮 Ход: ${currentSymbol}`;
        statusClass = '';
    } else if (gameState.status === 2) {
        if (gameState.winnerId === "00000000-0000-0000-0000-000000000000") {
            statusText = '🤝 Ничья!';
        } else {
            const winnerSymbol = gameState.winnerId === gameState.playerX ? 'X' : 'O';
            statusText = `🏆 Победитель: ${winnerSymbol}`;
        }
        statusClass = 'finished';
    }
    
    statusElement.textContent = statusText;
    statusElement.className = `status ${statusClass}`;
    
    // Обновление информации об игроках
    playerXElement.textContent = gameState.playerX ? `Игрок X` : 'Ожидание...';
    playerOElement.textContent = gameState.playerO ? `Игрок O` : 'Ожидание...';
}

// Выход из игры
function leaveGame() {
    if (gameInterval) {
        clearInterval(gameInterval);
        gameInterval = null;
    }
    currentGame = null;
    showScreen('main-menu');
}

// Выход из системы
function logout() {
    currentUser = null;
    localStorage.removeItem('user');
    showScreen('auth-screen');
    document.getElementById('signin-form').reset();
    document.getElementById('signup-form').reset();
    showAuthTab('signin');
}

// Проверка авторизации при загрузке
window.addEventListener('load', () => {
    const savedUser = localStorage.getItem('user');
    if (savedUser) {
        currentUser = JSON.parse(savedUser);
        document.getElementById('current-user').textContent = currentUser.username;
        showScreen('main-menu');
    } else {
        showScreen('auth-screen');
    }
});

// Везде: currentUser.authHeader -> Authorization: Bearer currentUser.accessToken
// Апдейт: глобальная функция для заголовков:
function getAuthHeaders() {
    if (currentUser && currentUser.accessToken) {
        return {
            'Authorization': `Bearer ${currentUser.accessToken}`,
            'Content-Type': 'application/json'
        };
    }
    return { 'Content-Type': 'application/json' };
}

// === Новое: История игр ===
async function loadGameHistory() {
    const container = document.getElementById('game-history-list');
    container.innerHTML = '<div class="empty-state">Загрузка истории...</div>';
    try {
        const res = await fetch(`${API_BASE}/game/history`, {
            method: 'GET',
            headers: getAuthHeaders()
        });
        if (res.ok) {
            const games = await res.json();
            if (!games || games.length === 0) {
                container.innerHTML = '<div class="empty-state">История игр пуста</div>';
                return;
            }
            container.innerHTML = '';
            games.forEach(g => {
                const status = g.status === 2 ? (g.winnerId === currentUser.userId ? 'Победа' : (g.winnerId === '00000000-0000-0000-0000-000000000000' ? 'Ничья' : 'Поражение')) : (g.status === 1 ? 'В процессе' : 'Ожидание');
                const date = new Date(g.createdAt).toLocaleString('ru-RU');
                const item = document.createElement('div');
                item.className = 'game-item';
                item.innerHTML = `<div class='game-id'>${date}</div><div class='game-side'>${status}</div>`;
                container.appendChild(item);
            });
        } else {
            container.innerHTML = '<div class="empty-state">Ошибка загрузки истории</div>';
        }
    } catch {
        container.innerHTML = '<div class="empty-state">Ошибка соединения</div>';
    }
}

// === Новое: Лидерборд ===
async function loadLeaderboard() {
    const container = document.getElementById('leaderboard-list');
    container.innerHTML = '<div class="empty-state">Загрузка рейтинга...</div>';
    try {
        const res = await fetch(`${API_BASE}/game/leaderboard/10`, {
            headers: getAuthHeaders()
        });
        if (res.ok) {
            const leaders = await res.json();
            if (!leaders || leaders.length === 0) {
                container.innerHTML = '<div class="empty-state">Пока никто не сыграл</div>';
                return;
            }
            container.innerHTML = '';
            leaders.forEach((l, i) => {
                const item = document.createElement('div');
                item.className = 'game-item';
                item.innerHTML = `<div class='game-id'>#${i+1} ${l.login}</div><div class='game-side'>Winrate: ${(l.winrate*100).toFixed(1)}%</div>`;
                container.appendChild(item);
            });
        } else {
            container.innerHTML = '<div class="empty-state">Ошибка загрузки лидерборда</div>';
        }
    } catch {
        container.innerHTML = '<div class="empty-state">Ошибка соединения</div>';
    }
}

// Универсальный fetch с авто-обновлением accessToken при 401
async function apiFetch(url, opts = {}, autoRetry = true) {
    if (!opts.headers) opts.headers = {};
    const token = currentUser && currentUser.accessToken;
    if (token) opts.headers['Authorization'] = `Bearer ${token}`;
    opts.headers['Content-Type'] = 'application/json';
    let res = await fetch(url, opts);
    if (res.status === 401 && autoRetry && currentUser && currentUser.refreshToken) {
        // Пытаемся обновить accessToken через refreshToken
        const refreshRes = await fetch(`${API_BASE}/auth/updateAccessToken`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ refreshToken: currentUser.refreshToken })
        });
        if (refreshRes.ok) {
            const data = await refreshRes.json();
            currentUser.accessToken = data.accessToken;
            currentUser.refreshToken = data.refreshToken;
            localStorage.setItem('user', JSON.stringify(currentUser));
            // Повторяем исходный запрос с новым токеном
            opts.headers['Authorization'] = `Bearer ${data.accessToken}`;
            return await fetch(url, opts);
        } else {
            // Не удалось обновить — logout
            logout();
            throw new Error('Сессия истекла. Войдите заново.');
        }
    }
    return res;
}