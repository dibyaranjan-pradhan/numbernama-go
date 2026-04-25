(() => {
  const selectedClass = 'selected-box';
  const removedText = ' ';
  const rowLen = 9;
  let ws;
  let gameType = 'elem1to18';
  let selectedElems = [];
  let selectedKeys = [];

  const $ = (sel) => document.querySelector(sel);
  const $$ = (sel) => document.querySelectorAll(sel);

  function wsURL() {
    const p = location.protocol === 'https:' ? 'wss' : 'ws';
    return `${p}://${location.host}/ws/gameplay`;
  }

  function send(event, payload) {
    if (!ws || ws.readyState !== WebSocket.OPEN) return;
    ws.send(JSON.stringify({ event, payload: payload === undefined ? null : payload }));
  }

  function connect() {
    ws = new WebSocket(wsURL());
    ws.onmessage = (ev) => {
      let msg;
      try {
        msg = JSON.parse(ev.data);
      } catch {
        return;
      }
      const { event, payload } = msg;
      switch (event) {
        case 'go_gameplay_connected':
          console.log('connected', payload);
          break;
        case 'initiateGamePlay':
          onInitiate(payload);
          break;
        case 'match':
          onMatch(payload);
          break;
        case 'check':
          onCheck(payload);
          break;
        case 'clear':
          onClear(payload);
          break;
        case 'undo':
          onUndo(payload);
          break;
        default:
          break;
      }
    };
    ws.onclose = () => {
      setTimeout(connect, 1500);
    };
  }

  function onInitiate(data) {
    if (data.error) {
      console.error(data.error);
      return;
    }
    renderBoard(data.gameplayArray);
  }

  function renderBoard(gameplayArray) {
    const root = $('.numberDiv');
    root.innerHTML = '';
    gameplayArray.forEach((row, i) => {
      const ul = document.createElement('ul');
      ul.className = 'list-group list-group-horizontal';
      ul.dataset.row = String(i);
      row.forEach((num, j) => {
        const li = document.createElement('li');
        li.className = `list-group-item n${num}`;
        li.dataset.col = String(j);
        li.textContent = String(num);
        ul.appendChild(li);
      });
      root.appendChild(ul);
    });
  }

  function onMatch(data) {
    if (data.err && data.err.length) console.log(data.err.join(' '));
    if (data.matched && data.selectedElems) {
      data.selectedElems.forEach(([x, y]) => {
        const li = $(`*[data-row='${x}']`)?.querySelector(`*[data-col='${y}']`);
        if (li) {
          li.classList.remove(selectedClass);
          li.textContent = removedText;
        }
      });
      $('.undo').classList.remove('disabled');
    }
    selectedElems = [];
    selectedKeys = [];
  }

  function onCheck(response) {
    const ne = Array.isArray(response.newElems) ? [...response.newElems] : [];
    let html = '';
    if (response.lastRowElemCount > 0) {
      for (let i = 0; i < response.lastRowElemCount; i++) {
        const num = ne.shift();
        html += `<li class="list-group-item n${num}" data-col="${rowLen - response.lastRowElemCount + i}">${num}</li>`;
      }
      const rowEl = $(`*[data-row='${response.lastRow}']`);
      if (rowEl) rowEl.insertAdjacentHTML('beforeend', html);
    }
    let lastRow = response.lastRow;
    while (ne.length > 0) {
      const chunk = ne.splice(0, rowLen);
      lastRow += 1;
      let rowHtml = `<ul class="list-group list-group-horizontal" data-row="${lastRow}">`;
      chunk.forEach((num, i) => {
        rowHtml += `<li class="list-group-item n${num}" data-col="${i}">${num}</li>`;
      });
      rowHtml += `</ul>`;
      $('.numberDiv').insertAdjacentHTML('beforeend', rowHtml);
    }
  }

  function onClear(response) {
    const [flag, cleared] = response;
    if (!flag) return;
    let rowCount = -1;
    $$('.numberDiv > ul').forEach((ul) => {
      const dr = +ul.dataset.row;
      if (cleared.includes(dr)) ul.remove();
      else {
        rowCount += 1;
        ul.dataset.row = String(rowCount);
      }
    });
    $('.undo').classList.add('disabled');
  }

  function onUndo(response) {
    const [ok, lastMatched] = response;
    if (!ok || !lastMatched) return;
    const c1 = $(`*[data-row='${lastMatched.elem1.x}']`)?.querySelector(`*[data-col='${lastMatched.elem1.y}']`);
    const c2 = $(`*[data-row='${lastMatched.elem2.x}']`)?.querySelector(`*[data-col='${lastMatched.elem2.y}']`);
    if (c1) c1.textContent = String(lastMatched.val1);
    if (c2) c2.textContent = String(lastMatched.val2);
    $('.undo').classList.add('disabled');
  }

  document.addEventListener('click', (e) => {
    const t = e.target;
    if (t.matches('.modal-footer .btn')) {
      e.preventDefault();
      gameType = t.id;
      send('initiateGamePlay', { resetFlag: false, gameType });
      document.querySelector('button.gameModal')?.remove();
    }
    if (t.matches('.restart')) {
      e.preventDefault();
      send('initiateGamePlay', { resetFlag: true, gameType: null });
    }
    if (t.matches('.check')) {
      e.preventDefault();
      send('check', {});
    }
    if (t.matches('.clear')) {
      e.preventDefault();
      send('clear', {});
    }
    if (t.matches('.undo') && !t.classList.contains('disabled')) {
      e.preventDefault();
      send('undo', {});
    }
  });

  document.querySelector('.numberDiv')?.addEventListener('click', (e) => {
    const li = e.target.closest('.numberDiv ul li');
    if (!li) return;
    e.preventDefault();
    const text = li.textContent.trim();
    if (text === '') return;
    const ul = li.parentElement;
    const x = +ul.dataset.row;
    const y = +li.dataset.col;
    const key = `${x},${y}`;
    if (!li.classList.contains(selectedClass)) {
      li.classList.add(selectedClass);
      li.dataset.sel = String(selectedElems.length);
      selectedElems.push([x, y]);
      selectedKeys.push(key);
    } else {
      const idx = selectedKeys.indexOf(key);
      if (idx >= 0) {
        selectedElems.splice(idx, 1);
        selectedKeys.splice(idx, 1);
      }
      li.classList.remove(selectedClass);
      delete li.dataset.sel;
    }
    if (selectedElems.length === 2) {
      $$('.selected-box').forEach((el) => {
        el.classList.remove(selectedClass);
        delete el.dataset.sel;
      });
      $('.undo').classList.remove('disabled');
      send('match', selectedElems);
    }
  });

  connect();
})();
