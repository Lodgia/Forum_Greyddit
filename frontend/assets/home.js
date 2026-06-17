document.addEventListener('DOMContentLoaded', function(){
  // Fade in cards
  document.querySelectorAll('.card').forEach((el, i) => {
    setTimeout(() => el.classList.add('show'), 80 * i);
  });

  // Simple button press animation
  document.querySelectorAll('.btn').forEach(btn => {
    btn.addEventListener('click', (e) => {
      btn.style.transform = 'scale(0.98)';
      setTimeout(() => btn.style.transform = '', 120);
    });
  });
});
