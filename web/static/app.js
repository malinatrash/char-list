window.App = (() => {
  const ed = document.body.dataset.ed;
  const form = document.getElementById('sheet');
  if (!ed || !form) return {};
  const key = 'char-list:' + ed;
  const hooks = { renderDerived: () => {}, afterApply: () => {} };

  function collect() {
    const c = { scores: {} };
    form.querySelectorAll('[data-field]').forEach(el => {
      const f = el.dataset.field;
      if (f === 'score') { c.scores[el.dataset.ab] = Number(el.value) || 0; return; }
      if (el.type === 'checkbox') { (c[f] = c[f] || []); if (el.checked) c[f].push(el.value); return; }
      c[f] = el.type === 'number' ? Number(el.value) || 0 : el.value;
    });
    return c;
  }

  function apply(c) {
    form.querySelectorAll('[data-field]').forEach(el => {
      const f = el.dataset.field;
      if (f === 'score') { el.value = (c.scores && c.scores[el.dataset.ab]) || ''; return; }
      if (el.type === 'checkbox') { el.checked = Array.isArray(c[f]) && c[f].includes(el.value); return; }
      if (c[f] !== undefined && c[f] !== null) el.value = c[f];
    });
    hooks.afterApply(c);
  }

  function save() { try { localStorage.setItem(key, JSON.stringify(collect())); } catch (e) {} }
  function load() { try { const s = localStorage.getItem(key); return s ? JSON.parse(s) : null; } catch (e) { return null; } }

  async function post(url, body) {
    const r = await fetch(url, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body) });
    if (!r.ok) throw new Error(await r.text());
    return r;
  }

  let t;
  async function derive() {
    const d = await (await post(`/api/${ed}/derive`, collect())).json();
    hooks.renderDerived(d, collect());
    return d;
  }
  function deriveSoon() { clearTimeout(t); t = setTimeout(() => { save(); derive().catch(console.error); }, 250); }

  async function random(fields, title) {
    const res = await (await post(`/api/${ed}/random`, { char: collect(), fields })).json();
    if (res.rolls && res.rolls.length) await Dice.show(res.rolls, title);
    apply(res.char);
    hooks.renderDerived(res.derived, res.char);
    save();
    return res;
  }

  function toast(msg) {
    const el = document.createElement('div'); el.className = 'toast'; el.textContent = msg;
    document.body.appendChild(el); setTimeout(() => el.remove(), 3500);
  }

  async function pdf(mode) {
    const btns = document.querySelectorAll('.pdf-btn'); btns.forEach(b => b.disabled = true);
    try {
      const r = await post(`/pdf/${ed}${mode === 'download' ? '?download=1' : ''}`, collect());
      const blob = await r.blob();
      const url = URL.createObjectURL(blob);
      const name = (collect().name || 'character') + '.pdf';
      if (mode === 'download') {
        const a = document.createElement('a'); a.href = url; a.download = name; document.body.appendChild(a); a.click(); a.remove();
      } else if (mode === 'print') {
        let fr = document.getElementById('pdf-frame');
        if (!fr) { fr = document.createElement('iframe'); fr.id = 'pdf-frame'; fr.className = 'pdf-wrap'; document.body.appendChild(fr); }
        fr.onload = () => { try { fr.contentWindow.focus(); fr.contentWindow.print(); } catch (e) { window.open(url, '_blank'); } };
        fr.src = url;
        toast('Отправляю на печать. Если диалог не открылся — нажми «Открыть PDF».');
      } else {
        window.open(url, '_blank');
      }
    } catch (e) { toast('Ошибка: ' + e.message); }
    finally { btns.forEach(b => b.disabled = false); }
  }

  form.addEventListener('input', deriveSoon);
  form.addEventListener('change', deriveSoon);
  form.addEventListener('submit', e => e.preventDefault());
  form.querySelectorAll('[data-rnd]').forEach(b => b.addEventListener('click', () => {
    const fields = b.dataset.rnd.split(',');
    api.random(fields, b.dataset.title || 'Бросок').catch(e => toast(e.message));
  }));
  document.getElementById('rnd-all')?.addEventListener('click', () => api.random(null, 'Судьба решает всё').catch(e => toast(e.message)));
  document.getElementById('btn-print')?.addEventListener('click', () => pdf('print'));
  document.getElementById('btn-open')?.addEventListener('click', () => pdf('open'));
  document.getElementById('btn-download')?.addEventListener('click', () => pdf('download'));
  document.getElementById('btn-clear')?.addEventListener('click', () => {
    if (!confirm('Очистить лист?')) return;
    form.reset(); form.querySelectorAll('[data-field]').forEach(el => { if (el.type === 'checkbox') el.checked = false; else if (!el.dataset.keep) el.value = ''; });
    save(); derive().catch(() => {});
  });

  const api = {
    ed, form, hooks, collect, apply, random, derive, save, toast,
    init() { const c = load(); if (c) apply(c); derive().catch(console.error); }
  };
  return api;
})();
