window.Music = (() => {
  const btn = document.getElementById('music-btn');
  const audio = new Audio('/static/music/tavern.mp3');
  audio.loop = true;
  audio.preload = 'none';
  audio.volume = 0;
  let on = false, fade;

  function ramp(to, ms, done) {
    clearInterval(fade);
    const from = audio.volume, start = Date.now();
    fade = setInterval(() => {
      const k = Math.min(1, (Date.now() - start) / ms);
      audio.volume = from + (to - from) * k;
      if (k >= 1) { clearInterval(fade); done && done(); }
    }, 50);
  }

  function setOn(v) {
    on = v;
    btn.setAttribute('aria-pressed', v ? 'true' : 'false');
    try { localStorage.setItem('music', v ? '1' : '0'); } catch (e) {}
    if (v) {
      audio.play().then(() => ramp(0.7, 1500)).catch(() => { on = false; btn.setAttribute('aria-pressed', 'false'); });
    } else {
      ramp(0, 600, () => audio.pause());
    }
  }
  btn.addEventListener('click', () => setOn(!on));
  let want = false;
  try { want = localStorage.getItem('music') === '1'; } catch (e) {}
  if (want) {
    btn.setAttribute('aria-pressed', 'true');
    const once = () => { setOn(true); document.removeEventListener('pointerdown', once); };
    document.addEventListener('pointerdown', once);
  }
  return { setOn };
})();
