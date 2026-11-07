// Поле с выпадающим списком: варианты с описанием, поиск по вводу, своё значение остаётся как есть.
window.Picker = (() => {
  const pickers = [];
  let open = null;

  function esc(s) { return String(s).replace(/[&<>"]/g, c => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;' }[c])); }

  function create(wrap, options, customNote) {
    const input = wrap.querySelector('input');
    const desc = wrap.querySelector('.pick-desc');
    const panel = document.createElement('div');
    panel.className = 'pick-panel';
    panel.hidden = true;
    wrap.appendChild(panel);
    const byName = new Map(options.map(o => [o.name.toLowerCase(), o]));
    let active = -1;

    function known() { return byName.get(input.value.trim().toLowerCase()); }

    function showDesc() {
      const o = known();
      if (o) {
        desc.innerHTML = `<b>${esc(o.desc || '')}</b>${o.details ? ' ' + esc(o.details) : ''}`;
      } else if (input.value.trim()) {
        desc.innerHTML = `<i>Свой вариант.</i> ${esc(customNote || '')}`;
      } else {
        desc.innerHTML = '';
      }
    }

    function render() {
      const q = input.value.trim().toLowerCase();
      const list = q && !known() ? options.filter(o => o.name.toLowerCase().includes(q) || (o.desc || '').toLowerCase().includes(q)) : options;
      let html = '';
      if (q && !known()) html += `<div class="pick-item pick-custom" data-custom="1"><span class="pick-name">Своё: «${esc(input.value.trim())}»</span><span class="pick-sub">${esc(customNote || 'Оставить как есть')}</span></div>`;
      html += list.map(o => `<div class="pick-item${o.name === input.value ? ' sel' : ''}" data-name="${esc(o.name)}"><span class="pick-name">${esc(o.name)}</span><span class="pick-sub">${esc(o.desc || '')}</span></div>`).join('');
      if (!html) html = '<div class="pick-empty">Ничего не нашлось, но своё значение можно оставить</div>';
      panel.innerHTML = html;
      active = -1;
    }

    function show() {
      if (open && open !== api) open.hide();
      open = api;
      render();
      panel.hidden = false;
      wrap.classList.add('open');
    }
    function hide() {
      panel.hidden = true;
      wrap.classList.remove('open');
      if (open === api) open = null;
      showDesc();
    }
    function choose(name) {
      input.value = name;
      input.dispatchEvent(new Event('input', { bubbles: true }));
      hide();
    }
    function move(d) {
      const items = [...panel.querySelectorAll('.pick-item')];
      if (!items.length) return;
      active = (active + d + items.length) % items.length;
      items.forEach((el, i) => el.classList.toggle('active', i === active));
      items[active].scrollIntoView({ block: 'nearest' });
    }

    input.addEventListener('focus', () => { show(); input.select(); });
    input.addEventListener('click', show);
    input.addEventListener('input', () => { if (panel.hidden) show(); else render(); showDesc(); });
    input.addEventListener('keydown', e => {
      if (panel.hidden) return;
      if (e.key === 'ArrowDown') { e.preventDefault(); move(1); }
      else if (e.key === 'ArrowUp') { e.preventDefault(); move(-1); }
      else if (e.key === 'Enter') { e.preventDefault(); const el = panel.querySelector('.pick-item.active'); if (el && el.dataset.name) choose(el.dataset.name); else hide(); }
      else if (e.key === 'Escape') { hide(); }
    });
    panel.addEventListener('pointerdown', e => {
      const el = e.target.closest('.pick-item');
      if (!el) return;
      e.preventDefault();
      if (el.dataset.name) choose(el.dataset.name); else hide();
    });
    wrap.querySelector('.pick-toggle').addEventListener('click', () => { if (panel.hidden) { input.focus(); show(); } else hide(); });
    input.addEventListener('blur', () => setTimeout(() => { if (!wrap.contains(document.activeElement)) hide(); }, 120));

    const api = { hide, refresh: showDesc };
    showDesc();
    pickers.push(api);
    return api;
  }

  document.addEventListener('pointerdown', e => { if (open && !e.target.closest('.pick')) open.hide(); });

  return {
    init(field, options, customNote) {
      const wrap = document.querySelector(`.pick[data-pick="${field}"]`);
      if (wrap) create(wrap, options, customNote);
    },
    refresh() { pickers.forEach(p => p.refresh()); }
  };
})();
