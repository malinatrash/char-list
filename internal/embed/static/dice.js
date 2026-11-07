window.Dice = (() => {
  const tray = document.getElementById('dice-tray');
  const label = document.getElementById('tray-label');
  const box = document.getElementById('tray-dice');
  document.getElementById('tray-close').onclick = () => { tray.hidden = true; };

  const pipMap = {1:[4],2:[0,8],3:[0,4,8],4:[0,2,6,8],5:[0,2,4,6,8],6:[0,2,3,5,6,8]};

  function cube(value) {
    const el = document.createElement('div');
    el.className = 'die';
    const side = [value, value % 6 + 1, (value + 1) % 6 + 1, 7 - value, (value + 2) % 6 + 1, (value + 3) % 6 + 1];
    for (let f = 1; f <= 6; f++) {
      const face = document.createElement('div');
      face.className = 'face f' + f;
      for (let i = 0; i < 9; i++) {
        const c = document.createElement('div');
        if (pipMap[side[f - 1]].includes(i)) c.className = 'pip';
        face.appendChild(c);
      }
      el.appendChild(face);
    }
    const rx = Math.floor(Math.random() * 360), ry = Math.floor(Math.random() * 360);
    el.style.transition = 'none';
    el.style.transform = `rotateX(${rx}deg) rotateY(${ry}deg)`;
    requestAnimationFrame(() => requestAnimationFrame(() => {
      el.style.transition = '';
      const turns = 2 + Math.floor(Math.random() * 2);
      const tx = -12 - Math.random() * 12, ty = 12 + Math.random() * 16;
      el.style.transform = `rotateX(${tx + turns * 360}deg) rotateY(${ty + turns * 360}deg)`;
    }));
    return el;
  }

  const shapes = {
    4: '30,4 56,52 4,52', 8: '30,2 56,30 30,58 4,30', 12: '30,2 52,14 58,38 42,58 18,58 2,38 8,14',
    20: '30,2 54,16 54,44 30,58 6,44 6,16', 10: '30,2 56,26 44,58 16,58 4,26', 100: '30,2 56,26 44,58 16,58 4,26'
  };
  function poly(sides, value) {
    const el = document.createElement('div');
    el.className = 'poly spin';
    const pts = shapes[sides] || shapes[20];
    el.innerHTML = `<svg viewBox="0 0 60 60"><polygon points="${pts}" fill="url(#g${sides})" stroke="#6b4a1e" stroke-width="2"/><defs><linearGradient id="g${sides}" x1="0" y1="0" x2="1" y2="1"><stop offset="0" stop-color="#fdf6e3"/><stop offset="1" stop-color="#c9a24a"/></linearGradient></defs></svg><span class="num"></span>`;
    const num = el.querySelector('.num');
    let t = 0;
    const iv = setInterval(() => { num.textContent = 1 + Math.floor(Math.random() * sides); if (++t > 14) { clearInterval(iv); } }, 80);
    setTimeout(() => { clearInterval(iv); num.textContent = value; el.classList.remove('spin'); }, 1300);
    return el;
  }

  function show(rolls, title) {
    return new Promise(resolve => {
      if (!rolls || !rolls.length) { resolve(); return; }
      box.innerHTML = '';
      label.textContent = title || 'Кубики брошены';
      tray.hidden = false;
      for (const r of rolls) {
        const g = document.createElement('div');
        g.className = 'tray-group';
        const gl = document.createElement('div'); gl.className = 'gl'; gl.textContent = `${r.label} · ${r.dice.length}к${r.sides}${r.bonus ? '+' + r.bonus : ''}`;
        const gd = document.createElement('div'); gd.className = 'gd';
        r.dice.forEach((v, i) => {
          const el = r.sides === 6 ? cube(v) : poly(r.sides, v);
          if (r.dropped && r.dropped.includes(i)) setTimeout(() => el.classList.add('dropped'), 1500);
          gd.appendChild(el);
        });
        const gt = document.createElement('div'); gt.className = 'gt'; gt.textContent = '';
        setTimeout(() => { gt.textContent = '= ' + r.total; }, 1500);
        g.append(gl, gd, gt);
        box.appendChild(g);
      }
      setTimeout(resolve, 1700);
    });
  }
  return { show };
})();
